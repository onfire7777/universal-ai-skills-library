#!/usr/bin/env python3
"""
rename_files.py - Intelligent file renaming with clean, accurate, specific names.

Cleans up messy filenames by removing clutter, standardizing formats, and
applying context-aware naming conventions. Supports dry-run mode.

Usage:
    python rename_files.py <scan_report.json> --strategy <strategy> [--dry-run] [--output <plan.json>]

Strategies:
    clean       - Remove clutter (extra spaces, special chars, redundant tokens)
    standardize - Enforce consistent naming convention (snake_case, kebab-case, etc.)
    date-prefix - Prepend YYYY-MM-DD based on file modified date
    descriptive - Use LLM-assisted renaming for vague/generic names (requires plan review)
    sequential  - Number files within each category sequentially

Safety: Generates a rename plan JSON. Does not rename files directly.
        Use --execute flag to apply. Undo log is always generated.
"""

import argparse
import json
import os
import re
import sys
import unicodedata
from datetime import datetime, timezone
from pathlib import Path


# ── Name Cleaning Rules ──────────────────────────────────────────────────────

def normalize_unicode(name: str) -> str:
    """Normalize unicode characters to ASCII equivalents where possible."""
    # NFKD decomposition + ASCII encoding to strip accents
    normalized = unicodedata.normalize("NFKD", name)
    ascii_name = normalized.encode("ascii", "ignore").decode("ascii")
    # If we lost the entire name, keep original
    return ascii_name if ascii_name.strip() else name


def remove_clutter(name: str) -> str:
    """Remove common clutter patterns from filenames."""
    stem = Path(name).stem
    ext = Path(name).suffix.lower()

    # Remove common junk patterns
    patterns_to_remove = [
        r"\s*\(\d+\)\s*",           # (1), (2), etc. - copy indicators
        r"\s*-\s*Copy\s*",          # - Copy
        r"\s*copy\s*\d*\s*",        # copy, copy 2 (case insensitive handled below)
        r"\s*\(copy\)\s*",          # (copy)
        r"\s*-\s*копия\s*",         # Russian copy indicator
        r"\s*-\s*Kopie\s*",         # German copy indicator
        r"\s*-\s*copia\s*",         # Spanish/Italian copy indicator
        r"\s*\(conflicted copy.*?\)", # Dropbox/sync conflicts
        r"\s*\(.*?'s conflicted.*?\)", # Sync conflict markers
        r"^\~\$",                    # Office temp file prefix
        r"\s+final\s*\d*\s*$",      # "final", "final2"
        r"\s+FINAL\s*\d*\s*$",
        r"\s+v\d+\s*$",             # "v1", "v2" at end
        r"\s+draft\s*\d*\s*$",      # "draft", "draft3"
        r"\s+revised\s*$",          # "revised"
        r"\s+edited\s*$",           # "edited"
        r"\s+new\s*$",              # trailing "new"
        r"\s+old\s*$",              # trailing "old"
    ]

    for pattern in patterns_to_remove:
        stem = re.sub(pattern, "", stem, flags=re.IGNORECASE)

    return stem.strip() + ext


def clean_separators(name: str) -> str:
    """Standardize separators and whitespace in filenames."""
    stem = Path(name).stem
    ext = Path(name).suffix.lower()

    # Replace multiple consecutive separators with single space
    stem = re.sub(r"[_\-\s]+", " ", stem)
    # Remove leading/trailing separators
    stem = stem.strip(" _-")
    # Collapse multiple spaces
    stem = re.sub(r"\s{2,}", " ", stem)

    return stem.strip() + ext


def remove_redundant_extensions(name: str) -> str:
    """Remove double extensions like .pdf.pdf or .jpg.jpeg."""
    parts = name.split(".")
    if len(parts) > 2:
        # Check for redundant extensions
        seen = set()
        clean_parts = [parts[0]]
        for part in parts[1:]:
            lower = part.lower()
            if lower not in seen:
                seen.add(lower)
                clean_parts.append(part)
        return ".".join(clean_parts)
    return name


def remove_hash_strings(name: str) -> str:
    """Remove embedded hash strings (common in downloads)."""
    stem = Path(name).stem
    ext = Path(name).suffix.lower()

    # Remove hex strings that look like hashes (8+ hex chars surrounded by separators)
    stem = re.sub(r"[\s_\-][0-9a-f]{8,}[\s_\-]?", " ", stem, flags=re.IGNORECASE)
    # Remove UUIDs
    stem = re.sub(r"[\s_\-]?[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}[\s_\-]?",
                  " ", stem, flags=re.IGNORECASE)
    # Remove trailing random strings (common in web downloads)
    stem = re.sub(r"[\s_\-][a-zA-Z0-9]{20,}$", "", stem)

    return stem.strip() + ext


def remove_url_encoding(name: str) -> str:
    """Decode URL-encoded characters in filenames."""
    stem = Path(name).stem
    ext = Path(name).suffix.lower()

    # Replace %20 and similar URL encodings
    stem = re.sub(r"%20", " ", stem)
    stem = re.sub(r"%28", "(", stem)
    stem = re.sub(r"%29", ")", stem)
    stem = re.sub(r"%2[Bb]", "+", stem)
    stem = re.sub(r"%2[Cc]", ",", stem)
    stem = re.sub(r"%26", "&", stem)
    stem = re.sub(r"%27", "'", stem)
    stem = re.sub(r"%23", "#", stem)
    stem = re.sub(r"%5[Bb]", "[", stem)
    stem = re.sub(r"%5[Dd]", "]", stem)

    return stem.strip() + ext


# ── Naming Convention Converters ─────────────────────────────────────────────

def to_snake_case(name: str) -> str:
    """Convert filename to snake_case."""
    stem = Path(name).stem
    ext = Path(name).suffix.lower()
    # Split on transitions: camelCase, spaces, hyphens, underscores
    tokens = re.sub(r"([a-z])([A-Z])", r"\1 \2", stem)
    tokens = re.sub(r"[_\-\s]+", " ", tokens)
    words = tokens.strip().split()
    return "_".join(w.lower() for w in words if w) + ext


def to_kebab_case(name: str) -> str:
    """Convert filename to kebab-case."""
    stem = Path(name).stem
    ext = Path(name).suffix.lower()
    tokens = re.sub(r"([a-z])([A-Z])", r"\1 \2", stem)
    tokens = re.sub(r"[_\-\s]+", " ", tokens)
    words = tokens.strip().split()
    return "-".join(w.lower() for w in words if w) + ext


def to_title_case(name: str) -> str:
    """Convert filename to Title Case with spaces."""
    stem = Path(name).stem
    ext = Path(name).suffix.lower()
    tokens = re.sub(r"([a-z])([A-Z])", r"\1 \2", stem)
    tokens = re.sub(r"[_\-\s]+", " ", tokens)
    words = tokens.strip().split()
    # Title case but keep short words lowercase (except first)
    small_words = {"a", "an", "the", "and", "or", "but", "in", "on", "at", "to", "for", "of", "with", "by"}
    titled = []
    for i, w in enumerate(words):
        if i == 0 or w.lower() not in small_words:
            titled.append(w.capitalize())
        else:
            titled.append(w.lower())
    return " ".join(titled) + ext


def to_camel_case(name: str) -> str:
    """Convert filename to camelCase."""
    stem = Path(name).stem
    ext = Path(name).suffix.lower()
    tokens = re.sub(r"([a-z])([A-Z])", r"\1 \2", stem)
    tokens = re.sub(r"[_\-\s]+", " ", tokens)
    words = tokens.strip().split()
    if not words:
        return name
    result = words[0].lower() + "".join(w.capitalize() for w in words[1:])
    return result + ext


CONVENTION_MAP = {
    "snake_case": to_snake_case,
    "kebab-case": to_kebab_case,
    "title_case": to_title_case,
    "camelCase": to_camel_case,
}


# ── Strategies ───────────────────────────────────────────────────────────────

def strategy_clean(files: list, **kwargs) -> list:
    """Clean strategy: remove clutter, fix separators, remove hashes."""
    plan = []
    for f in files:
        original = f["name"]
        cleaned = original

        # Apply cleaning pipeline in order
        cleaned = remove_url_encoding(cleaned)
        cleaned = remove_clutter(cleaned)
        cleaned = remove_hash_strings(cleaned)
        cleaned = remove_redundant_extensions(cleaned)
        cleaned = clean_separators(cleaned)
        cleaned = normalize_unicode(cleaned)

        # Only add to plan if name actually changed
        if cleaned != original and cleaned.strip("."):
            plan.append({
                "source": f["path"],
                "original_name": original,
                "new_name": cleaned,
                "reason": "cleaned_clutter",
            })
    return plan


def strategy_standardize(files: list, convention: str = "snake_case", **kwargs) -> list:
    """Standardize strategy: enforce a naming convention."""
    converter = CONVENTION_MAP.get(convention, to_snake_case)
    plan = []
    for f in files:
        original = f["name"]
        # First clean, then convert
        cleaned = remove_url_encoding(original)
        cleaned = remove_clutter(cleaned)
        cleaned = remove_hash_strings(cleaned)
        cleaned = clean_separators(cleaned)
        standardized = converter(cleaned)

        if standardized != original and standardized.strip("."):
            plan.append({
                "source": f["path"],
                "original_name": original,
                "new_name": standardized,
                "convention": convention,
                "reason": "standardized_convention",
            })
    return plan


def strategy_date_prefix(files: list, **kwargs) -> list:
    """Date prefix strategy: prepend YYYY-MM-DD from modified date."""
    plan = []
    for f in files:
        original = f["name"]
        # Parse the modified date
        try:
            dt = datetime.fromisoformat(f["modified"])
            date_str = dt.strftime("%Y-%m-%d")
        except (ValueError, KeyError):
            continue

        # Skip if already has a date prefix
        if re.match(r"^\d{4}-\d{2}-\d{2}", original):
            continue

        # Clean the name first
        cleaned = remove_clutter(original)
        cleaned = clean_separators(cleaned)
        new_name = f"{date_str} {cleaned}"

        if new_name != original:
            plan.append({
                "source": f["path"],
                "original_name": original,
                "new_name": new_name,
                "date_added": date_str,
                "reason": "date_prefix_added",
            })
    return plan


def strategy_sequential(files: list, **kwargs) -> list:
    """Sequential strategy: number files within each category."""
    # Group by category and parent directory
    groups = {}
    for f in files:
        key = (os.path.dirname(f["path"]), f["category"])
        if key not in groups:
            groups[key] = []
        groups[key].append(f)

    plan = []
    for (directory, category), group in groups.items():
        # Sort by modified date
        group.sort(key=lambda x: x["modified"])
        pad = len(str(len(group)))
        for i, f in enumerate(group, 1):
            original = f["name"]
            ext = Path(original).suffix.lower()
            # Clean the stem
            stem = remove_clutter(original)
            stem = clean_separators(stem)
            stem = Path(stem).stem
            new_name = f"{str(i).zfill(pad)}_{stem}{ext}"

            if new_name != original:
                plan.append({
                    "source": f["path"],
                    "original_name": original,
                    "new_name": new_name,
                    "sequence": i,
                    "category": category,
                    "reason": "sequential_numbering",
                })
    return plan


STRATEGY_MAP = {
    "clean": strategy_clean,
    "standardize": strategy_standardize,
    "date-prefix": strategy_date_prefix,
    "sequential": strategy_sequential,
}


# ── Conflict Resolution ─────────────────────────────────────────────────────

def resolve_conflicts(plan: list) -> list:
    """Resolve naming conflicts where multiple files would get the same name."""
    # Group by destination directory + new_name
    dest_names = {}
    for entry in plan:
        dest_dir = os.path.dirname(entry["source"])
        key = os.path.join(dest_dir, entry["new_name"])
        if key not in dest_names:
            dest_names[key] = []
        dest_names[key].append(entry)

    resolved = []
    for key, entries in dest_names.items():
        if len(entries) == 1:
            resolved.append(entries[0])
        else:
            # Add counter suffix to resolve conflicts
            for i, entry in enumerate(entries):
                if i == 0:
                    resolved.append(entry)
                else:
                    stem = Path(entry["new_name"]).stem
                    ext = Path(entry["new_name"]).suffix
                    entry["new_name"] = f"{stem}_{i+1}{ext}"
                    entry["conflict_resolved"] = True
                    resolved.append(entry)
    return resolved


# ── Execution ────────────────────────────────────────────────────────────────

def execute_renames(plan: list, dry_run: bool = True, log_file: str = None) -> dict:
    """Execute rename operations from the plan."""
    results = {
        "dry_run": dry_run,
        "timestamp": datetime.now(tz=timezone.utc).isoformat(),
        "total": len(plan),
        "completed": 0,
        "skipped": 0,
        "errors": 0,
        "operations": [],
        "undo_operations": [],
    }

    for entry in plan:
        src = entry["source"]
        new_name = entry["new_name"]
        dest = os.path.join(os.path.dirname(src), new_name)

        op = {
            "source": src,
            "destination": dest,
            "original_name": entry["original_name"],
            "new_name": new_name,
            "status": "pending",
        }

        if not os.path.isfile(src):
            op["status"] = "skipped"
            op["reason"] = "source not found"
            results["skipped"] += 1
            results["operations"].append(op)
            continue

        # Safety: never overwrite
        if os.path.exists(dest) and dest != src:
            stem = Path(new_name).stem
            ext = Path(new_name).suffix
            counter = 1
            while os.path.exists(os.path.join(os.path.dirname(src), f"{stem}_{counter}{ext}")):
                counter += 1
            new_name = f"{stem}_{counter}{ext}"
            dest = os.path.join(os.path.dirname(src), new_name)
            op["destination"] = dest
            op["new_name"] = new_name
            op["conflict_resolved"] = True

        if dry_run:
            op["status"] = "dry_run"
            results["completed"] += 1
            results["operations"].append(op)
            continue

        try:
            os.rename(src, dest)
            op["status"] = "completed"
            results["completed"] += 1
            results["undo_operations"].append({
                "action": "rename",
                "source": dest,
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
        print(f"Rename log saved to {log_file}", file=sys.stderr)

    return results


# ── Descriptive Rename Suggestions (LLM-assisted) ───────────────────────────

def identify_vague_names(files: list) -> list:
    """Identify files with vague or generic names that need descriptive renaming."""
    vague_patterns = [
        r"^untitled",
        r"^document\s*\d*$",
        r"^file\s*\d*$",
        r"^image\s*\d*$",
        r"^photo\s*\d*$",
        r"^video\s*\d*$",
        r"^audio\s*\d*$",
        r"^recording\s*\d*$",
        r"^screenshot\s*\d*$",
        r"^screen\s*shot",
        r"^img[_\-]?\d+$",
        r"^dsc[_\-]?\d+$",
        r"^dcim",
        r"^mov[_\-]?\d+$",
        r"^vid[_\-]?\d+$",
        r"^rec[_\-]?\d+$",
        r"^new\s*(file|document|folder)",
        r"^download",
        r"^attachment",
        r"^temp",
        r"^tmp",
        r"^\d+$",  # Just numbers
        r"^[a-f0-9]{8,}$",  # Hash-like names
    ]

    vague_files = []
    for f in files:
        stem = Path(f["name"]).stem.lower().strip()
        for pattern in vague_patterns:
            if re.match(pattern, stem, re.IGNORECASE):
                vague_files.append({
                    "path": f["path"],
                    "name": f["name"],
                    "category": f["category"],
                    "size_human": f["size_human"],
                    "modified": f["modified"],
                    "matched_pattern": pattern,
                    "suggestion": "Needs descriptive rename based on file content",
                })
                break
    return vague_files


def main():
    parser = argparse.ArgumentParser(description="Intelligent file renaming.")
    parser.add_argument("scan_report", help="Path to scan report JSON from scan_files.py")
    parser.add_argument("--strategy", "-s", default="clean",
                        choices=list(STRATEGY_MAP.keys()) + ["descriptive"],
                        help="Renaming strategy")
    parser.add_argument("--convention", "-c", default="snake_case",
                        choices=list(CONVENTION_MAP.keys()),
                        help="Naming convention (for 'standardize' strategy)")
    parser.add_argument("--output", "-o", default=None, help="Output rename plan JSON")
    parser.add_argument("--dry-run", action="store_true", default=True,
                        help="Preview only (default: True)")
    parser.add_argument("--execute", action="store_true",
                        help="Execute renames (overrides --dry-run)")
    parser.add_argument("--log", "-l", default=None, help="Operation log path")
    parser.add_argument("--exclude-categories", nargs="*", default=[],
                        help="Categories to exclude from renaming")
    args = parser.parse_args()

    if not os.path.isfile(args.scan_report):
        print(f"Error: '{args.scan_report}' not found.", file=sys.stderr)
        sys.exit(1)

    with open(args.scan_report, "r", encoding="utf-8") as f:
        scan_report = json.load(f)

    files = scan_report.get("files", [])

    # Filter out excluded categories
    if args.exclude_categories:
        files = [f for f in files if f["category"] not in args.exclude_categories]

    # Filter out junk files (they should be cleaned, not renamed)
    files = [f for f in files if not f.get("is_junk", False)]

    if args.strategy == "descriptive":
        # Output list of vague files for LLM-assisted renaming
        vague = identify_vague_names(files)
        result = {
            "strategy": "descriptive",
            "vague_files_count": len(vague),
            "vague_files": vague,
            "instruction": "Review these files and provide descriptive names based on content analysis.",
        }
        output = json.dumps(result, indent=2, ensure_ascii=False)
        if args.output:
            with open(args.output, "w", encoding="utf-8") as f:
                f.write(output)
            print(f"Vague files report saved to {args.output}", file=sys.stderr)
        else:
            print(output)
        return

    # Apply strategy
    strategy_fn = STRATEGY_MAP[args.strategy]
    kwargs = {"convention": args.convention}
    plan = strategy_fn(files, **kwargs)
    plan = resolve_conflicts(plan)

    if args.output and not args.execute:
        # Save plan for review
        result = {
            "strategy": args.strategy,
            "total_renames": len(plan),
            "renames": plan,
        }
        with open(args.output, "w", encoding="utf-8") as f:
            json.dump(result, f, indent=2, ensure_ascii=False)
        print(f"Rename plan saved to {args.output} ({len(plan)} renames)", file=sys.stderr)
        return

    # Execute or dry-run
    dry_run = not args.execute
    if dry_run:
        print("=== DRY RUN MODE (no files will be renamed) ===", file=sys.stderr)

    results = execute_renames(plan, dry_run=dry_run, log_file=args.log)

    print(f"\nResults:", file=sys.stderr)
    print(f"  Total renames: {results['total']}", file=sys.stderr)
    print(f"  Completed: {results['completed']}", file=sys.stderr)
    print(f"  Skipped: {results['skipped']}", file=sys.stderr)
    print(f"  Errors: {results['errors']}", file=sys.stderr)

    if not args.log and not args.output:
        print(json.dumps(results, indent=2, ensure_ascii=False))


if __name__ == "__main__":
    main()
