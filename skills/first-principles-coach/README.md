# First Principles Thinking

An agent skill for Socratic questioning and first-principles problem solving.

## What it does

Guides users through structured questioning to surface hidden assumptions, reach fundamental truths, and rebuild solutions from scratch.

**Use when:**
- Stuck on a problem where conventional solutions aren't working
- Making a high-stakes decision that warrants deeper analysis
- Building something new (not optimizing existing)
- Facing "industry standard" constraints that feel arbitrary

## Install

```bash
npx skills add akshat10/skills --skill first-principles-thinking
```

## Example

```
User: "I want to start a podcast but it seems expensive."

Claude: "Let's break this down. When you say 'expensive',
        what specific costs are you imagining?"

User: "Microphone, editing software, hosting, marketing..."

Claude: "What's your assumption about what a microphone needs to cost?"

User: "I've seen podcasters use $300-400 mics."

Claude: "Why do you believe you need that tier of equipment?"

[Continue until reaching bedrock truths, then rebuild from fundamentals]
```

## Based on

- Richard Paul's 6 types of Socratic questions
- Elon Musk's first principles approach at SpaceX
- Kim Scott's Radical Candor (for delivering challenging questions with care)
