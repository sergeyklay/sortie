---
status: accepted
date: 2026-03-17
decision-makers: Serghei Iakovlev
---

# Use Go as Core Runtime

## Context and Problem Statement

Sortie is a long-running orchestration service that polls an issue tracker, creates isolated
per-issue workspaces, dispatches coding agent sessions, and monitors their execution through
retries, reconciliation, and observability. The architecture is informed by
[OpenAI Symphony](https://github.com/openai/symphony) and adapted for multi-agent,
multi-tracker extensibility.

The core runtime must support concurrent subprocess management, context-driven cancellation,
predictable memory behavior under sustained load, and zero-dependency deployment to developer
machines, CI environments, and remote SSH hosts.

## Decision Drivers

1. **Runtime fitness.** Concurrency primitives, subprocess lifecycle management, and memory
   predictability under sustained load are non-negotiable requirements.
2. **Deployment simplicity.** Minimizing runtime dependencies on target hosts directly
   reduces operational burden.
3. **Agent-assisted development.** The codebase will be primarily written and maintained by
   AI coding agents. The stack must produce correct, idiomatic output from current-generation
   models with minimal iteration cycles.
4. **Long-term maintainability.** The project targets public release. Technology choices must
   favor a large contributor pool, stable tooling, and predictable upgrade paths.

## Considered Options

- Go
- Node.js/TypeScript
- Elixir/OTP
- Rust

## Decision Outcome

Chosen option: **Go**, because it provides the best combination of runtime fitness,
deployment simplicity, and long-term maintainability for this workload.

**Concurrency model.** Goroutines and channels are the native abstraction for the
orchestrator's workload: spawn a goroutine per agent session, use `context.Context` for
cancellation propagation, coordinate through typed channels. Go's scheduler distributes
work across OS threads without application-level coordination.

**Subprocess management.** `os/exec.CommandContext` integrates process lifecycle with
context cancellation: if a ticket moves to a terminal state, `cancel()` propagates through
the process tree. Signal handling and process cleanup on all target platforms (Linux, macOS,
Windows via SSH) are well-tested in the standard library.

**Deployment.** `go build` produces a single statically-linked binary with zero runtime
dependencies. Cross-compilation is built in. For SSH worker hosts, deployment reduces to
copying one file.

**Memory predictability.** The orchestrator is a 24/7 daemon. Go's garbage collector has
tunable controls (`GOGC`, `GOMEMLIMIT`) and goroutine stacks start at 2-8 KB.

**Agent generation quality.** Current LLMs produce higher pass@1 rates on TypeScript than
Go. This gap is real but transient: it narrows with each model generation, and Go's
uniformity (`gofmt`, single error handling idiom, minimal stylistic variation) partially
compensates by reducing the space for inconsistent output. The runtime characteristics
above are permanent architectural properties; the generation quality difference is a
snapshot of the current moment.

### Considered Options in Detail

**Node.js/TypeScript.** Strongest AI generation quality today. Best ecosystem for GraphQL
clients and template engines. However, the single-threaded event loop serializes all
orchestration logic; heavy JSON parsing or token accounting blocks stall detection and
reconciliation. Deployment requires a runtime on every target host. Long-running daemon
reliability demands active memory management. These are properties of the execution model,
not risks that agents can mitigate through better code.

**Elixir/OTP.** OpenAI's Symphony reference implementation uses Elixir, and BEAM's supervision
trees are an ideal fit for this workload. However, the ecosystem is small, the contributor
pool is narrow, and LLM generation quality for Elixir is significantly below Go or
TypeScript. Elixir is the right choice for a team of BEAM experts; it is the wrong default
for a public project targeting broad adoption.

**Rust.** Superior safety guarantees and performance. However, the borrow checker creates
long iteration cycles for agent-written code. For a project where development speed and
contributor accessibility matter, Rust's costs outweigh its benefits at this stage.
