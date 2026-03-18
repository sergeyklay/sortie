---
status: accepted
date: 2026-03-17
decision-makers: Serghei Iakovlev
---

# Use SQLite for Persistence

## Context and Problem Statement

The orchestrator must persist retry queues, session metadata, workspace registry, token
accounting, and run history across process restarts.
[The Symphony spec](https://github.com/openai/symphony/blob/main/SPEC.md) uses in-memory
state with no persistent database, accepting that this data is lost on restart. Sortie
corrects this.

The persistence layer must operate without external dependencies, match the orchestrator's
single-authority write pattern, and remain inspectable with standard tooling.

## Decision Drivers

1. **Zero external dependencies.** The orchestrator deploys as a single binary; adding a
   database server contradicts this constraint.
2. **Single-writer concurrency.** The orchestrator serializes all state mutations through
   one authority; the persistence layer must match this pattern without unnecessary
   complexity.
3. **Operational inspectability.** Operators must be able to query state directly during
   debugging and incident response.

## Considered Options

- SQLite (embedded, WAL mode)
- In-memory only (Symphony approach)
- PostgreSQL/MySQL
- Embedded key-value stores (BoltDB, BadgerDB)

## Decision Outcome

Chosen option: **SQLite in WAL mode**, because it provides durable state with zero external
dependencies while matching the orchestrator's single-writer pattern.

SQLite in WAL mode provides concurrent reads with a single writer. The entire state lives
in a single file scoped to the project directory (derived from the `WORKFLOW.md` location),
not alongside the binary. The default path is `.sortie.db` in the same directory as the
workflow file. This ensures per-project isolation (different projects never share state) and
avoids permission errors when the binary is installed in a system directory like
`/usr/local/bin/`. There is no external database to provision and no network dependency.

Go has mature SQLite bindings: `modernc.org/sqlite` (pure Go, no CGo) provides full SQLite
functionality without a C toolchain on the build host.

### Considered Options in Detail

**In-memory only (Symphony approach).** Simpler, but a process restart during active work
loses all retry state and requires a full cold start from the issue tracker.

**PostgreSQL/MySQL.** Adds an external dependency, connection management, and migration
tooling. Unnecessary for a single-process orchestrator that serializes all writes through
one authority.

**Embedded key-value stores (BoltDB, BadgerDB).** Viable, but SQLite provides relational
queries, schema migrations, and standard tooling (`sqlite3` CLI) that simplify debugging
and operational inspection.
