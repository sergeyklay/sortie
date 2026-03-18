#!/usr/bin/env python3
"""
Validate a project context file (AGENTS.md, CLAUDE.md, GEMINI.md) against
the golden rule and structural best practices.

Usage:
    validate_context_file.py <path-to-context-file> [--project-dir <path>]

Checks:
    - Line count (warn >60, error >100 for root files)
    - Anti-patterns: project structure descriptions, tech stack summaries,
      framework documentation, linter-enforceable style rules
    - Section structure (Commands, Gotchas, Boundaries)
    - Boundary subsections (Always, Ask first, Never)
    - Stale references to files that do not exist
    - Empty sections
    - Forbidden section names (Appendix, Quick Reference, Overview, Cheat Sheet)

Exit codes: 0 = pass (may have warnings), 1 = errors found
"""

import re
import sys
from pathlib import Path

# --- Anti-pattern detectors ---

# Phrases that indicate discoverable project structure descriptions
STRUCTURE_PATTERNS = [
    re.compile(r"##\s*project\s+structure", re.IGNORECASE),
    re.compile(r"##\s*directory\s+(structure|layout|tree)", re.IGNORECASE),
    re.compile(r"##\s*folder\s+structure", re.IGNORECASE),
    re.compile(r"^\s*(src|lib|app|test|docs)/\s+#", re.MULTILINE),
]

# Phrases that indicate discoverable tech stack summaries
STACK_PATTERNS = [
    re.compile(r"##\s*tech(nology)?\s+stack", re.IGNORECASE),
    re.compile(r"##\s*stack\s+overview", re.IGNORECASE),
    re.compile(
        r"^\s*-\s*(React|Vue|Angular|Next\.js|Express|Django|Flask|Laravel|Rails|Spring)"
        r"\s+(for|as|is)\s+",
        re.MULTILINE | re.IGNORECASE,
    ),
]

# Phrases that indicate framework documentation the agent already knows
FRAMEWORK_DOC_PATTERNS = [
    re.compile(r"use\s+(functional\s+)?components\s+with\s+hooks", re.IGNORECASE),
    re.compile(r"prefer\s+`?useState`?\s+for\s+local\s+state", re.IGNORECASE),
    re.compile(r"use\s+`?useEffect`?\s+for\s+side\s+effects", re.IGNORECASE),
    re.compile(r"follow\s+RESTful\s+conventions", re.IGNORECASE),
    re.compile(r"use\s+async/await\s+instead\s+of\s+callbacks", re.IGNORECASE),
    re.compile(r"##\s*(react|vue|angular|django|rails)\s+patterns", re.IGNORECASE),
]

# Phrases that indicate linter-enforceable style rules
LINTER_STYLE_PATTERNS = [
    re.compile(r"use\s+\d+-space\s+indentation", re.IGNORECASE),
    re.compile(r"use\s+single\s+quotes", re.IGNORECASE),
    re.compile(r"add\s+trailing\s+commas", re.IGNORECASE),
    re.compile(r"maximum\s+line\s+length", re.IGNORECASE),
    re.compile(r"##\s*style\s+guide", re.IGNORECASE),
    re.compile(r"follow\s+(ESLint|Prettier|Biome)\s+rules", re.IGNORECASE),
]

# Phrases that indicate codebase overview duplicating README
OVERVIEW_PATTERNS = [
    re.compile(r"##\s*about\s+this\s+project", re.IGNORECASE),
    re.compile(r"##\s*what\s+is\s+this", re.IGNORECASE),
    re.compile(
        r"this\s+(is\s+)?(a|an)\s+(microservice|platform|application|service|tool)"
        r"\s+(built|written|made)\s+with",
        re.IGNORECASE,
    ),
]

# Forbidden section names
FORBIDDEN_SECTIONS = [
    re.compile(r"##\s*appendix", re.IGNORECASE),
    re.compile(r"##\s*quick\s+reference", re.IGNORECASE),
    re.compile(r"##\s*overview", re.IGNORECASE),
    re.compile(r"##\s*cheat\s+sheet", re.IGNORECASE),
]

# Expected top-level sections
RECOMMENDED_SECTIONS = {"commands", "gotchas", "boundaries"}

# Expected boundary subsections
BOUNDARY_SUBSECTIONS = {"always", "ask first", "never"}


def validate(file_path: str, project_dir: str | None = None) -> list[dict]:
    """Validate context file. Returns list of {level, message, line} dicts."""
    issues: list[dict] = []
    path = Path(file_path).resolve()

    def error(msg: str, line: int = 0) -> None:
        issues.append({"level": "ERROR", "message": msg, "line": line})

    def warn(msg: str, line: int = 0) -> None:
        issues.append({"level": "WARN", "message": msg, "line": line})

    def info(msg: str, line: int = 0) -> None:
        issues.append({"level": "INFO", "message": msg, "line": line})

    # --- File existence ---
    if not path.exists():
        error(f"File not found: {path}")
        return issues

    content = path.read_text(encoding="utf-8")
    lines = content.splitlines()
    line_count = len(lines)

    # --- Line count ---
    if line_count > 100:
        error(
            f"File has {line_count} lines (max 100 for root context file). "
            "Split reference material into separate docs."
        )
    elif line_count > 60:
        warn(
            f"File has {line_count} lines. Aim for 30-60 lines. "
            "Every extra line dilutes the instructions that matter."
        )
    elif line_count < 5:
        warn(f"File has only {line_count} lines. Likely too short to be useful.")

    # --- Anti-pattern checks ---
    for i, line in enumerate(lines, 1):
        for pattern in STRUCTURE_PATTERNS:
            if pattern.search(line):
                error(
                    "Project structure description detected. "
                    "The agent reads the directory tree.",
                    i,
                )
                break

        for pattern in STACK_PATTERNS:
            if pattern.search(line):
                error(
                    "Tech stack summary detected. The agent reads config files.",
                    i,
                )
                break

        for pattern in FRAMEWORK_DOC_PATTERNS:
            if pattern.search(line):
                warn(
                    "Possible framework documentation. "
                    "The agent already knows framework conventions.",
                    i,
                )
                break

        for pattern in LINTER_STYLE_PATTERNS:
            if pattern.search(line):
                warn(
                    "Linter-enforceable style rule detected. "
                    "Use ESLint/Prettier/Biome, not the context file.",
                    i,
                )
                break

        for pattern in OVERVIEW_PATTERNS:
            if pattern.search(line):
                warn(
                    "Codebase overview detected. This likely duplicates the README.",
                    i,
                )
                break

        for pattern in FORBIDDEN_SECTIONS:
            if pattern.search(line):
                error(
                    f"Forbidden section name detected: '{line.strip()}'. "
                    "Do not use Appendix, Quick Reference, Overview, or Cheat Sheet.",
                    i,
                )
                break

    # --- Section structure ---
    headings = []
    for i, line in enumerate(lines, 1):
        match = re.match(r"^(#{1,3})\s+(.+)", line)
        if match:
            level = len(match.group(1))
            title = match.group(2).strip().lower()
            headings.append({"level": level, "title": title, "line": i})

    found_sections = {h["title"] for h in headings if h["level"] == 2}
    missing_sections = RECOMMENDED_SECTIONS - found_sections
    if missing_sections:
        info(
            f"Missing recommended sections: {', '.join(sorted(missing_sections))}. "
            "Omit only if nothing non-discoverable exists for that category."
        )

    # --- Boundary subsections ---
    has_boundaries = any(h["title"] == "boundaries" for h in headings)
    if has_boundaries:
        h3_titles = {h["title"] for h in headings if h["level"] == 3}
        missing_boundaries = BOUNDARY_SUBSECTIONS - h3_titles
        if missing_boundaries:
            warn(
                f"Boundaries section missing subsections: "
                f"{', '.join(sorted(missing_boundaries))}. "
                "Use Always / Ask first / Never."
            )

    # --- Empty sections ---
    for idx, heading in enumerate(headings):
        next_heading_line = (
            headings[idx + 1]["line"] if idx + 1 < len(headings) else line_count + 1
        )
        section_lines = lines[heading["line"] : next_heading_line - 1]
        content_lines = [
            ln for ln in section_lines if ln.strip() and not ln.startswith("#")
        ]
        if not content_lines:
            warn(
                f"Empty section: '{heading['title']}'. "
                "An empty section is worse than no section.",
                heading["line"],
            )

    # --- Stale file references ---
    if project_dir:
        proj = Path(project_dir).resolve()
        file_ref_pattern = re.compile(r"`([^`]+\.\w{1,5})`")
        for i, line in enumerate(lines, 1):
            for match in file_ref_pattern.finditer(line):
                ref = match.group(1)
                # Skip obvious non-paths
                if ref.startswith("http") or ref.startswith("@") or " " in ref:
                    continue
                ref_path = proj / ref
                if not ref_path.exists() and not any(proj.glob(f"**/{Path(ref).name}")):
                    warn(f"Referenced file may not exist: `{ref}`", i)

    # --- Summary ---
    error_count = sum(1 for i in issues if i["level"] == "ERROR")
    warn_count = sum(1 for i in issues if i["level"] == "WARN")
    info_count = sum(1 for i in issues if i["level"] == "INFO")

    info(
        f"Validation complete: {error_count} errors, "
        f"{warn_count} warnings, {info_count} info."
    )

    return issues


def main() -> None:
    if len(sys.argv) < 2:
        print(f"Usage: {sys.argv[0]} <path-to-context-file> [--project-dir <path>]")
        sys.exit(1)

    file_path = sys.argv[1]
    project_dir = None

    if "--project-dir" in sys.argv:
        idx = sys.argv.index("--project-dir")
        if idx + 1 < len(sys.argv):
            project_dir = sys.argv[idx + 1]

    issues = validate(file_path, project_dir)

    has_errors = False
    for issue in issues:
        level = issue["level"]
        msg = issue["message"]
        line = issue.get("line", 0)
        prefix = f"L{line}: " if line else ""
        symbol = {"ERROR": "✗", "WARN": "⚠", "INFO": "ℹ"}.get(level, " ")
        print(f"  {symbol} [{level}] {prefix}{msg}")

        if level == "ERROR":
            has_errors = True

    sys.exit(1 if has_errors else 0)


if __name__ == "__main__":
    main()
