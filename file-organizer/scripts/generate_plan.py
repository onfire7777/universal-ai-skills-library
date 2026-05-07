#!/usr/bin/env python3
"""
generate_plan.py - Generate an organization plan from a scan report.

Creates a structured plan JSON that maps files to their organized destinations
based on the chosen strategy. The plan is reviewed before execution.

Usage:
    python generate_plan.py <scan_report.json> --strategy <strategy> --dest <dest_dir> [--output <plan.json>]

Strategies:
    category   - Organize by file type category (Documents, Images, etc.)
    extension  - Organize by file extension (.pdf, .jpg, etc.)
    date       - Organize by year/month based on modified date
    hybrid     - Organize by category, then by date within each category

Safety: Read-only. Generates a plan JSON for review before execution.
"""

import argparse
import json
import os
import sys
from datetime import datetime, timezone
from pathlib import Path


def strategy_category(files: list, dest_dir: str) -> list:
    """Organize files into category folders."""
    actions = []
    for f in files:
        if f.get("is_junk", False):
            continue
        category = f["category"]
        name = f["name"]
        dest = os.path.join(category, name)
        actions.append({
            "source": f["path"],
            "destination": dest,
            "action": "move",
            "category": category,
        })
    return actions


def strategy_extension(files: list, dest_dir: str) -> list:
    """Organize files into extension-based folders."""
    actions = []
    for f in files:
        if f.get("is_junk", False):
            continue
        ext = f["extension"].lstrip(".").upper() or "NO_EXTENSION"
        name = f["name"]
        dest = os.path.join(ext, name)
        actions.append({
            "source": f["path"],
            "destination": dest,
            "action": "move",
            "extension_folder": ext,
        })
    return actions


def strategy_date(files: list, dest_dir: str) -> list:
    """Organize files into year/month folders."""
    actions = []
    for f in files:
        if f.get("is_junk", False):
            continue
        try:
            dt = datetime.fromisoformat(f["modified"])
            year = dt.strftime("%Y")
            month = dt.strftime("%m-%B")  # e.g., "03-March"
        except (ValueError, KeyError):
            year = "Unknown_Year"
            month = "Unknown_Month"

        name = f["name"]
        dest = os.path.join(year, month, name)
        actions.append({
            "source": f["path"],
            "destination": dest,
            "action": "move",
            "year": year,
            "month": month,
        })
    return actions


def strategy_hybrid(files: list, dest_dir: str) -> list:
    """Organize by category first, then by year within each category."""
    actions = []
    for f in files:
        if f.get("is_junk", False):
            continue
        category = f["category"]
        try:
            dt = datetime.fromisoformat(f["modified"])
            year = dt.strftime("%Y")
        except (ValueError, KeyError):
            year = "Unknown_Year"

        name = f["name"]
        dest = os.path.join(category, year, name)
        actions.append({
            "source": f["path"],
            "destination": dest,
            "action": "move",
            "category": category,
            "year": year,
        })
    return actions


STRATEGY_MAP = {
    "category": strategy_category,
    "extension": strategy_extension,
    "date": strategy_date,
    "hybrid": strategy_hybrid,
}


def resolve_destination_conflicts(actions: list) -> list:
    """Resolve conflicts where multiple files map to the same destination."""
    dest_map = {}
    resolved = []

    for action in actions:
        dest = action["destination"]
        if dest not in dest_map:
            dest_map[dest] = 0
            resolved.append(action)
        else:
            dest_map[dest] += 1
            stem = Path(dest).stem
            ext = Path(dest).suffix
            parent = str(Path(dest).parent)
            new_dest = os.path.join(parent, f"{stem}_{dest_map[dest]}{ext}")
            action["destination"] = new_dest
            action["conflict_resolved"] = True
            resolved.append(action)

    return resolved


def main():
    parser = argparse.ArgumentParser(description="Generate file organization plan.")
    parser.add_argument("scan_report", help="Path to scan report JSON")
    parser.add_argument("--strategy", "-s", required=True,
                        choices=list(STRATEGY_MAP.keys()),
                        help="Organization strategy")
    parser.add_argument("--dest", "-d", required=True,
                        help="Destination directory for organized files")
    parser.add_argument("--output", "-o", default=None, help="Output plan JSON")
    parser.add_argument("--exclude-categories", nargs="*", default=[],
                        help="Categories to exclude")
    parser.add_argument("--include-hidden", action="store_true",
                        help="Include hidden files")
    args = parser.parse_args()

    if not os.path.isfile(args.scan_report):
        print(f"Error: '{args.scan_report}' not found.", file=sys.stderr)
        sys.exit(1)

    with open(args.scan_report, "r", encoding="utf-8") as f:
        scan_report = json.load(f)

    files = scan_report.get("files", [])

    # Apply filters
    if args.exclude_categories:
        files = [f for f in files if f["category"] not in args.exclude_categories]
    if not args.include_hidden:
        files = [f for f in files if not f.get("is_hidden", False)]

    # Generate actions
    strategy_fn = STRATEGY_MAP[args.strategy]
    actions = strategy_fn(files, args.dest)
    actions = resolve_destination_conflicts(actions)

    plan = {
        "source_dir": scan_report.get("scan_target", ""),
        "dest_dir": os.path.abspath(args.dest),
        "strategy": args.strategy,
        "timestamp": datetime.now(tz=timezone.utc).isoformat(),
        "total_files": len(actions),
        "actions": actions,
    }

    output = json.dumps(plan, indent=2, ensure_ascii=False)
    if args.output:
        with open(args.output, "w", encoding="utf-8") as f:
            f.write(output)
        print(f"Organization plan saved to {args.output} ({len(actions)} files)", file=sys.stderr)
    else:
        print(output)


if __name__ == "__main__":
    main()
