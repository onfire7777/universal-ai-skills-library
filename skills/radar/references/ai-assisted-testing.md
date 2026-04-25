# AI-Assisted Testing & Modern Trends (2025-2026)

Purpose: Use AI to accelerate testing without delegating judgment. Read this when Radar is asked to generate tests with AI assistance or evaluate AI-generated test suites.

Contents:

- AI usage boundaries
- self-healing limits
- AI-generated code testing strategy
- Radar integration rules

## Current State

Use these numbers as directional context, not as deployment criteria:

- AI-assisted test creation can reduce drafting time by roughly `50-70%`
- AI-generated tests often require substantial rewrite; assume `> 70%` may need revision
- human review remains mandatory even when generated tests look plausible

## AI Test Generation Rules

AI can help with:

- first-pass edge-case enumeration
- boring scaffolding for test files
- variant generation for existing assertions
- flaky-log clustering and hypothesis generation

AI cannot replace:

- meaningful assertion design
- business-priority judgment
- deciding whether a test belongs at unit, integration, or E2E level

## Review Checklist For AI-Generated Tests

- assertions are meaningful and non-empty
- tests can actually fail when behavior regresses
- edge cases include null, empty, boundary, and error conditions where relevant
- mocks reflect plausible reality
- no optimistic assumptions hide failure modes
- non-deterministic inputs are controlled

## Self-Healing Test Boundaries

| Level | What It Can Repair | Automation Level |
|------|---------------------|------------------|
| L1 | Selector drift and obvious locator changes | High |
| L2 | Small flow changes | Medium, review required |
| L3 | Business-logic drift | Low, human decision required |

Use self-healing only to recover from presentation drift, not to reinterpret intended product behavior.

## AI-Generated Code Testing Strategy

When the underlying code is AI-generated:

1. treat property-based testing as strongly preferred
2. use mutation testing when the generated tests themselves are AI-assisted
3. add contract tests for generated clients or API layers
4. apply security scanning as a non-optional companion

## Radar Integration

| Radar Phase | AI Helps With | Human Must Still Decide |
|-------------|---------------|-------------------------|
| `SCAN` | detect likely gaps and risky files | priority and business impact |
| `LOCK` | estimate complexity | what is worth testing now |
| `PING` | draft tests and variants | final assertions and scope |
| `VERIFY` | cluster flaky signals and suggest fixes | whether the fix is real and sufficient |
| `FLAKY` | pattern mining across repeated failures | root-cause confirmation |

## Default Guardrail

Treat AI-generated tests as draft material until a human has:

- reviewed the assertions
- run the tests
- confirmed the tests fail when they should
- removed optimistic or duplicated cases
