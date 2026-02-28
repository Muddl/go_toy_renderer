# CI Fragility Reference

This document catalogues fragile patterns found in the CI pipeline during the
`ci-hardening_20260228` audit (2026-02-28). Each entry records the problem,
root cause, risk, and the fix applied. Consult this before modifying either
workflow file.

---

## Patterns Found & Fixed

### P1 — Race detector on `windows-latest`

| Field | Detail |
|-------|--------|
| **File** | `.github/workflows/ci.yml` |
| **Job** | `test` |
| **Line** | `go test -tags=headless -race -v ./...` — no OS guard |
| **Root cause** | The Go race detector requires CGo to build its runtime DLL. On Windows runners (MSVC/MinGW), the CGo toolchain is not reliably present in the headless path, causing the race-detector DLL to fail to load. |
| **Risk** | Intermittent `fatal error: race` or linker failures on every Windows test run. |
| **Fix applied** | Added `if: matrix.os != 'windows-latest'` to the race step; added a separate Windows-only step running `go test -tags=headless -v ./...` without `-race`. Race coverage is preserved on Linux and macOS. |

---

### P2 — Non-existent Go version `1.25`

| Field | Detail |
|-------|--------|
| **File** | `.github/workflows/ci.yml` |
| **Jobs** | `format-validate`, `lint`, `test` (all matrix rows), `security` |
| **Root cause** | Go 1.25 did not exist at the time of authoring. `actions/setup-go` may silently resolve it to a different version or fail entirely depending on runner availability. |
| **Risk** | Silent version substitution (running 1.24 while thinking 1.25) or outright job failure. The build matrix then tests 1.24 × 1.25 but effectively tests 1.24 × 1.24. |
| **Fix applied** | All `go-version: '1.25'` replaced with `go-version: '1.24'`. Update to a real future version when it is released and tested. |

---

### P3 — `golangci-lint` binary version floats (`version: latest`)

| Field | Detail |
|-------|--------|
| **File** | `.github/workflows/ci.yml` |
| **Job** | `lint` |
| **Step** | `golangci/golangci-lint-action@v8`, `with: version: latest` |
| **Root cause** | `version: latest` in the `with:` block controls the linter binary version, not the action version. Any new golangci-lint release can introduce stricter defaults, new enabled linters, or changed rule behaviour, silently breaking CI on a commit that touched no Go code. |
| **Risk** | Unexpected lint failures on days when golangci-lint releases, unrelated to code changes. |
| **Fix applied** | `version: latest` → `version: v2.1.6` (latest stable at time of hardening). Pin to a new version intentionally when upgrading. |

---

### P4 — `govulncheck@latest` in `ci.yml`

| Field | Detail |
|-------|--------|
| **File** | `.github/workflows/ci.yml` |
| **Job** | `security` |
| **Step** | `go install golang.org/x/vuln/cmd/govulncheck@latest` |
| **Root cause** | `@latest` resolves the latest published module version at the time the job runs. Behaviour can change between runs without any repo change. |
| **Risk** | New govulncheck version may report new findings or change exit codes, breaking CI unpredictably. |
| **Fix applied** | `@latest` → `@v1.1.3`. Update intentionally when a new stable version is verified. |

---

### P5 — `securego/gosec@master` in `security.yml`

| Field | Detail |
|-------|--------|
| **File** | `.github/workflows/security.yml` |
| **Job** | `gosec-scan` |
| **Step** | `uses: securego/gosec@master` |
| **Root cause** | Referencing `@master` on a third-party action means any commit to the gosec repo's master branch — including breaking changes, removed flags, or changed SARIF schema — immediately affects CI. This is the most dangerous pattern: it gives external maintainers implicit write access to your pipeline behaviour. |
| **Risk** | Scan step breaks on any gosec release; or, worse, silently changes scan behaviour. Also a supply-chain risk. |
| **Fix applied** | Pinned to `securego/gosec@v2.22.0` (latest tagged release at time of hardening). |

---

### P6 — `go-licenses` steps silenced with `|| true`

| Field | Detail |
|-------|--------|
| **File** | `.github/workflows/security.yml` |
| **Job** | `license-check` |
| **Steps** | `go-licenses report ./... \|\| true` and `go-licenses save ./... --save_path=licenses \|\| true` |
| **Root cause** | Both license check commands were suffixed with `|| true`, meaning any failure — including finding a license violation — was swallowed. The job always passed, providing no actual enforcement. |
| **Risk** | A dependency with a GPL or other incompatible license could be introduced undetected. The CI badge would remain green regardless. |
| **Fix applied** | `|| true` removed from both lines. Added `continue-on-error: false` (default) so violations now fail the job. Investigated current dependencies to confirm they pass before removing the suppression. |

---

### P7 — `govulncheck@latest` in `security.yml`

| Field | Detail |
|-------|--------|
| **File** | `.github/workflows/security.yml` |
| **Job** | `vulnerability-scan` |
| **Step** | `go install golang.org/x/vuln/cmd/govulncheck@latest` |
| **Root cause** | Same as P4 — unpinned tool install. |
| **Risk** | Identical to P4: version drift causing silent behaviour change or unexpected failures. |
| **Fix applied** | `@latest` → `@v1.1.3` (same version as ci.yml for consistency). |

---

### P8 — Go version inconsistency between `ci.yml` and `security.yml`

| Field | Detail |
|-------|--------|
| **File** | `.github/workflows/security.yml` |
| **Jobs** | All jobs used `go-version: '1.24'` while `ci.yml` used `'1.25'` |
| **Root cause** | The two files were authored/updated at different times without cross-referencing each other. `ci.yml` had been updated to 1.25 (which doesn't exist — see P2); `security.yml` was left on 1.24. |
| **Risk** | Security scans run against a different Go toolchain than the main CI, meaning govulncheck or gosec could use a different module graph. |
| **Fix applied** | Both files aligned to `go-version: '1.24'`. Update in sync going forward. |

---

## Future Phase Guidance

### Version Pinning Policy

**Rule:** Never use `@latest`, `@master`, or `version: latest` for any tool installed or invoked in CI.

| Context | Correct pattern | Example |
|---------|----------------|---------|
| `go install` tool | Pin to exact module version | `go install tool@v1.2.3` |
| GitHub Action | Pin to major version OR exact SHA | `uses: foo/bar@v4` or `uses: foo/bar@abc1234` |
| golangci-lint binary | Pin to semver in `with:` block | `version: v2.1.6` |
| Go toolchain | Pin to existing released version | `go-version: '1.24'` |

When upgrading a pinned version:
1. Check the release notes for breaking changes.
2. Update the pin in all workflow files at once.
3. Note the version bump in the commit message.

---

### Race Detector Rules

- **Linux / macOS:** Always run `go test -race`. The race detector is fully supported.
- **Windows:** Never run `go test -race` in CI when CGo is involved (even with `--tags=headless`, the race runtime DLL link can be fragile depending on the runner's C toolchain). Run tests without `-race` on Windows.
- **Pattern to use:**

```yaml
- name: Run tests (with race detector)
  if: matrix.os != 'windows-latest'
  run: go test -tags=headless -race -v ./...

- name: Run tests (Windows, no race detector)
  if: matrix.os == 'windows-latest'
  run: go test -tags=headless -v ./...
```

---

### `|| true` Policy

**Rule:** Never append `|| true` to a CI step that is intended to enforce a quality gate.

`|| true` means "this step can never fail CI". It is only appropriate for:
- Generating optional reports that are uploaded as artifacts (i.e., the CI doesn't care if the report tool errors, only if the artifact was produced).
- Cleanup steps where failure is acceptable.

If a quality gate is genuinely informational (e.g., license scanning is advisory), use `continue-on-error: true` at the step level with an explicit comment explaining why it is non-blocking. This is clearer than `|| true` and is visible in the GitHub Actions UI.

```yaml
# Correct: informational step — non-blocking, intent is explicit
- name: Generate license report (informational)
  continue-on-error: true  # advisory only — does not block merge
  run: go-licenses report ./...

# Wrong: silently swallows failures
- name: Check licenses
  run: go-licenses report ./... || true
```

---

### Cross-File Consistency

Both `ci.yml` and `security.yml` must use the same Go version. When updating one, update both. Add a grep check before opening a PR:

```bash
grep -r "go-version:" .github/workflows/
```

All versions should be identical across all workflow files.

---

_Last updated: 2026-02-28 (ci-hardening_20260228)_
