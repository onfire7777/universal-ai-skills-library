#!/usr/bin/env python3
"""Diff two bench snapshots into a markdown delta table.

  bench/compare.py results/before.json results/after.json [--out results/delta.md]

Prints (and optionally writes) before/after/delta for every comparable metric,
so the final Goal-3 report can quote hard numbers. Safe to run before `after.json`
exists (it just reports that the after snapshot is missing).
"""
import argparse
import json
import os
import sys

# (json path, label, unit) — path is dotted into the snapshot dict.
METRICS = [
    ("latency_ms.route_decision_preflight.median", "Route decision (preflight)", "ms"),
    ("latency_ms.manifest_load_list.median", "Manifest load (skills list)", "ms"),
    ("latency_ms.router_init_version.median", "Router init (--version)", "ms"),
    ("latency_ms.routing_overhead", "Routing overhead", "ms"),
    ("manifest_parse_go.parse_ms_median", "Go manifest parse", "ms"),
    ("manifest_parse_go.alloc_bytes_parse", "Parse alloc", "bytes"),
    ("manifest_parse_go.total_skills", "Manifest skills", "count"),
    ("size_bytes.router_binary", "Router binary", "bytes"),
    ("size_bytes.manifest_json", "manifest.json", "bytes"),
    ("size_bytes.build_manifest_json", "build_manifest.json", "bytes"),
    ("size_bytes.repo_excl_git", "Repo (excl .git)", "bytes"),
    ("size_bytes.skills_dir", "skills/ dir", "bytes"),
    ("footprint.file_bytes_total", "Tracked file bytes", "bytes"),
    ("footprint.json_count", "JSON files", "count"),
    ("footprint.json_bytes", "JSON bytes", "bytes"),
    ("footprint.font_count", "Font files", "count"),
    ("footprint.font_bytes", "Font bytes", "bytes"),
    ("footprint.skills_count", "Skill dirs", "count"),
    ("memory.max_rss_preflight", "Max RSS (preflight)", "bytes"),
]


def dig(d, path):
    for k in path.split("."):
        if not isinstance(d, dict) or k not in d:
            return None
        d = d[k]
    return d


def human(v, unit):
    if v is None:
        return "—"
    if unit == "bytes":
        if v >= 1_000_000:
            return f"{v/1_000_000:.2f} MB"
        if v >= 1_000:
            return f"{v/1_000:.1f} KB"
        return f"{v} B"
    if unit == "ms":
        return f"{v:.2f} ms"
    return str(v)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("before")
    ap.add_argument("after")
    ap.add_argument("--out")
    a = ap.parse_args()

    with open(a.before) as f:
        before = json.load(f)
    if not os.path.exists(a.after):
        print(f"after snapshot not found: {a.after}\nRun `bench/run.sh after` first.",
              file=sys.stderr)
        sys.exit(2)
    with open(a.after) as f:
        after = json.load(f)

    lines = [
        f"| Metric | Before (`{before.get('git_head','?')}`) | After (`{after.get('git_head','?')}`) | Delta |",
        "|--------|--------|-------|-------|",
    ]
    for path, label, unit in METRICS:
        b, c = dig(before, path), dig(after, path)
        if b is None and c is None:
            continue
        delta = ""
        if isinstance(b, (int, float)) and isinstance(c, (int, float)) and b:
            pct = (c - b) / b * 100.0
            delta = f"{human(c-b, unit)} ({pct:+.1f}%)"
        lines.append(f"| {label} | {human(b, unit)} | {human(c, unit)} | {delta} |")

    table = "\n".join(lines)
    print(table)
    if a.out:
        with open(a.out, "w") as f:
            f.write(table + "\n")
        print(f"\nwrote {a.out}", file=sys.stderr)


if __name__ == "__main__":
    main()
