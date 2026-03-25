---
name: typescript-clean-code
description: |
  Clean Code principles, professional practices, and workflows for TypeScript developers. Based on Robert C. Martin's "Clean Code" and "The Clean Coder" books.

  IMPORTANT: When this skill is active, always load and consult the reference files (rules.md, examples.md) before giving advice or writing code. Reference content takes precedence over general knowledge.

  Use this skill when:
  - Writing TypeScript/JavaScript code
  - Reviewing code or pull requests
  - Refactoring existing code
  - Following test-driven development (TDD)
  - Fixing bugs with proper test coverage
  - Planning test strategy for features
  - Estimating tasks accurately
  - Handling deadlines and commitments professionally
  - Working effectively with teams
---

# Clean Code

Principles, practices, and workflows for TypeScript developers.

## Critical: Reference-First Approach

**Always load and consult the reference files before applying any principle or making any recommendation.** The references in this skill contain curated, authoritative knowledge from Robert C. Martin's books, adapted for TypeScript. When this skill is active:

1. **Read references before responding** - For any code quality or professional practice topic, load the relevant `rules.md` and `examples.md` files from `references/` before giving advice or writing code. Do not rely on general knowledge alone.
2. **Reference content overrides internal knowledge** - If your general knowledge conflicts with what the reference files state, follow the reference files. They contain the specific rules, thresholds, and patterns this skill enforces.
3. **Cite specific rules** - When making recommendations, reference the specific rule (e.g., "per `references/functions/rules.md` Rule 1: Keep Functions Small, 2-5 lines ideal") so the user can trace the guidance back to its source.
4. **Use examples from reference files** - Prefer the bad/good code examples in `references/[topic]/examples.md` over generating your own. These examples are curated for TypeScript and demonstrate the exact patterns intended.
5. **Follow workflows step-by-step** - When executing a task (review, refactoring, TDD, etc.), load the corresponding workflow file and follow each step, loading the reference files each step points to.

**Do not skip loading references.** Even if you "know" Clean Code principles, the reference files contain specific TypeScript adaptations, thresholds, checklists, and smell catalogs that your general knowledge may not match exactly.

## Quick Start

1. **For a task**: Check `guidelines.md` → find the right workflow → load it → follow each step (loading referenced files)
2. **For reference**: Load the specific `rules.md` and `examples.md` files relevant to your work → apply them
3. **Follow the workflow**: Step-by-step process for consistent results — always load the files each step references

## Workflows

Step-by-step processes for common tasks:

| Workflow | When to Use |
|----------|-------------|
| `workflows/code-review/workflow.md` | Reviewing code for quality |
| `workflows/pr-review/workflow.md` | Reviewing pull requests |
| `workflows/tdd.md` | Test-driven development cycle |
| `workflows/refactoring/workflow.md` | Safe refactoring with tests |
| `workflows/new-feature.md` | Building new functionality |
| `workflows/bug-fix.md` | Fixing bugs properly |
| `workflows/test-strategy.md` | Planning test coverage |
| `workflows/estimation.md` | Estimating tasks (PERT) |
| `workflows/deadline-negotiation.md` | Handling unrealistic deadlines |

### Step-File Architecture (Code Review, PR Review, Refactoring)

The code review, PR review, and refactoring workflows use a **step-file architecture** for context-safe execution:

- Each workflow has a `workflow.md` entry point that describes steps and loads `steps/step-01-init.md`
- Each step is a separate file in `steps/`, loaded sequentially
- Progress is tracked via `stepsCompleted` array in the output document's YAML frontmatter
- If context is compacted mid-workflow, `step-01-init.md` detects the existing output and `step-01b-continue.md` resumes from the last completed step
- Each step loads specific reference files before analysis and cites rules in findings
- The refactoring workflow includes a loop (steps 4-7) for iterative change-test-commit cycles

## Reference Categories

### Part 1: Code Quality (Clean Code book)

| Category | Files | Purpose |
|----------|-------|---------|
| naming | 3 | Variable, function, class naming |
| functions | 4 | Function design and review |
| classes | 3 | Class/module design |
| comments | 3 | Comment best practices |
| error-handling | 3 | Exception handling |
| unit-tests | 3 | Clean test principles |
| formatting | 3 | Code layout |
| smells | 3 | Code smell catalog (50+) |

### Part 2: Professional Practices (Clean Coder book)

| Category | Files | Purpose |
|----------|-------|---------|
| professionalism | 3 | Professional ethics |
| saying-no | 3 | Declining requests |
| commitment | 3 | Making promises |
| coding-practices | 3 | Daily habits, flow, debugging |
| tdd | 3 | TDD workflow and benefits |
| practicing | 3 | Deliberate practice |
| acceptance-testing | 3 | Requirements as tests |
| testing-strategies | 3 | Test pyramid |
| time-management | 3 | Meetings, focus |
| estimation | 3 | PERT estimation |
| pressure | 3 | Working under pressure |
| collaboration | 3 | Working with teams |

## Key Principles (Summary Only — Always Load Full References)

These are abbreviated reminders. **Always load the corresponding reference files for the full rules, thresholds, and examples before applying.**

### Code Quality
1. **Readability** → `references/formatting/rules.md`, `references/naming/rules.md`
2. **Single Responsibility** → `references/classes/rules.md`, `references/functions/rules.md`
3. **Small Units** → `references/functions/rules.md` (Rule 1: 2-5 lines ideal)
4. **Meaningful Names** → `references/naming/rules.md`
5. **DRY** → `references/smells/rules.md` (G5)
6. **Clean Tests** → `references/unit-tests/rules.md`

### Professional Practices
1. **Take Responsibility** → `references/professionalism/rules.md`
2. **Say No** → `references/saying-no/rules.md`
3. **Commit Clearly** → `references/commitment/rules.md`
4. **Estimates != Commitments** → `references/estimation/rules.md`
5. **Stay Clean Under Pressure** → `references/pressure/rules.md`

## Guidelines

See `guidelines.md` for:
- Task → workflow mapping
- Situation → reference file mapping
- Decision tree for common scenarios

## Reference Loading Checklist

Before giving any code advice or writing code, verify:
- [ ] Identified which reference categories apply to the current task
- [ ] Loaded the `rules.md` for each applicable category
- [ ] Loaded `examples.md` if demonstrating patterns or reviewing code
- [ ] Loaded the relevant `workflow/*.md` if executing a multi-step task
- [ ] Will cite specific rules/files in recommendations
