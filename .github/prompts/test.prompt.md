---
name: test
description: Generate comprehensive test coverage for implemented features
argument-hint: Path to spec/plan, or description of what to test
agent: Tester
---

Generate test coverage for the implemented feature or changed code.

## Process

1. **Review the specification and implementation plan** (if provided) to understand the intended behavior and contracts.
2. **Analyze the actual implementation** across all layers — read the source files before writing any tests.
3. **Identify what requires test coverage:** services, domain logic, state transitions, adapters, edge cases, regression risks, error paths, and concurrency safety.
4. **Generate appropriate tests** following the project's testing conventions.

Apply coding standards from: [Go documentation guidelines](../instructions/go-documentation.instructions.md) and [Go environment guidelines](../instructions/go-environment.instructions.md)

${input:request:Path to spec/plan or description of what to test}
