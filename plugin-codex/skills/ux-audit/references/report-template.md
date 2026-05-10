# UX Audit Report Template

Use this template when writing audit reports to `docs/`.

## UX Walkthrough Report

```markdown
# UX Audit: [App Name]

**Date**: YYYY-MM-DD
**Scenario**: "[realistic task the user walked through]"
**Persona**: "[who the user is — role, context, tech comfort, time pressure]"
**Browser**: Chrome MCP / Playwright
**Viewport**: Desktop (1280px) + Mobile (375px)

## Summary

[2-3 sentences: overall impression, biggest concerns, what works well]

## Flow Score

| Metric | Value |
|--------|-------|
| Task completed? | Yes / No / Partially |
| Click count | [number] |
| Decision points (had to stop and think) | [number] |
| Dead ends (wrong path, had to backtrack) | [number] |
| Would come back? | Yes / Maybe / No |
| Could teach someone in 2 minutes? | Yes / No |
| One thing to make it twice as easy | [description] |

## Findings

### Critical (blocks user task)

- **[Short title]**: [What happened, what was expected]
  - *Where*: [page/component]
  - *Screenshot*: [filename if captured]
  - *Fix*: [concrete recommendation]

### High (causes confusion or friction)

- **[Short title]**: [description]
  - *Where*: [page/component]
  - *Fix*: [recommendation]

### Medium (suboptimal but workable)

- **[Short title]**: [description]
  - *Fix*: [recommendation]

### Low (polish)

- **[Short title]**: [description]

## Workflow Comprehension

How well the app supports the user's real work:

### Threads Tested
- **[Thread name]**: [e.g. "Client calls about renewal → find policy → send renewal email"]
  - Completable end-to-end? [Yes/No/Partially — detail]
  - Steps where the user would get lost: [list]
  - Missing connections between pages: [list]

### Wayfinding
- Can the user always tell where they are? [Yes/No — specifics]
- Do nav labels match user language? [Yes/No — mismatches]
- Post-action destinations logical? [Yes/No — worst offenders]

### Visual Cues
- Does the page guide the eye to the right action? [assessment]
- Status indicators clear? (What needs attention, what's done, what's pending)
- Empty states helpful? [Yes/No — which pages]

## Scenario Test Results

### New Hire Onboarding
- Time to first value: [clicks/pages]
- Self-explanatory score: [1-5]
- Terminology barriers: [list]

### Interrupted Workflow
| Interruption | Data preserved? | Could resume? |
|-------------|----------------|---------------|
| Navigate away mid-form | | |
| Close tab + reopen | | |
| Page refresh | | |

### Wrong Turn Recovery
| Wrong turn | Steps to recover | Context lost? |
|-----------|-----------------|---------------|
| Wrong nav section | | |
| Wrong record | | |

### Day Two Efficiency
| Thread | First time | Second time | Shortcuts found |
|--------|-----------|-------------|-----------------|
| | | | |

### Explain It (2-minute guides)
[Include the actual plain-English explanations — they double as documentation]

### What Changed?
- Awareness score: [1-5]
- Notification quality: [assessment]
- Missing signals: [list]

## Emotional Friction

Moments where the user felt uncertain, anxious, or frustrated:

- **[Moment]**: [What happened, why it felt bad, what would feel better]

## Resilience

| Test | Result |
|------|--------|
| Navigate away mid-form | [Pass/Fail — detail] |
| Bad input on forms | [Pass/Fail — detail] |
| Back button behaviour | [Pass/Fail — detail] |
| Page refresh | [Pass/Fail — detail] |

## Network Errors (detected during browsing)

| Endpoint | Status | Page | Severity | Notes |
|----------|--------|------|----------|-------|
| [URL] | [status] | [page you were on] | [Critical/High/Medium/Low] | [likely cause] |

## What Works Well

- [Positive findings — patterns to preserve and replicate]

## Recommendations (priority order)

1. [Highest impact fix — addresses critical/high findings]
2. [Second priority]
3. [Third priority]
```

## QA Sweep Report

```markdown
# QA Sweep: [App Name]

**Date**: YYYY-MM-DD
**Browser**: Chrome MCP / Playwright
**Areas Tested**: [count]

## Summary

[passed] / [total] areas passed. [brief overview of failures]

## Results

| Area | Status | Issues |
|------|--------|--------|
| [Page/Feature] | Pass / Fail | [brief note or "—"] |

## Failed Areas — Detail

### [Area Name]

- **Issue**: [what failed]
- **Steps**: [how to reproduce]
- **Expected**: [what should happen]
- **Actual**: [what happened]

## Cross-Cutting

| Check | Status | Notes |
|-------|--------|-------|
| Dark mode | Pass/Fail | |
| Mobile (375px) | Pass/Fail | |
| Search & filters | Pass/Fail | |
| Notifications | Pass/Fail | |
| Empty states | Pass/Fail | |
```

## Guidelines

- Keep findings concrete — "Submit button doesn't respond" not "form is broken"
- Include page/route for every finding so developers can locate it
- Screenshots are optional but make critical/high findings much more actionable
- "What Works Well" section prevents the report from being only negative
- Priority recommendations should be actionable in one sprint
