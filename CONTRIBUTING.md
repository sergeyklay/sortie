# Contributing to Sortie

## Project Layout

Sortie uses the standard Go `cmd/internal` pattern:

```
cmd/sortie/            — main entry point, CLI flag parsing
internal/domain/       — domain types: Issue, TrackerAdapter, AgentAdapter interfaces
internal/config/       — workflow loader, typed config, prompt template rendering
internal/orchestrator/ — poll loop, dispatch, reconciliation, retry, state machine
internal/tracker/      — tracker adapter implementations (jira, file, etc.)
internal/agent/        — agent adapter implementations (claude-code, mock, etc.)
internal/workspace/    — workspace creation, path safety, hook execution
internal/server/       — HTTP API and dashboard
```

The `internal/` directory enforces package-level encapsulation at the compiler
level — external consumers cannot import internal packages. Each architecture
component maps to one internal sub-package, keeping dependencies explicit and
testable in isolation.
