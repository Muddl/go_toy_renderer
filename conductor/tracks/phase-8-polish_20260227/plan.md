# Implementation Plan: Phase 8 — Testing & Polish

**Track ID:** phase-8-polish_20260227
**Spec:** [spec.md](./spec.md)
**Status:** [x] Complete
**Completed:** 2026-02-27

## Overview

Complete the MVP with a golden image regression test, integration test suite, performance benchmarks, and comprehensive README.

---

## Phase 1: Polish & Documentation

### Tasks

- [x] Task 1.1: Write `TestRender_GoldenImage_Triangle` in `pkg/render/integration_test.go` — deterministic RGB triangle, byte-exact comparison against `testdata/golden_triangle.png`; `-update` flag to regenerate. Commit `test:`.
- [x] Task 1.2: Write 7 integration tests covering all 3 shaders, multiple camera positions, multiple meshes, depth shader grayscale. Commit `test:`.
- [x] Task 1.3: Write benchmarks in `pkg/render/bench_test.go` (`BenchmarkRender_Cube_640x480`) and `pkg/rasterize/rasterizer_bench_test.go` (`BenchmarkTriangle_Rasterize_Large`, `BenchmarkEdgeFunction`). Commit `test:`.
- [x] Task 1.4: Write comprehensive `README.md` — quick start, API example, shader table, build/test/benchmark commands, performance table, coordinate system notes. Commit `docs:`.
- [x] Task 1.5: Update all documentation (roadmap Phase 8 COMPLETED, features checklist, docs hub). Commit `docs:`.
- [x] Task 1.6: Merge PR to `main`; confirm CI green with golden image test passing on all 3 OS.

### Verification

- [x] All benchmarks produce stable numbers (~1.2 ms/frame, ~58 ns/vertex, ~0.38 ns edge fn).
- [x] Golden image test byte-exact on all platforms.
- [x] PR merged to `main`, CI green.
- [x] MVP declared complete.

---

_Archived. All tasks complete. MVP (Phases 0–8) is done._
