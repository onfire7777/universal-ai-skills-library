#!/usr/bin/env python3
"""
scan_files.py - Scan a directory and produce a structured JSON inventory.

Collects file metadata (name, extension, size, dates, hash) for all files
in the target directory. Output is a JSON report used by other scripts.

Usage:
    python scan_files.py <target_directory> [--output <report.json>] [--max-depth <N>]

Safety: Read-only. Does not modify any files.
"""

import argparse
import hashlib
import json
import os
import stat
import sys
from datetime import datetime, timezone
from pathlib import Path


# File category mappings
CATEGORY_MAP = {
    # Documents
    ".pdf": "Documents", ".doc": "Documents", ".docx": "Documents",
    ".odt": "Documents", ".rtf": "Documents", ".tex": "Documents",
    ".txt": "Documents", ".md": "Documents", ".rst": "Documents",
    ".pages": "Documents", ".epub": "Documents", ".mobi": "Documents",
    # Spreadsheets
    ".xls": "Spreadsheets", ".xlsx": "Spreadsheets", ".csv": "Spreadsheets",
    ".ods": "Spreadsheets", ".tsv": "Spreadsheets", ".numbers": "Spreadsheets",
    # Presentations
    ".ppt": "Presentations", ".pptx": "Presentations", ".odp": "Presentations",
    ".key": "Presentations",
    # Images
    ".jpg": "Images", ".jpeg": "Images", ".png": "Images", ".gif": "Images",
    ".bmp": "Images", ".tiff": "Images", ".tif": "Images", ".svg": "Images",
    ".webp": "Images", ".ico": "Images", ".heic": "Images", ".heif": "Images",
    ".raw": "Images", ".cr2": "Images", ".nef": "Images", ".arw": "Images",
    ".psd": "Images", ".ai": "Images", ".eps": "Images",
    # Videos
    ".mp4": "Videos", ".avi": "Videos", ".mkv": "Videos", ".mov": "Videos",
    ".wmv": "Videos", ".flv": "Videos", ".webm": "Videos", ".m4v": "Videos",
    ".mpg": "Videos", ".mpeg": "Videos", ".3gp": "Videos", ".ogv": "Videos",
    # Audio
    ".mp3": "Audio", ".wav": "Audio", ".flac": "Audio", ".aac": "Audio",
    ".ogg": "Audio", ".wma": "Audio", ".m4a": "Audio", ".opus": "Audio",
    ".aiff": "Audio", ".mid": "Audio", ".midi": "Audio",
    # Archives
    ".zip": "Archives", ".rar": "Archives", ".7z": "Archives",
    ".tar": "Archives", ".gz": "Archives", ".bz2": "Archives",
    ".xz": "Archives", ".tgz": "Archives", ".dmg": "Archives",
    ".iso": "Archives",
    # Code
    ".py": "Code", ".js": "Code", ".ts": "Code", ".jsx": "Code",
    ".tsx": "Code", ".java": "Code", ".c": "Code", ".cpp": "Code",
    ".h": "Code", ".hpp": "Code", ".cs": "Code", ".go": "Code",
    ".rs": "Code", ".rb": "Code", ".php": "Code", ".swift": "Code",
    ".kt": "Code", ".scala": "Code", ".r": "Code", ".m": "Code",
    ".sh": "Code", ".bash": "Code", ".ps1": "Code", ".bat": "Code",
    ".cmd": "Code", ".sql": "Code", ".html": "Code", ".css": "Code",
    ".scss": "Code", ".sass": "Code", ".less": "Code", ".vue": "Code",
    ".svelte": "Code",
    # Data
    ".json": "Data", ".xml": "Data", ".yaml": "Data", ".yml": "Data",
    ".toml": "Data", ".ini": "Data", ".cfg": "Data", ".conf": "Data",
    ".env": "Data", ".properties": "Data", ".db": "Data", ".sqlite": "Data",
    ".sqlite3": "Data", ".parquet": "Data", ".avro": "Data",
    # Fonts
    ".ttf": "Fonts", ".otf": "Fonts", ".woff": "Fonts", ".woff2": "Fonts",
    ".eot": "Fonts",
    # Executables / Installers
    ".exe": "Executables", ".msi": "Executables", ".deb": "Executables",
    ".rpm": "Executables", ".app": "Executables", ".apk": "Executables",
    ".pkg": "Executables", ".snap": "Executables", ".flatpak": "Executables",
    # Design
    ".fig": "Design", ".sketch": "Design", ".xd": "Design",
    # 3D / CAD
    ".stl": "3D_Models", ".obj": "3D_Models", ".fbx": "3D_Models",
    ".blend": "3D_Models", ".dwg": "3D_Models", ".dxf": "3D_Models",
    # Ebooks
    ".azw": "Ebooks", ".azw3": "Ebooks", ".cbr": "Ebooks", ".cbz": "Ebooks",
    # Backups / Temp
    ".bak": "Backups", ".tmp": "Temporary", ".temp": "Temporary",
    ".swp": "Temporary", ".swo": "Temporary", ".log": "Logs",
    ".cache": "Temporary",
}

# Known junk / clutter patterns (case-insensitive basenames)
JUNK_PATTERNS = {
    "thumbs.db", "desktop.ini", ".ds_store", "._*",
    "__macosx", ".spotlight-v100", ".trashes",
    ".fseventsd", ".temporaryitems", ".apdisk",
}


def file_hash(filepath: str, chunk_size: int = 65536) -> str:
    """Compute SHA-256 hash of a file for deduplication."""
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


def quick_hash(filepath: str, sample_size: int = 8192) -> str:
    """Quick hash using file size + first/last bytes for fast pre-filtering."""
    try:
        size = os.path.getsize(filepath)
        h = hashlib.md5()
        h.update(str(size).encode())
        with open(filepath, "rb") as f:
            h.update(f.read(sample_size))
            if size > sample_size * 2:
                f.seek(-sample_size, 2)
                h.update(f.read(sample_size))
        return h.hexdigest()
    except (PermissionError, OSError):
        return "UNREADABLE"


def is_junk(name: str) -> bool:
    """Check if a filename matches known junk/clutter patterns."""
    lower = name.lower()
    if lower in JUNK_PATTERNS:
        return True
    if lower.startswith("._"):
        return True
    if lower.startswith("~$"):  # Office temp files
        return True
    return False


def categorize(ext: str) -> str:
    """Return the category for a given file extension."""
    return CATEGORY_MAP.get(ext.lower(), "Other")


def scan_directory(target: str, max_depth: int = -1) -> list:
    """Recursively scan directory and collect file metadata."""
    results = []
    target = os.path.abspath(target)

    for root, dirs, files in os.walk(target):
        # Respect max depth
        depth = root[len(target):].count(os.sep)
        if max_depth >= 0 and depth >= max_depth:
            dirs.clear()
            continue

        # Skip hidden directories (but still scan if target itself is hidden)
        dirs[:] = [d for d in dirs if not d.startswith(".") or os.path.join(root, d) == target]

        for fname in files:
            fpath = os.path.join(root, fname)
            try:
                st = os.stat(fpath)
            except (PermissionError, OSError):
                continue

            if not stat.S_ISREG(st.st_mode):
                continue

            ext = Path(fname).suffix.lower()
            rel_path = os.path.relpath(fpath, target)

            entry = {
                "name": fname,
                "path": fpath,
                "relative_path": rel_path,
                "extension": ext,
                "category": categorize(ext),
                "size_bytes": st.st_size,
                "size_human": human_size(st.st_size),
                "modified": datetime.fromtimestamp(st.st_mtime, tz=timezone.utc).isoformat(),
                "created": datetime.fromtimestamp(st.st_ctime, tz=timezone.utc).isoformat(),
                "is_junk": is_junk(fname),
                "is_hidden": fname.startswith("."),
                "is_empty": st.st_size == 0,
                "quick_hash": quick_hash(fpath),
            }
            results.append(entry)

    return results


def human_size(nbytes: int) -> str:
    """Convert bytes to human-readable size."""
    for unit in ("B", "KB", "MB", "GB", "TB"):
        if abs(nbytes) < 1024:
            return f"{nbytes:.1f} {unit}"
        nbytes /= 1024
    return f"{nbytes:.1f} PB"


def build_report(target: str, files: list) -> dict:
    """Build a structured report from scanned files."""
    # Category summary
    categories = {}
    for f in files:
        cat = f["category"]
        if cat not in categories:
            categories[cat] = {"count": 0, "total_size": 0}
        categories[cat]["count"] += 1
        categories[cat]["total_size"] += f["size_bytes"]

    for cat in categories:
        categories[cat]["total_size_human"] = human_size(categories[cat]["total_size"])

    # Find potential duplicates (by quick_hash)
    hash_groups = {}
    for f in files:
        qh = f["quick_hash"]
        if qh == "UNREADABLE":
            continue
        if qh not in hash_groups:
            hash_groups[qh] = []
        hash_groups[qh].append(f["relative_path"])

    potential_duplicates = {k: v for k, v in hash_groups.items() if len(v) > 1}

    # Junk files
    junk_files = [f["relative_path"] for f in files if f["is_junk"]]

    # Empty files
    empty_files = [f["relative_path"] for f in files if f["is_empty"]]

    total_size = sum(f["size_bytes"] for f in files)

    report = {
        "scan_target": os.path.abspath(target),
        "scan_time": datetime.now(tz=timezone.utc).isoformat(),
        "total_files": len(files),
        "total_size": total_size,
        "total_size_human": human_size(total_size),
        "categories": categories,
        "potential_duplicate_groups": len(potential_duplicates),
        "potential_duplicates": potential_duplicates,
        "junk_files_count": len(junk_files),
        "junk_files": junk_files,
        "empty_files_count": len(empty_files),
        "empty_files": empty_files,
        "files": files,
    }
    return report


def main():
    parser = argparse.ArgumentParser(description="Scan directory and produce file inventory.")
    parser.add_argument("target", help="Directory to scan")
    parser.add_argument("--output", "-o", default=None, help="Output JSON file (default: stdout)")
    parser.add_argument("--max-depth", "-d", type=int, default=-1, help="Max directory depth (-1 = unlimited)")
    args = parser.parse_args()

    if not os.path.isdir(args.target):
        print(f"Error: '{args.target}' is not a valid directory.", file=sys.stderr)
        sys.exit(1)

    files = scan_directory(args.target, args.max_depth)
    report = build_report(args.target, files)

    output = json.dumps(report, indent=2, ensure_ascii=False)
    if args.output:
        with open(args.output, "w", encoding="utf-8") as f:
            f.write(output)
        print(f"Report saved to {args.output}", file=sys.stderr)
    else:
        print(output)


if __name__ == "__main__":
    main()
