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

## Baseline → Phase 1 (hybrid semantic routing)

55 cases (47 positive, 8 negative). `baseline/metrics.json` now records the
**Phase-1 hybrid** numbers (the no-regression gate floor); the Phase-0 lexical
column is the comparison point.

| Metric | Phase 0 (lexical) | Phase 1 (hybrid) | Δ |
|--------|------------------:|-----------------:|---:|
| P@1 | 0.5319 | 0.5957 | **+6.4 pts** |
| MRR | 0.6193 | 0.7220 | **+10.3 pts** |
| Recall@5 | 0.6879 | 0.8582 | **+17.0 pts** |
| abstention | 1.0000 | 1.0000 | — |

Recall@5 clears the plan's `+≥15 pts` bar (+17.0) and MRR clears `+10`. P@1
improves +6.4 — short of the aspirational +10, which the plan defers to the
**Phase-3 trained reranker** (RRF ships as the identity reranker in Phase 1).
Exact name/alias routing is preserved by the deterministic guardrail (0
exact-match regressions).

### Determinism / how the gain is reproduced offline

The hybrid path needs (a) `routing-index.bin` (committed, hash-pinned) and (b) a
query embedding. For CI/eval, `query-vectors.json` (built by `build_query_cache.py`,
committed) supplies the query vectors via `SKILL_ROUTER_QUERY_CACHE`, so
`run_eval.py` reproduces the hybrid metrics **deterministically and without
Ollama**. At production query time the router embeds live (local Ollama) and
falls back to lexical-only routing when the embedder or index is absent.

Rebuild the index + cache after corpus/case changes:

```bash
skill-router index build                                   # -> routing-index.bin (+ .sha256)
python3 tests/routing-eval/build_query_cache.py             # -> query-vectors.json
python3 tests/routing-eval/run_eval.py --baseline           # re-record metrics.json
```
