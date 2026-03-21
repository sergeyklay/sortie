# Contributing to Sortie

## Project Layout

Sortie uses the standard Go `cmd/internal` pattern:

```
cmd/sortie/            — main entry point, CLI flag parsing
internal/agent/        — agent adapter implementations (claude-code, mock, etc.)
internal/config/       — typed config structs, defaults, $VAR resolution, validation
internal/domain/       — domain types: Issue, TrackerAdapter, AgentAdapter interfaces
internal/logging/      — structured logging utilities
internal/orchestrator/ — poll loop, dispatch, reconciliation, retry, state machine
internal/persistence/  — SQLite-backed durable storage (retry queues, run history, metrics)
internal/prompt/       — prompt template rendering (text/template, strict mode, FuncMap)
internal/registry/     — adapter registry for tracker and agent implementations
internal/server/       — HTTP API and dashboard
internal/tracker/      — tracker adapter implementations (jira, file, etc.)
internal/workflow/     — WORKFLOW.md loader (front matter + prompt body split), file watcher
internal/workspace/    — workspace creation, path safety, hook execution
```

The `internal/` directory enforces package-level encapsulation at the compiler
level — external consumers cannot import internal packages. Each architecture
component maps to one internal sub-package, keeping dependencies explicit and
testable in isolation.
