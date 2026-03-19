#!/usr/bin/env python3
"""
Validate an Agent Skill directory against the agentskills.io specification.

Usage:
    validate_skill.py <path-to-skill-directory>

Checks:
    - SKILL.md exists and has valid YAML frontmatter
    - name: format, length, reserved words, directory match
    - description: non-empty, length, no XML tags, perspective
    - compatibility: length constraint
    - Body size (warns if over 500 lines)
    - File reference depth (warns if nested references)
    - Directory structure
    - Path separators (warns on backslashes)

Exit codes: 0 = pass, 1 = errors found
"""

import re
import sys
from pathlib import Path

NAME_PATTERN = re.compile(r"^[a-z0-9]+(-[a-z0-9]+)*$")
XML_TAG_PATTERN = re.compile(r"<[^>]+>")
# Matches Windows-style paths: drive letters (C:\) or backslash + path character
# Excludes common escape sequences: \n \t \r \b \f \v \" \' \\
WINDOWS_PATH_PATTERN = re.compile(
    r'[A-Za-z]:\\'  # Drive letter pattern (C:\, D:\, etc.)
    r'|'
    r'\\(?![ntrfvb"\'\\])[A-Za-z0-9_]'  # Backslash + path char, not escape sequence
)
MAX_NAME = 64
MAX_DESC = 1024
MAX_COMPAT = 500
MAX_BODY_LINES = 500
RESERVED_WORDS = ["anthropic", "claude"]


def parse_frontmatter(content: str):
    """Extract YAML frontmatter from markdown content."""
    if not content.startswith("---"):
        return None, content

    # Find closing delimiter: must appear alone on its own line
    lines = content.splitlines(True)  # Keep line endings
    if len(lines) < 2:
        return None, content

    offset = len(lines[0])  # Start after opening "---" line
    end = -1
    for line in lines[1:]:
        if line.strip() == "---":
            end = offset
            break
        offset += len(line)

    if end == -1:
        return None, content

    try:
        import yaml

        fm = yaml.safe_load(content[3:end])
        body = content[end + 3 :].strip()
        return fm, body
    except ImportError:
        # Fallback: basic parsing without PyYAML
        lines = content[3:end].strip().split("\n")
        fm = {}
        current_key = None
        for line in lines:
            if line.startswith("  ") and current_key:
                fm[current_key] = fm.get(current_key, "") + " " + line.strip()
            elif ":" in line:
                key, _, val = line.partition(":")
                key = key.strip()
                val = val.strip().strip(">").strip("|").strip()
                fm[key] = val
                current_key = key
        body = content[end + 3 :].strip()
        return fm, body
    except Exception:
        return None, content


def validate(skill_dir: str) -> list[dict]:
    """Validate skill directory. Returns list of {level, message} dicts."""
    issues = []
    path = Path(skill_dir).resolve()

    def error(msg):
        issues.append({"level": "ERROR", "message": msg})

    def warn(msg):
        issues.append({"level": "WARN", "message": msg})

    def info(msg):
        issues.append({"level": "INFO", "message": msg})

    if not path.is_dir():
        error(f"Not a directory: {path}")
        return issues

    skill_md = path / "SKILL.md"
    if not skill_md.exists():
        error("SKILL.md not found")
        return issues

    content = skill_md.read_text(encoding="utf-8")

    fm, body = parse_frontmatter(content)
    if fm is None:
        error("SKILL.md must start with YAML frontmatter (--- delimiters)")
        return issues

    if not isinstance(fm, dict):
        error("Frontmatter must be a YAML mapping")
        return issues

    # Validate name
    name = fm.get("name")
    if not name:
        error("Missing required field: name")
    else:
        name = str(name).strip()
        if len(name) > MAX_NAME:
            error(f"name exceeds {MAX_NAME} chars ({len(name)})")
        if not NAME_PATTERN.match(name):
            error(f"name '{name}' must be lowercase alphanumeric with single hyphens")
        if "--" in name:
            error("name cannot contain consecutive hyphens (--)")
        for word in RESERVED_WORDS:
            if word in name.lower():
                error(f"name contains reserved word '{word}'")
        if XML_TAG_PATTERN.search(name):
            error("name cannot contain XML tags")
        if name != path.name:
            warn(f"name '{name}' does not match directory name '{path.name}'")

    # Validate description
    desc = fm.get("description")
    if not desc:
        error("Missing required field: description")
    else:
        desc = str(desc).strip()
        if len(desc) == 0:
            error("description cannot be empty")
        elif len(desc) > MAX_DESC:
            error(f"description exceeds {MAX_DESC} chars ({len(desc)})")
        if XML_TAG_PATTERN.search(desc):
            error("description cannot contain XML tags")
        if desc.startswith(("I ", "I can", "You ", "You can")):
            warn("description should use third person, not first/second person")
        if "[TODO" in desc:
            warn("description contains [TODO - complete before using")

        # Check for trigger phrases
        trigger_keywords = ["use when", "use this", "trigger", "whenever", "also use"]
        has_trigger = any(kw in desc.lower() for kw in trigger_keywords)
        if not has_trigger and len(desc) < 200:
            warn(
                "description may lack trigger phrases - agents need explicit guidance on when to activate"
            )

    # Validate optional fields
    compat = fm.get("compatibility")
    if compat and len(str(compat)) > MAX_COMPAT:
        error(f"compatibility exceeds {MAX_COMPAT} chars")

    # Check body size
    body_lines = body.count("\n") + 1 if body else 0
    if body_lines > MAX_BODY_LINES:
        warn(
            f"SKILL.md body is {body_lines} lines (recommended: under {MAX_BODY_LINES}). Consider splitting into reference files."
        )
    else:
        info(f"SKILL.md body: {body_lines} lines")

    # Check for nested references (one level deep only)
    md_links = re.findall(r"\[([^\]]*)\]\(([^)]+)\)", body)
    for label, href in md_links:
        if href.startswith(("http://", "https://", "#")):
            continue
        ref_path = path / href
        if ref_path.exists() and ref_path.suffix == ".md":
            ref_content = ref_path.read_text(encoding="utf-8", errors="ignore")
            nested_links = re.findall(r"\[([^\]]*)\]\(([^)]+\.md)\)", ref_content)
            nested_local = [
                h for _, h in nested_links if not h.startswith(("http://", "https://"))
            ]
            if nested_local:
                warn(
                    f"Nested reference: {href} links to {', '.join(nested_local)}. Keep references one level deep from SKILL.md."
                )

    # Check for Windows-style paths in body
    if WINDOWS_PATH_PATTERN.search(body):
        warn(
            "Body may contain backslash paths. Use forward slashes for cross-platform compatibility."
        )

    # Check directory structure
    known_dirs = {"scripts", "references", "assets", "evals", "agents"}
    for item in path.iterdir():
        if (
            item.is_dir()
            and item.name not in known_dirs
            and not item.name.startswith(".")
        ):
            info(f"Non-standard directory: {item.name}/")

    # Check reference files have TOC if large
    refs_dir = path / "references"
    if refs_dir.is_dir():
        for ref_file in refs_dir.glob("*.md"):
            line_count = (
                ref_file.read_text(encoding="utf-8", errors="ignore").count("\n") + 1
            )
            if line_count > 100:
                ref_content = ref_file.read_text(encoding="utf-8", errors="ignore")
                has_toc = any(
                    kw in ref_content[:500].lower()
                    for kw in ["contents", "table of contents", "## toc"]
                )
                if not has_toc:
                    warn(
                        f"references/{ref_file.name} is {line_count} lines but has no table of contents"
                    )

    # Summary
    errors_count = sum(1 for i in issues if i["level"] == "ERROR")
    warns_count = sum(1 for i in issues if i["level"] == "WARN")

    if errors_count == 0:
        info(
            f"Validation passed ({warns_count} warning{'s' if warns_count != 1 else ''})"
        )
    else:
        info(
            f"Validation failed: {errors_count} error{'s' if errors_count != 1 else ''}, {warns_count} warning{'s' if warns_count != 1 else ''}"
        )

    return issues


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: validate_skill.py <path-to-skill-directory>")
        sys.exit(1)

    issues = validate(sys.argv[1])

    for issue in issues:
        level = issue["level"]
        msg = issue["message"]
        prefix = {"ERROR": "x", "WARN": "!", "INFO": "-"}[level]
        print(f"  [{prefix}] {msg}")

    has_errors = any(i["level"] == "ERROR" for i in issues)
    sys.exit(1 if has_errors else 0)
