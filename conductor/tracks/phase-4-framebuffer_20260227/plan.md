# Implementation Plan: Phase 4 — Framebuffer

**Track ID:** phase-4-framebuffer_20260227
**Spec:** [spec.md](./spec.md)
**Status:** [x] Complete
**Completed:** 2026-02-27

## Overview

Implement `pkg/framebuffer` with colour and depth buffers, depth-tested pixel writes, and PNG export using TDD.

---

## Phase 1: Framebuffer Implementation

### Tasks

- [x] Task 1.1: Write failing tests for `New`, `Clear`, `GetPixel`, `GetDepth`. Commit `test:`.
- [x] Task 1.2: Implement `Framebuffer` struct and basic accessors. Commit `feat:`.
- [x] Task 1.3: Write and pass tests for `SetPixel` — depth test, out-of-bounds guard, colour write. Commit `test:` then `feat:`.
- [x] Task 1.4: Write and pass test for `SavePNG` — creates file, verify pixel correctness with known colour values. Commit `test:` then `feat:`.
- [x] Task 1.5: Run full test suite (21 tests, 94.6% coverage), open PR, merge to `main`.

### Verification

- [x] 21 framebuffer tests pass; coverage ≥ 90%.
- [x] PR merged to `main`, CI green.

---

_Archived. All tasks complete._
