# Routing eval

Makes skill routing **measurable**. This is Phase 0 of
[`docs/ARCHITECTURE_IMPROVEMENT_PLAN.md`](../../docs/ARCHITECTURE_IMPROVEMENT_PLAN.md)
(§3.5, §9.2–9.3): a labeled case set plus a harness that reports **P@1, MRR,
Recall@5** so any future routing change (Phase 1 semantic routing, Phase 3
reranker) can be proven to improve recall **without regressing** today's
exact-match behavior.

> **Phase 0 changes nothing at runtime.** The harness only *observes* the
> existing router through its documented, network-free `preflight --json
> --explain` contract. No router source is modified.

## Files

| File | What |
|------|------|
| `cases.jsonl` | Curated, hand-labeled eval set (committed, reviewed). The CI gate runs on this. |
| `seed_corpus.py` | Deterministic generator: derives a broad regression net from the corpus (descriptions + author `triggers:`). Output (`cases.corpus.jsonl`) is **gitignored** — regenerate on demand. |
| `run_eval.py` | The harness: routes each case and computes metrics. |
| `baseline/metrics.json` | Committed baseline snapshot. `run_eval.py --check` gates against it. |

## Case format

One JSON object per line (`#`-comment and blank lines ignored):

```json
{"id": "para-pytest", "prompt": "write pytest fixtures …", "expect": ["python-testing-patterns"], "reject": [], "tags": ["paraphrase"]}
```

- **`expect` non-empty → positive case.** Any listed id is an acceptable top
  route. `expect` is a *human* ground-truth label of the correct skill —
  authored independently of what the router currently does. Multi-skill `expect`
  captures genuinely ambiguous prompts.
- **`expect` empty → negative / abstention case.** The router should *not* route
  (or not to anything in `reject`). Covers off-domain chit-chat and the
  uninstall-intent guardrail.

Paraphrase cases the current lexical router misses are **intentional** — they
quantify the recall gap Phase 1 must close. The baseline records them as-is.

## Running

```bash
make eval                 # human report over cases.jsonl
make eval-check           # CI gate: fail on regression vs baseline/metrics.json
make eval-baseline        # re-record baseline (ONLY after an intended change)

# or directly:
python3 tests/routing-eval/run_eval.py --show-misses     # list what the router got wrong
python3 tests/routing-eval/run_eval.py --json            # machine-readable metrics
python3 tests/routing-eval/seed_corpus.py                # write cases.corpus.jsonl (2k+ cases)
python3 tests/routing-eval/run_eval.py --cases tests/routing-eval/cases.corpus.jsonl --limit 300
```

## Determinism

Routing is a pure function of `(prompt, manifest.json, skills/)` once external
roots are excluded. The harness pins the router to the **live corpus** but
redirects `HOME`/config to throwaway dirs, so host-specific external skill roots
(`~/.claude`, `~/.codex`, …) cannot leak in — the same isolation
`tests/characterization` uses, over the full library instead of the fixture. The
router binary is obtained via the characterization harness (`SKILL_ROUTER_BIN`
in CI, else built once from source). If the Go toolchain is absent the run
**skips** (exit 0), matching the characterization policy.

## Baseline (recorded at Phase 0)

| Metric | Value |
|--------|-------|
| P@1 | 0.5319 |
| MRR | 0.6193 |
| Recall@5 | 0.6879 |
| abstention accuracy | 1.0000 |

55 cases (47 positive, 8 negative). These are the numbers Phase 1 must beat on
recall (`+≥15pts Recall@5`, `+≥10pts P@1`) with **zero** exact-match regressions.
