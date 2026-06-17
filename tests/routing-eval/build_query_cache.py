#!/usr/bin/env python3
"""Build the offline query-vector cache for deterministic hybrid eval (Phase 1).

Embeds every prompt in a cases file with the local model (Ollama) and writes
{prompt: vector} JSON. The route path reads this via SKILL_ROUTER_QUERY_CACHE so
the semantic eval runs deterministically and Ollama-free in CI (the committed
routing-index.bin + this cache fully determine routing). Re-run after editing
cases.jsonl. NOT used at production query time (that embeds live).

Usage:
  python3 tests/routing-eval/build_query_cache.py                 # cases.jsonl -> query-vectors.json
  python3 tests/routing-eval/build_query_cache.py --model nomic-embed-text
"""
from __future__ import annotations

import argparse
import json
import os
import urllib.request

HERE = os.path.dirname(os.path.abspath(__file__))
DEFAULT_CASES = os.path.join(HERE, "cases.jsonl")
DEFAULT_OUT = os.path.join(HERE, "query-vectors.json")
DEFAULT_MODEL = "nomic-embed-text"


def embed(text: str, model: str, host: str) -> list:
    req = urllib.request.Request(
        host.rstrip("/") + "/api/embeddings",
        data=json.dumps({"model": model, "prompt": text}).encode(),
        headers={"Content-Type": "application/json"},
    )
    with urllib.request.urlopen(req, timeout=30) as resp:
        return json.load(resp)["embedding"]


def main() -> int:
    ap = argparse.ArgumentParser(description="Build the offline query-vector cache.")
    ap.add_argument("--cases", default=DEFAULT_CASES)
    ap.add_argument("--out", default=DEFAULT_OUT)
    ap.add_argument("--model", default=DEFAULT_MODEL)
    ap.add_argument("--host", default=os.environ.get("OLLAMA_HOST", "http://localhost:11434"))
    args = ap.parse_args()
    if not args.host.startswith("http"):
        args.host = "http://" + args.host

    prompts = []
    with open(args.cases, encoding="utf-8") as fh:
        for line in fh:
            line = line.strip()
            if not line or line.startswith("#"):
                continue
            p = json.loads(line).get("prompt")
            if p and p not in prompts:
                prompts.append(p)

    cache = {p: embed(p, args.model, args.host) for p in prompts}
    with open(args.out, "w", encoding="utf-8") as fh:
        json.dump(cache, fh, separators=(",", ":"), sort_keys=True)
        fh.write("\n")
    print(f"wrote {len(cache)} query vectors ({args.model}) -> {os.path.relpath(args.out, os.path.join(HERE, '..', '..'))}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
