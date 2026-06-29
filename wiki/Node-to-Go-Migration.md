# Node → Go Migration

The registry generator is migrating from a legacy **Node.js** implementation to the **Go** router (`skill-router registry build`). The migration is staged and guarded by a **byte-parity** guarantee: the Go generator must emit output that is byte-for-byte identical to the Node generator before Node is removed.

Source: [`docs/MIGRATION_NODE_TO_GO.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/docs/MIGRATION_NODE_TO_GO.md)

---

## Why migrate

- Collapse the toolchain to a **single Go binary** — the same binary that routes also generates the registry, so there's no separate Node dependency at runtime.
- Keep the [registry contract](Manifest-and-Registry.md) identical while swapping the implementation underneath it — exactly what the [decoupled architecture](Architecture.md) was built to allow.

---

## The byte-parity guarantee

The safety net is `scripts/registry/parity-check.sh`, run locally via `make parity` and in CI by the **registry-parity** workflow. It proves the Go builder produces **byte-identical** output to the Node generator (`generate-registry.mjs`) for both the committed "optimize" mode and the "faithful" mode. If they ever diverge, the gate turns red.

```bash
make parity                              # local byte-parity check
bash scripts/registry/parity-check.sh    # same, directly
```

This is what makes the cut-over safe: nothing ships unless Go == Node, byte-for-byte.

---

## Cut-over stages

| Stage | State |
|---|---|
| **0–2** | Foundations — Go builder implemented and brought to parity. ✅ Complete. |
| **3** | **Current.** Go is the **authoritative** builder; Node runs as a **non-blocking parity oracle**. `registry build --check` gates the committed artifacts; `parity-check.sh` guards against divergence. |
| **4** | Remove the Node generator after one clean release soak. ⏳ Pending. |
| **5** | Post-removal cleanup. ⏳ Pending. |

In Stage 3, the **registry-parity** and **release** workflows both run `registry build --check` plus `parity-check.sh`, so a release is refused if Go and Node disagree.

---

## What this means for contributors

- Treat **Go as the owner** of the registry. Regenerate with `make registry-write` (the Go path).
- The Node files in `scripts/registry/` (`generate-registry.mjs`, `lib/frontmatter.mjs`, their tests) are the **oracle** — keep them working until Stage 4 removes them.
- Any change to skills, the generator, or the router that affects the manifest must keep `make parity` green.

See [Testing & CI](Testing-and-CI.md#registry-parity) for the gate details and [Manifest & Registry](Manifest-and-Registry.md) for the artifacts involved.
