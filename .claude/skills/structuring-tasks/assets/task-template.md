---
id: "{ID}"
title: {TITLE}
status: draft
created: {DATE}
complexity: {low|medium|high}
platform: {PLATFORM}
scope:
  - {FILE_OR_DIR}
depends-on: []
---

# Task {ID}: {TITLE}

## What

{One paragraph describing the change. State what capability is added or modified.}

### Current state

- {How the system behaves today. One bullet per relevant observation.}

### Target state

- {How the system should behave after implementation. One bullet per change.}

## Why

- {Business or technical reason for the change. One bullet per reason.}

## Where

### Files to modify

| File | Change |
|------|--------|
| `{path}` | {What changes in this file.} |

### Files to create

| File | Purpose |
|------|---------|
| `{path}` | {Why this file is needed.} |

### Files out of scope

- `{path}` -- {Why this file should not be touched.}

## Technical reference

{API details, configuration syntax, protocol specifics, or external documentation summaries the implementer needs. Use subsections as needed. Include code blocks for syntax examples.}

## Acceptance criteria

- [ ] {Verifiable condition that proves the task is done.}

## References

- [{Title}]({URL}) -- {One-line context for why this link matters.}
