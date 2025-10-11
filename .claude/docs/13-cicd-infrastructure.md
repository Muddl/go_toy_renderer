# CI/CD Infrastructure

**Status:** ✅ Implemented (Phase 0)
**Implementation Date:** 2025-10-10

## Overview

Comprehensive GitHub Actions CI/CD pipeline providing automated testing, security scanning, and quality gates for the toy renderer project. The pipeline ensures code quality, cross-platform compatibility, and security before code reaches the main branch.

## Purpose

**Why CI/CD for a Learning Project?**
- **Automated quality checks** - Catch issues before they become problems
- **Professional workflow** - Learn industry-standard practices
- **Fast feedback** - Know within 5-8 minutes if changes break anything
- **Confidence to refactor** - Safety net for learning and experimentation
- **Multi-platform verification** - Ensure code works on Linux, macOS, Windows

## Pipeline Architecture

### Main CI Pipeline (`.github/workflows/ci.yml`)

**6 Specialized Jobs:**

1. **format-validate** (~30 seconds)
   - Code formatting check (`gofmt`)
   - Go vet for suspicious constructs
   - Dependency verification (`go mod tidy`)
   - Fast feedback - runs first

2. **lint** (~1-2 minutes)
   - 48+ linters via golangci-lint
   - Code quality, best practices, security
   - Configuration: `.golangci.yml`

3. **build** (~2-3 minutes)
   - **Matrix:** 3 OS × 2 Go versions = 6 combinations
     - Linux, macOS, Windows
     - Go 1.22 and 1.23
   - Cross-platform compatibility verification
   - Build artifacts published (7-day retention)

4. **test** (~3-5 minutes)
   - Full test suite with race detector
   - Coverage calculation (Linux + Go 1.23 only)
   - **Coverage enforcement:**
     - Overall project: ≥70%
     - Math package: ≥90%
   - Coverage reports uploaded (30-day retention)

5. **security** (~1-2 minutes)
   - govulncheck vulnerability scanning
   - Dependency security analysis
   - JSON reports generated

6. **ci-success** (~5 seconds)
   - Aggregates all job results
   - Single status check for branch protection
   - ✅ Pass = all jobs succeeded

**Total Runtime:** ~5-8 minutes (jobs run in parallel)

### Security Pipeline (`.github/workflows/security.yml`)

**Automated Security Scanning:**

1. **Vulnerability Scan** (govulncheck)
   - Checks for known CVEs in dependencies
   - JSON report generation
   - 90-day artifact retention

2. **Dependency Review** (PRs only)
   - Reviews dependency changes in pull requests
   - Fails on moderate+ severity issues
   - Uses GitHub's dependency graph

3. **Security Code Scan** (gosec)
   - Static security analysis
   - SARIF report upload to GitHub Security tab
   - Checks: SQL injection, hardcoded secrets, etc.

4. **License Compliance**
   - Scans all dependency licenses
   - Generates license report
   - Legal compliance tracking

5. **Security Summary**
   - Aggregates scan results
   - Creates workflow summary
   - Links to Security tab

**Triggers:**
- Weekly schedule (Mondays at 9:00 AM UTC)
- Manual workflow dispatch
- Push to `main` when go.mod/go.sum changes

## Linter Configuration (`.golangci.yml`)

**48+ Enabled Linters:**

**Core Analysis:**
- errcheck, gosimple, govet, staticcheck, unused

**Code Quality:**
- revive, gocritic, gocyclo (≤15), gocognit (≤20)

**Security:**
- gosec, forbidigo, forcetypeassert

**Best Practices:**
- bodyclose, errorlint, nilerr, nilnil

**Style:**
- gofmt, goimports, stylecheck, gofumpt

**Testing:**
- thelper, tparallel

**Project-Appropriate Settings:**
- Cyclomatic complexity: ≤15 (reasonable for algorithms)
- Cognitive complexity: ≤20 (allows learning-focused code)
- Nesting depth: ≤5
- Relaxed rules for test files
- US English spell checking

**Disabled Linters:**
- Too strict for learning: funlen, lll, gochecknoglobals
- Too opinionated: wsl, nlreturn, varnamelen
- Not needed yet: depguard, interfacebloat

## Files Created

```
.github/
├── workflows/
│   ├── ci.yml              # Main CI pipeline (~8 KB)
│   └── security.yml        # Security scanning (~5 KB)

.golangci.yml              # Linter configuration (~7 KB)
```

**Total:** 3 configuration files, ~20 KB

## Workflow Triggers

**Main CI Pipeline:**
- Push to `main` branch
- All pull requests (any branch)
- Manual workflow dispatch

**Security Pipeline:**
- Weekly schedule (Mondays 9 AM UTC)
- Manual workflow dispatch
- Push to `main` with go.mod/go.sum changes

**Concurrency Control:**
```yaml
concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true
```
Prevents wasted resources by canceling outdated workflow runs.

## GitHub Actions Features

**Actions Used:**
- `actions/checkout@v4` - Repository checkout
- `actions/setup-go@v5` - Go setup with caching
- `golangci/golangci-lint-action@v6` - Official linter
- `actions/upload-artifact@v4` - Artifact publishing
- `github/codeql-action/upload-sarif@v3` - Security findings
- `actions/dependency-review-action@v4` - Dependency security

**Performance Optimizations:**
- Go module caching (speeds up dependency downloads)
- Build caching (speeds up compilation)
- Concurrent job execution
- Strategic artifact retention policies

## CI/CD Quick Start

### What Happens on Push/PR

When you push code or create a PR:

1. ✅ **Format & Validate** - Code formatting, go vet, mod tidy
2. ✅ **Lint** - 48+ linters report issues
3. ✅ **Build** - 6-way matrix (3 OS × 2 Go versions)
4. ✅ **Test** - Race detector + coverage enforcement
5. ✅ **Security** - Vulnerability scanning
6. ✅ **Final Status** - Aggregate pass/fail

**You'll see:** Status checks on PR, downloadable artifacts, Security tab findings

### Viewing Results

**GitHub Actions Tab:**
```
https://github.com/muddl/go_toy_renderer/actions
```
- All workflow runs
- Downloadable artifacts (builds, coverage, security reports)
- Detailed job logs

**Security Tab:**
```
https://github.com/muddl/go_toy_renderer/security
```
- Security findings (SARIF reports)
- Dependency alerts
- Vulnerability dashboard

### Local Pre-Commit Checks

Run before pushing to catch issues early:

```bash
# Format code
go fmt ./...

# Run go vet
go vet ./...

# Run linter (if installed)
golangci-lint run

# Run tests with race detector
go test -race ./...

# Check coverage
go test -cover ./...

# Security scan (if installed)
govulncheck ./...
```

### Installing Tools Locally (Optional)

**golangci-lint:**
```bash
# macOS/Linux
brew install golangci-lint

# Windows (Chocolatey)
choco install golangci-lint

# Or with Go
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

**govulncheck:**
```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
```

## Adding CI Badges to README

### Badge Markdown

**CI Pipeline Status:**
```markdown
[![CI Pipeline](https://github.com/muddl/go_toy_renderer/actions/workflows/ci.yml/badge.svg)](https://github.com/muddl/go_toy_renderer/actions/workflows/ci.yml)
```

**Security Scanning:**
```markdown
[![Security](https://github.com/muddl/go_toy_renderer/actions/workflows/security.yml/badge.svg)](https://github.com/muddl/go_toy_renderer/actions/workflows/security.yml)
```

**Go Version:**
```markdown
![Go Version](https://img.shields.io/badge/Go-1.22%20%7C%201.23-00ADD8?logo=go)
```

**License:**
```markdown
![License](https://img.shields.io/github/license/muddl/go_toy_renderer)
```

### Recommended Badge Layout

```markdown
# Go Toy 3D Renderer

[![CI Pipeline](https://github.com/muddl/go_toy_renderer/actions/workflows/ci.yml/badge.svg)](https://github.com/muddl/go_toy_renderer/actions/workflows/ci.yml)
[![Security](https://github.com/muddl/go_toy_renderer/actions/workflows/security.yml/badge.svg)](https://github.com/muddl/go_toy_renderer/actions/workflows/security.yml)
![Go Version](https://img.shields.io/badge/Go-1.22%20%7C%201.23-00ADD8?logo=go)
![License](https://img.shields.io/github/license/muddl/go_toy_renderer)

A learning project: CPU-based 3D software renderer built from scratch in Go.
```

### Optional External Services

**Go Report Card** (free, automatic analysis):
1. Visit https://goreportcard.com/
2. Submit repository URL
3. Add badge:
```markdown
[![Go Report Card](https://goreportcard.com/badge/github.com/muddl/go_toy_renderer)](https://goreportcard.com/report/github.com/muddl/go_toy_renderer)
```

**Codecov** (detailed coverage with PR comments):
1. Sign up at https://codecov.io
2. Add repository
3. Uncomment Codecov step in `ci.yml`
4. Add `CODECOV_TOKEN` secret
5. Add badge:
```markdown
[![codecov](https://codecov.io/gh/muddl/go_toy_renderer/branch/main/graph/badge.svg)](https://codecov.io/gh/muddl/go_toy_renderer)
```

## Post-Merge Setup

### 1. Enable Branch Protection (Recommended)

**Prevent merging broken code:**

```
GitHub → Settings → Branches → Add rule

Branch name pattern: main
☑ Require status checks to pass before merging
☑ Require branches to be up to date before merging

Status checks to require:
  ☑ ci-success
```

This ensures all CI checks must pass before PR merge.

### 2. Monitor Security Tab

- Visit https://github.com/muddl/go_toy_renderer/security
- Review security findings (should be clean)
- Enable Dependabot alerts (recommended)
- Configure security advisories if desired

### 3. Add Badges to README

- Use badge markdown from section above
- Add at top of README for visibility
- Update after first successful workflow run

## Troubleshooting

### CI Failure: Format Issues

**Problem:** Code not formatted correctly

**Fix:**
```bash
go fmt ./...
git add .
git commit -m "fix: apply code formatting"
git push
```

### CI Failure: Lint Issues

**Problem:** Linter reports style/quality issues

**Fix:**
```bash
# See what's wrong locally
golangci-lint run

# Fix issues or adjust .golangci.yml if too strict
```

### CI Failure: Build Issues

**Problem:** Build fails on specific platform

**Fix:**
```bash
# Test build locally
go build -v ./...

# Check for platform-specific code issues
# Review build matrix logs for details
```

### CI Failure: Test Failures

**Problem:** Tests fail or race conditions detected

**Fix:**
```bash
# Run tests with race detector
go test -race -v ./...

# Fix failing tests
# Check for race conditions in parallel code
```

### CI Failure: Coverage Below Threshold

**Problem:** Coverage is below 70% (or 90% for math package)

**Fix:**
```bash
# Check current coverage
go test -cover ./...

# Add tests to increase coverage
# Focus on untested code paths
```

### CI Failure: Security Vulnerabilities

**Problem:** govulncheck finds vulnerabilities

**Fix:**
```bash
# Check vulnerabilities locally
govulncheck ./...

# Update dependencies
go get -u ./...
go mod tidy

# Verify fix
govulncheck ./...
```

## Integration with TDD Workflow

The CI/CD pipeline enhances the Test-Driven Development workflow:

**Local Development (TDD):**
1. RED - Write failing test → commit
2. GREEN - Make test pass → commit
3. REFACTOR - Improve code → commit
4. Push feature branch

**Automated CI (Quality Gates):**
5. CI runs all checks automatically
6. Verify coverage thresholds met
7. Ensure cross-platform compatibility
8. Security scan passes

**Merge to Main:**
9. All CI checks ✅
10. Human code review (optional)
11. Merge PR
12. Clean up feature branch

## Best Practices Applied

**1. Matrix Strategy:**
- Tests on 3 platforms × 2 Go versions = 6 combinations
- Catches platform-specific bugs early
- Ensures broad compatibility

**2. Efficient Caching:**
- Go module cache speeds up dependency downloads
- Build cache reduces compilation time
- ~50% faster CI runs

**3. Smart Job Dependencies:**
- Fast feedback first (format-validate)
- Parallel execution (lint, build, test, security)
- Final aggregation (ci-success)

**4. Coverage Strategy:**
- Calculate coverage only on canonical platform (Linux + Go 1.23)
- Reduces CI time while maintaining quality
- Artifacts available for detailed analysis

**5. Security Integration:**
- SARIF upload enables GitHub Security tab
- Dependency review prevents vulnerable deps
- Weekly scans catch newly discovered CVEs

**6. Artifact Management:**
- Build artifacts: 7-day retention (debugging)
- Coverage reports: 30-day retention (analysis)
- Security reports: 90-day retention (compliance)

## Impact on Development

### Before CI/CD
- ❌ Manual testing required
- ❌ No automated quality checks
- ❌ Potential for unnoticed bugs
- ❌ No security scanning
- ❌ Coverage unknown

### After CI/CD
- ✅ Automated testing on every push/PR
- ✅ Multi-platform compatibility verified
- ✅ Code quality enforced (48+ linters)
- ✅ Coverage thresholds enforced (70%/90%)
- ✅ Security vulnerabilities detected
- ✅ Fast feedback (5-8 minutes)
- ✅ Professional development workflow

## Success Metrics

**Quality:** 48+ automated checks on every commit
**Speed:** 5-8 minute feedback loop
**Coverage:** Automated enforcement (70%/90%)
**Security:** Weekly scans + PR dependency review
**Compatibility:** 6-way matrix (3 OS × 2 Go versions)
**Visibility:** Badges, artifacts, Security tab integration

## Common Gotchas

**Concurrency Issues:**
- Race detector may catch issues only in CI (timing-dependent)
- Run `go test -race` locally before pushing

**Platform-Specific Code:**
- File path separators differ (use `filepath` package)
- Line endings may cause format issues (configure git autocrlf)

**Coverage Calculation:**
- Only runs on Linux + Go 1.23 (canonical platform)
- Download artifacts to see detailed coverage report
- Math package has stricter threshold (90% vs 70%)

**Linter Strictness:**
- Some linters may be too strict for learning code
- Adjust `.golangci.yml` if needed (document why)
- Don't disable important linters just to pass CI

**Artifact Retention:**
- Build artifacts auto-delete after 7 days
- Download before expiration if needed
- Coverage reports kept for 30 days

## Reference

**Workflow Files:**
- `.github/workflows/ci.yml` - Main CI pipeline
- `.github/workflows/security.yml` - Security scanning

**Configuration:**
- `.golangci.yml` - Linter settings

**Documentation:**
- This file (comprehensive reference)
- CLAUDE.md (development commands and workflow)
- [Development Roadmap](12-development-roadmap.md) (Phase 0 details)

**External Resources:**
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [golangci-lint Linters](https://golangci-lint.run/usage/linters/)
- [govulncheck Guide](https://go.dev/doc/security/vulncheck)
