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

## Resolved bugs

- **`infrastructure/mcp-bridges/bridge_context_mode.ps1`** previously launched
  Node via the full path `C:\Program Files\nodejs\node.exe`, which could
  space-split the `-Command` argument under `mcp-proxy`. **Fixed (2026-06-17):**
  invocation now uses the 8.3 short path `C:\PROGRA~1\nodejs\node.exe` (the fix
  that had existed in the since-removed orphan `launch_context_mode.ps1`).
