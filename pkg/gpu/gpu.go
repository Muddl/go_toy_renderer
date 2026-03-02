//go:build !headless

package gpu

import (
	"errors"
	"fmt"
	gomath "math"

	"github.com/go-webgpu/webgpu/wgpu"
	"github.com/gogpu/gputypes"
	"github.com/muddl/go_toy_renderer/pkg/geometry"
)

// vertexStride is the byte stride for one packed vertex: 3×f32 position + 3×f32 color.
const vertexStride = 24 // 6 × float32 = 6 × 4 bytes

// makeCubeShaderWGSL generates the cube vertex+fragment WGSL shader with a
// hardcoded MVP matrix computed from camera pos=(3,2,5) looking at origin.
// The MVP is recomputed from the actual window dimensions to keep the correct
// aspect ratio (uniforms are introduced in Phase 14).
func makeCubeShaderWGSL(width, height uint32) string {
	mvp := computeHardcodedMVP(width, height)
	return fmt.Sprintf(`struct VertexInput {
    @location(0) position: vec3<f32>,
    @location(1) color: vec3<f32>,
}

struct VertexOutput {
    @builtin(position) clip_position: vec4<f32>,
    @location(0) color: vec3<f32>,
}

// MVP = Projection * View; camera at (3,2,5) looking at origin, fov=60 deg.
// Columns of the column-major mat4x4<f32>:
const mvp = mat4x4<f32>(
    vec4<f32>(%f, %f, %f, %f),
    vec4<f32>(%f, %f, %f, %f),
    vec4<f32>(%f, %f, %f, %f),
    vec4<f32>(%f, %f, %f, %f),
);

@vertex
fn vs_main(in: VertexInput) -> VertexOutput {
    var out: VertexOutput;
    out.clip_position = mvp * vec4<f32>(in.position, 1.0);
    out.color = in.color;
    return out;
}

@fragment
fn fs_main(in: VertexOutput) -> @location(0) vec4<f32> {
    return vec4<f32>(in.color, 1.0);
}
`,
		// column 0
		mvp[0], mvp[1], mvp[2], mvp[3],
		// column 1
		mvp[4], mvp[5], mvp[6], mvp[7],
		// column 2
		mvp[8], mvp[9], mvp[10], mvp[11],
		// column 3
		mvp[12], mvp[13], mvp[14], mvp[15],
	)
}

// computeHardcodedMVP computes MVP = Projection * View for a fixed camera.
// Camera: pos=(3,2,5), target=(0,0,0), up=(0,1,0), fov=60 deg.
// Projection uses WebGPU NDC (Z in [0,1]). Result is column-major [16]float64.
func computeHardcodedMVP(width, height uint32) [16]float64 {
	// Camera parameters.
	eyeX, eyeY, eyeZ := 3.0, 2.0, 5.0
	upX, upY, upZ := 0.0, 1.0, 0.0
	fov := gomath.Pi / 3.0 // 60 degrees
	aspect := float64(width) / float64(height)
	near, far := 0.1, 100.0

	// forward = normalize(target - eye); target is origin.
	fwdX, fwdY, fwdZ := -eyeX, -eyeY, -eyeZ
	fwdLen := gomath.Sqrt(fwdX*fwdX + fwdY*fwdY + fwdZ*fwdZ)
	fwdX, fwdY, fwdZ = fwdX/fwdLen, fwdY/fwdLen, fwdZ/fwdLen

	// right = normalize(forward × up)
	rx := fwdY*upZ - fwdZ*upY
	ry := fwdZ*upX - fwdX*upZ
	rz := fwdX*upY - fwdY*upX
	rLen := gomath.Sqrt(rx*rx + ry*ry + rz*rz)
	rx, ry, rz = rx/rLen, ry/rLen, rz/rLen

	// correctedUp = right × forward
	ux := ry*fwdZ - rz*fwdY
	uy := rz*fwdX - rx*fwdZ
	uz := rx*fwdY - ry*fwdX

	// View matrix translation components.
	rdotE := rx*eyeX + ry*eyeY + rz*eyeZ
	udotE := ux*eyeX + uy*eyeY + uz*eyeZ
	fdotE := fwdX*eyeX + fwdY*eyeY + fwdZ*eyeZ

	// View matrix (column-major):
	// Row i of the view matrix = [r, u, -fwd, 0] with translation column.
	view := [16]float64{
		rx, ux, -fwdX, 0, // col 0
		ry, uy, -fwdY, 0, // col 1
		rz, uz, -fwdZ, 0, // col 2
		-rdotE, -udotE, fdotE, 1, // col 3
	}

	// Projection matrix (WebGPU NDC Z in [0,1], column-major):
	// m[2][2] = far/(near-far), m[2][3] = near*far/(near-far), m[3][2] = -1
	tanHalf := gomath.Tan(fov / 2.0)
	proj := [16]float64{
		1.0 / (aspect * tanHalf), 0, 0, 0, // col 0
		0, 1.0 / tanHalf, 0, 0, // col 1
		0, 0, far / (near - far), -1, // col 2
		0, 0, near * far / (near - far), 0, // col 3
	}

	// MVP = proj × view (column-major multiply).
	// mvp[col*4+row] = Σ_k proj[k*4+row] * view[col*4+k]
	var mvp [16]float64
	for col := 0; col < 4; col++ {
		for row := 0; row < 4; row++ {
			var sum float64
			for k := 0; k < 4; k++ {
				sum += proj[k*4+row] * view[col*4+k]
			}
			mvp[col*4+row] = sum
		}
	}
	return mvp
}

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

	// Geometry cache (Phase 12).
	vertexBuf     *wgpu.Buffer
	vertexBufSize uint64
	indexBuf      *wgpu.Buffer
	indexBufSize  uint64
	indexCount    uint32
	cachedMesh    *geometry.Mesh

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

	// Step 9: Compile the cube WGSL shader (MVP hardcoded from camera at (3,2,5)).
	cubeWGSL := makeCubeShaderWGSL(width, height)
	shader := dev.CreateShaderModuleWGSL(cubeWGSL)
	if shader == nil {
		return errors.New("gpu: create shader module: returned nil")
	}
	d.shader = shader

	// Step 10: Create the render pipeline with vertex buffer layout and depth stencil.
	attrs := []wgpu.VertexAttribute{
		{Format: gputypes.VertexFormatFloat32x3, Offset: 0, ShaderLocation: 0},  // position
		{Format: gputypes.VertexFormatFloat32x3, Offset: 12, ShaderLocation: 1}, // color
	}
	pipeline := d.device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
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

// RenderFrame renders one frame of cube geometry to the GPU.
// LoadGeometry must be called before the first RenderFrame call to upload mesh data.
// Returns an error if Init has not been called successfully.
func (d *Device) RenderFrame() error {
	if !d.ready {
		return errors.New("gpu: Device not initialized — call Init first")
	}
	if d.vertexBuf == nil || d.indexBuf == nil {
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
	pass.SetVertexBuffer(0, d.vertexBuf, 0, d.vertexBufSize)
	pass.SetIndexBuffer(d.indexBuf, gputypes.IndexFormatUint32, 0, d.indexBufSize)
	pass.DrawIndexed(d.indexCount, 1, 0, 0, 0)
	// Optional overlay: renders additional draw calls into the same pass
	// (alpha-blended, no depth write) after the geometry draw.
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
	if d.vertexBuf != nil {
		d.vertexBuf.Destroy()
		d.vertexBuf.Release()
		d.vertexBuf = nil
	}
	if d.indexBuf != nil {
		d.indexBuf.Destroy()
		d.indexBuf.Release()
		d.indexBuf = nil
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
	d.cachedMesh = nil
	d.ready = false
}
