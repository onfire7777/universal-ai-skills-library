# Manus Prompt Engineering Techniques Reference

## Manus Context Engineering Principles

These principles come directly from the Manus engineering team's blog on context engineering. Apply them when crafting or optimizing prompts.

### 1. One-Shot Density

Credit usage grows almost multiplicatively with each follow-up message. A single dense, well-explained prompt is dramatically more efficient than iterative refinement. Even minimal follow-ups like "change the color" cost nearly as much as the original prompt.

**Application**: Pack all requirements, constraints, deliverables, format specifications, and success criteria into ONE prompt. Anticipate follow-up questions and answer them preemptively.

### 2. Agent-First Framing

Manus is an autonomous agent with a sandbox VM, not a chatbot. Prompts should leverage multi-step execution, tool use, file creation, web research, and structured outputs. If a chatbot could answer it in one reply, the prompt is underutilizing Manus.

**Agent trigger phrases**: "Take ownership", "Track this over time", "Guide me step by step", "Follow up when needed", "Figure out what needs to be done", "Manage the process until done".

### 3. File System as Memory

Manus uses the file system as externalized, persistent memory. For complex tasks, instruct it to create planning files (todo.md), save intermediate results, and organize outputs in clear directory structures. This prevents context loss over long tasks.

### 4. Attention Manipulation via Recitation

Manus creates and updates todo.md files to keep goals in its recent attention span. For complex multi-step tasks, explicitly request this behavior to prevent goal drift over ~50+ tool calls.

### 5. Error Prevention Through Specificity

Most Manus errors are preventable through prompting alone. Specify constraints, edge cases, fallback behaviors, and validation criteria upfront. Be imperative — Manus is an "unfocused lens" that produces clear results with enough adjustment.

### 6. Leverage Manus Tools Explicitly

Reference specific capabilities when relevant:

| Tool | When to Reference |
|---|---|
| Shell/Terminal | System operations, package installation, file manipulation |
| Python/Code | Data processing, automation, calculations, API calls |
| Browser | Web research, form filling, screenshots, web scraping |
| Search | Information gathering, fact-checking, finding resources |
| File System | Creating documents, organizing outputs, persistent storage |
| Image Generation | Visual content, diagrams, illustrations |
| Slides | Presentations, slide decks |
| Scheduling | Recurring tasks, timed execution |
| MCP | External service integrations |
| GitHub | Repository management, code deployment |
| Google Drive | Document management, collaboration |
| Map (Parallel) | Batch processing of similar items |

## Prompt Structure Templates

### Agent Task Template
```
**Objective**: [1-2 sentence clear goal]

**Context**: [Background, constraints, existing resources]

**Requirements**:
1. [Specific deliverable with quality standard]
2. [Specific deliverable with quality standard]
3. [Specific deliverable with quality standard]

**Process**:
- [Step 1 with tool hints]
- [Step 2 with tool hints]
- [Validation step]

**Output Format**: [Exact structure, file types, naming]

**Success Criteria**: [How to verify completion]

**Constraints**: [What to avoid, edge cases, limits]
```

### Research Task Template
```
**Research Goal**: [What to investigate]

**Scope**: [Boundaries, depth, time period]

**Sources**: [Where to look, priority sources]

**Analysis Framework**: [How to evaluate findings]

**Deliverable**: [Format, structure, length]

**Key Questions**: [Specific questions to answer]
```

### Code Development Template
```
**Build**: [What to create]

**Tech Stack**: [Languages, frameworks, dependencies]

**Architecture**: [Structure, patterns, file organization]

**Features**: [Prioritized feature list with specs]

**Quality Standards**: [Testing, error handling, documentation]

**Deployment**: [Where and how to deploy]
```

## Advanced Patterns for Manus

### Tree of Thoughts
Use when: Strategic planning, complex decisions, creative problem-solving.
Pattern: "Generate 3 different approaches to [goal], evaluate each on [criteria], then implement the best one."

### Self-Consistency
Use when: Critical decisions, reducing errors.
Pattern: "Approach this problem 3 different ways, then identify the most consistent and reliable solution."

### Reverse Engineering
Use when: Goal planning, roadmaps.
Pattern: "I want to achieve [end state]. Work backwards to create the step-by-step path from where I am now."

### Red Team / Blue Team
Use when: Strategy validation, risk assessment.
Pattern: "First, find every flaw in [plan]. Then, defend or improve each weakness."

### Context Injection
Use when: Personalized or domain-specific tasks.
Pattern: "Before starting, ask me for: [list of context needed]. Then proceed with the task."

### Prompt Chaining
Use when: Complex multi-stage workflows.
Pattern: "Step 1: [analyze]. Step 2: [generate based on analysis]. Step 3: [refine based on generation]. Show work at each step."

## Common Anti-Patterns to Avoid

| Anti-Pattern | Problem | Fix |
|---|---|---|
| Vague goals | Manus guesses wrong, wastes credits | State exact deliverables |
| No format spec | Output doesn't match expectations | Specify file types, structure |
| Iterative refinement | Each follow-up costs full credits | Pack everything in one prompt |
| Chatbot-style questions | Underutilizes agent capabilities | Frame as multi-step tasks |
| Missing constraints | Manus makes unwanted choices | List what to avoid explicitly |
| No success criteria | No way to verify completion | Define "done" conditions |
| Overly long prompts | Dilutes focus, wastes tokens | Be dense, not verbose |
| No error handling | Loops on failures | Specify fallback behaviors |
