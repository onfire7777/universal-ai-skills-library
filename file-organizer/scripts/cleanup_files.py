#!/usr/bin/env python3
"""
cleanup_files.py - Safely remove junk, temporary, and empty files.

Reads a scan report and identifies files safe to remove: OS junk files,
temporary files, empty files, and orphaned system files.

Usage:
    python cleanup_files.py <scan_report.json> [--output <cleanup_plan.json>]
    python cleanup_files.py <scan_report.json> --execute [--log <log.json>]

Safety:
  - Default mode is plan-only (no deletions)
  - Moves to trash directory instead of permanent deletion (configurable)
  - Generates undo log for every operation
  - Requires --execute flag to perform actual cleanup
"""

import argparse
import json
import os
import shutil
import sys
from datetime import datetime, timezone
from pathlib import Path


# Definitive junk files (always safe to remove)
DEFINITE_JUNK = {
    "thumbs.db", "desktop.ini", ".ds_store", "._.ds_store",
    ".spotlight-v100", ".trashes", ".fseventsd",
    ".temporaryitems", ".apdisk", ".com.apple.timemachine.donotpresent",
    "icon\r", ".localized", ".documentrevisions-v100",
    ".pkinstallsandboxmanager", ".pkinstallsandboxmanager-systemsoftware",
}

# Junk filename patterns (regex)
JUNK_PATTERNS = [
    r"^\._",                     # macOS resource forks
    r"^~\$",                     # Office temp files
    r"^~.*\.tmp$",               # Generic temp files
    r"^\.~lock\.",               # LibreOffice lock files
    r".*\.tmp$",                 # .tmp files
    r".*\.temp$",                # .temp files
    r".*\.swp$",                 # Vim swap files
    r".*\.swo$",                 # Vim swap files
    r".*\.bak$",                 # Backup files
    r".*\.orig$",                # Original backup files
    r".*\.cache$",               # Cache files
    r"^__MACOSX$",               # macOS archive artifacts
]

# Junk directory names (these entire directories are junk)
JUNK_DIRECTORIES = {
    "__macosx", ".spotlight-v100", ".trashes", ".fseventsd",
    ".temporaryitems", "$recycle.bin", "system volume information",
    ".thumbnails", "@eadir",  # Synology NAS thumbnails
}

import re

def classify_file(filepath: str, name: str, size: int, is_junk: bool) -> dict:
    """Classify a file for cleanup with reason and confidence."""
    lower_name = name.lower()
    result = {"removable": False, "reason": None, "confidence": "low", "category": None}

    # Check definite junk
    if lower_name in DEFINITE_JUNK:
        return {"removable": True, "reason": "os_junk_file", "confidence": "high", "category": "junk"}

    # Check junk patterns
    for pattern in JUNK_PATTERNS:
        if re.match(pattern, lower_name):
            return {"removable": True, "reason": f"matches_junk_pattern: {pattern}",
                    "confidence": "high", "category": "junk"}

    # Check if inside a junk directory
    path_parts = Path(filepath).parts
    for part in path_parts:
        if part.lower() in JUNK_DIRECTORIES:
            return {"removable": True, "reason": f"inside_junk_directory: {part}",
                    "confidence": "high", "category": "junk"}

    # macOS resource fork files
    if lower_name.startswith("._"):
        return {"removable": True, "reason": "macos_resource_fork", "confidence": "high", "category": "junk"}

    # Office temp/lock files
    if lower_name.startswith("~$"):
        return {"removable": True, "reason": "office_temp_file", "confidence": "high", "category": "temporary"}

    # Empty files (0 bytes)
    if size == 0:
        # Some empty files are intentional (e.g., __init__.py, .gitkeep)
        intentional_empty = {".gitkeep", ".gitignore", ".npmignore", ".keep",
                             "__init__.py", ".env", ".env.example", ".nojekyll",
                             ".placeholder"}
        if lower_name in intentional_empty:
            return {"removable": False, "reason": "intentional_empty_file",
                    "confidence": "high", "category": "empty"}
        return {"removable": True, "reason": "empty_file_0_bytes",
                "confidence": "medium", "category": "empty"}

    # Flagged as junk by scanner
    if is_junk:
        return {"removable": True, "reason": "flagged_by_scanner",
                "confidence": "medium", "category": "junk"}

    return result


def build_cleanup_plan(scan_report: dict) -> dict:
    """Build a cleanup plan from scan report."""
    files = scan_report.get("files", [])
    removable = []
    kept = []

    for f in files:
        classification = classify_file(
            f["path"], f["name"], f["size_bytes"], f.get("is_junk", False)
        )
        entry = {
            "path": f["path"],
            "relative_path": f["relative_path"],
            "name": f["name"],
            "size_bytes": f["size_bytes"],
            "size_human": f["size_human"],
            **classification,
        }
        if classification["removable"]:
            removable.append(entry)
        else:
            kept.append(entry)

    # Sort by confidence (high first), then by size (largest first)
    confidence_order = {"high": 0, "medium": 1, "low": 2}
    removable.sort(key=lambda x: (confidence_order.get(x["confidence"], 3), -x["size_bytes"]))

    total_recoverable = sum(f["size_bytes"] for f in removable)

    # Group by category
    by_category = {}
    for f in removable:
        cat = f["category"]
        if cat not in by_category:
            by_category[cat] = {"count": 0, "total_size": 0}
        by_category[cat]["count"] += 1
        by_category[cat]["total_size"] += f["size_bytes"]

    plan = {
        "scan_target": scan_report.get("scan_target", ""),
        "timestamp": datetime.now(tz=timezone.utc).isoformat(),
        "total_files_scanned": len(files),
        "removable_count": len(removable),
        "recoverable_bytes": total_recoverable,
        "recoverable_human": human_size(total_recoverable),
        "by_category": by_category,
        "removable_files": removable,
    }
    return plan


def execute_cleanup(plan: dict, trash_dir: str = None, dry_run: bool = True,
                    log_file: str = None, confidence_threshold: str = "medium") -> dict:
    """Execute cleanup by moving files to trash or deleting them."""
    confidence_levels = {"high": 0, "medium": 1, "low": 2}
    threshold = confidence_levels.get(confidence_threshold, 1)

    results = {
        "dry_run": dry_run,
        "timestamp": datetime.now(tz=timezone.utc).isoformat(),
        "trash_dir": trash_dir,
        "total": 0,
        "completed": 0,
        "skipped": 0,
        "errors": 0,
        "bytes_recovered": 0,
        "operations": [],
        "undo_operations": [],
    }

    removable = plan.get("removable_files", [])
    results["total"] = len(removable)

    for f in removable:
        src = f["path"]
        conf = confidence_levels.get(f["confidence"], 2)

        op = {
            "source": src,
            "name": f["name"],
            "confidence": f["confidence"],
            "reason": f["reason"],
            "status": "pending",
        }

        # Skip if below confidence threshold
        if conf > threshold:
            op["status"] = "skipped"
            op["skip_reason"] = f"confidence '{f['confidence']}' below threshold '{confidence_threshold}'"
            results["skipped"] += 1
            results["operations"].append(op)
            continue

        if not os.path.exists(src):
            op["status"] = "skipped"
            op["skip_reason"] = "file not found"
            results["skipped"] += 1
            results["operations"].append(op)
            continue

        if dry_run:
            op["status"] = "dry_run"
            results["completed"] += 1
            results["bytes_recovered"] += f["size_bytes"]
            results["operations"].append(op)
            continue

        try:
            if trash_dir:
                # Move to trash directory (safer)
                trash_path = os.path.join(trash_dir, f["relative_path"])
                os.makedirs(os.path.dirname(trash_path), exist_ok=True)
                shutil.move(src, trash_path)
                op["moved_to"] = trash_path
                results["undo_operations"].append({
                    "action": "move",
                    "source": trash_path,
                    "destination": src,
                })
            else:
                # Permanent delete
                os.remove(src)
                results["undo_operations"].append({
                    "action": "undelete_impossible",
                    "original_path": src,
                })

            op["status"] = "completed"
            results["completed"] += 1
            results["bytes_recovered"] += f["size_bytes"]
        except Exception as e:
            op["status"] = "error"
            op["error"] = str(e)
            results["errors"] += 1

        results["operations"].append(op)

    results["bytes_recovered_human"] = human_size(results["bytes_recovered"])

    if log_file:
        with open(log_file, "w", encoding="utf-8") as f:
            json.dump(results, f, indent=2, ensure_ascii=False)
        print(f"Cleanup log saved to {log_file}", file=sys.stderr)

    return results


def human_size(nbytes: int) -> str:
    for unit in ("B", "KB", "MB", "GB", "TB"):
        if abs(nbytes) < 1024:
            return f"{nbytes:.1f} {unit}"
        nbytes /= 1024
    return f"{nbytes:.1f} PB"


def main():
    parser = argparse.ArgumentParser(description="Safely clean up junk and temporary files.")
    parser.add_argument("scan_report", help="Path to scan report JSON")
    parser.add_argument("--output", "-o", default=None, help="Output cleanup plan JSON")
    parser.add_argument("--execute", action="store_true", help="Execute cleanup")
    parser.add_argument("--trash-dir", default=None,
                        help="Move files here instead of deleting (recommended)")
    parser.add_argument("--confidence", default="medium", choices=["high", "medium", "low"],
                        help="Minimum confidence to act on (default: medium)")
    parser.add_argument("--log", "-l", default=None, help="Operation log path")
    args = parser.parse_args()

    if not os.path.isfile(args.scan_report):
        print(f"Error: '{args.scan_report}' not found.", file=sys.stderr)
        sys.exit(1)

    with open(args.scan_report, "r", encoding="utf-8") as f:
        scan_report = json.load(f)

    plan = build_cleanup_plan(scan_report)

    if not args.execute:
        output = json.dumps(plan, indent=2, ensure_ascii=False)
        if args.output:
            with open(args.output, "w", encoding="utf-8") as f:
                f.write(output)
            print(f"Cleanup plan saved to {args.output}", file=sys.stderr)
            print(f"  Removable files: {plan['removable_count']}", file=sys.stderr)
            print(f"  Recoverable space: {plan['recoverable_human']}", file=sys.stderr)
        else:
            print(output)
        return

    # Execute
    if not args.trash_dir:
        print("WARNING: No --trash-dir specified. Files will be PERMANENTLY deleted.", file=sys.stderr)
        print("Recommendation: Use --trash-dir /path/to/trash for safe cleanup.", file=sys.stderr)

    results = execute_cleanup(
        plan, trash_dir=args.trash_dir, dry_run=False,
        log_file=args.log, confidence_threshold=args.confidence
    )

    print(f"\nCleanup Results:", file=sys.stderr)
    print(f"  Total: {results['total']}", file=sys.stderr)
    print(f"  Completed: {results['completed']}", file=sys.stderr)
    print(f"  Skipped: {results['skipped']}", file=sys.stderr)
    print(f"  Errors: {results['errors']}", file=sys.stderr)
    print(f"  Space recovered: {results['bytes_recovered_human']}", file=sys.stderr)


if __name__ == "__main__":
    main()
