# Performance & Bloat — Before / After (Goal 3)

Hard numbers for the *"reduce bloat / improve performance"* goal, measured on canonical
`main` before (`eec826d`) and after (`3c285e0`) B2 (router decouple), B4 (registry single-source)
and B5 (cleanup) landed. Methodology and the harness live in [`../bench/`](../bench/);
the raw diff is [`../bench/results/delta.md`](../bench/results/delta.md).

### Results (TL;DR)

| Headline | Before → After | Δ |
|----------|----------------|---|
| **Route decision** (`preflight`) | 167.9 → 129.0 ms | **−23.2%** |
| **Registry JSON** (manifest + build_manifest) | 1.50 MB → 0.67 MB | **−774 KB / −20.8%** |
| `build_manifest.json` | 750 → 5.6 KB | −99.2% |
| Go manifest parse | 4.78 → 4.03 ms | −15.7% |
| Repo (excl `.git`) | 136.34 → 135.55 MB | −782 KB |

Wins concentrate in **registry slimming** (B4) and **routing latency** (B2). Corpus (132 MB)
and fonts (21.6 MB) are unchanged — see *Outcomes by hot-spot* for why, and the remaining opportunity.

```bash
bench/run.sh before    # bench/results/before.json   (eec826d)
bench/run.sh after     # bench/results/after.json    (3c285e0)
bench/compare.py bench/results/before.json bench/results/after.json --out bench/results/delta.md
```

## BEFORE snapshot

- Captured: 2026-06-17 · `git_head` **eec826d** · go 1.26.3 · Darwin arm64
- Note: "before" = canonical `main` at capture time (other agents' early commits may
  already be present); the harness is re-runnable so AFTER is an apples-to-apples diff.

### Latency (wall-clock, median of 15 runs, hermetic)

| Metric | Command | BEFORE | AFTER | Δ |
|--------|---------|-------:|------:|---|
| Router init | `skill-router --version` | 8.17 ms | 10.84 ms | +2.67 ms (+32.7%) |
| Manifest load + render | `skill-router skills list` | 15.78 ms | 15.82 ms | +0.04 ms (≈0%) |
| **Route decision** | `skill-router preflight "<prompt>"` | **167.88 ms** (min 135.16) | **128.97 ms** (min 120.30) | **−38.91 ms (−23.2%)** |
| ↳ routing overhead | `preflight − version` | 159.71 ms | 118.13 ms | −41.58 ms (−26.0%) |
| ↳ manifest overhead | `list − version` | 7.61 ms | 4.98 ms | −2.63 ms |

> **Route decision (the key UX metric) dropped 23%** after B2's config-driven resolver +
> B4's slimmer registry. Router init rose ~2.7 ms (resolver does a little more config work
> at startup, and process-spawn jitter dominates a ~10 ms figure) — a sub-3 ms absolute move
> that is dwarfed by the 39 ms cut in route decision.

### Manifest parse (pure Go `json.Unmarshal`, IO excluded, median of 60)

| Metric | BEFORE | AFTER | Δ |
|--------|-------:|------:|---|
| Parse time | 4.78 ms (min 4.36) | 4.03 ms (min 3.77) | −0.75 ms (−15.7%) |
| Heap alloc / parse | 1.30 MB | 1.27 MB | −38.5 KB (−3.0%) |
| Skills in manifest | 1,813 (18 core + 1,795 library) | 1,813 (18 + 1,795) | 0 |

### Size & on-disk footprint

| Metric | BEFORE | AFTER | Δ |
|--------|-------:|------:|---|
| Router binary (`go build`) | 10.46 MB | 10.46 MB | −128 B (≈0%) |
| `manifest.json` | 747.2 KB | 668.8 KB | **−78.5 KB (−10.5%)** |
| `docs/build_manifest.json` | 750.0 KB | 5.6 KB | **−744.4 KB (−99.2%)** |
| Repo (excl `.git`) | 136.34 MB | 135.55 MB | −782 KB (−0.6%) |
| `skills/` directory | 132.00 MB | 132.00 MB | 0 (see note) |
| Tracked file bytes (excl `.git`) | 116.03 MB | 115.29 MB | −734 KB (−0.6%) |
| JSON files | 1,171 files · 3.72 MB | 1,170 files · 2.94 MB | **−1 file · −774 KB (−20.8%)** |
| **Font files** | 216 files · 21.65 MB | 216 files · 21.65 MB | 0 (see note) |
| Skill directories | 1,813 | 1,813 | 0 |
| Max RSS (preflight) | 21.79 MB | 21.66 MB | −131 KB (−0.6%) |

(`.git` is excluded from refactor metrics; it grew 59.07 → 60.17 MB from the refactor's own commit history.)

## Outcomes by hot-spot

What the BEFORE analysis flagged vs. what the refactor actually moved (✅ landed · ◻ remaining):

1. **Dual registry → single source ✅ (B4).** `build_manifest.json` collapsed 750 KB → 5.6 KB
   (−99.2%) and `manifest.json` slimmed 747 → 669 KB (−10.5%); total JSON **−774 KB (−20.8%)**,
   with Go parse time −15.7% and per-parse alloc −3.0%. Drift source removed (single generator + CI guard).
2. **Routing latency ✅ (B2).** The config-driven resolver cut **route decision 23.2%**
   (167.9 → 129.0 ms) and routing overhead 26% — while keeping `go test` green (no behavior change).
3. **`skills/` corpus — 132 MB, unchanged ◻ (B3).** B3 consolidation was a **verified no-op**:
   the former marketplace source proved a fully-redundant mirror (every skill already in `skills/`;
   166 collisions all kept-canonical), so the win was *confirming zero duplication* rather than
   deleting bytes. Consolidation is achieved at the **repo** level by keeping `skills/` as the
   only canonical corpus. Largest standing footprint; no reduction this round.
4. **Fonts — 216 files / 21.65 MB, unchanged ◻.** Out of scope for B2/B4/B5 this round; the single
   biggest **remaining** bloat opportunity (dedupe / externalize embedded font assets in skills).
5. **Router binary — 10.46 MB, unchanged ◻ (B2, optional).** `-ldflags "-s -w"` would trim it;
   deliberately not applied to preserve debuggability.

**Net:** the refactor's wins are concentrated in the **registry** (−774 KB JSON, single source)
and **routing latency** (−23%), exactly the "reduce bloat + improve performance + clean separation"
goal. The corpus/font footprint is the clear next target if a follow-up pass is wanted.

## Acceptance for this task

- [x] BEFORE snapshot captured (`bench/results/before.json`, git `eec826d`) + reproducible harness committed.
- [x] AFTER snapshot captured (`bench/results/after.json`, git `3c285e0`).
- [x] Delta table generated (`bench/results/delta.md`) and folded into this report.
