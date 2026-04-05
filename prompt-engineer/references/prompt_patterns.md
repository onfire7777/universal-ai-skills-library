# Prompt Patterns Quick Reference

## Pattern Selection Guide

Determine which patterns to apply based on the task type:

| Task Type | Primary Patterns | Secondary Patterns |
|---|---|---|
| Agent Task (build/create/deploy) | Agent-First Framing, One-Shot Density, Tool Leverage | File System Memory, Error Prevention |
| Research/Analysis | Context Injection, Tree of Thoughts, SWOT+ | Reverse Engineering, Progressive Disclosure |
| Code Development | First Principles, Red Team/Blue Team, Specificity | Prompt Chaining, Error Recovery |
| Content Creation | Perspective Shifting, Analogical Reasoning, Constraint Creativity | Iterative Refinement, Format Control |
| Data Processing | Format Control, Template Generator, Specificity | Prompt Chaining, Error Recovery |
| Workflow Automation | Prompt Chaining, Agent-First, Error Recovery | Scheduling, MCP Integration |
| Decision Making | Tree of Thoughts, Self-Consistency, SWOT+ | Red Team/Blue Team, First Principles |

## Optimization Checklist

Apply to every prompt before sending to Manus:

1. **Objective**: Is the goal stated in 1-2 clear sentences?
2. **Deliverables**: Are exact outputs specified (file types, names, structure)?
3. **Context**: Is all necessary background provided?
4. **Constraints**: Are boundaries and exclusions defined?
5. **Format**: Is the output format explicitly specified?
6. **Tools**: Are relevant Manus tools referenced?
7. **Validation**: Are success criteria defined?
8. **Error Handling**: Are fallback behaviors specified?
9. **One-Shot**: Is everything packed into a single prompt?
10. **Agent-First**: Does it leverage multi-step agent capabilities?

## Intent-Specific Optimization Strategies

### Agent Task
Add: clear end-goal, deliverables list, tool hints, directory structure, validation steps.
Remove: unnecessary pleasantries, chatbot-style phrasing, ambiguous language.

### Code Development
Add: language, framework, error handling, testing requirements, architecture, deployment target.
Remove: implementation details Manus can figure out, over-specification of obvious patterns.

### Content Creation
Add: tone, audience, style, format, length, structure, examples of desired quality.
Remove: vague style descriptors, contradictory requirements.

### Research Analysis
Add: scope, depth, sources, evaluation criteria, output format, actionable recommendations.
Remove: open-ended exploration without boundaries.

### Data Processing
Add: input format, transformations, output format, visualization preferences, sample data.
Remove: manual step-by-step instructions for standard operations.

### Workflow Automation
Add: trigger conditions, steps, error handling, success verification, integrations.
Remove: implementation details for standard patterns.

## Credit Optimization Rules

1. **Never iterate when you can specify**: Every follow-up costs nearly as much as the original
2. **Batch related requests**: Combine related tasks into one prompt
3. **Specify don't describe**: "Create a React app with..." not "I'd like something that..."
4. **Include examples**: One good example prevents 3 clarification rounds
5. **Define done**: Clear success criteria prevent unnecessary continuation
6. **Anticipate questions**: Answer likely clarifications preemptively
7. **Use chat mode first**: Refine your prompt in chat mode before running as agent task
