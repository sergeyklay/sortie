---
name: test
description: Generate comprehensive test coverage for implemented features
argument-hint: Path to spec/plan, or description of what to test
agent: Tester
---

Generate test coverage for the implemented feature or changed code.

You MUST use the Agent tool with subagent_type="Tester".

Process:

1. Review the specification and implementation plan (if provided)
2. Analyze the actual implementation across all layers
3. Identify what requires test coverage (Services, Actions, Components, edge cases, regression risks, etc)
4. Generate appropriate tests following the project's testing conventions

Follow its testing guidelines and coverage requirements strictly.

$ARGUMENTS
