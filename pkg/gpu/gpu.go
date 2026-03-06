//go:build !headless

package gpu

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/go-webgpu/webgpu/wgpu"
	"github.com/gogpu/gputypes"
	"github.com/muddl/go_toy_renderer/assets/shaders"
	"github.com/muddl/go_toy_renderer/pkg/geometry"
	"github.com/muddl/go_toy_renderer/pkg/math"
)

// vertexStride is the byte stride for one packed vertex: 3×f32 position + 3×f32 color + 3×f32 normal.
const vertexStride = 36 // 9 × float32 = 9 × 4 bytes

// uniformBufferSize is the byte size of a single mat4x4<f32> uniform (64 bytes).
const uniformBufferSize = 64

// Device holds the wgpu initialization chain and render resources.
// Create with [New] and initialize with [Device.Init] before calling
// [Device.RenderFrame].
type Device struct {
	instance *wgpu.Instance
	adapter  *wgpu.Adapter
	device   *wgpu.Device
	queue    *wgpu.Queue
	surface  *wgpu.Surface
	shader   *wgpu.ShaderModule
	pipeline *wgpu.RenderPipeline
	format   gputypes.TextureFormat
	width    uint32
	height   uint32
	ready    bool

	// Depth buffer resources (Phase 12).
	depthTexture *wgpu.Texture
	depthView    *wgpu.TextureView

	// Geometry cache: maps mesh pointer → GPU buffers.
	// Multiple meshes stay resident so buffers remain valid during a render pass.
	meshCache map[*geometry.Mesh]*meshBufferEntry

	// Active geometry for the current draw call (set by LoadGeometry).
	vertexBuf     *wgpu.Buffer
	vertexBufSize uint64
	indexBuf      *wgpu.Buffer
	indexBufSize  uint64
	indexCount    uint32

	// Uniform buffers (Phase 14).
	cameraUniformBuf *wgpu.Buffer
	meshUniformBuf   *wgpu.Buffer
	bindGroupLayout  *wgpu.BindGroupLayout
	cameraBindGroup  *wgpu.BindGroup
	meshBindGroup    *wgpu.BindGroup
	pipelineLayout   *wgpu.PipelineLayout

	// Optional overlay renderer (Phase 13 / perf-debug-overlay).
	// Set via SetOverlayRenderer; called after the geometry draw in each frame.
	overlayRenderer OverlayRenderer
}

// New returns a new, uninitialised Device.
func New() *Device { return &Device{} }

// IsReady returns true if Init completed successfully and the device is ready to render.
func (d *Device) IsReady() bool { return d.ready }

// HasQueue returns true if the wgpu queue was acquired during Init.
// Used in integration tests to verify the init chain ran up to queue acquisition.
func (d *Device) HasQueue() bool { return d.queue != nil }

// Init initializes wgpu and prepares the Device for rendering.
// handle contains the platform-specific native window pointers (populated by
// pkg/renderer using GLFW native handle functions).
// width and height are the initial swap-chain dimensions in pixels.
func (d *Device) Init(width, height uint32, handle NativeWindowHandle) error {
	d.width = width
	d.height = height

	// Step 1: Load the wgpu-native shared library.
	// Set WGPU_NATIVE_PATH to the platform library path (e.g.
	// assets/windows-x86_64-gnu/lib/wgpu_native.dll).
	if err := wgpu.Init(); err != nil {
		return fmt.Errorf("gpu: wgpu library: %w", err)
	}

	// Step 2: Create the WebGPU instance.
	inst, err := wgpu.CreateInstance(nil)
	if err != nil {
		return fmt.Errorf("gpu: create instance: %w", err)
	}
	d.instance = inst

	// Step 3: Create the platform-specific surface from the native handle.
	surface, err := createPlatformSurface(inst, handle)
	if err != nil {
		return fmt.Errorf("gpu: create surface: %w", err)
	}
	d.surface = surface

	// Step 4: Request a GPU adapter that is compatible with our surface.
	// CompatibleSurface ensures the selected adapter's queue family supports
	// presenting to the surface, preventing the "queue family" validation error.
	adapter, err := inst.RequestAdapter(&wgpu.RequestAdapterOptions{
		CompatibleSurface: surface.Handle(),
	})
	if err != nil {
		return fmt.Errorf("gpu: request adapter: %w", err)
	}
	d.adapter = adapter

	// Step 5: Request a logical device from the adapter.
	dev, err := adapter.RequestDevice(nil)
	if err != nil {
		return fmt.Errorf("gpu: request device: %w", err)
	}
	d.device = dev

	// Step 6: Obtain the default command queue.
	q := dev.GetQueue()
	if q == nil {
		return fmt.Errorf("gpu: get queue: returned nil")
	}
	d.queue = q

	// Step 7: Choose the surface format and configure the swap-chain.
	d.format = gputypes.TextureFormatBGRA8Unorm
	surface.Configure(&wgpu.SurfaceConfiguration{
		Device:      d.device,
		Format:      d.format,
		Usage:       gputypes.TextureUsageRenderAttachment,
		Width:       width,
		Height:      height,
		PresentMode: gputypes.PresentModeFifo,
		AlphaMode:   gputypes.CompositeAlphaModeAuto,
	})

	// Step 8: Create the depth texture (Depth24Plus).
	depthTex := d.device.CreateDepthTexture(width, height, gputypes.TextureFormatDepth24Plus)
	if depthTex == nil {
		return errors.New("gpu: create depth texture: returned nil")
	}
	d.depthTexture = depthTex
	depthView := depthTex.CreateView(nil)
	if depthView == nil {
		return errors.New("gpu: create depth texture view: returned nil")
	}
	d.depthView = depthView

	// Step 9: Create uniform buffers for camera (viewProj) and mesh (model).
	d.cameraUniformBuf = dev.CreateBuffer(&wgpu.BufferDescriptor{
		Size:  uniformBufferSize,
		Usage: gputypes.BufferUsageUniform | gputypes.BufferUsageCopyDst,
	})
	if d.cameraUniformBuf == nil {
		return errors.New("gpu: create camera uniform buffer: returned nil")
	}
	d.meshUniformBuf = dev.CreateBuffer(&wgpu.BufferDescriptor{
		Size:  uniformBufferSize,
		Usage: gputypes.BufferUsageUniform | gputypes.BufferUsageCopyDst,
	})
	if d.meshUniformBuf == nil {
		return errors.New("gpu: create mesh uniform buffer: returned nil")
	}

	// Step 9b: Create bind group layout with two uniform bindings.
	bglEntries := []wgpu.BindGroupLayoutEntry{
		{
			Binding:    0,
			Visibility: gputypes.ShaderStageVertex,
			Buffer: wgpu.BufferBindingLayout{
				Type:           gputypes.BufferBindingTypeUniform,
				MinBindingSize: uniformBufferSize,
			},
		},
		{
			Binding:    1,
			Visibility: gputypes.ShaderStageVertex,
			Buffer: wgpu.BufferBindingLayout{
				Type:           gputypes.BufferBindingTypeUniform,
				MinBindingSize: uniformBufferSize,
			},
		},
	}
	bgl := dev.CreateBindGroupLayout(&wgpu.BindGroupLayoutDescriptor{
		EntryCount: uintptr(len(bglEntries)),
		Entries:    uintptr(unsafe.Pointer(&bglEntries[0])),
	})
	if bgl == nil {
		return errors.New("gpu: create bind group layout: returned nil")
	}
	d.bindGroupLayout = bgl

	// Step 9c: Create bind groups using the helper that handles Handle() conversion.
	cameraBGEntries := []wgpu.BindGroupEntry{
		wgpu.BufferBindingEntry(0, d.cameraUniformBuf, 0, uniformBufferSize),
		wgpu.BufferBindingEntry(1, d.meshUniformBuf, 0, uniformBufferSize),
	}
	cameraBindGroup := dev.CreateBindGroup(&wgpu.BindGroupDescriptor{
		Layout:     bgl.Handle(),
		EntryCount: uintptr(len(cameraBGEntries)),
		Entries:    uintptr(unsafe.Pointer(&cameraBGEntries[0])),
	})
	if cameraBindGroup == nil {
		return errors.New("gpu: create bind group: returned nil")
	}
	d.cameraBindGroup = cameraBindGroup

	// Step 9d: Create pipeline layout.
	bglHandles := [1]uintptr{bgl.Handle()}
	plLayout := dev.CreatePipelineLayout(&wgpu.PipelineLayoutDescriptor{
		BindGroupLayoutCount: 1,
		BindGroupLayouts:     uintptr(unsafe.Pointer(&bglHandles[0])),
	})
	if plLayout == nil {
		return errors.New("gpu: create pipeline layout: returned nil")
	}
	d.pipelineLayout = plLayout

	// Step 10: Compile the cube WGSL shader with uniform bindings.
	shader := dev.CreateShaderModuleWGSL(shaders.CubeWGSL)
	if shader == nil {
		return errors.New("gpu: create shader module: returned nil")
	}
	d.shader = shader

	// Step 11: Create the render pipeline with vertex buffer layout and depth stencil.
	attrs := []wgpu.VertexAttribute{
		{Format: gputypes.VertexFormatFloat32x3, Offset: 0, ShaderLocation: 0},  // position
		{Format: gputypes.VertexFormatFloat32x3, Offset: 12, ShaderLocation: 1}, // color
		{Format: gputypes.VertexFormatFloat32x3, Offset: 24, ShaderLocation: 2}, // normal
	}
	pipeline := d.device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
		Layout: plLayout,
		Vertex: wgpu.VertexState{
			Module:     shader,
			EntryPoint: "vs_main",
			Buffers: []wgpu.VertexBufferLayout{
				{
					ArrayStride:    vertexStride,
					StepMode:       gputypes.VertexStepModeVertex,
					AttributeCount: uintptr(len(attrs)),
					Attributes:     &attrs[0],
				},
			},
		},
		Fragment: &wgpu.FragmentState{
			Module:     shader,
			EntryPoint: "fs_main",
			Targets: []wgpu.ColorTargetState{
				{Format: d.format, WriteMask: gputypes.ColorWriteMaskAll},
			},
		},
		DepthStencil: &wgpu.DepthStencilState{
			Format:            gputypes.TextureFormatDepth24Plus,
			DepthWriteEnabled: true,
			DepthCompare:      gputypes.CompareFunctionLess,
		},
		Primitive: wgpu.PrimitiveState{
			Topology:  gputypes.PrimitiveTopologyTriangleList,
			FrontFace: gputypes.FrontFaceCCW,
			CullMode:  gputypes.CullModeBack,
		},
		Multisample: wgpu.MultisampleState{
			Count: 1,
			Mask:  0xFFFFFFFF,
		},
	})
	if pipeline == nil {
		return errors.New("gpu: create render pipeline: returned nil")
	}
	d.pipeline = pipeline

	d.ready = true
	return nil
}

// UpdateCameraUniforms uploads the view-projection matrix to the camera uniform buffer.
func (d *Device) UpdateCameraUniforms(viewProj math.Mat4x4) {
	if d.queue == nil || d.cameraUniformBuf == nil {
		return
	}
	cu := CameraUniforms{ViewProj: viewProj}
	d.queue.WriteBuffer(d.cameraUniformBuf, 0, cu.Bytes())
}

// UpdateMeshUniforms uploads the model matrix to the mesh uniform buffer.
func (d *Device) UpdateMeshUniforms(model math.Mat4x4) {
	if d.queue == nil || d.meshUniformBuf == nil {
		return
	}
	mu := MeshUniforms{Model: model}
	d.queue.WriteBuffer(d.meshUniformBuf, 0, mu.Bytes())
}

// RenderFrame renders one frame of cube geometry to the GPU.
// LoadGeometry must be called before the first RenderFrame call to upload mesh data.
// Camera and mesh uniforms should be uploaded via UpdateCameraUniforms/UpdateMeshUniforms
// before calling RenderFrame.
// Returns an error if Init has not been called successfully.
func (d *Device) RenderFrame() error {
	return d.RenderFrameMulti(nil)
}

// RenderFrameMulti renders one frame with multiple draw nodes.
// If nodes is nil or empty, falls back to the single cached geometry (Phase 12 compat).
// Camera uniform should be uploaded via UpdateCameraUniforms before calling.
// Each node's model matrix is uploaded to the mesh uniform buffer before its draw call.
func (d *Device) RenderFrameMulti(nodes []DrawNode) error {
	if !d.ready {
		return errors.New("gpu: Device not initialized — call Init first")
	}

	useLegacy := len(nodes) == 0
	if useLegacy && (d.vertexBuf == nil || d.indexBuf == nil) {
		return errors.New("gpu: no geometry loaded — call LoadGeometry before RenderFrame")
	}

	// Acquire the current swap-chain texture.
	surfTex, err := d.surface.GetCurrentTexture()
	if err != nil {
		return fmt.Errorf("gpu: get surface texture: %w", err)
	}
	view := surfTex.Texture.CreateView(nil)
	if view == nil {
		surfTex.Texture.Release()
		return errors.New("gpu: create texture view: returned nil")
	}
	defer view.Release()
	defer surfTex.Texture.Release()

	// Record GPU commands.
	enc := d.device.CreateCommandEncoder(nil)
	if enc == nil {
		return errors.New("gpu: create command encoder: returned nil")
	}

	pass := enc.BeginRenderPass(&wgpu.RenderPassDescriptor{
		ColorAttachments: []wgpu.RenderPassColorAttachment{
			{
				View:    view,
				LoadOp:  gputypes.LoadOpClear,
				StoreOp: gputypes.StoreOpStore,
				ClearValue: wgpu.Color{
					R: 0.1, G: 0.1, B: 0.1, A: 1.0, // dark grey background
				},
			},
		},
		DepthStencilAttachment: &wgpu.RenderPassDepthStencilAttachment{
			View:            d.depthView,
			DepthLoadOp:     gputypes.LoadOpClear,
			DepthStoreOp:    gputypes.StoreOpStore,
			DepthClearValue: 1.0, // far plane
		},
	})
	if pass == nil {
		enc.Release()
		return errors.New("gpu: begin render pass: returned nil")
	}

	pass.SetPipeline(d.pipeline)
	pass.SetBindGroup(0, d.cameraBindGroup, nil)

	if useLegacy {
		// Single-mesh legacy path (Phase 12 compat).
		pass.SetVertexBuffer(0, d.vertexBuf, 0, d.vertexBufSize)
		pass.SetIndexBuffer(d.indexBuf, gputypes.IndexFormatUint32, 0, d.indexBufSize)
		pass.DrawIndexed(d.indexCount, 1, 0, 0, 0)
	} else {
		// Multi-mesh path: load geometry + update model uniform per node.
		for i := range nodes {
			if err := d.LoadGeometry(nodes[i].Mesh); err != nil {
				continue
			}
			d.UpdateMeshUniforms(nodes[i].Model)
			pass.SetVertexBuffer(0, d.vertexBuf, 0, d.vertexBufSize)
			pass.SetIndexBuffer(d.indexBuf, gputypes.IndexFormatUint32, 0, d.indexBufSize)
			pass.DrawIndexed(d.indexCount, 1, 0, 0, 0)
		}
	}

	// Optional overlay.
	if d.overlayRenderer != nil {
		_ = d.overlayRenderer.RenderIntoPass(pass)
	}
	pass.End()
	pass.Release()

	cmd := enc.Finish(nil)
	enc.Release()
	if cmd == nil {
		return errors.New("gpu: finish command encoder: returned nil")
	}

	d.queue.Submit(cmd)
	cmd.Release()

	d.surface.Present()
	return nil
}

// Shutdown releases all wgpu resources held by the Device.
// Safe to call on an uninitialised Device and to call multiple times.
func (d *Device) Shutdown() {
	// Release all cached mesh GPU buffers.
	for _, entry := range d.meshCache {
		if entry.vertexBuf != nil {
			entry.vertexBuf.Destroy()
			entry.vertexBuf.Release()
		}
		if entry.indexBuf != nil {
			entry.indexBuf.Destroy()
			entry.indexBuf.Release()
		}
	}
	d.meshCache = nil
	d.vertexBuf = nil
	d.indexBuf = nil
	if d.cameraBindGroup != nil {
		d.cameraBindGroup.Release()
		d.cameraBindGroup = nil
	}
	if d.cameraUniformBuf != nil {
		d.cameraUniformBuf.Destroy()
		d.cameraUniformBuf.Release()
		d.cameraUniformBuf = nil
	}
	if d.meshUniformBuf != nil {
		d.meshUniformBuf.Destroy()
		d.meshUniformBuf.Release()
		d.meshUniformBuf = nil
	}
	if d.bindGroupLayout != nil {
		d.bindGroupLayout.Release()
		d.bindGroupLayout = nil
	}
	if d.pipelineLayout != nil {
		d.pipelineLayout.Release()
		d.pipelineLayout = nil
	}
	if d.depthView != nil {
		d.depthView.Release()
		d.depthView = nil
	}
	if d.depthTexture != nil {
		d.depthTexture.Destroy()
		d.depthTexture.Release()
		d.depthTexture = nil
	}
	if d.pipeline != nil {
		d.pipeline.Release()
		d.pipeline = nil
	}
	if d.shader != nil {
		d.shader.Release()
		d.shader = nil
	}
	if d.surface != nil {
		d.surface.Unconfigure()
		d.surface.Release()
		d.surface = nil
	}
	if d.queue != nil {
		d.queue.Release()
		d.queue = nil
	}
	if d.device != nil {
		d.device.Release()
		d.device = nil
	}
	if d.adapter != nil {
		d.adapter.Release()
		d.adapter = nil
	}
	if d.instance != nil {
		d.instance.Release()
		d.instance = nil
	}
	d.ready = false
}
