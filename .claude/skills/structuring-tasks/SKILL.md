---
name: structuring-tasks
description: >
  Transform raw feature descriptions, bug reports, or ideas into structured task
  documents with YAML frontmatter and standardized sections. Use when asked to
  'create a task', 'write a task file', 'structure this feature', 'turn this into
  a task', 'make a ticket', or when a .tasks/ directory exists and the user provides
  an unstructured feature description. Also use when refining or reformatting
  existing task files. Do NOT use for Jira ticket creation, spec writing, or
  implementation planning.
metadata:
  author: airSlate Inc.
  version: "1.0"
  category: project-management
---

# Structuring tasks

Convert unstructured feature descriptions into task documents that an implementer
can execute without asking follow-up questions. Each task file lives in `.tasks/`
and follows a fixed template with YAML frontmatter and seven sections.

## When to use

- A user pastes a raw idea, feature request, or bug report and wants a task file
- A `.tasks/` directory exists and needs a new entry
- An existing task file needs restructuring to match the standard format
- The user says "create a task", "structure this", "turn this into a task"

## Workflow

### Phase 1: Analyze the input

Read the raw description. Identify:

1. **Core change** -- What capability is added, modified, or fixed
2. **Affected files** -- Which files or directories will change
3. **Knowledge gaps** -- Technical details, API specs, or configuration syntax the
   implementer will need but the description does not provide

Do not generate the task yet. Gaps drive Phase 2.

If the input lacks enough detail to fill the frontmatter or the What section, ask
the user to clarify before proceeding. Focus on: complexity assessment, platform
scope, and dependency relationships between tasks.

### Phase 2: Research

Fill knowledge gaps before writing. For each gap:

- Search the web for official documentation, release notes, or API references
- Read relevant project files to understand current state
- Check existing `.tasks/` files for ID numbering and naming patterns

Research determines the content of the Technical reference section. Skip this phase
only when the input already contains complete technical details.

### Phase 3: Determine frontmatter values

Assign each field:

| Field | How to determine |
|-------|-----------------|
| `id` | Next sequential number in `.tasks/`, zero-padded to 3 digits (`"001"`, `"002"`) |
| `title` | Imperative sentence, under 80 characters. "Add X to Y", "Fix Z in W" |
| `status` | Always `draft` for new tasks |
| `created` | Today's date in `YYYY-MM-DD` format |
| `complexity` | `low` (single file, < 1 hour), `medium` (multiple files, new patterns), `high` (cross-cutting, new architecture) |
| `platform` | Primary platform affected: `github-copilot`, `cursor`, `claude-code`, `cross-platform`, or omit if stack-agnostic |
| `scope` | List of file paths that will be modified. Directories end with `/` |
| `depends-on` | List of task IDs that must complete first. Empty array if none |

### Phase 4: Write the task

Copy the template from `assets/task-template.md`. Fill each section following
these rules.

**Language constraint:** Write every part of the task document in English —
frontmatter values, section headings, body text, code comments, and acceptance
criteria. Translate the user's input into English during this phase regardless of
the language the user communicates in.

**Title line** (`# Task {ID}: {TITLE}`)
Matches the frontmatter `id` and `title` exactly.

**What**
One paragraph, 2-4 sentences. State the change, not the motivation. Mention
specific tools, APIs, or standards by name.

**Current state**
Bullet list of how the system behaves today. Each bullet is a factual observation
the implementer can verify. Include file names where relevant.

**Target state**
Bullet list of how the system should behave after implementation. Mirror the
structure of Current state so diffs are obvious.

**Why**
Bullet list of reasons. Lead with the strongest reason. Each bullet connects a
problem (current state) to a benefit (target state). No generic statements like
"improves developer experience" without saying how.

**Where: Files to modify**
Table with two columns: File (path) and Change (what changes). One row per file.
Be specific: "Add `vscode/askQuestions` to `tools:` array" not "Update tools".

**Where: Files to create**
Table with two columns: File (path) and Purpose (why it exists). Omit this
subsection if no files are created.

**Where: Files out of scope**
Bullet list of files that might seem relevant but should not be touched, with a
reason. Prevents scope creep.

**Technical reference**
The section the implementer reads when they need syntax, API details, or
configuration examples. Use subsections for distinct topics. Include code blocks
with exact syntax. This section comes from Phase 2 research. If the input already
contained complete technical details, reproduce them here with proper formatting.

**Acceptance criteria**
Checklist (`- [ ]`) of verifiable conditions. Each criterion answers "how do I
prove this is done?" Not "implement feature X" but "Feature X renders in the UI
when setting Y is enabled." Aim for 5-8 criteria.

**References**
Bullet list of links with context. Format:
`- [Title](URL) -- One-line reason this link matters`

### Phase 5: Validate

Before saving, verify:

- [ ] Frontmatter `id` does not collide with existing tasks in `.tasks/`
- [ ] Frontmatter `scope` matches the files listed in the Where section
- [ ] Every file in "Files to modify" exists in the project (or is explicitly
      marked as new in "Files to create")
- [ ] Technical reference contains no placeholder text
- [ ] Acceptance criteria are verifiable without reading the task author's mind
- [ ] No banned words from AGENTS.md content rules

## Filename convention

```
.tasks/{ID}-{slug}.md
```

- `{ID}` -- Zero-padded 3-digit number matching frontmatter `id`
- `{slug}` -- Lowercase, hyphen-separated, derived from the title
- Example: `.tasks/001-ask-user-questions.md`

## Complexity calibration

| Level | Signals |
|-------|---------|
| `low` | Single file change. No new patterns introduced. Implementer needs < 30 minutes. |
| `medium` | 2-6 files. May introduce a new tool, frontmatter field, or section. Requires some research. |
| `high` | 7+ files or cross-cutting concern. New architectural pattern. Multiple platforms affected. Requires significant research. |

## Anti-patterns

| Pattern | Problem | Fix |
|---------|---------|-----|
| Vague acceptance criteria | "It works" is not verifiable | State the observable behavior and how to trigger it |
| Missing Technical reference | Implementer searches docs during coding | Front-load the research into the task |
| Scope creep in Where | Listing files "just in case" | Only list files that require changes for this task |
| Copy-pasting the raw input as What | No analysis performed | Synthesize: extract the core change, discard noise |
| Placeholder references | "[TODO: find link]" | Research now or remove the reference |

## Language

All task output is in English. This is non-negotiable regardless of the language
the user writes in. Frontmatter, section content, code comments, acceptance
criteria, and references are always English. During Phases 1-2 (analysis,
clarification questions, research) respond in the user's language so
communication stays natural. Switch to English only in the written artifact.
