# Architecture Reference — go_toy_renderer

**Last updated:** 2026-02-28 (Phase 10 complete — pkg/renderer interface implemented)

---

## CPU Pipeline Overview (Phases 0–8)

```
┌──────────────┐
│ Scene Setup  │  geometry.NewCube() → scene.AddMesh()
└──────┬───────┘
       │
┌──────▼───────┐
│ Vertex Stage │  transformVertex(): Local → Clip → NDC → Screen
│  (CPU)       │  mvp.MultiplyVec4 → perspective divide → viewport
└──────┬───────┘
       │
┌──────▼───────┐
│ Primitive    │  Triangle rejection: any vertex w ≤ 0 → skip
│ Assembly     │  (no frustum clip; handled by rejection + fb bounds)
└──────┬───────┘
       │
┌──────▼───────┐
│ Rasterizer   │  rasterize.TriangleShaded(): barycentric bounding-box
│  (CPU)       │  sweep; interpolate depth + color per pixel
└──────┬───────┘
       │
┌──────▼───────┐
│ Fragment     │  shader.Func called per pixel with Attributes
│ Shader       │  (VertexColor, NewFlatColor, Depth)
└──────┬───────┘
       │
┌──────▼───────┐
│ Framebuffer  │  SetPixel: depth < current → write; else discard
│  Depth Test  │  ColorBuffer []Vec3, DepthBuffer []float64
└──────┬───────┘
       │
┌──────▼───────┐
│ Image Output │  SavePNG: float→uint8 with clamping via image/png
└──────────────┘
```

### Vertex Coordinate Transform (transformVertex)

```
clip  = mvp.MultiplyVec4(px, py, pz, 1.0)   // clip space
                                              // reject if w ≤ 0
ndcX  = cx / cw                              // NDC x ∈ [-1, 1]
ndcY  = cy / cw                              // NDC y ∈ [-1, 1]
ndcZ  = cz / cw                              // NDC z ∈ [-1, 1]
depth = (ndcZ + 1) / 2                       // [0,1]: 0=near, 1=far
sX    = (ndcX + 1) / 2 * width              // screen pixels, origin TL
sY    = (1 - ndcY) / 2 * height             // Y flipped (screen +Y down)
```

---

## Coordinate System & Matrix Conventions

| Convention | Value |
|-----------|-------|
| Handedness | Right-handed (+X=Right, +Y=Up, +Z=Out of screen) |
| Matrix storage | Column-major `[16]float64` |
| Multiply order | `result = matrix × vector` (right-multiply) |
| Camera space | Camera looks down **−Z**; +Y up, +X right |
| VP matrix | `ProjectionMatrix.Multiply(ViewMatrix)` |
| Depth convention | 0.0 = near plane, 1.0 = far plane |
| Depth test | `depth < current` (closer fragment wins) |
| Depth init | 1.0 (far) — initialised by `Framebuffer.New` and `Clear` |
| Triangle winding | CCW = front face (rasterizer is winding-agnostic) |
| Screen origin | Top-left (0,0); +X right, +Y down |
| Pixel center | `(float64(ix)+0.5, float64(iy)+0.5)` |

---

## Package Dependency Rules

Dependencies flow in **one direction only**. Never import a higher-level package from a lower-level one.

```
cmd/renderer       → pkg/render, pkg/framebuffer, pkg/camera, pkg/geometry, pkg/shader
cmd/renderer-rt    → pkg/renderer*, pkg/window*         (* Phase 9+)
pkg/render         → pkg/rasterize, pkg/shader, pkg/framebuffer, pkg/camera, pkg/geometry, pkg/math
pkg/renderer*      → pkg/render, pkg/gpu*               (* Phase 10+)
pkg/gpu*           → pkg/geometry, pkg/math              (* Phase 11+)
pkg/rasterize      → pkg/shader, pkg/framebuffer, pkg/math
pkg/shader         → pkg/math
pkg/framebuffer    → pkg/math
pkg/camera         → pkg/math
pkg/geometry       → pkg/math
pkg/math           → (Go stdlib only)
```

---

## Package API Summary — MVP (Phases 0–8)

### `pkg/math`

**Status:** ✅ Complete | **Tests:** 74 (46 Vec3 + 28 Mat4x4) | **Coverage:** 100%

| Type / Function | Signature | Notes |
|----------------|-----------|-------|
| `Vec3` | `{X, Y, Z float64}` | Immutable by convention |
| `Vec3.Add` | `(Vec3) Vec3` | |
| `Vec3.Sub` | `(Vec3) Vec3` | |
| `Vec3.Scale` | `(float64) Vec3` | |
| `Vec3.Dot` | `(Vec3) float64` | |
| `Vec3.Cross` | `(Vec3) Vec3` | Right-handed |
| `Vec3.Normalize` | `() Vec3` | Guard: zero-length → undefined |
| `Vec3.Length` | `() float64` | |
| `Vec3.Distance` | `(Vec3) float64` | |
| `Vec3.Lerp` | `(Vec3, float64) Vec3` | |
| `Vec3.Equals` | `(Vec3, float64) bool` | Epsilon comparison |
| `Mat4x4` | `[16]float64` | Column-major |
| `NewIdentity` | `() Mat4x4` | |
| `NewZero` | `() Mat4x4` | |
| `Mat4x4.Multiply` | `(Mat4x4) Mat4x4` | |
| `Mat4x4.MultiplyVec3` | `(Vec3) Vec3` | Treats Vec3 as point (w=1) |
| `Mat4x4.MultiplyVec4` | `(x,y,z,w float64) (x,y,z,w)` | Full homogeneous transform |
| `Mat4x4.Transpose` | `() Mat4x4` | |
| `Mat4x4.Get/Set` | `(row,col int)` | Bounds-checked |
| `NewTranslation` | `(tx,ty,tz float64) Mat4x4` | |
| `NewScale` | `(sx,sy,sz float64) Mat4x4` | |
| `NewRotationX/Y/Z` | `(angle float64) Mat4x4` | Angle in radians, right-handed |

**Gotchas:** multiply order easy to reverse; divide by zero in Normalize; perspective divide W=0.

---

### `pkg/geometry`

**Status:** ✅ Complete | **Tests:** 21 (5 Vertex + 8 Mesh + 8 Primitives) | **Coverage:** ~100%

| Type / Function | Signature | Notes |
|----------------|-----------|-------|
| `Vertex` | `{Position, Color math.Vec3}` | Color in [0,1] range |
| `NewVertex` | `(pos, color Vec3) Vertex` | |
| `Vertex.Equals` | `(Vertex, float64) bool` | Epsilon on both fields |
| `Mesh` | `{Vertices []Vertex, Indices []int}` | |
| `NewMesh` | `() *Mesh` | |
| `Mesh.AddVertex` | `(Vertex) int` | Returns index |
| `Mesh.AddTriangle` | `(i0,i1,i2 int)` | CCW winding expected |
| `Mesh.GetTriangle` | `(i int) (int,int,int)` | Raw indices |
| `Mesh.GetTriangleVertices` | `(i int) (Vertex,Vertex,Vertex)` | Dereferenced |
| `Mesh.TriangleCount` | `() int` | `len(Indices)/3` |
| `Mesh.ValidateIndices` | `() error` | Bounds check all indices |
| `NewCube` | `() *Mesh` | 8 vertices, 12 triangles (36 indices), CCW |
| `NewTetrahedron` | `() *Mesh` | 4 vertices, 4 triangles (12 indices), CCW |

**Gotchas:** index out of bounds; degenerate (zero-area) triangles; inconsistent winding in OBJ files.

---

### `pkg/camera`

**Status:** ✅ Complete | **Tests:** 11 + 3 (Mat4x4) | **Coverage:** 100%

| Type / Function | Signature | Notes |
|----------------|-----------|-------|
| `Camera` | `{Position,Target,Up Vec3; FOV,Aspect,Near,Far float64}` | |
| `New` | `(pos,target,up Vec3, fov,aspect,near,far float64) Camera` | |
| `Camera.ViewMatrix` | `() Mat4x4` | LookAt; camera at origin looking down -Z |
| `Camera.ProjectionMatrix` | `() Mat4x4` | OpenGL-style perspective; near→z=-1, far→z=+1 |
| `Camera.ViewProjectionMatrix` | `() Mat4x4` | `Projection.Multiply(View)` |

**Gotchas:** Near must be > 0; far > near; up vector must not be parallel to forward vector.

---

### `pkg/framebuffer`

**Status:** ✅ Complete | **Tests:** 21 | **Coverage:** 94.6%

| Type / Function | Signature | Notes |
|----------------|-----------|-------|
| `Framebuffer` | `{Width,Height int; ColorBuffer []Vec3; DepthBuffer []float64}` | Linear: idx = y*W+x |
| `New` | `(width,height int) *Framebuffer` | Depth init to 1.0 |
| `Framebuffer.Clear` | `(color Vec3, depth float64)` | Resets all pixels |
| `Framebuffer.SetPixel` | `(x,y int, color Vec3, depth float64)` | Depth test; OOB silently ignored |
| `Framebuffer.GetPixel` | `(x,y int) Vec3` | OOB returns zero Vec3 |
| `Framebuffer.GetDepth` | `(x,y int) float64` | OOB returns 1.0 |
| `Framebuffer.SavePNG` | `(filename string) error` | Float→uint8 with clamping |

**Gotchas:** depth buffer not initialised = garbage; reversed depth test = wrong ordering; colour values > 1.0 clamped on export.

---

### `pkg/rasterize`

**Status:** ✅ Complete | **Tests:** 13 | **Coverage:** 100%

| Type / Function | Signature | Notes |
|----------------|-----------|-------|
| `ScreenVertex` | `{X,Y,Z float64; Color Vec3}` | Pixel coords; center at `ix+0.5, iy+0.5` |
| `Triangle` | `(v0,v1,v2 ScreenVertex, fb *Framebuffer)` | Calls `TriangleShaded` with `shader.VertexColor` |
| `TriangleShaded` | `(v0,v1,v2 ScreenVertex, fn shader.Func, fb *Framebuffer)` | Barycentric bounding-box rasteriser |

**Algorithm:** Bounding box clamped to framebuffer → barycentric coords per pixel → interpolate depth + color → call `shaderFn` → `fb.SetPixel`.

**Degenerate guard:** `area² < 1e-16` → silent skip.

**Winding:** Winding-agnostic; both CCW and CW render identically.

**Gotchas:** integer vs float coordinates; pixel-center sampling (`ix+0.5`); interpolated colours can exceed [0,1].

---

### `pkg/shader`

**Status:** ✅ Complete | **Tests:** 10 | **Coverage:** 100%

| Type / Function | Signature | Notes |
|----------------|-----------|-------|
| `Attributes` | `{Color Vec3; Depth float64}` | Interpolated per-fragment data |
| `Func` | `func(Attributes) Vec3` | Function type; no interface boilerplate |
| `VertexColor` | `func(Attributes) Vec3` | Pass-through: returns `attr.Color` |
| `NewFlatColor` | `(color Vec3) Func` | Closure; returns constant colour |
| `Depth` | `func(Attributes) Vec3` | Grayscale depth (0=black, 1=white) |

**Naming:** `shader.Func` not `shader.ShaderFunc` — revive linter: package name already provides context.

---

### `pkg/render`

**Status:** ✅ Complete | **Tests:** 12 unit + 7 integration + 1 golden image | **Coverage:** 100%

| Type / Function | Signature | Notes |
|----------------|-----------|-------|
| `Scene` | `{Meshes []*geometry.Mesh}` | Identity model matrix assumed for all meshes |
| `NewScene` | `() *Scene` | |
| `Scene.AddMesh` | `(m *geometry.Mesh)` | |
| `Render` | `(scene *Scene, cam Camera, fb *Framebuffer, shaderFn shader.Func)` | Stateless; clears fb, transforms all verts, rasterises |

**Internal:** `transformVertex(mvp Mat4x4, v Vertex, w,h int) (ScreenVertex, bool)` — returns `false` if `w ≤ 0`.

**Triangle rejection:** Any vertex with `w ≤ 0` → whole triangle skipped (no per-triangle clipping).

---

## Performance Baselines (Phase 8 Benchmarks)

Measured on a single thread, 640×480, coloured cube (12 triangles):

| Operation | Time |
|-----------|------|
| Full frame render (640×480 cube) | ~1.2 ms/frame |
| Vertex MVP transform | ~58 ns/vertex |
| Edge function | ~0.38 ns |

Run benchmarks: `go test -bench=. -benchmem ./pkg/render/... ./pkg/rasterize/...`

---

## GPU Architecture (Phases 9–16)

The CPU renderer is preserved as a **reference implementation and fallback**. A `Renderer` interface (Phase 10) abstracts CPU vs GPU backends.

### Backend Abstraction (Phase 10) ✅

```
┌──────────────────────────────────────────┐
│  pkg/renderer.Renderer interface         │
│  Init(width, height int) error           │
│  RenderFrame(scene *render.Scene) error  │
│  Shutdown()                              │
└───────────────────┬──────────────────────┘
                    │
         ┌──────────┴──────────┐
         │                     │
  ┌──────▼──────┐       ┌──────▼──────┐
  │ CPUBackend  │       │  GPUBackend │
  │ (cpu.go /   │       │  stub       │
  │  cpu_       │       │  Phase 11+  │
  │  headless)  │       └─────────────┘
  └─────────────┘
  pkg/render + GLFW
```

Factory: `renderer.New(backend string) (Renderer, error)` — `backend` is `"cpu" | "gpu" | "auto"`.

`ErrWindowClosed` — sentinel returned by `RenderFrame` when window is dismissed.

**Phase 10 package:** `pkg/renderer` | **Coverage:** 100% | **Tests:** 7 (GPU + factory + headless CPU)

| Type / Function | Signature | Notes |
|----------------|-----------|-------|
| `Renderer` | interface | `Init / RenderFrame / Shutdown` |
| `ErrWindowClosed` | `error` | RenderFrame sentinel; main loop breaks on this |
| `CPUBackend` | struct (`!headless`) | GLFW window + OpenGL blit; camera fixed at (3,2,5) looking at origin |
| `CPUBackend` | struct (`headless`) | Stub; all methods return error |
| `GPUBackend` | struct | Stub; `Init`→nil, `RenderFrame`→"GPU not yet implemented" |
| `New` | `(backend string) (Renderer, error)` | `"auto"` → CPUBackend (GPU-first in Phase 11) |

### GPU Pipeline (Phase 11+, wgpu-native / WebGPU)

```
┌──────────────┐
│  Scene Setup  │  geometry.Mesh → GPUScene (vertex/index buffers in VRAM)
└──────┬───────┘
       │
┌──────▼───────┐
│ Vertex Shader│  HLSL → WGSL (naga); MVP matrix from uniform buffer (binding 0)
│  (HLSL/WGSL) │
└──────┬───────┘
       │
┌──────▼───────┐
│ HW Rasteriser│  GPU barycentric interpolation (replaces CPU rasteriser)
│  (WebGPU)    │
└──────┬───────┘
       │
┌──────▼───────┐
│ Fragment     │  HLSL → WGSL; per-pixel colour, lighting, texture
│ Shader       │
└──────┬───────┘
       │
┌──────▼───────┐
│ HW Depth Test│  32-bit float depth; `Less` comparison; Depth32Float format
└──────┬───────┘
       │
┌──────▼───────┐
│ Swap Chain   │  BGRA8Unorm back buffer; present via D3D12/Metal/Vulkan
└──────────────┘
```

### Platform → Backend

| Platform | wgpu-native Backend | Notes |
|----------|-------------------|-------|
| Windows  | D3D12             | Primary; falls back to Vulkan |
| macOS    | Metal             | Only option on Apple Silicon |
| Linux    | Vulkan            | Primary; falls back to OpenGL |

### Shader Authoring (Phase 13)

```
HLSL source (.hlsl)       ← edit this
      │
  go generate
      │
   naga-cli                ← cargo install naga-cli
      │
WGSL output (.wgsl)       ← committed to repo
      │
  wgpu runtime             ← loads WGSL, compiles to native ISA
      │
  GPU execution
```

Compiled WGSL files are embedded in the binary via `//go:embed`.

### Post-GPU Package Layout (Phase 11+)

```
pkg/
├── math/           # unchanged
├── geometry/       # + Normal, UV in Phase 15/16
├── camera/         # unchanged
├── render/         # CPU pipeline preserved as reference
├── rasterize/      # CPU rasteriser preserved
├── shader/         # CPU shaders preserved
├── framebuffer/    # CPU framebuffer preserved
├── renderer/       # NEW (Phase 10): Renderer interface + factory
├── gpu/            # NEW (Phase 11): device, swapchain, pipeline, buffers
├── window/         # NEW (Phase 9): GLFW window, input, camera controller
├── scene/          # NEW (Phase 14): Transform component, GPUScene
└── loader/         # NEW (Phase 16): OBJ file loading

cmd/
├── renderer/       # PNG output (unchanged)
└── renderer-rt/    # NEW (Phase 9): real-time windowed renderer
```

### GPU Testing Notes

GPU tests must be skippable in CI (no GPU on runners):

```go
func TestGPUDevice_Init(t *testing.T) {
    if os.Getenv("GPU_TESTS") == "" {
        t.Skip("set GPU_TESTS=1 to run GPU integration tests")
    }
}
```

Run locally: `GPU_TESTS=1 go test ./pkg/gpu/...`

GPU visual output is validated against the CPU renderer with ±2 per-channel tolerance.

---

## Common Gotchas (Quick Reference)

See `conductor/product-guidelines.md` for the full annotated list.

| Area | Gotcha |
|------|--------|
| Math | Matrix multiply order easy to reverse |
| Math | W=0 divide → NaN; reject triangles with w≤0 |
| Transforms | VP = Projection × View (not View × Projection) |
| Transforms | Y-axis flipped: screen `+Y` down, NDC `+Y` up |
| Rasterizer | Pixel center at `ix+0.5`, not `ix` |
| Pipeline | Depth init must be 1.0 (far) before each frame |
| GPU | GLFW must run on OS main thread (`runtime.LockOSThread`) |
| GPU | Surface created before device on macOS Metal |
| GPU | Depth texture recreated on every window resize |
| Windows build | CGo links against C compiler; 32-bit MinGW → 32-bit PE → error 193 on x64. Use 64-bit MinGW-w64 and set `GOARCH=amd64`. CI pins this automatically. |
