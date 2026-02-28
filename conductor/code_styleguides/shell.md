# Shell Scripting Style Guide — go_toy_renderer

## Scope

Shell scripts in this project:
- `.githooks/pre-commit` — pre-commit hook (bash)
- Any future CI helper scripts under `.github/scripts/`

## Shell Selection

- Use **bash** (not sh) for hooks and CI scripts — the project targets environments where bash is available (Linux, macOS, Git Bash on Windows).
- Declare explicitly: `#!/usr/bin/env bash`.

## Defensive Defaults

Every script must start with:
```bash
#!/usr/bin/env bash
set -euo pipefail
```

| Flag | Effect |
|------|--------|
| `-e` | Exit immediately on error |
| `-u` | Treat unset variables as errors |
| `-o pipefail` | Propagate errors through pipes |

## Variable Quoting

- Always double-quote variables: `"$VAR"`, `"$@"`, `"${array[@]}"`.
- Never use unquoted `$*`.
- Use `${VAR:-default}` for optional variables with defaults.

## Functions

- Name functions with lowercase and underscores: `run_linter`, `check_format`.
- Keep functions small and single-purpose.
- Local variables inside functions: use `local`.

```bash
run_linter() {
  local target="${1:-.}"
  golangci-lint run "$target"
}
```

## Error Messages

- Print errors to stderr: `echo "error: something failed" >&2`.
- Use meaningful exit codes (0 = success, 1 = general error).

## Command Checks

Check for required tools before using them:
```bash
if ! command -v golangci-lint &>/dev/null; then
  echo "error: golangci-lint not found. See https://golangci-lint.run/usage/install/" >&2
  exit 1
fi
```

## Pre-commit Hook Pattern

The `.githooks/pre-commit` hook follows this structure:
```bash
#!/usr/bin/env bash
set -euo pipefail

echo "Running go fmt..."
go fmt ./...

echo "Running go vet..."
go vet ./...

echo "Running golangci-lint..."
golangci-lint run

echo "Pre-commit checks passed."
```

- No `set +e` or error suppression.
- Each step prints what it's doing.
- Exits non-zero on any failure (which aborts the commit).

## Cross-Platform Notes

- Use forward slashes in paths, even on Windows (Git Bash normalises them).
- Avoid Windows-specific commands (`cmd.exe`, `PowerShell`); keep scripts portable.
- Line endings: scripts must use LF (`\n`), not CRLF. Configure git: `git config core.autocrlf false` for script files.

## ShellCheck

Run ShellCheck on any new shell script before committing:
```bash
shellcheck .githooks/pre-commit
```

Fix all warnings — do not suppress without a comment explaining why.
