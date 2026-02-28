# CLAUDE.md

This file provides concise guidance to Claude Code. For full project context, see the Conductor artifacts below.

## Context Sources

| Artifact | Contents |
|----------|----------|
| [conductor/index.md](./conductor/index.md) | Navigation hub — start here |
| [conductor/product.md](./conductor/product.md) | Project vision, pipeline overview, conventions |
| [conductor/architecture.md](./conductor/architecture.md) | Package APIs, pipeline diagram, GPU architecture |
| [conductor/workflow.md](./conductor/workflow.md) | TDD policy, commit strategy, branching, coverage thresholds |
| [conductor/tracks.md](./conductor/tracks.md) | All 18 tracks (9 archived MVP + 8 GPU pending + 1 active) |
| [conductor/product-guidelines.md](./conductor/product-guidelines.md) | Design principles, common gotchas |

## Project Overview

A toy 3D renderer in Go — CPU software renderer MVP complete (Phases 0–8), producing a 640×480 PNG of a coloured cube via barycentric rasterization. Now advancing toward cross-platform GPU acceleration (Phases 9–16) using wgpu-native (WebGPU), GLFW, and HLSL shaders compiled to WGSL via naga-cli. See [conductor/product.md](./conductor/product.md) for full details.

## Development Commands

```bash
# One-time: activate pre-commit hook (fmt → vet → lint)
git config core.hooksPath .githooks

# Build CPU renderer
go build -o renderer ./cmd/renderer

# Build real-time renderer (Phase 9+)
go build -o renderer-rt ./cmd/renderer-rt

# Run CPU renderer (produces output.png)
go run ./cmd/renderer

# Run real-time renderer
go run ./cmd/renderer-rt --backend cpu   # CPU blit
go run ./cmd/renderer-rt --backend gpu   # wgpu-native
go run ./cmd/renderer-rt --backend auto  # GPU with CPU fallback

# Test
go test ./...
go test -cover ./...
go test -v -run TestName ./path/to/pkg

# Quality gates (must pass before PR)
go fmt ./...
go vet ./...
golangci-lint run           # requires golangci-lint v2
go test -race -coverprofile=coverage.out ./...
govulncheck ./...

# GPU shaders (Phase 13+) — requires: cargo install naga-cli
go generate ./assets/shaders/
```

## Windows Build Prerequisites (renderer-rt)

`cmd/renderer-rt` uses CGO (GLFW + OpenGL) for the windowed mode. On Windows,
the C toolchain determines the output binary's architecture — a 32-bit compiler
produces a 32-bit PE that fails with **error 193** on x64 systems.

**Required:** 64-bit MinGW-w64 GCC. Install via MSYS2:
```bash
pacman -S mingw-w64-x86_64-gcc
# Add to PATH: C:\msys64\mingw64\bin
```

Always set `GOARCH=amd64` explicitly when building for Windows x64:
```bash
GOARCH=amd64 go build -o renderer-rt.exe ./cmd/renderer-rt
```

CI enforces this automatically for `windows-latest` builds.

## Common Gotchas

**Math:**
- Matrix multiplication order (easy to reverse — result = Projection × View × Model)
- Left vs right-handed coordinates — project uses right-handed (+X right, +Y up, +Z out)
- Float epsilon comparisons: use `±0.0001`

**Transforms:**
- Perspective divide by zero: reject triangles with W ≤ 0 before NDC conversion
- Depth convention: 0 = near, 1 = far; depth buffer initialised to 1.0
- Screen Y is flipped: `sY = (1 − ndcY) / 2 × height`

**Rasterization:**
- Use float coordinates for interpolation, not integers
- Depth test direction: `incoming < stored` means closer wins
- Colour clamping: barycentric results can exceed [0, 1]

**Pipeline:**
- Triangles behind camera (negative W) must be rejected before rasterization
- Backface winding: project uses CCW for front faces
- MVP order: `transformVertex` applies `ViewProjectionMatrix()` only (no separate model matrix until Phase 14)

## External Resources

- [Scratchapixel](https://www.scratchapixel.com/) — 3D graphics fundamentals
- [Learn OpenGL](https://learnopengl.com/) — modern GPU pipeline concepts
- [WebGPU Spec](https://www.w3.org/TR/webgpu/) — wgpu-native reference
- [naga-cli docs](https://github.com/gfx-rs/wgpu/tree/trunk/naga) — HLSL→WGSL compiler

## Recent Updates

**2026-02-28 (Latest):** Phase 10 complete — `pkg/renderer` package: `Renderer` interface, `CPUBackend` (GLFW blit extracted from Phase 9), `GPUBackend` stub, `New()` factory. `cmd/renderer-rt` refactored to use `renderer.New()`. `--backend auto` falls back to CPU (GPU-first in Phase 11). `pkg/renderer` at 100% test coverage.

**2026-02-28:** CI hardening complete — pinned Go to 1.24, golangci-lint to v2.1.6, govulncheck to v1.1.3, gosec to v2.22.0; split Windows race-detector step; removed `|| true` from license enforcement. See `conductor/ci-fragility.md`.

**2026-02-28:** Bug fix — Windows x64 `renderer-rt` build pinned to `GOARCH=amd64` in CI; PowerShell PE-header check added. Local build requires 64-bit MinGW-w64 (see Windows Build Prerequisites above).
