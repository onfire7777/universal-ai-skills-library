# Analysis Framework

Before writing the infographic prompt, create `analysis.md` first.

## Output Template

```markdown
# Infographic Analysis

## Topic
- Main topic: <topic>
- Audience: <target audience>
- Language: <zh|en|ja|ko>

## Content Shape
- Density: <low|medium|high>
- Primary structure: <list|comparison|timeline|hierarchy|process|mixed>
- Data presence: <none|light|heavy>
- Text-in-image needed: <yes|no>

## Core Message
- Main message: <one-sentence summary>
- Supporting points:
  - <point 1>
  - <point 2>
  - <point 3>

## Visual Implications
- Best layout candidates:
  - <layout 1>
  - <layout 2>
- Best style candidates:
  - <style 1>
  - <style 2>
- Recommended aspect:
  - <3:4|4:3|16:9>

## Risks
- <too much density, too much image text, incomplete data, or similar>
```

## Decision Rules

- If the core structure is step-by-step, prefer `process`.
- If the core structure is side-by-side contrast, prefer `comparison`.
- If the page needs many small modules, mark the density as `high`.
- If the model will need to render a lot of on-canvas text, note that explicitly in `Risks`.
