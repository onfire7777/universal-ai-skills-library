# Skill-Router ↔ Skill-Library Decoupling — Resolution Contract

> **Status:** PUBLISHED (FOUNDATION). **Implementers:** B2 (router), B4 (registry).
> **Refactor-only — no behavior change.** Companion: [`../STRUCTURE.md`](../STRUCTURE.md).

## 1. Goal

Today the router (`skill-router-cli`) finds the corpus (`skills/`) and registry
(`manifest.json`) by **co-location / repo-relative directory walking**. This refactor makes
the router a **standalone component** that locates the library purely through
**configuration / environment + the registry manifest** — so the router can be installed and
run independently of where the corpus lives.

## 2. Roles & boundaries

| Role | Artifact | Owns | Must NOT |
|------|----------|------|----------|
| **Router** | `skill-router-cli` → binary `skill-router` (alias `manus`) | Prompt→skill routing/scoring, CLI UX | Embed skill content; assume co-location; hardcode skill-dir names |
| **Registry** | `manifest.json` (root) | The single index the router consumes | Exist in more than one source |
| **Corpus** | `skills/<kebab-name>/` | `SKILL.md` + `scripts/` + `references/` | Know about the router |
| **Consumer** | host AI agent / hooks | Calls `skill-router route …` | — |

## 3. Resolution contract (B2 implements)

The router resolves three things, each via an explicit override chain.
**Repo-relative `..` walking is a *last-resort fallback only*, never the primary mechanism.**

### 3.1 Skills corpus directory — `SkillsDir()`
1. `SKILL_ROUTER_SKILLS_DIR` (env)
2. `MANUS_SKILLS_DIR` (env, legacy alias — keep)
3. `skills_dir` in `~/.skill-router/config.json`
4. Installed default `~/.agent/skills` (OpenSkills standard)

### 3.2 Registry manifest — `loadManifest()`
1. `SKILL_ROUTER_REPO_DIR` → `<repo>/manifest.json` (env).
   *Recommended addition:* an explicit `SKILL_ROUTER_MANIFEST` file override for full decoupling.
2. `MANUS_REPO_DIR` (env, legacy)
3. `repo_dir` in config.json (validated)
4. Upward search from cwd / exe dir for a **repo marker** = a directory containing **both**
   `manifest.json` and `skills/` (`isRepoDir`) — *fallback only*.
5. Standard home locations (`~/universal-ai-skills-library`, …).

Loader today: `loadManifest()` → `os.ReadFile(RepoDir()/manifest.json)`.

### 3.3 External skill roots (overlay)
- `SKILL_ROUTER_EXTERNAL_SKILL_ROOTS` (path list); cache TTL via `SKILL_ROUTER_EXTERNAL_CACHE_TTL_MINUTES`.

### 3.4 Environment-variable surface (authoritative)

| Variable | Purpose | Status |
|----------|---------|--------|
| `SKILL_ROUTER_SKILLS_DIR` | Skills corpus root | primary |
| `SKILL_ROUTER_REPO_DIR` | Library repo root (holds `manifest.json` + `skills/`) | primary |
| `SKILL_ROUTER_CONFIG_DIR` | Config dir (default `~/.skill-router`) | primary |
| `SKILL_ROUTER_MCP_DIR` | MCP bridges data dir | primary |
| `SKILL_ROUTER_EXTERNAL_SKILL_ROOTS` | Extra skill roots overlay | primary |
| `SKILL_ROUTER_EXTERNAL_CACHE_TTL_MINUTES` | External-root cache TTL | primary |
| `SKILL_ROUTER_HOOK_EVENT` | Hook event context | primary |
| `SKILL_ROUTER_PAPERCLIP_SKILLS_DIR` / `..._INSTRUCTIONS_DIR` | Paperclip adapter roots | primary |
| `MANUS_SKILLS_DIR` | Legacy alias of `SKILL_ROUTER_SKILLS_DIR` | **legacy — keep** |
| `MANUS_REPO_DIR` | Legacy alias of `SKILL_ROUTER_REPO_DIR` | **legacy — keep** |
| `MANUS_API_KEY` | Manus API key | config |

## 4. Registry manifest schema (B4 produces — single source)

`manifest.json` (currently v2.2.8):

- **Top-level:** `version, generated, description, canonical_id_policy, core_skills[],
  library_skills[], total_skills, alias_count, routing`.
- **Entry:** `{ name, directory ("skills/<name>"), description, has_scripts, scripts[], aliases[]? }`.
- **`routing`:** `{ primary_access, legacy_access ("manus skill <name>"), search, list }`.

**Requirements:**
- Exactly **one** manifest is authoritative. B4 merges the `manus-skills-marketplace` registry
  **into** this file (dedup by canonical kebab id; old/source names become aliases).
- Preserve `routing.legacy_access` and `manus` alias semantics.
- `directory` values are **relative to the corpus root**, so they survive relocation
  (top-level *or* `packages/`).

## 5. Decoupling work items — B2 (router)

- Replace co-location assumptions with the resolution chains in §3; keep `isRepoDir` /
  upward-search **only** as a last-resort fallback.
- **Remove hardcoded skill-dir references baked into commands** — drive them from the
  manifest/config instead of compiled-in constants:
  `cmd/music`→`music-prompter`, `cmd/chat`→`chat-summarizer`, `cmd/web`→`similarweb-analytics`,
  `cmd/oracle`→`multi-model-oracle`, `cmd/models`→`model-selector`, `cmd/files`→`file-organizer`,
  `cmd/audit` (generic join), and the `universal-ai-config` references.
- Preserve primary binary `skill-router` **and** the legacy `manus` alias.

## 6. Testing requirement (eliminates baseline failures)

- Routing/scoring tests **MUST run against a pinned fixture library** (a small, committed
  fixture manifest + skills under `testdata/`), **NOT** the live `skills/` tree. This makes
  assertions deterministic and removes drift failures
  (`TestPreflightRoutesAcrossCanonicalLibrary`, `TestPreflightProvidesHostReviewForAmbiguousRoute`).
- Inject the fixture via `SKILL_ROUTER_SKILLS_DIR` / `SKILL_ROUTER_REPO_DIR` — the **same
  contract** the runtime uses, so tests exercise the real resolution path.
- Make platform-sensitive tests OS-aware: `TestIsUnsafeManifestDir` (separator semantics) and
  `TestRepoDirFindsCurrentCheckoutFromNestedDirectory` (resolve symlinks before comparing).

## 7. Invariants (see STRUCTURE.md §5)

`manus` alias · single registry · kebab canonical ids · refactor-only / local-only.

## 8. Contract acceptance checklist

- [ ] Router builds & runs with the corpus at a **non-default** location set **only** via env/config (no relative co-location).
- [ ] `manus skill <name>` still works.
- [ ] Exactly one `manifest.json` is read.
- [ ] Routing tests are green against a **pinned fixture**; the 4 baseline failures are gone.
- [ ] No new `go vet` / `go build` breakage.
