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
   - Relevant `.claude/docs/` component doc updated (status, API examples, coverage).
   - `.claude/docs/12-development-roadmap.md` — phase marked ✅ COMPLETED.
   - `.claude/docs/10-mvp-features.md` — feature statuses updated.
   - `.claude/docs/README.md` — Last Updated + Next Steps refreshed.
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
