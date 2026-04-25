# Full Routing Matrix

**Purpose:** Canonical full task-type to chain mapping.
**Read when:** The quick-start matrix is insufficient and you need the full mapping.

Complete task type → agent chain mapping. The SKILL.md contains the top 5 most common patterns; this file contains the full matrix.

---

| Task Type | Primary Chain | Additions |
|-----------|---------------|-----------|
| BUG | Scout → Builder → Radar | +Sentinel (security), +Sherpa (complex) |
| INCIDENT | Triage → Scout → Builder | +Mend (known pattern), +Radar, +Triage (postmortem), +Flux (deep postmortem), +Matrix (failure scenarios) |
| FEATURE | Forge → Builder → Radar | +Sherpa (complex), +Muse (UI), +Artisan (frontend), +Matrix (variant exploration), +Flux (lateral thinking) |
| INVESTIGATE | Lens | +Scout (bug-related), +Canvas (viz), +Rewind (git) |
| DECISION | Magi | +Accord (biz-tech), +Flux (reframe) |
| SECURITY | Sentinel → Builder → Radar | +Probe (dynamic), +Specter (concurrency), +Breach (red-team), +Vigil (detection) |
| REFACTOR | Zen → Radar | +Atlas (architectural), +Grove (structure) |
| OPTIMIZE | Bolt/Tuner → Radar | +Schema (DB), +Flux (first-principles), +Matrix (target combos) |
| ANALYSIS | Ripple → Builder → Radar | +Canon (standards), +Sweep (cleanup) |
| API | Gateway → Builder → Radar | +Quill, +Schema |
| DEPLOY | Guardian → Launch | +Harvest (reporting) |
| MODERNIZE | Horizon → Builder → Radar | +Polyglot (i18n), +Grove (structure), +Flux (first-principles), +Matrix (migration paths) |
| DOCS | Quill | +Canvas, +Morph (convert), +Scribe (specs) |
| STRATEGY | Spark → Builder → Radar | +Growth/Compete/Voice/Pulse/Retain/Experiment, +Helm (simulation) |
| STRATEGY_SIM | Helm | +Compete (intel), +Pulse (KPI), +Magi (decision), +Scribe (docs), +Canvas (viz), +Sherpa (execution) |
| INFRA | Scaffold → Gear → Radar | +Anvil (CLI), +Pipe (GHA workflows) |
| GHA_WORKFLOW | Pipe | +Gear (maintenance), +Launch (release), +Sentinel (security) |
| PARALLEL | Rally | +Sherpa (decomposition), see Rally escalation |
| PROJECT | Titan | Full product lifecycle — Titan orchestrates 9 phases, issues chains to Nexus |
| MESSAGING | Relay → Builder → Radar | +Sentinel (security), +Scaffold (infra) |
| BOT | Relay → Builder → Radar | +Sentinel (security) |
| REALTIME | Relay → Scaffold → Builder | +Radar (tests) |
| WEBHOOK | Gateway → Relay → Builder | +Radar (tests), +Sentinel (security) |
| HOOKS | Latch | +Gear (Git hooks), +Sentinel (security) |
| SKILL_GEN | Sigil | +Lens (codebase analysis), +Grove (structure) |
| EVOLUTION | Darwin | +Architect (improvement), +Void (sunset), +Lore (knowledge), +Canvas (viz) |
| KNOWLEDGE_SYNC | Lore | +Darwin (evolution input), +Architect (design insights), +Nexus (routing feedback) |
| QUALITY | Judge → Canvas | +Zen (smells), +Radar (coverage), +Sentinel (security), +Atlas (arch), +Sweep (dead code), +Matrix (combinatorial) |
| COMPARE | Arena | +Scout (bug-fix), +Sentinel (security), +Guardian (quality gate), +Matrix (dimension analysis), +Flux (reframe) |
| UX_RESEARCH | Researcher → Echo → Palette | +Cast (persona), +Trace (session data) |
| E2E | Voyager → Lens | +Gear (CI), +Echo (persona-based), +Matrix (test matrix) |
| BROWSER | Navigator → Builder | +Scout (bug repro), +Bolt (perf), +Lens (evidence) |
| DB_DESIGN | Schema → Builder → Radar | +Tuner (optimize), +Atlas (arch review) |
| OBSERVABILITY | Beacon → Gear → Builder | +Triage (incident link), +Scaffold (capacity) |
| AI_FEATURE | Oracle → Builder → Radar | +Gateway (API), +Stream (pipeline), +Sentinel (safety) |
| PRERELEASE | Warden → Guardian → Launch | +Sentinel (security gate), +Radar (test gate) |
| REQUIREMENTS | Accord → Scribe → Sherpa | +Canvas (diagram), +Magi (decision), +Saga (narrative), +Cast (persona) |
| DESIGN_SYSTEM | Vision → Muse → Showcase → Quill | +Palette (tokens), +Artisan (impl) |
| CONTENT | Prose → Echo → Artisan | +Polyglot (i18n), +Researcher (insights) |
| DEV_EXPERIENCE | Hearth → Gear → Latch | +Anvil (CLI), +Sigil (project skills), +Hone (CLI audit) |
| LOAD_TEST | Siege → Bolt → Builder | +Beacon (SLO), +Triage (resilience), +Matrix (scenario combos) |
| DEMO | Director/Reel → Quill | +Showcase (catalog), +Growth (marketing) |
| SPRINT_RETRO | Harvest → Canvas | +Quill (publish), +Triage (incident link) |
| KNOWLEDGE | Scribe → Prism | +Quill (polish), +Morph (format convert) |
| AITUBER | Cast → Aether → Builder | +Artisan (avatar UI), +Scaffold (infra), +Beacon (monitoring) |
| REVIEW | Judge → Builder | +Zen (refactor), +Sentinel (security), +Matrix (impact dimensions), +Flux (blind-spot) |
| YAGNI | Void → Sweep/Zen | +Magi (approval), +Pulse (usage data) |
| REMEDIATE | Mend → Radar | +Beacon (SLO check), +Gear (infra config), +Triage (escalation) |
| SPEC_VERIFY | Attest | +Scribe (spec gaps), +Radar (BDD→tests), +Builder (violation fixes), +Warden (release gate) |
| LOOP_OPS | Orbit | +Builder (script changes), +Guardian (commit policy), +Radar (verification closure) |
| GAME | Quest → Forge → Builder → Radar | +Tone (audio), +Dot (pixel), +Clay (3D), +Lyric (music), +Saga (narrative), +Matrix (balance) |
| DESIGN | Frame → Artisan → Radar | +Muse (tokens), +Loom (guidelines), +Vision (direction), +Forge (prototype) |
| ARCHITECTURE | Stratum → Canvas | +Lens (analysis), +Atlas (review), +Scribe (docs), +Ripple (impact) |
| CREATIVE | Vision → Sketch → Artisan | +Clay (3D), +Dot (pixel), +Growth (marketing) |
| BUSINESS | Levy → Scribe | +Schema (data), +Builder (calc), +Canvas (flow) |
| ECOSYSTEM | Darwin → Gauge → Canvas | +Architect (design), +Realm (viz), +Void (prune), +Lore (knowledge) |
