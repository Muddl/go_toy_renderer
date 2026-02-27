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

## Phase 3: Camera System ✅ **IMPLEMENTED** (2026-02-27)

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

**Status:** ✅ **IMPLEMENTED** (2026-02-27) — CI validation in progress

**Next:** Phase 4 - Framebuffer

---

## Phase 4: Framebuffer (Days 8-9)

**Goal:** Store rendered pixels and output images.

### Tasks
1. Implement Framebuffer type
   - Color buffer (RGB per pixel)
   - Depth buffer (Z per pixel)

2. Implement Clear method
   - Reset to background color and depth

3. Implement SetPixel with depth test
   - Check if new pixel is closer
   - Update color and depth if closer

4. Implement GetPixel (for testing)
   - Retrieve color at coordinates

5. Implement image export
   - Convert float RGB to byte RGB
   - Save as PNG using Go image libraries
   - Test file output

### Completion Criteria
- [ ] Can create framebuffer of any size
- [ ] Clear works correctly
- [ ] Depth test prevents farther pixels from overwriting closer ones
- [ ] Can save valid PNG file
- [ ] Framebuffer tests pass

**When complete:** You can store and save rendered images.

---

## Phase 5: Rasterization (Days 10-12)

**Goal:** Convert triangles to pixels with interpolation.

### Tasks
1. Choose rasterization algorithm
   - Recommended: Barycentric (simpler code)
   - Alternative: Scanline (more classic)

2. Implement barycentric coordinate calculation
   - Test with known points inside/outside triangle

3. Implement triangle bounding box
   - Find min/max X and Y coordinates

4. Implement rasterization loop
   - For each pixel in bounding box
   - Check if inside triangle
   - Calculate barycentric weights

5. Implement attribute interpolation
   - Interpolate depth (Z)
   - Interpolate color (RGB)
   - Test interpolation correctness

6. Connect to framebuffer
   - For each pixel covered
   - Call framebuffer.SetPixel with depth test

7. Test rasterization
   - Single pixel triangle
   - Axis-aligned triangles
   - Arbitrary triangles
   - Verify color interpolation

### Completion Criteria
- [ ] Can rasterize any triangle to pixels
- [ ] Attribute interpolation is correct at vertices
- [ ] Depth interpolation works
- [ ] No crashes on degenerate triangles
- [ ] Rasterizer tests pass

**When complete:** You can draw triangles to framebuffer.

---

## Phase 6: Shading (Days 13-14)

**Goal:** Calculate per-pixel colors.

### Tasks
1. Define Shader interface or function type
   - Input: interpolated attributes
   - Output: final color

2. Implement vertex color shader
   - Simply return interpolated color
   - Test with known inputs

3. Implement flat color shader (for testing)
   - Return constant color
   - Use for debugging

4. Optional: Implement depth visualization shader
   - Convert depth to grayscale
   - Useful for debugging

5. Test shaders
   - Verify correct colors returned
   - Ensure called once per pixel

### Completion Criteria
- [ ] Shader interface defined
- [ ] At least one shader implemented (vertex color)
- [ ] Shader receives correct interpolated data
- [ ] Shader tests pass

**When complete:** You can control per-pixel appearance.

---

## Phase 7: Render Pipeline Integration (Days 15-18)

**Goal:** Connect all components into working renderer.

### Tasks
1. Implement vertex transformation stage
   - For each vertex in mesh:
   - Apply model matrix (identity for now)
   - Apply view-projection matrix
   - Perform perspective divide
   - Convert to screen coordinates

2. Implement primitive assembly
   - Group indices into triangles (every 3 indices)
   - Fetch transformed vertices for each triangle

3. Connect rasterization and shading
   - For each triangle:
   - Rasterize to pixels
   - For each pixel, call shader
   - Write to framebuffer

4. Implement main Render function
   - Clear framebuffer
   - Transform all vertices
   - Process all triangles
   - Return/save framebuffer

5. Test end-to-end
   - Render tetrahedron
   - Render cube
   - Save output image
   - Visually verify correctness

6. Debug visual issues
   - Check matrix multiplication order
   - Verify depth test direction
   - Check winding order
   - Fix any artifacts

### Completion Criteria
- [ ] Can render complete mesh to image
- [ ] Perspective looks correct
- [ ] Depth ordering is correct
- [ ] Colors interpolate smoothly
- [ ] Output image matches expectations
- [ ] No crashes on valid input

**When complete:** MVP is functionally complete!

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
- ✅ Camera system implemented (Phase 3 — CI validation in progress)
- [ ] Framebuffer stores and saves images (Phase 4)
- [ ] Rasterizer fills triangles correctly (Phase 5)
- [ ] Shader calculates colors (Phase 6)
- [ ] Render pipeline connects all components (Phase 7)
- [ ] Can render colored 3D object with perspective (Phase 7)
- [ ] Output image saved as PNG (Phase 7)
- [ ] Tests pass (>70% coverage) (Phase 8)
- ✅ All CI checks pass on main (Phase 0 - complete)
- [ ] Demo app works (Phase 8)
- [ ] No critical bugs (Phase 8)

**Progress: 5/14 major milestones complete (35.7%)**

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
