# Scenario Design Guidelines

Principles, templates, and best practices for demo video scenario design.

Purpose: Read this when Director must design a story arc, tune pacing, select overlay timing, adapt the demo to an audience, or review scenario quality before recording.

Contents:
- `Storytelling Structure`: problem → solution → result framing
- `Scenario Templates`: full and quick planning templates
- `Operation Granularity Design`: step sizing rules for viewer comprehension
- `Wait Strategy`: explicit waits vs pacing pauses
- `Overlay Display Patterns`: overlay timing and style rules
- `Audience-Specific Adjustments`: new user, existing user, investor, developer variants
- `Time Allocation Guidelines`: duration guidance by scope
- `Scenario Anti-Patterns`: must-avoid pacing and structure mistakes
- `Test Data Realism`: realism and privacy rules
- `Scenario Review Checklist`: pre-recording quality review

---

## Storytelling Structure

Design demos as **stories**, not operation sequences. Follow the Problem > Solution > Result arc.

```
Problem (Setup)      Solution (Action)     Result (Resolution)
-----------------    -----------------     -----------------
Explain context      Address task          Achieve goal
5-10 sec             20-40 sec             5-10 sec
```

### Four-Act Application

| Phase | Purpose | Example (Login Feature) |
|-------|---------|------------------------|
| **Setup** | Set context | "Let's log in to the service" |
| **Rising** | Execute action | Form input, submission |
| **Climax** | Highlight | Loading, auth success |
| **Resolution** | Result & satisfaction | Dashboard display |

---

## Scenario Templates

### Standard Template

```markdown
## Demo Request: [Feature Name]

### Target Audience
- [ ] New users (onboarding)
- [ ] Existing users (new feature introduction)
- [ ] Investors / Stakeholders
- [ ] Sales / Marketing
- [ ] Internal documentation

### Demo Objective
What should viewers understand after watching this demo?
> [Describe in 1-2 sentences]

### Prerequisites
- Login state: [ ] Not logged in [ ] Logged in [ ] Admin
- Initial data: [Description of required data]
- Environment: [ ] Development [ ] Staging [ ] Demo-dedicated

### Story Flow

#### 1. Opening (5-10 seconds)
**Scene**: [First screen to display]
**Message**: [Context to convey to viewers]
**Overlay**: [ ] Yes [ ] No
> Overlay text: "[Text]"

#### 2. Main Action (20-40 seconds)
**Step list**:
1. [Action 1] -> [Expected result]
2. [Action 2] -> [Expected result]
3. [Action 3] -> [Expected result]

**Emphasis points**:
- [X seconds]: [What to emphasize]

#### 3. Closing (5-10 seconds)
**Scene**: [Final screen to display]
**Message**: [Impression to leave with viewers]
**Overlay**: [ ] Yes [ ] No
> Overlay text: "[Text]"

### Test Data Requirements
| Data Type | Content | Notes |
|-----------|---------|-------|
| User | demo@example.com | Display name: Demo User |

### Recording Settings
- Resolution: [ ] 1280x720 (recommended) [ ] 1920x1080 [ ] 375x667 (mobile)
- slowMo: [ ] 500ms (standard) [ ] 700ms (form-heavy) [ ] 1000ms (slow)
- Max duration: [XX] seconds
- Output formats: [ ] WebM [ ] MP4 [ ] GIF
```

### Quick Template

```markdown
## Quick Demo: [Feature Name]

**Audience**: [Who is this for?]
**Objective**: [What to convey?]
**Duration**: [XX seconds]

**Flow**:
1. [Screen] - [Action] - [Result]
2. [Screen] - [Action] - [Result]
3. [Screen] - [Action] - [Result]

**Test Data**: [Required data]
**Settings**: [Resolution] / slowMo [X]ms
```

---

## Operation Granularity Design

| Operation Type | Recommended Granularity | Reason |
|---------------|------------------------|--------|
| Button click | 1 action = 1 step | Clear separation |
| Form input | Split by field | Input content is visible |
| Page transition | Wait for completion | Recognize screen change |
| Animation | Wait until complete | Avoid incomplete states |

```typescript
// Bad: Too coarse (viewer can't follow)
await page.fill('#email', 'demo@example.com');
await page.fill('#password', 'password');
await page.click('#submit');

// Good: Appropriate granularity with locator-based waits
await page.fill('#email', 'demo@example.com');
await expect(page.locator('#email')).toHaveValue('demo@example.com');

await page.fill('#password', 'password');
await page.waitForTimeout(300); // Deliberate pacing pause

await page.click('#submit');
await expect(page.locator('#dashboard')).toBeVisible();
```

---

## Wait Strategy

### Prefer Locator-Based Waits for State Changes

Use `waitForTimeout()` **only** for deliberate pacing pauses. For all state changes, use explicit waits:

| Scene | Strategy | Example |
|-------|----------|---------|
| Element appears | `expect(locator).toBeVisible()` | After click, modal appears |
| Page navigation | `page.waitForURL()` | After login redirect |
| Network idle | `page.waitForLoadState('networkidle')` | After page load |
| Animation complete | `expect(locator).toHaveCSS()` | After transition |

### Pacing Pauses (waitForTimeout)

Use **only** for viewer comprehension — these are intentional delays, not state waits:

| Scene | Recommended Wait | Purpose |
|-------|-----------------|---------|
| After screen display | 500-1000ms | Viewer recognizes screen |
| After input | 300-500ms | Verify input content |
| Before button click | 200-300ms | Prepare for next action |
| After page transition | 1000-1500ms | Recognize new screen |
| Important result display | 1500-2000ms | Emphasize result |
| During overlay display | Based on text length | Until reading complete |

### Overlay Display Time Formula

```
Overlay display time = (character count x 100ms) + 500ms
```

Example: "Login successful" (16 chars) = 16 x 100 + 500 = 2100ms

---

## Overlay Display Patterns

### Step Explanation Overlay

```typescript
await showOverlay(page, 'Step 1: Enter email address', 2000);
```

### Highlight Overlay

```typescript
await showHighlight(page, '#submit-button', 'Click here!');
```

### Success/Error Overlay

```typescript
await showSuccessOverlay(page, 'Registration complete!');
await showErrorOverlay(page, 'An error occurred');
```

### Overlay Style Guide

| Property | Recommended Value | Reason |
|----------|------------------|--------|
| Background | rgba(0,0,0,0.8) | Readability |
| Text color | #FFFFFF | Contrast |
| Border radius | 8px | Soft impression |
| Padding | 16px 32px | Comfortable appearance |
| Font size | 18-24px | Readability |
| Position | Bottom center | Doesn't interfere with operation |

---

## Audience-Specific Adjustments

### For New Users
- Don't skip basic operations
- Avoid or explain technical terms
- Emphasize success experience

### For Existing Users
- Quick basic operations
- Focus on new features
- Emphasize differences from previous version

### For Investors/Stakeholders
- Emphasize business value
- Differentiation from competitors
- Imply scalability

### For Developers
- Include technical details
- Show API integration
- Customization points

---

## Time Allocation Guidelines

| Duration | Use Case |
|----------|----------|
| Under 30 seconds | Simple single operation |
| 30-60 seconds | Standard feature demo |
| 60-120 seconds | Complex flow |
| Over 120 seconds | Consider splitting |

---

## Scenario Anti-Patterns

| Anti-Pattern | Problem | Fix |
|-------------|---------|-----|
| Feature dump | 3+ features in 1 demo | 1 demo = 1 feature |
| Too fast | slowMo: 100ms | slowMo: 500-700ms |
| No context | Start immediately with form input | "Let's try XX" to set context |
| Incomplete ending | End on button click | Display result 1-2 seconds |
| Fake data | email: test@test.com, name: aaa | email: demo@example.com, name: Demo User |
| Timeout-only waits | All waits are `waitForTimeout` | Use `toBeVisible()` for state changes |

---

## Test Data Realism

- Fictional but realistic names and emails
- Meaningful numbers ($39.80 instead of $100)
- Use appropriate language for content
- Never use production data or real user information
- Keep data consistent across scenes

---

## Scenario Review Checklist

### Story
- [ ] Has clear starting point
- [ ] What viewer wants to achieve is clear
- [ ] Has satisfying conclusion

### Pacing
- [ ] Uses locator-based waits for state changes
- [ ] Uses `waitForTimeout` only for deliberate pauses
- [ ] No redundant waits

### Data
- [ ] Test data is realistic
- [ ] No confidential information included
- [ ] Data is consistent

### Technical
- [ ] Reproducible scenario
- [ ] No flaky elements
- [ ] All selectors are stable
