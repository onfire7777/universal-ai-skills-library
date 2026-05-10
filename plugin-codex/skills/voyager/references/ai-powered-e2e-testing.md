# AI-Powered E2E Testing & Playwright MCP/Agents (2026)

Purpose: Use this file when Voyager must plan, generate, review, or heal tests with AI assistance.

Contents:
- MCP vs CLI tradeoffs and token costs
- Planner / Generator / Healer roles
- Auth blockers and safe adoption rules
- Review checklist for AI-generated tests
- Rollout phases and cost metrics

## AI Delivery Modes

| Mode | Primary use | Token cost | Tradeoff |
|------|-------------|------------|----------|
| Playwright MCP | Interactive exploration and complex flows | `~114K/test` | Highest context quality, highest cost |
| Playwright CLI | Large suites and lower-cost iteration | `~27K/test` | Cheaper, snapshot-oriented, less interactive |
| Playwright Test Agents (`v1.56+`) | Planned generation and healing loops | Variable | Powerful, but needs strong review guardrails |

Rule: use MCP for exploration-heavy tasks, CLI for scale-sensitive workflows, and agents only when the suite already has stable architecture.

## MCP Model

The MCP approach routes AI actions through browser tools backed by the accessibility tree instead of brittle DOM selectors.

```typescript
// Brittle DOM-style interaction
await page.click('div.checkout-btn-v3');

// Stable semantic target
// Role: button, Name: "Checkout"
```

Why this matters:
- Better resilience to selector churn
- Lower need for vision-only tooling
- Closer alignment with accessibility-first locator strategy

### MCP Server Options

| Server | Notes |
|--------|-------|
| Microsoft Playwright MCP | Default official option, accessibility-tree-first |
| ExecuteAutomation MCP | Alternative implementation with extra features |
| Midscene | Computer-vision-focused option with `aiAct`, `aiQuery`, `aiAssert` |

## Playwright Test Agents

### Agent Roles

| Agent | Responsibility | Input | Output |
|-------|----------------|-------|--------|
| Planner | Explore the app and draft a test plan | Natural-language request | Markdown plan in `specs/` |
| Generator | Convert the plan into test code | Markdown plan | TypeScript tests in `tests/` |
| Healer | Attempt to repair a failing test | Failing test plus error context | Updated test code |

### Typical Workflow

```text
1. Plan: npx playwright test --plan "Test the login flow"
   -> generates specs/login.md

2. Generate: npx playwright test --generate
   -> generates tests/login.spec.ts

3. Run: npx playwright test
   -> execute and review the result

4. Heal: npx playwright test --heal
   -> repair and rerun when a failure is reviewable

5. Loop: --loop
   -> repeats plan -> generate -> run -> heal automatically
```

## Primary Blocker: Authentication

Largest adoption blocker:
- AI agents cannot explore behind auth walls without prior setup.
- `storageState` or equivalent auth fixtures must already exist.
- Re-authenticating on every run can trigger rate limits or security alerts.

Required mitigations:
- Save auth state in `auth.setup.ts`.
- Centralize auth in `globalSetup`.
- Prefer dedicated test-only auth tokens or accounts.

## Review Rules for AI-Generated Tests

### Unresolved Risks

| Risk | Effect | Required control |
|------|--------|------------------|
| Test explosion | CI cost grows faster than value | Tagging and pruning |
| Weak business assertions | Tests verify rendering but not correctness | Human-defined acceptance criteria |
| Architecture drift | Duplicated helpers, mixed locator styles, weak assertions | Mandatory human review |
| Debug hallucination | Wrong root-cause explanations sound plausible | Human validation of evidence |
| Stateful workflow complexity | Permissions, onboarding, and multi-step flows are mis-modeled | Human-designed plan plus AI assistance |

### AI Test Review Checklist

```markdown
## AI-Generated E2E Test Review

### Required checks
- [ ] Do assertions verify business behavior?
- [ ] Can the test run independently?
- [ ] Do selectors use `getByRole`, `getByLabel`, or `getByTestId`?
- [ ] Does the test avoid `waitForTimeout()`?
- [ ] Does the test name describe user behavior?

### Architecture checks
- [ ] Does the test follow the existing Page Object Model?
- [ ] Does it reuse existing helpers and fixtures?
- [ ] Is the locator strategy consistent with the suite?
- [ ] Is test data created through APIs or factories?
```

## Team-Based AI Testing

Parallel AI roles can be useful when responsibilities are explicit:
- Functional agent: happy-path verification
- Accessibility agent: WCAG checks
- Performance agent: Core Web Vitals measurement
- Security-focused probing should be coordinated with specialized security tooling, not embedded casually in Voyager tests

Commercial platforms in this area include ZeroStep, Octomind, Currents, Bug0, and Momentic. Use them only when their workflow fits the project constraints better than local Playwright tooling.

## Adoption Roadmap

| Phase | Action | Risk |
|-------|--------|------|
| Phase 1: Foundations | Standardize on `getByRole` / `getByLabel`, stabilize `playwright.config.ts` | Low |
| Phase 2: MCP review flow | Explore -> review -> generate -> run, without autonomy | Medium |
| Phase 3: Planner / Generator | Generate new tests with required human review | Medium |
| Phase 4: Healer integration | Attempt controlled auto-repair for known failures | Medium-High |
| Phase 5: Full loop | Use `--loop` for end-to-end autonomous cycles | High |

Rule: do not jump to full-loop execution before the suite already has stable auth, selectors, data strategy, and review guardrails.

## Cost Metrics

Track:
- Token cost per generated or repaired test
- Monthly AI test spend
- Human review time
- Flaky rate and false-positive rate

Sources: [Currents: State of Playwright AI Ecosystem 2026](https://currents.dev/posts/state-of-playwright-ai-ecosystem-in-2026) · [Bug0: Playwright MCP Servers](https://bug0.com/blog/playwright-mcp-servers-ai-testing) · [TestLeaf: Playwright MCP Explained](https://www.testleaf.com/blog/playwright-mcp-ai-test-automation-2026/) · [Awesome Testing: Playwright CLI & Skills](https://www.awesome-testing.com/2026/03/playwright-cli-skills-and-isolated-agentic-testing) · [BrowserStack: Playwright AI Test Generator](https://www.browserstack.com/guide/playwright-ai-test-generator)
