# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **toy 3D software renderer** implemented in Go - a learning project that demonstrates fundamental 3D graphics concepts without GPU acceleration. The renderer performs all calculations on the CPU, converting 3D geometry to 2D images through a complete graphics pipeline.

**Current Status:** Phase 0 (CI/CD) ✅ | Phase 1 (Math) ✅ | Phase 2 (Geometry) ✅ | Phase 3 (Camera) ✅ | Phase 4 (Framebuffer) ✅ | Phase 5 (Rasterization) ✅ | Phase 6 (Shading) ✅ | Ready for Phase 7 (Pipeline Integration)

**Progress Summary:**
- ✅ Phase 0 complete: CI/CD infrastructure (GitHub Actions, golangci-lint v2) merged
- ✅ Phase 1 complete: Vec3 and Mat4x4 with 100% test coverage
- ✅ Phase 2 complete: Vertex, Mesh, Tetrahedron, Cube with full test coverage
- ✅ Phase 3 complete: Camera with ViewMatrix (LookAt), ProjectionMatrix (Perspective), ViewProjectionMatrix
- ✅ Phase 4 complete: Framebuffer with color/depth buffers, depth test, PNG export (21 tests, 94.6% coverage)
- ✅ Phase 5 complete: Rasterization with barycentric triangle rasterizer (9 tests, 100% coverage)
- ✅ Phase 6 complete: Shader package with VertexColor, NewFlatColor, Depth shaders (10 tests, 100% coverage)
- 📋 13 comprehensive MVP documentation files established
- 🔄 Active development using TDD and trunk-based Git workflow

**Learning Goals:**
- Understand 3D graphics pipeline from first principles
- Implement core rendering algorithms (transformation, rasterization, shading)
- Master coordinate space transformations
- Build testable, modular graphics code

**MVP Target:** Render a simple 3D object (cube/tetrahedron) with perspective projection and save as PNG image.

**📚 For detailed vision and success criteria:** See [MVP Vision & Scope](./.claude/docs/01-mvp-vision.md)

## Development Commands

### Setup
```bash
# Initialize Go module (if not already done)
go mod init github.com/muddl/go_toy_renderer

# Download dependencies
go mod tidy
```

### Building
```bash
# Build the project
go build -o renderer ./cmd/renderer

# Build with race detector (for testing)
go build -race -o renderer ./cmd/renderer
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run a specific test
go test -run TestName ./path/to/package
```

### Running
```bash
# Run without building
go run ./cmd/renderer

# Run with arguments
go run ./cmd/renderer [args]
```

### Code Quality
```bash
# Format code
go fmt ./...

# Run linter (requires golangci-lint v1.61.0)
# Note: Project uses golangci-lint v1.x config format
# Install: https://golangci-lint.run/usage/install/
golangci-lint run

# Check for common issues
go vet ./...

# Security scanning (requires govulncheck)
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

**Important:** This project uses **golangci-lint v1.61.0**. The configuration file (`.golangci.yml`) is in v1 format. If you have v2.x installed locally, you may see configuration format errors. The CI pipeline uses v1.61.0 which will work correctly.

### Continuous Integration

The project uses GitHub Actions for automated CI/CD. The pipeline automatically runs on:
- Pushes to `main` branch
- Pull requests to `main` and feature branches
- Manual workflow dispatch

**CI Pipeline includes:**
- **Format & Validation:** Code formatting, `go vet`, dependency verification
- **Linting:** golangci-lint with 48+ linters
- **Build Matrix:** Multi-platform builds (Linux, macOS, Windows) on Go 1.24 & 1.25
- **Test Matrix:** Comprehensive testing with race detector and coverage reporting
- **Coverage Enforcement:** Fails if <70% overall or <90% for math package
- **Security Scanning:** govulncheck for vulnerability detection

**View status:**
- Check PR status checks before merging
- View workflow runs at: `https://github.com/muddl/go_toy_renderer/actions`
- CI badges displayed in README.md

**Pre-commit hook** (runs automatically on every `git commit`):

The repository ships a pre-commit hook at `.githooks/pre-commit` that runs
`go fmt`, `go vet`, and `golangci-lint` before each commit. Activate it once
per clone:

```bash
git config core.hooksPath .githooks
```

The hook runs: `go fmt ./...` → `go vet ./...` → `golangci-lint run`. If any
step fails the commit is aborted. Fix the reported issues and re-commit.

**Manual pre-push checks** (mimics full CI):
```bash
go fmt ./...
go vet ./...
golangci-lint run
go test -race -coverprofile=coverage.out ./...
govulncheck ./...
```

### Git Workflow
```bash
# One-time setup after cloning: activate the pre-commit hook
git config core.hooksPath .githooks

# Create a new feature branch
git checkout -b feature/descriptive-name

# Create a new bugfix branch
git checkout -b bugfix/issue-description

# Create a new release branch
git checkout -b release/version-number

# Check status and stage changes
git status
git add .

# Commit with descriptive message
git commit -m "type: brief description"

# Push branch to remote
git push -u origin branch-name

# Create pull request (using GitHub CLI)
gh pr create --title "Title" --body "Description"
```

## Git Development Guidelines

**CRITICAL: Always use trunk-based development with feature branches**

### Branching Strategy

**Never commit directly to `main`.** All changes must go through feature/bugfix/release branches and pull requests.

**Branch naming conventions:**
- `feature/` - New functionality (e.g., `feature/add-rasterizer`)
- `bugfix/` - Bug fixes (e.g., `bugfix/fix-matrix-multiply`)
- `release/` - Release preparation (e.g., `release/v0.1.0`)
- `hotfix/` - Critical production fixes (e.g., `hotfix/depth-buffer-crash`)

**Before making ANY code changes:**
1. Ensure you're on `main` branch: `git checkout main`
2. Pull latest changes: `git pull origin main`
3. Create a new branch: `git checkout -b feature/your-feature-name`
4. Make changes, commit frequently with clear messages
5. Push branch: `git push -u origin feature/your-feature-name`
6. Create pull request for review

### Commit Message Format

Use conventional commit format:
```
<type>: <brief description>

<optional detailed description>

<optional footer>
```

**Types:**
- `feat:` - New feature
- `fix:` - Bug fix
- `test:` - Add or update tests
- `refactor:` - Code refactoring
- `docs:` - Documentation changes
- `style:` - Code formatting (no logic change)
- `perf:` - Performance improvements
- `chore:` - Build/tooling changes

**Examples:**
```
feat: implement Vec3 cross product operation

test: add unit tests for matrix multiplication

fix: correct perspective divide by zero handling

docs: update CLAUDE.md with git workflow guidelines
```

### Pull Request Workflow

**Required for all changes:**
1. Create feature/bugfix branch from latest `main`
2. Implement changes with tests
3. Ensure all tests pass: `go test ./...`
4. Format code: `go fmt ./...`
5. Push branch to remote
6. Create pull request with:
   - Clear title describing the change
   - Description of what and why
   - Reference to any related issues
   - Test coverage summary
7. **Wait for CI checks to pass** - GitHub Actions will automatically:
   - Run all tests on multiple platforms
   - Check code formatting and linting
   - Verify test coverage meets thresholds (70% overall, 90% math)
   - Scan for security vulnerabilities
8. Address review feedback or CI failures if any
9. Merge to `main` only after all CI checks pass ✅

**Self-review checklist before creating PR:**
- All tests pass locally (`go test ./...`)
- Code is formatted (`go fmt ./...`)
- No linter warnings (`go vet ./...`)
- Coverage meets thresholds (check with `go test -cover ./...`)
- Changes match the branch purpose (feature/bugfix/etc.)
- Commit messages are clear and descriptive
- CI pipeline will verify all of the above automatically

### Branch Lifecycle

**Short-lived branches:** Aim to merge within 1-3 days to minimize divergence from `main`.

**Keep branches focused:** One branch = one feature/fix. Don't mix unrelated changes.

**Delete after merge:** Clean up merged branches to keep repository tidy.

**Sync with main regularly:** If working on long-running branch, periodically merge `main` into your branch to avoid conflicts.

## Test-Driven Development (TDD)

**CRITICAL: Always use Test-Driven Development for all code changes**

### TDD Workflow (Red-Green-Refactor)

**Every feature and bug fix must follow this cycle:**

1. **RED** - Write a failing test first
   - Write test that describes desired behavior
   - Run test to confirm it fails (proves test is valid)
   - Commit: `test: add failing test for [feature]`

2. **GREEN** - Write minimal code to pass the test
   - Implement simplest solution that makes test pass
   - Run test to confirm it passes
   - Commit: `feat/fix: implement [feature] to pass test`

3. **REFACTOR** - Improve code quality
   - Clean up implementation while keeping tests green
   - Optimize, remove duplication, improve readability
   - Commit: `refactor: improve [component] implementation`

### TDD Rules

**Never write production code without a failing test first.** The only exception is when setting up initial project structure (package declarations, empty files).

**One test at a time.** Focus on one behavior, write one test, make it pass, then move to next.

**Tests are documentation.** Test names should clearly describe what behavior is being tested.

**Test naming convention:**
```go
func TestComponentName_Behavior_ExpectedOutcome(t *testing.T) {
    // Example: TestVec3_Add_ReturnsSumOfVectors
    // Example: TestMatrix_Multiply_HandlesIdentityMatrix
    // Example: TestRasterizer_DrawTriangle_InterpolatesColors
}
```

### TDD Benefits for This Project

**Immediate feedback:** Know instantly if transformations are correct
**Confidence:** Refactor rendering pipeline without fear of breaking things
**Design quality:** TDD forces modular, testable architecture
**Regression prevention:** Catch bugs before they reach production
**Documentation:** Tests show how to use each component

### TDD Examples

**Example 1: Vector Addition**
```go
// RED - Write failing test first
func TestVec3_Add_ReturnsSumOfVectors(t *testing.T) {
    v1 := Vec3{1.0, 2.0, 3.0}
    v2 := Vec3{4.0, 5.0, 6.0}
    expected := Vec3{5.0, 7.0, 9.0}

    result := v1.Add(v2)

    if !result.Equals(expected, 0.0001) {
        t.Errorf("Add() = %v, want %v", result, expected)
    }
}

// GREEN - Implement minimal solution
func (v Vec3) Add(other Vec3) Vec3 {
    return Vec3{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

// REFACTOR - Already clean, no refactoring needed
```

**Example 2: Matrix Multiplication**
```go
// RED - Write failing test with known result
func TestMatrix_Multiply_WithIdentityMatrix(t *testing.T) {
    m := NewMatrix4x4(/* some values */)
    identity := NewIdentityMatrix()

    result := m.Multiply(identity)

    if !result.Equals(m, 0.0001) {
        t.Errorf("Multiply with identity should return original matrix")
    }
}

// GREEN - Implement multiplication
// REFACTOR - Optimize after all multiplication tests pass
```

### When to Write Tests

**Always before implementation:**
- New functions/methods
- Bug fixes (write test that reproduces bug, then fix)
- Edge cases (test boundary conditions)
- Mathematical operations (verify with known results)

**Integration tests:**
- After individual components work
- Test full pipeline stages (e.g., vertex transform → clip → project)

**Golden image tests:**
- After rendering pipeline complete
- Compare output against reference images

### Test Organization

```
pkg/math/
  vec3.go
  vec3_test.go        # Tests for Vec3
  matrix.go
  matrix_test.go      # Tests for Matrix
pkg/rasterize/
  rasterizer.go
  rasterizer_test.go  # Unit tests
  integration_test.go # Integration tests
```

**Test files alongside source:** Keep `_test.go` files next to the code they test.

**📚 For comprehensive testing strategy:** See [Test Strategy](./.claude/docs/11-test-strategy.md)

## Architecture Guidelines

### MVP Documentation

**📚 Documentation Hub:** [.claude/docs/README.md](./.claude/docs/README.md) - **Start here for complete navigation**

**Quick Links:**
- 🎯 [MVP Vision & Scope](./.claude/docs/01-mvp-vision.md) - Project goals and success criteria
- 🏗️ [Architecture Overview](./.claude/docs/02-architecture-overview.md) - Pipeline design and patterns
- 📋 [MVP Features Checklist](./.claude/docs/10-mvp-features.md) - Required vs optional features
- 🗺️ [Development Roadmap](./.claude/docs/12-development-roadmap.md) - 8-phase implementation plan (12-20 days)

**Component Specifications:**
- 🔢 [Math Component](./.claude/docs/03-math-component.md) - Vec3, Mat4x4, transformations
- 📐 [Geometry Component](./.claude/docs/04-geometry-component.md) - Vertex, Mesh, primitives
- 📷 [Camera Component](./.claude/docs/05-camera-component.md) - View & projection matrices
- 🎨 [Rasterizer Component](./.claude/docs/06-rasterizer-component.md) - Triangle rasterization
- 🖼️ [Framebuffer Component](./.claude/docs/07-framebuffer-component.md) - Pixel storage & depth testing
- 💡 [Shader Component](./.claude/docs/08-shader-component.md) - Shading system
- ⚙️ [Render Pipeline](./.claude/docs/09-render-pipeline.md) - End-to-end pipeline

### Current Project Structure

**Implemented packages:**
- `pkg/math/` - ✅ **Complete** - Vector, matrix, and mathematical operations
  - `vec3.go` - 3D vector operations (106 lines, 100% coverage)
  - `vec3_test.go` - Comprehensive Vec3 tests (464 lines, 46 tests)
  - `mat4x4.go` - 4x4 matrix operations (162 lines, 100% coverage)
  - `mat4x4_test.go` - Comprehensive Mat4x4 tests (338 lines, 28 tests)
- `pkg/geometry/` - ✅ **Complete** - Vertex, Mesh, and primitive types
  - `vertex.go` - Vertex with position + color (27 lines, 100% coverage)
  - `vertex_test.go` - Vertex tests (149 lines, 5 tests)
  - `mesh.go` - Mesh with vertex/index buffers (63 lines)
  - `mesh_test.go` - Mesh tests (246 lines, 8 tests)
  - `primitives.go` - Cube and Tetrahedron primitives (105 lines)
  - `primitives_test.go` - Primitive tests (150 lines, 8 tests)

- `pkg/camera/` - ✅ **Complete** - Camera with view and projection matrix generation
  - `camera.go` - Camera struct, ViewMatrix (LookAt), ProjectionMatrix (Perspective), ViewProjectionMatrix (55 lines)
  - `camera_test.go` - 11 tests covering all camera operations and edge cases
- `pkg/framebuffer/` - ✅ **Complete** - Framebuffer with depth testing and PNG export
  - `framebuffer.go` - Framebuffer struct, New, Clear, SetPixel (depth test), GetPixel, GetDepth, SavePNG (118 lines)
  - `framebuffer_test.go` - 21 tests covering all operations, 94.6% coverage
- `pkg/rasterize/` - ✅ **Complete** - Barycentric triangle rasterizer
  - `rasterizer.go` - ScreenVertex type, Triangle() function, edgeFunction/insideTriangle helpers
  - `rasterizer_test.go` - 9 tests covering all behaviors, 100% coverage

**Pending packages for MVP:**
- `pkg/shader/` - ✅ **Complete** - Shader package (Attributes, Func, VertexColor, NewFlatColor, Depth) (Phase 6)
- `pkg/render/` - ⏳ Core rendering pipeline and algorithms (Phase 7)
- `cmd/renderer/` - ⏳ Main application entry point (Phase 8)

**Post-MVP packages:** `pkg/scene/`, `pkg/light/`, `pkg/texture/` (not needed initially)

**Infrastructure:**
- `.github/workflows/` - ✅ CI/CD pipeline (merged)
- `.golangci.yml` - ✅ Linter configuration (golangci-lint v2)
- `.claude/docs/` - ✅ Complete MVP documentation (13 files)
- `.claude/tasks/` - ✅ Task tracking summaries

### Performance Considerations

**For MVP:** Focus on correctness and clarity first, not performance.

**Post-MVP optimizations:**
- Use efficient data structures for geometric operations (preallocate slices where possible)
- Consider using `sync.Pool` for frequently allocated objects
- Profile with `go test -bench` and `pprof` for performance-critical rendering code
- Use SIMD-friendly data layouts (struct of arrays vs array of structs)
- Consider concurrency for parallel rendering (scanlines, tiles, or ray batches)

### Testing Strategy

**MANDATORY: Test-Driven Development (TDD) for all code.**

**Test pyramid approach:**
1. **Unit tests (many)** - Math operations, individual components with known inputs/outputs
   - Write BEFORE implementation (Red-Green-Refactor cycle)
   - One test at a time, make it pass, move to next
2. **Integration tests (some)** - Transform pipeline, end-to-end rendering workflows
   - Write after unit tests pass for individual components
3. **Golden image tests (few)** - Compare rendered output against reference images for visual correctness
   - Write after rendering pipeline is complete

**Coverage goals:**
- Math package: >90% (critical foundation)
- Core packages: >80%
- Overall: >70% minimum for MVP

**Key testing practices (TDD-focused):**
- **ALWAYS write tests BEFORE implementation** (Red-Green-Refactor)
- Commit failing test first: `test: add failing test for [feature]`
- Commit passing implementation: `feat: implement [feature]`
- Use table-driven tests for edge cases and multiple scenarios
- Test with epsilon comparison for floats (±0.0001)
- Benchmark performance-critical paths after functionality works
- Create reference images for regression testing

**Test execution order:**
1. Write failing test → commit
2. Make test pass → commit
3. Refactor (if needed) → commit
4. Run full test suite: `go test ./...`
5. Check coverage: `go test -cover ./...`

See `.claude/docs/11-test-strategy.md` for detailed testing approach and TDD section above for workflow examples.

## Development Workflow

### Phase-Based Implementation

**📚 Complete 9-phase roadmap:** See [Development Roadmap](./.claude/docs/12-development-roadmap.md)

**Current Progress:** Phase 0 pending merge, Phase 1 complete ✅

**Quick phase overview:**
0. **Phase 0** (Day 0-1) - ⏳ **Pending** - CI/CD Infrastructure (GitHub Actions, linting, security scanning)
1. **Phase 1** (Days 1-3) - ✅ **Complete** - Math Foundation → [Details](./.claude/docs/03-math-component.md)
2. **Phase 2** (Days 4-5) - ✅ **Complete** - Geometry & Scene → [Details](./.claude/docs/04-geometry-component.md)
3. **Phase 3** (Days 6-7) - ✅ **Complete** - Camera System → [Details](./.claude/docs/05-camera-component.md)
4. **Phase 4** (Days 8-9) - ✅ **Complete** - Framebuffer → [Details](./.claude/docs/07-framebuffer-component.md)
5. **Phase 5** (Days 10-12) - ✅ **Complete** - Rasterization → [Details](./.claude/docs/06-rasterizer-component.md)
6. **Phase 6** (Days 13-14) - ✅ **Complete** - Shading → [Details](./.claude/docs/08-shader-component.md)
7. **Phase 7** (Days 15-18) - Pipeline Integration → [Details](./.claude/docs/09-render-pipeline.md)
8. **Phase 8** (Days 19-21) - Testing & Polish

**Daily workflow (TDD-focused):**
- Start: Pull latest `main`, create feature branch, identify first behavior to test
- RED: Write failing test for behavior → commit: `test: add failing test for X`
- GREEN: Write minimal code to pass → commit: `feat: implement X`
- REFACTOR: Clean up code (optional) → commit: `refactor: improve X`
- Repeat: Next behavior/test in same feature
- End: All tests passing (`go test ./...`), push branch, create PR
- Cleanup: Merge PR, delete branch, pull updated `main`

### Phase Completion Criteria

**CRITICAL: Before considering a phase complete, ALL of the following must be done:**

1. **Implementation Complete**
   - All code implemented and tested with TDD
   - All tests passing (`go test ./...`)
   - Code formatted (`go fmt ./...`) and linted (`go vet ./...`)

2. **Pull Request Merged**
   - Feature branch merged to `main` via PR
   - Branch deleted after merge
   - Latest `main` pulled locally

3. **Task Summary Created**
   - Task summary document created in `.claude/tasks/`
   - Includes objective, actions, decisions, results, next steps
   - Documents files modified and lessons learned
   - See [task documentation format](#task-tracking) below

4. **Documentation Updated** ✨ **NEW REQUIREMENT**
   - Review component documentation (`.claude/docs/##-component-name.md`)
     - Update implementation status (✅ complete, ⏳ pending, ❌ not needed)
     - Add API examples matching actual implementation
     - Update test coverage information
   - Update development roadmap (`.claude/docs/12-development-roadmap.md`)
     - Mark phase as ✅ **COMPLETED** with completion date
     - List all completed tasks with checkmarks
     - Note any deferred tasks and target phases
     - Update completion criteria
     - Add actual time taken vs estimate
   - Update MVP features checklist (`.claude/docs/10-mvp-features.md`)
     - Update feature status for affected components
     - Add checkmarks for completed requirements
     - Note pending items with phase numbers
   - Update documentation hub (`.claude/docs/README.md`)
     - Update "Last updated" section
     - Update "Next Steps" to reflect new current phase
   - Verify all documentation cross-references are accurate

5. **Ready for Next Phase**
   - All blocking dependencies resolved
   - Clear understanding of next phase requirements
   - No critical bugs or known issues

**Why this matters:** Keeping documentation synchronized with implementation ensures:
- Future phases have accurate reference material
- Progress tracking is up-to-date
- Other developers (or future you) understand what's been completed
- Component specifications reflect reality, not just plans

### Design Principles

**📚 Detailed architecture:** See [Architecture Overview](./.claude/docs/02-architecture-overview.md)

**Core principles:**
- **Separation of concerns:** Each pipeline stage isolated and testable
- **Data flow:** Geometry → Transform → Project → Rasterize → Shade → Framebuffer → Output
- **Simplicity first:** Optimize for readability and correctness over performance (for MVP)
- **Right-handed coordinates:** +X=Right, +Y=Up, +Z=Out (OpenGL style)
- **Column-major matrices:** Multiply on right (result = matrix × vector)

## Task Tracking

**📋 Task History:** All completed tasks are documented in `.claude/tasks/`

**Current task summaries:**
- [2025-10-10: Math Foundation Implementation](./.claude/tasks/2025-10-10_math-foundation-implementation.md) ⭐ **Latest**
- [2025-10-10: Project Setup & Documentation](./.claude/tasks/2025-10-10_project-setup-and-documentation.md)

**Task documentation format:**
- Structured sections: Objective, Actions, Decisions, Results, Next Steps
- Parsable YAML metadata for automated tracking
- Complete file modification tracking
- Lessons learned and technical notes

**When to create task summaries:**
- After completing major features or milestones
- At end of each development phase
- When significant architectural decisions are made
- After resolving complex bugs or issues

## Recent Updates

**2026-02-27 (Latest):** Phase 6 Complete - Shader Package Implemented ✅
- **Implemented `pkg/shader/` package** with `Attributes` struct, `Func` type, and three shaders
- **`VertexColor`**: pass-through shader returning interpolated vertex color unchanged
- **`NewFlatColor(color)`**: constructor returning a `shader.Func` capturing a constant color (useful for debug/testing)
- **`Depth`**: grayscale depth visualize shader (depth 0→black, 1→white)
- **`shader.Func` type**: `func(Attributes) math.Vec3` — simple function type, no interface boilerplate
- **`Attributes` struct**: holds interpolated `Color math.Vec3` and `Depth float64` per fragment
- **10 tests, 100% statement coverage** across all shader functions
- **Type/function naming** follows revive linter: `shader.Func` (not `shader.ShaderFunc`), `shader.VertexColor` (not `shader.VertexColorShader`)
- **Ready for Phase 7** (Render Pipeline Integration)

**2026-02-27:** Phase 5 Complete - Rasterization Implemented ✅
- **Implemented `rasterize.Triangle()`** (`pkg/rasterize/rasterizer.go`) using barycentric coordinates
- **ScreenVertex type** with X, Y (screen coords), Z (depth), Color (RGB Vec3)
- **edgeFunction** and **insideTriangle** helpers — winding-order agnostic (CCW and CW support)
- **Degenerate triangle guard**: area² < 1e-16 silently skips zero-area triangles
- **Bounding-box clamping**: tight box around vertices, clamped to framebuffer dimensions
- **Barycentric interpolation** for color and depth with 100% coverage across all attributes
- **9 rasterizer tests** covering: degenerate, single-pixel, axis-aligned, vertex color, centroid
  interpolation, depth ordering, off-screen, partially off-screen, CW winding
- **Function renamed** from `RasterizeTriangle` to `Triangle` (revive linter: package name provides context)
- **Ready for Phase 7** (Pipeline Integration — previously ready for Phase 6)

**2026-02-27:** Phase 4 Complete - Framebuffer Implemented ✅
- **Implemented Framebuffer type** (`pkg/framebuffer/framebuffer.go`) with Width, Height, ColorBuffer ([]Vec3), DepthBuffer ([]float64)
- **New(width, height)** allocates linear buffers; depth initialized to 1.0 (far plane)
- **Clear(color, depth)** resets all pixels in single pass
- **SetPixel(x, y, color, depth)** with strict depth test (`depth < current`); out-of-bounds silently ignored
- **GetPixel / GetDepth** with safe out-of-bounds returns (zero Vec3 / 1.0)
- **SavePNG(filename)** converts float RGB [0,1] → uint8 [0,255] with clamping via Go stdlib `image/png`
- **21 framebuffer tests** covering all operations: 94.6% package coverage
- **go.mod updated** from `go 1.25.2` → `go 1.24` / `toolchain go1.24.7` (matches available toolchain)
- **Documentation updated** across CLAUDE.md, roadmap, features checklist, README to reflect Phase 4 complete
- **Ready for Phase 5** (Rasterization: barycentric triangles, attribute interpolation)

**2026-02-27:** Phase 3 Complete - Camera System Implemented ✅
- **Implemented Camera type** (`pkg/camera/camera.go`) with Position, Target, Up, FOV, Aspect, Near, Far
- **ViewMatrix()** using LookAt algorithm: right-handed, camera looks down -Z
- **ProjectionMatrix()** using OpenGL-style perspective: near→NDC z=-1, far→NDC z=+1
- **ViewProjectionMatrix()** = Projection × View (pre-multiplied for efficiency)
- **Extended Mat4x4** with `MultiplyVec4(x,y,z,w)` for homogeneous coordinate transforms
- **11 camera tests + 3 Mat4x4 tests** covering basis orthonormality, depth mapping, FOV/aspect scaling
- **Documentation updated** across CLAUDE.md, roadmap, features checklist to reflect Phase 3 complete
- **Ready for Phase 4** (Framebuffer: color/depth buffer, PNG export)

**2026-02-27:** Phase 2 Complete - Geometry Component Implemented ✅
- **Implemented Vertex type** (`pkg/geometry/vertex.go`) with position and color, Equals with epsilon comparison, 5 tests
- **Implemented Mesh type** (`pkg/geometry/mesh.go`) with vertex/index buffers, AddVertex, AddTriangle, GetTriangle, GetTriangleVertices, ValidateIndices, 8 tests
- **Implemented Tetrahedron primitive** (4 vertices, 4 triangles, CCW winding, per-vertex colors)
- **Implemented Cube primitive** (8 vertices, 12 triangles, CCW winding for all faces)
- **Merged via PR #5** — all CI checks passing on main
- **Documentation updated** to reflect Phase 0, 1, 2 complete — ready for Phase 3 (Camera)

**2025-10-10:** Go Version Standardized to 1.25.2 🚀
- **Updated go.mod:** Changed Go directive from 1.22.5 to 1.25.2
- **Updated CI workflow:** All jobs now use Go 1.24 & 1.25 in build/test matrices
- **Updated golangci-lint config:** Changed unused linter setting from 1.23 to 1.25
- **Updated documentation:** Reflected Go 1.24 & 1.25 in CI pipeline description
- **Root cause:** Go version mismatch between local environment (1.25.2) and configurations (1.22.5/1.23)
- **Branch:** `feature/geometry-component` - resolving persistent CI lint failures
- **Impact:** CI pipeline will now run with consistent Go versions across all jobs

**2025-10-10:** CI Linting Configuration Fixed 🔧
- **Fixed deprecated linter in .golangci.yml:** Replaced `exportloopref` with `copyloopvar` (deprecated in golangci-lint v1.60+)
- **Updated Go version in linter config:** Changed from 1.22 to 1.23 to match CI environment
- **Pinned golangci-lint version in CI:** Changed from `latest` to `v1.61.0` for reproducible builds
- **Updated documentation:** Added golangci-lint version requirements to Code Quality section
- **Branch:** `feature/geometry-component` - fixes applied to resolve CI failures on PR #5
- **Impact:** CI pipeline lint job will now pass consistently

**2025-10-10:** CI/CD Infrastructure Ready for Implementation 🚀
- **CI/CD pipeline configuration files ready** (.github/workflows/, .golangci.yml)
- Implementation adds Phase 0 (CI/CD Infrastructure) as first phase
- Multi-platform testing (Linux, macOS, Windows) on Go 1.22 & 1.23
- Automated coverage enforcement (70% overall, 90% math package)
- Security scanning with govulncheck and 48+ golangci-lint linters
- Branch: `feature/gh_actions_integration` → PR #4 created
- **Status:** Resolving merge conflicts, pending merge

**2025-10-10:** Phase 1 Complete - Math Foundation Implemented ✅
- **Implemented Vec3 type with comprehensive operations (46 tests, 100% coverage)**
  - Merged via PR #1: feat: Implement Vec3 type with basic and advanced operations
  - Basic ops: Add, Sub, Mul, Div, Dot, Cross, Normalize
  - Advanced ops: Length, Distance, Lerp, Reflect, Project
- **Implemented Mat4x4 type with transformation matrices (100% coverage)**
  - Merged via PR #2: feat: implement Mat4x4 with transformation matrices
  - Core ops: NewIdentity, Multiply, Transpose
  - Transformations: Translate, Scale, RotateX/Y/Z
- **Applied TDD methodology throughout (tests written first, 3.2:1 test-to-code ratio)**
- See [task summary](./.claude/tasks/2025-10-10_math-foundation-implementation.md) for details
- Ready for Phase 2 (Geometry & Scene component) after CI/CD merge

**2025-10-10:** Task tracking system established
- Created `.claude/tasks/` directory for task summaries
- Added standardized task documentation format
- Task summaries now complement CLAUDE.md for historical tracking

**2025-10-10:** Git workflow and TDD methodology established
- **Added trunk-based Git workflow with feature branches and PR requirements**
- **Established Test-Driven Development (TDD) as mandatory practice**
- **CRITICAL:** Never commit directly to `main` - all changes via PRs

**2025-10-10:** Initial project setup and comprehensive documentation
- Created comprehensive MVP documentation (13 files in `.claude/docs/`)
- Established 8-phase development plan with ~12-20 day timeline
- Project ready for Phase 1 implementation (Math Foundation)

## Important Notes for Development

### MVP Scope (Don't Over-Engineer!)

**📚 Complete feature list:** See [MVP Features Checklist](./.claude/docs/10-mvp-features.md)

**Must have for MVP:**
- Basic 3D math (vectors, matrices, transformations)
- Simple geometry (hardcoded cube or tetrahedron)
- Camera with perspective projection
- Framebuffer with depth testing
- Triangle rasterization with color interpolation
- Simple vertex color shader
- Complete render pipeline producing PNG output

**Explicitly out of scope for MVP:**
- Textures, advanced lighting, shadows
- Transparency, post-processing
- Multi-threading, SIMD optimization
- OBJ file loading (nice to have, not required)
- Animation, materials, scene graph

**Success criteria:** Render a colored 3D object with correct perspective and depth ordering to a 640x480 PNG image.

### Common Gotchas to Avoid

**Math:**
- Matrix multiplication order (easy to reverse)
- Left vs right-handed coordinate systems
- Floating point precision in normalize operations

**Transforms:**
- Perspective divide by zero (W ≤ 0)
- Forgetting to apply perspective divide before NDC conversion
- Matrix concatenation order (projection × view × model)

**Rasterization:**
- Integer vs float coordinates (use float for interpolation)
- Depth test direction (closer < farther, not opposite)
- Color clamping (values can exceed [0,1] range)

**Pipeline:**
- Triangles behind camera (negative W) must be clipped or rejected
- Backface winding order affects culling
- Depth buffer initialization (start at 1.0 or far plane)

### Reference Documentation

**When stuck:**
1. Check relevant component doc in `.claude/docs/`
2. Review "Common Gotchas" sections
3. Verify test coverage for the failing component
4. Consult roadmap for phase completion criteria

**External resources:**
- [Scratchapixel](https://www.scratchapixel.com/) - Comprehensive 3D graphics tutorials
- [Learn OpenGL](https://learnopengl.com/) - Modern graphics pipeline concepts
- Go standard library docs for image handling
