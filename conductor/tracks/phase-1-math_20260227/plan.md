# Implementation Plan: Phase 1 — Math Foundation

**Track ID:** phase-1-math_20260227
**Spec:** [spec.md](./spec.md)
**Status:** [x] Complete
**Completed:** 2025-10-10

## Overview

Implement `pkg/math` using strict TDD (Red-Green-Refactor). Vec3 first, then Mat4x4 with transformation matrices.

---

## Phase 1: Vec3 and Mat4x4

### Tasks

- [x] Task 1.1: Write failing tests for Vec3 basic ops (Add, Sub, Mul, Div, Dot, Cross, Normalize, Length). Commit `test:`.
- [x] Task 1.2: Implement Vec3 basic ops to pass tests. Commit `feat:`.
- [x] Task 1.3: Write and pass tests for Vec3 advanced ops (Distance, Lerp, Reflect, Project) and epsilon equals. Commit `test:` then `feat:`.
- [x] Task 1.4: Write failing tests for Mat4x4 (NewIdentity, Multiply, Transpose, MultiplyVec4). Commit `test:`.
- [x] Task 1.5: Implement Mat4x4 core ops. Commit `feat:`.
- [x] Task 1.6: Write and pass tests for transformation constructors (Translate, Scale, RotateX/Y/Z). Commit `test:` then `feat:`.

### Verification

- [x] `go test -cover ./pkg/math/...` ≥ 90%.
- [x] All 74 tests pass (46 Vec3 + 28 Mat4x4).
- [x] PR merged to `main`, CI green.

---

_Archived. All tasks complete._
