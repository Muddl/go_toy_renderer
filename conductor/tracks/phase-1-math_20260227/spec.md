# Specification: Phase 1 — Math Foundation

**Track ID:** phase-1-math_20260227
**Type:** Feature
**Status:** Archived (Complete)
**Completed:** 2025-10-10

## Summary

Implement the core 3D math library (`pkg/math`) — `Vec3` for 3D vector operations and `Mat4x4` for 4×4 matrix operations including transformation matrices — with 100% test coverage using strict TDD.

## Context

The math package is the critical foundation that all other packages depend on. Every vertex, camera, and transform operation flows through `Vec3` and `Mat4x4`. Getting this right with full test coverage is the precondition for every subsequent phase.

## Acceptance Criteria

- [x] `Vec3` type with Add, Sub, Mul, Div, Dot, Cross, Normalize, Length, Distance, Lerp, Reflect, Project
- [x] `Mat4x4` type with NewIdentity, Multiply, Transpose, MultiplyVec4
- [x] Transformation constructors: Translate, Scale, RotateX, RotateY, RotateZ
- [x] Epsilon comparison (`±0.0001`) for all float operations
- [x] 100% test coverage (46 Vec3 tests + 28 Mat4x4 tests)
- [x] All tests written before implementation (strict TDD)

## Source Reference

- Task summary: [`.claude/tasks/2025-10-10_math-foundation-implementation.md`](../../../../.claude/tasks/2025-10-10_math-foundation-implementation.md)

---

_Archived. This track is a historical record of completed work._
