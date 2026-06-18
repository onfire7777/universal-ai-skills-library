# schemas/

The build-time frontmatter contract for skills. Phase 0 of
[`docs/ARCHITECTURE_IMPROVEMENT_PLAN.md`](../docs/ARCHITECTURE_IMPROVEMENT_PLAN.md)
(§3.2, §9.1).

## `skill.schema.json`

JSON Schema (draft 2020-12) for the YAML frontmatter of every
`skills/<id>/SKILL.md`. It is the documented contract that Phase 1+ consumers
(the Go `registry build`, editors, future validators) can share.

**Required today:** only `name` and `description` — the two fields all current
corpus skills already carry and the only two the router reads. Everything else
(`version`, `domain`, `tags`, `triggers`, `capabilities`, `requires`,
`quality_tier`, `supersedes`, `conflicts`, `aliases`) is **optional** and
backfilled incrementally. `additionalProperties` is `true`, so existing keys
(`license`, `metadata`, `allowed-tools`, …) and unschema'd skills keep
validating unchanged.

> Adding the schema changes **nothing** at runtime — the router's behavior is
> unaffected. The schema only makes the metadata Phases 1–5 depend on
> *checkable*.

## Validating

```bash
make schema                                   # --warn (default): report + coverage, exit 0
node scripts/registry/validate-schema.mjs --error   # strict: exit 1 on any violation
node scripts/registry/validate-schema.mjs --json     # machine-readable report
```

The validator (`scripts/registry/validate-schema.mjs`) is **zero-dependency**
(matching the registry toolchain) and is driven *by* this schema file, so it
cannot drift from the contract. It prints a **per-field coverage report** — the
backfill tracker — and runs in `--warn` mode in CI (the `registry-drift` job)
during incremental adoption. It flips to `--error` in a later phase once the
corpus is backfilled.

### `--warn`-then-`--error` ramp

| Field | Coverage at Phase 0 | Notes |
|-------|--------------------|-------|
| `name`, `description` | 100% | required |
| `requires` | ~45% | nested `{mcp,tools,env}` |
| `triggers` | ~3% | author example prompts → seed `tests/routing-eval` |
| `version`, `aliases` | <1% | |
| `domain`, `tags`, `capabilities`, `quality_tier`, `supersedes`, `conflicts` | 0% | new in Phase 0 |

A handful of skills currently mis-file natural-language phrases as kebab
`aliases` (e.g. `crawl4ai`, `instagram-cli`); the validator surfaces these as
warnings for backfill — they are real hygiene issues, not false positives.
