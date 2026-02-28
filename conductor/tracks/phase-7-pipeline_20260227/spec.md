# Specification: Phase 7 — Render Pipeline Integration

**Track ID:** phase-7-pipeline_20260227
**Type:** Feature
**Status:** Archived (Complete)
**Completed:** 2026-02-27

## Summary

Wire all packages into an end-to-end render pipeline (`pkg/render`) and a working binary (`cmd/renderer`) that renders a coloured cube to a 640×480 PNG.

## Context

Phase 7 is the integration milestone: `Scene`, `Render()`, and `transformVertex()` connect camera, geometry, rasterizer, shader, and framebuffer into a single pipeline. The binary `cmd/renderer/main.go` serves as the acceptance test for the full MVP rendering stack.

## Acceptance Criteria

- [x] `Scene` struct: meshes []geometry.Mesh, Camera, Shader shader.Func
- [x] `NewScene(camera, shader)` constructor
- [x] `AddMesh(mesh)` appends to scene
- [x] `Render(fb)` clears fb, transforms each vertex via `cam.ViewProjectionMatrix()`, rejects w≤0 triangles, calls `rasterize.TriangleShaded`
- [x] `transformVertex` — clip → perspective divide → NDC → depth `(ndcZ+1)/2` → screen (Y-flipped)
- [x] `cmd/renderer/main.go` — cube at origin, camera at {3,2,5}, 640×480, VertexColor shader, saves `output.png`
- [x] 12 tests; 100% coverage

## Source Reference

- Task summary: [`.claude/tasks/2026-02-27_phase-7-pipeline-integration.md`](../../../../.claude/tasks/2026-02-27_phase-7-pipeline-integration.md)

---

_Archived. This track is a historical record of completed work._
