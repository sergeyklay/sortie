---
status: accepted
date: 2026-03-17
decision-makers: Serghei Iakovlev
---

# Use Adapter Interfaces for Integration Extensibility

## Context and Problem Statement

The orchestrator must integrate with external systems in two dimensions: issue trackers
(source of work) and coding agents (execution of work). The system must support multiple
implementations in each dimension without modifying core orchestration logic.

## Decision Drivers

1. **Multi-tracker support.** Organizations use different issue trackers (Jira, Linear,
   GitHub Issues). The orchestrator must not be coupled to any single tracker's API.
2. **Agent agnosticism.** The coding agent landscape is evolving rapidly. The orchestrator
   must treat agents as interchangeable runtimes behind a stable contract.
3. **Extensibility without rewrite.** Adding a new tracker or agent must be an additive
   change — a new package implementing an existing interface.

## Considered Options

- Adapter interfaces (Go interfaces per integration dimension)
- Plugin system (dynamic loading via `plugin` package or RPC)
- Direct integration (hardcoded clients)

## Decision Outcome

Chosen option: **Adapter interfaces**, because they provide compile-time safety and simple
extensibility without the operational complexity of dynamic plugin systems.

The orchestrator defines Go interfaces for issue tracker access and coding agent
communication. Each tracker (Jira, Linear, GitHub Issues, File System) and each agent
runtime (Claude Code, Codex, generic HTTP) is implemented as a separate package behind
its respective interface. Issue and event data is normalized into common types at the
adapter boundary.

The initial implementation targets Jira as the primary tracker and Claude Code as the
primary agent runtime. The agent adapter communicates over stdio, allowing straightforward
substitution of alternative runtimes (Codex, Copilot, or any agent exposing a compatible
CLI interface) without changes to orchestration logic.

Linear support and additional agent adapters are planned as subsequent implementations.

### Considered Options in Detail

**Plugin system.** Go's `plugin` package is limited to Linux and macOS, produces fragile
ABI coupling, and complicates the single-binary deployment model. An RPC-based plugin
system adds network overhead and failure modes unnecessary for a single-process
orchestrator.

**Direct integration.** Hardcoding tracker and agent clients into the orchestrator creates
tight coupling and requires modifying core code for every new integration. This contradicts
the extensibility driver.
