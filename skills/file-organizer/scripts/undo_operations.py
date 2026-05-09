#!/usr/bin/env python3
"""
undo_operations.py - Reverse file organization operations using undo logs.

Reads an operation log (from organize_files.py, rename_files.py, or
cleanup_files.py) and reverses all completed operations.

Usage:
    python undo_operations.py <operation_log.json> [--dry-run]

Safety: Reverses operations in reverse order to maintain consistency.
"""

import argparse
import json
import os
import shutil
import sys
from datetime import datetime, timezone


def undo_operations(log_data: dict, dry_run: bool = True) -> dict:
    """Reverse operations from an operation log."""
    undo_ops = log_data.get("undo_operations", [])

    if not undo_ops:
        print("No undo operations found in log.", file=sys.stderr)
        return {"status": "nothing_to_undo", "total": 0, "completed": 0, "skipped": 0, "errors": 0, "operations": []}

    # Reverse order for consistency
    undo_ops = list(reversed(undo_ops))

    results = {
        "dry_run": dry_run,
        "timestamp": datetime.now(tz=timezone.utc).isoformat(),
        "original_log": log_data.get("timestamp", "unknown"),
        "total": len(undo_ops),
        "completed": 0,
        "skipped": 0,
        "errors": 0,
        "operations": [],
    }

    for op in undo_ops:
        action = op.get("action", "unknown")
        entry = {"action": action, "status": "pending", **op}

        if action == "move" or action == "rename":
            src = op.get("source", "")
            dest = op.get("destination", "")

            if not src or not dest:
                entry["status"] = "skipped"
                entry["reason"] = "missing source or destination"
                results["skipped"] += 1
                results["operations"].append(entry)
                continue

            if not os.path.exists(src):
                entry["status"] = "skipped"
                entry["reason"] = f"source not found: {src}"
                results["skipped"] += 1
                results["operations"].append(entry)
                continue

            if dry_run:
                entry["status"] = "dry_run"
                results["completed"] += 1
                results["operations"].append(entry)
                continue

            try:
                os.makedirs(os.path.dirname(dest), exist_ok=True)
                shutil.move(src, dest)
                entry["status"] = "completed"
                results["completed"] += 1
            except Exception as e:
                entry["status"] = "error"
                entry["error"] = str(e)
                results["errors"] += 1

        elif action == "delete":
            target = op.get("target", "")
            if os.path.exists(target):
                if dry_run:
                    entry["status"] = "dry_run"
                    results["completed"] += 1
                else:
                    try:
                        os.remove(target)
                        entry["status"] = "completed"
                        results["completed"] += 1
                    except Exception as e:
                        entry["status"] = "error"
                        entry["error"] = str(e)
                        results["errors"] += 1
            else:
                entry["status"] = "skipped"
                entry["reason"] = "target not found"
                results["skipped"] += 1

        elif action == "undelete_impossible":
            entry["status"] = "skipped"
            entry["reason"] = "permanent deletion cannot be undone"
            results["skipped"] += 1

        else:
            entry["status"] = "skipped"
            entry["reason"] = f"unknown action type: {action}"
            results["skipped"] += 1

        results["operations"].append(entry)

    return results


def main():
    parser = argparse.ArgumentParser(description="Undo file organization operations.")
    parser.add_argument("log_file", help="Path to operation log JSON")
    parser.add_argument("--dry-run", action="store_true", default=True,
                        help="Preview undo without executing (default: True)")
    parser.add_argument("--execute", action="store_true",
                        help="Actually execute undo (overrides --dry-run)")
    args = parser.parse_args()

    if not os.path.isfile(args.log_file):
        print(f"Error: '{args.log_file}' not found.", file=sys.stderr)
        sys.exit(1)

    with open(args.log_file, "r", encoding="utf-8") as f:
        log_data = json.load(f)

    dry_run = not args.execute
    if dry_run:
        print("=== DRY RUN MODE (no changes will be made) ===", file=sys.stderr)

    results = undo_operations(log_data, dry_run=dry_run)

    print(f"\nUndo Results:", file=sys.stderr)
    print(f"  Total: {results['total']}", file=sys.stderr)
    print(f"  Completed: {results['completed']}", file=sys.stderr)
    print(f"  Skipped: {results['skipped']}", file=sys.stderr)
    print(f"  Errors: {results['errors']}", file=sys.stderr)

    print(json.dumps(results, indent=2, ensure_ascii=False))


if __name__ == "__main__":
    main()
