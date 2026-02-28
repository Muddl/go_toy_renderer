# Specification: Phase 3 — Camera System

**Track ID:** phase-3-camera_20260227
**Type:** Feature
**Status:** Archived (Complete)
**Completed:** 2026-02-27

## Summary

Implement the camera package (`pkg/camera`) providing a right-handed LookAt view matrix and an OpenGL-style perspective projection matrix, with a convenience `ViewProjectionMatrix()` pre-multiply.

## Context

The camera transforms 3D world-space vertices into 2D clip-space coordinates. The `ViewProjectionMatrix()` is applied per-vertex in `pkg/render`, and the conventions established here (right-handed, column-major, NDC z ∈ [−1, +1]) must match the rasterizer and depth buffer.

## Acceptance Criteria

- [x] `Camera` struct: Position, Target, Up, FOV, Aspect, Near, Far
- [x] `ViewMatrix()` — right-handed LookAt; camera looks down −Z
- [x] `ProjectionMatrix()` — OpenGL-style perspective; near→−1, far→+1 in NDC z
- [x] `ViewProjectionMatrix()` — returns `Projection × View`
- [x] `Mat4x4.MultiplyVec4(x,y,z,w)` added to `pkg/math`
- [x] 11 camera tests + 3 Mat4x4 tests covering basis orthonormality, depth mapping, FOV/aspect scaling

## Source Reference

- No `.claude/tasks/` summary for this phase; content derived from `12-development-roadmap.md`.

---

_Archived. This track is a historical record of completed work._
