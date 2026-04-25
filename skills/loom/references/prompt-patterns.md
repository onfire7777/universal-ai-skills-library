# Prompt Patterns

Purpose: Use this reference when Loom must design prompt sequences for Figma Make.

## Contents
- Prompt design principles
- `TC-EBC` framework
- Complexity tiers
- Attachment handling
- Sequencing patterns
- Failure recovery

## Prompt Design Principles

- Split complex work into staged prompts.
- State the target screen or component explicitly.
- Reference package components and token names exactly.
- Use Figma vocabulary: `Auto Layout`, `Fill container`, `Hug contents`, `Fixed`.
- Separate exact replication from stylistic inspiration.
- Prefer one intention per prompt.

## `TC-EBC` Framework

Use this structure for generation prompts:

| Field | Purpose |
|---|---|
| `Task` | What Make should generate now |
| `Context` | Relevant project, package, or screen context |
| `Examples` | Short examples or attached references |
| `Boundaries` | Constraints, prohibitions, and non-goals |
| `Checks` | What must be true in the result |

Minimal shape:

```text
Task:
Create the [screen/component].

Context:
Use the project Guidelines package and existing component names.

Examples:
Attachment A is style-only.
Attachment B is the 1:1 reference for structure.

Boundaries:
Use design-system tokens only.
Keep Auto Layout nesting to 3 levels or fewer.
Do not invent variants not listed in Guidelines.

Checks:
Name layers semantically.
Use Auto Layout throughout.
Match responsive behavior described in Guidelines.
```

## Complexity Tiers

| Tier | Screens | Strategy | Prompt count |
|---|---:|---|---:|
| `Simple` | `1-3` | single-pass with Guidelines | `1-3` |
| `Medium` | `4-7` | components first, then assembly | `5-10` |
| `Complex` | `8-15` | system -> pattern -> screen -> polish | `12-25` |
| `Large` | `15+` | ask first; split by module | `25+` |

## Attachment Handling

For every attachment, label it as one of:
- `Exact`: reproduce structure closely
- `Style-only`: borrow tone or visual direction, not structure
- `Content-only`: reuse text or information hierarchy

Never leave the attachment role implicit.

## Sequencing Patterns

### Component-first sequence
Use when screens share reusable parts.
1. Generate component set
2. Validate names, variants, and tokens
3. Assemble screen from known parts

### Screen-first sequence
Use when a one-off surface matters more than reusable parts.
1. Generate the target screen
2. Validate structure and naming
3. Normalize reusable parts afterward

### Flow sequence
Use for multi-step journeys.
1. Generate shared layout shell
2. Generate each step separately
3. Align transitions, hierarchy, and states

## Prompt Recovery Patterns

If Make output is weak:
- reduce scope
- isolate one screen
- replace prose with explicit constraints
- replace ambiguous references with labeled attachments
- split variants into a separate step

If output ignores tokens:
- restate exact token names
- remove decorative prose
- put token rules under `Boundaries`

If output ignores layout:
- specify Auto Layout direction, gap, padding, and sizing rules explicitly
- avoid mixed layout requests in one step

## Prompt Checklist

- [ ] Prompt has one clear task
- [ ] Attachment roles are explicit
- [ ] Token names are exact
- [ ] Auto Layout direction, gap, padding, and sizing are explicit
- [ ] Variant set is limited
- [ ] Constraints and non-goals are explicit
- [ ] Checks define what success looks like
