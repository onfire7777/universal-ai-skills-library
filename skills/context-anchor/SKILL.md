---
name: context-anchor
description: Set and maintain a persistent core topic/purpose for the entire chat session to keep Manus focused. Use when the user sets a core topic, anchors context, says "let's focus on X", "today we're building Y", "we're working on Z", or any implicit/explicit session purpose framing. Also use when asked to maintain focus, stay on track, or remember the overarching goal.
---

# Context Anchor

Persistent attention anchoring for Manus sessions. Externalizes the core session purpose to `~/.context_anchor.md` and enforces periodic recitation to maintain deep focus across 50+ tool calls, preventing cognitive drift from tangential requests.

## Activation

Detect **explicit** triggers: user says "set anchor", "anchor this", "core topic is", "focus on".
Detect **implicit** triggers: user frames session purpose with "today we're building X", "let's work on Y", "this project is about Z". Summarize the detected intent and confirm with the user before setting.

On session start, check if `~/.context_anchor.md` exists. If it does, run `anchor.py show` immediately and recite the anchor to orient yourself.

## Quick Start

```bash
# Set anchor with full context
python3 /home/ubuntu/skills/context-anchor/scripts/anchor.py set "Building a SaaS analytics dashboard" \
  -o "Design database schema" "Build REST API" "Create React frontend" \
  -b "No mobile app" "No payment processing" \
  -s "All CRUD endpoints work" "Dashboard renders real-time metrics"

# View current anchor
python3 /home/ubuntu/skills/context-anchor/scripts/anchor.py show

# Check if a sub-task aligns with the anchor (AI-powered relevance scoring)
python3 /home/ubuntu/skills/context-anchor/scripts/anchor.py check "Should I add Stripe integration?"
```

## Recitation Discipline

This is the most critical behavior. You MUST enforce periodic recitation to keep the anchor in your active attention window.

**Every 5-10 tool calls** (aim for ~7), run `anchor.py show` and mentally reaffirm: "My core focus is [anchor topic]." If a `todo.md` also exists, recite both together — anchor first (WHY), then todo (WHAT). If 10 tool calls elapse without recitation, force an immediate refresh.

On recitation failure (file missing or inaccessible), fall back to your last known anchor from memory. Log the failure silently and retry next interval. Notify the user only if the anchor file is persistently unavailable.

## Handling Tangential Requests

When the user sends a request that diverges from the anchor:

1. Run `anchor.py check "user's request"` if relevance is unclear.
2. If the score is below 5/10 (divergent), acknowledge: "I notice this diverges from our core focus on [anchor topic]. I'll handle this, then return to the main thread."
3. Fulfill the request completely — the user always has final authority over the anchor.
4. After completing the tangent, explicitly restate the anchor: "Returning to our core focus: [anchor topic]."

Never refuse a user request because of the anchor. The anchor is a guide, not a constraint.

## Evolving the Anchor

As projects progress, refine the anchor rather than clearing it:

```bash
# Add a new objective
python3 /home/ubuntu/skills/context-anchor/scripts/anchor.py update --add-objective "Implement WebSocket for real-time updates"

# Add a boundary
python3 /home/ubuntu/skills/context-anchor/scripts/anchor.py update --add-boundary "Skip authentication for MVP"

# Narrow focus to current phase
python3 /home/ubuntu/skills/context-anchor/scripts/anchor.py update --refine "Phase 2: Backend API development for the analytics dashboard"
```

Confirm with the user before making anchor changes. Do not update more than once per 3 tool calls to prevent instability.

## Integration with todo.md

The anchor and `todo.md` are complementary:

| File | Purpose | Contains |
|---|---|---|
| `~/.context_anchor.md` | WHY — the overarching purpose | Core topic, objectives, boundaries, success criteria |
| `todo.md` | WHAT — the specific tasks | Current tasks, progress, next steps |

When both exist, recite them together during periodic recitation. On anchor changes, review `todo.md` to ensure task alignment with the updated purpose.

## Commands Reference

| Command | Usage |
|---|---|
| `set` | `anchor.py set "topic" [-o objectives...] [-b boundaries...] [-s success...]` |
| `show` | `anchor.py show` |
| `check` | `anchor.py check "message to evaluate"` |
| `update` | `anchor.py update [--add-objective text] [--add-boundary text] [--refine text]` |
| `history` | `anchor.py history` |
| `clear` | `anchor.py clear` |

## Advanced Patterns

For deep insights into why anchoring works, recitation schedules by session length, drift detection patterns, multi-layer anchoring, and anchor conflict resolution, read: `/home/ubuntu/skills/context-anchor/references/attention_patterns.md`
