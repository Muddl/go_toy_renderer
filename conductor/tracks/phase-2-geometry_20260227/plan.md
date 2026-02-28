# Implementation Plan: Phase 2 — Geometry Component

**Track ID:** phase-2-geometry_20260227
**Spec:** [spec.md](./spec.md)
**Status:** [x] Complete
**Completed:** 2026-02-27

## Overview

Implement `pkg/geometry` with `Vertex`, `Mesh`, and primitive constructors using TDD.

---

## Phase 1: Vertex, Mesh, Primitives

### Tasks

- [x] Task 1.1: Write and pass tests for `Vertex` type (Position, Color, Equals with epsilon). Commit `test:` then `feat:`.
- [x] Task 1.2: Write and pass tests for `Mesh` type — AddVertex, AddTriangle, GetTriangle, GetTriangleVertices, ValidateIndices. Commit `test:` then `feat:`.
- [x] Task 1.3: Write and pass tests for `NewTetrahedron()` — 4 vertices, 4 triangles, CCW winding, per-vertex colours. Commit `test:` then `feat:`.
- [x] Task 1.4: Write and pass tests for `NewCube()` — 8 vertices, 12 triangles, correct winding for all 6 faces. Commit `test:` then `feat:`.
- [x] Task 1.5: Run full test suite, check coverage, open PR, merge to `main`.

### Verification

- [x] All 21 geometry tests pass.
- [x] PR merged to `main`, CI green.

---

_Archived. All tasks complete._
