# Manifest & Registry

`manifest.json` is the **single contract** between the [corpus](Skills-Corpus.md) and the [router](Skill-Router-CLI.md). It is generated — never hand-edited — and CI enforces that the committed copy matches what the generator produces.

Source: [`manifest.json`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/manifest.json) · [`scripts/registry/`](https://github.com/onfire7777/universal-ai-skills-library/tree/main/scripts/registry) · [`schemas/`](https://github.com/onfire7777/universal-ai-skills-library/tree/main/schemas)

---

## What the manifest is

A machine-readable catalog of the whole corpus. The router scores prompts against this compact index instead of reading every `SKILL.md`, which is what keeps routing fast (see [Performance](Performance-and-Benchmarks.md)).

### Top-level keys

| Key | Meaning |
|---|---|
| `version` | Schema/catalog version — currently `2.2.8`. |
| `generated` | Generation date — `2026-05-15`. |
| `description` | Human description of the catalog. |
| `canonical_id_policy` | The kebab-case ID rule (below). |
| `core_skills` | Array of the 18 core skills. |
| `library_skills` | Array of the ~1,794 library skills. |
| `total_skills` | `1812`. |
| `alias_count` | Number of legacy aliases. |
| `routing` | Routing config: `primary_access`, `compatibility_access`, `search`, `list`. |

### Per-skill entries

| Entry type | Fields |
|---|---|
| Core skill | `name`, `directory`, `description`, `has_scripts`, `scripts` |
| Library skill | `name`, `directory`, `description` |

There are two generated artifacts, kept in lockstep:

| Artifact | Role |
|---|---|
| `manifest.json` | The full router catalog (~0.67 MB after slimming). |
| `docs/build_manifest.json` | A slim provenance/build snapshot (~5.6 KB). |

---

## Canonical ID policy

> "Canonical skill ids are kebab-case top-level directory names. Legacy display names and old source names are aliases only."

So `skills/banner-design/` has canonical ID `banner-design`, and any old names (`ckm-banner-design`, `ckm:banner-design`, "Banner Design") resolve to it through `compatibility_aliases.json`.

### `compatibility_aliases.json`

| Field | Purpose |
|---|---|
| `schema` | `universal-ai-skills-aliases/v1`. |
| `generated_at` | Build timestamp. |
| `aliases` | Map of canonical ID → array of alternate names. |
| `merged_directories` | Metadata on merged source directories. |
| `disabled_colliding_aliases` | Aliases dropped due to collisions. |

Alias resolution is exact-match: a legacy or display name is looked up and rewritten to the canonical kebab ID.

---

## Generating & checking the registry

The generator lives in [`scripts/registry/`](https://github.com/onfire7777/universal-ai-skills-library/tree/main/scripts/registry) and is owned by the Go router (`skill-router registry build`). During the [Node → Go migration](Node-to-Go-Migration.md), the legacy Node generator (`generate-registry.mjs`) remains as a **byte-parity oracle**.

```bash
make registry-write    # regenerate manifest.json + docs/build_manifest.json (Go owner)
make registry-check    # verify committed artifacts match (skill-router registry build --check)
make schema            # validate every SKILL.md frontmatter against schemas/skill.schema.json
make parity            # prove Go output is byte-identical to the Node generator
```

Key generator files:

| File | Role |
|---|---|
| `generate-registry.mjs` | Legacy Node generator (parity oracle). |
| `parity-check.sh` | Proves Go output == Node output, byte-for-byte. |
| `validate-schema.mjs` | Strict frontmatter validation. |
| `registry.config.json` | Generation config — the single source of truth. |
| `lib/frontmatter.mjs` | Shared YAML frontmatter parser. |

> [!NOTE]
> The **registry-parity** CI workflow runs `registry build --check` and `parity-check.sh` on any change to `skills/**`, `scripts/registry/**`, `skill-router-cli/**`, or the artifacts themselves. A divergence turns the gate red. See [Testing & CI](Testing-and-CI.md).

---

## Schemas

`schemas/skill.schema.json` validates the `SKILL.md` frontmatter contract (the field table is on the [Skills Corpus](Skills-Corpus.md#skillmd-frontmatter-schema) page). It requires `name` and `description` and defines the optional Phase-0 routing fields (`domain`, `tags`, `triggers`, `capabilities`, `quality_tier`, `supersedes`, `conflicts`, …), currently enforced in warn mode.

---

## Related pages

- [Skills Corpus](Skills-Corpus.md) — what the manifest indexes
- [Skill Router CLI](Skill-Router-CLI.md#registry-commands) — the `registry build` command
- [Node → Go Migration](Node-to-Go-Migration.md) — why parity matters
