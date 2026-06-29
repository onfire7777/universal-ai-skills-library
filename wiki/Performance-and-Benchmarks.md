# Performance & Benchmarks

Routing is meant to be fast and hermetic — the deterministic [preflight](Architecture.md#stage-2--route--preflight-deterministic-decision) makes its decision with **no network calls**. This page collects the measured numbers and how to reproduce them.

Source: [`docs/PERFORMANCE.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/docs/PERFORMANCE.md) · [`bench/`](https://github.com/onfire7777/universal-ai-skills-library/tree/main/bench)

> All latencies are **median of 15 runs (3 warmup)**, hermetic (API keys emptied, repo dir pinned, no network).

---

## Latency

| Metric | Command | Before | After | Δ |
|---|---|---|---|---|
| Router init | `skill-router --version` | 8.17 ms | 10.84 ms | +2.67 ms |
| Manifest load | `skill-router skills list` | 15.78 ms | 15.82 ms | ≈0 |
| **Route decision** (key UX) | `skill-router preflight "<prompt>"` | **167.88 ms** | **128.97 ms** | **−38.91 ms (−23.2%)** |
| Routing overhead | preflight − version | 159.71 ms | 118.13 ms | −41.58 ms (−26.0%) |
| Manifest overhead | list − version | 7.61 ms | 4.98 ms | −2.63 ms (−34.6%) |

The headline result: the **route decision dropped 23.2%** after the config-driven resolver and a slimmer registry. Router init rose ~2.7 ms (sub-3 ms absolute), dwarfed by the ~39 ms gain on the decision path.

---

## Registry size

| Artifact | Before | After | Δ |
|---|---|---|---|
| `docs/build_manifest.json` | 750 KB | 5.6 KB | **−99.2%** |
| `manifest.json` | 747 KB | 669 KB | −10.5% |
| **Total registry JSON** | 1.50 MB | 0.67 MB | **−774 KB (−20.8%)** |
| Go manifest parse | 4.78 ms | 4.03 ms | −15.7% |

Standing footprint (unchanged, noted for future work): skills corpus ~132 MB (no duplication found), font assets ~21.65 MB across 216 files (largest remaining dedupe target).

---

## Routing-quality (eval) baseline

The Phase-0 baseline over a curated, hand-labeled set of **55 cases** (47 positive, 8 negative):

| Metric | Value |
|---|---|
| **P@1** (precision @ rank 1) | 0.5319 |
| **MRR** (mean reciprocal rank) | 0.6193 |
| **Recall@5** | 0.6879 |
| **Abstention accuracy** | 1.0000 |

Paraphrase cases are intentionally missed at this baseline — they quantify the recall gap that [Phase 1](Roadmap-and-Phases.md#phase-1-target-routing-quality) semantic routing must close. See [Testing & CI](Testing-and-CI.md#routing-eval-harness) for the harness.

---

## Reproducing the benchmarks

The suite lives in [`bench/`](https://github.com/onfire7777/universal-ai-skills-library/tree/main/bench):

| File | Role |
|---|---|
| `run.sh` | Orchestrator: build → measure → write JSON. |
| `measure.py` | Captures latency, size, and footprint metrics. |
| `parsebench/` | Standalone Go module timing pure `manifest.json` `Unmarshal` (IO excluded). |
| `results/before.json`, `results/after.json` | Committed snapshots. |
| `compare.py` | Diffs two snapshots into a `delta.md` report. |

```bash
bench/run.sh before     # → bench/results/before.json
bench/run.sh after      # → bench/results/after.json
bench/compare.py bench/results/before.json bench/results/after.json --out bench/results/delta.md
```

**Methodology:** hermetic and deterministic — AI API keys emptied, `SKILL_ROUTER_REPO_DIR` pinned, corpus hash fixed so results are reproducible across runs. Route-decision latency scales with corpus size, which is why corpus dedupe and manifest slimming both move this metric.

---

## Related pages

- [Architecture](Architecture.md) — why search scores against the manifest, not the corpus
- [Testing & CI](Testing-and-CI.md) — the eval gate that protects these numbers
- [Roadmap & Phases](Roadmap-and-Phases.md) — the Phase-1 routing-quality targets
