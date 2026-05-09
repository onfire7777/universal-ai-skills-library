#!/usr/bin/env python3
"""
generate_report.py - Generate a human-readable Markdown report from scan data.

Produces a comprehensive summary of the directory state, including category
breakdown, duplicate analysis, junk files, naming issues, and recommendations.

Usage:
    python generate_report.py <scan_report.json> [--duplicates <dup_report.json>] [--output <report.md>]

Safety: Read-only. Produces a Markdown report file.
"""

import argparse
import json
import os
import sys
from datetime import datetime, timezone


def human_size(nbytes: int) -> str:
    for unit in ("B", "KB", "MB", "GB", "TB"):
        if abs(nbytes) < 1024:
            return f"{nbytes:.1f} {unit}"
        nbytes /= 1024
    return f"{nbytes:.1f} PB"


def generate_report(scan: dict, duplicates: dict = None) -> str:
    """Generate a Markdown report from scan and duplicate data."""
    lines = []
    files = scan.get("files", [])

    lines.append("# File Organization Report")
    lines.append("")
    lines.append(f"**Scanned Directory:** `{scan.get('scan_target', 'N/A')}`")
    lines.append(f"**Scan Time:** {scan.get('scan_time', 'N/A')}")
    lines.append(f"**Total Files:** {scan.get('total_files', 0):,}")
    lines.append(f"**Total Size:** {scan.get('total_size_human', 'N/A')}")
    lines.append("")

    # Category breakdown
    lines.append("## Category Breakdown")
    lines.append("")
    categories = scan.get("categories", {})
    if categories:
        lines.append("| Category | Files | Size |")
        lines.append("|----------|------:|-----:|")
        sorted_cats = sorted(categories.items(), key=lambda x: x[1]["total_size"], reverse=True)
        for cat, info in sorted_cats:
            lines.append(f"| {cat} | {info['count']:,} | {info.get('total_size_human', human_size(info['total_size']))} |")
        lines.append("")

    # Duplicates
    if duplicates:
        lines.append("## Duplicate Files")
        lines.append("")
        lines.append(f"**Duplicate Groups:** {duplicates.get('confirmed_duplicate_groups', 0)}")
        lines.append(f"**Total Duplicate Files:** {duplicates.get('total_duplicate_files', 0)}")
        lines.append(f"**Wasted Space:** {duplicates.get('total_wasted_human', 'N/A')}")
        lines.append("")

        groups = duplicates.get("groups", [])
        if groups:
            lines.append("### Top Duplicate Groups")
            lines.append("")
            for i, g in enumerate(groups[:20], 1):
                lines.append(f"**Group {i}** ({g['count']} copies, {g['size_human']} each, {g['wasted_human']} wasted)")
                lines.append(f"- Keep: `{g['keep']['relative_path']}`")
                for rc in g["remove_candidates"]:
                    lines.append(f"- Remove: `{rc['relative_path']}`")
                lines.append("")

    # Junk files
    junk = scan.get("junk_files", [])
    if junk:
        lines.append("## Junk Files")
        lines.append("")
        lines.append(f"**Total Junk Files:** {len(junk)}")
        lines.append("")
        for j in junk[:50]:
            lines.append(f"- `{j}`")
        if len(junk) > 50:
            lines.append(f"- ... and {len(junk) - 50} more")
        lines.append("")

    # Empty files
    empty = scan.get("empty_files", [])
    if empty:
        lines.append("## Empty Files (0 bytes)")
        lines.append("")
        lines.append(f"**Total Empty Files:** {len(empty)}")
        lines.append("")
        for e in empty[:30]:
            lines.append(f"- `{e}`")
        if len(empty) > 30:
            lines.append(f"- ... and {len(empty) - 30} more")
        lines.append("")

    # Naming issues
    import re
    naming_issues = []
    for f in files:
        name = f["name"]
        issues = []
        if len(name) > 200:
            issues.append("excessively long name")
        if re.search(r"\s{2,}", name):
            issues.append("multiple consecutive spaces")
        if re.search(r"[_\-]{2,}", name):
            issues.append("multiple consecutive separators")
        if name.startswith(" ") or name.endswith(" "):
            issues.append("leading/trailing whitespace")
        if re.search(r"%[0-9a-fA-F]{2}", name):
            issues.append("URL-encoded characters")
        if re.search(r"\(\d+\)", name):
            issues.append("copy indicator (e.g., (1))")
        if re.search(r"copy\s*\d*", name, re.IGNORECASE):
            issues.append("copy marker in name")
        stem_lower = os.path.splitext(name)[0].lower()
        if re.match(r"^(untitled|document|file|image|photo|screenshot|img_?\d*|dsc_?\d*|new\s*file)$", stem_lower):
            issues.append("vague/generic name")
        if issues:
            naming_issues.append({"file": f["relative_path"], "issues": issues})

    if naming_issues:
        lines.append("## Naming Issues")
        lines.append("")
        lines.append(f"**Files with naming issues:** {len(naming_issues)}")
        lines.append("")
        lines.append("| File | Issues |")
        lines.append("|------|--------|")
        for ni in naming_issues[:50]:
            lines.append(f"| `{ni['file']}` | {', '.join(ni['issues'])} |")
        if len(naming_issues) > 50:
            lines.append(f"| ... | {len(naming_issues) - 50} more files |")
        lines.append("")

    # Recommendations
    lines.append("## Recommendations")
    lines.append("")
    recs = []
    if junk:
        recs.append(f"1. **Clean junk files** — {len(junk)} junk files detected. Run `cleanup_files.py` to remove them safely.")
    if duplicates and duplicates.get("total_duplicate_files", 0) > 0:
        recs.append(f"2. **Remove duplicates** — {duplicates['total_duplicate_files']} duplicate files wasting {duplicates['total_wasted_human']}. Review the deduplication report and remove confirmed duplicates.")
    if naming_issues:
        recs.append(f"3. **Fix naming issues** — {len(naming_issues)} files have naming problems. Run `rename_files.py --strategy clean` to fix.")
    if empty:
        recs.append(f"4. **Review empty files** — {len(empty)} empty (0-byte) files found. Remove unless intentional.")
    if len(categories) > 5:
        recs.append(f"5. **Organize by category** — Files span {len(categories)} categories. Run `generate_plan.py --strategy category` to sort them.")
    if not recs:
        recs.append("Directory appears well-organized. No immediate action needed.")
    lines.extend(recs)
    lines.append("")

    return "\n".join(lines)


def main():
    parser = argparse.ArgumentParser(description="Generate organization report.")
    parser.add_argument("scan_report", help="Path to scan report JSON")
    parser.add_argument("--duplicates", "-d", default=None, help="Path to duplicate report JSON")
    parser.add_argument("--output", "-o", default=None, help="Output Markdown file")
    args = parser.parse_args()

    if not os.path.isfile(args.scan_report):
        print(f"Error: '{args.scan_report}' not found.", file=sys.stderr)
        sys.exit(1)

    with open(args.scan_report, "r", encoding="utf-8") as f:
        scan = json.load(f)

    duplicates = None
    if args.duplicates and os.path.isfile(args.duplicates):
        with open(args.duplicates, "r", encoding="utf-8") as f:
            duplicates = json.load(f)

    report = generate_report(scan, duplicates)

    if args.output:
        with open(args.output, "w", encoding="utf-8") as f:
            f.write(report)
        print(f"Report saved to {args.output}", file=sys.stderr)
    else:
        print(report)


if __name__ == "__main__":
    main()
