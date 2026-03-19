# Frontmatter Fields Reference

Complete guide to SKILL.md YAML frontmatter fields per the agentskills.io specification.

## Contents

- Required fields (name, description)
- Optional fields (license, compatibility, metadata)
- allowed-tools (experimental, with platform support matrix)
- Description writing strategy
- Platform-specific extensions

## name (required)

| Constraint  | Rule                                                   |
| ----------- | ------------------------------------------------------ |
| Length      | 1-64 characters                                        |
| Characters  | Lowercase alphanumeric and hyphens (`a-z`, `0-9`, `-`) |
| Start/end   | Cannot start or end with `-`                           |
| Consecutive | No `--` allowed                                        |
| Match       | Must match parent directory name                       |
| Reserved    | Cannot contain "anthropic" or "claude"                 |
| XML         | Cannot contain XML tags                                |

**Naming conventions**:

Prefer gerund form for clarity about the capability:

- `processing-pdfs`, `analyzing-data`, `testing-code`, `writing-documentation`

Acceptable alternatives:

- Noun phrases: `pdf-processing`, `spreadsheet-analysis`
- Action verbs: `process-pdfs`, `analyze-spreadsheets`

Avoid vague names: `helper`, `utils`, `tools`, `documents`, `data`, `files`.

## description (required)

| Constraint  | Rule                                       |
| ----------- | ------------------------------------------ |
| Length      | 1-1024 characters                          |
| Content     | Non-empty, no XML tags                     |
| Perspective | Third person (injected into system prompt) |

### Why descriptions matter

The description is the primary discovery mechanism. At startup, agents pre-load ONLY the name and description of every installed skill (~100 tokens each). When a user's request arrives, the agent scans all descriptions to decide which skill to activate. A poor description means your skill never triggers, regardless of how good the instructions are.

### Formula

```plaintext
[What the skill does] + [When to use it, with specific triggers and edge cases]
```

### Effective examples

```yaml
description: Extract text and tables from PDF files, fill forms, merge documents. Use when working with PDF files or when the user mentions PDFs, forms, or document extraction.
```

```yaml
description: Generate descriptive commit messages by analyzing git diffs. Use when the user asks for help writing commit messages, reviewing staged changes, or preparing a pull request description.
```

```yaml
description: "Apply Acme Corp brand guidelines to presentations and documents, including official colors (#FF6B35, #004E89), fonts (Montserrat, Open Sans), and logo usage rules. Use whenever creating external-facing materials, slide decks, reports, or any document representing Acme Corp."
```

### Making descriptions pushy

Agents tend to under-trigger skills. Be explicit about edge cases that should trigger activation:

```yaml
# Instead of:
description: Build dashboards to display data.

# Write:
description: Build fast dashboards to display data. Use this skill whenever the user mentions dashboards, data visualization, metrics, charts, graphs, or wants to display any kind of data visually, even if they don't explicitly ask for a "dashboard."
```

### Common mistakes

```yaml
description: Helps with PDFs.
# Too vague. No triggers. Agent won't know when to activate.

description: I can help you process Excel files and generate reports.
# First person. Descriptions are injected into system prompts in third person.

description: Processes data.
# No triggers. No specifics. Useless for discovery.
```

## license (optional)

Short license name or reference to bundled file:

```yaml
license: Apache-2.0
license: Proprietary. LICENSE.txt has complete terms
```

## compatibility (optional)

Max 500 characters. Only include when the skill has specific environment requirements:

```yaml
compatibility: Requires git, docker, jq, and access to the internet
compatibility: Designed for Claude Code (or similar products with filesystem access)
```

Most skills don't need this field.

## metadata (optional)

Arbitrary key-value mapping for additional properties. Use reasonably unique key names to avoid conflicts:

```yaml
metadata:
  author: your-org
  version: "1.0"
  category: development
```

Codex uses `metadata.short-description` as a user-facing description in the IDE.

## allowed-tools (optional, experimental)

Restricts which tools the agent can use while the skill is active. Use this for read-only skills (code review, architecture analysis) to prevent accidental modifications.

```yaml
allowed-tools:
  - Read
  - Grep
  - Glob
  - WebSearch
```

**Platform support matrix:**

| Platform          | Support Level |
| ----------------- | ------------- |
| Claude Code       | Full          |
| VS Code / Copilot | No            |
| Cursor            | No            |
| Gemini CLI        | Partial       |
| Codex CLI         | No            |

"Partial" means the field is recognized but behavior may differ from specification. Only include `allowed-tools` when you need to enforce safety constraints AND your target platform supports it.

## Platform-Specific Extensions

### OpenAI Codex

**Storage note**: Codex uses `.agents/skills/` (not `.codex/skills/`). Place skills in `.agents/skills/<skill-name>/` for project scope or `~/.agents/skills/<skill-name>/` for user scope.

**agents/openai.yaml**: Codex supports an optional `agents/openai.yaml` file alongside SKILL.md for UI metadata and invocation control:

```yaml
interface:
  display_name: "Human-Friendly Name"
  short_description: "Shown in Codex UI"
  icon_small: "./assets/icon.svg"
  icon_large: "./assets/icon-large.png"
  brand_color: "#3B82F6"
  default_prompt: "Optional surrounding prompt"
policy:
  allow_implicit_invocation: false # require explicit $ invocation
dependencies:
  tools:
    - type: "mcp"
      value: "serverName"
      description: "Required MCP server"
      transport: "streamable_http"
      url: "https://example.com/mcp"
```

When `allow_implicit_invocation` is false, Codex won't activate the skill based on user prompt alone - it requires explicit `$skill-name` invocation. This is useful for skills that should only run when deliberately chosen.

This file is Codex-specific and safely ignored by all other platforms. Include it only when building for Codex and needing UI customization or invocation control.
