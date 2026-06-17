#!/usr/bin/env python3
"""Routing eval harness — Phase 0 of ARCHITECTURE_IMPROVEMENT_PLAN.md (§3.5, §9.3).

Makes routing *measurable*. It runs the deterministic, offline ``preflight``
command over a labeled case set and reports **P@1, MRR, Recall@5** so any future
routing change (Phase 1 semantic routing, Phase 3 reranker) can be proven to
improve recall *without regressing* exact-match behavior.

Why this is safe for Phase 0:
  * It changes NOTHING at runtime — it only *observes* the existing router via
    its documented, network-free ``preflight --json --explain`` contract.
  * It is deterministic: the router is pinned to the live ``skills/`` corpus with
    HOME/config redirected to throwaway dirs, so host-specific external skill
    roots (~/.claude, ~/.codex, …) cannot leak in and perturb results — the same
    isolation the characterization suite uses, but over the full live corpus
    instead of the tiny fixture.
  * Zero third-party deps; reuses tests/characterization/harness.py for the
    single source of truth on *how to obtain a runnable router binary*
    (SKILL_ROUTER_BIN in CI, else build once with the Go toolchain).

Usage:
  python3 tests/routing-eval/run_eval.py                      # human report, curated cases.jsonl
  python3 tests/routing-eval/run_eval.py --json               # machine-readable metrics
  python3 tests/routing-eval/run_eval.py --cases <file.jsonl> # eval a different set (e.g. corpus seed)
  python3 tests/routing-eval/run_eval.py --limit 100          # cap cases (useful for the big corpus set)
  python3 tests/routing-eval/run_eval.py --baseline           # write baseline/metrics.json
  python3 tests/routing-eval/run_eval.py --check              # CI gate: fail on regression vs baseline
  python3 tests/routing-eval/run_eval.py --show-misses        # list cases the router got wrong
"""
from __future__ import annotations

import argparse
import json
import os
import sys
import tempfile
from typing import Dict, List, Optional

HERE = os.path.dirname(os.path.abspath(__file__))
REPO_ROOT = os.path.abspath(os.path.join(HERE, "..", ".."))
CASES_DEFAULT = os.path.join(HERE, "cases.jsonl")
BASELINE_PATH = os.path.join(HERE, "baseline", "metrics.json")

# Reuse the characterization harness purely for router-binary discovery/build.
sys.path.insert(0, os.path.join(REPO_ROOT, "tests", "characterization"))
import harness  # noqa: E402  (single source of truth for router_binary())

RANK_DEPTH = 8  # `preflight --explain` emits up to 8 ranked candidates ("top").


def eval_env() -> Dict[str, str]:
    """Isolated env that pins the router to the LIVE corpus, deterministically.

    Points the router at the real ``skills/`` tree (so the full 1,812-skill
    corpus is in play) but redirects HOME / USERPROFILE / config to throwaway
    temp dirs so no host-specific *external* skill roots are discovered. Colour
    is disabled for stable text. Mirrors harness.fixture_env() but for the live
    repo rather than the frozen fixture.
    """
    iso = tempfile.mkdtemp(prefix="routing-eval-home-")
    cfg = os.path.join(iso, "cfg")
    os.makedirs(cfg, exist_ok=True)
    env = {
        "SKILL_ROUTER_REPO_DIR": REPO_ROOT,
        "SKILL_ROUTER_SKILLS_DIR": os.path.join(REPO_ROOT, "skills"),
        "SKILL_ROUTER_CONFIG_DIR": cfg,
        "HOME": iso,
        "USERPROFILE": iso,
        "NO_COLOR": "1",
        "CLICOLOR": "0",
    }
    # Phase 1: when the committed offline query-vector cache exists, point the
    # router at it so the hybrid semantic path runs deterministically and without
    # Ollama (CI-safe). The router auto-loads routing-index.bin next to the
    # manifest; with no cache AND no live embedder it falls back to lexical.
    cache = os.path.join(HERE, "query-vectors.json")
    if os.path.isfile(cache):
        env["SKILL_ROUTER_QUERY_CACHE"] = cache
    return env


def load_cases(path: str, limit: Optional[int]) -> List[Dict]:
    cases: List[Dict] = []
    with open(path, encoding="utf-8") as fh:
        for line_no, line in enumerate(fh, 1):
            line = line.strip()
            if not line or line.startswith("#"):
                continue
            try:
                obj = json.loads(line)
            except json.JSONDecodeError as e:
                raise SystemExit(f"{path}:{line_no}: invalid JSON: {e}")
            obj.setdefault("id", f"case-{line_no}")
            obj.setdefault("expect", [])
            obj.setdefault("reject", [])
            obj.setdefault("tags", [])
            if not obj.get("prompt"):
                raise SystemExit(f"{path}:{line_no}: case missing 'prompt'")
            cases.append(obj)
            if limit and len(cases) >= limit:
                break
    return cases


def route(prompt: str, env: Dict[str, str]) -> Dict:
    """Return parsed ``preflight --json --explain`` output for one prompt."""
    proc = harness.run_router(
        ["preflight", "--json", "--explain", prompt],
        env=env, cwd=REPO_ROOT, timeout=120,
    )
    if not proc.stdout.strip():
        raise AssertionError(f"empty preflight output for {prompt!r}: {proc.stderr}")
    return json.loads(proc.stdout)


def ranked_names(result: Dict) -> List[str]:
    """Ranked candidate ids (score-desc). Prefer `top` (explain); else best/second."""
    top = result.get("top") or []
    if top:
        return [c["name"] for c in top]
    out = []
    for key in ("best", "second"):
        c = result.get(key)
        if c and c.get("name"):
            out.append(c["name"])
    return out


def score_case(case: Dict, result: Dict) -> Dict:
    """Compute per-case metrics from a routing result."""
    decision = result.get("decision")
    best = (result.get("best") or {}).get("name")
    ranks = ranked_names(result)
    expect = list(case.get("expect") or [])
    reject = list(case.get("reject") or [])

    rec = {
        "id": case["id"], "prompt": case["prompt"], "tags": case.get("tags", []),
        "decision": decision, "best": best, "expect": expect, "reject": reject,
        "ranked": ranks[:RANK_DEPTH],
    }

    if expect:  # positive case
        # First-hit rank (1-based) of any expected skill in the ranked list.
        first_rank = 0
        for i, name in enumerate(ranks, 1):
            if name in expect:
                first_rank = i
                break
        rec["kind"] = "positive"
        rec["p_at_1"] = 1.0 if ranks[:1] and ranks[0] in expect else 0.0
        rec["mrr"] = (1.0 / first_rank) if first_rank else 0.0
        hits5 = sum(1 for n in set(expect) if n in ranks[:5])
        rec["recall_at_5"] = hits5 / len(set(expect))
        rec["routed_to_expected"] = bool(decision == "route" and best in expect)
        rec["hit"] = rec["p_at_1"] == 1.0
    else:  # negative / abstention case
        rec["kind"] = "negative"
        routed_to_rejected = bool(decision == "route" and (best in reject if reject else True))
        rec["abstained"] = not routed_to_rejected
        rec["hit"] = rec["abstained"]
    return rec


def aggregate(records: List[Dict]) -> Dict:
    pos = [r for r in records if r["kind"] == "positive"]
    neg = [r for r in records if r["kind"] == "negative"]
    decisions: Dict[str, int] = {}
    for r in records:
        decisions[r["decision"]] = decisions.get(r["decision"], 0) + 1

    def mean(xs):
        return round(sum(xs) / len(xs), 4) if xs else 0.0

    return {
        "n_cases": len(records),
        "n_positive": len(pos),
        "n_negative": len(neg),
        "p_at_1": mean([r["p_at_1"] for r in pos]),
        "mrr": mean([r["mrr"] for r in pos]),
        "recall_at_5": mean([r["recall_at_5"] for r in pos]),
        "routed_to_expected_rate": mean([1.0 if r["routed_to_expected"] else 0.0 for r in pos]),
        "abstention_accuracy": mean([1.0 if r["abstained"] else 0.0 for r in neg]),
        "decision_distribution": dict(sorted(decisions.items())),
    }


# Metrics that are gated for no-regression in --check (higher is better).
GATED = ["p_at_1", "mrr", "recall_at_5", "abstention_accuracy"]


def main() -> int:
    ap = argparse.ArgumentParser(description="Skill-router routing eval (P@1 / MRR / Recall@5).")
    ap.add_argument("--cases", default=CASES_DEFAULT, help="path to a cases.jsonl file")
    ap.add_argument("--limit", type=int, default=None, help="evaluate at most N cases")
    ap.add_argument("--json", action="store_true", help="emit machine-readable metrics JSON")
    ap.add_argument("--baseline", action="store_true", help="write metrics to baseline/metrics.json")
    ap.add_argument("--check", action="store_true", help="CI gate: fail if metrics regress vs baseline")
    ap.add_argument("--tolerance", type=float, default=0.0, help="allowed regression per gated metric")
    ap.add_argument("--show-misses", action="store_true", help="print cases the router got wrong")
    args = ap.parse_args()

    cases = load_cases(args.cases, args.limit)
    try:
        harness.router_binary()  # resolve/build once up-front; clearer error if absent
    except harness.RouterUnavailable as e:
        print(f"SKIP: {e}", file=sys.stderr)
        return 0  # router toolchain absent → skip, do not fail (matches characterization policy)

    env = eval_env()
    records = [score_case(c, route(c["prompt"], env)) for c in cases]
    metrics = aggregate(records)
    metrics["cases_file"] = os.path.relpath(args.cases, REPO_ROOT)

    if args.baseline:
        os.makedirs(os.path.dirname(BASELINE_PATH), exist_ok=True)
        with open(BASELINE_PATH, "w", encoding="utf-8") as fh:
            json.dump(metrics, fh, indent=2, sort_keys=True)
            fh.write("\n")
        print(f"baseline written: {os.path.relpath(BASELINE_PATH, REPO_ROOT)}")

    if args.json:
        print(json.dumps(metrics, indent=2, sort_keys=True))
    else:
        print(f"\nRouting eval — {metrics['n_cases']} cases "
              f"({metrics['n_positive']} positive, {metrics['n_negative']} negative)")
        print("─" * 60)
        print(f"  P@1                    {metrics['p_at_1']:.4f}")
        print(f"  MRR                    {metrics['mrr']:.4f}")
        print(f"  Recall@5               {metrics['recall_at_5']:.4f}")
        print(f"  routed-to-expected     {metrics['routed_to_expected_rate']:.4f}")
        print(f"  abstention accuracy    {metrics['abstention_accuracy']:.4f}")
        print(f"  decisions              {metrics['decision_distribution']}")
        print("─" * 60)

    if args.show_misses:
        misses = [r for r in records if not r["hit"]]
        print(f"\n{len(misses)} miss(es):")
        for r in misses:
            tgt = "/".join(r["expect"]) if r["kind"] == "positive" else f"reject {r['reject']}"
            print(f"  [{r['decision']:>9}] {r['id']}: want {tgt}; "
                  f"best={r['best']}; top={r['ranked'][:3]}")
            print(f"            prompt: {r['prompt']}")

    if args.check:
        if not os.path.isfile(BASELINE_PATH):
            print(f"--check: no baseline at {BASELINE_PATH}; run --baseline first", file=sys.stderr)
            return 1
        with open(BASELINE_PATH, encoding="utf-8") as fh:
            base = json.load(fh)
        regressions = []
        for k in GATED:
            delta = metrics.get(k, 0.0) - base.get(k, 0.0)
            if delta < -args.tolerance:
                regressions.append(f"{k}: {base.get(k):.4f} -> {metrics.get(k):.4f} ({delta:+.4f})")
        if regressions:
            print("REGRESSION vs baseline:", file=sys.stderr)
            for r in regressions:
                print(f"  {r}", file=sys.stderr)
            return 1
        print("OK: no regression vs baseline.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
