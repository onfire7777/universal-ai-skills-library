# Nexus Execution Phases Reference

**Purpose:** Phase-by-phase execution flow for AUTORUN modes.
**Read when:** You need the exact flow for `AUTORUN_FULL`, `AUTORUN`, or proactive mode.

## Contents
- Phase 0: PROACTIVE_ANALYSIS (Optional)
- AUTORUN_FULL (7 Phases)
- AUTORUN (5 Phases - Simple Tasks Only)

Detailed phase descriptions for AUTORUN modes.

---

## Simplified 3-Phase Path (cross-model)

When 7-phase state tracking is difficult, use this simplified path:

| Simplified | Maps To | Actions |
|------------|---------|---------|
| PLAN | PLAN + PREPARE + CHAIN_SELECT | Classify task, assess complexity, select chain |
| DO | EXECUTE + AGGREGATE | Execute chain steps, merge parallel results |
| CHECK | VERIFY + DELIVER | Run tests, deliver final output |

Each phase completes before the next begins. Track only: current phase (PLAN/DO/CHECK) and current step (X/Y).

---

## Phase 0: PROACTIVE_ANALYSIS (Optional)

Automatically activates when `/Nexus` is invoked by itself. Skip this phase when a normal task instruction is present.

### 0-A: Project State Scan
Collect the current state of the project:

```bash
# Git status
git status --porcelain

# Recent commits
git log --oneline -10

# Activity Log (if exists)
.agents/PROJECT.md → Activity Log section
```

### 0-B: Health Assessment
Assess project health across four indicators:

| Indicator | Checks | Rating |
|-----------|--------|--------|
| `test_health` | Test execution, coverage | 🟢/🟡/🔴 |
| `security_health` | `npm audit`, dependencies | 🟢/🟡/🔴 |
| `code_health` | Linting, type checks | 🟢/🟡/🔴 |
| `doc_health` | README freshness, JSDoc | 🟢/🟡/🔴 |

### 0-C: Recommendation Generation
Generate recommended actions with priorities:

| Priority | Conditions |
|----------|------------|
| 🔴 High | Security issues, test failures, build errors |
| 🟡 Medium | Lint warnings, coverage regression, missing documentation |
| 🟢 Low | Refactoring opportunities, optimization suggestions |

### Flow After Phase 0

```
Phase 0 Complete
    ↓
User Selection (ON_PROACTIVE_START)
    ↓
├─ Recommended action selected → Phase 1: PLAN (AUTORUN_FULL)
├─ Continue previous work → Phase 1: PLAN (AUTORUN_FULL)
└─ New task specified → Standard routing → Phase 1
```

See `references/proactive-mode.md` for detailed specifications.

---

## AUTORUN_FULL (7 Phases)

### Phase 1: PLAN
Classify and analyze the task:

**Task Classification:**
- **BUG**: Error fix, defect response, "not working", "broken"
- **INCIDENT**: Production outage, service degradation, "down", "emergency", "SEV1/2/3/4"
- **API**: API design, endpoint creation, OpenAPI spec
- **FEATURE**: New feature, "I want to...", "add..."
- **REFACTOR**: Code cleanup (behavior unchanged)
- **OPTIMIZE**: Performance improvement
- **SECURITY**: Security response, vulnerability
- **DOCS**: Documentation
- **INFRA**: Infrastructure provisioning

**Complexity Assessment:**
- **SIMPLE**: 1-2 steps to complete
- **MEDIUM**: 3-5 steps
- **COMPLEX**: 6+ steps (decompose with Sherpa)

**Analysis:**
- Identify independent tasks (parallelizable)
- Identify dependent tasks (sequential required)
- Map file ownership per branch
- Determine guardrail requirements

### Phase 2: PREPARE
Set up execution environment:

1. **Context Snapshot Creation** - Capture initial goal and acceptance criteria
2. **Rollback Point Definition** - Create git stash or branch for recovery
3. **Guardrail Configuration** - Set appropriate levels per step
4. **Parallel Branch Preparation** - Split independent tasks, assign file ownership

### Phase 3: CHAIN_SELECT
Auto-select agent chain based on classification.

For parallel execution:
```
_PARALLEL_CHAINS:
  - branch_id: A
    chain: [Agent1, Agent2]
    files: [file1.ts, file2.ts]
  - branch_id: B
    chain: [Agent3, Agent4]
    files: [file3.ts, file4.ts]
  merge_point: Radar
```

### Phase 4: EXECUTE
Execute steps with guardrail checkpoints:

**Sequential:**
1. Execute agent role for current step
2. Perform work according to SKILL.md
3. Guardrail Check at configured checkpoints
4. Record result as `_STEP_COMPLETE`
5. Verify success conditions
6. Proceed to next step OR trigger recovery

**Parallel:**
1. Launch parallel branches simultaneously
2. Each branch executes independently
3. Monitor for conflicts
4. Wait for all branches at merge point

### Phase 5: AGGREGATE
Merge parallel results:

1. Collect Branch Results - Gather outputs, check for conflicts
2. Conflict Resolution - Resolve or escalate file conflicts
3. Context Consolidation - Update L1_GLOBAL, prepare unified state

### Phase 6: VERIFY
Verify acceptance criteria:

1. Run tests (Radar equivalent)
2. Confirm build passes
3. Security scan if applicable (Sentinel)
4. Final Guardrail Check (L2_CHECKPOINT minimum)

### Phase 7: DELIVER
Finalize and present results:

1. Integrate final output
2. Generate change summary
3. Present verification steps
4. Cleanup rollback points (on success)

---

## AUTORUN (5 Phases - Simple Tasks Only)

| Phase | Description |
|-------|-------------|
| **CLASSIFY** | Same as AUTORUN_FULL Phase 1 |
| **CHAIN_SELECT** | Auto-select agent chain |
| **EXECUTE_LOOP** | Execute each agent role, record _STEP_COMPLETE |
| **VERIFY** | Run tests, confirm build |
| **DELIVER** | Integrate output, generate summary |

COMPLEX tasks are downgraded to GUIDED mode.
