# GPU Backend: WebGPU / wgpu-native

**Status:** 📋 Planned — Phase 11 (WebGPU Integration)

---

## Overview

The GPU backend uses [wgpu-native](https://github.com/gfx-rs/wgpu-native) via Go bindings to provide cross-platform GPU rendering. wgpu-native implements the WebGPU standard and automatically selects the best native API:

| Platform | Native Backend | Notes |
|----------|---------------|-------|
| Windows  | D3D12         | Primary; falls back to Vulkan |
| macOS    | Metal         | Only option on Apple Silicon |
| Linux    | Vulkan        | Primary; falls back to OpenGL |

This means a single codebase runs natively on all major platforms with first-class GPU performance.

---

## Architecture

### Package Layout

```
pkg/gpu/
├── device.go        # GPU instance, adapter, device, queue
├── swapchain.go     # Surface creation, swap chain, resize handling
├── pipeline.go      # Render pipeline (vertex layout, shaders, depth-stencil)
├── buffer.go        # Vertex buffer, index buffer, uniform buffer management
├── shader.go        # Shader module loading (from compiled WGSL)
├── texture.go       # GPU textures, depth attachment, sampler (Phase 16)
└── gpu_test.go      # Integration tests (require GPU; skip if unavailable)
```

### Renderer Interface (Phase 10)

The GPU backend implements the `pkg/renderer.Renderer` interface:

```go
// pkg/renderer/renderer.go
type Renderer interface {
    Render(scene *render.Scene, cam camera.Camera) error
    Resize(width, height int)
    Destroy()
}
```

Factory:
```go
func New(backend Backend, width, height int) (Renderer, error) {
    switch backend {
    case CPU:  return newCPURenderer(width, height)
    case GPU:  return newGPURenderer(width, height)
    case Auto: // try GPU, fall back to CPU
    }
}
```

---

## Device Initialisation (Phase 11)

```
wgpu.CreateInstance()
    └── RequestAdapter(surface, PowerPreference: HighPerformance)
            └── RequestDevice(limits, features)
                    ├── device.Queue        ← command submission
                    └── device.CreateXxx()  ← resource creation
```

**Key points:**
- `RequestAdapter` is async (Go binding provides synchronous wrapper)
- Log adapter name, backend, and driver version at startup
- Register device-lost and uncaptured-error callbacks for debugging

---

## Surface & Swap Chain (Phase 11)

The surface is created from the GLFW window handle:

```go
// pkg/gpu/swapchain.go
type SwapChain struct {
    surface   *wgpu.Surface
    config    wgpu.SurfaceConfiguration
    depthView *wgpu.TextureView
}

func (sc *SwapChain) Resize(width, height int) {
    sc.config.Width  = uint32(width)
    sc.config.Height = uint32(height)
    sc.surface.Configure(device, &sc.config)
    sc.recreateDepth(device, width, height)
}
```

**Swap chain format:** `BGRA8Unorm` (preferred on most platforms)
**Present mode:** `Fifo` (vsync) by default; `Mailbox` for uncapped frame rate

---

## Render Pass (Phase 11+)

Each frame follows this pattern:

```go
// Acquire current frame texture
frame, err := swapChain.surface.GetCurrentTexture()
view := frame.texture.CreateView(nil)

// Begin command encoding
encoder := device.CreateCommandEncoder(nil)

// Render pass with color + depth attachments
pass := encoder.BeginRenderPass(&wgpu.RenderPassDescriptor{
    ColorAttachments: []wgpu.RenderPassColorAttachment{{
        View:       view,
        LoadOp:     wgpu.LoadOpClear,
        StoreOp:    wgpu.StoreOpStore,
        ClearValue: wgpu.Color{R: 0.1, G: 0.1, B: 0.1, A: 1.0},
    }},
    DepthStencilAttachment: &wgpu.RenderPassDepthStencilAttachment{
        View:            depthView,
        DepthLoadOp:     wgpu.LoadOpClear,
        DepthStoreOp:    wgpu.StoreOpStore,
        DepthClearValue: 1.0,
    },
})

// Draw calls here (Phase 12)
pass.End()

// Submit and present
device.Queue.Submit(encoder.Finish(nil))
frame.Present()
```

---

## Geometry Buffers (Phase 12)

### Vertex Layout

Interleaved format (24 bytes/vertex):

| Offset | Size | Attribute | WGSL Type |
|--------|------|-----------|-----------|
| 0      | 12   | Position  | `vec3f`   |
| 12     | 12   | Color     | `vec3f`   |

```go
// pkg/gpu/buffer.go
func UploadMesh(device *wgpu.Device, mesh *geometry.Mesh) (*GPUMesh, error) {
    // Interleave position + color into []float32
    // CreateBuffer with BufferUsageVertex | CopyDst
    // CreateBuffer with BufferUsageIndex | CopyDst
}
```

### Dirty Flag

`GPUScene` wraps `render.Scene` and tracks whether geometry has changed since the last upload. Vertex/index buffers are only re-uploaded when the mesh is modified.

---

## Depth Buffer (Phase 11)

- Format: `Depth32Float` (32-bit float depth, no stencil)
- Depth range: 0.0 (near) to 1.0 (far)
- Comparison: `Less` (closer fragment wins)
- Recreated on window resize alongside the swap chain

---

## Testing Strategy

**GPU tests must be skippable** — CI runners may not have a GPU:

```go
func TestGPUDevice_Init(t *testing.T) {
    if os.Getenv("GPU_TESTS") == "" {
        t.Skip("set GPU_TESTS=1 to run GPU integration tests")
    }
    // … test GPU initialisation
}
```

Run locally with: `GPU_TESTS=1 go test ./pkg/gpu/...`

**CPU renderer** continues to serve as the reference for visual correctness — GPU output is compared against CPU output with a tolerance (±2 pixel values per channel) to account for GPU floating-point rounding differences.

---

## Common Gotchas

- **Surface creation order:** The wgpu surface must be created before the device on some platforms (especially macOS Metal). Create the instance and surface from the GLFW window before requesting an adapter.
- **Depth texture lifetime:** The depth texture view must be recreated every time the window resizes. Rendering with a stale depth view causes validation errors.
- **Adapter selection:** On machines with integrated + discrete GPUs, `HighPerformance` preference selects the discrete GPU. Use `LowPower` for testing battery impact.
- **wgpu-native cgo:** Building requires the wgpu-native static library (`.a` / `.lib`) in the platform-specific lib path. Document the exact setup in README.

---

## References

- [wgpu-native C API](https://github.com/gfx-rs/wgpu-native)
- [go-webgpu Go bindings](https://github.com/rajveermalviya/go-webgpu)
- [WebGPU specification](https://www.w3.org/TR/webgpu/)
- [wgpu Examples (Rust)](https://github.com/gfx-rs/wgpu/tree/trunk/examples)
