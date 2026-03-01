# Tech Stack — go_toy_renderer

## Language

| Language | Version | Notes |
|----------|---------|-------|
| Go | 1.24 (toolchain 1.24.7) | Primary language; no CGo for MVP phases |

## Dependencies

### Current (Phases 0–8 — MVP)

| Package | Source | Purpose |
|---------|--------|---------|
| `image/png` | Go stdlib | PNG export from framebuffer |
| `math` | Go stdlib | Float ops (sin, cos, sqrt, etc.) |
| `image` | Go stdlib | RGBA image construction |
| `os` | Go stdlib | File I/O |
| `testing` | Go stdlib | Unit and integration tests |

**No external dependencies.** The MVP is stdlib-only.

### Planned (Phase 9+)

| Package | Phase | Purpose |
|---------|-------|---------|
| `github.com/go-gl/glfw` | 9 | Window creation, input, event loop |
| `go-webgpu/webgpu` v0.4.0 (Zero-CGo FFI) | 11 | WebGPU backend (D3D12/Metal/Vulkan); set `WGPU_NATIVE_PATH` at runtime |
| `naga-cli` (external tool) | 13 | HLSL → WGSL shader compilation |
| DXC or naga | 13 | HLSL cross-compilation toolchain |

## Output / Display

| Phase | Output |
|-------|--------|
| 0–8 (MVP) | `output.png` (640×480 PNG file, CPU rendered) |
| 9+ | Real-time window via GLFW (CPU blit initially, then GPU) |

## Database

None. Assets will eventually be loaded from OBJ files (Phase 16) — no database.

## Package Dependency Graph

Dependencies flow in **one direction only** — lower packages never import higher ones.

```
cmd/renderer       → pkg/render, pkg/framebuffer, pkg/camera, pkg/geometry, pkg/shader
cmd/renderer-rt    → pkg/renderer (Phase 9+), pkg/window (Phase 9+)
pkg/render         → pkg/rasterize, pkg/shader, pkg/framebuffer, pkg/camera, pkg/geometry, pkg/math
pkg/renderer       → pkg/render, pkg/gpu (Phase 10+)
pkg/gpu            → pkg/geometry, pkg/math (Phase 11+)
pkg/rasterize      → pkg/shader, pkg/framebuffer, pkg/math
pkg/shader         → pkg/math
pkg/framebuffer    → pkg/math
pkg/camera         → pkg/math
pkg/geometry       → pkg/math
pkg/math           → (stdlib only)
```

## Infrastructure & CI

### CI Pipeline — GitHub Actions (`.github/workflows/ci.yml`)

**6 jobs, run in parallel after format-validate:**

| Job | Duration | What it does |
|-----|----------|-------------|
| `format-validate` | ~30 s | `gofmt`, `go vet`, `go mod tidy` |
| `lint` | ~1–2 min | golangci-lint v2 (30+ linters) |
| `build` | ~2–3 min | 3 OS × 2 Go versions = 6 matrix combinations |
| `test` | ~3–5 min | Race detector + coverage enforcement |
| `security` | ~1–2 min | govulncheck vulnerability scan |
| `ci-success` | ~5 s | Aggregate pass/fail for branch protection |

**Build matrix:** Linux × macOS × Windows, Go 1.24 × Go 1.25
**Windows `renderer-rt` build:** `GOARCH=amd64` pinned explicitly; PowerShell PE-header check (machine type `0x8664`) verifies the artifact is 64-bit. `macos-latest` excluded from the pin (ARM64 runner).
**Coverage enforcement:** Overall ≥70%; math package ≥90%
**Total runtime:** ~5–8 minutes

### Security Pipeline (`.github/workflows/security.yml`)

Triggers: weekly (Mon 9 AM UTC) + manual dispatch + push to `main` with go.mod changes.

| Scan | Tool |
|------|------|
| Known CVEs | govulncheck |
| Dependency review | GitHub dependency-review-action (PRs only) |
| Static security | gosec (SARIF → GitHub Security tab) |
| License compliance | license scan |

### Tooling

| Tool | Purpose |
|------|---------|
| GitHub Actions | CI/CD pipeline |
| golangci-lint v2 | 30+ linters (`.golangci.yml`) |
| gofmt / gofumpt / goimports | Code formatting (enforced in pre-commit + CI) |
| govulncheck | Security vulnerability scanning |
| `.githooks/pre-commit` | Local hook: `go fmt` → `go vet` → `golangci-lint` |

## Development Environment

- **Platform:** Cross-platform (Windows primary dev machine; CI covers Linux + macOS)
- **Shell:** Bash (Unix syntax even on Windows via Git Bash / WSL)
- **Editor:** Any; project ships GoLand run configs
