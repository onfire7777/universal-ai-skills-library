# Nexus Orchestration Patterns Reference

**Purpose:** Reference patterns for sequential, parallel, and gated execution.
**Read when:** You need to choose or explain an orchestration pattern.

## Contents
- Pattern A: Sequential Chain
- Pattern B: Parallel Branches (with Auto-Conflict Resolution)
- Pattern C: Conditional Routing
- Pattern D: Recovery Loop
- Pattern E: Escalation Path
- Pattern F: Verification Gate
- Hub Communication Protocol

Detailed patterns for agent chain execution.

---

## Pattern A: Sequential Chain

```
Nexus → NEXUS_ROUTING → Agent1 → NEXUS_HANDOFF
                           ↓
Nexus → NEXUS_ROUTING → Agent2 → NEXUS_HANDOFF
                           ↓
Nexus → NEXUS_ROUTING → Agent3 → NEXUS_HANDOFF
                           ↓
Nexus → VERIFY → DELIVER
```

**Use when**: Steps have strict dependencies (output of one is input of next)

---

## Pattern B: Parallel Branches (with Auto-Conflict Resolution)

```
Nexus → NEXUS_ROUTING (Branch A) → [Agent1 → Agent2] → NEXUS_HANDOFF
      → NEXUS_ROUTING (Branch B) → [Agent3 → Agent4] → NEXUS_HANDOFF
                                        ↓
                            ┌───────────────────────┐
                            │ CONFLICT DETECTION    │
                            │ - Identify overlaps   │
                            │ - Classify conflicts  │
                            └───────────┬───────────┘
                                        │
                        ┌───────────────┴───────────────┐
                        ▼                               ▼
                   No Conflicts                    Has Conflicts
                        │                               │
                        │                   ┌───────────┴───────────┐
                        │                   ▼                       ▼
                        │              Auto-Resolvable         Needs User
                        │              (ADJACENT,              (SEMANTIC
                        │               FORMATTING,             unclear,
                        │               SEMANTIC clear)         STRUCTURAL)
                        │                   │                       │
                        │                   ▼                       ▼
                        │              Auto-Resolve            ESCALATE
                        │                   │                       │
                        └───────────┬───────┘                       │
                                    ▼                               │
                            AGGREGATE                               │
                                    │                               │
                                    ▼                               │
                              VERIFY (tests)                        │
                                    │                               │
                        ┌───────────┴───────────┐                   │
                        ▼                       ▼                   │
                      PASS                    FAIL                  │
                        │                       │                   │
                        ▼                       ▼                   │
                    DELIVER              RECOVERY ←─────────────────┘
```

**Use when**: Independent tasks can execute simultaneously (e.g., separate features)

### Auto-Conflict Resolution Rules

| Conflict Type | Auto-Resolve? | Method |
|---------------|---------------|--------|
| ADJACENT | ✅ Always | Accept both, merge in order |
| FORMATTING | ✅ Always | Regenerate with formatter |
| SEMANTIC (owner clear) | ✅ If ownership >= 0.70 | Ownership priority |
| SEMANTIC (owner unclear) | ❌ | Escalate to user |
| STRUCTURAL | ❌ Never | Always escalate |

### Ownership Priority

When two branches modify the same code semantically:

```yaml
ownership_score:
  primary_agent_role: 0.40  # Domain specialist bonus
  change_volume: 0.30       # More changes = more ownership
  task_alignment: 0.30      # Is file central to task?

  auto_resolve_if: ownership_score >= 0.70
```

See `references/conflict-resolution.md` for detailed resolution strategies.

---

## Pattern C: Conditional Routing

```
Nexus → NEXUS_ROUTING → Agent1 → NEXUS_HANDOFF
                           ↓
Nexus → Analyze findings
           │
           ├─ [Security issue] → Sentinel → NEXUS_HANDOFF
           ├─ [Performance issue] → Bolt → NEXUS_HANDOFF
           └─ [No issues] → Continue to next step
```

**Use when**: Next agent depends on findings (e.g., Judge → Builder OR Sentinel)

---

## Pattern D: Recovery Loop

```
Nexus → NEXUS_ROUTING → Agent → NEXUS_HANDOFF
                           │
                           ├─ [SUCCESS] → Continue
                           │
                           └─ [FAILED] → Error Handler
                                    ↓
                              ┌─────────────────┐
                              │ Recovery Action │
                              │ - Retry (L1)    │
                              │ - Inject fix (L2)│
                              │ - Rollback (L3) │
                              └────────┬────────┘
                                       ↓
                              Re-execute or Escalate
```

**Use when**: Errors occur during execution (auto-recovery enabled)

---

## Pattern E: Escalation Path

```
Nexus → NEXUS_ROUTING → Agent → NEXUS_HANDOFF (Pending Confirmation)
                                        ↓
Nexus → Present to User (AskUserQuestion)
                                        ↓
User → Select option
                                        ↓
Nexus → NEXUS_ROUTING (with User Confirmation) → Agent continues
```

**Use when**: Agent encounters decision requiring user input (L4 guardrail or GUIDED mode)

---

## Pattern F: Verification Gate

```
Nexus → Chain execution complete
                   ↓
          ┌───────────────────┐
          │ VERIFICATION GATE │
          │ - Tests pass?     │
          │ - Build OK?       │
          │ - Security OK?    │
          └─────────┬─────────┘
                    │
          ┌────────┴────────┐
          ↓ PASS            ↓ FAIL
      DELIVER          RECOVERY
                           │
                    ┌──────┴──────┐
                    │ Rollback OR │
                    │ Re-execute  │
                    └─────────────┘
```

**Use when**: Critical verification before final delivery (always used in AUTORUN_FULL)

---

## Hub Communication Protocol

```
User Request
     ↓
  NEXUS (Classify & Design Chain)
     ↓
  ┌──────────────────────────────────────────────────────────────┐
  │                    NEXUS_ROUTING                             │
  │  (Context, Goal, Step, Constraints, Expected Output)         │
  └──────────────────────────────────────────────────────────────┘
     ↓
  Agent A executes
     ↓
  ┌──────────────────────────────────────────────────────────────┐
  │                    NEXUS_HANDOFF                             │
  │  (Summary, Artifacts, Risks, Suggested Next, _STEP_COMPLETE) │
  └──────────────────────────────────────────────────────────────┘
     ↓
  NEXUS (Aggregate, Route, or Verify)
     ↓
  Next Agent or DELIVER
```
