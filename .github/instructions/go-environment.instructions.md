---
name: 'Go Environment'
description: 'Makefile targets, asdf toolchain resolution, and constraints on direct Go binary invocation'
applyTo: '**/*.go,go.mod,go.sum,Makefile,.golangci.yml,.github/workflows/*.yml'
---

# Go Environment

Go is managed by asdf. The `go` binary resolves through `~/.asdf/shims/go`.

## Commands

Use Makefile targets for all build, test, and lint operations:

- Format: `make fmt`
- Lint: `make lint`
- Build: `make build`
- Run tests: `make test`
- Run package tests: `make test PKG=./internal/persistence`
- Run single test: `make test RUN=TestOpenStore`
- Run single test in package: `make test PKG=./internal/persistence RUN=TestOpenStore`

Read the Makefile to discover available targets before running any Go toolchain commands directly.

## Constraints

- NEVER prefix commands with `GOPATH=...`, `GOMODCACHE=...`, or any Go environment overrides. The asdf shim configures everything.
- NEVER run `go test`, `go build`, `go vet`, or `golangci-lint` directly. Use the corresponding `make` target.
- NEVER use `/usr/local/go/bin/go`, `/usr/bin/go`, or any absolute path to a Go binary.
- NEVER downgrade the `go` directive in `go.mod`. NEVER add or modify `toolchain` directives in `go.mod` unless explicitly asked.

## CI and Dockerfiles

When writing CI workflows or Dockerfiles, pin the Go version to match `.tool-versions`. Do not use `latest` or `1.x`.
