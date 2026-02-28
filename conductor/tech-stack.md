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
| `wgpu-native` (CGo binding) | 11 | WebGPU backend (D3D12/Metal/Vulkan) |
| `naga-cli` (external tool) | 13 | HLSL → WGSL shader compilation |
| DXC or naga | 13 | HLSL cross-compilation toolchain |

## Output / Display

| Phase | Output |
|-------|--------|
| 0–8 (MVP) | `output.png` (640×480 PNG file, CPU rendered) |
| 9+ | Real-time window via GLFW (CPU blit initially, then GPU) |

## Database

None. Assets will eventually be loaded from OBJ files (Phase 16) — no database.

## Infrastructure & CI

| Tool | Purpose |
|------|---------|
| GitHub Actions | CI/CD — build matrix (Linux/macOS/Windows × Go 1.24/1.25) |
| golangci-lint v2 | Static analysis (30+ linters, see `.golangci.yml`) |
| govulncheck | Security vulnerability scanning |
| `.githooks/pre-commit` | Local pre-commit hook (fmt → vet → lint) |

## Development Environment

- **Platform:** Cross-platform (Windows primary dev machine; CI covers Linux + macOS)
- **Shell:** Bash (Unix syntax even on Windows via Git Bash / WSL)
- **Editor:** Any; project ships GoLand run configs
