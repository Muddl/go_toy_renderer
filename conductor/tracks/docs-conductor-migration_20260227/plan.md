# Implementation Plan: Migrate .claude/ Context into Conductor Artifacts

**Track ID:** docs-conductor-migration_20260227
**Spec:** [spec.md](./spec.md)
**Created:** 2026-02-27
**Status:** [~] In Progress

## Overview

Systematically migrate 16 `.claude/docs/` files and 7 `.claude/tasks/` summaries into Conductor artifacts. Work in 6 phases: audit, enrich core artifacts, create architecture doc, archive completed phases (0–8), create pending GPU tracks (9–16), and trim CLAUDE.md. Historical `.claude/` files are preserved — never deleted.

---

## Phase 1: Audit & Mapping

Read every `.claude/docs/` and `.claude/tasks/` file; produce a migration mapping table at `conductor/migration-audit.md`.

### Tasks

- [x] Task 1.1: Read `01-mvp-vision.md` → map sections to `product.md` and `product-guidelines.md`.
- [x] Task 1.2: Read `02-architecture-overview.md` → map to new `conductor/architecture.md`.
- [x] Task 1.3: Read component docs `03`–`09` → identify content for `architecture.md` (API, coverage) vs. archived tracks.
- [x] Task 1.4: Read `10-mvp-features.md` and `11-test-strategy.md` → map to `workflow.md` and `architecture.md`.
- [x] Task 1.5: Read `12-development-roadmap.md` fully → extract per-phase task lists for archived + pending tracks.
- [x] Task 1.6: Read GPU docs `13-cicd-infrastructure.md`, `14-gpu-backend.md`, `15-hlsl-shader-pipeline.md`, `16-realtime-display.md` → map to `tech-stack.md` and pending tracks 9–16.
- [x] Task 1.7: Read all 7 `.claude/tasks/` summaries → note key decisions, deferred items, and lessons for archived tracks.
- [x] Task 1.8: Write `conductor/migration-audit.md` — two-column table: source file → target Conductor artifact(s); flag any content with no clear home.

### Verification

- [x] `migration-audit.md` exists and covers all 16 docs + 7 tasks.
- [x] Every source file maps to ≥1 target artifact.
- [x] Content with no clear home is explicitly flagged.

---

## Phase 2: Enrich Core Conductor Artifacts

Update the 5 artifacts created during setup with richer content sourced from the audit.

### Tasks

- [x] Task 2.1: Enrich `conductor/product.md` — add architecture summary, pipeline data-flow diagram (ASCII from `02`), coordinate system and matrix convention.
- [x] Task 2.2: Enrich `conductor/tech-stack.md` — add package dependency table (from `02`), CI matrix details (from `13`), GPU planned deps table (from `14`/`15`/`16`).
- [x] Task 2.3: Enrich `conductor/workflow.md` — add test naming conventions, coverage thresholds (>90% math, >80% core, >70% overall), phase completion checklist, golden image test guidance.
- [x] Task 2.4: Enrich `conductor/product-guidelines.md` — add "Common Gotchas" section (math, transforms, rasterization, pipeline).

### Verification

- [x] Each core artifact reviewed against its source docs — no key content missing.
- [x] `conductor/tech-stack.md` GPU table matches `14`/`15`/`16`.
- [x] `conductor/workflow.md` coverage thresholds correct.

---

## Phase 3: Create `conductor/architecture.md`

Consolidate the 7 component spec docs (03–09) and architecture overview into a single reference document.

### Tasks

- [ ] Task 3.1: Write `conductor/architecture.md` with sections:
  - Pipeline overview (ASCII diagram from `02`)
  - Package dependency rules (one-direction flow)
  - Coordinate system & matrix conventions
  - Per-component API summary table (pkg, key types/funcs, test coverage) for all 8 MVP packages
  - Performance baselines (from Phase 8 benchmarks)
  - GPU architecture additions (Phase 9+ package layout from `02` and `14`)
- [ ] Task 3.2: Link `conductor/architecture.md` from `conductor/index.md`.

### Verification

- [ ] All 8 MVP packages have an entry in the API summary table.
- [ ] Pipeline diagram matches current implementation.
- [ ] GPU package layout section matches `14-gpu-backend.md`.

---

## Phase 4: Archived Tracks for Completed Phases (0–8)

Create minimal Conductor tracks for Phases 0–8 with status `archived`. These are historical records, not active work items.

### Tasks

- [ ] Task 4.1: Create `conductor/tracks/phase-0-cicd_20260227/` — spec summarising CI/CD goals; plan listing completed tasks; metadata `status: archived`.
- [ ] Task 4.2: Create `conductor/tracks/phase-1-math_20260227/` — link to `.claude/tasks/2025-10-10_math-foundation-implementation.md`.
- [ ] Task 4.3: Create `conductor/tracks/phase-2-geometry_20260227/`.
- [ ] Task 4.4: Create `conductor/tracks/phase-3-camera_20260227/`.
- [ ] Task 4.5: Create `conductor/tracks/phase-4-framebuffer_20260227/` — link to `.claude/tasks/2026-02-27_framebuffer-implementation.md`.
- [ ] Task 4.6: Create `conductor/tracks/phase-5-rasterizer_20260227/` — link to `.claude/tasks/2026-02-27_rasterizer-implementation.md`.
- [ ] Task 4.7: Create `conductor/tracks/phase-6-shading_20260227/` — link to `.claude/tasks/2026-02-27_phase-6-shading-implementation.md`.
- [ ] Task 4.8: Create `conductor/tracks/phase-7-pipeline_20260227/` — link to `.claude/tasks/2026-02-27_phase-7-pipeline-integration.md`.
- [ ] Task 4.9: Create `conductor/tracks/phase-8-polish_20260227/` — link to `.claude/tasks/2026-02-27_phase-8-testing-polish.md`.
- [ ] Task 4.10: Register all 9 archived tracks in `conductor/tracks.md`.

### Verification

- [ ] 9 archived track directories exist.
- [ ] Each has `spec.md`, `plan.md`, `metadata.json`, `index.md`.
- [ ] All 9 registered in `conductor/tracks.md` with status `archived`.
- [ ] Each links to its `.claude/tasks/` source where one exists.

---

## Phase 5: Pending Tracks for GPU Roadmap (9–16)

Create full Conductor tracks (spec + plan) for each GPU phase.

### Tasks

- [ ] Task 5.1: Create `conductor/tracks/phase-9-window_20260227/` — full spec + plan from roadmap Phase 9 (GLFW, event loop, CPU blit, 60 fps) + `16-realtime-display.md`.
- [ ] Task 5.2: Create `conductor/tracks/phase-10-gpu-abstraction_20260227/` — `Renderer` interface, CPU/GPU backends from roadmap + `14`.
- [ ] Task 5.3: Create `conductor/tracks/phase-11-webgpu_20260227/` — wgpu-native integration from roadmap + `14`.
- [ ] Task 5.4: Create `conductor/tracks/phase-12-gpu-geometry_20260227/`.
- [ ] Task 5.5: Create `conductor/tracks/phase-13-hlsl_20260227/` — HLSL → WGSL pipeline from roadmap + `15`.
- [ ] Task 5.6: Create `conductor/tracks/phase-14-uniforms_20260227/`.
- [ ] Task 5.7: Create `conductor/tracks/phase-15-lighting_20260227/`.
- [ ] Task 5.8: Create `conductor/tracks/phase-16-advanced_20260227/`.
- [ ] Task 5.9: Register all 8 pending tracks in `conductor/tracks.md`.

### Verification

- [ ] 8 pending track directories exist with spec + plan.
- [ ] Phase 9 plan is detailed enough to start `/conductor:implement`.
- [ ] All 8 registered in `conductor/tracks.md` with status `pending`.

---

## Phase 6: Trim CLAUDE.md

Reduce CLAUDE.md to a lean pointer document; move substantive content to Conductor artifacts.

### Tasks

- [ ] Task 6.1: Remove or condense sections now fully covered by Conductor:
  - "Project Overview" → 1-paragraph summary + link to `conductor/product.md`
  - TDD section → link to `conductor/workflow.md`
  - Git workflow section → link to `conductor/workflow.md`
  - Architecture guidelines → link to `conductor/architecture.md`
  - Phase-based implementation → link to `conductor/tracks.md`
- [ ] Task 6.2: Retain in CLAUDE.md (not duplicated in Conductor):
  - Development commands (build, test, run, lint)
  - "Common Gotchas" block (quick-reference for Claude context injection)
  - External resource links
  - "Recent Updates" changelog (last 3 entries)
- [ ] Task 6.3: Add "Context Sources" section at top of CLAUDE.md pointing to `conductor/index.md` as the canonical reference.

### Verification

- [ ] CLAUDE.md is ≤150 lines after trimming.
- [ ] All retained sections provide value not covered in Conductor.
- [ ] No broken cross-references in CLAUDE.md.

---

## Final Verification

- [ ] All acceptance criteria in `spec.md` are met.
- [ ] `conductor/migration-audit.md` confirms no significant content lost.
- [ ] `conductor/tracks.md` lists 9 archived + 8 pending + 1 this chore track (18 total).
- [ ] `conductor/index.md` links to `architecture.md` and shows active tracks.
- [ ] All tests still pass: `go test ./...` (no Go source code changed).
- [ ] CLAUDE.md ≤150 lines and references Conductor throughout.

---

_Generated by Conductor. Tasks marked [~] when in progress, [x] when complete._
