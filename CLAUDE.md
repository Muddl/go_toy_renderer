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

**2026-02-27 (Latest):** Conductor migration complete — all `.claude/docs/` and `.claude/tasks/` content migrated into `conductor/` artifacts. 9 archived tracks (Phases 0–8) + 8 pending GPU tracks (Phases 9–16) created. CLAUDE.md trimmed to pointer document.

**2026-02-27:** GPU Roadmap set — WebGPU (wgpu-native) + HLSL→WGSL via naga-cli across Phases 9–16. New docs: `14-gpu-backend.md`, `15-hlsl-shader-pipeline.md`, `16-realtime-display.md`. Architecture and CLAUDE.md updated.

**2026-02-27:** Phase 8 Complete — MVP done. Golden image test `TestRender_GoldenImage_Triangle`, 7 integration tests, benchmarks (~1.2 ms/frame, ~58 ns/vertex), full README written.
