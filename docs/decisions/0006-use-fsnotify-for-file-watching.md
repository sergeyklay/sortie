---
status: accepted
date: 2026-03-18
decision-makers: Serghei Iakovlev
---

# Use fsnotify for Filesystem Event Watching

## Context and Problem Statement

The orchestrator must detect changes to `WORKFLOW.md` at runtime and reload configuration
without a process restart (Section 6.2 — dynamic reload semantics). When a reload fails,
the last known-good configuration must be retained and an error logged.

The detection mechanism must work reliably across the editors and save patterns used in
practice — including the atomic write (write-to-temp + rename) strategy used by Vim, Emacs,
most JetBrains IDEs, and VS Code.

## Decision Drivers

1. **Reliability across editor save patterns.** Polling misses rename-based atomic writes
   unless the poll interval is very short; kernel event notification handles them natively.
2. **Zero external infrastructure.** File watching must work on any host: developer
   laptops, CI containers, remote SSH sessions — no inotify wrappers, daemons, or OS
   packages required at runtime.
3. **Single-binary deployment constraint.** Any addition must be pure Go with no CGo
   requirement (same constraint that drove ADR-0002).
4. **Maintenance burden.** A hand-rolled `syscall.InotifyInit` implementation requires
   per-platform branches (Linux inotify, macOS FSEvents/kqueue, Windows ReadDirectoryChanges);
   this is non-trivial maintenance surface for a small project.

## Considered Options

- `github.com/fsnotify/fsnotify` — pure-Go, cross-platform, idiomatic
- `os.Stat` polling loop — no dependency, but unreliable and high-latency
- Manual `syscall.InotifyInit` (Linux only) — no dependency, but not portable
- `golang.org/x/exp/filenotify` — experimental, not suitable for production

## Decision Outcome

Chosen option: **`github.com/fsnotify/fsnotify`**, because it is the de-facto standard
cross-platform file watching library for Go, is pure Go (no CGo), satisfies the
single-binary deployment constraint, and eliminates per-platform syscall maintenance.

`fsnotify` uses OS-native kernel event APIs (inotify on Linux, kqueue on macOS/BSDs,
ReadDirectoryChanges on Windows) behind a unified Go interface. The only transitive
dependency is `golang.org/x/sys`, which is a quasi-standard-library package maintained
by the Go team and already present in most Go dependency graphs.

The implementation watches the **parent directory** rather than the file itself. This
handles atomic-rename editor saves: editors write to a temp file and rename it into place,
which generates a `CREATE` event on the parent directory rather than a `WRITE` event on
the original file. Filtering by `filepath.Base` ensures only `WORKFLOW.md` events trigger
a reload.

### Considered Options in Detail

**`os.Stat` polling loop.** Avoids the dependency entirely. However, detecting a rename-based
atomic write requires polling at sub-100 ms intervals to be practically responsive, which
wastes CPU and still has a latency window. It also cannot distinguish a file being replaced
from a file being unchanged.

**Manual `syscall.InotifyInit` (Linux only).** Pure Go, no external dependency. However,
inotify is Linux-specific; separate implementations are required for macOS (kqueue/FSEvents)
and Windows (ReadDirectoryChanges). This is exactly what `fsnotify` abstracts — reimplementing
it would reproduce the library without its test coverage or community maintenance.

**`golang.org/x/exp/filenotify`.** Experimental (`x/exp`), API is unstable and not
recommended for production use.

## Consequences

- `go.mod` gains `github.com/fsnotify/fsnotify v1.9.0` (direct) and
  `golang.org/x/sys` (indirect, quasi-stdlib, Go-team maintained).
- The binary remains CGo-free and statically linkable.
- Kubernetes ConfigMap volume mounts use a double-symlink strategy that causes the parent
  directory watch to miss some update events. This is a known limitation documented in
  `Manager.Start`. Operators using ConfigMap mounts should call `Reload()` via an external
  trigger or accept the next poll-cycle latency as the effective reload mechanism.
