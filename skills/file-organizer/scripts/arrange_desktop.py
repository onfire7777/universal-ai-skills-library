#!/usr/bin/env python3
"""
arrange_desktop.py - Clean desktop arrangement and organization.

Analyzes the user's Desktop folder and creates a plan to declutter it by:
  1. Identifying what belongs on the desktop vs. what should be filed away
  2. Grouping loose files into organized folders or moving to proper locations
  3. Preserving essential shortcuts and pinned items
  4. Reducing visual clutter to a clean, minimal desktop

Usage:
    python arrange_desktop.py <desktop_path> [--output <plan.json>] [--max-desktop-items <N>]
    python arrange_desktop.py <desktop_path> --execute --dest <organized_dir> [--log <log.json>]

Strategies:
    minimal    - Keep only shortcuts/links on desktop, move everything else
    tidy       - Group loose files into categorized desktop folders, keep shortcuts
    archive    - Move all non-essential items to a dated archive folder

Safety:
  - Default mode is plan-only (no changes)
  - Never removes shortcuts or .lnk files without explicit instruction
  - Preserves pinned/essential items
  - Generates undo log for every operation
"""

import argparse
import json
import os
import re
import shutil
import stat
import sys
from datetime import datetime, timezone
from pathlib import Path


# ── Desktop-Specific Classification ─────────────────────────────────────────

# Items that should ALWAYS stay on the desktop
DESKTOP_KEEP_PATTERNS = [
    r".*\.lnk$",           # Windows shortcuts
    r".*\.url$",           # Internet shortcuts
    r".*\.desktop$",       # Linux desktop entries
    r".*\.webloc$",        # macOS web shortcuts
    r"^desktop\.ini$",     # Windows desktop config
    r"^\.localized$",      # macOS localization
]

# Common desktop folder names that indicate intentional organization
INTENTIONAL_FOLDERS = {
    "shortcuts", "links", "pinned", "important", "workspace",
    "current project", "current_project", "active", "todo",
}

# File categories with desktop-specific handling
DESKTOP_CATEGORIES = {
    "shortcuts": {
        "extensions": {".lnk", ".url", ".desktop", ".webloc"},
        "desktop_action": "keep",
        "reason": "Application/web shortcuts belong on desktop",
    },
    "screenshots": {
        "patterns": [r"^screenshot", r"^screen\s*shot", r"^capture", r"^snip",
                     r"^scr_", r"^grab_"],
        "desktop_action": "move",
        "target_folder": "Screenshots",
        "reason": "Screenshots accumulate quickly and clutter the desktop",
    },
    "downloads_overflow": {
        "patterns": [r"^download", r"^attachment", r"\(\d+\)\."],
        "desktop_action": "move",
        "target_folder": "Downloads",
        "reason": "Downloaded files should be in Downloads folder",
    },
    "temp_documents": {
        "patterns": [r"^untitled", r"^new\s*(document|file|text)", r"^temp",
                     r"^tmp", r"^test", r"^draft"],
        "desktop_action": "move",
        "target_folder": "Drafts",
        "reason": "Temporary/draft documents clutter the desktop",
    },
    "media_files": {
        "extensions": {".jpg", ".jpeg", ".png", ".gif", ".bmp", ".mp4", ".mov",
                       ".avi", ".mp3", ".wav", ".heic", ".webp", ".tiff"},
        "desktop_action": "move",
        "target_folder": "Media",
        "reason": "Media files are better organized in dedicated folders",
    },
    "archives": {
        "extensions": {".zip", ".rar", ".7z", ".tar", ".gz", ".dmg", ".iso"},
        "desktop_action": "move",
        "target_folder": "Archives",
        "reason": "Archive files should be extracted or stored elsewhere",
    },
    "installers": {
        "extensions": {".exe", ".msi", ".pkg", ".deb", ".rpm", ".app", ".apk"},
        "desktop_action": "move",
        "target_folder": "Installers",
        "reason": "Installers should be removed or stored after use",
    },
}

# Category map for remaining files (reuse from scan_files.py logic)
EXT_TO_CATEGORY = {
    ".pdf": "Documents", ".doc": "Documents", ".docx": "Documents",
    ".odt": "Documents", ".rtf": "Documents", ".txt": "Documents",
    ".md": "Documents", ".pages": "Documents",
    ".xls": "Spreadsheets", ".xlsx": "Spreadsheets", ".csv": "Spreadsheets",
    ".ods": "Spreadsheets", ".numbers": "Spreadsheets",
    ".ppt": "Presentations", ".pptx": "Presentations", ".key": "Presentations",
    ".py": "Code", ".js": "Code", ".ts": "Code", ".html": "Code",
    ".css": "Code", ".java": "Code", ".c": "Code", ".cpp": "Code",
    ".json": "Data", ".xml": "Data", ".yaml": "Data", ".yml": "Data",
    ".sql": "Data", ".db": "Data", ".sqlite": "Data",
    ".psd": "Design", ".ai": "Design", ".fig": "Design", ".sketch": "Design",
    ".ttf": "Fonts", ".otf": "Fonts", ".woff": "Fonts",
}


def is_shortcut(name: str) -> bool:
    """Check if a file is a shortcut/link."""
    for pattern in DESKTOP_KEEP_PATTERNS:
        if re.match(pattern, name, re.IGNORECASE):
            return True
    return False


def classify_desktop_item(filepath: str, name: str, size: int, is_dir: bool) -> dict:
    """Classify a desktop item and determine what to do with it."""
    lower_name = name.lower()
    ext = Path(name).suffix.lower()

    # Directories on desktop
    if is_dir:
        if lower_name in INTENTIONAL_FOLDERS:
            return {
                "action": "keep",
                "reason": "intentional_desktop_folder",
                "category": "organized_folder",
            }
        # Check if it's a junk directory
        junk_dirs = {"__macosx", "$recycle.bin", ".spotlight-v100", ".trashes"}
        if lower_name in junk_dirs:
            return {
                "action": "remove",
                "reason": "junk_directory",
                "category": "junk",
            }
        return {
            "action": "evaluate",
            "reason": "user_folder_on_desktop",
            "category": "folder",
        }

    # Always keep shortcuts
    if is_shortcut(name):
        return {
            "action": "keep",
            "reason": "shortcut_or_link",
            "category": "shortcuts",
        }

    # Check desktop-specific categories
    for cat_name, cat_info in DESKTOP_CATEGORIES.items():
        # Check by extension
        if "extensions" in cat_info and ext in cat_info["extensions"]:
            return {
                "action": cat_info["desktop_action"],
                "reason": cat_info["reason"],
                "category": cat_name,
                "target_folder": cat_info.get("target_folder"),
            }
        # Check by pattern
        if "patterns" in cat_info:
            stem = Path(name).stem.lower()
            for pattern in cat_info["patterns"]:
                if re.match(pattern, stem, re.IGNORECASE):
                    return {
                        "action": cat_info["desktop_action"],
                        "reason": cat_info["reason"],
                        "category": cat_name,
                        "target_folder": cat_info.get("target_folder"),
                    }

    # Remaining files: categorize by extension
    general_cat = EXT_TO_CATEGORY.get(ext, "Other")
    return {
        "action": "move",
        "reason": f"general_{general_cat.lower()}_file",
        "category": general_cat,
        "target_folder": general_cat,
    }


def scan_desktop(desktop_path: str) -> list:
    """Scan the desktop directory (non-recursive, top-level only)."""
    items = []
    desktop_path = os.path.abspath(desktop_path)

    try:
        entries = os.listdir(desktop_path)
    except PermissionError:
        print(f"Error: Cannot access '{desktop_path}'.", file=sys.stderr)
        return items

    for entry_name in entries:
        full_path = os.path.join(desktop_path, entry_name)
        try:
            st = os.stat(full_path)
        except (PermissionError, OSError):
            continue

        is_dir = stat.S_ISDIR(st.st_mode)
        size = st.st_size if not is_dir else 0

        # For directories, calculate total size
        if is_dir:
            for root, dirs, files in os.walk(full_path):
                for f in files:
                    try:
                        size += os.path.getsize(os.path.join(root, f))
                    except (PermissionError, OSError):
                        pass

        classification = classify_desktop_item(full_path, entry_name, size, is_dir)

        items.append({
            "name": entry_name,
            "path": full_path,
            "is_directory": is_dir,
            "size_bytes": size,
            "size_human": human_size(size),
            "extension": Path(entry_name).suffix.lower() if not is_dir else "",
            "modified": datetime.fromtimestamp(st.st_mtime, tz=timezone.utc).isoformat(),
            "is_hidden": entry_name.startswith("."),
            **classification,
        })

    return items


def build_arrangement_plan(desktop_path: str, items: list, strategy: str = "tidy",
                           dest_dir: str = None, max_desktop_items: int = 15) -> dict:
    """Build a desktop arrangement plan."""
    desktop_path = os.path.abspath(desktop_path)
    if dest_dir is None:
        dest_dir = os.path.join(os.path.dirname(desktop_path), "Desktop_Organized")

    keep_items = []
    move_items = []
    remove_items = []

    for item in items:
        if item.get("is_hidden", False) and item["action"] != "remove":
            keep_items.append(item)
            continue

        if item["action"] == "keep":
            keep_items.append(item)
        elif item["action"] == "remove":
            remove_items.append(item)
        elif item["action"] == "move":
            move_items.append(item)
        elif item["action"] == "evaluate":
            # For "evaluate" items (user folders), strategy determines behavior
            if strategy == "minimal":
                move_items.append(item)
            else:
                keep_items.append(item)

    # Build move actions
    actions = []
    for item in move_items:
        target_folder = item.get("target_folder", item.get("category", "Other"))

        if strategy == "archive":
            # Everything goes into a dated archive
            date_str = datetime.now().strftime("%Y-%m-%d")
            rel_dest = os.path.join(f"Desktop_Archive_{date_str}", target_folder, item["name"])
        elif strategy == "minimal":
            # Move to organized directory by category
            rel_dest = os.path.join(target_folder, item["name"])
        else:  # tidy
            # Create categorized folders ON the desktop for small groups,
            # move to dest_dir for large groups
            rel_dest = os.path.join(target_folder, item["name"])

        actions.append({
            "source": item["path"],
            "destination": rel_dest,
            "action": "move",
            "category": item.get("category", "Other"),
            "reason": item.get("reason", ""),
            "is_directory": item.get("is_directory", False),
        })

    # Summary
    plan = {
        "desktop_path": desktop_path,
        "dest_dir": os.path.abspath(dest_dir),
        "strategy": strategy,
        "timestamp": datetime.now(tz=timezone.utc).isoformat(),
        "current_item_count": len(items),
        "max_recommended": max_desktop_items,
        "items_to_keep": len(keep_items),
        "items_to_move": len(move_items),
        "items_to_remove": len(remove_items),
        "keep": [{"name": i["name"], "reason": i.get("reason", "")} for i in keep_items],
        "remove": [{"name": i["name"], "path": i["path"], "reason": i.get("reason", "")} for i in remove_items],
        "actions": actions,
        "result_desktop_items": len(keep_items),
        "is_clean": len(keep_items) <= max_desktop_items,
    }
    return plan


def execute_arrangement(plan: dict, dry_run: bool = True, log_file: str = None) -> dict:
    """Execute the desktop arrangement plan."""
    dest_dir = plan.get("dest_dir", "")
    actions = plan.get("actions", [])

    results = {
        "dry_run": dry_run,
        "timestamp": datetime.now(tz=timezone.utc).isoformat(),
        "total": len(actions),
        "completed": 0,
        "skipped": 0,
        "errors": 0,
        "operations": [],
        "undo_operations": [],
    }

    for action in actions:
        src = action["source"]
        rel_dest = action["destination"]
        full_dest = os.path.join(dest_dir, rel_dest)

        op = {
            "source": src,
            "destination": full_dest,
            "status": "pending",
        }

        if not os.path.exists(src):
            op["status"] = "skipped"
            op["reason"] = "source not found"
            results["skipped"] += 1
            results["operations"].append(op)
            continue

        # Prevent overwrite
        if os.path.exists(full_dest):
            stem = Path(full_dest).stem
            ext = Path(full_dest).suffix
            parent = os.path.dirname(full_dest)
            counter = 1
            while os.path.exists(os.path.join(parent, f"{stem}_{counter}{ext}")):
                counter += 1
            full_dest = os.path.join(parent, f"{stem}_{counter}{ext}")
            op["destination"] = full_dest

        if dry_run:
            op["status"] = "dry_run"
            results["completed"] += 1
            results["operations"].append(op)
            continue

        try:
            os.makedirs(os.path.dirname(full_dest), exist_ok=True)
            shutil.move(src, full_dest)
            op["status"] = "completed"
            results["completed"] += 1
            results["undo_operations"].append({
                "action": "move",
                "source": full_dest,
                "destination": src,
            })
        except Exception as e:
            op["status"] = "error"
            op["error"] = str(e)
            results["errors"] += 1

        results["operations"].append(op)

    if log_file:
        with open(log_file, "w", encoding="utf-8") as f:
            json.dump(results, f, indent=2, ensure_ascii=False)
        print(f"Arrangement log saved to {log_file}", file=sys.stderr)

    return results


def human_size(nbytes: int) -> str:
    for unit in ("B", "KB", "MB", "GB", "TB"):
        if abs(nbytes) < 1024:
            return f"{nbytes:.1f} {unit}"
        nbytes /= 1024
    return f"{nbytes:.1f} PB"


def main():
    parser = argparse.ArgumentParser(description="Clean desktop arrangement.")
    parser.add_argument("desktop_path", help="Path to Desktop folder")
    parser.add_argument("--strategy", "-s", default="tidy",
                        choices=["minimal", "tidy", "archive"],
                        help="Arrangement strategy (default: tidy)")
    parser.add_argument("--dest", "-d", default=None,
                        help="Destination for moved items (default: ../Desktop_Organized)")
    parser.add_argument("--max-desktop-items", type=int, default=15,
                        help="Max recommended items on desktop (default: 15)")
    parser.add_argument("--output", "-o", default=None, help="Output plan JSON")
    parser.add_argument("--execute", action="store_true", help="Execute the plan")
    parser.add_argument("--log", "-l", default=None, help="Operation log path")
    args = parser.parse_args()

    if not os.path.isdir(args.desktop_path):
        print(f"Error: '{args.desktop_path}' is not a valid directory.", file=sys.stderr)
        sys.exit(1)

    # Scan desktop
    items = scan_desktop(args.desktop_path)
    plan = build_arrangement_plan(
        args.desktop_path, items,
        strategy=args.strategy,
        dest_dir=args.dest,
        max_desktop_items=args.max_desktop_items,
    )

    if not args.execute:
        output = json.dumps(plan, indent=2, ensure_ascii=False)
        if args.output:
            with open(args.output, "w", encoding="utf-8") as f:
                f.write(output)
            print(f"Desktop arrangement plan saved to {args.output}", file=sys.stderr)
        else:
            print(output)

        print(f"\nDesktop Summary:", file=sys.stderr)
        print(f"  Current items: {plan['current_item_count']}", file=sys.stderr)
        print(f"  Items to keep: {plan['items_to_keep']}", file=sys.stderr)
        print(f"  Items to move: {plan['items_to_move']}", file=sys.stderr)
        print(f"  Items to remove: {plan['items_to_remove']}", file=sys.stderr)
        print(f"  Result: {'Clean' if plan['is_clean'] else 'Still cluttered'} "
              f"({plan['result_desktop_items']}/{plan['max_recommended']} items)", file=sys.stderr)
        return

    # Execute
    if not args.dest:
        args.dest = os.path.join(os.path.dirname(os.path.abspath(args.desktop_path)), "Desktop_Organized")

    plan["dest_dir"] = os.path.abspath(args.dest)
    results = execute_arrangement(plan, dry_run=False, log_file=args.log)

    print(f"\nArrangement Results:", file=sys.stderr)
    print(f"  Total: {results['total']}", file=sys.stderr)
    print(f"  Completed: {results['completed']}", file=sys.stderr)
    print(f"  Skipped: {results['skipped']}", file=sys.stderr)
    print(f"  Errors: {results['errors']}", file=sys.stderr)

    if not args.log:
        print(json.dumps(results, indent=2, ensure_ascii=False))


if __name__ == "__main__":
    main()
