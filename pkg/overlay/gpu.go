//go:build !headless

package overlay

import (
	"errors"

	"github.com/go-webgpu/webgpu/wgpu"
	"github.com/gogpu/gputypes"
	"github.com/muddl/go_toy_renderer/pkg/gpu"
)

// overlayWGSL is the WGSL shader for the fullscreen-quad overlay pass.
// The vertex shader generates a two-triangle quad from vertex_index (no VB
// needed). The fragment shader samples the overlay RGBA texture and outputs
// the result; the pipeline blend state handles SrcAlpha/OneMinusSrcAlpha.
const overlayWGSL = `
struct VertexOutput {
    @builtin(position) clip_pos: vec4<f32>,
    @location(0) uv: vec2<f32>,
}

@vertex
fn vs_main(@builtin(vertex_index) vi: u32) -> VertexOutput {
    // Two-triangle fullscreen quad (CCW winding, 6 vertices).
    // NDC: top-left (-1,+1), bottom-right (+1,-1).
    // UV:  top-left (0,0),   bottom-right (1,1).
    let pos = array<vec2<f32>, 6>(
        vec2<f32>(-1.0,  1.0),
        vec2<f32>(-1.0, -1.0),
        vec2<f32>( 1.0,  1.0),
        vec2<f32>( 1.0, -1.0),
        vec2<f32>( 1.0,  1.0),
        vec2<f32>(-1.0, -1.0),
    );
    let uv = array<vec2<f32>, 6>(
        vec2<f32>(0.0, 0.0),
        vec2<f32>(0.0, 1.0),
        vec2<f32>(1.0, 0.0),
        vec2<f32>(1.0, 1.0),
        vec2<f32>(1.0, 0.0),
        vec2<f32>(0.0, 1.0),
    );
    var out: VertexOutput;
    out.clip_pos = vec4<f32>(pos[vi], 0.0, 1.0);
    out.uv = uv[vi];
    return out;
}

@group(0) @binding(0) var overlay_tex: texture_2d<f32>;
@group(0) @binding(1) var overlay_sampler: sampler;

@fragment
fn fs_main(in: VertexOutput) -> @location(0) vec4<f32> {
    return textureSample(overlay_tex, overlay_sampler, in.uv);
}
`

// OverlayPass holds the GPU resources for compositing the debug overlay into
// the main render pass after the geometry draw call. It implements the
// gpu.OverlayRenderer interface so it can be registered via
// gpu.Device.SetOverlayRenderer.
//
// Lifecycle: NewOverlayPass → CreateOverlayPipeline → (per frame)
// UpdateOverlayTexture → (called by Device.RenderFrame) RenderIntoPass →
// Release.
type OverlayPass struct {
	gpuDevice *gpu.Device

	// Raw wgpu handles — populated by CreateOverlayPipeline.
	wgpuDevice     *wgpu.Device
	wgpuQueue      *wgpu.Queue
	format         gputypes.TextureFormat
	fbWidth        uint32
	fbHeight       uint32
	shader         *wgpu.ShaderModule
	bindLayout     *wgpu.BindGroupLayout
	pipelineLayout *wgpu.PipelineLayout
	pipeline       *wgpu.RenderPipeline
	sampler        *wgpu.Sampler
	texture        *wgpu.Texture
	textureView    *wgpu.TextureView
	bindGroup      *wgpu.BindGroup
}

// NewOverlayPass creates an OverlayPass wrapping the given gpu.Device.
// Call CreateOverlayPipeline after gpu.Device.Init to allocate GPU resources.
func NewOverlayPass(d *gpu.Device) *OverlayPass {
	return &OverlayPass{gpuDevice: d}
}

// CreateOverlayPipeline allocates all GPU resources for the overlay:
// shader, sampler, bind group layout, pipeline, framebuffer-sized RGBA texture.
// Returns an error if gpu.Device is not initialized.
func (op *OverlayPass) CreateOverlayPipeline() error {
	if !op.gpuDevice.IsReady() {
		return errors.New("overlay: CreateOverlayPipeline: gpu device not initialized — call Device.Init first")
	}

	op.wgpuDevice = op.gpuDevice.DeviceHandle()
	op.wgpuQueue = op.gpuDevice.QueueHandle()
	op.format = op.gpuDevice.SurfaceFormat()
	op.fbWidth, op.fbHeight = op.gpuDevice.Dimensions()

	// Compile the overlay WGSL shader.
	shader := op.wgpuDevice.CreateShaderModuleWGSL(overlayWGSL)
	if shader == nil {
		return errors.New("overlay: CreateOverlayPipeline: create shader module returned nil")
	}
	op.shader = shader

	// Nearest-filter sampler — text is rendered 1:1 on-screen, no interpolation needed.
	sampler := op.wgpuDevice.CreateNearestSampler()
	if sampler == nil {
		return errors.New("overlay: CreateOverlayPipeline: create sampler returned nil")
	}
	op.sampler = sampler

	// Bind group layout: binding 0 = texture, binding 1 = sampler.
	bindLayout := op.wgpuDevice.CreateBindGroupLayoutSimple([]wgpu.BindGroupLayoutEntry{
		{
			Binding:    0,
			Visibility: gputypes.ShaderStageFragment,
			Texture: wgpu.TextureBindingLayout{
				SampleType:    gputypes.TextureSampleTypeFloat,
				ViewDimension: gputypes.TextureViewDimension2D,
				Multisampled:  wgpu.False,
			},
		},
		{
			Binding:    1,
			Visibility: gputypes.ShaderStageFragment,
			Sampler: wgpu.SamplerBindingLayout{
				Type: gputypes.SamplerBindingTypeFiltering,
			},
		},
	})
	if bindLayout == nil {
		return errors.New("overlay: CreateOverlayPipeline: create bind group layout returned nil")
	}
	op.bindLayout = bindLayout

	// Pipeline layout wraps the bind group layout. Must be passed explicitly to
	// CreateRenderPipeline — omitting it causes wgpu to use auto-layout, which
	// creates exclusive bind group layouts incompatible with our bindGroup.
	pipelineLayout := op.wgpuDevice.CreatePipelineLayoutSimple([]*wgpu.BindGroupLayout{bindLayout})
	if pipelineLayout == nil {
		return errors.New("overlay: CreateOverlayPipeline: create pipeline layout returned nil")
	}
	op.pipelineLayout = pipelineLayout

	// Blend state: SrcAlpha / OneMinusSrcAlpha (standard alpha compositing).
	blend := &wgpu.BlendState{
		Color: wgpu.BlendComponent{
			Operation: gputypes.BlendOperationAdd,
			SrcFactor: gputypes.BlendFactorSrcAlpha,
			DstFactor: gputypes.BlendFactorOneMinusSrcAlpha,
		},
		Alpha: wgpu.BlendComponent{
			Operation: gputypes.BlendOperationAdd,
			SrcFactor: gputypes.BlendFactorOne,
			DstFactor: gputypes.BlendFactorZero,
		},
	}

	// Render pipeline: no vertex buffers (vertices generated in shader),
	// no depth write (overlay always draws on top), alpha blend.
	pipeline := op.wgpuDevice.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
		Layout: pipelineLayout,
		Vertex: wgpu.VertexState{
			Module:     shader,
			EntryPoint: "vs_main",
		},
		Fragment: &wgpu.FragmentState{
			Module:     shader,
			EntryPoint: "fs_main",
			Targets: []wgpu.ColorTargetState{
				{
					Format:    op.format,
					Blend:     blend,
					WriteMask: gputypes.ColorWriteMaskAll,
				},
			},
		},
		Primitive: wgpu.PrimitiveState{
			Topology: gputypes.PrimitiveTopologyTriangleList,
			CullMode: gputypes.CullModeNone,
		},
		// DepthStencil: depth test always passes, no depth write.
		// Required to be compatible with the main pass's depth attachment.
		DepthStencil: &wgpu.DepthStencilState{
			Format:            gputypes.TextureFormatDepth24Plus,
			DepthWriteEnabled: false,
			DepthCompare:      gputypes.CompareFunctionAlways,
		},
		Multisample: wgpu.MultisampleState{
			Count: 1,
			Mask:  0xFFFFFFFF,
		},
	})
	if pipeline == nil {
		return errors.New("overlay: CreateOverlayPipeline: create render pipeline returned nil")
	}
	op.pipeline = pipeline

	// Overlay texture: framebuffer-sized RGBA8Unorm, initially all zeros
	// (transparent). UpdateOverlayTexture writes the TextLayer region each frame.
	overlayTex := op.wgpuDevice.CreateTexture(&wgpu.TextureDescriptor{
		Label:     wgpu.EmptyStringView(),
		Usage:     gputypes.TextureUsageTextureBinding | gputypes.TextureUsageCopyDst,
		Dimension: gputypes.TextureDimension2D,
		Size: gputypes.Extent3D{
			Width:              op.fbWidth,
			Height:             op.fbHeight,
			DepthOrArrayLayers: 1,
		},
		Format:        gputypes.TextureFormatRGBA8Unorm,
		MipLevelCount: 1,
		SampleCount:   1,
	})
	if overlayTex == nil {
		return errors.New("overlay: CreateOverlayPipeline: create overlay texture returned nil")
	}
	op.texture = overlayTex

	textureView := overlayTex.CreateView(nil)
	if textureView == nil {
		return errors.New("overlay: CreateOverlayPipeline: create texture view returned nil")
	}
	op.textureView = textureView

	// Bind group: bind the texture view and sampler.
	bindGroup := op.wgpuDevice.CreateBindGroupSimple(bindLayout, []wgpu.BindGroupEntry{
		wgpu.TextureBindingEntry(0, textureView),
		wgpu.SamplerBindingEntry(1, sampler),
	})
	if bindGroup == nil {
		return errors.New("overlay: CreateOverlayPipeline: create bind group returned nil")
	}
	op.bindGroup = bindGroup

	return nil
}

// UpdateOverlayTexture uploads the TextLayer's RGBA pixel buffer to the GPU
// overlay texture at position (0, 0). The layer must not be nil. Call this
// each frame before Device.RenderFrame so the texture is current.
func (op *OverlayPass) UpdateOverlayTexture(layer *TextLayer) error {
	if layer == nil {
		return errors.New("overlay: UpdateOverlayTexture: layer is nil")
	}
	if op.texture == nil {
		return errors.New("overlay: UpdateOverlayTexture: pipeline not created — call CreateOverlayPipeline first")
	}

	w := uint32(layer.Width())
	h := uint32(layer.Height())
	pixels := layer.Pixels()

	op.wgpuQueue.WriteTexture(
		&wgpu.TexelCopyTextureInfo{
			Texture:  op.texture.Handle(),
			MipLevel: 0,
			Origin:   gputypes.Origin3D{X: 0, Y: 0, Z: 0},
			Aspect:   wgpu.TextureAspectAll,
		},
		pixels,
		&wgpu.TexelCopyBufferLayout{
			Offset:       0,
			BytesPerRow:  w * 4, // 4 bytes per RGBA texel
			RowsPerImage: h,
		},
		&gputypes.Extent3D{
			Width:              w,
			Height:             h,
			DepthOrArrayLayers: 1,
		},
	)
	return nil
}

// RenderIntoPass records the overlay draw into an active render pass encoder.
// It sets the overlay pipeline, binds the texture/sampler bind group, and
// draws a fullscreen quad (6 vertices, no vertex buffer). Implements
// gpu.OverlayRenderer. Returns nil if the pipeline is not yet created.
func (op *OverlayPass) RenderIntoPass(pass *wgpu.RenderPassEncoder) error {
	if op.pipeline == nil || op.bindGroup == nil {
		return nil // overlay not ready; skip silently
	}
	pass.SetPipeline(op.pipeline)
	pass.SetBindGroup(0, op.bindGroup, nil)
	pass.Draw(6, 1, 0, 0) // fullscreen quad: 6 vertices, 1 instance
	return nil
}

// Release frees all GPU resources held by the OverlayPass. Safe to call on an
// uninitialized OverlayPass and to call multiple times.
func (op *OverlayPass) Release() {
	if op.bindGroup != nil {
		op.bindGroup.Release()
		op.bindGroup = nil
	}
	if op.textureView != nil {
		op.textureView.Release()
		op.textureView = nil
	}
	if op.texture != nil {
		op.texture.Destroy()
		op.texture.Release()
		op.texture = nil
	}
	if op.sampler != nil {
		op.sampler.Release()
		op.sampler = nil
	}
	if op.bindLayout != nil {
		op.bindLayout.Release()
		op.bindLayout = nil
	}
	if op.pipeline != nil {
		op.pipeline.Release()
		op.pipeline = nil
	}
	if op.pipelineLayout != nil {
		op.pipelineLayout.Release()
		op.pipelineLayout = nil
	}
	if op.shader != nil {
		op.shader.Release()
		op.shader = nil
	}
	op.wgpuDevice = nil
	op.wgpuQueue = nil
}

// verifyOverlayRenderer ensures *OverlayPass satisfies gpu.OverlayRenderer
// at compile time. The blank assignment is discarded by the compiler.
var _ interface {
	RenderIntoPass(pass *wgpu.RenderPassEncoder) error
} = (*OverlayPass)(nil)
