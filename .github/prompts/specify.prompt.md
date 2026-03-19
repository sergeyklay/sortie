---
name: specify
description: Transform a feature request into a detailed technical specification
argument-hint: Describe the feature or problem to specify
agent: Architect
---

Transform the provided feature request into a comprehensive technical specification.

Before writing anything, read the relevant sections of [architecture.md](../../docs/architecture.md) — this is the authoritative specification for all domain models, state machines, algorithms, and validation rules. Your spec must conform to it; do not invent behavior that contradicts the architecture document.

Also review [TODO.md](../../TODO.md) to understand where this feature fits in the milestone sequence and what dependencies exist.

Apply coding standards from: [Go documentation guidelines](../instructions/go-documentation.instructions.md)

${input:request:Describe the feature or problem to specify}
