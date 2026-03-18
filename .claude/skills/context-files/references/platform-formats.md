# Platform-specific context file formats

Details on how each platform loads context files, their unique features, and bridge file setup.

## AGENTS.md (universal)

- **Platforms:** OpenAI Codex (native), Cursor (native), GitHub Copilot (opt-in), Gemini CLI (configurable)
- **Spec:** [agents.md](https://agents.md/) (Linux Foundation)
- **Hierarchy:** user-level → project root → subdirectories
- **Syntax:** Plain markdown. No import directives.

**Codex** searches `AGENTS.md` in three levels:

1. `~/.codex/AGENTS.md` (personal global guidance)
2. `AGENTS.md` at repository root
3. `AGENTS.md` in the current working directory

**Cursor** reads `AGENTS.md` from subdirectories. More specific nested files override parent directories. Also supports `.cursor/rules/*.md` with YAML frontmatter for structured rule scoping (always, intelligently, by glob, manually).

**GitHub Copilot** requires opt-in via `chat.useAgentsMdFile` setting. Enable `chat.useNestedAgentsMdFiles` for subfolder support. The primary file remains `.github/copilot-instructions.md` (always included). JetBrains, Eclipse, and Xcode only read `.github/copilot-instructions.md`.

## CLAUDE.md

- **Platform:** Claude Code
- **Hierarchy:** three levels
  1. `~/.claude/CLAUDE.md` (user-level, every project)
  2. `CLAUDE.md` at repository root (project-level)
  3. `CLAUDE.md` in the current working directory (directory-level)

**Unique features:**

- Supports `@path/to/file` import syntax (up to 5 hops deep)
- Injected inside `<system-reminder>` tags
- Content is conditionally relevant (model may deprioritize sections it considers irrelevant)
- `.claude/settings.json` for additional configuration

**Import tradeoff:** `@`-imported files load unconditionally on every session regardless of task relevance. Use `@` imports only for content that applies to every session (commands, global boundaries). For task-specific references, list files and let the agent decide what to read.

## GEMINI.md

- **Platform:** Gemini CLI
- **Default file name:** `GEMINI.md` (configurable in `.gemini/settings.json`)
- **Hierarchy:** three levels
  1. `~/.gemini/GEMINI.md` (global)
  2. `GEMINI.md` in project root and parent directories (workspace)
  3. Just-in-time loading when a tool accesses a file in a directory with `GEMINI.md`

**Configuration override:**

```json
{
  "context": {
    "fileName": ["AGENTS.md", "GEMINI.md"]
  }
}
```

**Unique features:**

- `@file.md` import syntax (similar to Claude Code)
- `/memory show` to inspect loaded context
- `/memory refresh` to reload after edits
- All found files are concatenated and sent with every prompt

## Bridge files

When a team uses multiple platforms, maintain one canonical file (`AGENTS.md`) and create bridges:

| Platform       | Bridge mechanism                                                                  |
| -------------- | --------------------------------------------------------------------------------- |
| Claude Code    | `CLAUDE.md` containing `@AGENTS.md`                                               |
| GitHub Copilot | `.github/copilot-instructions.md` containing `Refer to [AGENTS.MD](../AGENTS.md)` |
| Gemini CLI     | `.gemini/settings.json` with `"context": { "fileName": ["AGENTS.md"] }`           |
| Cursor         | Native support, no bridge needed                                                  |
| Codex          | Native support, no bridge needed                                                  |

Never duplicate content across platform files. One source of truth, references everywhere else. A bridge file that copies content from `AGENTS.md` will drift within weeks.

## Choosing your canonical file

| Scenario                     | Recommendation                        |
| ---------------------------- | ------------------------------------- |
| Team uses only Claude Code   | `CLAUDE.md` directly                  |
| Team uses only Gemini CLI    | `GEMINI.md` directly                  |
| Team uses multiple platforms | `AGENTS.md` + bridge files            |
| Team uses Copilot + Cursor   | `AGENTS.md` (native for both)         |
| Solo developer, any platform | Whichever the platform reads natively |
