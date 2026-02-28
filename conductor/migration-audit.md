# Migration Audit — .claude/ → conductor/

**Date:** 2026-02-27
**Track:** docs-conductor-migration_20260227
**Status:** Phase 1 complete — all files mapped

---

## Source → Target Mapping

### .claude/docs/ files

| Source File | Primary Target Artifact(s) | Notes |
|-------------|---------------------------|-------|
| `01-mvp-vision.md` | `conductor/product.md` | Description, problem, learning goals, success criteria — all present; "What Makes This Bare Bones" section not duplicated (not needed) |
| `02-architecture-overview.md` | `conductor/architecture.md` (new) | Pipeline diagram, coordinate system, matrix convention, depth convention (0=near,1=far), package layout, GPU pipeline diagram, backend table, shader authoring flow, post-GPU package layout |
| `03-math-component.md` | `conductor/architecture.md` — math section | API table (Vec3, Mat4x4, transforms), coverage (46 tests, 100%), gotchas |
| `04-geometry-component.md` | `conductor/architecture.md` — geometry section | API table (Vertex, Mesh, NewCube, NewTetrahedron), gotchas |
| `05-camera-component.md` | `conductor/architecture.md` — camera section | API table, conventions (looks down -Z, vertical FOV), gotchas |
| `06-rasterizer-component.md` | `conductor/architecture.md` — rasterizer section | Actual API (ScreenVertex, Triangle, TriangleShaded), barycentric algorithm, 9 tests 100%, gotchas |
| `07-framebuffer-component.md` | `conductor/architecture.md` — framebuffer section | Actual API (New, Clear, SetPixel, GetPixel, GetDepth, SavePNG), 21 tests 94.6%, gotchas |
| `08-shader-component.md` | `conductor/architecture.md` — shader section | Actual API (Attributes, Func, VertexColor, NewFlatColor, Depth), 10 tests 100%, gotchas |
| `09-render-pipeline.md` | `conductor/architecture.md` — render pipeline section | Actual API (Scene, NewScene, AddMesh, Render, transformVertex, TriangleShaded), coordinate transform formula, 12 tests 100% |
| `10-mvp-features.md` | `conductor/product.md` (success criteria) + `conductor/architecture.md` (coverage table) | Feature validation checklist → coverage table in architecture.md; GPU roadmap features → pending track specs |
| `11-test-strategy.md` | `conductor/workflow.md` (test naming, thresholds, golden image) + `conductor/architecture.md` (test pyramid) | Coverage thresholds (math >90%, core >80%, overall >70%), golden image pattern, table-driven test examples |
| `12-development-roadmap.md` | **Phases 0–8** → 9 archived Conductor tracks; **Phases 9–16** → 8 pending tracks | Per-phase task lists, completion criteria, time estimates all captured in track spec/plan files |
| `13-cicd-infrastructure.md` | `conductor/tech-stack.md` (CI matrix, linter config, tooling) | 6 CI jobs, build matrix (3 OS × 2 Go), security pipeline, linter list, troubleshooting |
| `14-gpu-backend.md` | `conductor/tech-stack.md` (GPU deps, platform table) + `conductor/architecture.md` (GPU architecture) + phase-11/phase-12 track specs | wgpu-native, device init, swap chain, geometry buffers, depth buffer, testing strategy (GPU_TESTS env var) |
| `15-hlsl-shader-pipeline.md` | `conductor/architecture.md` (shader authoring section) + phase-13 track spec | HLSL → WGSL via naga, directory layout, built-in shader sources, uniform buffer binding layout, shader loader pattern |
| `16-realtime-display.md` | `conductor/architecture.md` (window/display section) + phase-9 track spec | GLFW package layout, window lifecycle, frame loop, CPU blit path, camera controller, input bindings, resize handling, common gotchas |
| `README.md` (docs hub) | `conductor/index.md` | Navigation superseded; original remains for historical reference |

### .claude/tasks/ files

| Task Summary | Target Archived Track | Notes |
|-------------|----------------------|-------|
| `2025-10-10_project-setup-and-documentation.md` | `phase-0-cicd_20260227` (linked) | Documents initial project scaffolding and doc creation; linked from Phase 0 archived track |
| `2025-10-10_math-foundation-implementation.md` | `phase-1-math_20260227` (primary) | Vec3 (46 tests), Mat4x4 (100% coverage), PR #1 + #2 |
| `2026-02-27_framebuffer-implementation.md` | `phase-4-framebuffer_20260227` (primary) | 21 tests, 94.6%, go.mod Go version fix |
| `2026-02-27_rasterizer-implementation.md` | `phase-5-rasterizer_20260227` (primary) | 9 tests, 100%, barycentric algorithm decisions |
| `2026-02-27_phase-6-shading-implementation.md` | `phase-6-shading_20260227` (primary) | 10 tests, 100%, Func naming decision (revive linter) |
| `2026-02-27_phase-7-pipeline-integration.md` | `phase-7-pipeline_20260227` (primary) | 12 tests, 100%, TriangleShaded, transformVertex formula |
| `2026-02-27_phase-8-testing-polish.md` | `phase-8-polish_20260227` (primary) | Golden image test, 7 integration tests, benchmarks, README |

> **Note:** Phases 2 (Geometry) and 3 (Camera) have no dedicated `.claude/tasks/` summary files. The roadmap doc (`12-development-roadmap.md`) contains sufficient task detail for their archived tracks.

---

## Flagged Content — No Clear Single Home

| Content | Source(s) | Resolution |
|---------|-----------|-----------|
| **"Common Gotchas" per component** | `03`–`09`, `09-render-pipeline.md` | **Dual-place:** Consolidated per-package gotchas go in `conductor/architecture.md`; also retained in CLAUDE.md as a quick-reference "cheat sheet" section (only the 6-bullet summary from CLAUDE.md, not the full per-component lists) |
| **"Future Enhancements" per component** (MSAA, perspective-correct interp, lighting models, etc.) | `06`, `07`, `08`, `09` | **Route to GPU track specs:** Framebuffer enhancements (MSAA, HDR) → phase-16; Perspective-correct interpolation → phase-12; Phong/PBR shading → phase-15 |
| **Pre-implementation "Conceptual API" code** | `03`, `04`, `05` (early docs) | **Discard:** These contain pseudocode that predates the actual implementation. The actual API is documented in `06`–`09` and in the roadmap. Do not migrate pseudocode — only verified actual API. |
| **`.claude/docs/README.md`** (docs hub) | n/a | **Preserved in-place** as historical record; superseded by `conductor/index.md`. Not migrated. |
| **CI badge markdown** | `13-cicd-infrastructure.md` | **Stay in README.md** — CI badges belong in the project README, not in conductor. |
| **Phase 2 & 3 task summaries** | (missing) | **Not a problem:** Roadmap Phase 2/3 sections contain complete task lists. Archived tracks will be built from roadmap content only. |

---

## Content Verified Correctly Present in conductor/ Already

After audit, the following content from `.claude/` is already well-represented in the conductor artifacts created during setup:

| Content Area | Already In |
|-------------|-----------|
| TDD Red-Green-Refactor workflow | `conductor/workflow.md` |
| Conventional Commits format | `conductor/workflow.md` |
| Branch naming conventions | `conductor/workflow.md` |
| PR policy & self-review checklist | `conductor/workflow.md` |
| Phase completion criteria (5-step checklist) | `conductor/workflow.md` |
| Go 1.24 + stdlib-only dependency statement | `conductor/tech-stack.md` |
| Planned Phase 9–11 deps (GLFW, wgpu-native, naga) | `conductor/tech-stack.md` |
| CI matrix (Linux/macOS/Windows, Go 1.24/1.25) | `conductor/tech-stack.md` |
| golangci-lint v2 rules (30+ linters) | `conductor/code_styleguides/go.md` |
| Shell scripting conventions | `conductor/code_styleguides/shell.md` |
| Design principles (correctness first, TDD, right-handed coords) | `conductor/product-guidelines.md` |

---

## Migration Completeness Assessment

| Category | Files | Status |
|----------|-------|--------|
| `.claude/docs/` mapped | 16/16 | ✅ All mapped |
| `.claude/tasks/` mapped | 7/7 | ✅ All mapped |
| Flagged content resolved | 5 items | ✅ All resolved |
| Content with no target | 0 items | ✅ None |
| Information loss risk | Low | ✅ Pre-impl pseudocode intentionally excluded |

**Conclusion:** Ready to proceed to Phase 2 (Enrich Core Artifacts) and Phase 3 (architecture.md).
