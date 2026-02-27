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

## Phase 8: Testing & Polish (Days 19-21)

**Goal:** Ensure quality and reliability.

### Tasks
1. Write integration tests
   - Test complete render pipeline
   - Test various camera positions
   - Test multiple meshes

2. Create golden image tests
   - Render reference images
   - Compare output against references
   - Set up automated comparison

3. Improve test coverage
   - Aim for >80% overall
   - Focus on uncovered critical paths

4. Create demo application
   - cmd/renderer/main.go
   - Render cube from nice angle
   - Save as output.png

5. Write basic documentation
   - README with usage instructions
   - Code comments for public APIs
   - Example usage

6. Performance check
   - Benchmark critical paths
   - Identify bottlenecks (don't optimize yet)
   - Document performance characteristics

### Completion Criteria
- [ ] All tests passing
- [ ] >70% test coverage
- [ ] At least one golden image test
- [ ] Demo app produces nice output image
- [ ] README explains how to use
- [ ] Code is reasonably documented

**When complete:** MVP is ready to show!

---

## Beyond MVP (Future Phases)

### Phase 9: Enhancements (Optional)
- OBJ file loading
- Backface culling
- Multiple meshes in scene
- Wireframe rendering mode

### Phase 10: Lighting (Optional)
- Add normals to vertices
- Implement simple diffuse lighting
- Directional light source

### Phase 11: Optimization (Optional)
- Multi-threading (parallel triangle processing)
- SIMD vectorization
- Profile-guided optimization

### Phase 12: Advanced Features (Optional)
- Texture mapping
- Specular highlights
- Simple shadows
- Post-processing effects

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
- [ ] Tests pass (>70% coverage) (Phase 8 — already >70%, integration/golden images pending)
- ✅ All CI checks pass on main (Phase 0 - complete)
- [x] Demo app works (Phase 7 ✅ — `cmd/renderer/main.go` renders `output.png`)
- [ ] No critical bugs (Phase 8)

**Progress: 11/14 major milestones complete (79%)**

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
