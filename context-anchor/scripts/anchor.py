#!/usr/bin/env python3
"""
Context Anchor — Persistent session focus management for Manus AI.

Commands:
    set     Create or overwrite the context anchor
    show    Display the current anchor
    check   Score a message/task against the anchor for relevance
    update  Modify the anchor without full replacement
    history Show anchor change history
    clear   Remove the anchor

Usage:
    python3 anchor.py set "Core topic description"
    python3 anchor.py set --file /path/to/topic.txt
    python3 anchor.py set "topic" --objectives "obj1" "obj2" --boundaries "not X" --success "criteria"
    python3 anchor.py show
    python3 anchor.py check "Is this sub-task relevant?"
    python3 anchor.py update --add-objective "New objective"
    python3 anchor.py update --add-boundary "Stay away from X"
    python3 anchor.py update --refine "Narrower focus description"
    python3 anchor.py history
    python3 anchor.py clear
"""

import argparse
import json
import os
import re
import sys
from datetime import datetime
from pathlib import Path

ANCHOR_FILE = os.path.expanduser("~/.context_anchor.md")
HISTORY_FILE = os.path.expanduser("~/.context_anchor_history.json")

# Model for AI-powered relevance checking (configurable via env var)
ANCHOR_MODEL = os.environ.get("ANCHOR_MODEL", "gpt-4.1-nano")


def _load_history():
    """Load anchor change history."""
    if os.path.exists(HISTORY_FILE):
        try:
            with open(HISTORY_FILE, "r") as f:
                return json.load(f)
        except (json.JSONDecodeError, IOError):
            return []
    return []


def _save_history(entries):
    """Save anchor change history."""
    try:
        with open(HISTORY_FILE, "w") as f:
            json.dump(entries, f, indent=2)
    except IOError as e:
        print(f"Warning: Could not save history: {e}", file=sys.stderr)


def _log_change(action, summary):
    """Log an anchor change to history."""
    history = _load_history()
    history.append({
        "timestamp": datetime.now().isoformat(),
        "action": action,
        "summary": summary[:200]
    })
    # Keep last 50 entries
    if len(history) > 50:
        history = history[-50:]
    _save_history(history)


def _read_anchor():
    """Read the current anchor file. Returns content or None."""
    if not os.path.exists(ANCHOR_FILE):
        return None
    try:
        with open(ANCHOR_FILE, "r") as f:
            return f.read().strip()
    except IOError as e:
        print(f"Error reading anchor: {e}", file=sys.stderr)
        return None


def _write_anchor(content):
    """Write content to the anchor file."""
    try:
        with open(ANCHOR_FILE, "w") as f:
            f.write(content + "\n")
        return True
    except IOError as e:
        print(f"Error writing anchor: {e}", file=sys.stderr)
        return False


def _strip_markdown(text):
    """Strip markdown formatting from text for clean tokenization."""
    return re.sub(r'[#*`\-|>_\[\]()]', ' ', text.lower())


def _insert_into_section(content, section_header, new_item, numbered=False):
    """Insert an item at the end of a markdown section, before the next ## heading.

    Args:
        content: Full document content
        section_header: The ## header to find (e.g., "## Key Objectives")
        new_item: The text to insert (without prefix)
        numbered: If True, prefix with next number; if False, prefix with "- "

    Returns:
        Modified content string, or None if section not found
    """
    lines = content.split("\n")
    new_lines = []
    in_section = False
    inserted = False
    section_items = 0

    for line in lines:
        if line.strip() == section_header:
            in_section = True
            new_lines.append(line)
            continue

        if in_section and not inserted:
            if line.startswith("##"):
                # End of section — insert before next heading
                if numbered:
                    new_lines.append(f"{section_items + 1}. {new_item}")
                else:
                    new_lines.append(f"- {new_item}")
                new_lines.append("")
                new_lines.append(line)
                inserted = True
                in_section = False
            else:
                new_lines.append(line)
                # Count items in this section only
                stripped = line.strip()
                if stripped:
                    if numbered and stripped[0].isdigit() and ". " in stripped:
                        section_items += 1
                    elif not numbered and stripped.startswith("- "):
                        section_items += 1
        else:
            new_lines.append(line)

    # If section was the last one (no subsequent ##)
    if in_section and not inserted:
        if numbered:
            new_lines.append(f"{section_items + 1}. {new_item}")
        else:
            new_lines.append(f"- {new_item}")
        new_lines.append("")

    return "\n".join(new_lines)


def cmd_set(args):
    """Set the context anchor."""
    # Get the core topic
    if args.file:
        try:
            with open(args.file, "r") as f:
                topic = f.read().strip()
        except (IOError, FileNotFoundError) as e:
            print(f"Error reading file: {e}", file=sys.stderr)
            sys.exit(1)
    elif args.topic:
        topic = " ".join(args.topic)
    else:
        print("Error: Provide a topic string or --file path.", file=sys.stderr)
        sys.exit(1)

    if not topic:
        print("Error: Topic cannot be empty.", file=sys.stderr)
        sys.exit(1)

    # Build the anchor document
    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    lines = [
        "# Context Anchor",
        "",
        f"<!-- anchor-timestamp -->{timestamp}",
        "",
        "## Core Topic",
        "",
        topic,
        "",
    ]

    if args.objectives:
        lines.append("## Key Objectives")
        lines.append("")
        for i, obj in enumerate(args.objectives, 1):
            lines.append(f"{i}. {obj}")
        lines.append("")

    if args.boundaries:
        lines.append("## Boundaries (Out of Scope)")
        lines.append("")
        for b in args.boundaries:
            lines.append(f"- {b}")
        lines.append("")

    if args.success:
        lines.append("## Success Criteria")
        lines.append("")
        for s in args.success:
            lines.append(f"- {s}")
        lines.append("")

    lines.append("## Recitation Reminder")
    lines.append("")
    lines.append("Re-read this file every 5-10 tool calls to maintain focus alignment.")
    lines.append("Before starting any new sub-task, verify it serves the Core Topic above.")
    lines.append("")

    content = "\n".join(lines)

    if _write_anchor(content):
        _log_change("set", topic[:200])
        print("=" * 60)
        print("  CONTEXT ANCHOR SET")
        print("=" * 60)
        print(content)
        print("=" * 60)
        print(f"  Anchor saved to: {ANCHOR_FILE}")
        print("  Manus will re-read this file periodically to stay focused.")
        print("=" * 60)
    else:
        sys.exit(1)


def cmd_show(args):
    """Display the current anchor."""
    content = _read_anchor()
    if content is None:
        print("No context anchor is currently set.")
        print("Use: python3 anchor.py set \"Your core topic\"")
        sys.exit(0)

    print("=" * 60)
    print("  CURRENT CONTEXT ANCHOR")
    print("=" * 60)
    print(content)
    print("=" * 60)


def cmd_check(args):
    """Check a message/task against the anchor for relevance."""
    content = _read_anchor()
    if content is None:
        print("No context anchor set. Cannot check relevance.")
        sys.exit(1)

    message = " ".join(args.message)
    if not message:
        print("Error: Provide a message to check.", file=sys.stderr)
        sys.exit(1)

    # Try AI-powered relevance check
    api_key = os.environ.get("OPENAI_API_KEY")
    if api_key:
        try:
            from openai import OpenAI
            client = OpenAI(timeout=15.0)
            response = client.chat.completions.create(
                model=ANCHOR_MODEL,
                messages=[
                    {"role": "system", "content": (
                        "You are a relevance checker. Given a context anchor (the core topic/purpose "
                        "of a session) and a new message or task, evaluate how relevant the message "
                        "is to the anchor. Respond with ONLY a JSON object:\n"
                        '{"relevance_score": 0-10, "alignment": "aligned|tangential|divergent", '
                        '"reasoning": "brief explanation", "recommendation": "proceed|caution|redirect"}'
                    )},
                    {"role": "user", "content": (
                        f"CONTEXT ANCHOR:\n{content}\n\n"
                        f"MESSAGE TO CHECK:\n{message}"
                    )}
                ],
                temperature=0.1,
                max_tokens=300
            )
            result_text = response.choices[0].message.content.strip()
            # Strip code fences if present
            if result_text.startswith("```"):
                result_text = "\n".join(result_text.split("\n")[1:])
                if result_text.endswith("```"):
                    result_text = result_text[:-3]
                result_text = result_text.strip()

            try:
                result = json.loads(result_text)
            except json.JSONDecodeError:
                result = {"raw_response": result_text}

            print("=" * 60)
            print("  RELEVANCE CHECK")
            print("=" * 60)
            if "relevance_score" in result:
                score = result["relevance_score"]
                alignment = result.get("alignment", "unknown")
                reasoning = result.get("reasoning", "")
                recommendation = result.get("recommendation", "")

                bar_filled = min(10, max(0, int(score)))
                bar_empty = 10 - bar_filled
                bar = "█" * bar_filled + "░" * bar_empty

                print(f"  Score:          [{bar}] {score}/10")
                print(f"  Alignment:      {alignment}")
                print(f"  Recommendation: {recommendation}")
                print(f"  Reasoning:      {reasoning}")
            else:
                print(f"  Result: {json.dumps(result, indent=2)}")
            print("=" * 60)
            return

        except ImportError:
            print("Warning: openai package not available. Falling back to keyword check.", file=sys.stderr)
        except Exception as e:
            err_msg = str(e)
            # Redact API keys from error messages
            err_msg = re.sub(r'sk-[a-zA-Z0-9]{20,}', 'sk-***REDACTED***', err_msg)
            print(f"Warning: AI check failed ({err_msg}). Falling back to keyword check.", file=sys.stderr)

    # Fallback: keyword overlap check (relative to message size)
    clean_anchor = _strip_markdown(content)
    clean_message = _strip_markdown(message)

    anchor_words = set(clean_anchor.split())
    message_words = set(clean_message.split())

    # Remove common stop words
    stop_words = {"the", "a", "an", "is", "are", "was", "were", "be", "been", "being",
                  "have", "has", "had", "do", "does", "did", "will", "would", "could",
                  "should", "may", "might", "shall", "can", "to", "of", "in", "for",
                  "on", "with", "at", "by", "from", "as", "into", "through", "during",
                  "before", "after", "above", "below", "between", "and", "but", "or",
                  "not", "no", "nor", "so", "yet", "both", "either", "neither", "each",
                  "every", "all", "any", "few", "more", "most", "other", "some", "such",
                  "than", "too", "very", "just", "about", "up", "out", "if", "then",
                  "it", "its", "this", "that", "these", "those", "i", "me", "my", "we",
                  "our", "you", "your", "he", "him", "his", "she", "her", "they", "them"}
    anchor_keywords = anchor_words - stop_words
    message_keywords = message_words - stop_words

    # DBG-003 fix: measure overlap relative to message size, not anchor size
    if not message_keywords:
        overlap_pct = 0
    else:
        overlap = anchor_keywords & message_keywords
        overlap_pct = len(overlap) / len(message_keywords) * 100

    score = min(10, int(overlap_pct / 10))
    if score >= 7:
        alignment = "aligned"
        recommendation = "proceed"
    elif score >= 4:
        alignment = "tangential"
        recommendation = "caution"
    else:
        alignment = "divergent"
        recommendation = "redirect"

    bar_filled = score
    bar_empty = 10 - bar_filled
    bar = "█" * bar_filled + "░" * bar_empty

    print("=" * 60)
    print("  RELEVANCE CHECK (keyword-based fallback)")
    print("=" * 60)
    print(f"  Score:          [{bar}] {score}/10")
    print(f"  Alignment:      {alignment}")
    print(f"  Recommendation: {recommendation}")
    print(f"  Keyword overlap: {overlap_pct:.0f}%")
    print("=" * 60)


def cmd_update(args):
    """Update the anchor without full replacement."""
    content = _read_anchor()
    if content is None:
        print("No context anchor set. Use 'set' first.", file=sys.stderr)
        sys.exit(1)

    modified = False

    if args.add_objective:
        obj_text = " ".join(args.add_objective)
        if "## Key Objectives" in content:
            result = _insert_into_section(content, "## Key Objectives", obj_text, numbered=True)
            if result:
                content = result
        else:
            content += f"\n## Key Objectives\n\n1. {obj_text}\n"
        modified = True
        _log_change("add_objective", obj_text)

    if args.add_boundary:
        boundary_text = " ".join(args.add_boundary)
        if "## Boundaries (Out of Scope)" in content:
            result = _insert_into_section(content, "## Boundaries (Out of Scope)", boundary_text, numbered=False)
            if result:
                content = result
        else:
            content += f"\n## Boundaries (Out of Scope)\n\n- {boundary_text}\n"
        modified = True
        _log_change("add_boundary", boundary_text)

    if args.refine:
        new_topic = " ".join(args.refine)
        # Replace the Core Topic section content
        lines = content.split("\n")
        new_lines = []
        in_core = False
        replaced = False
        for line in lines:
            if line.strip() == "## Core Topic":
                in_core = True
                new_lines.append(line)
                continue
            if in_core and not replaced:
                if line.startswith("##"):
                    # End of Core Topic section — insert new topic
                    new_lines.append("")
                    new_lines.append(new_topic)
                    new_lines.append("")
                    new_lines.append(line)
                    replaced = True
                    in_core = False
                else:
                    # Skip ALL old topic content (blank and non-blank)
                    continue
            else:
                new_lines.append(line)
        # If Core Topic was the last section
        if in_core and not replaced:
            new_lines.append("")
            new_lines.append(new_topic)
            new_lines.append("")
        content = "\n".join(new_lines)
        modified = True
        _log_change("refine", new_topic)

    if not modified:
        print("No update flags provided. Use --add-objective, --add-boundary, or --refine.", file=sys.stderr)
        sys.exit(1)

    # Update timestamp using unique marker
    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    if "<!-- anchor-timestamp -->" in content:
        content = re.sub(
            r'<!-- anchor-timestamp -->.*',
            f'<!-- anchor-timestamp -->{timestamp} (updated)',
            content
        )
    elif "**Set:**" in content:
        # Legacy format support
        lines = content.split("\n")
        for i, line in enumerate(lines):
            if line.startswith("**Set:**"):
                lines[i] = f"<!-- anchor-timestamp -->{timestamp} (updated)"
                break
        content = "\n".join(lines)

    if _write_anchor(content):
        print("=" * 60)
        print("  CONTEXT ANCHOR UPDATED")
        print("=" * 60)
        print(content)
        print("=" * 60)
    else:
        sys.exit(1)


def cmd_history(args):
    """Show anchor change history."""
    history = _load_history()
    if not history:
        print("No anchor history found.")
        sys.exit(0)

    print("=" * 60)
    print("  CONTEXT ANCHOR HISTORY")
    print("=" * 60)
    for entry in history:
        ts = entry.get("timestamp", "unknown")
        action = entry.get("action", "unknown")
        summary = entry.get("summary", "")
        print(f"  [{ts}] {action}: {summary}")
    print("=" * 60)
    print(f"  Total changes: {len(history)}")
    print("=" * 60)


def cmd_clear(args):
    """Clear the context anchor."""
    if os.path.exists(ANCHOR_FILE):
        try:
            os.remove(ANCHOR_FILE)
            _log_change("clear", "Anchor removed")
            print("Context anchor cleared.")
        except IOError as e:
            print(f"Error clearing anchor: {e}", file=sys.stderr)
            sys.exit(1)
    else:
        print("No context anchor to clear.")


def main():
    parser = argparse.ArgumentParser(
        description="Context Anchor — Persistent session focus management for Manus AI",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=__doc__
    )
    subparsers = parser.add_subparsers(dest="command", help="Available commands")

    # set
    p_set = subparsers.add_parser("set", help="Set the context anchor")
    p_set.add_argument("topic", nargs="*", help="Core topic description")
    p_set.add_argument("--file", "-f", help="Read topic from a file")
    p_set.add_argument("--objectives", "-o", nargs="+", help="Key objectives for this session")
    p_set.add_argument("--boundaries", "-b", nargs="+", help="Out-of-scope boundaries")
    p_set.add_argument("--success", "-s", nargs="+", help="Success criteria")

    # show
    subparsers.add_parser("show", help="Display the current anchor")

    # check
    p_check = subparsers.add_parser("check", help="Check relevance of a message against the anchor")
    p_check.add_argument("message", nargs="+", help="Message or task to check")

    # update
    p_update = subparsers.add_parser("update", help="Update the anchor without replacement")
    p_update.add_argument("--add-objective", nargs="+", help="Add a new objective")
    p_update.add_argument("--add-boundary", nargs="+", help="Add a new boundary")
    p_update.add_argument("--refine", nargs="+", help="Refine the core topic")

    # history
    subparsers.add_parser("history", help="Show anchor change history")

    # clear
    subparsers.add_parser("clear", help="Clear the context anchor")

    args = parser.parse_args()

    if not args.command:
        parser.print_help()
        sys.exit(1)

    commands = {
        "set": cmd_set,
        "show": cmd_show,
        "check": cmd_check,
        "update": cmd_update,
        "history": cmd_history,
        "clear": cmd_clear,
    }

    commands[args.command](args)


if __name__ == "__main__":
    main()
