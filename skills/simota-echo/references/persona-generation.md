# Echo Persona Generation Workflow

Workflow for automatically generating service-specific personas from code/documentation.

---

## Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    ANALYZE                                   │
│  Collect user information from code/documentation            │
└─────────────────────┬───────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────────────────────┐
│                    EXTRACT                                   │
│  Extract persona elements (roles, goals, pain points)        │
└─────────────────────┬───────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────────────────────┐
│                    GENERATE                                  │
│  Generate personas following template                        │
└─────────────────────┬───────────────────────────────────────┘
                      ↓
...
```

---

## Trigger Methods

### 1. Explicit Command

```
/Echo generate personas
/Echo generate personas for [service-name]
/Echo generate personas from [file-path]
/Echo generate internal personas          # Generate development organization personas
/Echo generate internal personas for [service-name]
```

### 2. Auto-Suggestion

Auto-suggested when starting review without defined personas:

```yaml
questions:
  - question: "No service-specific personas found. Would you like to generate them?"
    header: "Persona"
    options:
      - label: "Yes, generate personas (Recommended)"
        description: "Auto-generate from code/documentation"
      - label: "Use Echo base personas"
        description: "Continue review with standard personas"
      - label: "I'll provide personas"
        description: "Define personas manually"
    multiSelect: false
```

---

## Phase 1: ANALYZE

### Analysis Targets

| File Type | Extraction Target | Priority |
|-----------|------------------|----------|
| `README.md` | Target users, usage scenarios | High |
| `docs/**/*.md` | User guides, tutorials | High |
| `src/**/auth*` | Auth flow, user types | Medium |
| `src/**/user*` | User models, role definitions | Medium |
| `src/**/checkout*` | Purchase flow, customer types | Medium |
| `tests/**/*` | Test scenarios, use cases | Medium |
| `*.config.*` | Config options, feature flags | Low |
| `package.json` | Project description | Low |

### Analysis Targets (Internal Persona)

Analysis targets when generating development organization personas.

| File Type | Extraction Target | Priority |
|-----------|------------------|----------|
| `CODEOWNERS` | Team structure, responsibility areas | High |
| `.github/CODEOWNERS` | Code owners, team structure | High |
| `docs/CONTRIBUTING.md` | Development flow, contribution guidelines | High |
| `.editorconfig` | Development environment settings | Medium |
| `.vscode/**` | IDE settings, recommended extensions | Medium |
| `.idea/**` | IDE settings (JetBrains) | Medium |
| `docker-compose*.yml` | Local development environment | Medium |
| `.github/workflows/*` | CI/CD workflows | Medium |
| `Makefile` | Development commands, task definitions | Medium |
| `scripts/**` | Development/operations scripts | Medium |
| `docs/runbook*` | Operations runbooks | Medium |
| `docs/onboarding*` | Onboarding materials | Medium |
| `ARCHITECTURE.md` | Architecture documentation | Low |
| `ADR/**` | Architecture Decision Records | Low |

### Analysis Patterns

#### Pattern A: README/Documentation Analysis

```markdown
## Extraction Keywords

### User Types
- "for developers", "for teams", "for enterprises"
- "beginner", "advanced", "admin"
- "customer", "admin", "guest", "member"

### Usage Scenarios
- "when you need to...", "use case:"
- "for example...", "in cases like..."

### Tech Level
- "no coding required" → Low
- "API integration" → Medium-High
- "advanced configuration" → High
```

#### Pattern B: Code Structure Analysis

```markdown
## Code Analysis

### User Models
- `user.role`, `user.type`, `user.tier`
- Enum values for user types
- Permission levels

### Flow Analysis
- Route definitions → user journeys
- Form components → required inputs
- Error handlers → friction points

### Feature Flags
- Feature toggles → user segments
- A/B test configs → user variations
```

#### Pattern C: Test Scenario Analysis

```markdown
## Test Analysis

### E2E Tests
- Test descriptions → user stories
- Test steps → expected flows
- Assertions → success criteria

### User Stories in Tests
- "as a [role], I want to [action]"
- describe blocks with user context
```

#### Pattern D: Extended Attribute Analysis

Patterns for extracting extended attributes (Demographics, Psychographics, etc.).

```markdown
## Extended Attribute Analysis

### Demographics Signals
- "for students", "enterprise-grade" → age/occupation
- Pricing tiers → income level
- i18n/l10n configs → location
- "mobile-first", "PWA" → device/region characteristics

### Psychographics Signals
- "quick setup", "instant results" → time-saver
- "detailed docs", "comprehensive guide" → detail-oriented
- "proven", "trusted by", "case studies" → track record focused
- "innovative", "cutting-edge" → early adopter

### Digital Behavior Signals
...
```

#### Pattern E: Internal Persona Analysis

Extraction patterns for development organization personas.

```markdown
## Internal Persona Analysis

### Team Structure Signals
- CODEOWNERS entries → team/responsibility areas
- Team mentions in docs → organization structure
- GitHub team references → team names
- Directory ownership patterns → responsibility scope

### Development Environment Signals
- .editorconfig → editor settings
- .vscode/settings.json → VSCode recommended settings
- .vscode/extensions.json → recommended extensions
- docker-compose.yml → local development environment
- Makefile targets → development commands

...
```

---

## Phase 2: EXTRACT

### Extraction Matrix (Core Attributes)

| Element | Source | Extraction Method |
|---------|--------|-------------------|
| **User Types** | README, user models | Role/type enumeration |
| **Goals** | Documentation, test descriptions | "want to", "need to" patterns |
| **Context** | Usage examples, tutorials | Scenario descriptions |
| **Pain Points** | Error handling, FAQ | Error messages, common questions |
| **Tech Level** | Setup complexity, API docs | Required prerequisites |
| **Devices** | Responsive code, mobile tests | Device support status |

### Extraction Matrix (Extended Attributes)

| Element | Source | Extraction Method |
|---------|--------|-------------------|
| **Demographics** | Pricing, i18n, marketing copy | Explicit/implicit target audience |
| **Psychographics** | Value propositions, CTAs | Value proposition analysis |
| **Digital Behavior** | Analytics, session configs | Behavior data/settings |
| **Literacy & Experience** | Onboarding, help system | Detail level/format of explanations |
| **Social Context** | Permissions, workflows | Presence/complexity of org features |
| **Life Stage** | Pricing, trials, messaging | Purchase stage assumptions |

### Extraction Matrix (Internal Persona)

| Element | Source | Extraction Method |
|---------|--------|-------------------|
| **Job Type** | CODEOWNERS, tech stack | Infer from directory structure and tech |
| **Team** | CODEOWNERS, docs | Extract team/department names |
| **Experience** | Onboarding complexity | Infer from documentation detail level |
| **Responsibility** | Directory ownership | Identify from code ownership scope |
| **Primary Tools** | .vscode, .editorconfig | Extract from IDE config files |
| **Work Style** | Docker, dev scripts | Infer from dev environment setup |
| **Daily Tasks** | CI/CD, scripts | Extract from workflow definitions |
| **Collaboration** | PR templates, CONTRIBUTING | Extract from review/coordination flows |
| **Pain Points** | Issue templates, docs/troubleshooting | Extract from common issues |

### Extraction Output Format

```yaml
extracted:
  # Core Attributes
  user_types:
    - name: "First-Time Buyer"
      evidence: "README.md line 42: 'perfect for first-time shoppers'"
    - name: "Enterprise Admin"
      evidence: "src/models/user.ts: role enum includes 'admin'"

  goals:
    - goal: "Purchase products quickly"
      type: functional
      evidence: "docs/quick-start.md: 'complete purchase in 3 steps'"

  pain_points:
    - pain: "Complex registration process"
# ...
```

### Default Values & Inference Rules

Default values and inference rules for attributes that couldn't be extracted:

```yaml
defaults:
  demographics:
    age_group:
      b2b: "30s-40s"
      b2c: "20s-30s"
    occupation:
      b2b: "Employee"
      b2c: null  # Omit if cannot be inferred
    income_level:
      b2b: "Middle"
      b2c: null

  psychographics:
    time_vs_cost:
      saas: "Time-saver"
# ...
```

---

## Phase 3: GENERATE

### Generation Rules

1. **Minimum 3 personas**: Primary, Secondary, Edge Case
2. **Map to Echo base personas**: Always establish mapping
3. **Emotion Triggers**: Define service-specific triggers
4. **Testing Focus**: Explicitly state flows to validate

### Generation Rules (Internal Persona)

1. **Team structure based**: Identify key teams from CODEOWNERS, generate personas for each team
2. **Minimum 2 personas**: At least 1 developer + 1 operations/business
3. **Map to Internal Base Personas**: Reference Internal Base Personas in persona-template.md
4. **Workflow Context required**: Define daily tasks, collaboration, pain points
5. **Testing Focus**: Explicitly state validation flows from DX perspective

### Extended Attribute Generation Rules

1. **B2B/B2C determination**: Auto-determine from service characteristics, select appropriate sections
2. **Required vs Optional**: Core attributes required, extended attributes based on information availability
3. **Inferred marker**: Add `[inferred]` marker to inferred attributes
4. **Unknown marker**: Use `[unknown]` for attributes that cannot be filled due to insufficient information

### Persona Priority

| Priority | Type | Description |
|----------|------|-------------|
| P0 | Primary | Main target user (required) |
| P1 | Secondary | Important secondary user |
| P2 | Edge | Special cases, accessibility |

### Generation Prompt Structure

```markdown
## Generate Persona

Based on extracted data:
- User type: [extracted user type]
- Goals: [extracted goals]
- Context: [extracted context]
- Pain points: [extracted pain points]

Extended attributes (if available):
- Demographics: [extracted or inferred]
- Psychographics: [extracted or inferred]
- Digital Behavior: [extracted or inferred]
- Literacy & Experience: [extracted or inferred]
- Social Context: [extracted or inferred]
- Life Stage: [extracted or inferred]
...
```

### Generation Prompt Structure (Internal Persona)

```markdown
## Generate Internal Persona

Based on extracted data:
- Job type: [extracted job type]
- Team: [extracted team name]
- Responsibilities: [extracted responsibilities]
- Tools: [extracted development tools]

Workflow context (from CI/CD, docs):
- Daily tasks: [extracted or inferred]
- Collaboration: [extracted or inferred]
- Pain points: [extracted or inferred]

Generate an internal persona following `persona-template.md`:
1. Set `type: internal` and appropriate `category`
...
```

### Attribute Detail Level Selection

Detail level selection during generation:

```yaml
questions:
  - question: "Select the detail level for extended attributes"
    header: "Detail"
    options:
      - label: "Minimal (Recommended)"
        description: "Required attributes only. Fast generation"
      - label: "Standard"
        description: "Include main extended attributes"
      - label: "Full"
        description: "Include all extended attributes (for B2B)"
      - label: "Auto"
        description: "Auto-select based on extracted information"
    multiSelect: false
```

---

## Phase 4: SAVE

### File Naming Convention

```
.agents/personas/{service}/{persona-name}.md

Examples:
.agents/personas/ec-platform/first-time-buyer.md
.agents/personas/ec-platform/power-shopper.md
.agents/personas/admin-dashboard/it-admin.md
```

### Directory Structure

```
.agents/
└── personas/
    ├── README.md              # Usage guide
    └── {service-name}/        # Per-service folder
        ├── primary-user.md
        ├── secondary-user.md
        └── edge-case-user.md
```

### Save Confirmation

```yaml
questions:
  - question: "Would you like to save the generated personas?"
    header: "Save"
    options:
      - label: "Yes, save all (Recommended)"
        description: "Save to .agents/personas/{service}/"
      - label: "Review and edit first"
        description: "Review content before saving"
      - label: "Save selected only"
        description: "Save only some personas"
    multiSelect: false
```

---

## Service-Specific Review

How to run UX reviews using saved personas.

### Load Personas

```
/Echo review with saved personas
/Echo review [flow] as [persona-name]
```

### Review Process

```
1. LOAD - Load personas from .agents/personas/{service}/
2. SELECT - Select review target flow and persona
3. WALK - Apply persona-specific Emotion Triggers
4. SCORE - Score in service-specific context
5. REPORT - Report based on persona-specific Testing Focus
```

### Extended Attribute Usage in Review

How extended attributes are used in reviews:

| Extended Attribute | Review Usage |
|-------------------|--------------|
| Demographics | Font size, price display format, privacy expression validation |
| Psychographics | CTA copy, social proof placement, risk mitigation expression validation |
| Digital Behavior | Auto-save timing, progress display, sync UI validation |
| Literacy & Experience | Terminology difficulty, help format, shortcut exposure |
| Social Context | Approval flow, sharing features, permission display validation |
| Life Stage | Onboarding length, pricing presentation, time value proposition |

### Cross-Persona Analysis

Compare flows across multiple saved personas:

```markdown
### Cross-Persona Matrix: Checkout Flow

| Step | First-Time Buyer | Power Shopper | Enterprise Admin |
|------|-----------------|---------------|------------------|
| 1    | +1              | +2            | +1               |
| 2    | -2              | +1            | -1               |
| 3    | -3              | +2            | -2               |

**Universal Issues**: Step 3 (all personas struggle)
**Segment Issues**: Step 2 (affects First-Time Buyer)
```

### Extended Attribute Insights

Segment-specific insights from extended attributes:

```markdown
### Attribute-Based Insights

| Attribute | First-Time Buyer | Power Shopper | Enterprise Admin |
|-----------|-----------------|---------------|------------------|
| Time Constraint | Limited | Normal | Limited |
| Decision Making | Self | Self | Formal approval |
| Interruption Frequency | High | Low | High |

**Design Implications**:
- Auto-save: Essential for First-Time Buyer, Enterprise Admin
- Approval flow: Needs optimization for Enterprise Admin
- Session recovery: Design for interruptions across all personas
```

---

## Integration with Echo Workflow

### SKILL.md Reference

This feature integrates with the following sections in `SKILL.md`:

- **PERSONA LIBRARY**: Mapping to base personas
- **EMOTION SCORING**: Custom Emotion Triggers
- **ECHO'S DAILY PROCESS**: Extended persona selection

### Journal Integration

Important discoveries during persona generation are recorded in `.agents/echo.md`:

```markdown
## YYYY-MM-DD - Persona Discovery: [Service Name]

**Generated Personas:**
- [Persona 1]: [Key insight]
- [Persona 2]: [Key insight]

**Extended Attributes Discovered:**
- Demographics: [B2B/B2C characteristics]
- Key Psychographic: [Decision-making pattern]
- Digital Behavior: [Distinctive usage pattern]

**Unexpected Finding:**
[Unexpected user types discovered from code analysis, etc.]
```

---

## Question Templates

### ON_PERSONA_GENERATION

```yaml
questions:
  - question: "Which sources should be used to generate personas?"
    header: "Source"
    options:
      - label: "Auto-detect (Recommended)"
        description: "Auto-analyze README, docs, src"
      - label: "Documentation only"
        description: "Analyze documentation files only"
      - label: "Code only"
        description: "Analyze source code only"
      - label: "Specify files"
        description: "Specify files to analyze"
    multiSelect: false
```

### ON_PERSONA_COUNT

```yaml
questions:
  - question: "How many personas should be generated?"
    header: "Count"
    options:
      - label: "3 (Recommended)"
        description: "Primary, Secondary, Edge Case"
      - label: "5"
        description: "More detailed segmentation"
      - label: "Auto"
        description: "Based on discovered user types"
    multiSelect: false
```

### ON_PERSONA_REVIEW

```yaml
questions:
  - question: "Would you like to review with saved personas?"
    header: "Persona"
    options:
      - label: "Use saved personas (Recommended)"
        description: "Load from .agents/personas/"
      - label: "Use Echo base personas"
        description: "Use standard personas"
      - label: "Generate new personas"
        description: "Generate new personas"
    multiSelect: false
```

### ON_EXTENDED_ATTRIBUTES

```yaml
questions:
  - question: "How much extended attribute detail should be included?"
    header: "Extended"
    options:
      - label: "Auto (Recommended)"
        description: "Auto-select based on extracted information"
      - label: "Minimal"
        description: "Required attributes only (fast generation)"
      - label: "Standard"
        description: "Include main extended attributes"
      - label: "Full"
        description: "Include all extended attributes"
    multiSelect: false
```

### ON_B2B_B2C_SELECTION

```yaml
questions:
  - question: "What is the primary target of the service?"
    header: "Target"
    options:
      - label: "B2B"
        description: "Enterprise/team-focused (Social Context emphasis)"
      - label: "B2C"
        description: "Individual consumer-focused (Life Stage emphasis)"
      - label: "Both"
        description: "Both B2B/B2C (consider all attributes)"
      - label: "Auto-detect"
        description: "Auto-determine from code/documentation"
    multiSelect: false
```

### ON_INTERNAL_PERSONA_GENERATION

```yaml
questions:
  - question: "What type of Internal Persona should be generated?"
    header: "Internal"
    options:
      - label: "Auto-detect (Recommended)"
        description: "Auto-detect from CODEOWNERS, documentation"
      - label: "Developer focused"
        description: "Prioritize engineering personas"
      - label: "Operations focused"
        description: "Prioritize operations/business personas"
      - label: "Select specific roles"
        description: "Select job types to generate"
    multiSelect: false
```

### ON_INTERNAL_PERSONA_ROLES

```yaml
questions:
  - question: "Which job types of Internal Persona should be generated?"
    header: "Roles"
    options:
      - label: "Frontend Developer"
        description: "Frontend developer"
      - label: "Backend Developer"
        description: "Backend developer"
      - label: "Infra Engineer"
        description: "Infrastructure engineer"
      - label: "QA Engineer"
        description: "QA engineer"
    multiSelect: true
```

### ON_PERSONA_TYPE_SELECTION

```yaml
questions:
  - question: "What type of personas should be generated?"
    header: "Type"
    options:
      - label: "User Personas (Recommended)"
        description: "Personas for service users"
      - label: "Internal Personas"
        description: "Personas for development organization"
      - label: "Both"
        description: "Generate both types"
    multiSelect: false
```

### ON_INTERNAL_REVIEW_TARGET

```yaml
questions:
  - question: "What should be reviewed with Internal Persona?"
    header: "Review"
    options:
      - label: "Admin Panel"
        description: "Admin panel UX validation"
      - label: "Developer Tools"
        description: "Dev tools/CI/CD validation"
      - label: "Documentation"
        description: "Documentation/specs validation"
      - label: "Error Messages / Logs"
        description: "Error messages/logs validation"
      - label: "API / SDK"
        description: "API/SDK design validation"
    multiSelect: true
```
