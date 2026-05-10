---
name: master-skill-orchestrator
description: The master conductor for all 630+ installed skills. Activates automatically on every task to determine which skills to apply, in what order, and how to combine them. Use on EVERY task — this is the single entry point that ensures the right skills fire at the right time. Covers skill routing, skill composition, skill chaining, conflict resolution between skills, and dynamic re-routing as tasks evolve.
license: Unspecified
---
# Master Skill Orchestrator

## Purpose

You have access to 630+ specialized skills spanning software engineering, UI/UX design, security, privacy, debugging, optimization, reasoning, creativity, planning, branding, color theory, and more. This skill is the **conductor** — it ensures you never miss a relevant skill, never apply the wrong one, and always combine them optimally.

> "A tool is only as good as the craftsman's judgment about when to use it."

## Activation Protocol

This skill activates on **every task**. Before beginning any work, run the Skill Triage Protocol below.

---

## Phase 1: Skill Triage Protocol

When receiving any task, immediately classify it across these dimensions before doing anything else.

### Step 1: Task Domain Classification

Identify ALL domains this task touches. Most real tasks span multiple domains.

| Domain | Signal Words | Primary Skill Clusters |
|--------|-------------|----------------------|
| **Frontend Development** | React, UI, component, layout, CSS, Tailwind, responsive | `frontend-design`, `react-best-practices`, `shadcn-ui-official`, `composition-patterns`, `next-best-practices` |
| **UI/UX Design** | design, user experience, usability, accessibility, interaction | `ui-ux-pro-max`, `impeccable-frontend-design`, `top-design`, `ux-audit`, `interaction-design`, `wondelai-refactoring-ui` |
| **Visual Design** | color, typography, branding, logo, visual, aesthetic | `color-palette`, `graphic-design`, `logo-designer`, `brand-analyzer`, `wondelai-web-typography` |
| **Backend Development** | API, server, database, auth, microservices, architecture | `architecture-patterns`, `api-design`, `microservices-patterns`, `fastapi-templates`, `auth-patterns` |
| **Database** | SQL, query, schema, migration, PostgreSQL, MySQL | `database-design-patterns`, `sql-optimization`, `postgresql-skills`, `mysql-skills`, `simota-tuner` |
| **Security** | auth, encryption, vulnerability, OWASP, pentest, hardening | `owasp-security`, `semgrep-code-security`, `tob-*` skills, `secure-coding`, `encryption-*` |
| **Privacy** | GDPR, PII, consent, data protection, compliance | `privacy-*` skills, `gdpr-*` skills, `pii-detection`, `differential-privacy` |
| **Testing** | test, TDD, coverage, E2E, unit, integration, QA | `test-driven-development`, `tyler-*-testing`, `playwright-testing`, `cypress-testing` |
| **Debugging** | bug, error, crash, fix, investigate, root cause | `systematic-debugging`, `interactive-debugger`, `simota-scout`, `simota-specter`, `simota-rewind` |
| **Performance** | slow, optimize, speed, cache, profiling, lighthouse | `simota-bolt`, `tlc-perf-*`, `web-quality-*`, `cloudflare-web-perf` |
| **DevOps/CI** | deploy, pipeline, Docker, Kubernetes, GitHub Actions | `deployment-pipeline-design`, `github-actions`, `vercel-deploy` |
| **Planning** | plan, roadmap, sprint, requirements, scope, estimate | `strategic-planning-engine`, `pm-roadmap-planning`, `pm-prd-development`, `writing-plans` |
| **Reasoning** | think, analyze, decide, evaluate, compare, assess | `meta-thinking-engine`, `thinking-model-router`, `thinking-*` models |
| **Creativity** | brainstorm, innovate, ideate, create, design, explore | `creative-connections-engine`, `lyn-brainstorm-diverge`, `thinking-lateral` |
| **Documentation** | document, README, ADR, spec, API docs | `lyn-adr-architecture`, `openapi-spec`, `code-documenter` |

### Step 2: Task Complexity Assessment

| Complexity | Characteristics | Skill Approach |
|-----------|----------------|----------------|
| **Simple** (1 domain, clear solution) | Single file change, known pattern, quick fix | Route to 1-2 primary skills |
| **Moderate** (2-3 domains, some ambiguity) | Feature implementation, design + code, multi-file | Route to 3-5 skills in sequence |
| **Complex** (4+ domains, significant ambiguity) | Full-stack feature, architecture decision, new system | Route to 5-10 skills, use Skill Composition Protocol |
| **Wicked** (many domains, high uncertainty, novel) | Greenfield project, major refactor, system redesign | Full orchestration with all phases below |

### Step 3: Task Phase Classification

Every task moves through phases. Different skills apply at different phases.

| Phase | Purpose | Key Skills to Activate |
|-------|---------|----------------------|
| **Understand** | Clarify what's being asked, identify constraints | `meta-thinking-engine` (Layer 2: Pre-Reasoning Audit), `pm-problem-statement`, `pm-problem-framing` |
| **Plan** | Design the approach, break down work | `strategic-planning-engine`, `writing-plans`, `thinking-model-router`, `pm-epic-breakdown` |
| **Research** | Gather information, evaluate options | `lyn-environmental-foresight`, `reasoning-tools-kb`, `lyn-bayesian-calibration` |
| **Design** | Make architectural and design decisions | `architecture-patterns`, `frontend-design`, `ui-ux-pro-max`, `lyn-adr-architecture` |
| **Implement** | Write the code | Domain-specific skills (React, backend, database, etc.) |
| **Verify** | Test, review, validate | `verification-before-completion`, `test-driven-development`, `code-review-excellence`, `ux-audit` |
| **Reflect** | Learn from the process | `lyn-reviews-retros`, `meta-thinking-engine` (Layer 5: Self-Correction) |

---

## Phase 2: Skill Selection Rules

### The Priority Stack

When multiple skills could apply, use this priority ordering:

1. **Safety-critical skills first.** If the task involves user data, auth, or security, ALWAYS activate security and privacy skills before implementation skills. Never ship insecure code because you were focused on features.

2. **Design skills before implementation skills.** Think before you code. Activate `architecture-patterns`, `frontend-design`, or `ui-ux-pro-max` before `react-best-practices` or `composition-patterns`.

3. **Thinking skills before domain skills.** When facing ambiguity, activate `meta-thinking-engine` or `thinking-model-router` before jumping into domain-specific work.

4. **Quality skills after implementation skills.** After writing code, activate `code-review-excellence`, `verification-before-completion`, testing skills, and `ux-audit`.

5. **Reflection skills last.** After delivery, activate `lyn-reviews-retros` to capture learnings.

### Conflict Resolution

When two skills give contradictory guidance:

**Rule 1: Specificity wins.** A skill designed for exactly this scenario beats a general-purpose skill. `react-best-practices` beats `frontend-design` on React-specific questions.

**Rule 2: Safety wins ties.** When specificity is equal, the more conservative/secure recommendation wins.

**Rule 3: Evidence wins opinion.** A skill backed by data or research beats one based on convention.

**Rule 4: Escalate to meta-thinking.** When you can't resolve the conflict, activate `meta-thinking-engine` to reason about which skill's guidance is more applicable.

---

## Phase 3: Common Task Patterns

### Pattern: "Build a Full-Stack Feature"

Skill chain: `pm-problem-statement` → `pm-prd-development` → `architecture-patterns` → `frontend-design` + `ui-ux-pro-max` → `react-best-practices` + `shadcn-ui-official` → `api-design` + `database-design-patterns` → `owasp-security` → `test-driven-development` → `verification-before-completion`

### Pattern: "Fix a Bug"

Skill chain: `systematic-debugging` → `simota-scout` (root cause) → `simota-specter` (if concurrency) OR `simota-rewind` (if regression) → domain-specific implementation skill → `test-driven-development` (write regression test) → `verification-before-completion`

### Pattern: "Design a UI"

Skill chain: `pm-problem-statement` → `ui-ux-pro-max` → `impeccable-frontend-design` → `color-palette` → `wondelai-web-typography` → `responsive-design` → `interaction-design` → `frontend-design` → `react-best-practices` + `shadcn-ui-official` → `accessibility-testing` → `ux-audit`

### Pattern: "Make a Strategic Decision"

Skill chain: `meta-thinking-engine` → `thinking-model-router` → selected thinking model(s) → `lyn-bayesian-calibration` → `thinking-red-team` → `thinking-pre-mortem` → `lyn-reviews-retros`

### Pattern: "Plan a Project"

Skill chain: `strategic-planning-engine` → `pm-product-strategy` → `pm-roadmap-planning` → `pm-epic-breakdown` → `pm-user-story-mapping` → `lyn-project-risk-register` → `thinking-pre-mortem` → `lyn-stakeholders-org`

### Pattern: "Optimize Performance"

Skill chain: `simota-bolt` → `tlc-perf-lighthouse` → `web-quality-performance` → `cloudflare-web-perf` → `simota-tuner` (if DB) → `sql-optimization` (if DB) → domain-specific implementation → `verification-before-completion`

### Pattern: "Security Review"

Skill chain: `owasp-security` → `semgrep-code-security` → `tob-codeql-analysis` → `secure-coding` → `auth-patterns` → `encryption-*` (as needed) → `stride-threat-model` → `continuous-security-validation`

### Pattern: "Creative Brainstorm"

Skill chain: `creative-connections-engine` → `lyn-brainstorm-diverge` → `thinking-lateral` → `thinking-first-principles` → `thinking-scamper` → `lyn-synthesis-analogy` → `pm-opportunity-tree`

---

## Phase 4: Dynamic Re-Routing

Tasks evolve. What starts as "build a feature" may reveal a "fix a bug" or "redesign the architecture." Monitor for these re-routing signals:

| Signal | Re-Route To |
|--------|------------|
| "This is more complex than I thought" | Activate `strategic-planning-engine`, break down further |
| "I'm not sure what the user wants" | Activate `pm-problem-statement`, `pm-discovery-process` |
| "The existing code is a mess" | Activate `jeffallan-spec-miner`, `jeffallan-legacy-modernizer` |
| "I'm stuck" | Activate `meta-thinking-engine` (Layer 6: Meta Mode), `thinking-model-router` |
| "This has security implications" | Immediately activate `owasp-security`, `semgrep-code-security` |
| "The design doesn't feel right" | Activate `ux-audit`, `lyn-design-evaluation`, `ui-ux-pro-max` |
| "Performance is bad" | Activate `simota-bolt`, `tlc-perf-lighthouse` |
| "I keep going in circles" | Activate `meta-thinking-engine` (Self-Correction Protocol), `block-rp-why` |

---

## Phase 5: Skill Stacking Rules

Some skills are force-multipliers when combined. These combinations are greater than the sum of their parts.

| Combination | Synergy |
|------------|---------|
| `meta-thinking-engine` + `thinking-model-router` | Metacognition selects the right thinking mode, router selects the right model within that mode |
| `frontend-design` + `ui-ux-pro-max` + `color-palette` | Design principles + 99 UX rules + scientifically-generated color palettes |
| `owasp-security` + `semgrep-code-security` + `tob-codeql-analysis` | Three layers of security analysis: guidelines + automated scanning + deep static analysis |
| `strategic-planning-engine` + `pm-roadmap-planning` + `lyn-project-risk-register` | Strategy + execution plan + risk management |
| `creative-connections-engine` + `thinking-lateral` + `lyn-synthesis-analogy` | Bisociation + lateral thinking + cross-domain analogy = maximum creative output |
| `systematic-debugging` + `simota-scout` + `interactive-debugger` | Methodology + root cause analysis + live debugging tools |
| `test-driven-development` + `code-review-excellence` + `verification-before-completion` | Write tests first + review quality + final verification = zero-defect delivery |

---

## Anti-Patterns

**Skill Overload:** Don't activate 20 skills for a simple task. Match skill count to complexity.

**Skill Tunnel Vision:** Don't only use skills from one domain. A "frontend task" might need security, accessibility, and performance skills too.

**Skill Cargo Culting:** Don't follow a skill's guidance blindly when it doesn't fit the context. Use `meta-thinking-engine` to evaluate applicability.

**Skill Skipping:** Don't skip the verification and reflection phases because you're in a hurry. They're where the most value is captured.
