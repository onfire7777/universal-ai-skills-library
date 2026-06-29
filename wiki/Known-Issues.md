# Known Issues & Intentional Design Notes

These notes were captured during the Goal 1–3 consolidation so that future cleanup doesn't re-introduce removed bloat or mistake intentional structure for redundancy.

Source: [`docs/KNOWN_ISSUES.md`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/docs/KNOWN_ISSUES.md)

---

## Intentional duplication — do NOT remove as "bloat"

Some per-skill vendoring is deliberate **self-containment**, verified during the refactor. Git stores identical blobs once, so the on-disk repetition is **not** a storage cost:

| What | Scope | Why it's intentional |
|---|---|---|
| **`canvas-fonts`** | vendored per skill (~4 skills × 54 files) | Keeps each canvas/design skill runnable standalone, without a shared font dependency. |
| **OOXML / XSD schemas** | vendored per skill (~5 skills × 39 files) | Keeps document-generation skills self-contained. |

> [!WARNING]
> Don't "deduplicate" these. They are preserved on purpose; removing them breaks the skills' standalone runnability.

---

## Resolved bugs

- **`infrastructure/mcp-bridges/bridge_context_mode.ps1`** previously launched Node via the full path `C:\Program Files\nodejs\node.exe`, which could space-split the `-Command` argument under `mcp-proxy`. **Fixed (2026-06-17):** the invocation now uses the 8.3 short path `C:\PROGRA~1\nodejs\node.exe`.

---

## Dependency vulnerabilities (Dependabot)

`skills/onefilellm/references/` pins onefilellm's upstream dependency set.

**Resolved 2026-06-17:**

| Package | Bump | Advisory |
|---|---|---|
| `yt-dlp` | `2026.3.17 → 2026.6.9` | GHSA-f7j3-774f-rfhj (cookie leak) |
| `bleach` | `6.3.0 → 6.4.0` | GHSA-8rfp-98v4-mmr6 (URI sanitization) + GHSA-g75f-g53v-794x (linkify CPU exhaustion, `==6.3.0` only) |

**Accepted risk (no upstream fix available):**

- **`nltk==3.9.4`** — GHSA-p4gq-832x-fm9v (URL-encoded path traversal in `nltk.data.load()`), severity **High**. `3.9.4` is the latest release; no patched version exists. onefilellm uses NLTK only for stopword removal via the hardcoded `stopwords` corpus name — it never passes untrusted input to `nltk.data.load()`, so the vulnerable path is **not reachable** in this usage. The Dependabot alert was dismissed as *tolerable risk*; it will be bumped as soon as upstream ships a fix.

---

## Related pages

- [Security](Security.md) — reporting and secret handling
- [Testing & CI](Testing-and-CI.md) — the gates that protect against regressions
- [Skills Corpus](Skills-Corpus.md) — why some vendoring is intentional
