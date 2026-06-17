# Registry generator — single source of truth

This directory unifies what used to be **multiple hand-maintained registries**
(which drifted: `manifest.json` listed 1,812 skills while `docs/build_manifest.json`
listed 1,811, and `marketplace.json` was byte-duplicated under `plugin/`) into
**one source → many generated artifacts**.

```
SOURCE                                        GENERATED (never hand-edit)
─────────────────────────────────────────    ───────────────────────────────────
skills/<id>/SKILL.md   (the catalog)      ┐   manifest.json                 (router catalog)
scripts/registry/registry.config.json     ├──▶ marketplace.json              (Claude plugin marketplace)
   (curated, non-derivable metadata)      │   .agents/plugins/marketplace.json (codex variant, lockstep)
                                          ┘   docs/build_manifest.json        (provenance / build report)
```

Because every artifact is produced from the same scan in one run, they can no
longer drift from each other or from the `skills/` tree.

> **The generator is the Go command `skill-router registry build`**
> (`skill-router-cli/internal/registry`). It replaced the former Node generator
> `generate-registry.mjs` at proven byte-parity; the Node tool was removed at the
> end of the cut-over (see `docs/MIGRATION_NODE_TO_GO.md`). `registry.config.json`
> and the shared `lib/frontmatter.mjs` (used by `validate-schema.mjs`) remain.

## What lives where

- **`skills/<id>/SKILL.md`** — the catalog. A skill's canonical id is its
  kebab-case directory name; its description is the front-matter `description`.
- **`registry.config.json`** — the only hand-maintained file. Holds the things
  that *cannot* be derived from disk:
  - `coreSkills` — ordered list of the 18 core skills (library is auto-sorted)
  - `aliases` — per-skill compatibility aliases
  - `descriptionOverrides` — the handful of router-tuned descriptions that
    intentionally differ from the skill's own `SKILL.md`
  - `manifest` / `marketplace` / `buildManifest` — package + provenance metadata
  - `groupings` — the 14 themed plugin collections imported from the former
    manus marketplace (its only non-duplicate value)

## Commands

```bash
# From the repo root use the `skill-router` binary, or from skill-router-cli/
# run `go run . registry build ...`.

# Validate that the committed registries are in sync with skills/ (CI guard).
# Default checks ALL four artifacts + the stale-duplicate guard.
skill-router registry build --check

# Regenerate the registries on disk.
skill-router registry build --write

# Inspect one artifact without writing.
skill-router registry build --print manifest

# Reproduce the legacy manifest/marketplace byte-for-byte.
skill-router registry build --check --faithful

# Tests: the shared frontmatter parser lib (the Go generator is tested under
# skill-router-cli/internal/registry/).
node --test scripts/registry/lib/frontmatter.test.mjs
```

Optimize is the **default** (the committed registries are the optimized output):
drop empty optional fields relying on the Go reader's `omitempty`; slim
`build_manifest` to provenance only — the catalog lives in `manifest.json`;
portable relative paths; recomputed counts; themed groupings in
`marketplace.json`. `--faithful` instead reproduces the legacy
`manifest.json`/`marketplace.json` byte-for-byte (refactor-only proof).

## Contract with the router (do not break)

`skill-router-cli` reads `manifest.json` at the repo root via `loadManifest()`
and validates it with `validate-manifest`. The generator therefore guarantees:

- schema `core_skills[]` / `library_skills[]`, entries `{name, directory,
  description, aliases?, has_scripts?, scripts?}`, `directory` relative to repo root
- `has_scripts` / `scripts` computed exactly like `validate_manifest.go`
  (`listSkillScripts`: recursive, skips `__pycache__`/`.git` and `.pyc`/`.pyo`, sorted)
- **the `manus` legacy command alias is preserved**:
  `manifest.routing.legacy_access = "manus skill <name>"` and
  `build_manifest.legacy_binary_alias = "manus"`
- `merged_legacy_directories`, `disabled_colliding_aliases`, `compatibility_policy`
  carried through unchanged

## Bootstrap provenance

`seed-config.mjs` captured the curated data from the legacy hand-authored
registries exactly once; `load-groupings.mjs` imported + validated the 14 themed
groupings. Neither needs to be run again in normal operation.

## Notable corrections (behaviour-neutral)

- `docs/build_manifest.json` no longer lags `manifest.json` (single scan ⇒ equal counts).
- `build_manifest` paths are portable (were hardcoded Windows `%USERPROFILE%`).
- `manifest.alias_count` is recomputed deterministically (the legacy `1917` was
  not reproducible from the catalog; it is now the actual alias-string count).
- `scripts[]` are emitted in the canonical sorted order the Go validator expects.
- `plugin/marketplace.json` (a stray byte-duplicate of the root aggregate) was
  collapsed — the root `marketplace.json` is the only canonical marketplace and
  the plugin is self-described by `plugin/plugin.json`. `--check` lists it under
  `STALE_REGISTRIES` and fails if it ever reappears.
