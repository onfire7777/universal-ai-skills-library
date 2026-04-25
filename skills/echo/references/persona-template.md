# Echo Persona Template

Use this template to define service-specific personas.

---

## Template

```markdown
---
name: [Persona Name]
service: [service-identifier]
type: custom  # custom | base | internal
category: user  # user | developer | designer | business | operations
created: [YYYY-MM-DD]
source: [analyzed files/documents]
---

# [Persona Name]

## Profile

| Attribute | Value |
|-----------|-------|
...
```
Physical: [Physical situation]
Temporal: [Time constraints]
Social: [Social situation]
Cognitive: [Cognitive state]
Technical: [Technical environment]
```

## JTBD (Jobs-to-be-Done)

### Functional Job
[Specific task to accomplish]

### Emotional Job
[Desired emotional state]

### Social Job
[How they want to be perceived]

## Source Analysis

Sources that informed this persona:
...
```

---

## Template Structure Overview

```
# [Persona Name]

## Profile (Required)                          # Basic information
## Internal Profile [Internal Persona]         # For development org (when type: internal)
## Demographics [Optional]                     # Demographics
## Quote (Required)                            # Symbolic statement
## Psychographics [Optional]                   # Psychographics
## Digital Behavior [Optional]                 # Digital behavior
## Literacy & Experience [Optional]            # Literacy details
## Social Context [Optional]                   # Social context
## Life Stage [Optional]                       # Life stage
## Cognitive Profile [Optional]                # CPM: Beliefs, Goals, Emotions, Values, Stance, Communication Style
## Workflow Context [Internal Persona]         # For development org (when type: internal)
## Goals / Frustrations / Behaviors (Required) # Core attributes
## Emotion Triggers (Required)                 # Emotion triggers
## Context Scenarios (Required)                # Context scenarios
...
```

---

## Field Descriptions

### Required Fields

| Field | Description | Example |
|-------|-------------|---------|
| `name` | Persona name (English recommended) | `First-Time Buyer` |
| `service` | Service identifier | `ec-platform` |
| `type` | Persona type | `custom` |
| `category` | Persona category | `user` |
| `created` | Creation date | `2026-01-31` |
| `source` | Analyzed files/documents | `[README.md, src/checkout/*]` |

### Type Values

| Value | Description |
|-------|-------------|
| `custom` | Service-specific persona (for users) |
| `base` | Echo standard persona |
| `internal` | Development organization persona |

### Category Values

| Value | Description | Target |
|-------|-------------|--------|
| `user` | Service users | End users, customers |
| `developer` | Developers | Frontend/Backend/Infra/QA engineers |
| `designer` | Designers | UI/UX designers, researchers |
| `business` | Business roles | PdM/PO/CS/Sales |
| `operations` | Operations roles | Ops/Content Editor |

### Profile Section

- **Role**: User's position (customer, admin, guest, etc.)
- **Tech Level**: Low (beginner) / Medium (general) / High (technical)
- **Device**: Primary device and percentage
- **Usage Context**: When and where they use it
- **Usage Frequency**: Frequency of service interaction

### Extended Attributes [Optional]

Six extended categories are optional. Use only the sections needed based on service characteristics.

| Category | B2B | B2C | Primary Use |
|----------|-----|-----|-------------|
| Demographics | ○ | ◎ | Font size, price sensitivity, localization |
| Psychographics | ◎ | ◎ | Copy, CTA, social proof |
| Digital Behavior | ◎ | ○ | Auto-save, progress display, sync design |
| Literacy & Experience | ◎ | ◎ | Terminology, onboarding, help format |
| Social Context | ◎ | ○ | Approval flow, sharing, permissions |
| Life Stage | ○ | ◎ | Onboarding length, upsell, pricing |
| Cognitive Profile | ◎ | ◎ | Deep persona simulation fidelity |

◎ = Recommended, ○ = Situational

### Cognitive Profile [Optional]

6-dimension deep simulation model. Use for high-fidelity persona simulations where implicit persona descriptions are insufficient. Full specification: `references/cognitive-persona-model.md`

| Sub-section | Description | Key Attributes |
|-------------|-------------|----------------|
| Beliefs | What the persona holds to be true | `technology_beliefs`, `self_efficacy`, `trust_disposition`, `mental_models`, `assumptions` |
| Goals | Hierarchical motivation (extends JTBD) | `immediate_goals`, `journey_goals`, `life_goals`, `goal_conflicts`, `goal_priority_triggers` |
| Emotions | Persona-specific emotional profile | `baseline_mood` (VAD), `emotional_reactivity`, `frustration_threshold`, `recovery_pattern`, `delight_sensitivity` |
| Values | What matters beyond the task | `value_axes`, `non_negotiables`, `willingness_to_pay`, `effort_tolerance` |
| Stance | Pre-formed opinions on UX patterns | `feature_stances`, `ux_pattern_preferences`, `risk_appetite`, `decision_style` |
| Communication Style | How they express their experience | `vocabulary_level`, `expression_style`, `complaint_pattern`, `question_style`, `reference_frame` |

### Emotion Triggers

Linked to Echo's Emotion Score:

| Score | State | Use Case |
|-------|-------|----------|
| +3 | Delighted | Experience exceeds expectations |
| +2 | Satisfied | Smooth progress |
| +1 | Relieved | Concerns resolved |
| 0 | Neutral | No particular emotion |
| -1 | Confused | Slight hesitation |
| -2 | Frustrated | Clear problem |
| -3 | Abandoned | Abandonment |

### Echo Testing Focus

Persona-specific validation priority flows. Progress can be tracked with checkboxes.

---

## Mapping to Echo Base Personas

Service-specific personas can be mapped to Echo's base personas:

| Base Persona | Example Specialized Personas |
|--------------|------------------------------|
| Newbie | First-Time Buyer, New Employee |
| Power User | Heavy Buyer, Admin User |
| Skeptic | Price-Conscious Shopper |
| Mobile User | Commuter Shopper |
| Senior | Accessibility-Focused User |
| Accessibility User | Screen Reader User |
| Privacy Paranoid | Data-Conscious User |

This mapping enables use of Echo's existing analysis frameworks (Mental Model Gap, Cognitive Load Index, etc.).

---

## Minimal vs Full Persona

### Minimal Persona (Required fields only)

For quick generation or simple services. Cognitive Profile is not needed at this level:


```markdown
# [Name]
## Profile
## Quote
## Goals / Frustrations / Behaviors
## Emotion Triggers
## Echo Testing Focus
## JTBD
## Source Analysis
```

### Full Persona (All fields)

For detailed analysis or complex B2B/B2C services. Add Cognitive Profile for deep simulation fidelity:

```markdown
# [Name]
## Profile
## Demographics
## Quote
## Psychographics
## Digital Behavior
## Literacy & Experience
## Social Context
## Life Stage
## Cognitive Profile
## Goals / Frustrations / Behaviors
## Emotion Triggers
## Context Scenarios
## JTBD
## Echo Testing Focus
## Source Analysis
```

### Internal Persona (Development organization)

For reviewing from development team member perspectives:

```markdown
# [Name]
## Profile
## Internal Profile [Internal Persona]
## Quote
## Workflow Context [Internal Persona]
## Goals / Frustrations / Behaviors
## Emotion Triggers
## Context Scenarios
## JTBD
## Echo Testing Focus
## Source Analysis
```

---

## Internal Base Personas

Base persona library for development organizations. Used with `type: internal`.

### Engineering

| Persona | Description | Primary Review Target |
|---------|-------------|----------------------|
| Frontend Developer | Frontend developer | Components, design system, dev tools |
| Backend Developer | Backend developer | API design, documentation, logs |
| Infra Engineer | Infrastructure engineer | Deploy, monitoring, ops tools |
| QA Engineer | QA engineer | Test environment, bug reports, quality metrics |
| New Engineer | Newly onboarded engineer | Onboarding, documentation |

### Design

| Persona | Description | Primary Review Target |
|---------|-------------|----------------------|
| UI Designer | UI designer | Design tool integration, components |
| UX Researcher | UX researcher | User data, analytics tools |

### Business

| Persona | Description | Primary Review Target |
|---------|-------------|----------------------|
| Product Manager | Product manager | Specs, roadmap, metrics |
| CS Representative | Customer support | Admin panel, FAQ, inquiry handling |
| Sales | Sales representative | Demo environment, materials, customer management |

### Operations

| Persona | Description | Primary Review Target |
|---------|-------------|----------------------|
| Ops Manager | Operations manager | Admin panel, monitoring, alerts |
| Content Editor | Content editor | CMS, publishing flow |

---

## User Persona vs Internal Persona

Selection guide for persona types.

### Selection Criteria

| Target | Recommended Persona Type | Primary Purpose |
|--------|-------------------------|-----------------|
| Customer-facing screens | User Persona | End user experience validation |
| Admin panel | Internal Persona | Operations usability validation |
| Dev tools/CI/CD | Internal Persona | Developer experience (DX) validation |
| Documentation/Specs | Internal Persona | Comprehensibility, new member perspective |
| Error messages/Logs | Internal Persona | Ops/debugging usefulness |
| API/SDK | Internal Persona | Developer interface validation |
| Marketing site | User Persona | Prospect experience validation |

### Combination Patterns

| Scenario | Recommended Composition |
|----------|------------------------|
| E-commerce site | User: Buyer + Internal: CS Rep |
| SaaS product | User: End user + Internal: Admin, Ops |
| Developer tool | User: External dev + Internal: Internal dev |
| B2B platform | User: Client company rep + Internal: Sales, CS |

### Internal Persona-Specific Review Perspectives

| Perspective | Description | Validation Examples |
|-------------|-------------|---------------------|
| Efficiency | Daily task operation efficiency | Bulk operations, shortcuts |
| Observability | Ease of state awareness | Dashboard, logs, alerts |
| Troubleshooting | Ease of problem resolution | Error details, search function |
| Onboarding | New member learning ease | Documentation, tooltips |
| Collaboration | Cross-team coordination ease | Sharing, handoff, comments |
