---
name: track-improve
description: "Track mistakes and improvement opportunities with 5-whys root cause analysis. Use when Claude makes an error, gives a wrong answer, misunderstands intent, or when you spot a pattern worth capturing. Triggers on: 'track this mistake', 'log this issue', 'remember this for next time', 'you keep doing this', 'improvement opportunity', 'track improvement', 'what went wrong'."
---

# Track & Improve

Capture mistakes and improvement opportunities with root cause analysis.

## Track a Mistake or Improvement

When invoked:

1. Extract the description of what went wrong
2. Identify what persona/role was active
3. Perform 5 whys root cause analysis
4. Extract last 5 message exchanges for context
5. Create a report in `~/.claude/trk-db/active/` with format:

```markdown
---
id: YYYY-MM-DD-HH-MM-SS
created: ISO8601 timestamp
project: current project name
status: active
category: rule-violation|improvement|confusion
---

# [Description]

## 5 Whys Analysis
1. **Why did this happen?** [Analysis]
2. **Why?** [Analysis]
3. **Why?** [Analysis]
4. **Why?** [Analysis]
5. **Root cause:** [Conclusion]

## Context (Last 5 Messages)
[Extracted conversation]

## Resolution
[Empty - filled when resolved]
```

6. Git add and commit

## Review Reports

Read all files from `~/.claude/trk-db/active/`, present organized by recency and category, discuss patterns.

## Resolve a Report

Move from `active/` to `resolved/`, fill in resolution, git commit.

## Rules
🚨 The 5 whys must be genuine root cause analysis, not surface restating.
🚨 Category must accurately reflect whether it's a rule violation, improvement idea, or confusion.
🚨 Reports are tracked in git for history and pattern analysis.