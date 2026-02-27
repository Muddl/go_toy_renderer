# MVP Feature Requirements

## Infrastructure (Must Have First)

These infrastructure features are REQUIRED before development begins. They ensure code quality and prevent regressions.

### 0. CI/CD Pipeline ✓
**What:** Automated testing, building, and quality enforcement

**Requirements:**
- GitHub Actions workflow with multi-platform builds (Linux, macOS, Windows)
- Automated test execution with race detector on all PRs
- Code coverage enforcement (70% overall, 90% math package)
- Linting with golangci-lint (48+ linters)
- Security scanning with govulncheck
- Status badges in README

**Success criteria:** All PRs automatically validated for tests, coverage, quality, and security before merge.

**Status:** ✅ **Complete** - Merged via PR #5 (2026-02-27)

---

## Core Features (Must Have)

These features are REQUIRED for MVP completion. Without these, the renderer is not functional.

### 1. Basic 3D Math ✅
**What:** Vector and matrix operations

**Requirements:**
- ✅ Vec3 add, subtract, scale, dot, cross, normalize, distance (46 tests, 100% coverage)
- ✅ Mat4x4 multiply, identity, zero, transpose, element access (100% coverage)
- ✅ Transform creation: translation, rotation (X/Y/Z), scale (all tested)
- ✅ LookAt matrix for camera — `Mat4x4` via `camera.ViewMatrix()`
- ✅ Perspective projection matrix — `Mat4x4` via `camera.ProjectionMatrix()`
- ✅ `MultiplyVec4(x,y,z,w)` for homogeneous transforms — `pkg/math/mat4x4.go`
- ⏳ Viewport transformation (Phase 5 - Rasterizer, converts NDC → screen pixels)

**Status:** ✅ Phase 1 complete (2025-10-10), Phase 3 camera matrices implemented (2026-02-27).

**Success criteria:**
- ✅ All basic math operations pass unit tests with known results
- ✅ Camera transformation matrices implemented and tested

---

### 2. Simple Geometry ✅
**What:** Basic 3D shapes

**Requirements:**
- ✅ Vertex type (position + color) — `pkg/geometry/vertex.go`
- ✅ Mesh type (vertex buffer + index buffer) — `pkg/geometry/mesh.go`
- ✅ Tetrahedron: 4 vertices, 4 triangles — `pkg/geometry/primitives.go`
- ✅ Cube: 8 vertices, 12 triangles — `pkg/geometry/primitives.go`

**Success criteria:** Can create a mesh with valid triangles.

**Status:** ✅ **Complete** (2026-02-27) — merged via PR #5

---

### 3. Camera System ✅
**What:** Viewpoint control

**Requirements:**
- ✅ Camera with position, target, up vector — `pkg/camera/camera.go`
- ✅ FOV, aspect ratio, near/far planes — `Camera` struct fields
- ✅ Generate view matrix — `ViewMatrix()` using LookAt algorithm
- ✅ Generate projection matrix — `ProjectionMatrix()` OpenGL-style perspective
- ✅ Combine into view-projection matrix — `ViewProjectionMatrix()` = Projection × View

**Success criteria:** Camera transformations produce correct screen-space coordinates.

**Status:** ✅ **Complete** (2026-02-27)

---

### 4. Framebuffer ✅
**What:** Pixel storage with depth testing

**Requirements:**
- ✅ 2D color buffer (RGB per pixel) — `ColorBuffer []math.Vec3`
- ✅ 2D depth buffer (Z per pixel) — `DepthBuffer []float64`
- ✅ SetPixel with depth test (closer wins) — strict `depth < current`
- ✅ Clear to background color — `Clear(color, depth)`
- ✅ Export to PNG file — `SavePNG(filename)` via stdlib `image/png`
- ✅ GetPixel and GetDepth for testing and debugging

**Success criteria:** Can write pixels and save valid image file.

**Status:** ✅ **Complete** (Phase 4 — 2026-02-27) — 21 tests, 94.6% coverage

---

### 5. Triangle Rasterization ✅
**What:** Convert triangles to pixels

**Requirements:**
- ✅ `ScreenVertex` type with X, Y (screen coords), Z (depth), Color (RGB Vec3) — `pkg/rasterize/rasterizer.go`
- ✅ `rasterize.Triangle(v0, v1, v2, fb)` using barycentric coordinates
- ✅ Barycentric color interpolation across triangle
- ✅ Depth interpolation for depth testing
- ✅ Degenerate triangle guard (zero-area triangles silently skipped)
- ✅ CW and CCW winding support (winding-order agnostic)
- ✅ Out-of-bounds pixel safety (clamped bounding box + silent ignore via framebuffer)

**Success criteria:** Rasterized triangle fills correct pixels with interpolated colors.

**Status:** ✅ **Complete** (Phase 5 — 2026-02-27) — 9 tests, 100% coverage

---

### 6. Basic Shader ⏳
**What:** Per-pixel color calculation

**Requirements:**
- Shader function that takes interpolated attributes
- Returns final RGB color
- Minimum: vertex color pass-through shader

**Success criteria:** Shader receives interpolated attributes and produces colors.

**Status:** ⏳ **Pending** (Phase 6)

---

### 7. Rendering Pipeline ⏳
**What:** Orchestrate all components

**Requirements:**
- Transform vertices (local → world → view → clip → screen)
- Process all triangles in mesh
- Rasterize each triangle
- Shade each pixel
- Write to framebuffer with depth test
- Save final image

**Success criteria:** Complete render of a simple mesh produces valid output image.

**Status:** ⏳ **Pending** (Phase 7)

---

## Nice-to-Have Features (Optional for MVP)

These features improve the renderer but are NOT required for initial MVP.

### Backface Culling
**Benefit:** Performance - skip triangles facing away
**Complexity:** Low
**Priority:** Medium

### OBJ File Loading
**Benefit:** Test with real models, not just hardcoded geometry
**Complexity:** Medium
**Priority:** Medium

### Simple Diffuse Lighting
**Benefit:** More realistic appearance
**Complexity:** Medium (requires normals)
**Priority:** Low-Medium

### Multiple Meshes in Scene
**Benefit:** Render more complex scenes
**Complexity:** Low
**Priority:** Medium

### Wireframe Mode
**Benefit:** Debugging, visualization
**Complexity:** Low
**Priority:** Low

---

## Explicitly Out of Scope (Post-MVP)

These features are valuable but too complex for initial MVP:

### ❌ Textures
- Requires UV coordinates
- Texture loading and sampling
- Filtering (bilinear, trilinear)

### ❌ Advanced Lighting
- Multiple lights
- Specular highlights
- Shadows

### ❌ Transparency
- Alpha blending
- Render order sorting

### ❌ Post-Processing
- Bloom, tone mapping
- Anti-aliasing (MSAA, FXAA)

### ❌ Optimization
- Multi-threading
- SIMD vectorization
- Tile-based rendering

### ❌ Advanced Features
- Normal mapping
- Animation/skinning
- Particle systems
- Scene graph hierarchy

---

## Feature Validation Checklist

Use this checklist to verify MVP is complete:

- [x] CI/CD pipeline operational with all quality gates (Phase 0 ✅)
- [x] Math library has all required operations (Phase 1 basic ops ✅, camera ops in Phase 3)
- [x] Can create at least one test mesh (cube/tetrahedron) (Phase 2 ✅)
- [x] Camera can be positioned and oriented (Phase 3 ✅)
- [x] Framebuffer can store and export pixels (Phase 4 ✅)
- [x] Triangle rasterizer fills correct pixels (Phase 5 ✅)
- [x] Colors interpolate smoothly across triangles (Phase 5 ✅)
- [x] Depth test prevents Z-fighting (Phase 4 ✅)
- [ ] Complete pipeline renders mesh to image file (Phase 7)
- [ ] Output image visually shows 3D object with perspective (Phase 7)
- [x] All core components have basic tests (math: 100%, geometry: ~95%+)
- [x] All CI checks pass on main branch (Phase 0 ✅)

**Progress: 10/12 items complete (83%)**

---

## MVP Demo Target

**Goal:** Render this scene and save as image:

```
Scene:
- Single colored cube (or tetrahedron)
- Each face/vertex has different color
- Camera positioned to see multiple faces
- Perspective projection shows depth
- Output: 640x480 PNG image

Result should show:
✓ Correct 3D perspective (far is smaller)
✓ Smooth color interpolation
✓ Faces in correct depth order (no Z-fighting)
✓ Clean edges (no major artifacts)
```

**When this works, MVP is complete.**

---

## Success Metrics

**MVP is successful when:**

1. **It renders** - Produces a viewable image
2. **It's correct** - Perspective and depth look right
3. **It's tested** - Core functionality has tests
4. **It's documented** - Code is reasonably clear
5. **It's complete** - All 7 core features implemented

**MVP is NOT successful if:**
- No image output
- Severe visual artifacts (completely broken)
- Missing any of the 7 core features
- No tests for math/core functions
