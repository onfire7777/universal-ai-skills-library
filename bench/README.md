# `bench/` — performance & bloat measurement harness

Reproducible measurement of the skill-router + library, used to quantify the
**Goal 3** ("reduce bloat / improve performance") refactor with hard before/after
numbers. See the rendered report in [`../docs/PERFORMANCE.md`](../docs/PERFORMANCE.md).

## Run

```bash
bench/run.sh before     # baseline on canonical main (capture BEFORE the refactor)
# ... B2 (router) / B4 (registry slim) / B5 (asset dedup) land ...
bench/run.sh after      # re-measure; then diff results/before.json vs results/after.json
```

Each run builds the router fresh with `go build` (matching CI) into a temp dir —
**no binary artifact is committed** — then writes `results/<label>.json`.

## What is measured

| Group | Metric | How |
|-------|--------|-----|
| Latency | `router_init_version` | `skill-router --version` — binary init, no manifest |
| Latency | `manifest_load_list` | `skill-router skills list` — manifest read + render |
| Latency | `route_decision_preflight` | `skill-router preflight "<prompt>"` — full route decision (init + manifest + corpus scoring), no skill body loaded |
| Latency | `*_overhead` (derived) | `list - version`, `preflight - version` |
| Parse | `manifest_parse_go` | standalone Go `json.Unmarshal` of `manifest.json`, median of N (IO excluded) — see `parsebench/` |
| Size | `router_binary`, `manifest_json`, `build_manifest_json`, `repo_total/excl_git`, `git`, `skills_dir` | `stat` / `du` |
| Footprint | `file_bytes_total`, `json_count/bytes`, `font_count/bytes`, `skills_count` | one `os.walk` (skips `.git`) |
| Memory | `max_rss_preflight` | `/usr/bin/time -l` peak RSS (macOS) |

**Methodology notes**
- Timings are wall-clock medians of 15 runs (3 warmup) via Python `subprocess`.
- Routing is made hermetic & deterministic: AI API keys are emptied and
  `SKILL_ROUTER_REPO_DIR` is pinned to the repo, so preflight never hits the network.
- Route-decision latency scales with **corpus size**, so B3 dedupe + B4 manifest
  slim are expected to move `route_decision_preflight` and `manifest_parse_go`.

## Files

- `run.sh` — orchestrator (build → measure → write JSON).
- `measure.py` — all measurements; writes `results/<label>.json`.
- `parsebench/` — standalone Go module timing pure `manifest.json` parse.
- `results/` — committed snapshots (`before.json`, later `after.json`).
