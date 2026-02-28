# Implementation Plan: Phase 6 — Shader Package

**Track ID:** phase-6-shading_20260227
**Spec:** [spec.md](./spec.md)
**Status:** [x] Complete
**Completed:** 2026-02-27

## Overview

Implement `pkg/shader` with a function-type shader interface and three built-in shaders using TDD.

---

## Phase 1: Shader Package

### Tasks

- [x] Task 1.1: Write failing tests for `Attributes` struct and `Func` type. Commit `test:`.
- [x] Task 1.2: Define `Attributes`, `Func` type in `shader.go`. Commit `feat:`.
- [x] Task 1.3: Write and pass test for `VertexColor` — returns interpolated colour unchanged. Commit `test:` then `feat:`.
- [x] Task 1.4: Write and pass test for `NewFlatColor(c)` — closure captures constant colour. Commit `test:` then `feat:`.
- [x] Task 1.5: Write and pass test for `Depth` — 0.0→black, 1.0→white, 0.5→mid-grey. Commit `test:` then `feat:`.

### Verification

- [x] 10 shader tests pass; 100% statement coverage.
- [x] revive linter: no `shader.ShaderFunc` or similar stutter.
- [x] PR merged to `main`, CI green.

---

_Archived. All tasks complete._
