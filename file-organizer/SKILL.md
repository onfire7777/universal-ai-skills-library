---
name: file-organizer
description: Comprehensive file organization suite. Use for cleaning up messy folders, arranging the desktop, finding duplicates, renaming files intelligently, and categorizing files by type or date. Includes safety features like dry-run mode and undo capabilities.
---

# File Organizer Skill

The ultimate suite for cleaning, organizing, and managing files. This skill provides robust tools for deduplication, safe cleanup, intelligent renaming, and desktop arrangement.

## Core Principles

- **Safety First**: All scripts support `--dry-run` by default (except report generation). You must explicitly pass `--execute` to make changes.
- **Reversible Actions**: All destructive actions (moves, renames, safe cleanups) generate an undo log.
- **Accurate Classification**: Uses both extensions and content signatures (via hashing) to classify and compare files accurately.

## Available Workflows

### 1. Full Folder Cleanup Workflow

When a user asks to "clean up a messy folder", follow these steps:

1. **Scan**: `python /home/ubuntu/skills/file-organizer/scripts/scan_files.py <target_dir> -o /tmp/scan.json`
2. **Report**: Generate a summary to understand the mess:
   `python /home/ubuntu/skills/file-organizer/scripts/generate_report.py /tmp/scan.json -o /tmp/report.md`
   *Read this report to inform the user.*
3. **Clean Junk**: `python /home/ubuntu/skills/file-organizer/scripts/cleanup_files.py /tmp/scan.json -o /tmp/cleanup.json`
   *If user approves, run with `--execute`.*
4. **Deduplicate**: `python /home/ubuntu/skills/file-organizer/scripts/find_duplicates.py /tmp/scan.json -o /tmp/dups.json`
   *Review duplicates and delete redundant copies.*
5. **Organize**: `python /home/ubuntu/skills/file-organizer/scripts/generate_plan.py /tmp/scan.json --strategy category --dest <target_dir>/Organized -o /tmp/plan.json`
   *Run `organize_files.py /tmp/plan.json --execute` to move files.*

### 2. Desktop Arrangement Workflow

When a user asks to "clean up my desktop":

1. **Plan**: `python /home/ubuntu/skills/file-organizer/scripts/arrange_desktop.py <desktop_path> --strategy tidy -o /tmp/desktop_plan.json`
   *Strategies: `tidy` (group into folders), `minimal` (move everything except shortcuts), `archive` (move to dated folder).*
2. **Review**: Check `/tmp/desktop_plan.json` to see what stays (shortcuts) and what moves.
3. **Execute**: `python /home/ubuntu/skills/file-organizer/scripts/arrange_desktop.py <desktop_path> --strategy tidy --execute`

### 3. Intelligent Renaming Workflow

When a user asks to "fix filenames" or "standardize names":

1. **Scan**: `python /home/ubuntu/skills/file-organizer/scripts/scan_files.py <target_dir> -o /tmp/scan.json`
2. **Plan**: `python /home/ubuntu/skills/file-organizer/scripts/rename_files.py /tmp/scan.json --strategy clean -o /tmp/rename.json`
   *Strategies: `clean` (remove clutter/copies), `standardize` (enforce snake_case, etc.), `date-prefix`, `sequential`.*
3. **Execute**: `python /home/ubuntu/skills/file-organizer/scripts/rename_files.py /tmp/scan.json --strategy clean --execute`

## Script Reference

### `scan_files.py`
The foundation of all other scripts. Scans a directory and produces a comprehensive JSON inventory.
**Usage:** `python /home/ubuntu/skills/file-organizer/scripts/scan_files.py <dir> -o <scan.json>`

### `find_duplicates.py`
Finds exact duplicates using SHA-256 hashing. Safe and read-only.
**Usage:** `python /home/ubuntu/skills/file-organizer/scripts/find_duplicates.py <scan.json> -o <dups.json>`

### `cleanup_files.py`
Identifies and safely removes OS junk (`.DS_Store`, `thumbs.db`), temp files, and empty files.
**Usage:** `python /home/ubuntu/skills/file-organizer/scripts/cleanup_files.py <scan.json> --execute --trash-dir <path>`

### `generate_plan.py` & `organize_files.py`
Generates an organization plan (by category, extension, or date) and executes it.
**Usage:**
1. `python /home/ubuntu/skills/file-organizer/scripts/generate_plan.py <scan.json> --strategy category --dest <dest> -o <plan.json>`
2. `python /home/ubuntu/skills/file-organizer/scripts/organize_files.py <plan.json> --execute --log <log.json>`

### `rename_files.py`
Cleans up messy names (e.g., "document copy (1).txt" -> "document.txt") or standardizes them.
**Usage:** `python /home/ubuntu/skills/file-organizer/scripts/rename_files.py <scan.json> --strategy clean --execute`

### `arrange_desktop.py`
Specialized script for desktops. Keeps shortcuts and pinned items, moves clutter to organized folders.
**Usage:** `python /home/ubuntu/skills/file-organizer/scripts/arrange_desktop.py <desktop_path> --strategy minimal --execute`

### `undo_operations.py`
Reverses any executed plan (moves, renames) using the generated log file.
**Usage:** `python /home/ubuntu/skills/file-organizer/scripts/undo_operations.py <log.json> --execute`

## Best Practices

1. **Always read the report**: Generate a Markdown report first (`generate_report.py`) to understand the scale of the mess before acting.
2. **Use `--dry-run`**: When in doubt, run the execute scripts without `--execute` first to preview the actions.
3. **Save logs**: Always use `--log <file.json>` when executing moves or renames so you can undo them if the user changes their mind.
