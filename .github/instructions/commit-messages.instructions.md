---
description: Commit message formatting rules following Conventional Commits specification
---

# Commit Message Format

Follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification for all commit messages.

## Structure

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

## Commit Types

| Type       | Use Case                                 |
| ---------- | ---------------------------------------- |
| `feat`     | New feature                              |
| `fix`      | Bug fix                                  |
| `docs`     | Documentation changes                    |
| `style`    | Formatting, no logic change              |
| `refactor` | Code restructuring, no behavior change   |
| `test`     | Adding or updating tests                 |
| `chore`    | Maintenance, dependencies, build changes |
| `perf`     | Performance improvements                 |
| `ci`       | CI/CD changes                            |
| `build`    | Build system changes                     |
| `revert`   | Reverting previous commit                |

## Description Rules

- Use imperative mood: "Add", "Fix", "Update", "Remove"
- **IMPORTANT:** Subject line must be under 72 characters
- No period at end
- Be specific about what changed
- Avoid vague words: "various", "some", "minor", "stuff"

## Body (Summary) Formatting

When to include a body:

- When the change cannot be fully explained in the subject line
- For non-trivial changes or behavior changes

Body formatting:

- Separate from subject with a blank line
- **IMPORTANT:** Wrap lines at 72 characters
- Explain **what** and **why**, not how
- Reference GitHub issues/Jira tasks if relevant (e.g., "Fixes #123", "Related to #456", "Resolves ABC-789")
- Can be multiple paragraphs

## Breaking Changes

Indicate breaking changes with:

1. Exclamation mark before colon: `feat!: remove deprecated API`
2. Footer: `BREAKING CHANGE: previous tokens are no longer valid`

## Common Scopes (This Project)

- `(deps)` - dependency updates
- `(deps-dev)` - dev dependency updates
- No scope - most other changes

## Examples

| Change             | Good                                               | Bad                |
| ------------------ | -------------------------------------------------- | ------------------ |
| New API endpoint   | `feat(api): add user preferences endpoint`         | `Added new stuff`  |
| Fix null reference | `fix: handle null user in auth middleware`         | `Fixed bug`        |
| Update README      | `docs: clarify installation steps for Docker`      | `Update README`    |
| Rename variable    | `refactor: rename userId to accountId for clarity` | `Refactoring`      |
| Upgrade dependency | `chore(deps): bump zod from 4.3.4 to 4.3.5`        | `Updated packages` |
| Add unit test      | `test: add coverage for PaymentService edge cases` | `Tests`            |

## Full Example

```plaintext
fix: resolve race condition in user auth

Prevents race condition during simultaneous login attempts.
Multiple concurrent authentications were causing session
corruption and authentication failures.

Fixes #123

BREAKING CHANGE: session tokens are now single-use only
```

## Anti-Patterns

| Pattern             | Problem            | Correct Form            |
| ------------------- | ------------------ | ----------------------- |
| `fix: added...`     | Past tense         | `fix: add...`           |
| `fix: adding...`    | Present continuous | `fix: add...`           |
| `fix: add feature.` | Trailing period    | `fix: add feature`      |
| `add feature`       | Missing type       | `feat: add feature`     |
| `fix: fix bug`      | Vague              | `fix: handle null in X` |
| 72+ char subject    | Too long           | Keep under 72 chars     |
| body line >72 chars | Hard to read       | Wrap at 72 chars        |
| `✨ feat: add`      | Emoji in subject   | `feat: add`             |

## Multi-line Messages

```bash
git commit -m "feat: add payment validation" -m "Validates amount is positive and currency is supported." -m "Closes #123"
```

## Language

ALWAYS write commit messages in **English** regardless of conversation language.
