# Performance & Bloat — Before / After (Goal 3)

Hard numbers for the *"reduce bloat / improve performance"* goal. The **BEFORE**
column is captured on canonical `main`; the **AFTER** column is filled once B2
(router), B4 (registry slim) and B5 (asset dedup) land. Methodology and the
measurement harness live in [`../bench/`](../bench/).

```bash
bench/run.sh before                                  # done — bench/results/before.json
bench/run.sh after                                   # after B2/B4/B5 land
bench/compare.py bench/results/before.json bench/results/after.json --out bench/results/delta.md
```

## BEFORE snapshot

- Captured: 2026-06-17 · `git_head` **eec826d** · go 1.26.3 · Darwin arm64
- Note: "before" = canonical `main` at capture time (other agents' early commits may
  already be present); the harness is re-runnable so AFTER is an apples-to-apples diff.

### Latency (wall-clock, median of 15 runs, hermetic)

| Metric | Command | BEFORE | AFTER | Δ |
|--------|---------|-------:|------:|---|
| Router init | `skill-router --version` | **8.17 ms** | _TBD_ | |
| Manifest load + render | `skill-router skills list` | **15.78 ms** | _TBD_ | |
| **Route decision** | `skill-router preflight "<prompt>"` | **167.88 ms** (min 135.16) | _TBD_ | |
| ↳ routing overhead | `preflight − version` | 159.71 ms | _TBD_ | |
| ↳ manifest overhead | `list − version` | 7.61 ms | _TBD_ | |

> **Route decision is the key UX metric** (time to pick a skill for a prompt) and it
> scales with corpus size — so B3 dedupe + B4 registry slim should move it the most.

### Manifest parse (pure Go `json.Unmarshal`, IO excluded, median of 60)

| Metric | BEFORE | AFTER | Δ |
|--------|-------:|------:|---|
| Parse time | **4.78 ms** (min 4.36) | _TBD_ | |
| Heap alloc / parse | 1.30 MB | _TBD_ | |
| Skills in manifest | 1,812 (18 core + 1,794 library) | _TBD_ | |

### Size & on-disk footprint

| Metric | BEFORE | AFTER | Δ |
|--------|-------:|------:|---|
| Router binary (`go build`) | **10.46 MB** | _TBD_ | |
| `manifest.json` | 747.2 KB | _TBD_ | |
| `docs/build_manifest.json` | 750.0 KB | _TBD_ | |
| Repo (excl `.git`) | **136.34 MB** | _TBD_ | |
| `skills/` directory | **132.00 MB** | _TBD_ | |
| Tracked file bytes (excl `.git`) | 116.03 MB | _TBD_ | |
| JSON files | 1,171 files · 3.72 MB | _TBD_ | |
| **Font files** | **216 files · 21.65 MB** | _TBD_ | |
| Skill directories | 1,812 | _TBD_ | |
| Max RSS (preflight) | 21.79 MB | _TBD_ | |

(`.git` itself is 59.07 MB; repo total 195.41 MB. `.git` is excluded from refactor metrics.)

## Bloat hot-spots → owners

Ranked by expected payoff. These are **targets for other builders** — this task only
measures; it edits no source.

1. **`skills/` corpus — 132 MB / 1,812 dirs (B3).** Dominates the tree and drives
   route-decision latency. Consolidating `manus-skills-marketplace` content is reported
   by Scout 1 as ~4,814/5,008 files identical/droppable → the largest single win in both
   size and routing speed.
2. **Fonts — 216 files / 21.65 MB (B5/B3).** Binary font assets embedded inside skills.
   Dedupe/externalize for a large, low-risk size reduction.
3. **Dual registry — `manifest.json` (747 KB) + `build_manifest.json` (750 KB) ≈ 1.5 MB (B4).**
   Two near-duplicate registries; collapsing to a single source cuts JSON bytes and
   per-parse allocation, and removes a drift source.
4. **Routing overhead — ~160 ms (B2 + B3/B4).** Config-driven resolution must not regress
   it; corpus slimming should reduce it. Re-measure with `preflight`.
5. **Router binary — 10.46 MB (B2, optional).** `-ldflags "-s -w"` would trim it if a
   smaller distributable is wanted (does not change behavior).

## Acceptance for this task

- [x] BEFORE snapshot captured (`bench/results/before.json`) + reproducible harness committed.
- [ ] AFTER snapshot captured once B2/B4/B5 land (`bench/run.sh after`).
- [ ] Delta table generated (`bench/compare.py … --out bench/results/delta.md`) and folded into the final report.
