# Go Style Guide — go_toy_renderer

## Toolchain

| Tool | Version | Config |
|------|---------|--------|
| Go | 1.24 | `go.mod` |
| golangci-lint | v2.x | `.golangci.yml` |
| gofmt | stdlib | enforced in pre-commit |
| gofumpt | latest | enforced in pre-commit |
| goimports | latest | enforced in pre-commit |

Run before every commit (also enforced by pre-commit hook):
```bash
go fmt ./...
go vet ./...
golangci-lint run
```

## Formatting

- **gofmt + gofumpt + goimports** are all active. Never submit unformatted code.
- Imports are grouped automatically by `goimports`: stdlib → external → internal.
- No line-length limit (`lll` is disabled), but keep lines readable.

## Naming

| Construct | Convention | Example |
|-----------|-----------|---------|
| Packages | short, lowercase, no underscores | `rasterize`, `shader`, `framebuffer` |
| Exported types | PascalCase | `Vec3`, `Framebuffer`, `ScreenVertex` |
| Exported functions | PascalCase | `NewScene`, `Triangle`, `SavePNG` |
| Unexported | camelCase | `edgeFunction`, `transformVertex` |
| Test functions | `TestComponent_Behaviour_Expected` | `TestVec3_Add_ReturnsSumOfVectors` |
| Benchmarks | `BenchmarkComponent_Operation` | `BenchmarkTriangle_Rasterize` |

**revive rule:** Avoid stutter — don't repeat the package name in the identifier.
- ✅ `shader.Func` (not `shader.ShaderFunc`)
- ✅ `shader.VertexColor` (not `shader.VertexColorShader`)

## Active Linters (key rules)

| Linter | What it enforces |
|--------|-----------------|
| `revive` | Exported doc comments, no stutter, confidence ≥ 0.8 |
| `staticcheck` | All checks except ST1000, ST1003 |
| `gosec` | Security issues (severity ≥ medium, excludes G104) |
| `errorlint` | Proper `errors.Is` / `errors.As` usage |
| `wrapcheck` | Errors from external packages must be wrapped |
| `godot` | Comments must end with a period |
| `misspell` | US English spelling |
| `nestif` | Max nesting complexity 5 |
| `gocognit` | Max cognitive complexity 20 |
| `gocyclo` | Max cyclomatic complexity 15 |
| `dupl` | No duplicated code blocks |
| `goconst` | Extract repeated string literals |
| `prealloc` | Pre-allocate slices where length is known |

**Test files** are exempt from: `dupl`, `errcheck`, `funlen`, `goconst`, `gocritic`, `gocyclo`, `gosec`, `staticcheck` (dot imports).

**`pkg/math/`** is exempt from `gocognit` "cognitive complexity" warnings (math operations are inherently complex).

## Error Handling

- Always check errors — `errcheck` is active.
- Wrap errors from other packages: `fmt.Errorf("rasterize: %w", err)`.
- Use `errors.Is` / `errors.As` for error inspection — not string comparison.
- Avoid naked returns; `nakedret` is active.
- Avoid `nilnil` — don't return `(nil, nil)` from functions that return `(value, error)`.

## Nolint Directives

If a `//nolint` is truly necessary:
- Always include a specific linter name: `//nolint:gosec` (not `//nolint`).
- Always include an explanation: `//nolint:gosec // G304: path is validated before use`.
- `nolintlint` is active and will reject unexplained or unspecific directives.

## Comments

- All exported types, functions, and methods must have doc comments (`revive:exported`).
- Comments must end with a period (`godot`).
- Use `// TODO(phase-N):` for known future work.
- Avoid `/* block comments */` in Go code; use `//` line comments.

## Testing

- Test file name: `<source>_test.go` alongside the source file.
- Use table-driven tests for multiple cases.
- Float comparisons: always use an epsilon (`±0.0001`).
- Avoid `t.Parallel()` unless the test is genuinely safe to parallelise (`tparallel` is active).
- Helper functions must call `t.Helper()` (`thelper` is active).
- Use `testify` only if stdlib `testing` is insufficient — prefer stdlib for this project.

## Package Structure

Each package owns exactly one pipeline stage. Dependencies flow in one direction:

```
cmd/renderer → pkg/render → pkg/rasterize, pkg/shader, pkg/framebuffer, pkg/camera, pkg/geometry, pkg/math
```

Never import a higher-level package from a lower-level one.

## Performance Notes

- Use `make([]T, 0, n)` when length is known — `prealloc` will catch missed opportunities.
- Avoid allocation in hot paths (rasterizer inner loop, vertex transform).
- Benchmark with `go test -bench=. -benchmem` after functionality is proven.
