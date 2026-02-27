# go_toy_renderer

A toy 3D renderer written in Go. The MVP is a complete CPU software renderer that converts 3D geometry to PNG images. The project is now advancing toward **cross-platform GPU acceleration** with real-time windowed rendering and HLSL shaders via WebGPU (wgpu-native).

![CI](https://github.com/muddl/go_toy_renderer/actions/workflows/ci.yml/badge.svg)

---

## What it does (MVP — complete)

- Transforms 3D geometry through world → camera → clip → NDC → screen space
- Rasterizes triangles using barycentric coordinates with correct depth ordering
- Interpolates vertex colors across triangle surfaces
- Outputs a PNG image (640×480)

## Where it's going (GPU roadmap)

- **Phase 9:** Real-time windowed display at 60 fps via GLFW
- **Phase 10–11:** WebGPU backend (wgpu-native) — D3D12 on Windows, Metal on macOS, Vulkan on Linux
- **Phase 12–13:** GPU hardware rasterization with HLSL shaders compiled to WGSL
- **Phase 14–15:** Per-mesh transforms, Phong/PBR lighting, vertex normals
- **Phase 16:** Texture mapping, OBJ loading, post-processing (tone mapping, FXAA)

---

## Quick start (CPU renderer)

```bash
# Render a 640×480 PNG of a colored cube
go run ./cmd/renderer
# Output: output.png
```

```bash
# Or build first, then run
go build -o renderer ./cmd/renderer
./renderer
```

## Project structure

```
cmd/
  renderer/       # PNG output — renders output.png (CPU, complete)
  renderer-rt/    # Real-time window — Phase 9+ (planned)

pkg/
  math/           # Vec3, Mat4x4 — vectors, matrices, transformations
  geometry/       # Vertex, Mesh, and hardcoded primitives (Cube, Tetrahedron)
  camera/         # Perspective camera with LookAt view matrix
  framebuffer/    # Color + depth buffers; PNG export
  rasterize/      # Barycentric triangle rasterizer with attribute interpolation
  shader/         # Per-pixel shader system (VertexColor, FlatColor, Depth)
  render/         # Pipeline orchestration — Scene, Render()
  renderer/       # Renderer interface + CPU/GPU factory (Phase 10, planned)
  gpu/            # WebGPU device, swapchain, buffers, pipelines (Phase 11+, planned)
  window/         # GLFW window, event loop, camera controller (Phase 9, planned)

assets/
  shaders/        # HLSL source + compiled WGSL (Phase 13+, planned)
```

---

## Using the CPU API

```go
package main

import (
    "log"
    "github.com/muddl/go_toy_renderer/pkg/camera"
    "github.com/muddl/go_toy_renderer/pkg/framebuffer"
    "github.com/muddl/go_toy_renderer/pkg/geometry"
    math "github.com/muddl/go_toy_renderer/pkg/math"
    "github.com/muddl/go_toy_renderer/pkg/render"
    "github.com/muddl/go_toy_renderer/pkg/shader"
)

func main() {
    // 1. Build a scene
    scene := render.NewScene()
    scene.AddMesh(geometry.NewCube())

    // 2. Position the camera
    cam := camera.New(
        math.Vec3{X: 3, Y: 2, Z: 5}, // position
        math.Vec3{X: 0, Y: 0, Z: 0}, // target
        math.Vec3{X: 0, Y: 1, Z: 0}, // up vector
        45.0,                         // FOV (degrees)
        640.0/480.0,                  // aspect ratio
        0.1,                          // near plane
        100.0,                        // far plane
    )

    // 3. Allocate a framebuffer and render
    fb := framebuffer.New(640, 480)
    render.Render(scene, cam, fb, shader.VertexColor)

    // 4. Save to disk
    if err := fb.SavePNG("output.png"); err != nil {
        log.Fatal(err)
    }
}
```

### Built-in shaders (CPU)

| Shader | Description |
|---|---|
| `shader.VertexColor` | Interpolates vertex colors across each triangle |
| `shader.NewFlatColor(color)` | Fills every pixel with a constant color |
| `shader.Depth` | Encodes depth as grayscale (black = near, white = far) |

Custom CPU shaders: any function with signature `func(shader.Attributes) math.Vec3`.

---

## HLSL shaders (Phase 13+)

HLSL is the authoring language for GPU shaders. At build time, shaders are compiled from HLSL to WGSL using `naga-cli` and embedded in the binary.

```hlsl
// assets/shaders/vertex_color.hlsl
cbuffer Transforms : register(b0) { float4x4 u_mvp; };

struct VSInput  { float3 pos : POSITION; float3 color : COLOR; };
struct VSOutput { float4 pos : SV_Position; float3 color : COLOR; };

VSOutput VSMain(VSInput input) {
    VSOutput o;
    o.pos   = mul(float4(input.pos, 1.0), u_mvp);
    o.color = input.color;
    return o;
}

float4 PSMain(VSOutput input) : SV_Target {
    return float4(input.color, 1.0);
}
```

Compile all shaders:
```bash
# Requires naga-cli: cargo install naga-cli
go generate ./assets/shaders/
```

---

## Building and testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./pkg/render/ ./pkg/rasterize/

# Update golden image reference (after intentional rendering changes)
go test ./pkg/render/ -run TestRender_GoldenImage_Triangle -update

# Lint (requires golangci-lint v2)
golangci-lint run

# Format
go fmt ./...
```

---

## Performance

### CPU (single-threaded, Phase 8 baseline)

Measured on Intel Xeon Platinum 8581C @ 2.10 GHz:

| Benchmark | Time |
|---|---|
| Full 640×480 cube render | ~1.2 ms/frame |
| 100×100 triangle render | ~30 µs/frame |
| Vertex MVP transform | ~58 ns/vertex |
| Small triangle rasterize | ~2 µs |
| Full-screen triangle rasterize | ~1.9 ms |
| Edge function | ~0.38 ns |

### GPU targets (Phase 11+)

| Target | Goal |
|---|---|
| Frame time at 1080p | <16.67 ms (60 fps) |
| Frame time at 4K | <16.67 ms (60 fps) |
| GPU rasterization vs CPU | >100× speedup expected |

---

## Coordinate system

- Right-handed: +X right, +Y up, +Z toward viewer
- Screen space: origin top-left, +X right, +Y down
- Depth buffer: 0 = near plane, 1 = far plane

---

## Learning goals

This project demonstrates the core stages of a 3D graphics pipeline:

1. **Math** — vectors, matrices, transformation composition
2. **Geometry** — vertex buffers, index buffers, triangle primitives
3. **Camera** — LookAt view matrix, perspective projection matrix
4. **Rasterization** — barycentric coordinates, interpolation, depth testing
5. **Shading** — per-pixel color from interpolated attributes (CPU) and HLSL (GPU)
6. **GPU pipeline** — WebGPU setup, vertex/index buffers, shader binding, render pass
7. **Real-time rendering** — frame loop, input, vsync, swap chain

---

## License

MIT
