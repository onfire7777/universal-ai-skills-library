# Registry generator — single source of truth

This directory unifies what used to be **multiple hand-maintained registries**
(which drifted: `manifest.json` listed 1,812 skills while
`docs/build_manifest.json` listed 1,811) into **one source → CLI-first generated
artifacts**.

```
SOURCE                                        GENERATED (never hand-edit)
─────────────────────────────────────────    ───────────────────────────────────
skills/<id>/SKILL.md   (the catalog)      ┐   manifest.json                 (router catalog)
scripts/registry/registry.config.json     ├──▶ docs/build_manifest.json      (provenance / build report)
   (curated, non-derivable metadata)      │
                                          ┘
```

Because every artifact is produced from the same scan in one run, they can no
longer drift from each other or from the `skills/` tree.

## What lives where

- **`skills/<id>/SKILL.md`** — the catalog. A skill's canonical id is its
  kebab-case directory name; its description is the front-matter `description`.
- **`registry.config.json`** — the only hand-maintained file. Holds the things
  that *cannot* be derived from disk:
  - `coreSkills` — ordered list of the 18 core skills (library is auto-sorted)
  - `aliases` — per-skill compatibility aliases
  - `descriptionOverrides` — the handful of router-tuned descriptions that
    intentionally differ from the skill's own `SKILL.md`
  - `manifest` / `buildManifest` — package + provenance metadata

## Commands

```bash
# Validate that the committed CLI-first artifacts are in sync with skills/ (CI guard).
# Default checks manifest.json + docs/build_manifest.json and fails if any
# retired marketplace JSON reappears.
node scripts/registry/generate-registry.mjs --check

# Regenerate the registries on disk.
node scripts/registry/generate-registry.mjs --write

# Inspect one artifact without writing.
node scripts/registry/generate-registry.mjs --print manifest

# Characterization: reproduce the legacy manifest byte-for-byte.
node scripts/registry/generate-registry.mjs --check --faithful

# Tests.
node --test scripts/registry/lib/frontmatter.test.mjs scripts/registry/generate-registry.test.mjs
```

Optimize is the **default** (the committed registries are the optimized output):
drop empty optional fields relying on the Go reader's `omitempty`; slim
`build_manifest` to provenance only — the catalog lives in `manifest.json`;
portable relative paths; recomputed counts. `--faithful` instead reproduces the
legacy `manifest.json` byte-for-byte (refactor-only proof).

## Contract with the router (do not break)

`skill-router-cli` reads `manifest.json` at the repo root via `loadManifest()`
and validates it with `validate-manifest`. The generator therefore guarantees:

- schema `core_skills[]` / `library_skills[]`, entries `{name, directory,
  description, aliases?, has_scripts?, scripts?}`, `directory` relative to repo root
- `has_scripts` / `scripts` computed exactly like `validate_manifest.go`
  (`listSkillScripts`: recursive, skips `__pycache__`/`.git` and `.pyc`/`.pyo`, sorted)
- **compatibility command aliases are preserved in generic fields**:
  `manifest.routing.compatibility_access[]` and
  `build_manifest.compatibility_binary_aliases[]`
- `merged_legacy_directories`, `disabled_colliding_aliases`, `compatibility_policy`
  carried through unchanged

## Notable corrections (behaviour-neutral)

- `docs/build_manifest.json` no longer lags `manifest.json` (single scan ⇒ equal counts).
- `build_manifest` paths are portable (were hardcoded Windows `%USERPROFILE%`).
- `manifest.alias_count` is recomputed deterministically (the legacy `1917` was
  not reproducible from the catalog; it is now the actual alias-string count).
- `scripts[]` are emitted in the canonical sorted order the Go validator expects.
- `marketplace.json`, `.agents/plugins/marketplace.json`, and
  `plugin/marketplace.json` are retired. The router is CLI-first; `--check`
  lists all three under `STALE_REGISTRIES` and fails if any reappears.
