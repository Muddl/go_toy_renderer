# Specification: Phase 6 — Shader Package

**Track ID:** phase-6-shading_20260227
**Type:** Feature
**Status:** Archived (Complete)
**Completed:** 2026-02-27

## Summary

Implement the shader package (`pkg/shader`) with a `Func` type and three built-in shaders: `VertexColor` (pass-through), `NewFlatColor` (constant colour), and `Depth` (grayscale depth visualiser).

## Context

The shader package decouples fragment colouring logic from the rasterizer. `shader.Func` is a plain function type (`func(Attributes) Vec3`) with no interface overhead. `Attributes` carries interpolated colour and depth per fragment.

## Acceptance Criteria

- [x] `Attributes` struct: `Color Vec3`, `Depth float64`
- [x] `Func` type: `func(Attributes) math.Vec3`
- [x] `VertexColor` — returns interpolated `attr.Color` unchanged
- [x] `NewFlatColor(color Vec3) Func` — returns a closure capturing the constant colour
- [x] `Depth` — maps depth 0→black, 1→white (grayscale)
- [x] Naming follows revive linter (`shader.Func` not `shader.ShaderFunc`)
- [x] 10 tests; 100% statement coverage

## Source Reference

- Task summary: [`.claude/tasks/2026-02-27_phase-6-shading-implementation.md`](../../../../.claude/tasks/2026-02-27_phase-6-shading-implementation.md)

---

_Archived. This track is a historical record of completed work._
