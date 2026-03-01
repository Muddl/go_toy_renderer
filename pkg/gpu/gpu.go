//go:build !headless

package gpu

import (
	"errors"
	"fmt"

	"github.com/go-webgpu/webgpu/wgpu"
	"github.com/gogpu/gputypes"
)

// helloTriangleWGSL is an inline WGSL shader for the Hello Triangle.
// Vertex: positions are computed from vertex_index (no vertex buffer needed).
// Fragment: solid orange color.
const helloTriangleWGSL = `
@vertex
fn vs_main(@builtin(vertex_index) in_vertex_index: u32) -> @builtin(position) vec4<f32> {
    var pos = array<vec2<f32>, 3>(
        vec2<f32>( 0.0,  0.5),
        vec2<f32>(-0.5, -0.5),
        vec2<f32>( 0.5, -0.5),
    );
    return vec4<f32>(pos[in_vertex_index], 0.0, 1.0);
}

@fragment
fn fs_main() -> @location(0) vec4<f32> {
    return vec4<f32>(1.0, 0.5, 0.2, 1.0);
}
`

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

	// Step 4: Request a GPU adapter (physical device selection).
	adapter, err := inst.RequestAdapter(nil)
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

	// Step 8: Compile the Hello Triangle WGSL shader.
	shader := dev.CreateShaderModuleWGSL(helloTriangleWGSL)
	if shader == nil {
		return errors.New("gpu: create shader module: returned nil")
	}
	d.shader = shader

	// Step 9: Create the render pipeline (no vertex buffers needed).
	pipeline := dev.CreateRenderPipelineSimple(
		nil, // auto layout
		shader, "vs_main",
		shader, "fs_main",
		d.format,
	)
	if pipeline == nil {
		return errors.New("gpu: create render pipeline: returned nil")
	}
	d.pipeline = pipeline

	d.ready = true
	return nil
}

// RenderFrame submits one Hello Triangle frame to the GPU.
// Returns an error if Init has not been called successfully.
func (d *Device) RenderFrame() error {
	if !d.ready {
		return errors.New("gpu: Device not initialized — call Init first")
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
	})
	if pass == nil {
		enc.Release()
		return errors.New("gpu: begin render pass: returned nil")
	}

	pass.SetPipeline(d.pipeline)
	pass.Draw(3, 1, 0, 0) // 3 vertices, 1 instance
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
