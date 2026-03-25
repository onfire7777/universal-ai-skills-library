# Acceptance Testing Knowledge

Core concepts and foundational understanding for acceptance testing as a communication and specification tool.

## Overview

Acceptance tests are formal specifications written collaboratively by stakeholders, testers, and developers to define when a requirement is truly "done." They serve primarily as communication tools that eliminate ambiguity, not just as verification mechanisms. The tests become executable requirements documents that cannot get out of sync with the application.

## Key Concepts

### Acceptance Tests as Communication

**Definition**: Tests written by collaboration of stakeholders and programmers to define when a requirement is complete.

Acceptance tests are documents first, tests second. Their primary purpose is to formally specify system behavior in unambiguous terms that all parties understand.

**Key points**:
- Eliminate communication errors between programmers and stakeholders
- Force precision in requirements discussions
- Create shared understanding of expected behavior

### The Definition of "Done"

**Definition**: Done means all code written, all tests pass, QA and stakeholders have accepted.

Professional developers have a single, unambiguous definition of done. When acceptance tests pass, the feature is complete and ready for deployment.

**Key points**:
- Not "done" vs "done-done" - just done
- Automated tests define the finish line
- Eliminates ambiguity about completion status

### Premature Precision

**Definition**: The trap of trying to specify requirements with exact detail too early in the process.

Both business and developers want precision before it's achievable. Business wants to know exactly what they'll get; developers want to know exactly what to build. This precision cannot be achieved early and wastes resources.

**Key points**:
- Requirements appear different on paper vs. running system
- Stakeholders change minds when they see working software
- Early precision becomes irrelevant as understanding evolves

### The Uncertainty Principle (Requirements)

**Definition**: Demonstrating a feature gives stakeholders new information that changes how they see the whole system.

When you show working software, stakeholders gain insight that impacts their view of requirements. The more precise early requirements are, the less relevant they become during implementation.

### Late Ambiguity

**Definition**: Unresolved disagreements or assumptions hidden in vague requirement language.

Stakeholders may "wordsmith" around disagreements rather than resolve them. Ambiguity in requirements often represents unresolved arguments or unstated assumptions.

**Key points**:
- Defer precision as long as possible, but resolve ambiguity before coding
- Context differs between stakeholders and developers
- Professional developers ensure all ambiguity is removed

## Terminology

| Term | Definition |
|------|------------|
| Acceptance Test | Automated test defining when a requirement is done |
| Happy Path | Tests describing features with business value |
| Unhappy Path | Tests for boundary conditions, exceptions, corner cases |
| Fixture | Code connecting test statements to the system under test |
| Scenario | Reusable test pattern matching statements to functions |

## How It Relates To

- **Unit Tests**: Different audience and purpose; unit tests are for programmers, acceptance tests are for business
- **Requirements Documents**: Acceptance tests ARE the requirements, but executable
- **Continuous Integration**: All acceptance tests should run on every commit
- **Definition of Done**: Acceptance tests formally define completion criteria

## Common Misconceptions

- **Myth**: Acceptance tests are extra work on top of requirements
  **Reality**: They ARE the work of specifying requirements, just in executable form

- **Myth**: Unit tests and acceptance tests are redundant
  **Reality**: They test through different pathways and serve different audiences

- **Myth**: Acceptance tests are primarily for verification
  **Reality**: Their primary purpose is specification and communication; verification is secondary

- **Myth**: Stakeholders should write all acceptance tests
  **Reality**: Collaboration is key; BAs write happy paths, QA writes edge cases, developers may help

## Quick Reference

| Concept | One-Line Summary |
|---------|-----------------|
| Acceptance Tests | Executable requirements written collaboratively |
| Definition of Done | All tests pass, QA and stakeholders accept |
| Premature Precision | Specifying details too early wastes effort |
| Uncertainty Principle | Seeing software changes stakeholder understanding |
| Late Ambiguity | Hidden disagreements in vague requirements |
| Communication Purpose | Tests eliminate ambiguity between all parties |
