# Implementation Plan: Phase 11 — WebGPU Integration via go-webgpu

**Track ID:** phase-11-webgpu_20260228
**Spec:** [spec.md](./spec.md)
**Created:** 2026-02-28
**Status:** [~] In Progress

## Overview

Add `go-webgpu/webgpu` v0.4.0 (Zero-CGo FFI) as a dependency, create
`pkg/gpu` with the wgpu init chain and Hello Triangle render path, and wire
it into `GPUBackend` in `pkg/renderer`. Since go-webgpu is Zero-CGo, `pkg/gpu`
needs no build tags for compilation — GPU integration tests are gated by
`GPU_TESTS=1` instead. Follow strict TDD throughout.

## Phase 1: Dependency & Package Scaffolding

Add `go-webgpu/webgpu` to the module and create the `pkg/gpu` package skeleton.

### Tasks

- [x] Task 1.1: `go get github.com/go-webgpu/webgpu@v0.4.0` — add to
      go.mod / go.sum. Note: wgpu-native shared libraries are already committed
      to the repo under `assets/{platform}/lib/` (v27.0.4.0) — no separate
      download required for local development
- [x] Task 1.2: Create `pkg/gpu/gpu.go` — `Device` struct holding
      `*wgpu.Instance`, `*wgpu.Adapter`, `*wgpu.Device`, `*wgpu.Queue`
      fields; `New()` constructor; stub `Init` and `RenderFrame` signatures
- [x] Task 1.3: Write failing tests in `pkg/gpu/gpu_test.go`:
      `TestDevice_New_ReturnsNonNil`, `TestDevice_Init_SkipsWithoutGPUTests`
      (skip-guard pattern using `GPU_TESTS` env var)

### Verification

- [x] `go build ./pkg/gpu/...` compiles on all platforms (no CGo)
- [x] `go test ./pkg/gpu/...` passes (GPU tests skip without `GPU_TESTS=1`)

## Phase 2: wgpu Initialization Chain (RED → GREEN)

Instance → Adapter → Device → Queue → Surface → SwapChain.

### Tasks

- [x] Task 2.1: Write failing GPU integration test for `Device.Init` (gated
      by `GPU_TESTS=1`; asserts non-nil fields after call)
- [x] Task 2.2: Implement `wgpu.Init()` call with graceful error return if
      library not found (descriptive message referencing `WGPU_NATIVE_PATH`).
      Local usage: point `WGPU_NATIVE_PATH` to committed asset for the current
      platform (e.g. `assets/windows-x86_64-gnu/lib/wgpu_native.dll`)
- [x] Task 2.3: Implement `wgpu.CreateInstance(nil)` → store instance
- [x] Task 2.4: Implement `instance.RequestAdapter(nil)` with surface
      compatibility hint
- [x] Task 2.5: Implement `adapter.RequestDevice(nil)` → device; call
      `device.GetQueue()` → queue
- [~] Task 2.6: Create wgpu Surface from GLFW window native handle
      (Win32/Cocoa/X11 platform branch)
- [x] Task 2.7: Configure Surface — `BGRA8Unorm` format, present mode `Fifo`

### Verification

- [ ] `GPU_TESTS=1 WGPU_NATIVE_PATH=<path> go test ./pkg/gpu/...` — init
      test passes on GPU machine
- [ ] `--backend gpu` opens GLFW window without panicking
- [ ] `go vet ./pkg/gpu/...` and `golangci-lint run ./pkg/gpu/...` clean

## Phase 3: Hello Triangle (RED → GREEN)

Hardcoded WGSL shader and full GPU render pass.

### Tasks

- [x] Task 3.1: Write failing GPU integration test for `Device.RenderFrame`
      (gated by `GPU_TESTS=1`)
- [x] Task 3.2: Define inline `const helloTriangleWGSL string` — vertex
      shader: NDC triangle from `vertex_index`; fragment: solid orange
      `vec4(1.0, 0.5, 0.2, 1.0)`
- [x] Task 3.3: Create `wgpu.ShaderModule` from inline WGSL
- [x] Task 3.4: Create `wgpu.RenderPipeline` (no vertex buffers; surface
      texture format as color target)
- [x] Task 3.5: Implement `RenderFrame`: acquire surface texture → begin
      render pass → set pipeline → `draw(3, 1, 0, 0)` → end pass → queue
      submit → present
- [~] Task 3.6: Wire `GPUBackend.RenderFrame` in `pkg/renderer/gpu.go` to
      delegate to `(*pkg/gpu.Device).RenderFrame()`

### Verification

- [ ] Hello Triangle visible in GLFW window with `--backend gpu` (manual test)
- [ ] No wgpu validation errors in debug output
- [ ] `go test -cover ./pkg/gpu/...` — coverage >80%

## Phase 4: Integration & CI Hardening

End-to-end wiring, graceful fallback, CI compliance.

### Tasks

- [ ] Task 4.1: Update `pkg/renderer/factory.go` — `--backend auto` tries
      `GPUBackend.Init()` first; falls back to `CPUBackend` on any error
- [ ] Task 4.2: Confirm `go test -tags=headless ./...` still passes (`pkg/gpu`
      needs no headless tag; CPU path unchanged)
- [ ] Task 4.3: Audit `.github/workflows/ci.yml` — verify `pkg/gpu` is covered
      by existing `-tags=headless` build and test steps; confirm no step
      requires `WGPU_NATIVE_PATH` to compile; document `assets/{platform}/lib/`
      paths in CI comments for future GPU-runner integration; update ci.yml if
      any step must be amended
- [ ] Task 4.4: Full quality gates: `go fmt ./...`, `go vet ./...`,
      `golangci-lint run --build-tags headless`, `go test -tags=headless
      -cover ./...`
- [ ] Task 4.5: Update `conductor/architecture.md` — add `pkg/gpu` API
      summary table, coverage, and Zero-CGo note
- [ ] Task 4.6: Update `conductor/tracks.md` status and `CLAUDE.md` Recent
      Updates

## Final Verification

- [ ] All 5 acceptance criteria met
- [ ] `go test -race -tags=headless ./...` passes
- [ ] `golangci-lint run --build-tags headless` passes
- [ ] CI green on all matrix entries
- [ ] Docs updated: architecture.md, tracks.md, CLAUDE.md
- [ ] `.github/workflows/ci.yml` reviewed and updated if needed

---
_Generated by Conductor. Tasks will be marked [~] in progress and [x] complete._
