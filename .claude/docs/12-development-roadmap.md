# Development Roadmap

## Overview

This roadmap breaks MVP development into phases, each building on the previous. Complete each phase before moving to the next.

**Estimated total time:** 2-4 weeks for hobbyist (few hours per day)

---

## Phase 0: CI/CD Infrastructure ✅ **COMPLETED** (2026-02-27)

**Goal:** Establish automated testing and quality gates before development begins.

### Tasks
1. Set up GitHub Actions CI/CD pipeline
   - Create `.github/workflows/ci.yml` with complete workflow
   - Configure multi-platform build matrix (Linux, macOS, Windows)
   - Set up test execution with race detector and coverage reporting
   - Add golangci-lint for code quality enforcement

2. Configure linter settings
   - Create `.golangci.yml` with recommended Go linters
   - Enable 48+ linters (gofmt, govet, errcheck, staticcheck, etc.)
   - Set reasonable complexity thresholds for learning project

3. Integrate security scanning
   - Add govulncheck for vulnerability detection
   - Configure scheduled security scans (weekly)
   - Set up GitHub Security tab integration

4. Set up coverage enforcement
   - Configure 70% overall coverage threshold
   - Configure 90% math package coverage threshold
   - Automatic coverage reporting to Codecov (optional)

5. Add CI badges to README
   - Build status badge
   - Coverage badge (if using Codecov)
   - Go version badge

### Completion Criteria
- [x] CI workflow file created and tested
- [x] All jobs passing on main branch (merged via PR #5)
- [x] Linter configuration in place (migrated to golangci-lint v2)
- [x] Coverage thresholds enforced
- [x] Security scanning operational
- [x] Local pre-commit checks documented

**When complete:** All future PRs will be automatically validated for quality, tests, and security.

**Status:** ✅ **COMPLETED** (2026-02-27) - Merged via PR #5

**Time taken:** ~1 day (configuration, golangci-lint v2 migration, CI fixes)

---

## Phase 1: Math Foundation ✅ **COMPLETED** (Day 1)

**Goal:** Rock-solid math library that everything else depends on.

**Status:** Completed 2025-10-10

### Tasks Completed ✅
1. ✅ Implement Vec3 type and operations
   - Add, subtract, scale, dot, cross, normalize, length, distance
   - 46 comprehensive tests with known results
   - 100% test coverage

2. ✅ Implement Mat4x4 type and operations
   - Identity, zero, multiply, transpose, element access
   - Comprehensive tests for matrix multiplication properties

3. ✅ Create basic transformation matrices
   - Translation, rotation (X/Y/Z axes), scale
   - Tests for known transformations (90° rotation, identity, etc.)

### Tasks Deferred to Phase 3 (Camera System)
- ⏳ LookAt matrix (view matrix) - part of Camera implementation
- ⏳ Perspective projection - part of Camera implementation
- ⏳ Viewport transformation - part of Camera implementation

### Completion Criteria
- ✅ All vector operations pass tests (46 tests)
- ✅ Matrix operations pass tests (comprehensive coverage)
- ✅ Can create basic transformation matrices (translation, rotation, scale)
- ✅ >90% test coverage in math package (100% achieved)
- ✅ No known bugs
- ⏳ Camera matrices deferred to Phase 3

**Achievement:** Foundation is solid. Ready for Phase 2 (Geometry & Scene).

**Time taken:** ~1 day (faster than 2-3 day estimate due to TDD efficiency)

**Next:** Phase 2 - Geometry & Scene component (after Phase 0 CI/CD merge)

**Status:** ✅ **COMPLETED** (2025-10-10)

---

## Phase 2: Geometry & Scene ✅ **COMPLETED** (2026-02-27)

**Goal:** Define 3D objects and organize them in a scene.

### Tasks Completed ✅
1. ✅ Implement Vertex type (`pkg/geometry/vertex.go`, 27 lines)
   - Position (Vec3) and Color (Vec3)
   - NewVertex constructor, Equals with epsilon comparison
   - 5 tests with 100% coverage

2. ✅ Implement Mesh type (`pkg/geometry/mesh.go`, 63 lines)
   - Vertex buffer ([]Vertex) and index buffer ([]int)
   - AddVertex, AddTriangle, GetTriangle, GetTriangleVertices, ValidateIndices
   - 8 comprehensive tests

3. ✅ Create hardcoded primitives (`pkg/geometry/primitives.go`, 105 lines)
   - Tetrahedron: 4 vertices, 4 triangles (12 indices), unique per-vertex colors
   - Cube: 8 vertices, 12 triangles (36 indices), counter-clockwise winding
   - 8 tests validating counts, indices, winding order, colors

### Completion Criteria
- [x] Can create tetrahedron (4 vertices, 4 triangles) and cube (8 vertices, 12 triangles)
- [x] Meshes have valid vertex/index data (ValidateIndices passes)
- [x] All tests pass — geometry package ~95%+ coverage
- [x] Code formatted and linted, merged to main

**When complete:** You have test objects to render.

**Status:** ✅ **COMPLETED** (2026-02-27) - Merged via PR #5

**Time taken:** ~1 day

**Next:** Phase 3 - Camera System

---

## Phase 3: Camera System ✅ **COMPLETED** (2026-02-27)

**Goal:** Flexible viewpoint control.

### Tasks Completed ✅
1. ✅ Implement Camera type (`pkg/camera/camera.go`)
   - Position, Target, Up, FOV, Aspect, Near, Far
   - `New()` constructor

2. ✅ Implement view matrix generation (`ViewMatrix()`)
   - LookAt algorithm: right-handed, camera looks down -Z
   - Basis vectors: forward = normalize(eye-target), right = cross(up,forward), newUp = cross(forward,right)
   - 3 tests: origin-at-camera, basis orthonormality, camera-position-maps-to-origin

3. ✅ Implement projection matrix generation (`ProjectionMatrix()`)
   - OpenGL-style perspective: near → NDC z=-1, far → NDC z=+1
   - Points behind camera have negative W in clip space
   - 4 tests: near/far plane mapping, FOV scaling, aspect scaling, behind-camera rejection

4. ✅ Combine into view-projection matrix (`ViewProjectionMatrix()`)
   - Returns `Projection.Multiply(View)`
   - 2 tests: equality to explicit product, world-origin projects to NDC center (0,0)

5. ✅ Extended `pkg/math/Mat4x4` with `MultiplyVec4(x,y,z,w)`
   - Required for projection where W encodes depth
   - 3 tests: identity, translation with w=1, direction vector with w=0

### Completion Criteria
- [x] Camera generates correct view matrix (camera position maps to origin in camera space)
- [x] Camera generates correct projection matrix (near/far map to correct NDC depth)
- [x] Combined VP matrix transforms points correctly (world origin projects to NDC center)
- [x] 13 camera/math tests covering all key properties

**When complete:** You can control viewpoint.

**Status:** ✅ **COMPLETED** (2026-02-27)

**Next:** Phase 4 - Framebuffer (complete)

---

## Phase 4: Framebuffer ✅ **COMPLETED** (2026-02-27)

**Goal:** Store rendered pixels and output images.

### Tasks Completed ✅
1. ✅ Implement Framebuffer type (`pkg/framebuffer/framebuffer.go`)
   - Color buffer (`[]math.Vec3`, RGB float64 [0,1] per pixel)
   - Depth buffer (`[]float64`, initialized to 1.0 far plane)
   - Linear layout: index = `y*Width + x`

2. ✅ Implement Clear method
   - Resets all pixels to given color and depth in one pass

3. ✅ Implement SetPixel with depth test
   - Strictly-less-than depth test (`depth < current`)
   - Out-of-bounds coordinates silently ignored

4. ✅ Implement GetPixel and GetDepth
   - Return zero Vec3 / 1.0 for out-of-bounds access

5. ✅ Implement SavePNG
   - Float RGB [0,1] → uint8 [0,255] with clamping via `clampToByte`
   - Uses `image.NewNRGBA` + `png.Encode` from Go stdlib
   - 21 tests across all operations: 94.6% coverage

### Completion Criteria
- [x] Can create framebuffer of any size
- [x] Clear works correctly
- [x] Depth test prevents farther pixels from overwriting closer ones
- [x] Can save valid PNG file
- [x] Framebuffer tests pass (21 tests, 94.6% coverage)

**When complete:** You can store and save rendered images.

**Status:** ✅ **COMPLETED** (2026-02-27)

**Time taken:** ~1 day

**Next:** Phase 5 - Rasterization

---

## Phase 5: Rasterization ✅ **COMPLETED** (2026-02-27)

**Goal:** Convert triangles to pixels with interpolation.

### Tasks Completed ✅
1. ✅ Chose rasterization algorithm — Barycentric coordinates (cleaner code, winding-agnostic)

2. ✅ Implemented `ScreenVertex` type (`pkg/rasterize/rasterizer.go`)
   - X, Y: screen coordinates (pixel units, center at float(ix)+0.5)
   - Z: depth value in [0, 1]
   - Color: RGB Vec3 in [0, 1]

3. ✅ Implemented barycentric coordinate calculation
   - `edgeFunction(ax, ay, bx, by, px, py)` — 2D cross product
   - `insideTriangle(area, w0, w1, w2)` — handles both CCW and CW winding

4. ✅ Implemented tight bounding box with framebuffer clamping
   - int(min) ... int(max)+1, clamped to [0, Width) × [0, Height)

5. ✅ Implemented rasterization loop with pixel-center sampling
   - px = float64(ix) + 0.5, py = float64(iy) + 0.5

6. ✅ Implemented barycentric attribute interpolation
   - Depth: b0*v0.Z + b1*v1.Z + b2*v2.Z
   - Color: per-channel linear interpolation

7. ✅ Connected to framebuffer via `fb.SetPixel(ix, iy, color, depth)`
   - Depth test handled transparently by framebuffer

8. ✅ 9 comprehensive tests covering all behaviors (100% coverage)
   - Degenerate triangle, single-pixel, axis-aligned, vertex color, centroid interpolation
   - Depth ordering, off-screen, partially off-screen, CW winding

### Completion Criteria
- [x] Can rasterize any triangle to pixels
- [x] Attribute interpolation is correct at vertices
- [x] Depth interpolation works
- [x] No crashes on degenerate triangles
- [x] Rasterizer tests pass (9 tests, 100% coverage)

**When complete:** You can draw triangles to framebuffer.

**Status:** ✅ **COMPLETED** (2026-02-27)

**Time taken:** ~1 day

**Next:** Phase 6 - Shading

---

## Phase 6: Shading ✅ **COMPLETED** (2026-02-27)

**Goal:** Calculate per-pixel colors.

### Tasks Completed ✅
1. ✅ Defined `Func` type (`shader.Func`): `func(Attributes) math.Vec3` — simple function type, no interface
   - `Attributes` struct: `Color math.Vec3` (interpolated vertex color) + `Depth float64`
   - Named `Func` (not `ShaderFunc`): revive linter flags `shader.ShaderFunc` as redundant

2. ✅ Implemented `VertexColor` shader (`pkg/shader/shader.go`)
   - Pass-through: returns `attr.Color` unchanged
   - Demonstrates correct color interpolation

3. ✅ Implemented `NewFlatColor(color math.Vec3) ShaderFunc`
   - Constructor returning closure capturing constant color
   - Useful for debugging and solid-fill testing

4. ✅ Implemented `Depth` shader
   - Converts interpolated depth to grayscale (depth 0→black, 1→white)
   - Useful for verifying depth buffer correctness

5. ✅ 10 tests across all shaders — 100% statement coverage

### Completion Criteria
- [x] Shader interface defined (`ShaderFunc` function type)
- [x] At least one shader implemented (`VertexColor`, plus `NewFlatColor` and `Depth`)
- [x] Shader receives correct interpolated data (via `Attributes` struct)
- [x] Shader tests pass (10 tests, 100% coverage)

**When complete:** You can control per-pixel appearance.

**Status:** ✅ **COMPLETED** (2026-02-27)

**Time taken:** ~1 day

**Next:** Phase 7 - Render Pipeline Integration

---

## Phase 7: Render Pipeline Integration ✅ **COMPLETED** (2026-02-27)

**Goal:** Connect all components into working renderer.

### Tasks Completed ✅
1. ✅ Extended `pkg/rasterize/rasterizer.go` with `TriangleShaded` function
   - Added `shader.Func` parameter for per-pixel shading
   - `Triangle` refactored to delegate to `TriangleShaded(shader.VertexColor)` — zero duplication
   - 4 new `TriangleShaded` tests (13 total, 100% coverage)

2. ✅ Implemented `pkg/render/` package
   - `Scene` struct: holds `[]*geometry.Mesh` for MVP (identity model matrix assumed)
   - `NewScene()`, `AddMesh()` for scene building
   - `Render(scene, cam, fb, shaderFn)` — stateless, clears framebuffer, transforms & rasterizes all meshes
   - `transformVertex()` helper: MVP × vertex → clip space → perspective divide → NDC → screen coords
   - Triangle rejection: any vertex with `w ≤ 0` (behind camera) causes the whole triangle to be skipped
   - 12 tests, 100% coverage

3. ✅ Implemented full `cmd/renderer/main.go`
   - Renders a colored cube (camera at {3,2,5} looking at origin)
   - 640×480 framebuffer, VertexColor shader
   - Saves `output.png`

4. ✅ Coordinate transform pipeline:
   - `clip = mvp.MultiplyVec4(px, py, pz, 1)` (MVP = ViewProjectionMatrix, identity model)
   - Reject if `w ≤ 0`
   - `ndcX = cx/cw`, `ndcY = cy/cw`, `ndcZ = cz/cw`
   - `depth = (ndcZ + 1) / 2` maps NDC z [-1,1] → [0,1]
   - `screenX = (ndcX+1)/2 * width`, `screenY = (1-ndcY)/2 * height` (Y flipped)

### Completion Criteria
- [x] Can render complete mesh to image
- [x] Perspective looks correct
- [x] Depth ordering is correct (closer triangle overwrites farther)
- [x] Colors interpolate smoothly (VertexColor shader pass-through)
- [x] Output image `output.png` saved successfully (640×480 PNG)
- [x] No crashes on valid input

**When complete:** MVP is functionally complete!

**Status:** ✅ **COMPLETED** (2026-02-27)

**Time taken:** ~1 day

---

## Phase 8: Testing & Polish ✅ **COMPLETED** (2026-02-27)

**Goal:** Ensure quality and reliability.

### Tasks Completed ✅
1. ✅ Write integration tests (`pkg/render/integration_test.go`)
   - 6 new integration tests covering multiple camera positions, all 3 shaders, multiple meshes
   - Tests verify visual correctness at pixel level

2. ✅ Create golden image tests
   - `TestRender_GoldenImage_Triangle` with `-update` flag support
   - Reference stored in `pkg/render/testdata/golden_triangle.png`
   - Byte-exact PNG comparison on every test run

3. ✅ Coverage well above thresholds
   - camera: 100% | geometry: 100% | math: 96.1% | rasterize: 100% | render: 100% | shader: 100%
   - framebuffer: 90.2% | overall: well above 70% minimum

4. ✅ Demo application complete (Phase 7)
   - `cmd/renderer/main.go` renders colored cube to `output.png` (640×480)

5. ✅ Write documentation
   - `README.md` fully written with usage instructions, API example, performance table
   - Code comments throughout all public APIs

6. ✅ Performance benchmarks
   - `pkg/render/bench_test.go`: full pipeline, triangle, vertex transform
   - `pkg/rasterize/rasterizer_bench_test.go`: small/large triangle, edge function
   - Baseline: 640×480 cube ~1.2 ms/frame; vertex transform ~58 ns

### Completion Criteria
- [x] All tests passing
- [x] >70% test coverage (all packages well above threshold)
- [x] At least one golden image test
- [x] Demo app produces nice output image
- [x] README explains how to use
- [x] Code is reasonably documented

**When complete:** MVP is ready to show!

**Status:** ✅ **COMPLETED** (2026-02-27)

**Time taken:** ~1 day

---

## GPU Acceleration & Real-time Rendering Roadmap

**North Star:** Cross-platform GPU acceleration with real-time rendering and HLSL shaders.

**GPU Strategy:** WebGPU (via wgpu-native) as the backend — uses D3D12 on Windows, Metal on macOS, and Vulkan on Linux, providing a single unified API. HLSL shaders are the authoring language, transpiled to WGSL (WebGPU Shading Language) for the runtime via the `naga` toolchain.

**Architecture Principle:** The CPU renderer (Phases 0–8) is preserved as a reference implementation and fallback. A `Renderer` interface abstracts CPU vs. GPU backends, so tests and tooling continue to work against the CPU backend while GPU development proceeds.

---

## Phase 9: Window & Real-time Display

**Goal:** Replace PNG file output with a live window at 60 fps.

**Status:** 📋 Planned

### Tasks
1. **Choose windowing library**
   - Evaluate `github.com/go-gl/glfw/v3.3/glfw` (GLFW) or an SDL2 wrapper
   - GLFW is preferred: maintained, cross-platform (Windows/macOS/Linux), integrates with WebGPU surface creation
   - Create `pkg/window/` package with `Window` struct (width, height, title, event loop)

2. **Implement event loop**
   - Main loop: poll events → update → render → present
   - Handle window close, ESC key
   - Expose `OnUpdate(dt float64)` and `OnRender()` callbacks

3. **CPU blit path (bridge to GPU)**
   - Render to existing `framebuffer.Framebuffer` (CPU)
   - Upload pixel data to a GPU texture each frame
   - Blit GPU texture to the swap chain back buffer
   - Allows running existing CPU renderer in a window immediately (no GPU pipeline required yet)

4. **Basic camera controls**
   - WASD + mouse look for interactive navigation
   - Camera controller in `pkg/window/` or `pkg/input/`

5. **Frame timing & vsync**
   - Track frame time and FPS counter
   - Vsync on by default; `--no-vsync` flag for benchmarking

6. **New entry point**
   - `cmd/renderer-rt/main.go` — real-time windowed renderer
   - Original `cmd/renderer/main.go` (PNG output) retained unchanged

### Completion Criteria
- [ ] Window opens on Windows, macOS, and Linux
- [ ] CPU-rendered framebuffer displayed live at ≥60 fps
- [ ] Camera moves with WASD/mouse
- [ ] Window closes cleanly with ESC or X button
- [ ] Frame time displayed in title bar
- [ ] `cmd/renderer` (PNG) still works unchanged

**Estimated duration:** 3–5 days

**Dependencies:** Phase 0–8 (MVP complete)

---

## Phase 10: GPU Backend Abstraction Layer

**Goal:** Introduce a `Renderer` interface so the CPU and future GPU backends are interchangeable.

**Status:** 📋 Planned

### Tasks
1. **Define `Renderer` interface** in `pkg/renderer/`
   ```go
   type Renderer interface {
       Render(scene *render.Scene, cam camera.Camera) error
       Resize(width, height int)
       Destroy()
   }
   ```

2. **Wrap existing pipeline as `CPURenderer`**
   - Moves logic from `pkg/render/render.go` into `pkg/renderer/cpu.go`
   - Holds `framebuffer.Framebuffer` internally; exposes `ReadPixels()` for display
   - All existing tests continue to pass unchanged

3. **Create `GPURenderer` stub**
   - `pkg/renderer/gpu.go` — returns `ErrNotImplemented`
   - Placeholder accepted by the interface, replaced in Phase 11

4. **Backend selection at startup**
   - `--backend cpu|gpu|auto` CLI flag (auto = GPU if available, CPU fallback)
   - Factory function: `renderer.New(backend, width, height) (Renderer, error)`

5. **Update `cmd/renderer-rt/main.go`**
   - Use `Renderer` interface, selectable backend

### Completion Criteria
- [ ] `Renderer` interface defined and documented
- [ ] `CPURenderer` passes all existing render tests
- [ ] `GPURenderer` stub compiles and returns correct error
- [ ] Backend selection flag works in `cmd/renderer-rt`
- [ ] No regression in `cmd/renderer` (PNG path)

**Estimated duration:** 2–3 days

**Dependencies:** Phase 9

---

## Phase 11: WebGPU Integration (wgpu-native)

**Goal:** Initialize the GPU and render the first triangle via WebGPU.

**Status:** 📋 Planned

### Tasks
1. **Add wgpu-native dependency**
   - Library: `github.com/rajveermalviya/go-webgpu/wgpu` (Go bindings for wgpu-native)
   - wgpu-native provides WebGPU on top of D3D12 (Windows), Metal (macOS), Vulkan (Linux)
   - Add cgo build constraints and document native library requirements

2. **GPU device initialisation** (`pkg/gpu/device.go`)
   - Create `wgpu.Instance` → request adapter → request device
   - Log adapter info (GPU name, driver version, backend)
   - Handle device lost / error callbacks

3. **Surface & swap chain** (`pkg/gpu/swapchain.go`)
   - Create `wgpu.Surface` from GLFW window handle
   - Configure swap chain: BGRA8Unorm preferred, mailbox/vsync present mode
   - Handle swap chain resize

4. **Depth attachment**
   - Create `Depth32Float` texture + view for depth testing
   - Resize depth texture with window

5. **Hello triangle render pass**
   - Hardcoded triangle vertices in WGSL vertex shader
   - Solid color WGSL fragment shader
   - `wgpu.RenderPassDescriptor` with color + depth attachments
   - Submit command buffer, present

6. **Integrate with `GPURenderer`**
   - `GPURenderer.Render()` submits a render pass each frame
   - `GPURenderer.Resize()` recreates swap chain + depth texture

### Completion Criteria
- [ ] GPU device and surface initialised without errors on all three platforms
- [ ] Hardcoded triangle visible in window
- [ ] Depth buffer attached and functional
- [ ] Swap chain resizes correctly on window resize
- [ ] `GPURenderer` implements full `Renderer` interface
- [ ] `--backend gpu` selects GPU path successfully

**Estimated duration:** 4–5 days

**Dependencies:** Phase 10

**New packages:** `pkg/gpu/` (device, swapchain, depth buffer)

---

## Phase 12: GPU Geometry & Draw Pipeline

**Goal:** Upload scene geometry to VRAM and render meshes with hardware rasterization.

**Status:** 📋 Planned

### Tasks
1. **GPU vertex buffer** (`pkg/gpu/buffer.go`)
   - Convert `geometry.Mesh` vertices to a flat `[]float32` (interleaved position + color)
   - `wgpu.Device.CreateBufferInit()` — upload to VRAM
   - Vertex layout: `location 0` = position (vec3f), `location 1` = color (vec3f)
   - Stride: 24 bytes per vertex

2. **GPU index buffer**
   - Convert `geometry.Mesh` index slice to `[]uint32`
   - Upload to VRAM as index buffer

3. **Vertex buffer layout descriptor**
   - `wgpu.VertexBufferLayout` matching the interleaved format above

4. **Basic WGSL vertex shader** (`assets/shaders/basic.wgsl`)
   ```wgsl
   @vertex
   fn vs_main(@location(0) pos: vec3f, @location(1) color: vec3f) -> VertexOutput { … }
   ```

5. **Basic WGSL fragment shader**
   ```wgsl
   @fragment
   fn fs_main(in: VertexOutput) -> @location(0) vec4f { return vec4f(in.color, 1.0); }
   ```

6. **Render pipeline** (`pkg/gpu/pipeline.go`)
   - `wgpu.RenderPipelineDescriptor` with vertex + fragment shader modules
   - Primitive topology: triangle list
   - Depth-stencil state: `Less`, write enabled
   - Backface culling: CCW front face, cull back

7. **Draw call submission**
   - `render_pass.SetVertexBuffer(0, vertex_buffer, 0, size)`
   - `render_pass.SetIndexBuffer(index_buffer, Uint32, 0, size)`
   - `render_pass.DrawIndexed(index_count, 1, 0, 0, 0)`

8. **Scene upload management**
   - `GPUScene` struct mirrors `render.Scene` with GPU buffers
   - Re-upload buffers only when geometry changes (dirty flag)

### Completion Criteria
- [ ] Colored cube rendered via GPU hardware rasterization
- [ ] Index buffer reduces vertex redundancy correctly
- [ ] Depth test orders faces correctly
- [ ] Backface culling removes hidden faces
- [ ] Visual output matches CPU renderer (same camera, same geometry)

**Estimated duration:** 4–5 days

**Dependencies:** Phase 11

**New packages:** `pkg/gpu/buffer.go`, `pkg/gpu/pipeline.go`

**New assets:** `assets/shaders/basic.wgsl`

---

## Phase 13: HLSL Shader Pipeline

**Goal:** Author shaders in HLSL; compile to WGSL for the WebGPU runtime.

**Status:** 📋 Planned

### Why HLSL → WGSL?
WebGPU's native language is WGSL. HLSL is the authoring target because it is expressive, widely documented, and supported by DirectX tooling. The `naga` library (from the wgpu project) can transpile HLSL → SPIR-V or WGSL at build time, so HLSL files are compiled once and the resulting WGSL is embedded in the binary.

### Tasks
1. **Shader compilation toolchain** (`tools/compile_shaders.go` or `Makefile`)
   - Install `dxc` (DirectX Shader Compiler) for HLSL → SPIR-V
   - Or use `naga-cli` for HLSL → WGSL directly
   - Build step: `go generate ./assets/shaders/` runs the compiler
   - Output: `assets/shaders/compiled/*.wgsl` (committed to repo)

2. **Port vertex color shader to HLSL** (`assets/shaders/vertex_color.hlsl`)
   ```hlsl
   struct VSInput  { float3 pos : POSITION; float3 color : COLOR; };
   struct VSOutput { float4 pos : SV_Position; float3 color : COLOR; };

   VSOutput VSMain(VSInput input) {
       VSOutput output;
       output.pos   = mul(float4(input.pos, 1.0), u_mvp);
       output.color = input.color;
       return output;
   }

   float4 PSMain(VSOutput input) : SV_Target {
       return float4(input.color, 1.0);
   }
   ```

3. **Port flat-color and depth shaders to HLSL**
   - `flat_color.hlsl` — constant colour from push constant
   - `depth.hlsl` — depth value → grayscale output

4. **Uniform buffer (cbuffer) for MVP matrix**
   - `cbuffer Transforms { float4x4 u_mvp; };` in every vertex shader
   - `wgpu.BindGroupLayout` with a uniform buffer binding
   - Update uniform buffer per-frame with camera VP matrix

5. **Shader loader** (`pkg/gpu/shader.go`)
   - Load compiled WGSL from `assets/shaders/compiled/`
   - `wgpu.Device.CreateShaderModule()` from WGSL source
   - Cache compiled modules

6. **Custom shader support**
   - Users drop `.hlsl` files in `assets/shaders/` and run `go generate`
   - `GPURenderer` accepts a shader name at construction time

### Completion Criteria
- [ ] HLSL compilation toolchain documented in README and `go generate`
- [ ] `vertex_color.hlsl`, `flat_color.hlsl`, `depth.hlsl` compile and run correctly
- [ ] MVP matrix passed via uniform buffer, perspective correct
- [ ] Custom shader can be specified by filename
- [ ] `assets/shaders/compiled/` checked in for reproducible builds

**Estimated duration:** 4–5 days

**Dependencies:** Phase 12

**New packages:** `pkg/gpu/shader.go`

**New assets:** `assets/shaders/*.hlsl`, `assets/shaders/compiled/*.wgsl`

---

## Phase 14: Uniform Buffers & Per-Mesh Transforms

**Goal:** Support multiple meshes with independent model matrices; full MVP per mesh.

**Status:** 📋 Planned

### Tasks
1. **Per-mesh model matrix**
   - Add `ModelMatrix math.Mat4x4` to `render.Scene` mesh entries
   - Upload per-mesh model matrix to uniform buffer (dynamic offset or push constant)
   - HLSL: `cbuffer PerObject { float4x4 u_model; float4x4 u_vp; };`

2. **Bind group per draw call**
   - Option A: One large uniform buffer, dynamic offsets per mesh (preferred — fewer bind group switches)
   - Option B: One `wgpu.BindGroup` per mesh (simpler, acceptable for low mesh counts)

3. **Camera uniform buffer**
   - Separate `cbuffer Camera { float4x4 u_view; float4x4 u_proj; float3 u_camPos; };`
   - Updated once per frame, shared by all meshes

4. **Multi-mesh scene rendering**
   - Loop: for each mesh → bind per-object uniform → draw
   - Verify each mesh at independent position/rotation

5. **Transform component** (`pkg/scene/transform.go`)
   - `Transform` struct: Position, Rotation (Euler angles), Scale (Vec3)
   - `ModelMatrix() math.Mat4x4` → TRS composition

### Completion Criteria
- [ ] Multiple cubes renderable at different positions simultaneously
- [ ] Per-mesh model matrix applied correctly (rotation, scale, translation)
- [ ] Camera uniform updated correctly every frame
- [ ] No visible transform artifacts between frames

**Estimated duration:** 3–4 days

**Dependencies:** Phase 13

**New packages:** `pkg/scene/transform.go`

---

## Phase 15: Lighting, Normals & Physical Shading

**Goal:** Implement Phong and/or PBR shading with vertex normals in HLSL.

**Status:** 📋 Planned

### Tasks
1. **Vertex normals** (`pkg/geometry/vertex.go`)
   - Add `Normal math.Vec3` field to `Vertex`
   - Update GPU vertex layout: stride grows to 36 bytes (pos + color + normal)
   - Recompute/hardcode normals for Cube and Tetrahedron primitives

2. **Normal transformation**
   - In vertex shader: transform normal by inverse-transpose of model matrix
   - HLSL: `float3 worldNormal = normalize(mul(input.normal, (float3x3)u_normalMatrix));`

3. **Phong lighting shader** (`assets/shaders/phong.hlsl`)
   - Ambient + diffuse + specular terms
   - Directional light: direction, color, intensity from `cbuffer Lights { … }`
   - Material properties: albedo, shininess from push constant or cbuffer

4. **Light uniform buffer**
   - `cbuffer Lights { float3 u_lightDir; float3 u_lightColor; float u_lightIntensity; float3 u_ambient; };`
   - Updated per frame

5. **Normal map support (optional stretch goal)**
   - Tangent space normal mapping
   - Requires tangent vector in vertex layout

6. **Simple PBR shader** (`assets/shaders/pbr.hlsl`) — stretch goal
   - GGX specular BRDF
   - Metallic-roughness workflow
   - Image-based lighting (IBL) deferred to Phase 16

### Completion Criteria
- [ ] Vertex normals correct for all primitive types
- [ ] Phong shader compiles and applies lighting visually
- [ ] Directional light controllable from `cmd/renderer-rt`
- [ ] Specular highlight visible on shiny materials
- [ ] CPU renderer unaffected (normals stored in Vertex but ignored by CPU shaders)

**Estimated duration:** 4–5 days

**Dependencies:** Phase 14

---

## Phase 16: Advanced Features & Polish

**Goal:** Production-ready features: textures, OBJ loading, post-processing, performance tuning.

**Status:** 📋 Planned

### Tasks
1. **Texture mapping** (`pkg/gpu/texture.go`)
   - Load PNG/JPG/KTX2 textures into `wgpu.Texture`
   - Add `UV math.Vec2` to `Vertex` layout
   - HLSL sampler: `Texture2D t_albedo; SamplerState s_albedo;`
   - Mipmapping via `wgpu.Texture.GenerateMipmaps()` or pre-generated

2. **OBJ file loader** (`pkg/loader/obj.go`)
   - Parse `.obj` + `.mtl` files into `geometry.Mesh` slices
   - Support positions, normals, UV coordinates, face indices
   - Basic material extraction (diffuse color, texture path)

3. **Post-processing pass** (`pkg/gpu/postprocess.go`)
   - Full-screen triangle pass over HDR render target
   - Tone mapping (Reinhard or ACES filmic) in HLSL
   - Optional: FXAA anti-aliasing
   - Optional: Bloom (threshold → blur → composite)

4. **Screen-space ambient occlusion (SSAO)** — stretch goal
   - Render G-buffer (position, normal, depth) in geometry pass
   - SSAO compute pass reading G-buffer
   - Blur pass and composite

5. **Performance profiling & optimization**
   - GPU timestamp queries (`wgpu.QuerySet`) for per-pass timing
   - CPU-side frame breakdown (update, upload, render, present)
   - Batching: merge static meshes into one draw call where possible
   - Instance rendering (`DrawIndexed` with instance count) for repeated geometry

6. **Documentation & examples**
   - `examples/` directory with minimal working examples per feature
   - `docs/shader-guide.md` — how to write custom HLSL shaders
   - Performance guide: CPU vs. GPU comparison table
   - Updated README with screenshots, build prerequisites, and GPU feature list

### Completion Criteria
- [ ] PNG textures applied to geometry via GPU sampler
- [ ] OBJ files loadable and renderable
- [ ] Tone mapping eliminates HDR clipping
- [ ] GPU frame time measured and displayed
- [ ] README shows GPU screenshots and build instructions for all platforms

**Estimated duration:** 7–10 days

**Dependencies:** Phase 15

---

## GPU Roadmap Summary

| Phase | Goal | Est. Duration | Depends On |
|-------|------|--------------|------------|
| 9  | Window & Real-time Display | 3–5 days | MVP (8) |
| 10 | GPU Backend Abstraction | 2–3 days | Phase 9 |
| 11 | WebGPU Integration (wgpu) | 4–5 days | Phase 10 |
| 12 | GPU Geometry & Draw Pipeline | 4–5 days | Phase 11 |
| 13 | HLSL Shader Pipeline | 4–5 days | Phase 12 |
| 14 | Uniform Buffers & Per-Mesh Transforms | 3–4 days | Phase 13 |
| 15 | Lighting, Normals & Physical Shading | 4–5 days | Phase 14 |
| 16 | Advanced Features & Polish | 7–10 days | Phase 15 |
| **Total** | | **~31–42 days** | |

---

## Daily Development Checklist

**Start of day:**
- [ ] Review what you completed yesterday
- [ ] Identify today's goal (1-2 tasks from current phase)
- [ ] Set up test environment

**During development:**
- [ ] Write tests BEFORE implementation (TDD)
- [ ] Commit working code frequently
- [ ] Run tests after each change
- [ ] Take breaks when stuck

**End of day:**
- [ ] All tests passing
- [ ] Code committed to git
- [ ] Update roadmap with progress
- [ ] Note any blockers or questions

---

## MVP Success Checklist

MVP is complete when:

- ✅ CI/CD pipeline operational (Phase 0 - merged PR #5)
- ✅ Math library fully implemented and tested (Phase 1 complete)
- ✅ Can create test geometry (cube/tetrahedron) (Phase 2 complete)
- ✅ Camera system implemented (Phase 3 — complete)
- ✅ Framebuffer stores and saves images (Phase 4 — complete)
- [x] Rasterizer fills triangles correctly (Phase 5 ✅)
- [x] Shader calculates colors (Phase 6 ✅)
- [x] Render pipeline connects all components (Phase 7 ✅)
- [x] Can render colored 3D object with perspective (Phase 7 ✅)
- [x] Output image saved as PNG (Phase 7 ✅)
- [x] Tests pass (>70% coverage) (Phase 8 ✅ — all packages above threshold)
- ✅ All CI checks pass on main (Phase 0 - complete)
- [x] Demo app works (Phase 7 ✅ — `cmd/renderer/main.go` renders `output.png`)
- [x] No critical bugs (Phase 8 ✅)

**Progress: 14/14 major milestones complete (100%) — MVP DONE! 🎉**

**When all checked: MVP is done!**

---

## Estimated Timeline

**Conservative estimate (beginner):** 20-25 days
**Moderate estimate (some 3D experience):** 12-15 days
**Optimistic estimate (experienced):** 8-10 days

**Factors affecting timeline:**
- Prior 3D graphics experience
- Time available per day
- Debugging time needed
- Test-writing discipline
- Scope creep (adding non-MVP features)

**Tips to stay on track:**
1. Stick to MVP scope (resist feature creep!)
2. Write tests first (saves debugging time)
3. Commit working code frequently
4. Ask for help when stuck >1 hour
5. Celebrate completing each phase

**Most important:** Complete is better than perfect. Finish MVP first, then improve.
