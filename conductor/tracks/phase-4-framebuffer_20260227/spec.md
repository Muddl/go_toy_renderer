# Specification: Phase 4 — Framebuffer

**Track ID:** phase-4-framebuffer_20260227
**Type:** Feature
**Status:** Archived (Complete)
**Completed:** 2026-02-27

## Summary

Implement the framebuffer package (`pkg/framebuffer`) — a 2D pixel store with a colour buffer (`[]Vec3`), a depth buffer (`[]float64`), a depth test on `SetPixel`, and PNG export via the Go standard library.

## Context

The framebuffer is the rendering target. `SetPixel` only updates a pixel if the incoming depth value is strictly less than the stored value, implementing per-pixel depth testing without a Z-buffer library. `SavePNG` converts float RGB [0, 1] to uint8 [0, 255] with clamping.

## Acceptance Criteria

- [x] `Framebuffer` struct: Width, Height, ColorBuffer ([]Vec3), DepthBuffer ([]float64)
- [x] `New(width, height)` — allocates linear buffers; depth initialised to 1.0 (far plane)
- [x] `Clear(color, depth)` — resets all pixels in a single pass
- [x] `SetPixel(x, y, color, depth)` — depth test `depth < current`; out-of-bounds silently ignored
- [x] `GetPixel / GetDepth` — safe out-of-bounds returns (zero Vec3 / 1.0)
- [x] `SavePNG(filename)` — float→uint8 with clamping; stdlib `image/png`
- [x] 21 tests; 94.6% package coverage

## Source Reference

- Task summary: [`.claude/tasks/2026-02-27_framebuffer-implementation.md`](../../../../.claude/tasks/2026-02-27_framebuffer-implementation.md)

---

_Archived. This track is a historical record of completed work._
