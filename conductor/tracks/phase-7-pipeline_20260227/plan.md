# Implementation Plan: Phase 7 — Render Pipeline Integration

**Track ID:** phase-7-pipeline_20260227
**Spec:** [spec.md](./spec.md)
**Status:** [x] Complete
**Completed:** 2026-02-27

## Overview

Wire all packages into a complete render pipeline (`pkg/render`) and ship a working `cmd/renderer` binary using TDD.

---

## Phase 1: Pipeline & Binary

### Tasks

- [x] Task 1.1: Add `TriangleShaded(v0, v1, v2, shaderFn, fb)` to `pkg/rasterize`. Write tests first; refactor `Triangle` as thin wrapper. Commit `test:` then `feat:`.
- [x] Task 1.2: Write failing tests for `Scene`, `NewScene`, `AddMesh`. Commit `test:`.
- [x] Task 1.3: Implement `Scene` struct and constructors. Commit `feat:`.
- [x] Task 1.4: Write and pass tests for `transformVertex` — clip → perspective divide → NDC → depth → screen (Y-flipped). Commit `test:` then `feat:`.
- [x] Task 1.5: Write and pass tests for `Render(fb)` — clears fb, transforms vertices, rejects w≤0 triangles, calls `TriangleShaded`. Commit `test:` then `feat:`.
- [x] Task 1.6: Implement `cmd/renderer/main.go` — cube at origin, camera at {3,2,5}, 640×480, VertexColor shader, saves `output.png`. Commit `feat:`.

### Verification

- [x] 12 render tests pass; 100% coverage.
- [x] `go run ./cmd/renderer` produces `output.png` with coloured cube visible.
- [x] PR merged to `main`, CI green.

---

_Archived. All tasks complete._
