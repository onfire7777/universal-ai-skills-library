# Phase 5 — Distribution: Versioned, Signed, Pull-on-Demand Skills (RFC)

> **STATUS: DESIGN — NOT YET IMPLEMENTED. Depends on Phases 0–4.**
>
> This is an implementation-ready RFC for Phase 5 of `docs/ARCHITECTURE_IMPROVEMENT_PLAN.md`
> (roadmap §4, row 5; design §3.7). It describes per-skill semver, content-addressed
> skill packages, a signed index, and subset pull. **No behavior changes ship from this
> document.** Phase 5 is **HIGH risk, ~3–4 weeks**, and is the *last* phase precisely because
> it composes everything before it.
>
> Authoritative invariants this design must not break (from plan §6):
> universal `skill-router` binary surface · exactly one authoritative registry (`manifest.json`) ·
> deterministic offline no-remote-LLM **query** path · kebab-case canonical ids ·
> telemetry local-only + opt-in · the `--no-git` gitleaks posture ·
> do **not** relocate `skill-router-cli/` → `packages/`.

---

## 1. Goals / Non-Goals, and why Phase 5 cannot ship first

### 1.1 Goals

1. **Per-skill versioning.** Every skill carries an explicit, machine-checkable `version`
   (semver). The single `manifest.json` contract gains an additive `version` field per entry.
2. **Content-addressed packages.** Each skill is buildable into a deterministic, reproducible
   archive whose name embeds its content hash. Same bytes in → same bytes out → same hash.
3. **A signed index.** One index file enumerates every published skill `{name, version, hash,
   size, deps}` plus a root hash, signed by **both** minisign (air-gapped) and cosign (keyless
   OIDC via Sigstore). Clients verify integrity *before* install.
4. **Subset pull.** A client resolves "I need skills X, Y" → transitive dep closure → pulls only
   those packages, verifies each against the signed index, installs locally. Offline-tolerant.
5. **Replace physical-copy `DefaultSync`** (which Phase 2 deprecates) with pull-on-demand, while
   the MCP `serve` surface (Phase 2) remains the primary integration.

### 1.2 Non-Goals (explicit — won't do)

- **No network/remote-LLM at query/route time.** Distribution adds *install-time* network
  (opt-in, explicit `pull`); the route/search/load query path stays offline and deterministic.
- **No second registry.** The signed index is *derived from* `manifest.json`, not a competing
  source of truth. `manifest.json` remains the one authoritative catalog.
- **No `packages/` relocation** of the Go CLI. (Rejected upstream for blast radius.)
- **No mandatory migration.** Whole-repo install (today's model) stays supported through GA and
  beyond. Subset pull is additive opt-in.
- **No per-skill independent publishing cadence in v1.** Publishing is a single repo-level release
  event that snapshots all skills; per-skill semver tracks *change*, not *release independence*.

### 1.3 The hard dependency on Phases 0–4

Phase 5 is unbuildable today and unbuildable first. Concretely:

| Needs from prior phase | Why Phase 5 cannot proceed without it |
|---|---|
| **Phase 0** — `skill-router registry build` reproduces `manifest.json` at **byte-parity** under `--check` | Packaging + index generation must be owned by the Go binary so build, package, and sign are one deterministic pipeline. Today the build is owned by the Node generator `scripts/registry/generate-registry.mjs` (4 artifacts; the only `--check` gate). Without Phase 0, there is no Go-side build to hang packaging off, and no proof the migration is behavior-preserving. |
| **Phase 0** — JSON-Schema frontmatter (`schemas/skill.schema.json`) validates 100% of corpus | `version:` is a new frontmatter field; it needs the schema-gating machinery from Phase 0 to be enforced and backfilled across 1,813 skills. |
| **Phase 2** — `skill-router serve` (MCP) + CLI-as-thin-client; physical-copy sync deprecated | Subset pull only makes sense once the resolution/serve surface exists and `DefaultSync` is on its way out. Pull-on-demand is the *replacement* for the deprecation Phase 2 starts. |
| **Phase 1/3/4** — routing-index, telemetry, capability DAG | The index's `deps[]` edges are exactly the Phase-4 capability DAG. Without it, "transitive deps" has no source. |

**Current state note:** the `skill-router` Go binary owns registry generation
and emits the CLI-first artifacts (`manifest.json`, `docs/build_manifest.json`).
Marketplace JSON outputs are retired and guarded against reappearing. Per-skill
packaging remains a design-level extension until the package/pull lifecycle is
implemented.

---

## 2. Per-skill semver

### 2.1 The `version` frontmatter field

Today a `skills/<kebab-name>/SKILL.md` frontmatter looks like (verified):

```yaml
---
name: 21risk-automation
description: "Automate 21risk tasks via Rube MCP (Composio). Always search tools first for current schemas."
requires:
  mcp: [rube]
---
```

Phase 5 adds **one** field, `version`, and (optionally) explicit `deps`:

```yaml
---
name: 21risk-automation
version: 1.4.2                 # NEW — semver MAJOR.MINOR.PATCH, required at GA
description: "Automate 21risk tasks via Rube MCP (Composio). Always search tools first for current schemas."
requires:
  mcp: [rube]
deps:                          # NEW (optional) — other canonical skill ids this skill loads/composes
  - rube-connection-setup
---
```

### 2.2 JSON-Schema gating

Extend `schemas/skill.schema.json` (introduced in Phase 0). `version` is **optional during
alpha/beta** (warn) and **required at GA** (error):

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://onfire7777.dev/schemas/skill.schema.json",
  "type": "object",
  "required": ["name", "description"],
  "properties": {
    "name":        { "type": "string", "pattern": "^[a-z0-9]+(-[a-z0-9]+)*$" },
    "version":     { "type": "string", "pattern": "^(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(?:-[0-9A-Za-z-.]+)?$" },
    "description": { "type": "string", "minLength": 1 },
    "deps": {
      "type": "array",
      "items": { "type": "string", "pattern": "^[a-z0-9]+(-[a-z0-9]+)*$" },
      "uniqueItems": true
    }
  }
}
```

Enforced by the Phase-0 verb `skill-router skills validate-schema` (`--warn` during rollout,
hard-fail at GA). Backfill is one pass of the existing backfill script seeding `version: 1.0.0`.

### 2.3 Manual vs content-hash-driven bumping — **recommendation: content-hash-*gated*, manual bump**

Two pure options and the chosen hybrid:

- **(A) Fully manual semver** — authors hand-edit `version`. *Risk:* humans forget; 1,813 skills
  guarantees drift between "content changed" and "version bumped".
- **(B) Fully content-hash-driven** — derive the whole version from the package hash. *Problem:* a
  content hash is not ordered, so it cannot express MAJOR vs MINOR vs PATCH *intent* (a breaking
  rewrite vs a typo fix hash identically "differently"). Semantics are lost.

**Decision: hybrid — author declares the semver bump; CI *gates* it against the content hash.**
A build-time check (`skill-router registry build --check`, extended) computes each skill's
content hash, compares to the hash recorded for that `version` in the previous published index, and:

- content changed **and** version unchanged → **CI fails**: "skill `X` content changed but
  `version` still `1.4.2` — bump it."
- version changed **and** content unchanged → **CI warns** (allowed: metadata-only re-tag).
- both changed → pass.

This keeps human-meaningful semver (intent) while making "forgot to bump" impossible to merge.
The `+<shorthash>` build-metadata suffix on the *package* (see §3) is always derived, never typed.

### 2.4 Per-skill CHANGELOG expectation

Optional `skills/<name>/CHANGELOG.md` (Keep-a-Changelog format). Not gated in v1 — recommended,
lint-warned if a MAJOR bump lands with no changelog entry for that version. Kept out of the
package hash input is **wrong**; it *is* part of the skill payload, so a changelog edit alone is a
PATCH-worthy content change. Authors should bump PATCH when they touch the changelog, or batch it
with the substantive change.

### 2.5 Interaction with the single `manifest.json` contract

The router catalog gains `version` *additively* per entry. Verified current entry shape:

```json
{
  "name": "chat-summarizer",
  "directory": "skills/chat-summarizer",
  "description": "Generate a comprehensive, AI-optimized summary of the current AI session...",
  "has_scripts": true,
  "scripts": ["scripts/format_summary.py"]
}
```

Phase-5 entry shape (additive — `version` + optional `deps`; everything else byte-identical, and
the Go reader's `omitempty` keeps un-versioned entries clean during rollout):

```json
{
  "name": "chat-summarizer",
  "directory": "skills/chat-summarizer",
  "version": "2.1.0",
  "description": "Generate a comprehensive, AI-optimized summary of the current AI session...",
  "has_scripts": true,
  "scripts": ["scripts/format_summary.py"],
  "deps": ["context-anchor"]
}
```

`manifest.json` stays the **one authoritative registry**. The signed index (§4) is *generated
from* it in the same build run, so the two can never drift — same guarantee the Node generator
already provides for its four artifacts.

---

## 3. Content-addressed skill packages

### 3.1 Package tar layout

A package is the skill directory plus a per-package manifest, nothing else:

```
chat-summarizer@2.1.0+9f3a2c1.tar.zst
└── (decompressed tar)
    ├── package.json          # per-package manifest (see §3.4)
    ├── SKILL.md
    ├── scripts/
    │   └── format_summary.py
    └── references/           # present only if the skill ships references
```

`references/` and `scripts/` are included only when they exist (matching the existing
`has_scripts`/`scripts` computation rules, so packaging stays consistent with
`skill-router validate-manifest`). Nothing outside the skill directory is ever included — no repo
metadata, no `.git`, no sibling skills.

### 3.2 Content-hash algorithm (canonical tar → sha256)

The hash is `sha256` over a **canonicalized** tar stream. Canonicalization removes every source
of nondeterminism so the same skill content always hashes identically across machines, clocks,
and users:

1. **Entry order:** all paths sorted lexicographically (bytewise), `package.json` first.
2. **Timestamps:** `mtime = 0` (Unix epoch) on every entry.
3. **Ownership:** `uid = 0`, `gid = 0`, `uname = ""`, `gname = ""`.
4. **Permissions:** fixed — files `0644`, directories `0755`, no setuid/setgid/sticky bits.
5. **Type:** regular files + directories only; symlinks/devices/hardlinks rejected at build (fail
   loud — they are an exfiltration/path-escape vector).
6. **Format:** PAX/ustar with a fixed header field order; no GNU extensions; no global PAX records.
7. **Hash input is the *uncompressed* canonical tar**, not the compressed `.tar.zst`. Compression
   (zstd, default level, single-threaded for determinism) is applied *after* hashing for transport
   only. Two clients can use different zstd builds and still agree on the content hash.

```
content_hash(skill) = sha256( canonical_tar( sorted, mtime=0, uid=0, gid=0, perms=fixed ) )
shorthash           = first 7 hex chars of content_hash
```

Implementation note: Go's `archive/tar` with an explicit `tar.Header` written field-by-field
(do **not** use `tar.FileInfoHeader`, which leaks real mtime/uid/perms). This is the same
discipline the build already applies to keep `manifest.json` byte-stable.

### 3.3 Package naming

```
<canonical-skill-id>@<semver>+<shorthash>.tar.zst
```

- `<canonical-skill-id>` — kebab-case, == directory name == `manifest.json` `name` (invariant).
- `@<semver>` — the author-declared `version` (§2).
- `+<shorthash>` — semver build-metadata, == first 7 chars of the content hash. By semver rules
  build metadata is ignored for precedence (correct: precedence is by `version`; the hash is for
  *integrity*, not ordering).

Example: `chat-summarizer@2.1.0+9f3a2c1.tar.zst`.

### 3.4 Per-package manifest (`package.json` inside the tar)

```json
{
  "schema": "uasl-skill-package/v1",
  "name": "chat-summarizer",
  "version": "2.1.0",
  "content_hash": "sha256:9f3a2c1e8b7d6a5f4c3b2a1908f7e6d5c4b3a29180716253849abcdef0123456",
  "files": [
    { "path": "SKILL.md",                 "size": 4096, "sha256": "..." },
    { "path": "scripts/format_summary.py", "size": 2310, "sha256": "..." }
  ],
  "deps": ["context-anchor"],
  "requires": { "mcp": [] },
  "legacy_aliases": ["claude-native-chat-summarizer"]
}
```

`content_hash` is the §3.2 hash and is **self-excluding**: it is computed over the tar in which
`package.json` carries `content_hash: ""`, then written back. (Standard self-reference handling.)
`legacy_aliases` is carried only as historical lookup metadata so renamed skills
remain discoverable after pulled installs. It must not reintroduce legacy
ecosystem branding into package names, command names, or install targets.

### 3.5 Reproducibility guarantees

- Byte-identical skill directory → byte-identical canonical tar → identical `content_hash` and
  identical package filename, on any OS/arch (the §3.2 rules erase all platform-specific tar
  fields — directly closing the Windows/8.3-path/`pwsh` fragility class the plan calls out).
- The signed index (§4) records `content_hash`, so `skill-router pkg verify` proves a fetched
  package is bit-for-bit what was signed.
- `--check` reproducibility test in CI: build every package twice in separate temp dirs; assert
  identical hashes. Any drift fails the build (same posture as the existing `manifest.json`
  byte-parity gate).

---

## 4. Signed index

### 4.1 Index JSON schema (`skills-index.json`)

Generated in the *same* `skill-router registry build` run that emits `manifest.json`, so it can
never drift from the catalog:

```json
{
  "schema": "uasl-skills-index/v1",
  "generated_at": "2026-07-01T00:00:00.000Z",
  "registry_version": "2.3.0",
  "source_of_truth": "manifest.json",
  "primary_binary": "skill-router",
  "compatibility_binary_aliases": [],
  "root_hash": "sha256:1a2b3c...e2",
  "package_format": "uasl-skill-package/v1",
  "hash_algorithm": "sha256",
  "skills": [
    {
      "name": "chat-summarizer",
      "version": "2.1.0",
      "hash": "sha256:9f3a2c1...3456",
      "size": 6406,
      "package": "chat-summarizer@2.1.0+9f3a2c1.tar.zst",
      "deps": ["context-anchor"]
    },
    {
      "name": "context-anchor",
      "version": "1.0.3",
      "hash": "sha256:771be4...90ab",
      "size": 5120,
      "package": "context-anchor@1.0.3+771be4a.tar.zst",
      "deps": []
    }
  ]
}
```

**`root_hash`** = `sha256` over the canonical JSON serialization of the `skills` array **excluding
the `root_hash` field itself** (entries sorted by `name`, stable key order, no insignificant
whitespace). It is a Merkle-ish single-shot digest: changing any package hash, version, or dep
edge changes `root_hash`. (A full Merkle tree is deliberately *not* used in v1 — see §7; a flat
root over a sorted list is enough for "verify the whole index is intact" and is far simpler. We
revisit per-entry proofs only if partial-index streaming becomes a requirement.)

The index intentionally mirrors fields already present in `docs/build_manifest.json`
(`compatibility_binary_aliases: []`, `primary_binary: "skill-router"`) so the
supply-chain artifact and the provenance artifact stay self-consistent.

### 4.2 Signing — minisign **and** cosign (both, by design)

We sign with **two** mechanisms, serving two threat audiences:

| | **minisign** | **cosign (keyless, Sigstore OIDC)** |
|---|---|---|
| Trust root | A long-lived ed25519 public key we publish | Sigstore Fulcio CA + Rekor transparency log, tied to a GitHub Actions OIDC identity |
| Network at verify | **None** — fully air-gapped | Online (Rekor inclusion proof); can be made offline with a bundled inclusion proof |
| Audience | Air-gapped / regulated / offline-first installs | CI/CD-native installs wanting "who/what built this" provenance with no key custody |
| Key custody | We hold a private key (rotation burden) | Ephemeral keys; nothing long-lived to leak |

Artifacts produced per release:

```
skills-index.json
skills-index.json.minisig          # minisign detached signature (ed25519)
skills-index.json.cosign.bundle    # cosign keyless bundle (cert + sig + Rekor proof)
```

**We sign the index, not every package.** The index records each package's `content_hash`, so
one signature transitively covers all 1,813 packages — sign once, verify any subset. This is the
key scalability decision: per-package signatures would be 1,813 signing ops per release for zero
added security over a signed index of hashes.

### 4.3 Key management + rotation

- **minisign key** stored only in the release CI secret store; **never** committed (honors the
  `--no-git` gitleaks posture — the gitleaks gate must allowlist the *public* key file only).
- Publish the minisign public key in three places for cross-checking: `docs/SIGNING_KEYS.md`, a
  GitHub release asset, and the repo's `SECURITY.md`. Clients pin the key on first use (TOFU) and
  warn on change.
- **Rotation:** publish a new public key signed by the *old* key (a signed key-transition
  statement `keys.json` → `keys.json.minisig`) so existing clients can chain trust to the new key
  without manual re-pinning. cosign needs no rotation (ephemeral keys); its trust is the GitHub
  workflow identity, which we pin by repo + workflow ref in the client verify policy.

### 4.4 Client-side verification flow

```
fetch skills-index.json (+ .minisig and/or .cosign.bundle)
   │
   ├─ minisign: verify .minisig against pinned public key            ── fail → ABORT
   │  (offline path; default for `--offline` / air-gapped clients)
   │
   ├─ cosign (optional, if online + policy enabled):
   │     verify .cosign.bundle: cert chains to Fulcio,
   │     identity == github.com/onfire7777/...:.github/workflows/release.yml,
   │     Rekor inclusion proof valid                                 ── fail → ABORT
   │
   ├─ recompute root_hash over skills[] and compare to index.root_hash ── mismatch → ABORT
   │
   └─ index trusted ✔  → proceed to resolve + pull (§5)
```

At least one signature must verify; policy decides whether both are required. Default:
minisign-required, cosign-required-if-online.

---

## 5. Subset pull — CLI design

### 5.1 New CLI verbs

All under the existing `skill-router` binary. Two groups:

```
# Authoring / release side (CI + maintainers)
skill-router pkg build   [--only <ids>] [--out <dir>] [--check]
skill-router pkg sign    --index <skills-index.json> [--minisign-key <path>] [--cosign-keyless]
skill-router pkg verify  <package.tar.zst> --index <skills-index.json>
skill-router pkg publish --index <skills-index.json> --target <github|oci|cdn> [--dry-run]

skill-router index build [--check|--write]
skill-router index sign  [--minisign-key <path>] [--cosign-keyless]
skill-router index verify --index <url|path> [--offline]

# Client / install side (end users)
skill-router pkg pull <id|id@version> ... [--with-deps] [--offline] [--dest <dir>]
skill-router pkg ls   [--installed|--available] [--outdated]
skill-router pkg add  <id> ...     # resolve + pull + verify + install (the high-level verb)
skill-router pkg remove <id> ...   # uninstall a pulled skill (manifest-tracked, not rm -rf guesswork)
```

`pkg add` is the verb users touch; `pull` is its lower-level half (fetch + verify, no install).
`pkg build`/`index build` reuse the Phase-0 deterministic build pipeline; they are *not* a new
source of truth.

### 5.2 Resolution + transitive deps

```
add(requested_ids):
  index = index_verify(fetch_index())          # §4.4 — trust before anything else
  closure = ∅; frontier = requested_ids
  while frontier:
    id = frontier.pop()
    entry = index.skills[id]                    # missing id → hard error (no silent skip)
    closure.add(entry)
    frontier += [d for d in entry.deps if d not in closure]   # deps[] == Phase-4 capability DAG
  detect_cycles(closure)                        # DAG invariant; cycle → hard error
  for entry in closure:
    pkg = fetch(entry.package)                  # from configured host (§5.4)
    assert sha256(canonical_untar_then_retar(pkg)) == entry.hash   # INTEGRITY BEFORE INSTALL
    verify_internal(pkg.package_json)           # per-file sha256 match
  for entry in topo_order(closure):             # install deps first
    install(pkg)                                # atomic: temp dir → fsync → rename
```

Integrity is checked **before** any file lands in the install destination. A package whose hash
doesn't match the signed index is never written.

### 5.3 Offline behavior

- `--offline` (and air-gapped default): use a previously fetched, cached `skills-index.json` +
  `.minisig`; verify with the pinned minisign key (no network). Packages are pulled from a local
  cache or a mounted mirror directory; if a needed package isn't cached, **fail loudly** with the
  exact missing `name@version+hash` — never silently fall back to whole-repo or to the network.
- The **route/search/load query path remains 100% offline regardless** — pull is the only verb
  that ever touches the network, and only on explicit invocation. This preserves the hard
  invariant: no network/remote-LLM at query time.

### 5.4 Hosting options

| Option | How | Pros | Cons |
|---|---|---|---|
| **(A) GitHub Releases assets** | Each release uploads `skills-index.json`, its signatures, and the `.tar.zst` packages as release assets | Zero new infra; free; integrates natively with the cosign GitHub-OIDC identity; works with existing `gh` CLI auth | 1,813 assets per release is awkward; per-asset download; release-asset API rate limits; no content-addressed CDN semantics |
| **(B) OCI registry (ghcr.io)** | Push packages as OCI artifacts; the index is an OCI image index; cosign signs OCI-natively | First-class cosign support; content-addressable by digest *natively*; pull only the layers you need; mature CDN-backed delivery | Requires OCI tooling on the client; auth for private pulls; a heavier mental model for "it's just skills" |
| **(C) Static CDN** (e.g. R2/S3 + CDN) | Upload packages to a content-addressed key space `/<sha256>.tar.zst` + index at a stable URL | Dead-simple client (plain HTTPS GET); cacheable; cheap; trivially mirrorable for air-gap | New infra to own + pay for; cosign keyless still works but provenance ties to CI, not the CDN; we manage retention/GC |

**Recommendation: (A) GitHub Releases for alpha/beta, migrate the *package store* to (B) ghcr.io
OCI for GA**, keeping the human-readable `skills-index.json` (GitHub-hosted) as the entry point.
Rationale: (A) gets us shipping with zero infra and native cosign OIDC; OCI (B) is the right
long-term home because it is content-addressable by digest *natively* (our hashes map straight to
OCI digests), CDN-backed, and the cosign integration is best-in-class. (C) stays the documented
escape hatch for fully self-hosted/air-gapped mirrors. The client abstracts the host behind a
`--target` / config so users never care which is in use.

---

## 6. Invariant preservation

| Invariant | How Phase 5 preserves it |
|---|---|
| **Exactly one authoritative registry (`manifest.json`)** | `skills-index.json` is *generated from* `manifest.json` in the same build run. It is a derived artifact, peer to `build_manifest.json` — never a second source of truth. `--check` proves no drift. |
| **Universal binary surface** | `routing.compatibility_access[]` in `manifest.json` and `compatibility_binary_aliases[]` in `build_manifest.json` remain metadata fields, but default releases keep them empty. The index carries historical lookup aliases only as non-install metadata; pulled installs use `skill-router`. |
| **Deterministic, offline, no-remote-LLM query path** | `pull`/`add` are the *only* verbs that touch the network, and only when explicitly invoked. `route`/`search`/`load_skill` never change. Embeddings/index for routing stay build-time and hash-pinned. |
| **kebab-case canonical ids** | Package names, index `name`, and `manifest.json` `name` are the same kebab id (schema `pattern` enforces it). |
| **Telemetry local-only + opt-in** | Pull adds no telemetry. No "phone home" on install. Any pull metrics (if added later) follow the existing `~/.skill-router/telemetry.jsonl` local-only + opt-in rule. |
| **`--no-git` gitleaks posture** | Private signing keys never enter git (CI secret store only). Only *public* keys are committed, and the gitleaks gate allowlists those specific public-key files; intentional teaching-example hits stay gated behind `--no-git` as today. |
| **No `packages/` relocation** | All new code lives under `skill-router-cli/cmd/pkg/` and `skill-router-cli/cmd/index/`. The CLI directory is not moved. |

### 6.1 Replacing physical-copy `DefaultSync` with pull-on-demand

Phase 2 deprecates the physical-copy `DefaultSync` roots (the source of the `.claude=309` /
`.codex=162` partial-copy drift warnings) in favor of the MCP `serve` surface. Phase 5 supplies
the *positive* replacement for the cases where a client genuinely needs files on disk:

- **Before:** `DefaultSync` bulk-copies the entire corpus into ~30 client roots → drift, partial
  copies, stale skills, no integrity story.
- **After:** `skill-router pkg add <ids>` pulls *only* the needed skills + transitive deps,
  verifies each against the signed index, and installs into a tracked local store
  (`~/.skill-router/pkgs/<name>@<version>/`). `pkg ls`/`pkg remove` make the install set explicit
  and reversible. No bulk copy, no drift class, integrity guaranteed.

The MCP `serve` path remains primary for MCP-capable clients; pull-on-demand serves the
instruction-only / files-on-disk clients without resurrecting bulk physical copy.

---

## 7. Threat model + supply-chain integrity

### 7.1 Threats considered

| Threat | Mitigation |
|---|---|
| **Tampered package in transit / at rest** | Content hash in signed index; `pkg verify` recomputes and compares before install. A flipped byte changes the hash → rejected. |
| **Tampered / forged index** | Index signed by minisign (pinned key) + cosign (OIDC identity). `root_hash` binds all entries; editing one entry invalidates the signature. |
| **Compromised hosting (GitHub/OCI/CDN)** | Hosting is untrusted by design — integrity comes from the *signature*, not the transport. A malicious mirror cannot produce a validly signed altered index without the signing identity. |
| **Stolen minisign private key** | cosign keyless provides a second, independent trust path (CI OIDC). Rotation via signed key-transition statement. cosign has no long-lived key to steal. |
| **Path traversal / symlink escape in a package** | Canonical tar rejects symlinks/devices/hardlinks at build (§3.2.5); the installer refuses any entry path that escapes the destination (`..`, absolute paths). |
| **Downgrade / rollback attack** | The index records the current `version` per skill; the client persists the highest version+`generated_at` it has ever trusted and **refuses an index older than what it has seen** unless `--allow-downgrade` is passed explicitly. Prevents an attacker pinning a victim to a known-vulnerable older skill. |
| **Mix-and-match (valid old package, current index)** | Index entry binds `name@version` to `hash`; installing a package whose embedded `package.json` version/hash doesn't match the index entry is rejected. |

### 7.2 keyless cosign vs minisign — and is TUF warranted?

Both are used (§4.2): minisign for the offline/air-gapped audience, cosign for CI-native
provenance + no key custody. They are complementary, not redundant.

**TUF (The Update Framework): not warranted in v1 — recommended explicitly *against*.** TUF's
value (separation of root/targets/snapshot/timestamp roles, threshold signing, freeze-attack
resistance via timestamp roles) is aimed at large, multi-maintainer, continuously-published
repositories with many signing keys. Our model is a single repo with a single release pipeline
publishing a single snapshot index. A signed index + downgrade protection + a signed
key-transition statement covers the realistic threats at a fraction of the operational cost.
**Revisit TUF only if** we move to independent per-skill publishing cadences or multiple
independent signers — at which point the role separation earns its complexity. Document this as a
deliberate, revisitable decision, not an oversight.

### 7.3 Integrity-before-install (non-negotiable)

No package byte is ever written to the install destination before its hash matches the
signed-and-verified index entry. Installs are atomic (temp dir → fsync → rename) so a crash mid-
install cannot leave a half-written, unverified skill on disk.

---

## 8. Rollout

### 8.1 Phased rollout

| Stage | Scope | Gate to advance |
|---|---|---|
| **Alpha — signing only** | `index build` + `index sign` + `index verify`. Packages produced and signed in CI; published as release assets. No client pull yet. `version` schema-gated as `--warn`. | Index reproducible under `--check`; both signatures verify; reproducibility test (build-twice-equal) green. |
| **Beta — subset pull** | `pkg pull`/`add`/`ls`/`remove`; resolution + transitive deps; offline cache. Hosted on GitHub Releases (option A). `DefaultSync` marked deprecated in docs (started in Phase 2). | A client pulls + verifies a real subset end-to-end; downgrade protection tested; whole-repo install still works unchanged. |
| **GA** | `version` schema-gated **required**; package store migrated to ghcr.io OCI (option B); `pkg add` is the documented default for files-on-disk clients; whole-repo install still supported as legacy path. | 100% corpus versioned + validated; OCI publish + pull verified; migration script run over all 1,813 skills. |

### 8.2 Backwards-compat with today's whole-repo install

The current "clone the repo / DefaultSync the whole tree" model **keeps working unchanged through
GA and after.** Subset pull is strictly additive. The router reads `manifest.json` exactly as
today; `version`/`deps` are additive fields a pre-Phase-5 reader simply ignores (Go `omitempty` /
unknown-field tolerance). No flag day.

### 8.3 Migrating the 1,813-skill corpus into versioned packages

1. **Backfill** `version: 1.0.0` into every `SKILL.md` lacking one (Phase-0 backfill script; one
   commit, schema-validated).
2. **Seed `deps[]`** from the Phase-4 capability DAG where edges exist; empty otherwise.
3. **First index** (`registry_version` bumped, e.g. `2.3.0`): every skill at `1.0.0`, hashes
   computed, `root_hash` set, signed.
4. Thereafter §2.3's content-hash-gated bump enforcement keeps versions honest on every PR.

### 8.4 CI/CD release workflow shape

`.github/workflows/release.yml` (sketch — names illustrative):

```yaml
on:
  push:
    tags: ["v*"]
permissions:
  contents: write        # upload release assets
  id-token: write        # cosign keyless OIDC
jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: go build ./skill-router-cli/...

      # 1. Deterministic build — proves no drift vs committed manifest.json (Phase 0 gate)
      - run: skill-router registry build --check

      # 2. Build content-addressed packages + the index (deterministic)
      - run: skill-router pkg build   --out dist/pkgs
      - run: skill-router index build --write
      - run: skill-router pkg build   --out dist/pkgs2   # build twice...
      - run: diff -r dist/pkgs dist/pkgs2                # ...assert reproducible

      # 3. Sign the index (NOT each package): minisign (offline) + cosign (keyless OIDC)
      - run: skill-router index sign --minisign-key "$MINISIGN_KEY" --cosign-keyless
        env: { MINISIGN_KEY: ${{ secrets.MINISIGN_KEY }} }

      # 4. Self-verify before publishing — never ship an unverifiable index
      - run: skill-router index verify --index dist/skills-index.json

      # 5. Publish (alpha/beta: GitHub Releases; GA: also push OCI to ghcr.io)
      - run: skill-router pkg publish --index dist/skills-index.json --target github
```

No goreleaser is introduced (none exists today, and the deterministic Go pipeline already owns
the build); if goreleaser is later adopted it wraps this same `skill-router` pipeline rather than
replacing it.

---

## 9. Open questions, risks, effort

### 9.1 Open questions

- **zstd determinism across versions.** We hash the *uncompressed* tar to sidestep this, but the
  published `.tar.zst` filename's reproducibility still wants a pinned zstd level + single-thread.
  Confirm the pinned settings produce identical compressed bytes across the zstd versions we
  support, or accept that only the *content hash* (not the compressed blob) is guaranteed
  reproducible. (Leaning: only the content hash is a hard guarantee; compressed bytes are best-
  effort.)
- **`deps[]` authority.** Are dep edges authored in frontmatter, derived from the Phase-4 DAG, or
  both? Recommendation: derived-from-DAG is authoritative; frontmatter `deps` is an author hint
  the build reconciles against the DAG (mismatch → warn).
- **OCI vs Releases for beta.** Beta ships on Releases for simplicity; revisit whether to move to
  OCI earlier if the 1,813-asset-per-release ergonomics prove painful.
- **Changelog gating strictness at GA.** Warn-only (recommended) vs hard-fail on MAJOR-without-
  changelog.
- **cosign offline story.** Do we bundle Rekor inclusion proofs for offline cosign verification,
  or declare minisign the sole offline path (recommended: minisign is the offline path; cosign is
  the online/provenance path)?

### 9.2 Risks (HIGH overall)

- **Build-ownership migration risk.** Packaging hangs off the Go `registry build` (Phase 0). If
  byte-parity isn't airtight, packaging inherits drift. *Mitigation:* Phase 0's `--check` gate is
  a hard prerequisite; packaging adds its own build-twice reproducibility gate.
- **Reproducibility is unforgiving.** A single leaked mtime/uid/perm bit silently changes every
  hash. *Mitigation:* hand-written `tar.Header` (never `FileInfoHeader`); cross-OS CI matrix
  asserting identical hashes on Linux/macOS/Windows.
- **Key management is operational debt.** Minisign private key custody + rotation is a standing
  responsibility. *Mitigation:* cosign keyless as the no-custody second path; signed
  key-transition statements; keys live only in CI secrets (honors `--no-git`).
- **Scope creep toward TUF/per-skill publishing.** *Mitigation:* §7.2 explicitly defers TUF;
  v1 publishes a single snapshot index.
- **Client diversity.** ~30 historical client roots; pull-on-demand must serve the files-on-disk
  ones without resurrecting bulk copy. *Mitigation:* `serve` (Phase 2) remains primary;
  `pkg add` is the narrow files-on-disk replacement.

### 9.3 Rough effort

**~3–4 weeks (HIGH risk)**, gated entirely behind Phases 0–4:

| Workstream | Est. |
|---|---|
| `version`/`deps` schema + backfill + content-hash bump gate | ~3 days |
| Canonical-tar packager + reproducibility test (`pkg build`) | ~4 days |
| Index generation (`index build`) wired into `registry build` | ~2 days |
| minisign + cosign signing + verify + key mgmt (`index sign/verify`) | ~4 days |
| Resolution + subset pull + offline cache + atomic install (`pkg pull/add/ls/remove`) | ~5 days |
| Hosting (GitHub Releases → OCI) + release workflow | ~3 days |
| Threat-model hardening (downgrade protection, traversal guards) + docs + migration | ~3 days |

---

## Appendix A — File layout added by Phase 5 (additive; no relocation)

```
skill-router-cli/
  cmd/
    registry/                 # Phase 0 — extended here to also emit skills-index.json
    pkg/                      # NEW — pkg build|sign|verify|publish|pull|ls|add|remove
    index/                    # NEW — index build|sign|verify
  internal/
    pkgformat/                # NEW — canonical tar, content hashing, package.json
    signing/                  # NEW — minisign + cosign adapters, key pinning/rotation
    pull/                     # NEW — resolution, transitive deps, atomic install, offline cache
schemas/
  skill.schema.json           # Phase 0 — gains `version`, `deps`
docs/
  PHASE_5_DISTRIBUTION_DESIGN.md   # this file
  SIGNING_KEYS.md             # NEW (at impl time) — public keys + transition statements
dist/                         # build output (gitignored) — packages + index + signatures
~/.skill-router/              # client-side (runtime, gitignored)
  pkgs/<name>@<version>/      # installed pulled skills
  cache/                      # fetched index + packages for offline reuse
```

## Appendix B — Decision summary

1. **Hybrid semver:** authors declare the bump; CI content-hash-gates it (forgot-to-bump can't merge).
2. **Hash the uncompressed canonical tar** (sorted, mtime/uid/gid=0, fixed perms); zstd is transport-only.
3. **Sign the index, not 1,813 packages** — with **both** minisign (offline) and cosign (keyless OIDC).
4. **Subset pull replaces physical-copy `DefaultSync`**; integrity verified before install; route path stays offline.
5. **GitHub Releases for alpha/beta → ghcr.io OCI for GA**; TUF deferred; `manifest.json` stays the one source of truth.
