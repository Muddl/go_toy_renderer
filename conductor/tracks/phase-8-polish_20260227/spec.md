# Specification: Phase 8 — Testing & Polish

**Track ID:** phase-8-polish_20260227
**Type:** Chore
**Status:** Archived (Complete)
**Completed:** 2026-02-27

## Summary

Complete the MVP with a golden image regression test, a full integration test suite, performance benchmarks, and a comprehensive README — confirming the CPU renderer is production-ready for the project's learning goals.

## Context

Phase 8 is the MVP completion milestone. All core functionality was proven in Phases 1–7; Phase 8 locks it down with pixel-exact regression testing, quantified performance baselines, and documentation for future contributors and phases.

## Acceptance Criteria

- [x] `TestRender_GoldenImage_Triangle` in `pkg/render/integration_test.go` — renders deterministic RGB triangle, byte-exact comparison against `testdata/golden_triangle.png`; `-update` flag regenerates reference
- [x] 7 integration tests covering all 3 shaders, multiple camera positions, multiple meshes, depth shader grayscale
- [x] Benchmarks in `pkg/render/bench_test.go` and `pkg/rasterize/rasterizer_bench_test.go`:
  - Full 640×480 cube render: ~1.2 ms/frame
  - Vertex MVP transform: ~58 ns/vertex
  - Edge function: ~0.38 ns
- [x] `README.md` fully written: quick start, API example, shader table, build/test/benchmark commands, performance table, coordinate system notes
- [x] All documentation updated; roadmap Phase 8 marked COMPLETED (14/14 milestones)

## Source Reference

- Task summary: [`.claude/tasks/2026-02-27_phase-8-testing-polish.md`](../../../../.claude/tasks/2026-02-27_phase-8-testing-polish.md)

---

_Archived. This track is a historical record of completed work._
