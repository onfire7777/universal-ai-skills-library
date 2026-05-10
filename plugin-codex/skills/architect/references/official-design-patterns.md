# Official Design Patterns Reference

> Source: "The Complete Guide to Building Skills for Claude" (Anthropic, 2025)

Official design pattern reference for Architect's ENVISION / DESIGN phases.

---

## 1. Three Use Case Categories

### Category 1: Document & Asset Creation

**Used for**: Creating consistent, high-quality output — documents, presentations, apps, designs, code

**Key Techniques**:
- Embedded style guides and brand standards
- Template structures for consistent output
- Quality checklists before finalizing
- No external tools required — uses Claude's built-in capabilities

**Real example**: `frontend-design` skill, Office skills (docx, pptx, xlsx, ppt)

### Category 2: Workflow Automation

**Used for**: Multi-step processes that benefit from consistent methodology, including multi-MCP coordination

**Key Techniques**:
- Step-by-step workflow with validation gates
- Templates for common structures
- Built-in review and improvement suggestions
- Iterative refinement loops

**Real example**: `skill-creator` skill

### Category 3: MCP Enhancement

**Used for**: Workflow guidance to enhance tool access an MCP server provides

**Key Techniques**:
- Coordinates multiple MCP calls in sequence
- Embeds domain expertise
- Provides context users would otherwise need to specify
- Error handling for common MCP issues

**Real example**: `sentry-code-review` skill (Sentry)

### Approach Selection: Problem-first vs Tool-first

| Approach | When | User Mindset |
|----------|------|-------------|
| **Problem-first** | "I need to set up a project workspace" | Describes outcomes; skill handles tools |
| **Tool-first** | "I have Notion MCP connected" | Has access; skill provides expertise |

> Most skills lean one direction. Knowing which framing fits the use case helps choose the right pattern.

---

## 2. Five Official Patterns

> These are implementation patterns for Claude skill design. For inter-agent structural patterns, see Section 7.

### Pattern 1: Sequential Workflow Orchestration

**Use when**: Users need multi-step processes in a specific order.

```markdown
## Workflow: Onboard New Customer
### Step 1: Create Account
Call MCP tool: `create_customer`
Parameters: name, email, company
### Step 2: Setup Payment
Call MCP tool: `setup_payment_method`
Wait for: payment method verification
### Step 3: Create Subscription
Call MCP tool: `create_subscription`
Parameters: plan_id, customer_id (from Step 1)
### Step 4: Send Welcome Email
Call MCP tool: `send_email`
Template: welcome_email_template
```

**Key Techniques**: Explicit step ordering, dependencies between steps, validation at each stage, rollback instructions for failures

---

### Pattern 2: Multi-MCP Coordination

**Use when**: Workflows span multiple services.

**Example**: Design-to-development handoff

| Phase | Service | Actions |
|-------|---------|---------|
| 1. Design Export | Figma MCP | Export assets, generate specs, create manifest |
| 2. Asset Storage | Drive MCP | Create folder, upload assets, generate links |
| 3. Task Creation | Linear MCP | Create tasks, attach links, assign team |
| 4. Notification | Slack MCP | Post summary to #engineering with links and refs |

**Key Techniques**: Clear phase separation, data passing between MCPs, validation before next phase, centralized error handling

---

### Pattern 3: Iterative Refinement

**Use when**: Output quality improves with iteration.

```markdown
## Iterative Report Creation
### Initial Draft
1. Fetch data via MCP
2. Generate first draft report
3. Save to temporary file
### Quality Check
1. Run validation script: `scripts/check_report.py`
2. Identify issues (missing sections, formatting, data errors)
### Refinement Loop
1. Address each identified issue
2. Regenerate affected sections
3. Re-validate
4. Repeat until quality threshold met
### Finalization
1. Apply final formatting
2. Generate summary
3. Save final version
```

**Key Techniques**: Explicit quality criteria, iterative improvement, validation scripts, knowing when to stop

---

### Pattern 4: Context-Aware Tool Selection

**Use when**: Same outcome, different tools depending on context.

**Example**: Smart file storage with decision tree

| Condition | Tool |
|-----------|------|
| Large files (>10MB) | Cloud storage MCP |
| Collaborative docs | Notion/Docs MCP |
| Code files | GitHub MCP |
| Temporary files | Local storage |

**Key Techniques**: Clear decision criteria, fallback options, transparency about choices (explain to user why that tool was chosen)

---

### Pattern 5: Domain-Specific Intelligence

**Use when**: Skill adds specialized knowledge beyond tool access.

**Example**: Financial compliance in payment processing

```markdown
### Before Processing (Compliance Check)
1. Fetch transaction details via MCP
2. Apply compliance rules:
   - Check sanctions lists
   - Verify jurisdiction allowances
   - Assess risk level
3. Document compliance decision
### Processing
IF compliance passed → process transaction
ELSE → flag for review, create compliance case
### Audit Trail
- Log all compliance checks
- Record processing decisions
- Generate audit report
```

**Key Techniques**: Domain expertise embedded in logic, compliance before action, comprehensive documentation, clear governance

---

### Pattern Selection Guide

| Your Need | Primary Pattern | Secondary Pattern |
|-----------|----------------|-------------------|
| Ordered multi-step process | Sequential Workflow | — |
| Cross-service coordination | Multi-MCP Coordination | Sequential Workflow |
| Quality-sensitive output | Iterative Refinement | Sequential Workflow |
| Context-dependent routing | Context-Aware Tool Selection | Domain-Specific Intelligence |
| Regulatory/domain constraints | Domain-Specific Intelligence | Sequential Workflow |
| Complex workflow spanning services + quality | Multi-MCP Coordination | Iterative Refinement |

---

## 3. Planning Methodology

### Use Case Definition Template

```
Use Case: [Name]
Trigger: [What user says or does]
Steps:
1. [Action with tool/MCP reference]
2. [Action]
3. [Action]
Result: [Expected outcome]
```

### Technical Requirements Checklist

1. **Tools inventory**: Built-in capabilities vs MCP requirements
2. **File structure**: `SKILL.md` + `scripts/` + `references/` + `assets/`
3. **Dependencies**: Environment requirements → `compatibility` field
4. **Folder naming**: kebab-case, matching `name` field

### Planning Questions

- What does a user want to accomplish?
- What multi-step workflows does this require?
- Which tools are needed (built-in or MCP)?
- What domain knowledge or best practices should be embedded?

---

## 4. Success Criteria Framework

### Quantitative Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Trigger accuracy | 90%+ on relevant queries | 10-20 test queries, track auto-load rate |
| Workflow efficiency | X tool calls (baseline comparison) | Same task with/without skill |
| API reliability | 0 failed calls per workflow | MCP server logs, retry rates |

### Qualitative Metrics

| Metric | Assessment |
|--------|-----------|
| Self-guided completion | Users don't need to prompt Claude about next steps |
| Correction-free execution | Same request 3-5 times yields consistent quality |
| First-try accessibility | New user accomplishes task with minimal guidance |

### Baseline Comparison Template

| Dimension | Without Skill | With Skill |
|-----------|--------------|------------|
| User instructions | Provided each time | Automatic |
| Back-and-forth | ~15 messages | ~2 clarifying questions |
| Failed API calls | ~3 requiring retry | 0 |
| Token consumption | ~12,000 | ~6,000 |

> Note: These are aspirational targets — rough benchmarks rather than precise thresholds. Anthropic is actively developing more robust measurement guidance and tooling.

---

## 5. Progressive Disclosure Design Principle

### Three-Level System

| Level | Content | Loading | Design Goal |
|-------|---------|---------|-------------|
| **1st** | YAML frontmatter | Always in system prompt | Minimal: just enough for Claude to know WHEN to use |
| **2nd** | SKILL.md body | When skill judged relevant | Full instructions, core workflow |
| **3rd** | `references/`, `scripts/`, `assets/` | On-demand navigation | Detailed docs, validation, templates |

### Application in Design

- **SKILL.md**: Keep under 5,000 words. Focus on core instructions.
- **references/**: Detailed API patterns, rate limiting, pagination, error codes
- **scripts/**: Deterministic validation (code > language instructions for critical checks)
- **assets/**: Templates, fonts, icons used in output

---

## 6. Composability Principle

### Portability

- Skills work identically across **Claude.ai**, **Claude Code**, and **API**
- Create once, works across all surfaces
- Note environment dependencies in `compatibility` field

### MCP Complementarity (Kitchen Analogy)

| | MCP | Skills |
|---|-----|--------|
| **Metaphor** | Professional kitchen | Recipes |
| **Provides** | Tool access, real-time data | Workflows, best practices |
| **Answers** | What Claude **can** do | How Claude **should** do it |

### Multi-Skill Coexistence

- Skills should work alongside others
- Don't assume exclusive capability
- Design for composition, not isolation

### Distribution Considerations

| Surface | Capability |
|---------|-----------|
| Claude.ai | Upload via Settings > Capabilities > Skills |
| Claude Code | Place in skills directory |
| API | `/v1/skills` endpoint, `container.skills` parameter |
| Organization | Admin workspace-wide deployment |

---

## 7. Agentic Composable Patterns (Anthropic 2025-2026)

> Source: Anthropic "Building effective agents" (2025) + subsequent updates (2026)

The 5 patterns in Section 2 are intra-skill implementation patterns. This section covers inter-agent and inter-workflow structural patterns.

### Design Philosophy: Simplicity First

The key to successful agent systems is to start with the simplest possible configuration and only add complexity when there is measurable improvement. Rather than introducing frameworks or orchestration layers upfront, build up from composable primitive patterns as needed.

### Pattern A: Prompt Chaining

**Structure**: Fixed-step sequential pipeline. Each step's output becomes the next step's input.

```
Step 1 → Gate ✓ → Step 2 → Gate ✓ → Step 3 → Output
                ✗ → Abort/Retry
```

**When to use**: Tasks naturally decompose into independent subtasks, and you want to gate quality at each step.
**Skill design implication**: Superset of Sequential Workflow Orchestration (Pattern 1). Explicitly insert gates (quality checks, approvals, validations) between steps.

### Pattern B: Routing

**Structure**: Classify input and dispatch to specialized handlers.

```
Input → Classifier → Handler A (domain 1)
                   → Handler B (domain 2)
                   → Handler C (domain 3)
```

**When to use**: Different input types require different optimal processing paths. Each path can be optimized independently.
**Skill design implication**: Generalization of Context-Aware Tool Selection (Pattern 4). Nexus routing logic is a canonical example. Use `Output Routing` tables to make branching conditions explicit.

### Pattern C: Parallelization

**Structure**: Execute tasks concurrently and aggregate results. Two variants:

| Variant | Description | Example |
|---------|-------------|---------|
| **Sectioning** | Split task into independent parts for parallel processing | Concurrent review of multiple files |
| **Voting** | Execute the same task multiple times and aggregate results | Multi-perspective evaluation, consensus judgment |

**When to use**: Subtasks are independent (Sectioning), or reliability/diversity is needed (Voting).
**Skill design implication**: Arena (Voting) and Nexus chains (Sectioning) are existing implementations. Specify parallelizable partners in `COLLABORATION_PATTERNS`.

### Pattern D: Orchestrator-Worker

**Structure**: Central orchestrator dynamically decomposes tasks and delegates to workers.

```
Orchestrator → [dynamic task decomposition]
  → Worker 1 (Task A)
  → Worker 2 (Task B)
  → Worker N (Task N)
Orchestrator ← [result synthesis]
```

**When to use**: Task decomposition cannot be predicted in advance and requires runtime judgment.
**Skill design implication**: Nexus Hub-and-Spoke model and Titan's 9-phase lifecycle exemplify this pattern. Ensure hub compatibility via the `Nexus Compatibility` section.

### Pattern E: Evaluator-Optimizer

**Structure**: Generate → Evaluate → Improve feedback loop.

```
Generator → Output → Evaluator → Feedback
    ↑                                 │
    └─────────── Refine ──────────────┘
```

**When to use**: Clear evaluation criteria exist, and iterative refinement improves quality.
**Skill design implication**: Structured version of Iterative Refinement (Pattern 3). Judge→Builder feedback loop is a canonical example. Extract evaluation criteria into `references/` and separate generation from evaluation responsibilities.

### Pattern F: Autonomous Agent

**Structure**: Autonomously loop tool usage based on environment feedback.

```
while (!done) {
  Observe → Plan → Act (tool use) → Evaluate result
}
```

**When to use**: Exploratory tasks, multi-stage problem solving, or tasks requiring environment interaction.
**Skill design implication**: Scout (bug investigation) and Navigator (browser automation) are existing implementations. Always design loop termination conditions and guardrails (max iterations, timeouts, safety constraints).

### Existing ↔ Agentic Pattern Mapping

| Agentic Pattern | Closest Existing Pattern (Section 2) | Relationship | Ecosystem Example |
|----------------|--------------------------------------|--------------|-------------------|
| A: Prompt Chaining | P1: Sequential Workflow | P1 is a concrete impl of A | Orbit (runner scripts) |
| B: Routing | P4: Context-Aware Tool Selection | P4 is B specialized for tool selection | Nexus (Output Routing) |
| C: Parallelization | — (new) | No existing pattern | Arena (Voting), Rally (Sectioning) |
| D: Orchestrator-Worker | P2: Multi-MCP Coordination | P2 is D specialized for MCP | Nexus, Titan |
| E: Evaluator-Optimizer | P3: Iterative Refinement | P3 is a self-contained version of E | Judge↔Builder loop |
| F: Autonomous Agent | P5: Domain-Specific Intelligence | P5 adds domain knowledge to F | Scout, Navigator |

---

## 8. Simplicity-First Design Principle

### Decision Ladder

A framework for progressively escalating agent system complexity. Do not advance to a higher level if the lower level can solve the problem.

| Level | Composition | When to Use |
|-------|-------------|-------------|
| **L0** | Single prompt + retrieval | One-shot Q&A, knowledge lookup |
| **L1** | Single prompt + tools | Single task with external integration |
| **L2** | Prompt Chaining / Routing | Fixed-flow multi-stage processing |
| **L3** | Orchestrator-Worker / Evaluator-Optimizer | Dynamic task decomposition or iterative refinement needed |
| **L4** | Multi-Agent Autonomous | Multiple autonomous agents cooperating |

### Application to Skill Design

- **ENVISION phase**: Map requirements to the Decision Ladder and identify the minimum level.
- **DESIGN phase**: Apply the Agentic Pattern corresponding to the chosen level. Document specific reasons before escalating to a higher level.
- **VALIDATE phase**: Verify the final design uses the minimum necessary complexity.

---

## 9. Interoperability Awareness

### MCP (Model Context Protocol)

Industry-standard protocol connecting AI models to external services and data sources.

- When a skill depends on MCP servers, declare it explicitly in the `compatibility` field.
- Skills in the MCP Enhancement category (Section 1) assume MCP server availability.
- Recommend **graceful degradation** so basic functionality works without MCP.

### A2A Protocol (Agent-to-Agent)

Agent-to-agent communication protocol proposed by Google. Defines capability advertisement via Agent Cards and task-based asynchronous communication.

- No direct impact on skill design at present, but recognized as a future consideration.
- The ecosystem's `CAPABILITIES_SUMMARY` is conceptually equivalent to Agent Card capability advertisement.
- The Nexus Hub-and-Spoke model is structurally similar to A2A's task delegation pattern.

### Agent Harness Pattern (Anthropic 2026)

State management pattern for long-running agents. Structures interruption, resumption, and state recovery.

- Limited direct applicability since the current ecosystem uses a 1-session = 1-task model.
- However, Orbit's `state files` and Rally's `multi-session` design are partial implementations of this pattern.
- Recorded as a reference pattern for future cross-session design.
