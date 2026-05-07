#!/usr/bin/env python3
"""
find_duplicates.py - Find exact duplicate files using SHA-256 verification.

Takes a scan report (from scan_files.py) and performs full-hash verification
on potential duplicates. Outputs a deduplication plan.

Usage:
    python find_duplicates.py <scan_report.json> [--output <duplicates.json>]

Safety: Read-only. Does not modify or delete any files.
"""

import argparse
import hashlib
import json
import os
import sys
from collections import defaultdict


def full_hash(filepath: str, chunk_size: int = 65536) -> str:
    """Compute SHA-256 hash of a file."""
    h = hashlib.sha256()
    try:
        with open(filepath, "rb") as f:
            while True:
                chunk = f.read(chunk_size)
                if not chunk:
                    break
                h.update(chunk)
        return h.hexdigest()
    except (PermissionError, OSError):
        return "UNREADABLE"


def find_duplicates(scan_report: dict) -> dict:
    """Identify exact duplicates from scan report using full SHA-256 hashing."""
    files = scan_report.get("files", [])

    # Step 1: Group by size (files of unique size cannot be duplicates)
    size_groups = defaultdict(list)
    for f in files:
        if f["size_bytes"] > 0:  # Skip empty files
            size_groups[f["size_bytes"]].append(f)

    candidates = []
    for size, group in size_groups.items():
        if len(group) > 1:
            candidates.extend(group)

    # Step 2: Group by quick_hash (pre-filter)
    quick_groups = defaultdict(list)
    for f in candidates:
        qh = f.get("quick_hash", "")
        if qh and qh != "UNREADABLE":
            quick_groups[qh].append(f)

    pre_filtered = []
    for qh, group in quick_groups.items():
        if len(group) > 1:
            pre_filtered.extend(group)

    # Step 3: Full SHA-256 verification
    hash_groups = defaultdict(list)
    for f in pre_filtered:
        fh = full_hash(f["path"])
        if fh != "UNREADABLE":
            hash_groups[fh].append(f)

    # Step 4: Build duplicate groups (only groups with 2+ files)
    duplicate_groups = []
    total_wasted = 0

    for fhash, group in hash_groups.items():
        if len(group) < 2:
            continue

        # Sort: keep the oldest file (by modified date) as the original
        group.sort(key=lambda x: x["modified"])

        original = group[0]
        duplicates = group[1:]
        wasted = sum(d["size_bytes"] for d in duplicates)
        total_wasted += wasted

        duplicate_groups.append({
            "hash": fhash,
            "size_bytes": original["size_bytes"],
            "size_human": original["size_human"],
            "count": len(group),
            "wasted_bytes": wasted,
            "wasted_human": human_size(wasted),
            "keep": {
                "path": original["path"],
                "relative_path": original["relative_path"],
                "name": original["name"],
                "modified": original["modified"],
                "reason": "oldest_modified_date",
            },
            "remove_candidates": [
                {
                    "path": d["path"],
                    "relative_path": d["relative_path"],
                    "name": d["name"],
                    "modified": d["modified"],
                }
                for d in duplicates
            ],
        })

    # Sort by wasted space (largest first)
    duplicate_groups.sort(key=lambda x: x["wasted_bytes"], reverse=True)

    result = {
        "scan_target": scan_report.get("scan_target", ""),
        "total_files_scanned": len(files),
        "candidates_by_size": len(candidates),
        "candidates_by_quick_hash": len(pre_filtered),
        "confirmed_duplicate_groups": len(duplicate_groups),
        "total_duplicate_files": sum(len(g["remove_candidates"]) for g in duplicate_groups),
        "total_wasted_bytes": total_wasted,
        "total_wasted_human": human_size(total_wasted),
        "groups": duplicate_groups,
    }
    return result


def human_size(nbytes: int) -> str:
    """Convert bytes to human-readable size."""
    for unit in ("B", "KB", "MB", "GB", "TB"):
        if abs(nbytes) < 1024:
            return f"{nbytes:.1f} {unit}"
        nbytes /= 1024
    return f"{nbytes:.1f} PB"


def main():
    parser = argparse.ArgumentParser(description="Find exact duplicate files.")
    parser.add_argument("scan_report", help="Path to scan report JSON from scan_files.py")
    parser.add_argument("--output", "-o", default=None, help="Output JSON file (default: stdout)")
    args = parser.parse_args()

    if not os.path.isfile(args.scan_report):
        print(f"Error: '{args.scan_report}' not found.", file=sys.stderr)
        sys.exit(1)

    with open(args.scan_report, "r", encoding="utf-8") as f:
        scan_report = json.load(f)

    result = find_duplicates(scan_report)

    output = json.dumps(result, indent=2, ensure_ascii=False)
    if args.output:
        with open(args.output, "w", encoding="utf-8") as f:
            f.write(output)
        print(f"Duplicate report saved to {args.output}", file=sys.stderr)
    else:
        print(output)


if __name__ == "__main__":
    main()
