# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **toy 3D software renderer** implemented in Go - a learning project that demonstrates fundamental 3D graphics concepts without GPU acceleration. The renderer performs all calculations on the CPU, converting 3D geometry to 2D images through a complete graphics pipeline.

**Current Status:** Project is in initial stages with comprehensive MVP documentation completed.

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

# Run linter (requires golangci-lint)
golangci-lint run

# Check for common issues
go vet ./...
```

### Git Workflow
```bash
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
7. Address review feedback if any
8. Merge to `main` only after approval (if working with team) or all checks pass

**Self-review checklist before creating PR:**
- All tests pass
- Code is formatted (`go fmt`)
- No linter warnings (`go vet`)
- Changes match the branch purpose (feature/bugfix/etc.)
- Commit messages are clear and descriptive

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

### Recommended Project Structure

Organize code into these packages for MVP:

- `cmd/renderer/` - Main application entry point
- `pkg/math/` - Vector, matrix, and mathematical operations (Vec3, Mat4x4, transformations)
- `pkg/geometry/` - Mesh, vertex, and primitive types (Vertex, Triangle, Mesh)
- `pkg/camera/` - Camera and view transformation logic (LookAt, Perspective)
- `pkg/render/` - Core rendering pipeline and algorithms (vertex transform, primitive assembly)
- `pkg/rasterize/` - Triangle rasterization with attribute interpolation
- `pkg/shader/` - Shader implementations (vertex color, flat color)
- `pkg/framebuffer/` - Framebuffer with depth testing and image output
- `internal/` - Internal implementation details

**Post-MVP packages:** `pkg/scene/`, `pkg/light/`, `pkg/texture/` (not needed initially)

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

**📚 Complete 8-phase roadmap:** See [Development Roadmap](./.claude/docs/12-development-roadmap.md)

**Quick phase overview:**
1. **Phase 1** (Days 1-3) - Math Foundation → [Details](./.claude/docs/03-math-component.md)
2. **Phase 2** (Days 4-5) - Geometry & Scene → [Details](./.claude/docs/04-geometry-component.md)
3. **Phase 3** (Days 6-7) - Camera System → [Details](./.claude/docs/05-camera-component.md)
4. **Phase 4** (Days 8-9) - Framebuffer → [Details](./.claude/docs/07-framebuffer-component.md)
5. **Phase 5** (Days 10-12) - Rasterization → [Details](./.claude/docs/06-rasterizer-component.md)
6. **Phase 6** (Days 13-14) - Shading → [Details](./.claude/docs/08-shader-component.md)
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

**2025-10-10 (Latest):** Task tracking system established
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
