# go_toy_renderer

A toy 3D software renderer written in Go. All calculations run on the CPU — no GPU acceleration. The renderer implements a complete graphics pipeline from 3D geometry to a PNG image.

![CI](https://github.com/muddl/go_toy_renderer/actions/workflows/ci.yml/badge.svg)

## What it does

- Transforms 3D geometry through world → camera → clip → NDC → screen space
- Rasterizes triangles using barycentric coordinates with correct depth ordering
- Interpolates vertex colors across triangle surfaces
- Outputs a PNG image (640×480 by default)

## Quick start

```bash
# Build and run (outputs output.png in the current directory)
go run ./cmd/renderer

# Or build first, then run
go build -o renderer ./cmd/renderer
./renderer
```

The renderer produces `output.png`: a 640×480 image of a colored cube rendered with perspective projection.

## Project structure

```
cmd/renderer/       # Main application — renders output.png
pkg/
  math/             # Vec3, Mat4x4 — vectors, matrices, transformations
  geometry/         # Vertex, Mesh, and hardcoded primitives (Cube, Tetrahedron)
  camera/           # Perspective camera with LookAt view matrix
  framebuffer/      # Color + depth buffers; PNG export
  rasterize/        # Barycentric triangle rasterizer with attribute interpolation
  shader/           # Per-pixel shader system (VertexColor, FlatColor, Depth)
  render/           # Pipeline orchestration — Scene, Render()
```

## Using the API

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

### Shaders

Three built-in shaders are provided:

| Shader | Description |
|---|---|
| `shader.VertexColor` | Interpolates vertex colors across each triangle |
| `shader.NewFlatColor(color)` | Fills every pixel with a constant color |
| `shader.Depth` | Encodes depth as grayscale (black = near, white = far) |

Custom shaders are functions with the signature `func(shader.Attributes) math.Vec3`.

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

# Lint (requires golangci-lint v1.61.0)
golangci-lint run

# Format
go fmt ./...
```

## Performance (CPU, single-threaded)

Measured on an Intel Xeon Platinum 8581C @ 2.10 GHz:

| Benchmark | Time |
|---|---|
| Full 640×480 cube render | ~1.2 ms/frame |
| 100×100 triangle render | ~30 µs/frame |
| Vertex MVP transform | ~58 ns/vertex |
| Small triangle rasterize | ~2 µs |
| Full-screen triangle rasterize | ~1.9 ms |
| Edge function | ~0.38 ns |

## Coordinate system

- Right-handed: +X right, +Y up, +Z toward viewer
- Screen space: origin top-left, +X right, +Y down
- Depth buffer: 0 = near plane, 1 = far plane

## Learning goals

This project demonstrates the core stages of a 3D graphics pipeline:

1. **Math** — vectors, matrices, transformation composition
2. **Geometry** — vertex buffers, index buffers, triangle primitives
3. **Camera** — LookAt view matrix, perspective projection matrix
4. **Rasterization** — barycentric coordinates, interpolation, depth testing
5. **Shading** — per-pixel color calculation from interpolated attributes
6. **Output** — framebuffer, PNG export

## License

MIT
