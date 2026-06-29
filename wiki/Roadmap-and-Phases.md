# Roadmap & Phases

UASL's work is organized into phases, from the stable pre-phase baseline through router intelligence, an MCP service, a feedback loop, composition, and distribution.

Source docs: [`docs/PHASE_STATUS.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/docs/PHASE_STATUS.md) · [`docs/ADAPTER_DEPRECATION.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/docs/ADAPTER_DEPRECATION.md) · [`docs/DISTRIBUTION_STRATEGY.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/docs/DISTRIBUTION_STRATEGY.md)

---

## Phases

| Phase | Theme | Status |
|---|---|---|
| **PRE-0** | Baseline corpus + working router | ✅ Stable |
| **0** | **Foundations** — canonical router verbs, the Go CLI, layer decoupling, the [Node → Go migration](Node-to-Go-Migration.md), goreleaser packaging, schema + routing-eval harness | ⏳ In progress |
| **1** | Semantic routing — close the recall gap on paraphrased prompts | Planned |
| **2** | MCP service — the [MCP server](MCP-Server.md) surface over the shared engine | Planned / underway |
| **3** | Feedback loop — telemetry + reranker learning from routing decisions | Planned / underway |
| **4** | Composition — multi-step [DAG `compose --pipeline`](Skill-Router-CLI.md#core-routing-commands) | Planned / underway |
| **5** | Distribution — signed per-skill packages | Planned |

> [!NOTE]
> Phase numbers describe capability tracks, not a strict sequence — several later-phase capabilities (MCP serve, reranker, compose pipeline) already exist in the engine while Phase 0 foundations are still being finalized. Treat `docs/PHASE_STATUS.md` as the live source of truth.

### Phase 1 target (routing quality)

Phase 1 must improve routing recall while regressing nothing:

- **Recall@5:** +≥15 points over the Phase-0 baseline (0.6879).
- **P@1:** +≥10 points over baseline (0.5319).
- **Abstention:** zero regressions (baseline 1.0000).

These are measured by the routing-eval harness — see [Performance](Performance-and-Benchmarks.md#routing-quality-eval-baseline) and [Testing & CI](Testing-and-CI.md#routing-eval-harness).

---

## Adapter deprecation

The integration model changed fundamentally:

| | Old model | New model |
|---|---|---|
| **How clients get skills** | Physical-copy propagation — a copy of the corpus in each downstream agent | One canonical engine (`skillservice`) reached via CLI or MCP |
| **Drift** | Each copy drifts independently | Single source of truth; no drift |
| **Client footprint** | Full corpus per client | A compact wrapper per client |

The old copy-propagation adapters are deprecated in favor of the router + wrapper approach described in [Architecture](Architecture.md#adapter-model-from-copy-propagation-to-one-canonical-engine).

---

## Distribution strategy

Distribution is a **two-layer** model:

| Layer | What ships | How |
|---|---|---|
| **A — Binary** | The `skill-router` binary | **goreleaser** builds macOS (universal), Linux, and Windows artifacts; the **release** workflow signs checksums with **cosign** (keyless / sigstore OIDC) and can publish a Homebrew tap. |
| **B — Corpus** | The skills themselves | Today: whole-repo. Phase 5: per-skill **signed packages** (minisign + cosign). |

See [Testing & CI](Testing-and-CI.md#release) for the release workflow, and the [Node → Go migration](Node-to-Go-Migration.md) for how the registry generator consolidates into the single binary.

---

## Related pages

- [Architecture](Architecture.md) — the structure these phases build on
- [Node → Go Migration](Node-to-Go-Migration.md) — Phase-0 foundation work
- [Performance & Benchmarks](Performance-and-Benchmarks.md) — the metrics phases are measured against
