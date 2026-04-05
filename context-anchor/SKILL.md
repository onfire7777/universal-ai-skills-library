---
name: context-anchor
description: Set a persistent core topic/purpose for the entire chat session to keep Manus focused and on-track. Use when the user asks to set a core topic, anchor the context, maintain focus, or ensure Manus remembers the overarching goal regardless of tangential requests.
---

# Context Anchor

This skill manages a persistent "context anchor" — a core topic or overarching goal for the session. By externalizing the core purpose to a file and periodically reciting it, Manus maintains deep focus across long sessions (50+ tool calls) and prevents cognitive drift when handling tangential requests.

## How It Works

1. **Set**: You create an anchor file (`~/.context_anchor.md`) with the user's core topic.
2. **Recite**: You re-read the anchor periodically to keep it in your active attention.
3. **Check**: You verify if new tasks or user messages align with the anchor.
4. **Update**: You evolve the anchor as the project progresses.

## Quick Start

```bash
# Set a basic anchor
python3 /home/ubuntu/skills/context-anchor/scripts/anchor.py set "Building a SaaS analytics dashboard"

# Set a detailed anchor with objectives, boundaries, and success criteria
python3 /home/ubuntu/skills/context-anchor/scripts/anchor.py set "Building a SaaS analytics dashboard" \
  -o "Design database schema" "Build REST API" \
  -b "No mobile app" "No payment processing" \
  -s "All CRUD endpoints work"

# View current anchor
python3 /home/ubuntu/skills/context-anchor/scripts/anchor.py show
```

## Workflow: Maintaining Focus

Once an anchor is set, you MUST adopt the following attention patterns:

1. **Periodic Recitation**: Every 5-10 tool calls, run `python3 /home/ubuntu/skills/context-anchor/scripts/anchor.py show` to refresh your memory of the core topic. Do this alongside checking your `todo.md` if one exists.
2. **Pre-Task Check**: Before starting any major sub-task, mentally verify it serves the anchor. If unsure, run `python3 /home/ubuntu/skills/context-anchor/scripts/anchor.py check "description of task"`.
3. **Handling Tangents**: If the user asks for something outside the anchor's scope (e.g., asking about a recipe while building a web app):
   - Acknowledge the divergence: "I notice this diverges from our core focus on [anchor topic]."
   - Fulfill the user's request (the user always has final authority).
   - Explicitly state your return to the core focus afterward.

## Updating the Anchor

As the project evolves, refine the anchor instead of clearing it:

```bash
# Add a new objective
python3 /home/ubuntu/skills/context-anchor/scripts/anchor.py update --add-objective "Implement WebSocket support"

# Add a boundary (out of scope)
python3 /home/ubuntu/skills/context-anchor/scripts/anchor.py update --add-boundary "Skip user authentication for now"

# Refine the core topic (narrows focus)
python3 /home/ubuntu/skills/context-anchor/scripts/anchor.py update --refine "Phase 2: Backend API development for the analytics dashboard"
```

## Advanced Attention Patterns

For deep insights into why anchoring works and advanced drift-detection strategies, read:
`/home/ubuntu/skills/context-anchor/references/attention_patterns.md`

## Commands Reference

- `set [topic] [-o objectives...] [-b boundaries...] [-s success...]`: Create/overwrite anchor
- `show`: Display current anchor
- `check [message]`: Score relevance of a message against the anchor
- `update [--add-objective text] [--add-boundary text] [--refine text]`: Modify anchor
- `history`: Show change log of the anchor
- `clear`: Delete the anchor
