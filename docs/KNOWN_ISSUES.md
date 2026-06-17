# Known issues & intentional design notes

Captured during the Goal 1–3 consolidation so future cleanup does not
re-introduce removed bloat or mistake intentional structure for redundancy.

## Intentional duplication (do NOT remove as "bloat")

Some per-skill vendoring is deliberate **self-containment**, verified during the
refactor (Builder 5). Git stores identical blobs once, so the on-disk repetition
is not a storage cost:

- **`canvas-fonts`** — vendored per skill (~4 skills × 54 files). Keeps each
  canvas/design skill runnable standalone without a shared font dependency.
- **OOXML / XSD schemas** — vendored per skill (~5 skills × 39 files). Keeps
  document-generation skills self-contained.

These are intentional and should be preserved.

## Latent bugs (tracked separately, out of refactor-only scope)

- **`infrastructure/mcp-bridges/bridge_context_mode.ps1`** launches Node via the
  full path `C:\Program Files\nodejs\node.exe`. Under `mcp-proxy` this path can
  space-split. A latent fix (8.3 short path `C:\PROGRA~1\nodejs\node.exe`) existed
  in the now-removed orphan `launch_context_mode.ps1` but was **not** ported
  (this effort is refactor-only). Track as a separate bug-fix.
