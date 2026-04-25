# Production Workflow

GDD structure, milestone frameworks, sprint planning, prioritization methods, risk management, and playtest planning.

---

## GDD (Game Design Document) Canonical Structure

### Section Template

| # | Section | Contents |
|---|---------|----------|
| 1 | **Overview** | Title, genre, platform, target audience, elevator pitch (1-2 sentences), unique selling points |
| 2 | **Core Loop** | Primary gameplay loop diagram, session flow, moment-to-moment experience |
| 3 | **Mechanics** | Detailed system descriptions, rules, interactions, formulas |
| 4 | **Progression** | Level structure, unlock sequence, difficulty curve, XP/reward pacing |
| 5 | **Narrative** | Story synopsis, characters, world-building, quest structure, dialogue approach |
| 6 | **Economy** | Currency systems, taps & sinks, monetization model, pricing |
| 7 | **Art & Audio Direction** | Visual style, mood boards, audio direction, reference games |
| 8 | **Technical Requirements** | Platform specs, engine, networking, save systems, analytics |
| 9 | **Production Plan** | Milestones, team structure, timeline, risks, success metrics |

### GDD Writing Guidelines

- **Living document**: GDD evolves throughout production; version control it.
- **Audience**: Written for the development team, not players or investors.
- **Specificity**: Numbers, formulas, and examples over vague descriptions.
- **Visual**: Include diagrams, flowcharts, and mockups where possible.
- **Modular**: Each section should be independently readable and updatable.
- **Length**: Concept GDD = 5-10 pages; Full GDD = 30-100 pages depending on scope.

### GDD Section Template (Individual)

```markdown
## [Section Name]

### Overview
[1-2 sentence summary of this system]

### Design Intent
[Why this system exists; what player experience it creates]

### Detailed Design
[Rules, formulas, interactions, diagrams]

### Edge Cases
[What happens in unusual scenarios]

### Acceptance Criteria
[Testable conditions for "done"]

### Dependencies
[Other systems this interacts with]

### Open Questions
[Unresolved design decisions]
```

---

## Milestone Framework

### Standard Game Development Milestones

| Milestone | Goal | Duration (typical) | Deliverable |
|-----------|------|---------------------|-------------|
| **Concept** | Define the game | 2-4 weeks | Concept doc, pitch deck |
| **Pre-Production** | Prove the fun | 1-3 months | Vertical slice, core loop prototype |
| **Prototype** | Validate mechanics | 1-2 months | Playable prototype, first playtest |
| **Alpha** | Feature complete | 2-6 months | All features implemented (bugs ok) |
| **Beta** | Content complete | 1-3 months | All content in; polish and balance |
| **Release Candidate** | Ship quality | 2-4 weeks | Bug fixes only; certification |
| **Gold** | Ship it | 1 week | Final build submitted |
| **Post-Launch** | Live ops | Ongoing | Patches, content updates, events |

### Milestone Gate Criteria

Each milestone transition requires:

1. **Playtest report**: At least one playtest session with findings documented
2. **Risk update**: Re-evaluated risk matrix
3. **Scope check**: Features vs timeline reality check
4. **Go/No-Go decision**: Explicit stakeholder approval

---

## Sprint Planning (Game Dev Adapted)

### Sprint Template

```markdown
## Sprint [N]: [Theme/Goal]
Duration: [2 weeks]
Sprint Goal: [One sentence describing the deliverable]

### Committed Items
| ID | Task | Owner | Estimate | Priority |
|----|------|-------|----------|----------|
| | | | | |

### Stretch Goals
| ID | Task | Owner | Estimate | Notes |
|----|------|-------|----------|-------|
| | | | | |

### Dependencies
- [External dependency and mitigation]

### Risks
- [Sprint-specific risk and plan]

### Definition of Done
- [ ] Feature implemented and testable
- [ ] No P0/P1 bugs introduced
- [ ] Art/audio placeholders identified
- [ ] Playtest session scheduled (if applicable)
```

### Game-Specific Sprint Considerations

- **Art pipeline lag**: Art tasks often take 2-3× longer than code. Start art early.
- **Playtest cadence**: Playtest at least every 2 sprints; every sprint in Beta.
- **Polish budget**: Reserve 20-30% of total timeline for polish (not features).
- **Crunch prevention**: If sprint velocity drops 2 sprints in a row, cut scope, don't add hours.

---

## Prioritization Frameworks

### MoSCoW

| Category | Definition | Budget |
|----------|-----------|--------|
| **Must Have** | Game doesn't work without it | 60% of effort |
| **Should Have** | Significantly improves experience | 20% of effort |
| **Could Have** | Nice to have if time allows | 15% of effort |
| **Won't Have** | Explicitly deferred (not forgotten) | 5% (documentation only) |

### RICE Score

```
RICE = (Reach × Impact × Confidence) / Effort
```

| Factor | Scale | Game Context |
|--------|-------|-------------|
| **Reach** | Players affected (% of DAU) | 1-100% |
| **Impact** | Experience improvement | 0.25 (minimal) to 3 (massive) |
| **Confidence** | Design certainty | 0.5 (low) to 1.0 (high) |
| **Effort** | Person-sprints | 0.5 to 10+ |

**Thresholds**: >100 = High priority, 50-100 = Medium, <50 = Low

### Kano Model

| Category | Description | Game Example |
|----------|-------------|-------------|
| **Must-Be** | Expected, absence causes dissatisfaction | Save system, basic controls |
| **Performance** | More = more satisfaction (linear) | Content volume, graphics quality |
| **Attractive** | Unexpected delight | Easter eggs, emergent gameplay |
| **Indifferent** | No impact on satisfaction | Backend optimizations players don't notice |
| **Reverse** | More = less satisfaction | Excessive tutorials, forced social features |

---

## Risk Matrix

### Risk Assessment Template

| Risk | Probability | Impact | Severity | Mitigation | Owner |
|------|------------|--------|----------|------------|-------|
| | L/M/H | L/M/H | P×I | | |

### Common Game Development Risks

#### Technical Risks
- **Performance**: Target platform can't maintain 30/60 FPS
- **Networking**: Multiplayer sync issues, cheating, server costs
- **Third-party**: Engine bugs, SDK deprecation, store policy changes
- **Scope**: Feature complexity exceeds team capability

#### Design Risks
- **Fun factor**: Core loop isn't engaging in playtest
- **Balance**: Systems interact in unintended ways
- **Onboarding**: New players can't learn the game
- **Retention**: Players leave after initial sessions

#### Schedule Risks
- **Art pipeline**: Assets take longer than estimated
- **Scope creep**: Features keep expanding
- **Dependencies**: Blocked on external deliverables
- **Team**: Key person dependency, attrition

#### Business Risks
- **Market timing**: Competitor releases similar game
- **Platform changes**: Store policies, hardware shifts
- **Monetization**: Revenue model doesn't convert
- **Regulatory**: Loot box legislation, age rating issues

---

## Playtest Plan Template

### Playtest Session Template

```markdown
## Playtest [N]: [Focus Area]
Date: [YYYY-MM-DD]
Build: [version]
Participants: [N players, recruitment criteria]
Duration: [minutes]

### Goals
1. [Primary question to answer]
2. [Secondary question]
3. [Tertiary question]

### Setup
- Device/platform: [spec]
- Pre-test survey: [Y/N]
- Think-aloud protocol: [Y/N]
- Recording: [screen/audio/video]

### Tasks
| # | Task | Success Criteria | Time Limit |
|---|------|------------------|------------|
| 1 | | | |

### Metrics to Collect
- Task completion rate
- Time to complete
- Error count
- Subjective ratings (1-5 scale)
- [Custom metrics for this test]

### Post-Test
- Debrief interview: [questions]
- Survey: [SUS/custom]
- Follow-up: [when/how]
```

### Playtest Cadence

| Phase | Frequency | Focus | Participants |
|-------|-----------|-------|-------------|
| Pre-Production | Monthly | Core loop fun | Team + friends |
| Prototype | Bi-weekly | Mechanics clarity | Target audience |
| Alpha | Weekly | Feature completeness | Mix of new + returning |
| Beta | 2-3× per week | Balance, polish, UX | Large sample, diverse |
| Pre-Launch | Daily | Bug finding | QA team + external |

---

## Design Pillars Template

Design pillars are 3-5 non-negotiable principles that guide every design decision:

```markdown
## Design Pillars

### 1. [Pillar Name] (e.g., "Meaningful Choices")
**Definition**: [One sentence]
**Implication**: [How this constrains design]
**Test**: [How to verify a feature supports this pillar]

### 2. [Pillar Name]
...

### Conflict Resolution
When pillars conflict, priority order is:
1. [Most important pillar]
2. [Second]
3. [Third]
```

**Examples of strong pillars**: "Every death teaches something", "The world reacts to your choices", "Cooperation > competition", "Easy to learn, lifetime to master".
