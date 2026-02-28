# Workflow — go_toy_renderer

## TDD Policy — Strict

**Every feature and bug fix must follow Red-Green-Refactor.**

1. **RED** — Write a failing test that describes the desired behaviour. Run it to confirm it fails. Commit: `test: add failing test for <feature>`.
2. **GREEN** — Write the minimal implementation to make the test pass. Commit: `feat: implement <feature>`.
3. **REFACTOR** — Clean up without breaking tests. Commit: `refactor: improve <component>` (only if needed).

**Rules:**
- Never write production code without a prior failing test.
- One test at a time; make it pass; move to the next.
- Test names: `TestComponent_Behaviour_ExpectedOutcome` (e.g. `TestVec3_Add_ReturnsSumOfVectors`).
- Use epsilon comparison for floats: `±0.0001`.
- Benchmark performance-critical paths only after functionality is proven correct.

### Test Naming Convention

```go
// Unit tests
func TestComponent_Behaviour_ExpectedOutcome(t *testing.T)
// Examples:
//   TestVec3_Add_ReturnsSumOfVectors
//   TestMatrix_Multiply_HandlesIdentityMatrix
//   TestRasterizer_Triangle_SkipsDegenerateTriangle

// Table-driven tests
func TestFunctionName_Cases(t *testing.T)

// Integration tests
func TestRender_Integration_ScenarioName(t *testing.T)

// Golden image tests
func TestRender_GoldenImage_SceneName(t *testing.T)

// Benchmarks
func BenchmarkComponent_Operation(b *testing.B)
// Examples:
//   BenchmarkRender_Cube_640x480
//   BenchmarkTriangle_Rasterize_Large
```

### Coverage Thresholds

| Package | Minimum | Current |
|---------|---------|---------|
| `pkg/math` | **>90%** | 96.1% ✅ |
| `pkg/camera` | >80% | 100% ✅ |
| `pkg/geometry` | >80% | ~100% ✅ |
| `pkg/rasterize` | >80% | 100% ✅ |
| `pkg/shader` | >80% | 100% ✅ |
| `pkg/render` | >80% | 100% ✅ |
| `pkg/framebuffer` | >80% | 90.2% ✅ |
| Overall | **>70%** | Well above ✅ |

Check coverage: `go test -cover ./...`
Detailed report: `go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out`

### Golden Image Tests

After the rendering pipeline is complete or when visual output changes, golden image tests verify pixel-exact correctness:

1. Render a deterministic scene (fixed geometry, fixed camera, fixed seed).
2. Save as PNG.
3. Byte-compare against reference in `testdata/`.
4. Use `-update` flag to regenerate reference when rendering intentionally changes.

```go
// Pattern (see pkg/render/integration_test.go)
func TestRender_GoldenImage_Triangle(t *testing.T) {
    // render to 100×100 fb
    // if *update flag: overwrite testdata/golden_triangle.png
    // else: byte-compare against testdata/golden_triangle.png
}
```

**Reference files** are committed to the repo in `testdata/` alongside the test.

### Test Pyramid

```
        ┌─────────────┐
        │  Golden     │  ← Few; visual regression guard
        │  Image      │
        ├─────────────┤
        │ Integration │  ← Some; key pipeline workflows
        │   Tests     │
        ├─────────────┤
        │    Unit     │  ← Many; fast, one behaviour each
        │   Tests     │
        └─────────────┘
```

## Commit Strategy — Conventional Commits

```
<type>: <brief description>

<optional body>

<optional footer>
```

**Types:**
| Type | Use for |
|------|---------|
| `feat:` | New feature or behaviour |
| `fix:` | Bug fix |
| `test:` | Add or update tests |
| `refactor:` | Restructuring without behaviour change |
| `docs:` | Documentation only |
| `style:` | Formatting, no logic change |
| `perf:` | Performance improvement |
| `chore:` | Build / tooling / CI changes |

**Examples:**
```
test: add failing test for perspective divide
feat: implement perspective divide in transformVertex
refactor: extract NDC conversion to helper
docs: update Phase 9 roadmap with GLFW dependency
```

## Branching Strategy — Trunk-Based Development

- **Never commit directly to `main`.**
- All changes go through feature/bugfix/release branches and pull requests.

**Branch naming:**
| Prefix | Use for |
|--------|---------|
| `feature/` | New functionality |
| `bugfix/` | Bug fixes |
| `release/` | Release preparation |
| `hotfix/` | Critical production fixes |

**Daily workflow:**
1. `git checkout main && git pull origin main`
2. `git checkout -b feature/<name>`
3. Write failing test → commit (`test:`)
4. Implement → commit (`feat:`)
5. Refactor if needed → commit (`refactor:`)
6. `go test ./... && go fmt ./... && golangci-lint run`
7. `git push -u origin feature/<name>`
8. Open PR; wait for CI to pass; merge.

## Code Review Policy — Required for All Changes

- All changes require a PR — even single-author work.
- Self-review checklist before opening PR:
  - [ ] All tests pass: `go test ./...`
  - [ ] Code formatted: `go fmt ./...`
  - [ ] No vet warnings: `go vet ./...`
  - [ ] No lint warnings: `golangci-lint run`
  - [ ] Coverage meets thresholds (>90% math, >80% core, >70% overall)
  - [ ] Commit messages follow Conventional Commits format
  - [ ] Changes match branch purpose

## Verification Checkpoints — After Each Phase

A phase is only **complete** when ALL of the following are done:

1. **Implementation** — All code written, all tests passing.
2. **PR merged** — Feature branch merged to `main` and deleted.
3. **Task summary** — `.claude/tasks/<date>_<phase>.md` created.
4. **Documentation updated:**
   - `conductor/architecture.md` — update API summary table and coverage for affected packages.
   - Relevant archived track in `conductor/tracks/` marked complete in `plan.md` and `metadata.json`.
   - `conductor/tracks.md` — track status updated to `[x]`.
   - `conductor/index.md` — Active Tracks section updated.
   - `CLAUDE.md` — Recent Updates section reflects phase completion (last 3 entries only).
5. **CLAUDE.md** — Recent Updates section reflects phase completion.

## Task Lifecycle

```
pending → in_progress → completed
```

- Mark a task `in_progress` before starting work.
- Mark `completed` only when fully done (tests pass, PR merged, docs updated).
- If blocked, create a new task describing the blocker rather than leaving the original stalled.

## Pre-commit Hook

Activate once per clone:
```bash
git config core.hooksPath .githooks
```

Runs automatically on every `git commit`: `go fmt ./...` → `go vet ./...` → `golangci-lint run`.
