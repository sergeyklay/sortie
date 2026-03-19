#!/usr/bin/env python3
"""
Initialize a new Agent Skill directory from template.

Usage:
    init_skill.py <skill-name> --path <output-directory>

Examples:
    init_skill.py processing-pdfs --path .claude/skills
    init_skill.py analyzing-data --path ~/.cursor/skills
    init_skill.py my-workflow --path .codex/skills
"""

import re
import sys
from pathlib import Path

NAME_PATTERN = re.compile(r"^[a-z0-9]+(-[a-z0-9]+)*$")
MAX_NAME_LENGTH = 64
RESERVED_WORDS = ["anthropic", "claude"]


def validate_name(name: str) -> list[str]:
    """Validate skill name against specification. Returns list of errors."""
    errors = []
    if len(name) > MAX_NAME_LENGTH:
        errors.append(f"Name exceeds {MAX_NAME_LENGTH} characters ({len(name)})")
    if not NAME_PATTERN.match(name):
        errors.append(
            "Name must be lowercase alphanumeric with single hyphens, no leading/trailing hyphens"
        )
    if "--" in name:
        errors.append("Consecutive hyphens (--) not allowed")

    name_lower = name.lower().strip()
    for word in RESERVED_WORDS:
        if word in name_lower:
            errors.append(f"Reserved word '{word}' not allowed in name")
    return errors


def title_case(name: str) -> str:
    return " ".join(w.capitalize() for w in name.split("-"))


SKILL_MD = """---
name: {name}
description: >
  [TODO: What this skill does AND when to use it.
  Include specific trigger phrases, file types, and task contexts.
  Be slightly pushy - agents tend to under-trigger.
  Max 1024 characters.]
---

# {title}

[TODO: 1-2 sentence overview. What does this skill enable?]

## [First section]

[TODO: Choose a structure that fits your skill:
- Workflow: Step-by-step sequential process
- Task-based: Grouped by operation type
- Reference: Organized by domain
- Conditional: Decision tree with branches

Then delete this guidance block.]
"""


def init_skill(name: str, path: str) -> bool:
    errors = validate_name(name)
    if errors:
        for e in errors:
            print(f"  Error: {e}")
        return False

    skill_dir = Path(path).expanduser().resolve() / name

    if skill_dir.exists():
        print(f"  Error: directory already exists: {skill_dir}")
        return False

    try:
        skill_dir.mkdir(parents=True)
        (skill_dir / "SKILL.md").write_text(
            SKILL_MD.format(name=name, title=title_case(name))
        )

        print(f"  Created: {skill_dir}/")
        print("  Created: SKILL.md")
        print()
        print("Next steps:")
        print("  1. Edit SKILL.md - write the description and instructions")
        print("  2. Add scripts/, references/, or assets/ directories as needed")
        print("  3. Validate: python validate_skill.py " + str(skill_dir))
        return True

    except Exception as e:
        print(f"  Error: {e}")
        return False


if __name__ == "__main__":
    if len(sys.argv) < 4 or sys.argv[2] != "--path":
        print("Usage: init_skill.py <skill-name> --path <output-directory>")
        print()
        print("Name rules: lowercase, hyphens, max 64 chars")
        print("Prefer gerund form: processing-pdfs, analyzing-data")
        print()
        print("Examples:")
        print("  init_skill.py processing-pdfs --path .claude/skills")
        print("  init_skill.py my-workflow --path ~/.codex/skills")
        sys.exit(1)

    name = sys.argv[1]
    path = sys.argv[3]

    print(f"Initializing skill: {name}")
    print(f"Location: {path}")
    print()

    sys.exit(0 if init_skill(name, path) else 1)
