#!/usr/bin/env python3
"""
organize_files.py - Organize files into a clean folder structure.

Reads a scan report and an organization plan (JSON), then moves files
into categorized folders. Supports dry-run mode for safe previewing.

Usage:
    python organize_files.py <plan.json> [--dry-run] [--log <logfile.json>]

Plan JSON format:
{
    "source_dir": "/path/to/messy/folder",
    "dest_dir": "/path/to/organized/folder",
    "strategy": "category" | "extension" | "date" | "custom",
    "actions": [
        {
            "source": "/absolute/path/to/file.pdf",
            "destination": "Documents/file.pdf",
            "action": "move" | "copy"
        }
    ]
}

Safety:
  - Default mode is --dry-run (no changes made)
  - Generates an undo log for every operation
  - Never overwrites existing files (appends counter)
  - Validates all paths before execution
"""

import argparse
import json
import os
import shutil
import sys
from datetime import datetime, timezone


def safe_dest_path(dest: str) -> str:
    """Ensure destination path does not overwrite existing files."""
    if not os.path.exists(dest):
        return dest
    base, ext = os.path.splitext(dest)
    counter = 1
    while os.path.exists(f"{base}_{counter}{ext}"):
        counter += 1
    return f"{base}_{counter}{ext}"


def validate_plan(plan: dict) -> list:
    """Validate the organization plan and return list of issues."""
    issues = []
    if "actions" not in plan:
        issues.append("Missing 'actions' key in plan")
        return issues

    if "dest_dir" not in plan:
        issues.append("Missing 'dest_dir' key in plan")

    for i, action in enumerate(plan["actions"]):
        if "source" not in action:
            issues.append(f"Action {i}: missing 'source'")
        elif not os.path.isfile(action["source"]):
            issues.append(f"Action {i}: source file not found: {action['source']}")
        if "destination" not in action:
            issues.append(f"Action {i}: missing 'destination'")
        if action.get("action", "move") not in ("move", "copy"):
            issues.append(f"Action {i}: invalid action type '{action.get('action')}'")

    return issues


def execute_plan(plan: dict, dry_run: bool = True, log_file: str = None) -> dict:
    """Execute the organization plan."""
    dest_dir = plan.get("dest_dir", "")
    actions = plan.get("actions", [])
    results = {
        "dry_run": dry_run,
        "timestamp": datetime.now(tz=timezone.utc).isoformat(),
        "dest_dir": dest_dir,
        "total_actions": len(actions),
        "completed": 0,
        "skipped": 0,
        "errors": 0,
        "operations": [],
        "undo_operations": [],
    }

    for action in actions:
        src = action["source"]
        rel_dest = action["destination"]
        op_type = action.get("action", "move")
        full_dest = os.path.join(dest_dir, rel_dest)

        # Ensure no overwrite
        full_dest = safe_dest_path(full_dest)

        op = {
            "source": src,
            "destination": full_dest,
            "action": op_type,
            "status": "pending",
        }

        if not os.path.isfile(src):
            op["status"] = "skipped"
            op["reason"] = "source not found"
            results["skipped"] += 1
            results["operations"].append(op)
            continue

        if dry_run:
            op["status"] = "dry_run"
            results["completed"] += 1
            results["operations"].append(op)
            continue

        try:
            # Create destination directory
            os.makedirs(os.path.dirname(full_dest), exist_ok=True)

            if op_type == "move":
                shutil.move(src, full_dest)
                # Record undo operation
                results["undo_operations"].append({
                    "action": "move",
                    "source": full_dest,
                    "destination": src,
                })
            elif op_type == "copy":
                shutil.copy2(src, full_dest)
                results["undo_operations"].append({
                    "action": "delete",
                    "target": full_dest,
                })

            op["status"] = "completed"
            results["completed"] += 1
        except Exception as e:
            op["status"] = "error"
            op["error"] = str(e)
            results["errors"] += 1

        results["operations"].append(op)

    # Save log
    if log_file:
        with open(log_file, "w", encoding="utf-8") as f:
            json.dump(results, f, indent=2, ensure_ascii=False)
        print(f"Operation log saved to {log_file}", file=sys.stderr)

    return results


def main():
    parser = argparse.ArgumentParser(description="Organize files according to a plan.")
    parser.add_argument("plan", help="Path to organization plan JSON")
    parser.add_argument("--dry-run", action="store_true", default=True,
                        help="Preview changes without executing (default: True)")
    parser.add_argument("--execute", action="store_true",
                        help="Actually execute the plan (overrides --dry-run)")
    parser.add_argument("--log", "-l", default=None, help="Path to save operation log JSON")
    args = parser.parse_args()

    if not os.path.isfile(args.plan):
        print(f"Error: Plan file '{args.plan}' not found.", file=sys.stderr)
        sys.exit(1)

    with open(args.plan, "r", encoding="utf-8") as f:
        plan = json.load(f)

    # Validate
    issues = validate_plan(plan)
    if issues:
        print("Plan validation errors:", file=sys.stderr)
        for issue in issues:
            print(f"  - {issue}", file=sys.stderr)
        sys.exit(1)

    dry_run = not args.execute
    if dry_run:
        print("=== DRY RUN MODE (no files will be changed) ===", file=sys.stderr)

    results = execute_plan(plan, dry_run=dry_run, log_file=args.log)

    # Summary
    print(f"\nResults:", file=sys.stderr)
    print(f"  Total actions: {results['total_actions']}", file=sys.stderr)
    print(f"  Completed: {results['completed']}", file=sys.stderr)
    print(f"  Skipped: {results['skipped']}", file=sys.stderr)
    print(f"  Errors: {results['errors']}", file=sys.stderr)

    if not args.log:
        print(json.dumps(results, indent=2, ensure_ascii=False))


if __name__ == "__main__":
    main()
