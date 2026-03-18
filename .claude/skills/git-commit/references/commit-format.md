# Commit Message Format

Conventional Commits specification: https://www.conventionalcommits.org/en/v1.0.0/

## Structure

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

## Types

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

## Rules

- Imperative mood: "Add", "Fix", "Update" — not "Added", "Fixing"
- Subject line under 72 characters, no trailing period
- Body: wrap at 72 chars, explain what and why (not how)
- Breaking changes: `feat!: description` or `BREAKING CHANGE:` footer
- ALWAYS in English regardless of conversation language

## Anti-Patterns

| Pattern             | Problem          | Fix                     |
| ------------------- | ---------------- | ----------------------- |
| `fix: added...`     | Past tense       | `fix: add...`           |
| `fix: fix bug`      | Vague            | `fix: handle null in X` |
| `add feature`       | Missing type     | `feat: add feature`     |
| `fix: add feature.` | Trailing period  | `fix: add feature`      |
| 72+ char subject    | Too long         | Shorten                 |

## Examples

| Change             | Message                                            |
| ------------------ | -------------------------------------------------- |
| New API endpoint   | `feat(api): add user preferences endpoint`         |
| Fix null reference | `fix: handle null user in auth middleware`          |
| Update README      | `docs: clarify installation steps for Docker`      |
| Rename variable    | `refactor: rename userId to accountId for clarity`  |
| Upgrade dependency | `chore(deps): bump zod from 4.3.4 to 4.3.5`       |
