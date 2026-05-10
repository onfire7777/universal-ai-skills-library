---
name: universal-ai-setup
description: Manage this repository as the user's universal AI setup hub for OpenSkills, source repo inventory, bootstrap scripts, skill audits, and cross-agent operating rules.
---

# Universal AI Setup

Use this skill before modifying this repository as an AI setup hub.

## Operating Model

This repo is the control plane for:

- OpenSkills project-local skills in `.agent/skills/`
- setup manifests in `ai-setup/manifests/`
- repeatable audit and inventory scripts in `ai-setup/scripts/`
- installed source repo records for OpenUI, ECC, design resources, LLM ingestion, Lightpanda, MemPalace, Context Mode, and Skill Seekers

## Required First Steps

1. Run `git status --short`.
2. Read `ai-setup/manifests/source-repos.json`.
3. Read `ai-setup/manifests/curated-skills.json`.
4. Read `docs/UNIVERSAL_AI_SETUP.md` for the current restore and audit workflow.

## Safe Change Rules

- Be additive by default.
- Do not overwrite global skill roots.
- Do not delete duplicate skills without proving they are truly redundant.
- Do not delete non-English skills until an English replacement exists or the skill is proven useless.
- Do not run `npx openskills sync -y` while other agents may be editing `AGENTS.md`.
- Keep secrets out of manifests, docs, logs, and generated reports.

## Standard Commands

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\audit-ai-skills.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\collect-ai-setup-inventory.ps1 -WriteReports
powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\bootstrap-universal-ai-setup.ps1 -DryRun
npx openskills list
npm run typecheck
npm test
```

## Skill Work

- Use `skill-router skill skill-security-auditor` before accepting or installing a skill.
- Use `skill-router skill skill-judge` when evaluating SKILL.md quality.
- Use `skill-router skill skill-seekers` when generating or packaging skills from external sources.
- Use project-local English overrides for usability fixes instead of mutating global skill roots first.
