---
name: implement
description: Execute implementation plan or implement feature directly
argument-hint: Path to plan file, or feature description if working without plan
agent: Coder
---

Implement the requested feature.

You MUST use the Agent tool with subagent_type="Coder".

If a plan exists: Execute it strictly phase-by-phase, following all architectural constraints.

If no plan exists: Analyze the request, identify required changes across all layers, and implement atomically following the architecture guidelines.

Follow its implementation rules strictly.

$ARGUMENTS
