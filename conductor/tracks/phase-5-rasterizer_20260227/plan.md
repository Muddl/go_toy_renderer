# Implementation Plan: Phase 5 — Rasterization

**Track ID:** phase-5-rasterizer_20260227
**Spec:** [spec.md](./spec.md)
**Status:** [x] Complete
**Completed:** 2026-02-27

## Overview

Implement `pkg/rasterize` with a barycentric-coordinate triangle rasterizer using TDD.

---

## Phase 1: Rasterizer Implementation

### Tasks

- [x] Task 1.1: Define `ScreenVertex` type and write failing test for degenerate triangle (zero area silently skipped). Commit `test:`.
- [x] Task 1.2: Implement `Triangle(v0, v1, v2, fb)` with bounding-box scan and barycentric test. Commit `feat:`.
- [x] Task 1.3: Write and pass tests for colour interpolation (vertex colour, centroid interpolation). Commit `test:` then `feat:`.
- [x] Task 1.4: Write and pass tests for depth interpolation and depth ordering (nearer triangle wins). Commit `test:` then `feat:`.
- [x] Task 1.5: Write and pass tests for edge cases: off-screen triangle, partially off-screen, CW winding. Commit `test:` then `feat:`.

### Verification

- [x] 13 rasterizer tests pass; 100% coverage.
- [x] PR merged to `main`, CI green.

---

_Note: `TriangleShaded` variant added in Phase 7 (Pipeline Integration)._

_Archived. All tasks complete._
