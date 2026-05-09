---
name: chat-summarizer
description: Generate a comprehensive, AI-optimized summary of the current Manus chat session. Use when the user asks to summarize this chat, create a session summary, generate a handoff document, capture what was done, or save progress for continuation in a new chat. Also use when the user says "summarize", "recap", "save context", or "handoff".
---

# Chat Summarizer

Generate a structured, comprehensive summary of the current Manus chat session, specifically designed for AI consumption. The output enables seamless continuation in a new chat — another AI instance can read it and pick up exactly where this session left off.

## When to Use

Use this skill when the user asks to summarize the current chat, save session progress, create a handoff document, or capture what was accomplished. The summary is optimized for AI-to-AI context transfer, not human reading (though it is human-readable).

## Workflow

### Step 1: Extract and Write the Summary

Analyze the entire chat history in your context window. Write a Markdown file with exactly 7 sections, using the schema below. Save to `/home/ubuntu/chat_summary.md`.

**Critical extraction principles:**
- Extract FACTS, not impressions. Quote exact file paths, command outputs, error messages.
- Capture RATIONALE for every decision. "We chose X because Y" is 10x more valuable than "We did X."
- Record FAILURES and dead ends. The next AI must not repeat them.
- Be EXHAUSTIVE on artifacts. Every file created, modified, or deleted must be listed with its full path.
- Use PRECISE technical language. Eliminate filler. Every sentence must contain actionable information.

#### Required Schema

```markdown
# Chat Session Summary
Generated: [ISO 8601 timestamp]

## SESSION_METADATA

| Field | Value |
|---|---|
| Date | [Date range of session] |
| Primary project | [Project name and repo if applicable] |
| Working directory | [Primary working directory path] |
| User's machine | [OS, key details if known] |
| Sandbox environment | [Key tools, APIs, packages available] |
| Connected services | [GitHub repos, MCP servers, Google Drive, etc.] |

## OBJECTIVE_TREE

Root goal: [One sentence describing the overarching objective]

[For each sub-goal, use this format:]
- [COMPLETED] Sub-goal description — outcome achieved
- [IN-PROGRESS] Sub-goal description — current state
- [BLOCKED] Sub-goal description — what is blocking it
- [ABANDONED] Sub-goal description — why it was abandoned

## DECISIONS_LOG

[For each significant decision:]

**Decision: [What was decided]**
- Alternatives considered: [What else was evaluated]
- Rationale: [Why this choice was made]
- Confidence: [HIGH/MEDIUM/LOW]
- Impact: [What this decision affected downstream]

## ARTIFACTS_REGISTRY

[For each file created or modified:]

| File | Purpose | State | Dependencies |
|---|---|---|---|
| `/full/path/to/file` | What it does | created/modified/committed/deployed/tested | Other files it depends on |

## DISCOVERIES

[For each bug found, insight gained, or pattern observed:]

**[BUG/INSIGHT/PATTERN]: Title**
- Detail: [Precise technical description]
- Root cause: [If applicable]
- Resolution: [If resolved, how; if not, why]

## UNRESOLVED_ITEMS

[For each open item:]
- **[OPEN QUESTION]** Description — context for why it matters
- **[KNOWN ISSUE]** Description — impact and workaround if any
- **[PENDING FEATURE]** Description — what was requested, current state
- **[BLOCKED]** Description — what is blocking and what would unblock

## CONTINUATION_CONTEXT

**Current system state:** [Exact state of the project/system right now]

**What the next AI must know:**
[Critical context that is NOT obvious from the artifacts alone — implicit
knowledge, user preferences, gotchas, environment quirks]

**Recommended next steps (ordered by priority):**
1. [Most important next action]
2. [Second priority]
3. [Third priority]

**Warnings:**
[Anything the next AI should be careful about — fragile code, user
sensitivities, environment constraints]
```

### Step 2: Self-Review

After writing the summary, perform a completeness check by asking yourself:

1. Could an AI reading ONLY this summary continue the work without asking the user to repeat anything?
2. Are ALL file paths absolute and complete?
3. Is every decision's rationale captured, not just the outcome?
4. Are failures and dead ends documented so they won't be repeated?
5. Is the CONTINUATION_CONTEXT specific enough to be actionable?

If any answer is "no", revise the relevant section before proceeding.

### Step 3: Validate and Format

Run the validation script:

```bash
python3 /home/ubuntu/skills/chat-summarizer/scripts/format_summary.py /home/ubuntu/chat_summary.md --json --handoff
```

This validates all 7 sections, checks file path references, generates a compact JSON version, and a dense handoff prompt.

Optionally, add `--enrich` for AI-generated continuation recommendations (requires `OPENAI_API_KEY`):

```bash
python3 /home/ubuntu/skills/chat-summarizer/scripts/format_summary.py /home/ubuntu/chat_summary.md --json --handoff --enrich
```

If validation fails, fix the missing sections and re-run.

### Step 4: Deliver

Send the user these files (in order of importance):
1. `/home/ubuntu/chat_summary.md` — the full structured summary
2. `/home/ubuntu/handoff_prompt.txt` — dense single-block for pasting into new chats
3. `/home/ubuntu/chat_summary.json` — machine-parseable version
4. `/home/ubuntu/ai_recommendations.md` — AI-generated next steps (only if `--enrich` was used and succeeded)

## Quality Criteria

A good summary scores well on these dimensions:

| Dimension | Indicator |
|---|---|
| Completeness | All 7 sections populated with substantive content |
| Precision | Exact paths, exact error messages, exact decisions — no vague language |
| Continuability | Another AI can continue without asking user to repeat anything |
| Rationale density | Every decision has a "because" or "rationale" |
| Failure coverage | Dead ends and bugs documented with root causes |
| Actionability | CONTINUATION_CONTEXT has specific, ordered next steps |
