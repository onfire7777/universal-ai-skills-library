#!/usr/bin/env python3
"""Seed routing-eval cases from the existing corpus — Phase 0 (§9.2).

Deterministically derives labeled routing cases from data the corpus *already*
carries, so the broad regression net needs no hand-labeling and stays in sync
with the library:

  * one ``corpus-description`` case per skill — the skill's own description as the
    prompt, expecting that skill. This is a sanity/regression net: the router must
    be able to find a skill from its own description, and Phase 1 must not regress
    it. (It is intentionally easy — the *hard* paraphrase cases live in the curated
    cases.jsonl.)
  * one ``corpus-trigger`` case per entry of every skill's ``triggers:`` block
    (58 skills today) — author-provided example prompts == free positive labels,
    exactly as the plan intends triggers to seed the eval.

Output is JSONL sorted by (id) so the file is byte-stable across runs (no
timestamps, no randomness). NOTHING here runs the router or changes runtime.

Usage:
  python3 tests/routing-eval/seed_corpus.py                       # -> cases.corpus.jsonl
  python3 tests/routing-eval/seed_corpus.py --out -               # -> stdout
  python3 tests/routing-eval/seed_corpus.py --kind triggers        # only trigger-derived cases
  python3 tests/routing-eval/seed_corpus.py --max-desc-chars 200    # truncate long descriptions
"""
from __future__ import annotations

import argparse
import json
import os
import re
from typing import Dict, List

HERE = os.path.dirname(os.path.abspath(__file__))
REPO_ROOT = os.path.abspath(os.path.join(HERE, "..", ".."))
MANIFEST = os.path.join(REPO_ROOT, "manifest.json")
SKILLS_DIR = os.path.join(REPO_ROOT, "skills")
DEFAULT_OUT = os.path.join(HERE, "cases.corpus.jsonl")


def manifest_skills() -> List[Dict]:
    with open(MANIFEST, encoding="utf-8") as fh:
        m = json.load(fh)
    skills = list(m.get("core_skills") or []) + list(m.get("library_skills") or [])
    return sorted(skills, key=lambda s: s["name"])


def read_triggers(skill_dir: str) -> List[str]:
    """Parse the ``triggers:`` block (a YAML string sequence) from a SKILL.md."""
    path = os.path.join(SKILLS_DIR, skill_dir, "SKILL.md")
    try:
        text = open(path, encoding="utf-8").read().replace("\r\n", "\n").lstrip("﻿")
    except OSError:
        return []
    fm = re.match(r"^---\n(.*?)\n---", text, re.S)
    if not fm:
        return []
    body = fm.group(1).split("\n")
    out: List[str] = []
    in_block = False
    for line in body:
        if re.match(r"^triggers:\s*$", line):
            in_block = True
            continue
        if in_block:
            if re.match(r"^\S", line):  # next top-level key ends the block
                break
            m = re.match(r"^\s*-\s*(.*)$", line)
            if m:
                val = m.group(1).strip()
                q = re.match(r'^(["\'])(.*)\1$', val)
                out.append(q.group(2) if q else val)
    return out


def main() -> int:
    ap = argparse.ArgumentParser(description="Seed routing-eval cases from the corpus.")
    ap.add_argument("--out", default=DEFAULT_OUT, help="output JSONL path, or '-' for stdout")
    ap.add_argument("--kind", choices=["all", "descriptions", "triggers"], default="all")
    ap.add_argument("--max-desc-chars", type=int, default=240,
                    help="truncate description prompts to N chars (0 = no limit)")
    args = ap.parse_args()

    cases: List[Dict] = []
    for s in manifest_skills():
        name = s["name"]
        directory = s.get("directory", f"skills/{name}").rsplit("/", 1)[-1]
        if args.kind in ("all", "descriptions"):
            desc = re.sub(r"\s+", " ", (s.get("description") or "")).strip()
            if desc:
                if args.max_desc_chars and len(desc) > args.max_desc_chars:
                    desc = desc[: args.max_desc_chars].rstrip()
                cases.append({
                    "id": f"corpus-desc:{name}",
                    "prompt": desc,
                    "expect": [name],
                    "reject": [],
                    "tags": ["corpus-description"],
                })
        if args.kind in ("all", "triggers"):
            for i, trig in enumerate(read_triggers(directory)):
                trig = re.sub(r"\s+", " ", trig).strip()
                if trig:
                    cases.append({
                        "id": f"corpus-trig:{name}:{i}",
                        "prompt": trig,
                        "expect": [name],
                        "reject": [],
                        "tags": ["corpus-trigger"],
                    })

    cases.sort(key=lambda c: c["id"])
    lines = [json.dumps(c, ensure_ascii=False, sort_keys=True) for c in cases]
    payload = "\n".join(lines) + ("\n" if lines else "")

    if args.out == "-":
        import sys
        sys.stdout.write(payload)
    else:
        with open(args.out, "w", encoding="utf-8") as fh:
            fh.write(payload)
        print(f"wrote {len(cases)} cases -> {os.path.relpath(args.out, REPO_ROOT)}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
