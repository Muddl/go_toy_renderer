# go_toy_renderer

A toy 3D renderer written in Go. The CPU software renderer MVP is complete, and the full GPU roadmap (Phases 9–16) is now done: real-time windowed rendering via **wgpu-native** (WebGPU), WGSL shaders, Phong lighting, OBJ loading, texture mapping, and GPU timestamp profiling.

![CI](https://github.com/muddl/go_toy_renderer/actions/workflows/ci.yml/badge.svg)

---

## Features

### CPU renderer (Phases 0–8 — complete)
- World → camera → clip → NDC → screen space transforms
- Barycentric triangle rasterization with depth ordering
- Interpolated vertex colors; PNG output (640×480)

### GPU renderer (Phases 9–16 — complete)
- Real-time 1280×720 GLFW window at 60 fps
- WebGPU backend via **wgpu-native** (D3D12 on Windows, Metal on macOS, Vulkan on Linux)
- WGSL shaders embedded via `go:embed` — no external shader compiler needed
- Per-mesh model/normal matrix uniforms; scene graph with static/dynamic nodes
- **Phong lighting** — diffuse, specular (Blinn), ambient; orbiting camera
- **OBJ file loader** (`pkg/loader`) — positions, normals, UVs, fan triangulation
- **Texture mapping** — `pkg/gpu.Texture2D`; RGBA8 staging upload; 1×1 white fallback
- **GPU timestamp profiling** — `wgpu.QuerySet`; logs average GPU frame time to stdout once/second
- CPU/GPU overlay showing FPS, frame time, and backend

---

## Quick start

### CPU renderer (PNG output)

```bash
go run ./cmd/renderer
# Produces output.png — 640×480 PNG of a colored cube
```

### GPU real-time renderer

**Prerequisite on Windows:** 64-bit MinGW-w64 (see [Windows setup](#windows-setup) below).

```bash
# Auto-select best backend (GPU with CPU fallback)
go run ./cmd/renderer-rt

# Force GPU (wgpu-native)
go run ./cmd/renderer-rt --backend gpu

# Force CPU blit
go run ./cmd/renderer-rt --backend cpu

# Custom resolution
go run ./cmd/renderer-rt --width 1920 --height 1080
```

**Controls:** `ESC` or close the window to exit.

The demo scene renders a Utah Teapot at centre, a cube on the left, a ground plane, and a small cylinder at the orbiting camera position.

---

## Windows setup

`cmd/renderer-rt` uses CGo (GLFW). A 32-bit compiler produces a 32-bit PE that fails with **error 193** on x64 systems.

Install **64-bit MinGW-w64** via MSYS2:

```bash
pacman -S mingw-w64-x86_64-gcc
# Add to PATH: C:\msys64\mingw64\bin
```

Build with explicit architecture:

```bash
GOARCH=amd64 go build -o renderer-rt.exe ./cmd/renderer-rt
```

### wgpu-native setup

Pre-built wgpu-native v27.0.4.0 libraries are already committed under `assets/`:

| Platform | Path |
|---|---|
| Windows x64 | `assets/windows-x86_64-gnu/lib/wgpu_native.dll` |
| Linux ARM64 | `assets/linux-aarch64/lib/libwgpu_native.so` |
| macOS ARM64 | `assets/macos-aarch64/lib/libwgpu_native.dylib` |

Set the `WGPU_NATIVE_PATH` environment variable to the correct library for your platform, or place it in your PATH. `go run ./cmd/renderer-rt --backend gpu` will fail with a helpful error if the library is not found.

To run GPU integration tests locally:

```bash
GPU_TESTS=1 WGPU_NATIVE_PATH=assets/windows-x86_64-gnu/lib/wgpu_native.dll \
  go test ./pkg/gpu/...
```

---

## WGSL shader authoring

GPU shaders live in `assets/shaders/` and are embedded into the binary at build time via `//go:embed`. No external shader compiler is needed — edit the `.wgsl` file and rebuild.

### Current shader: `cube.wgsl`

**Vertex inputs** (stride 44 bytes):

| Location | Type | Meaning |
|---|---|---|
| `@location(0)` | `vec3<f32>` | Position |
| `@location(1)` | `vec3<f32>` | Vertex color |
| `@location(2)` | `vec3<f32>` | Normal (object space) |
| `@location(3)` | `vec2<f32>` | UV coordinates |

**Bind groups:**

| Group | Binding | Type | Contents |
|---|---|---|---|
| `@group(0)` | `@binding(0)` | uniform | `CameraUniforms` — `mat4x4<f32>` view-projection |
| `@group(0)` | `@binding(1)` | uniform | `MeshUniforms` — model matrix + normal matrix (128 bytes) |
| `@group(0)` | `@binding(2)` | uniform | `LightUniforms` — direction, color, ambient, camera pos (64 bytes) |
| `@group(1)` | `@binding(0)` | `texture_2d<f32>` | Albedo texture |
| `@group(1)` | `@binding(1)` | `sampler` | Linear sampler |

**Fragment output:** Blinn-Phong shading — `albedo × (ambient + diffuse·NdotL) + specular·(HdotN^32)`.

### Adding a new shader

1. Create `assets/shaders/myshader.wgsl`
2. Add an embed declaration in `assets/shaders/embed.go`:
   ```go
   //go:embed myshader.wgsl
   var MyShaderWGSL string
   ```
3. Pass `shaders.MyShaderWGSL` to `dev.CreateShaderModuleWGSL(...)` in `pkg/gpu/gpu.go`.

---

## Project structure

```
cmd/
  renderer/         # PNG output — CPU pipeline (complete)
  renderer-rt/      # Real-time GPU/CPU window — Phase 9–16 (complete)
  gen-assets/       # Standalone tool: generates assets/models/teapot.obj + .png

pkg/
  math/             # Vec3, Mat4x4, transformations, inverse, normal matrix
  geometry/         # Vertex (pos+color+normal+UV), Mesh, Cube, Tetrahedron,
                    #   Cylinder, Plane, primitives with correct normals
  camera/           # Perspective camera, LookAt, WebGPU projection (Z∈[0,1])
  framebuffer/      # Color + depth buffers; PNG export
  rasterize/        # Barycentric rasterizer with attribute interpolation
  shader/           # CPU per-pixel shaders (VertexColor, FlatColor, Depth)
  render/           # CPU pipeline — Scene, Render()
  renderer/         # Renderer interface + CPU/GPU/auto factory
  gpu/              # WebGPU device, buffers, pipelines, textures, profiling
  loader/           # OBJ file loader (positions, normals, UVs, fan triangulation)
  scene/            # Transform, Node, Scene — scene graph for GPU renderer
  overlay/          # Bitmap font, Metrics, FPS/frame-time debug overlay

assets/
  shaders/          # WGSL shaders embedded via go:embed
  models/           # teapot.obj (9248 vertices, 16384 triangles) + teapot.png
  windows-x86_64-gnu/lib/   # wgpu_native.dll
  linux-aarch64/lib/        # libwgpu_native.so
  macos-aarch64/lib/        # libwgpu_native.dylib
```

---

## Building and testing

```bash
# Run all tests (headless — no display or GPU required)
go test -tags=headless ./...

# Run with coverage
go test -tags=headless -cover ./...

# Lint (requires golangci-lint v2)
golangci-lint run --build-tags headless

# Format + vet
go fmt ./...
go vet ./...

# Vulnerability scan
govulncheck ./...
```

> **Note:** Use `-tags=headless` in any environment without a GPU or display.
> The CI pipeline uses this tag automatically.

---

## Performance

### CPU renderer (single-threaded, Phase 8 baseline)

Measured on Intel Xeon Platinum 8581C @ 2.10 GHz:

| Benchmark | Time |
|---|---|
| Full 640×480 cube render | ~1.2 ms/frame |
| 100×100 triangle render | ~30 µs/frame |
| Vertex MVP transform | ~58 ns/vertex |
| Small triangle rasterize | ~2 µs |
| Full-screen triangle rasterize | ~1.9 ms |

### GPU renderer (wgpu-native, Phase 16)

| Metric | Value |
|---|---|
| Target GPU frame time at 1080p | <1 ms |
| Target GPU frame time at 4K | <4 ms |
| GPU vs CPU speedup | >100× (hardware rasterization) |

GPU frame time is measured each frame using `wgpu.QuerySet` (timestamp queries) and logged to stdout as a per-second average:
```
GPU frame time: 0.42 ms avg (60 samples)
```
Falls back gracefully when the adapter does not expose `TIMESTAMP_QUERY`.

---

## CPU API example

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
    scene := render.NewScene()
    scene.AddMesh(geometry.NewCube())

    cam := camera.New(
        math.Vec3{X: 3, Y: 2, Z: 5},
        math.Vec3{X: 0, Y: 0, Z: 0},
        math.Vec3{X: 0, Y: 1, Z: 0},
        45.0, 640.0/480.0, 0.1, 100.0,
    )

    fb := framebuffer.New(640, 480)
    render.Render(scene, cam, fb, shader.VertexColor)

    if err := fb.SavePNG("output.png"); err != nil {
        log.Fatal(err)
    }
}
```

### Built-in CPU shaders

| Shader | Description |
|---|---|
| `shader.VertexColor` | Interpolates vertex colors across each triangle |
| `shader.NewFlatColor(color)` | Fills every pixel with a constant color |
| `shader.Depth` | Encodes depth as grayscale (black = near, white = far) |

Custom shader: any `func(shader.Attributes) math.Vec3`.

---

## Coordinate system

- Right-handed: +X right, +Y up, +Z toward viewer
- Depth buffer: 0 = near plane, 1 = far plane (WebGPU convention)
- Screen space: origin top-left, +X right, +Y down
- Winding: CCW front faces

---

## Learning goals

1. **Math** — vectors, matrices, transformation composition, inverse, normal matrix
2. **Geometry** — vertex buffers, index buffers, normals, UV coordinates
3. **Camera** — LookAt view matrix, perspective projection (OpenGL and WebGPU conventions)
4. **Rasterization** — barycentric coordinates, interpolation, depth testing
5. **Shading** — CPU per-pixel shaders and Blinn-Phong WGSL fragment shader
6. **GPU pipeline** — WebGPU init chain, swap chain, vertex/index/uniform buffers, bind groups
7. **Textures** — RGBA staging upload, sampler, UV-driven fragment sampling
8. **Scenes** — transform hierarchy, model/normal matrix uniforms, multi-mesh rendering
9. **OBJ loading** — Wavefront format parsing, vertex deduplication, fan triangulation
10. **Real-time rendering** — frame loop, vsync, input handling, GPU profiling

---

## License

MIT
