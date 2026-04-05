# Attention Patterns for Context Anchoring

## Why Anchoring Works

Manus operates in an agent loop where each tool call consumes attention. After 20-50 tool calls, earlier context fades from the model's effective attention window. The anchor file acts as externalized memory that can be re-read to restore focus, mirroring the `todo.md` recitation pattern from Manus's own context engineering research.

## Recitation Schedule

| Session Length | Recitation Frequency | Trigger |
|---|---|---|
| Short (< 10 tool calls) | Once at start | After planning |
| Medium (10-30 tool calls) | Every 10 calls | After completing a sub-task |
| Long (30-100+ tool calls) | Every 5-8 calls | Before starting any new sub-task |

## Drift Detection Patterns

**Pattern 1 — Pre-Task Alignment Check**
Before starting any new sub-task or responding to a new user message, mentally compare the task against the anchor. If alignment is unclear, run `anchor.py check "description of the sub-task"` to get a scored assessment.

**Pattern 2 — Post-Completion Reflection**
After completing a major deliverable, re-read the anchor and ask: "Does this output serve the core topic? Did I drift?" If drift occurred, note it and course-correct.

**Pattern 3 — Tangent Acknowledgment**
When the user sends a message that diverges from the anchor, acknowledge the divergence explicitly: "I notice this diverges from our core focus on [anchor topic]. I'll handle this, then return to the main thread." This keeps both Manus and the user aware of the anchor.

## Integration with todo.md

The context anchor and `todo.md` serve complementary purposes:

- **Context anchor** = WHY (the overarching purpose and topic)
- **todo.md** = WHAT (the specific tasks and progress)

When both exist, recite them together. Read the anchor first (purpose), then the todo (progress). This creates a complete orientation: "I know why I'm here and what I need to do next."

## Handling Anchor Conflicts

If a user explicitly asks for something outside the anchor's boundaries:

1. Acknowledge the anchor exists and note the divergence
2. Ask if the user wants to update the anchor to include the new direction
3. If yes, use `anchor.py update` to evolve the anchor
4. If no, handle the request as a one-off and return to the anchor

Never silently ignore the anchor. Never refuse a user request because of the anchor — the anchor is a guide, not a constraint. The user always has final authority.

## Advanced: Multi-Layer Anchoring

For complex projects, the anchor can define layers:

- **Layer 1 — Mission**: The overarching goal (e.g., "Build a SaaS analytics platform")
- **Layer 2 — Phase**: The current project phase (e.g., "Phase 2: Backend API development")
- **Layer 3 — Sprint**: The immediate focus (e.g., "Implementing user authentication endpoints")

Use `anchor.py update --refine` to narrow focus as the project progresses, while preserving the broader mission context.
