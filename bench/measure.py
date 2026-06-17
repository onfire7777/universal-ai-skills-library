#!/usr/bin/env python3
"""Capture a performance/bloat snapshot of the skill-router + library.

Measures router latency (process timings, median of N runs), pure Go manifest
parse cost, on-disk footprints (binary, manifest, repo, JSON, fonts), skill
counts, and peak memory. Writes a JSON snapshot used by docs/PERFORMANCE.md to
quantify before/after deltas for the Goal-3 ("reduce bloat / improve
performance") refactor.

Usage:
  measure.py --bin <router-binary> --repo <repo-root> --label before --out results/before.json
"""
import argparse
import json
import os
import re
import statistics
import subprocess
import sys
import time
from datetime import datetime, timezone

# Latency is dominated by the corpus scan, so empty the AI keys and pin the repo
# so routing is deterministic and never reaches the network.
HERMETIC_ENV = {
    "OPENAI_API_KEY": "",
    "OPENROUTER_API_KEY": "",
    "ANTHROPIC_API_KEY": "",
}
ROUTE_PROMPT = "organize and rename messy files in this folder"
WARMUP = 3
RUNS = 15


def time_cmd(args, cwd, env, runs=RUNS, warmup=WARMUP):
    """Return (median_ms, min_ms) wall time for a command, ignoring output."""
    for _ in range(warmup):
        subprocess.run(args, cwd=cwd, env=env, stdout=subprocess.DEVNULL,
                       stderr=subprocess.DEVNULL)
    samples = []
    for _ in range(runs):
        t0 = time.perf_counter()
        subprocess.run(args, cwd=cwd, env=env, stdout=subprocess.DEVNULL,
                       stderr=subprocess.DEVNULL)
        samples.append((time.perf_counter() - t0) * 1000.0)
    return round(statistics.median(samples), 2), round(min(samples), 2)


def du_kb(path):
    """Disk usage in bytes via `du -sk` (portable across BSD/GNU)."""
    try:
        out = subprocess.check_output(["du", "-sk", path], stderr=subprocess.DEVNULL)
        return int(out.split()[0]) * 1024
    except (subprocess.CalledProcessError, ValueError, IndexError):
        return 0


def walk_footprint(root, skip_git=True):
    """Single os.walk: total file bytes, JSON count/bytes, font count/bytes."""
    FONT_EXT = (".ttf", ".otf", ".woff", ".woff2", ".eot")
    total = json_n = json_b = font_n = font_b = 0
    for dirpath, dirnames, filenames in os.walk(root):
        if skip_git and ".git" in dirnames:
            dirnames.remove(".git")
        for name in filenames:
            try:
                size = os.path.getsize(os.path.join(dirpath, name))
            except OSError:
                continue
            total += size
            low = name.lower()
            if low.endswith(".json"):
                json_n += 1
                json_b += size
            if low.endswith(FONT_EXT):
                font_n += 1
                font_b += size
    return {"file_bytes_total": total, "json_count": json_n, "json_bytes": json_b,
            "font_count": font_n, "font_bytes": font_b}


def max_rss_bytes(bin_path, repo, env):
    """Peak resident set of a route decision via `/usr/bin/time -l` (macOS)."""
    try:
        proc = subprocess.run(
            ["/usr/bin/time", "-l", bin_path, "preflight", ROUTE_PROMPT],
            cwd=repo, env=env, stdout=subprocess.DEVNULL, stderr=subprocess.PIPE)
        m = re.search(r"(\d+)\s+maximum resident set size", proc.stderr.decode())
        return int(m.group(1)) if m else None
    except OSError:
        return None


def go_parse(repo, manifest_path):
    """Run the standalone Go manifest-parse microbench; return its JSON dict."""
    pb = os.path.join(repo, "bench", "parsebench")
    try:
        out = subprocess.check_output(["go", "run", ".", manifest_path, "60"],
                                      cwd=pb, stderr=subprocess.PIPE)
        return json.loads(out.decode())
    except (subprocess.CalledProcessError, ValueError) as e:
        return {"error": str(e)}


def stat_bytes(path):
    return os.path.getsize(path) if os.path.exists(path) else 0


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--bin", required=True)
    ap.add_argument("--repo", required=True)
    ap.add_argument("--label", default="before")
    ap.add_argument("--out", required=True)
    a = ap.parse_args()

    repo = os.path.abspath(a.repo)
    env = {**os.environ, **HERMETIC_ENV, "SKILL_ROUTER_REPO_DIR": repo}
    manifest = os.path.join(repo, "manifest.json")

    print(f"[measure] label={a.label} repo={repo}", file=sys.stderr)

    ver_med, ver_min = time_cmd([a.bin, "--version"], repo, env)
    list_med, list_min = time_cmd([a.bin, "skills", "list"], repo, env)
    pf_med, pf_min = time_cmd([a.bin, "preflight", ROUTE_PROMPT], repo, env)

    foot = walk_footprint(repo)
    skills_dir = os.path.join(repo, "skills")
    skills_count = sum(
        1 for n in os.listdir(skills_dir)
        if os.path.isdir(os.path.join(skills_dir, n))) if os.path.isdir(skills_dir) else 0

    snapshot = {
        "label": a.label,
        "timestamp": datetime.now(timezone.utc).isoformat(),
        "git_head": subprocess.getoutput(f"git -C {repo} rev-parse --short HEAD"),
        "host": {
            "go": subprocess.getoutput("go version"),
            "uname": subprocess.getoutput("uname -sm"),
        },
        "latency_ms": {
            "router_init_version": {"median": ver_med, "min": ver_min},
            "manifest_load_list": {"median": list_med, "min": list_min},
            "route_decision_preflight": {"median": pf_med, "min": pf_min},
            "manifest_load_overhead": round(list_med - ver_med, 2),
            "routing_overhead": round(pf_med - ver_med, 2),
            "_method": f"wall-clock subprocess, median of {RUNS} runs ({WARMUP} warmup)",
            "_prompt": ROUTE_PROMPT,
        },
        "manifest_parse_go": go_parse(repo, manifest),
        "size_bytes": {
            "router_binary": stat_bytes(a.bin),
            "manifest_json": stat_bytes(manifest),
            "build_manifest_json": stat_bytes(os.path.join(repo, "docs", "build_manifest.json")),
            "repo_total": du_kb(repo),
            "git": du_kb(os.path.join(repo, ".git")),
            "skills_dir": du_kb(skills_dir),
        },
        "footprint": {**foot, "skills_count": skills_count},
        "memory": {"max_rss_preflight": max_rss_bytes(a.bin, repo, env)},
    }
    snapshot["size_bytes"]["repo_excl_git"] = (
        snapshot["size_bytes"]["repo_total"] - snapshot["size_bytes"]["git"])

    os.makedirs(os.path.dirname(a.out), exist_ok=True)
    with open(a.out, "w") as f:
        json.dump(snapshot, f, indent=2)
        f.write("\n")
    json.dump(snapshot, sys.stdout, indent=2)
    print()


if __name__ == "__main__":
    main()
