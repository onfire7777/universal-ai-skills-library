---
name: awesome-design-md
description: Use when selecting, copying, adapting, or applying a curated DESIGN.md file for a project UI, brand-style reference, or design-system brief. Includes the VoltAgent awesome-design-md corpus locally.
license: MIT
source: https://github.com/VoltAgent/awesome-design-md
---

# Awesome DESIGN.md

Use this skill to choose and apply one relevant `DESIGN.md` reference from the local corpus. A `DESIGN.md` is a plain Markdown design-system document that AI design/coding agents can read directly. Do not load the full corpus by default.

## Workflow

1. Identify the requested product, brand, site, or aesthetic.
2. Check `CATALOG.txt` for exact or nearby matches.
3. Load only the selected `design-md/<slug>/DESIGN.md`.
4. If the user asks for a persistent project design brief, copy or adapt the chosen file into the target project root as `DESIGN.md`; this matches the upstream README usage flow.
5. Use the selected file as a visual-system reference covering identity, tokens, layout, components, motion, responsive behavior, and agent prompt guidance.
6. Treat upstream files as design guidance and public CSS-token references. Do not copy proprietary assets, logos, trademark claims, or imply brand affiliation.
7. Combine this with `ui-ux-pro-max` for broad UI/UX coverage and `impeccable` for final interaction and polish review.

## Corpus

Local corpus root: `design-md/`

Upstream README snapshot: `UPSTREAM_README.md`
