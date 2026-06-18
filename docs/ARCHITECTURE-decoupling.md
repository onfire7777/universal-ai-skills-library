# Skill-Router ↔ Skill-Library Decoupling — Resolution Contract

> **Status:** PUBLISHED (FOUNDATION). **Implementers:** B2 (router), B4 (registry).
> **Compatibility-preserving modernization.** Companion: [`../STRUCTURE.md`](../STRUCTURE.md).

## 1. Goal

Today the router (`skill-router-cli`) finds the corpus (`skills/`) and registry
(`manifest.json`) by **co-location / repo-relative directory walking**. This refactor makes
the router a **standalone component** that locates the library purely through
**configuration / environment + the registry manifest** — so the router can be installed and
run independently of where the corpus lives.

## 2. Roles & boundaries

| Role | Artifact | Owns | Must NOT |
|------|----------|------|----------|
| **Router** | `skill-router-cli` → binary `skill-router` plus compatibility aliases | Prompt→skill routing/scoring, CLI UX | Embed skill content; assume co-location; hardcode skill-dir names |
| **Registry** | `manifest.json` (root) | The single index the router consumes | Exist in more than one source |
| **Corpus** | `skills/<kebab-name>/` | `SKILL.md` + `scripts/` + `references/` | Know about the router |
| **Consumer** | host AI agent / hooks | Calls `skill-router route …` | — |

## 3. Resolution contract (B2 implements)

Authoritative resolution order (confirmed with B2). Each resolver is an explicit override
chain; **repo-relative `..` walking is a *last-resort fallback only*, never the primary
mechanism.** The two NEW resolvers (`SkillSourceDir`, `ManifestPath`) are **additive and
default to the canonical repo; old `MANUS_*` path aliases are retired in favor of
`SKILL_ROUTER_*` environment variables.

### 3.1 Skills corpus directory — `SkillsDir()`
1. `SKILL_ROUTER_SKILLS_DIR` (env)
2. `skills_dir` in config.json
3. Installed default `~/.agent/skills` (OpenSkills standard)

### 3.2 Library repo root — `RepoDir()`
1. `SKILL_ROUTER_REPO_DIR` (env)
2. `repo_dir` in config.json (validated)
3. Upward search from **cwd** for a **repo marker** = a dir containing **both** `manifest.json`
   and `skills/` (`isRepoDir`) — *fallback only*
4. Upward search from the **executable** dir — *fallback only*
6. Home candidates for the canonical repo name only (`~/universal-ai-skills-library`, `~/repos/universal-ai-skills-library`, `~/Documents/universal-ai-skills-library`)

### 3.3 Skills corpus source — `SkillSourceDir()` *(NEW, additive)*
1. `SKILL_ROUTER_SKILLS_SOURCE_DIR` (env)
2. `skills_source_dir` in config.json
3. Default `RepoDir()/skills`

### 3.4 Registry manifest path — `ManifestPath()` *(NEW, additive)*
1. `SKILL_ROUTER_MANIFEST` (env)
2. `manifest_path` in config.json
3. Default `RepoDir()/manifest.json`

Loader: `loadManifest()` → `os.ReadFile(ManifestPath())`.

### 3.5 Config dir — `ConfigDir()`
1. `SKILL_ROUTER_CONFIG_DIR` (env)
2. Default `~/.skill-router`

### 3.6 External skill roots (overlay)
- `SKILL_ROUTER_EXTERNAL_SKILL_ROOTS` (path list); cache TTL via `SKILL_ROUTER_EXTERNAL_CACHE_TTL_MINUTES`.

### 3.7 Environment-variable surface (authoritative)

| Variable | Purpose | Status |
|----------|---------|--------|
| `SKILL_ROUTER_SKILLS_DIR` | Skills corpus root | primary |
| `SKILL_ROUTER_REPO_DIR` | Library repo root (holds `manifest.json` + `skills/`) | primary |
| `SKILL_ROUTER_SKILLS_SOURCE_DIR` | Skills corpus source override (default `RepoDir()/skills`) | primary (additive) |
| `SKILL_ROUTER_MANIFEST` | Manifest file path override (default `RepoDir()/manifest.json`) | primary (additive) |
| `SKILL_ROUTER_CONFIG_DIR` | Config dir (default `~/.skill-router`) | primary |
| `SKILL_ROUTER_MCP_DIR` | MCP bridges data dir | primary |
| `SKILL_ROUTER_EXTERNAL_SKILL_ROOTS` | Extra skill roots overlay | primary |
| `SKILL_ROUTER_EXTERNAL_CACHE_TTL_MINUTES` | External-root cache TTL | primary |
| `SKILL_ROUTER_HOOK_EVENT` | Hook event context | primary |
| `SKILL_ROUTER_PAPERCLIP_SKILLS_DIR` / `..._INSTRUCTIONS_DIR` | Paperclip adapter roots | primary |
| `SKILL_ROUTER_API_BASE` | Compatible provider API base URL override | provider config |
| `SKILL_ROUTER_API_KEY` | Compatible provider API key | config |

## 4. Registry manifest schema (B4 produces — single source)

`manifest.json` (currently v2.2.8):

- **Top-level:** `version, generated, description, canonical_id_policy, core_skills[],
  library_skills[], total_skills, alias_count, routing`.
- **Entry:** `{ name, directory ("skills/<name>"), description, has_scripts, scripts[], aliases[]? }`.
- **`routing`:** `{ primary_access, compatibility_access[], search, list }`.

**Requirements:**
- Exactly **one** manifest is authoritative. Historical marketplace metadata is already
  folded into this file (dedup by canonical kebab id; old/source names become aliases).
- Preserve `routing.compatibility_access[]` and compatibility alias semantics.
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
- Preserve primary binary `skill-router` **and** compatibility aliases.

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

compatibility aliases · single registry · kebab canonical ids · universal active naming.

## 8. Contract acceptance checklist

- [ ] Router builds & runs with the corpus at a **non-default** location set **only** via env/config (no relative co-location).
- [ ] Compatibility command aliases still work.
- [ ] Exactly one `manifest.json` is read.
- [ ] Routing tests are green against a **pinned fixture**; the 4 baseline failures are gone.
- [ ] No new `go vet` / `go build` breakage.
