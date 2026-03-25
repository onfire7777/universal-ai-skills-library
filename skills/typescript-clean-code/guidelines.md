# Clean Code Guidelines

Quick reference for finding the right workflow or reference file.

**How to use**: 
1. Find your task below → use the recommended **workflow**
2. Need specific rules → load the **reference files**

---

## Workflows

Step-by-step processes for common tasks. **Start here for known tasks.**

### Code Tasks

| Task | Workflow |
|------|----------|
| Review code quality | `workflows/code-review/workflow.md` |
| Review a pull request | `workflows/pr-review/workflow.md` |
| Build new feature | `workflows/new-feature.md` |
| Fix a bug | `workflows/bug-fix.md` |
| Refactor code | `workflows/refactoring/workflow.md` |
| Write tests (TDD) | `workflows/tdd.md` |
| Plan test coverage | `workflows/test-strategy.md` |

### Professional Tasks

| Task | Workflow |
|------|----------|
| Estimate a task | `workflows/estimation.md` |
| Handle unrealistic deadline | `workflows/deadline-negotiation.md` |

---

## Part 1: Code Quality References

### By Task - Code Review

| Reviewing... | Load these files |
|--------------|------------------|
| Function quality | `references/functions/rules.md`, `references/functions/checklist.md` |
| Variable/function names | `references/naming/rules.md` |
| Class design | `references/classes/rules.md` |
| Comment quality | `references/comments/rules.md` |
| Error handling | `references/error-handling/rules.md` |
| Test quality | `references/unit-tests/rules.md` |
| Code formatting | `references/formatting/rules.md` |
| General code smells | `references/smells/rules.md` |

### By Task - Writing New Code

| Creating... | Load these files |
|-------------|------------------|
| New function | `references/functions/rules.md`, `references/naming/rules.md` |
| New class/module | `references/classes/rules.md`, `references/naming/rules.md` |
| Unit tests | `references/unit-tests/rules.md`, `references/tdd/rules.md` |
| Acceptance tests | `references/acceptance-testing/rules.md` |
| Error handling | `references/error-handling/rules.md` |
| Comments/docs | `references/comments/rules.md` |

### By Task - Refactoring

| Fixing... | Load these files |
|-----------|------------------|
| Long function | `references/functions/rules.md`, `references/functions/examples.md` |
| Large class | `references/classes/rules.md`, `references/classes/examples.md` |
| Bad names | `references/naming/rules.md`, `references/naming/examples.md` |
| Poor error handling | `references/error-handling/rules.md`, `references/error-handling/examples.md` |
| Messy tests | `references/unit-tests/rules.md`, `references/unit-tests/examples.md` |
| Any smell | `references/smells/rules.md`, `references/smells/examples.md` |

### By Code Element

| Working with... | Primary | Secondary |
|-----------------|---------|-----------|
| Functions | `references/functions/rules.md` | `references/naming/rules.md` |
| Classes | `references/classes/rules.md` | `references/functions/rules.md` |
| Modules | `references/classes/rules.md` | `references/naming/rules.md` |
| Variables | `references/naming/rules.md` | `references/comments/rules.md` |
| Tests | `references/unit-tests/rules.md` | `references/tdd/rules.md` |
| Error handling | `references/error-handling/rules.md` | `references/functions/rules.md` |
| Comments | `references/comments/rules.md` | - |

### By Problem/Symptom

| If you notice... | Load these files |
|------------------|------------------|
| Function > 20 lines | `references/functions/rules.md` |
| Too many arguments (> 3) | `references/functions/rules.md` |
| Unclear variable name | `references/naming/rules.md` |
| Class doing too much | `references/classes/rules.md` (SRP) |
| Duplicate code blocks | `references/smells/rules.md` (G5) |
| Commented-out code | `references/smells/rules.md` (C5) |
| Null checks everywhere | `references/error-handling/rules.md` |
| Hard-to-read tests | `references/unit-tests/rules.md` |
| Inconsistent formatting | `references/formatting/rules.md` |
| Feature envy | `references/smells/rules.md` (G14) |
| God class | `references/classes/rules.md`, `references/smells/rules.md` |

---

## Part 2: Professional Practices References

### By Situation - Communication

| Situation | Load these files |
|-----------|------------------|
| Asked to commit to deadline | `references/commitment/rules.md`, `references/saying-no/rules.md` |
| Need to say no | `references/saying-no/rules.md`, `references/saying-no/examples.md` |
| Making a promise | `references/commitment/rules.md`, `references/commitment/examples.md` |
| Providing estimates | `references/estimation/rules.md`, `references/estimation/examples.md` |
| Under pressure from stakeholders | `references/pressure/rules.md`, `references/saying-no/rules.md` |

### By Situation - Daily Work

| Situation | Load these files |
|-----------|------------------|
| Starting TDD workflow | `references/tdd/rules.md`, `references/tdd/examples.md` |
| Planning test strategy | `references/testing-strategies/rules.md` |
| Writing acceptance tests | `references/acceptance-testing/rules.md` |
| Struggling to code | `references/coding-practices/rules.md` |
| Too many meetings | `references/time-management/rules.md` |
| Feeling stuck | `references/coding-practices/rules.md` (writer's block) |
| Working under deadline | `references/pressure/rules.md`, `references/coding-practices/rules.md` |

### By Situation - Team & Career

| Situation | Load these files |
|-----------|------------------|
| Working with others | `references/collaboration/rules.md` |
| Pairing decisions | `references/collaboration/rules.md` |
| Career development | `references/professionalism/rules.md`, `references/practicing/rules.md` |
| Improving skills | `references/practicing/rules.md`, `references/practicing/examples.md` |
| New to team | `references/professionalism/rules.md`, `references/collaboration/rules.md` |

### By Problem - Professional

| If you're facing... | Load these files |
|---------------------|------------------|
| Unrealistic deadline | `references/saying-no/rules.md`, `references/estimation/rules.md` |
| Pressure to cut corners | `references/pressure/rules.md`, `references/professionalism/rules.md` |
| Estimate treated as commitment | `references/estimation/rules.md` |
| Too much in meetings | `references/time-management/rules.md` |
| Stuck in a mess/swamp | `references/time-management/rules.md` (messes section) |
| Code ownership conflict | `references/collaboration/rules.md` |
| Not improving skills | `references/practicing/rules.md` |

---

## Decision Tree

```
What do you need to do?
│
├─► KNOWN TASK (use workflow)
│   │
│   ├─► Code Tasks
│   │   ├─► Review code → workflows/code-review/workflow.md
│   │   ├─► Review PR → workflows/pr-review/workflow.md
│   │   ├─► Build feature → workflows/new-feature.md
│   │   ├─► Fix bug → workflows/bug-fix.md
│   │   ├─► Refactor → workflows/refactoring/workflow.md
│   │   ├─► Write tests → workflows/tdd.md
│   │   └─► Plan tests → workflows/test-strategy.md
│   │
│   └─► Professional Tasks
│       ├─► Estimate task → workflows/estimation.md
│       └─► Negotiate deadline → workflows/deadline-negotiation.md
│
├─► NEED REFERENCE (use reference files)
│   │
│   ├─► Code Quality
│   │   ├─► Names → references/naming/rules.md
│   │   ├─► Functions → references/functions/rules.md
│   │   ├─► Classes → references/classes/rules.md
│   │   ├─► Tests → references/unit-tests/rules.md
│   │   └─► Smells → references/smells/rules.md
│   │
│   └─► Professional
│       ├─► Saying no → references/saying-no/rules.md
│       ├─► Commitments → references/commitment/rules.md
│       └─► Estimation → references/estimation/rules.md
│
└─► LEARNING (use knowledge files)
    └─► Start with → references/[topic]/knowledge.md
```

---

## File Index

### Workflows (9 files)

| Workflow | Purpose |
|----------|---------|
| `workflows/code-review/workflow.md` | Review code for quality |
| `workflows/pr-review/workflow.md` | Review pull requests |
| `workflows/new-feature.md` | Build new functionality |
| `workflows/bug-fix.md` | Fix bugs with test coverage |
| `workflows/refactoring/workflow.md` | Safe refactoring process |
| `workflows/tdd.md` | Test-driven development |
| `workflows/test-strategy.md` | Plan test coverage |
| `workflows/estimation.md` | PERT-based estimation |
| `workflows/deadline-negotiation.md` | Handle deadline pressure |

### Code Quality References (25 files)

| Category | Files | Purpose |
|----------|-------|---------|
| naming | 3 | Variable, function, class naming |
| functions | 4 | Function design and review |
| classes | 3 | Class/module design |
| comments | 3 | Comment best practices |
| error-handling | 3 | Exception handling |
| unit-tests | 3 | Clean test principles |
| formatting | 3 | Code layout |
| smells | 3 | Code smell catalog |

### Professional Practice References (36 files)

| Category | Files | Purpose |
|----------|-------|---------|
| professionalism | 3 | Professional ethics |
| saying-no | 3 | Declining requests professionally |
| commitment | 3 | Making real commitments |
| coding-practices | 3 | Daily coding habits |
| tdd | 3 | Test-driven development workflow |
| practicing | 3 | Deliberate practice |
| acceptance-testing | 3 | Requirements as tests |
| testing-strategies | 3 | Test pyramid |
| time-management | 3 | Meetings, focus, priorities |
| estimation | 3 | Estimating tasks |
| pressure | 3 | Working under pressure |
| collaboration | 3 | Working with others |

---

## Common Combinations

| Scenario | Files |
|----------|-------|
| Full code review | `workflows/code-review/workflow.md` or `references/functions/checklist.md` + `references/smells/rules.md` |
| Writing a new feature | `workflows/new-feature.md` |
| Starting TDD | `workflows/tdd.md` |
| Test strategy planning | `workflows/test-strategy.md` |
| Estimate + negotiate | `workflows/estimation.md` + `workflows/deadline-negotiation.md` |
| Under pressure | `references/pressure/rules.md` + `references/coding-practices/rules.md` |
| Career growth | `references/professionalism/rules.md` + `references/practicing/rules.md` |
