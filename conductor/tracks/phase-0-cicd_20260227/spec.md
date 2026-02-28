# Specification: Phase 0 — CI/CD Infrastructure

**Track ID:** phase-0-cicd_20260227
**Type:** Chore
**Status:** Archived (Complete)
**Completed:** 2025-10-10

## Summary

Set up automated CI/CD infrastructure before writing any application code — establishing a quality gate that enforces formatting, linting, testing, and security scanning on every commit and pull request.

## Context

This was the first phase of the `go_toy_renderer` project. Infrastructure was established before any rendering code was written, ensuring all subsequent phases start with a green pipeline and that quality standards are enforced automatically.

## Acceptance Criteria

- [x] GitHub Actions workflow runs on push to `main` and all PRs
- [x] Multi-platform build matrix: Linux, macOS, Windows on Go 1.24 & 1.25
- [x] `golangci-lint v2` configured with 30+ linters (`.golangci.yml`)
- [x] Coverage enforcement: fails if <70% overall or <90% for `pkg/math`
- [x] Security scanning via `govulncheck` integrated into CI
- [x] Pre-commit hook at `.githooks/pre-commit` (fmt → vet → lint)

## Source Reference

- Task summary: [`.claude/tasks/2025-10-10_project-setup-and-documentation.md`](../../../../.claude/tasks/2025-10-10_project-setup-and-documentation.md)

---

_Archived. This track is a historical record of completed work._
