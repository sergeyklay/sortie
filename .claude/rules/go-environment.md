---
paths:
  - "**/*.go"
  - "go.mod"
  - "go.sum"
  - "Makefile"
  - ".golangci.yml"
  - ".github/workflows/*.yml"
---

# Go Environment

- Go version is managed by asdf (see `.tool-versions`).
- The `go` binary is at `~/.asdf/shims/go` (asdf shim). Do NOT use `/usr/local/go/bin/go` or any other system-installed Go binary.
- Always run `go` commands without an explicit path — the asdf shim resolves the correct version.
- When writing CI workflows or Dockerfiles, pin to `go <version>` — do not use `latest` or `1.x`.
- The `go` directive in `go.mod` is `go <version>`. Do not downgrade it.
- Do not add `toolchain` directives to `go.mod` unless explicitly asked.
- Use Go language features freely (range-over-func, generic type aliases, etc.) - all that supported in the version specified in `.tool-versions`. Do not use features from newer versions.

Replace '<version>' with the version specified in `.tool-versions`.
