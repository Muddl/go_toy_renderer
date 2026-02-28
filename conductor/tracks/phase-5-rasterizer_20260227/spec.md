# Specification: Phase 5 — Rasterization

**Track ID:** phase-5-rasterizer_20260227
**Type:** Feature
**Status:** Archived (Complete)
**Completed:** 2026-02-27

## Summary

Implement the rasterizer package (`pkg/rasterize`) with a barycentric-coordinate triangle rasterizer that interpolates per-vertex colour and depth, and a degenerate triangle guard.

## Context

The rasterizer converts screen-space triangles into per-pixel colour + depth values written to the framebuffer. Barycentric interpolation enables smooth per-vertex colour and correct depth interpolation. The `TriangleShaded` variant (added in Phase 7) accepts a shader function for extensible fragment processing.

## Acceptance Criteria

- [x] `ScreenVertex` type: X, Y (screen coords), Z (depth), Color (Vec3)
- [x] `Triangle(v0, v1, v2, fb)` — rasterizes with `shader.VertexColor`; degenerate guard (area² < 1e-16)
- [x] `TriangleShaded(v0, v1, v2, shaderFn, fb)` — per-pixel shader callback
- [x] Bounding-box clamping to framebuffer dimensions
- [x] Winding-order agnostic (CCW and CW supported)
- [x] 13 tests; 100% coverage

## Source Reference

- Task summary: [`.claude/tasks/2026-02-27_rasterizer-implementation.md`](../../../../.claude/tasks/2026-02-27_rasterizer-implementation.md)

---

_Archived. This track is a historical record of completed work._
