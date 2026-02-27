# Toy 3D Renderer Documentation

This directory contains comprehensive planning and architecture documentation for the go_toy_renderer project — a CPU software renderer (MVP complete) evolving toward cross-platform GPU acceleration with real-time rendering.

**Purpose:** These docs define WHAT to build and WHY, not HOW to implement it. Use them as a guide during development.

---

## Documentation Overview

### Getting Started

1. **[MVP Vision & Scope](01-mvp-vision.md)** - Start here
   - What we're building
   - Why this project exists
   - Success criteria
   - Constraints and non-goals

2. **[Architecture Overview](02-architecture-overview.md)**
   - High-level design philosophy
   - Rendering pipeline stages
   - Package organization
   - Key design decisions

---

### Core Components

Each component doc explains purpose, requirements, and testing approach:

3. **[Math Component](03-math-component.md)**
   - Vectors, matrices, transformations
   - Critical foundation for everything else

4. **[Geometry Component](04-geometry-component.md)**
   - Vertices, triangles, meshes
   - Hardcoded primitives

5. **[Camera Component](05-camera-component.md)**
   - Viewpoint control
   - View and projection matrices

6. **[Rasterizer Component](06-rasterizer-component.md)**
   - Triangle to pixel conversion
   - Attribute interpolation

7. **[Framebuffer Component](07-framebuffer-component.md)**
   - Pixel storage with depth testing
   - Image output

8. **[Shader Component](08-shader-component.md)**
   - Per-pixel color calculation
   - Simple shading models

9. **[Render Pipeline](09-render-pipeline.md)**
   - Orchestrating all components
   - End-to-end data flow

---

### Planning & Execution

10. **[MVP Features](10-mvp-features.md)**
    - CI/CD infrastructure requirements (Phase 0)
    - Required features checklist
    - Nice-to-have features
    - Out-of-scope items
    - Success validation

11. **[Test Strategy](11-test-strategy.md)**
    - Unit, integration, and golden image tests
    - Coverage goals
    - Testing best practices

12. **[Development Roadmap](12-development-roadmap.md)**
    - Phase 0: CI/CD Infrastructure ✅
    - Phase 1–8: MVP implementation phases ✅
    - Phase 9–16: GPU acceleration roadmap 📋
    - Time estimates and completion criteria per phase

---

### Infrastructure & Tooling

13. **[CI/CD Infrastructure](13-cicd-infrastructure.md)**
    - GitHub Actions pipeline architecture
    - Automated testing and security scanning
    - Linter configuration (48+ linters)
    - Coverage enforcement (70%/90%)
    - Badge setup and troubleshooting
    - Integration with TDD workflow

---

### GPU Acceleration (Post-MVP)

14. **[GPU Backend](14-gpu-backend.md)**
    - WebGPU / wgpu-native architecture
    - Cross-platform backend selection (D3D12 / Metal / Vulkan)
    - Device, swapchain, render pass design

15. **[HLSL Shader Pipeline](15-hlsl-shader-pipeline.md)**
    - HLSL authoring workflow
    - Compilation: HLSL → WGSL via naga/DXC
    - Uniform buffers and bind groups
    - Built-in and custom shader guide

16. **[Real-time Display](16-realtime-display.md)**
    - Window creation and event loop
    - Frame timing, vsync, camera controls
    - CPU-blit bridge and GPU swap chain

---

## How to Use These Docs

### If you're starting fresh:
1. Read [01-mvp-vision.md](01-mvp-vision.md) to understand the original goal
2. Skim [02-architecture-overview.md](02-architecture-overview.md) for big picture (includes GPU pipeline)
3. Follow [12-development-roadmap.md](12-development-roadmap.md) — MVP phases 0–8 are complete; GPU phases 9–16 are next
4. Reference component docs as you implement each part

### If you're mid-development:
- Use component docs as reference for what to implement
- Check [10-mvp-features.md](10-mvp-features.md) to verify you're on track
- Use [11-test-strategy.md](11-test-strategy.md) to guide testing

### If you're debugging:
- Check component docs for common gotchas
- Review architecture doc for data flow
- Verify test coverage with test strategy doc

---

## Documentation Principles

These docs are designed to be:

- **Concise** - No unnecessary detail
- **Actionable** - Clear what to build
- **Scoped** - Focused on MVP only
- **Practical** - Real examples and gotchas
- **Testable** - Clear success criteria

**Note:** These are planning docs, not tutorials. They describe WHAT and WHY, not step-by-step HOW.

---

## Quick Reference

### Must-Have Infrastructure (Phase 0)
- GitHub Actions CI/CD pipeline
- Multi-platform builds (Linux, macOS, Windows)
- Automated testing with coverage enforcement
- Security scanning (govulncheck)
- Linting (golangci-lint)

### Must-Have Features (MVP)
- 3D math (vectors, matrices) ✅
- Basic geometry (cube/tetrahedron) ✅
- Camera system ✅
- Framebuffer with depth test ✅ (Phase 4 — complete)
- Triangle rasterization ✅
- Simple shader ✅
- Complete render pipeline ✅ (Phase 7 — complete)

### Success Criteria
✓ Renders 3D object to image file
✓ Perspective looks correct
✓ Depth ordering works
✓ Colors interpolate smoothly
✓ Core tests pass

### Timeline Estimate
- Phase 0 (CI/CD): 0.5–1 day ✅
- Phases 1–8 (MVP): 12–20 days ✅ **DONE**
- Phases 9–16 (GPU): ~31–42 days 📋 planned

---

## Current Status

### MVP (Phases 0–8) — ✅ COMPLETE (2026-02-27)

| Phase | Description | Status |
|-------|-------------|--------|
| 0 | CI/CD Infrastructure | ✅ Complete |
| 1 | Math Foundation (Vec3, Mat4x4) | ✅ Complete |
| 2 | Geometry & Scene (Mesh, Cube, Tetrahedron) | ✅ Complete |
| 3 | Camera System (LookAt, Perspective) | ✅ Complete |
| 4 | Framebuffer (color/depth, PNG export) | ✅ Complete |
| 5 | Triangle Rasterization (barycentric) | ✅ Complete |
| 6 | Shading (VertexColor, FlatColor, Depth) | ✅ Complete |
| 7 | Render Pipeline Integration | ✅ Complete |
| 8 | Testing & Polish (golden image, benchmarks) | ✅ Complete |

### GPU Roadmap (Phases 9–16) — 📋 Planned

| Phase | Description | Status |
|-------|-------------|--------|
| 9  | Window & Real-time Display | 📋 Planned |
| 10 | GPU Backend Abstraction Layer | 📋 Planned |
| 11 | WebGPU Integration (wgpu-native) | 📋 Planned |
| 12 | GPU Geometry & Draw Pipeline | 📋 Planned |
| 13 | HLSL Shader Pipeline | 📋 Planned |
| 14 | Uniform Buffers & Per-Mesh Transforms | 📋 Planned |
| 15 | Lighting, Normals & Physical Shading | 📋 Planned |
| 16 | Advanced Features & Polish | 📋 Planned |

**Next step:** Begin Phase 9 — Window & Real-time Display

**Questions or stuck?**
- Review the relevant component doc
- Check "Common Gotchas" in CLAUDE.md
- Verify against test strategy
- See GPU docs in files 14, 15, 16

---

## Document Maintenance

These docs are written for MVP scope only. As the project evolves:

- MVP docs should remain stable (reference for basics)
- Create new docs for post-MVP features
- Update roadmap with actual vs. estimated timeline
- Keep architecture doc current with major changes

**Last updated:** 2026-02-27
- GPU roadmap (Phases 9–16) added with WebGPU/HLSL north star
- Architecture doc updated with GPU pipeline diagram and backend table
- New GPU component docs: 14-gpu-backend.md, 15-hlsl-shader-pipeline.md, 16-realtime-display.md
- Roadmap "Beyond MVP" placeholder replaced with detailed Phase 9–16 specs
- MVP features doc updated: viewport transform marked complete, out-of-scope section revised
- Depth buffer convention clarified (0 = near, 1 = far)
- Phase 8 (Testing & Polish) completed — MVP is fully done!

---

## Contributing to Docs

If you find errors or have suggestions:
- Document should be concise (not comprehensive)
- Focus on clarity over completeness
- Include practical examples
- Maintain consistent formatting
- Update README if adding new docs

---

Happy rendering! 🎨
