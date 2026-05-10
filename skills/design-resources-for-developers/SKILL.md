---
name: design-resources-for-developers
description: Use when selecting design and UI resources for frontend work, including UI graphics, fonts, color tools, icons, logos, favicons, stock photos, videos, sound effects, vectors, mockups, HTML/CSS templates, CSS frameworks, animation/chart libraries, UI kits, React/Vue/Angular/Svelte/React Native UI libraries, design systems, design tools, inspiration, compression tools, browser extensions, and AI graphic design tools. This is a locally integrated mirror of bradtraversy/design-resources-for-developers.
license: MIT
---

# Design Resources for Developers

Use this skill to find practical design resources before building, redesigning, or polishing frontend interfaces. Prefer it when the task needs third-party resource discovery, asset sourcing, design-tool suggestions, UI-library comparison, or inspiration links.

## Source

- Upstream repository: `https://github.com/bradtraversy/design-resources-for-developers`
- Integrated source commit: `6d9884e7260076e0eced4ac9f3be7665a91a9476`
- Local source checkout: `%USERPROFILE%\.design-resources-for-developers\design-resources-for-developers`

## Workflow

1. Start with `references/category-index.md` to choose the relevant category.
2. Use `scripts/search_design_resources.py` for targeted lookup without loading the full README.
3. Load `references/catalog.md` only when category-level detail is needed.
4. Load `references/source-readme.md` when the exact upstream wording or full list matters.
5. Check `references/contributing.md` before proposing additions back to the upstream list.

## Search Helper

Run from this skill directory:

```bash
python scripts/search_design_resources.py --query "tailwind" --limit 10 --format md
python scripts/search_design_resources.py --category "React UI Libraries" --limit 20 --format md
python scripts/search_design_resources.py --query "color contrast" --format json
```

The helper searches resource names, descriptions, URLs, and categories in `references/resources.json`.

## Use Rules

- Recommend resources by category and use case, not as a flat link dump.
- Prefer official, actively maintained, accessible, and production-friendly tools when the user is making implementation choices.
- Treat the mirrored catalogue as a starting point. If freshness matters, verify the resource live before relying on availability, pricing, license, or maintenance status.
- Do not paste the full catalogue into the answer. Cite only the relevant few resources and mention the category used.
- Keep licensing separate from availability: the upstream list focuses on free resources, but individual resources may have their own terms.
