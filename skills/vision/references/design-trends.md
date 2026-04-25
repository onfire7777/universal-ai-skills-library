# Design Trends & AI Tools Integration

Purpose: Use this file when the task involves trend selection, trend risk, or AI-assisted design workflows.

Contents:
- Visual trend risk buckets
- AI interface patterns
- Trend-application rules
- AI tool guardrails

## Visual Trends (2025-2026)

| Trend | Risk | Best For |
|-------|------|----------|
| Dark Mode | Low | all products |
| Micro-animations | Low | interactive elements |
| AI-Native Interfaces | Low | AI-powered products |
| Variable Fonts | Low | all products |
| Adaptive UI | Low | personalized apps |
| Bento Grid | Medium | dashboards, portfolios |
| Glassmorphism 2.0 | Medium | overlays, cards |
| Spatial Design | Medium | XR-oriented products |
| Sustainable Design | Medium | eco-conscious brands |
| Neo-Brutalism | High | creative or tech brands only |

## AI Interface Patterns

| Pattern | Use |
|---------|-----|
| `Chat-First` | AI assistants and search |
| `Inline Suggestions` | editors and forms |
| `Progressive Disclosure` | complex workflows |
| `Confidence Indicators` | AI-generated output |
| `Regeneration UI` | generative tools |

## Trend Application Rules

`Apply confidently`
- dark mode support
- micro-animations in the `100-300ms` range
- variable fonts
- adaptive layouts

`Apply carefully`
- glassmorphism
- spatial depth
- bento layouts
- oversized typography
- sustainable-design aesthetics

`Apply sparingly`
- neo-brutalism
- kinetic typography
- extreme minimalism
- heavy 3D

Before applying any trend:
- [ ] Brand fit is explicit
- [ ] Target users are likely to expect it
- [ ] `WCAG 2.2 AA` can still be met
- [ ] Performance cost is acceptable
- [ ] The pattern should age well for `2-3 years`

## AI Tool Guardrails

| Tool | Best For | Guardrail |
|------|----------|-----------|
| Figma AI | layout exploration | review for brand, tokens, and a11y |
| v0 | component scaffolding | route through `Forge`, then production agents |
| Claude Artifacts | concept generation | use for exploration, not direct production |
| Galileo AI | fast UI concepts | require Vision review before adoption |
| Midjourney / DALL-E | moodboards and assets | check brand fit and licensing |

Use AI generation for:
- initial concept exploration with `3+` variations
- moodboards
- layout alternatives
- placeholder assets

Do not use AI generation as the final authority for:
- production code
- brand-critical identity work
- accessibility-sensitive components
- performance-critical implementations
