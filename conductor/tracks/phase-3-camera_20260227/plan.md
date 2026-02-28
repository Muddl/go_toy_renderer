# Implementation Plan: Phase 3 — Camera System

**Track ID:** phase-3-camera_20260227
**Spec:** [spec.md](./spec.md)
**Status:** [x] Complete
**Completed:** 2026-02-27

## Overview

Implement `pkg/camera` with right-handed LookAt and OpenGL-style perspective projection using TDD. Also extend `pkg/math` with `MultiplyVec4`.

---

## Phase 1: Camera Implementation

### Tasks

- [x] Task 1.1: Extend `Mat4x4` with `MultiplyVec4(x,y,z,w float64) (float64, float64, float64, float64)`. Write test first. Commit `test:` then `feat:`.
- [x] Task 1.2: Write failing tests for `Camera` struct and `ViewMatrix()` — verify right vector, up vector, and forward vector are orthonormal. Commit `test:`.
- [x] Task 1.3: Implement `ViewMatrix()` using LookAt algorithm. Commit `feat:`.
- [x] Task 1.4: Write and pass tests for `ProjectionMatrix()` — verify near plane maps to NDC z=−1 and far plane to +1. Commit `test:` then `feat:`.
- [x] Task 1.5: Write and pass test for `ViewProjectionMatrix()` = Projection × View. Commit `test:` then `feat:`.

### Verification

- [x] 11 camera tests + 3 Mat4x4 tests pass.
- [x] PR merged to `main`, CI green.

---

_Archived. All tasks complete._
