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

## Architecture Guidelines

### MVP Documentation

Comprehensive planning documentation is available in `.claude/docs/`:

- **[MVP Vision & Scope](/.claude/docs/01-mvp-vision.md)** - Project goals and success criteria
- **[Architecture Overview](/.claude/docs/02-architecture-overview.md)** - Pipeline design and patterns
- **[Core Components](/.claude/docs/)** - Detailed specs for math, geometry, camera, rasterizer, framebuffer, shader, and pipeline
- **[MVP Features](/.claude/docs/10-mvp-features.md)** - Required vs optional features checklist
- **[Test Strategy](/.claude/docs/11-test-strategy.md)** - Unit, integration, and golden image testing approach
- **[Development Roadmap](/.claude/docs/12-development-roadmap.md)** - 8-phase implementation plan (est. 12-20 days)

**Start here:** Read `.claude/docs/README.md` for navigation guide.

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

**Test pyramid approach:**
1. **Unit tests (many)** - Math operations, individual components with known inputs/outputs
2. **Integration tests (some)** - Transform pipeline, end-to-end rendering workflows
3. **Golden image tests (few)** - Compare rendered output against reference images for visual correctness

**Coverage goals:**
- Math package: >90% (critical foundation)
- Core packages: >80%
- Overall: >70% minimum for MVP

**Key testing practices:**
- Write tests BEFORE implementation (TDD)
- Use table-driven tests for edge cases
- Test with epsilon comparison for floats (±0.0001)
- Benchmark performance-critical paths
- Create reference images for regression testing

See `.claude/docs/11-test-strategy.md` for detailed testing approach.

## Development Workflow

### Phase-Based Implementation (Recommended)

Follow the 8-phase roadmap in `.claude/docs/12-development-roadmap.md`:

1. **Phase 1: Math Foundation** (Days 1-3) - Vector3, Matrix4x4, transformations
2. **Phase 2: Geometry & Scene** (Days 4-5) - Vertex, Mesh, primitives
3. **Phase 3: Camera System** (Days 6-7) - View and projection matrices
4. **Phase 4: Framebuffer** (Days 8-9) - Pixel storage with depth test
5. **Phase 5: Rasterization** (Days 10-12) - Triangle to pixels with interpolation
6. **Phase 6: Shading** (Days 13-14) - Per-pixel color calculation
7. **Phase 7: Pipeline Integration** (Days 15-18) - Connect all components
8. **Phase 8: Testing & Polish** (Days 19-21) - Tests, docs, demo app

**Daily workflow:**
- Start: Review yesterday's progress, identify 1-2 tasks
- During: Write tests first, commit frequently, run tests after changes
- End: All tests passing, code committed, update progress

### Design Principles

**Separation of concerns:** Each pipeline stage isolated and testable
**Data flow:** Geometry → Transform → Project → Rasterize → Shade → Framebuffer → Output
**Simplicity first:** Optimize for readability and correctness over performance (for MVP)
**Right-handed coordinates:** +X=Right, +Y=Up, +Z=Out (OpenGL style)
**Column-major matrices:** Multiply on right (result = matrix × vector)

## Recent Updates

**2025-10-10:** Initial project setup
- Created comprehensive MVP documentation (13 files in `.claude/docs/`)
- Defined architecture, components, features, test strategy, and roadmap
- Established 8-phase development plan with ~12-20 day timeline
- Project ready for Phase 1 implementation (Math Foundation)

## Important Notes for Development

### MVP Scope (Don't Over-Engineer!)

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
