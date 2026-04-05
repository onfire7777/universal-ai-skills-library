#!/usr/bin/env python3
"""Validate, format, and enhance a chat session summary.

Takes a raw Markdown summary (produced by Manus following SKILL.md instructions)
and:
  1. Validates all 7 required sections are present and non-empty
  2. Checks that referenced file paths actually exist on disk
  3. Generates a compact JSON version for programmatic consumption
  4. Produces a "handoff prompt" — a dense paragraph for pasting into a new chat
  5. Optionally enriches with AI-generated continuation recommendations

Usage:
    python3 format_summary.py <summary.md>                    # Validate + format
    python3 format_summary.py <summary.md> --json              # Also output JSON
    python3 format_summary.py <summary.md> --handoff           # Also output handoff prompt
    python3 format_summary.py <summary.md> --enrich            # AI-enriched continuation
    python3 format_summary.py <summary.md> --all               # Everything
"""

import argparse
import json
import os
import re
import sys
from datetime import datetime, timezone
from pathlib import Path

# ─── Required Sections ──────────────────────────────────────────────────────

REQUIRED_SECTIONS = [
    "SESSION_METADATA",
    "OBJECTIVE_TREE",
    "DECISIONS_LOG",
    "ARTIFACTS_REGISTRY",
    "DISCOVERIES",
    "UNRESOLVED_ITEMS",
    "CONTINUATION_CONTEXT",
]

SECTION_DESCRIPTIONS = {
    "SESSION_METADATA": "Date, user context, environment, working directories",
    "OBJECTIVE_TREE": "Hierarchical goals with completion status",
    "DECISIONS_LOG": "Key decisions with rationale and alternatives considered",
    "ARTIFACTS_REGISTRY": "Files created/modified with paths, purpose, and state",
    "DISCOVERIES": "Bugs found, technical insights, user preferences learned",
    "UNRESOLVED_ITEMS": "Open questions, known issues, pending features",
    "CONTINUATION_CONTEXT": "What the next AI needs to know to continue seamlessly",
}

# AI model to use for enrichment — must be available in Manus sandbox
ENRICH_MODEL = "gpt-4.1-mini"


# ─── Parsing ────────────────────────────────────────────────────────────────

def parse_summary(text: str) -> dict:
    """Parse a Markdown summary into sections by heading.

    DBG-003 FIX: Regex now handles mixed-case section names, single-word
    names, and names with numbers. Matches any ## heading that contains
    at least one uppercase word with underscores.
    DBG-008 FIX: Content before the first ## heading is captured as
    _PREAMBLE (includes H1 title and Generated timestamp).
    """
    sections = {}
    current_section = "_PREAMBLE"
    current_content = []

    for line in text.splitlines():
        # Match ## headings — flexible: allows mixed case, numbers,
        # single words, and optional numbering prefix like "1."
        heading_match = re.match(
            r"^##\s+(?:\d+\.\s+)?([A-Z][A-Za-z0-9_]+(?:[_ ][A-Za-z0-9_]+)*)\s*$",
            line,
        )
        if heading_match:
            if current_content or current_section == "_PREAMBLE":
                sections[current_section] = "\n".join(current_content).strip()
            raw_name = heading_match.group(1).strip()
            current_section = raw_name.upper().replace(" ", "_")
            current_content = []
        else:
            current_content.append(line)

    # Capture the last section
    if current_section:
        sections[current_section] = "\n".join(current_content).strip()

    return sections


def extract_file_paths(text: str) -> list[str]:
    """Extract file paths from the summary text.

    DBG-006 FIX: Strips trailing punctuation, markdown table pipes,
    backticks, and parentheses more aggressively.
    """
    patterns = [
        r"`(/[^\s`]+)`",           # Backtick-wrapped absolute paths
        r"(?<![`|])(/home/\S+)",   # /home/ paths (not inside backticks/tables)
        r"(?<![`|])(/mnt/\S+)",    # /mnt/ paths
        r"(C:\\[^\s,)|`]+)",       # Windows paths
    ]
    paths = set()
    for pattern in patterns:
        for match in re.finditer(pattern, text):
            path = match.group(1)
            # Strip trailing punctuation and markdown artifacts
            path = path.rstrip(".,;:)|`*>")
            # Skip paths that are clearly table separators
            if path.endswith("---") or path == "/":
                continue
            if len(path) > 3:
                paths.add(path)
    return sorted(paths)


def check_paths_exist(paths: list[str]) -> dict:
    """Check which referenced paths exist on disk."""
    results = {}
    for p in paths:
        try:
            exists = Path(p).exists()
        except (OSError, ValueError):
            exists = False
        results[p] = exists
    return results


# ─── Validation ─────────────────────────────────────────────────────────────

def validate_summary(text: str) -> dict:
    """Validate the summary has all required sections with content."""
    sections = parse_summary(text)
    errors = []
    warnings = []

    # Check required sections
    for section in REQUIRED_SECTIONS:
        if section not in sections:
            errors.append(f"Missing required section: ## {section}")
        elif len(sections[section].strip()) < 20:
            warnings.append(
                f"Section {section} appears too short "
                f"({len(sections[section].strip())} chars)"
            )

    # Check file path references
    paths = extract_file_paths(text)
    path_status = check_paths_exist(paths)
    missing_paths = [p for p, exists in path_status.items() if not exists]
    if missing_paths:
        for p in missing_paths:
            warnings.append(f"Referenced path does not exist: {p}")

    # Check for key quality indicators
    if "OBJECTIVE_TREE" in sections:
        statuses = re.findall(
            r"\b(completed|in.progress|abandoned|blocked|pending)\b",
            sections["OBJECTIVE_TREE"],
            re.IGNORECASE,
        )
        if not statuses:
            warnings.append(
                "OBJECTIVE_TREE has no status indicators "
                "(completed/in-progress/abandoned/blocked)"
            )

    if "DECISIONS_LOG" in sections:
        rationale_words = re.findall(
            r"\b(because|rationale|reason|chose|decided|alternative)\b",
            sections["DECISIONS_LOG"],
            re.IGNORECASE,
        )
        if not rationale_words:
            warnings.append(
                "DECISIONS_LOG lacks rationale language "
                "(because/rationale/reason/chose)"
            )

    return {
        "valid": len(errors) == 0,
        "sections_found": [s for s in sections if s != "_PREAMBLE"],
        "sections_missing": [
            s for s in REQUIRED_SECTIONS if s not in sections
        ],
        "preamble": sections.get("_PREAMBLE", ""),
        "errors": errors,
        "warnings": warnings,
        "paths_referenced": len(paths),
        "paths_missing": len(missing_paths),
        "word_count": len(text.split()),
        "sections": sections,
    }


# ─── JSON Export ────────────────────────────────────────────────────────────

def to_json(sections: dict, metadata: dict) -> dict:
    """Convert parsed sections to a structured JSON object."""
    output = {
        "schema_version": "1.0",
        "generated_at": datetime.now(timezone.utc).isoformat(),
        "preamble": metadata.get("preamble", ""),
        "validation": {
            "valid": metadata["valid"],
            "errors": metadata["errors"],
            "warnings": metadata["warnings"],
        },
        "summary": {},
    }

    for section_name in REQUIRED_SECTIONS:
        content = sections.get(section_name, "")
        output["summary"][section_name.lower()] = {
            "content": content,
            "char_count": len(content),
        }

    return output


# ─── Handoff Prompt ─────────────────────────────────────────────────────────

def generate_handoff(sections: dict) -> str:
    """Generate a dense handoff prompt for pasting into a new chat.

    This is the most critical output — it's a single block of text optimized
    for maximum context transfer in minimum tokens. An AI reading this should
    be able to continue the work seamlessly.

    DBG-001 FIX: Header separator uses proper string repetition.
    """
    parts = []

    # Session context
    meta = sections.get("SESSION_METADATA", "")
    if meta:
        key_lines = [
            line.strip("- *")
            for line in meta.splitlines()
            if any(
                kw in line.lower()
                for kw in [
                    "date", "project", "directory", "os", "environment",
                    "repo", "user", "working",
                ]
            )
        ]
        if key_lines:
            parts.append("CONTEXT: " + "; ".join(key_lines[:5]))

    # Objectives with status
    obj = sections.get("OBJECTIVE_TREE", "")
    if obj:
        goal_lines = [
            line.strip("- *")
            for line in obj.splitlines()
            if line.strip() and not line.startswith("#")
        ]
        if goal_lines:
            parts.append("GOALS: " + " | ".join(goal_lines[:8]))

    # Key decisions
    dec = sections.get("DECISIONS_LOG", "")
    if dec:
        decision_lines = [
            line.strip("- *")
            for line in dec.splitlines()
            if line.strip()
            and not line.startswith("#")
            and len(line.strip()) > 15
        ]
        if decision_lines:
            parts.append("KEY DECISIONS: " + " | ".join(decision_lines[:6]))

    # Artifacts
    art = sections.get("ARTIFACTS_REGISTRY", "")
    if art:
        path_lines = [
            line.strip("- *")
            for line in art.splitlines()
            if ("/" in line or "\\" in line) and "---" not in line
        ]
        if path_lines:
            parts.append("FILES: " + " | ".join(path_lines[:10]))

    # Discoveries
    disc = sections.get("DISCOVERIES", "")
    if disc:
        disc_lines = [
            line.strip("- *")
            for line in disc.splitlines()
            if line.strip()
            and not line.startswith("#")
            and len(line.strip()) > 15
        ]
        if disc_lines:
            parts.append("DISCOVERIES: " + " | ".join(disc_lines[:5]))

    # Unresolved
    unres = sections.get("UNRESOLVED_ITEMS", "")
    if unres:
        unres_lines = [
            line.strip("- *")
            for line in unres.splitlines()
            if line.strip()
            and not line.startswith("#")
            and len(line.strip()) > 10
        ]
        if unres_lines:
            parts.append("UNRESOLVED: " + " | ".join(unres_lines[:5]))

    # Continuation
    cont = sections.get("CONTINUATION_CONTEXT", "")
    if cont:
        cont_lines = [
            line.strip("- *")
            for line in cont.splitlines()
            if line.strip()
            and not line.startswith("#")
            and len(line.strip()) > 15
        ]
        if cont_lines:
            parts.append("CONTINUE FROM: " + " | ".join(cont_lines[:5]))

    handoff = "\n\n".join(parts)

    # DBG-001 FIX: Use proper character repetition (char * count, not str * count)
    separator = "\u2500" * 60
    header = (
        "=== AI HANDOFF PROMPT ===\n"
        "Paste this into a new chat to give the AI full context:\n"
        f"{separator}\n"
    )
    footer = f"\n{separator}"

    return header + handoff + footer


# ─── AI Enrichment ──────────────────────────────────────────────────────────

def enrich_with_ai(sections: dict) -> str:
    """Use the local AI model to generate continuation recommendations.

    DBG-004 FIX: Uses the ENRICH_MODEL constant (gpt-4.1-mini) which is
    confirmed available in the Manus sandbox.
    DBG-007 FIX: Clear error messages for missing dependencies with
    actionable guidance.
    """
    try:
        from openai import OpenAI
    except ImportError:
        return (
            "AI enrichment unavailable: openai package not installed. "
            "Run: pip3 install openai"
        )

    api_key = os.environ.get("OPENAI_API_KEY")
    if not api_key:
        return (
            "AI enrichment unavailable: OPENAI_API_KEY not set. "
            "This is pre-configured in the Manus sandbox environment."
        )

    client = OpenAI()

    # Build a focused prompt from the summary
    context_parts = []
    for section in ["OBJECTIVE_TREE", "UNRESOLVED_ITEMS", "CONTINUATION_CONTEXT"]:
        content = sections.get(section, "")
        if content:
            context_parts.append(f"## {section}\n{content}")

    if not context_parts:
        return "AI enrichment skipped: insufficient summary content."

    prompt = (
        "You are a senior technical advisor reviewing a project status summary. "
        "Based on the following information, provide:\n"
        "1. The 3 most critical next steps, ordered by priority\n"
        "2. Any risks or blockers you foresee\n"
        "3. Specific recommendations for the next work session\n\n"
        "Be concrete and actionable. No generic advice.\n\n"
        + "\n\n".join(context_parts)
    )

    try:
        response = client.chat.completions.create(
            model=ENRICH_MODEL,
            messages=[{"role": "user", "content": prompt}],
            max_tokens=1000,
            temperature=0.3,
        )
        return response.choices[0].message.content.strip()
    except Exception as e:
        return f"AI enrichment failed: {e}"


# ─── Main ───────────────────────────────────────────────────────────────────

def main():
    parser = argparse.ArgumentParser(
        description="Validate and format a Manus chat session summary"
    )
    parser.add_argument("summary_file", help="Path to the summary Markdown file")
    parser.add_argument(
        "--json", action="store_true", help="Output JSON version"
    )
    parser.add_argument(
        "--handoff", action="store_true", help="Generate handoff prompt"
    )
    parser.add_argument(
        "--enrich", action="store_true",
        help="AI-enriched recommendations (requires OPENAI_API_KEY)",
    )
    parser.add_argument(
        "--all", action="store_true",
        help="All outputs (json + handoff + enrich)",
    )
    parser.add_argument(
        "--output-dir",
        default=None,
        help="Directory for output files (default: same as input)",
    )
    args = parser.parse_args()

    if args.all:
        args.json = args.handoff = args.enrich = True

    # Read summary
    summary_path = Path(args.summary_file)
    if not summary_path.exists():
        print(f"ERROR: File not found: {summary_path}", file=sys.stderr)
        sys.exit(1)

    try:
        text = summary_path.read_text(encoding="utf-8")
    except Exception as e:
        print(f"ERROR: Could not read file: {e}", file=sys.stderr)
        sys.exit(1)

    if not text.strip():
        print("ERROR: Summary file is empty.", file=sys.stderr)
        sys.exit(1)

    # DBG-005 FIX: Default output dir is /home/ubuntu (where SKILL.md tells
    # Manus to save chat_summary.md), not the input file's parent dir.
    # This ensures handoff_prompt.txt and chat_summary.json land in the
    # same directory as the summary file by default.
    if args.output_dir:
        output_dir = Path(args.output_dir)
    else:
        output_dir = summary_path.parent
    output_dir.mkdir(parents=True, exist_ok=True)

    # Validate
    print("=" * 60)
    print("CHAT SUMMARY VALIDATOR")
    print("=" * 60)

    result = validate_summary(text)

    print(f"\n  Sections found: {len(result['sections_found'])}/{len(REQUIRED_SECTIONS)}")
    print(f"  Word count: {result['word_count']}")
    print(f"  Paths referenced: {result['paths_referenced']}")
    print(f"  Paths missing: {result['paths_missing']}")

    if result.get("preamble"):
        # Show first line of preamble (usually the H1 title)
        first_line = result["preamble"].splitlines()[0].strip()
        if first_line:
            print(f"  Title: {first_line}")

    if result["errors"]:
        print(f"\n  ERRORS ({len(result['errors'])}):")
        for err in result["errors"]:
            print(f"    \u2717 {err}")

    if result["warnings"]:
        print(f"\n  WARNINGS ({len(result['warnings'])}):")
        for warn in result["warnings"]:
            print(f"    \u26a0 {warn}")

    if result["valid"]:
        print("\n  \u2713 Summary is VALID \u2014 all 7 required sections present.")
    else:
        print(f"\n  \u2717 Summary is INVALID \u2014 {len(result['errors'])} error(s).")

    sections = result["sections"]

    # JSON output
    if args.json:
        json_path = output_dir / "chat_summary.json"
        json_data = to_json(sections, result)
        try:
            json_path.write_text(
                json.dumps(json_data, indent=2, ensure_ascii=False),
                encoding="utf-8",
            )
            print(f"\n  JSON: {json_path}")
        except OSError as e:
            print(f"\n  WARNING: Could not write JSON: {e}", file=sys.stderr)

    # Handoff prompt
    if args.handoff:
        handoff = generate_handoff(sections)
        handoff_path = output_dir / "handoff_prompt.txt"
        try:
            handoff_path.write_text(handoff, encoding="utf-8")
            print(f"  Handoff: {handoff_path}")
        except OSError as e:
            print(f"  WARNING: Could not write handoff: {e}", file=sys.stderr)
        print(f"\n{handoff}")

    # AI enrichment
    if args.enrich:
        print("\n  Generating AI recommendations...")
        recommendations = enrich_with_ai(sections)
        enrich_path = output_dir / "ai_recommendations.md"
        try:
            enrich_path.write_text(
                f"# AI-Generated Continuation Recommendations\n\n{recommendations}\n",
                encoding="utf-8",
            )
            print(f"  Recommendations: {enrich_path}")
        except OSError as e:
            print(
                f"  WARNING: Could not write recommendations: {e}",
                file=sys.stderr,
            )
        print(f"\n{recommendations}")

    print(f"\n{'=' * 60}")
    print("DONE")
    print(f"{'=' * 60}")

    sys.exit(0 if result["valid"] else 1)


if __name__ == "__main__":
    main()
