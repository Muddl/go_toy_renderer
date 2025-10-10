# Toy 3D Renderer MVP Documentation

This directory contains comprehensive planning and architecture documentation for building a bare-bones toy 3D software renderer in Go.

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
    - Phase 0: CI/CD Infrastructure (NEW)
    - Phase 1-8: Implementation phases
    - Time estimates
    - Completion criteria per phase

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

## How to Use These Docs

### If you're starting fresh:
1. Read [01-mvp-vision.md](01-mvp-vision.md) to understand the goal
2. Skim [02-architecture-overview.md](02-architecture-overview.md) for big picture
3. Review [13-cicd-infrastructure.md](13-cicd-infrastructure.md) for CI/CD setup details
4. Follow [12-development-roadmap.md](12-development-roadmap.md) starting with Phase 0 (CI/CD)
5. Reference component docs as you implement each part

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
- Basic geometry (cube/tetrahedron)
- Camera system
- Framebuffer with depth test
- Triangle rasterization
- Simple shader
- Complete render pipeline

### Success Criteria
✓ Renders 3D object to image file
✓ Perspective looks correct
✓ Depth ordering works
✓ Colors interpolate smoothly
✓ Core tests pass

### Timeline Estimate
- Phase 0 (CI/CD): 0.5-1 day
- Phases 1-8 (Development): 12-20 days
- Total: ~13-21 days

---

## Next Steps

**Ready to start coding?**

1. Set up Go project structure (see CLAUDE.md in repo root)
2. **Begin with Phase 0: CI/CD Infrastructure** (ensure quality gates from day 1)
3. Continue with Phase 1: Math Foundation ✅ **COMPLETED**
4. Write tests first (TDD)
5. Commit frequently
6. Refer back to these docs as needed

**Questions or stuck?**
- Review the relevant component doc
- Check "Common Gotchas" sections
- Verify against test strategy
- Remember: MVP is minimal, not perfect

---

## Document Maintenance

These docs are written for MVP scope only. As the project evolves:

- MVP docs should remain stable (reference for basics)
- Create new docs for post-MVP features
- Update roadmap with actual vs. estimated timeline
- Keep architecture doc current with major changes

**Last updated:** 2025-10-10
- Added comprehensive CI/CD Infrastructure documentation (doc 13)
- Consolidated CI/CD docs from .github/ into .claude/docs/
- Added Phase 0 (CI/CD Infrastructure) to roadmap
- Updated feature status tracking (Phase 1 complete)
- Added CI/CD requirements to MVP features

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
