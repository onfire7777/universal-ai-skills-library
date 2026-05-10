# Content Strategy for Design

Purpose: Use this reference when writing content that will be integrated into visual design — particularly landing pages, marketing sites, and product pages where copy and layout are inseparable.

## Contents

- Product Language Principle
- The 30% Cut Rule
- Copy-First Design Process
- Hero Copy Contract
- Section Copy Patterns
- Content and Composition Alignment
- Content Strategy Checklist

---

## Product Language Principle

**Write about the product, not about the design.**

| Do | Don't |
|----|-------|
| "Ship faster with automated deployments" | "A beautiful deployment experience" |
| "Find anything in your codebase in < 1 second" | "A sleek, modern search interface" |
| "Your team's knowledge, organized" | "An elegantly designed knowledge base" |
| "Track every change, revert any mistake" | "A stunning version control dashboard" |

### Why This Matters

Design commentary ("beautiful", "elegant", "modern", "sleek") is:
- **Empty** — it describes the surface, not the value
- **Unverifiable** — the user decides if it's beautiful, not the copy
- **Generic** — every product claims to be beautiful
- **Self-serving** — it's about the maker, not the user

Product language:
- **Specific** — names the actual benefit or capability
- **Verifiable** — the user can confirm the claim
- **Differentiating** — unique to what this product actually does
- **User-centered** — frames value from the user's perspective

---

## The 30% Cut Rule

After writing any piece of copy, **cut 30% of the words.** Then check if meaning is preserved. If yes, the cut version ships.

### Process

1. Write the first draft naturally
2. Count the words
3. Remove 30% of the word count
4. Verify the core message survives
5. If meaning is lost, restore only the minimum words needed

### Examples

| First Draft (100%) | After 30% Cut |
|--------------------|---------------|
| "Our platform helps you to easily manage and organize all of your team's projects in one single place" (18 words) | "Manage all your team's projects in one place" (8 words) |
| "Get started today and see how our solution can transform the way your team collaborates on a daily basis" (19 words) | "Transform how your team collaborates daily" (6 words) |
| "We've built a powerful analytics dashboard that gives you real-time insights into your application's performance" (15 words) | "Real-time insights into your app's performance" (7 words) |

### Where the Cut Hits Hardest

| Word Type | Usually Cut | Sometimes Kept |
|-----------|------------|----------------|
| Adverbs (easily, quickly, simply) | Always cut | Never — they weaken, not strengthen |
| Filler phrases ("in order to", "that allows you to") | Always cut | — |
| Hedging ("helps you", "enables you to") | Usually cut | Keep only when directness would be inaccurate |
| Adjectives (powerful, advanced, comprehensive) | Usually cut | Keep only when specific and verifiable |

---

## Copy-First Design Process

**Write headlines before choosing layouts.** Copy determines composition, not the other way around.

### Why Copy-First

| Layout-First Problem | Copy-First Solution |
|---------------------|-------------------|
| "We need 3 cards here" → forces content into card format | Content determines whether cards are needed |
| "Fill this hero with text" → headline retrofitted to space | Headline length and message determine hero layout |
| "Add a subtitle here" → unnecessary text for layout symmetry | Only text that serves the user appears |
| "We need more sections" → padding with filler content | Page length matches actual content value |

### Process

1. **Write all headlines first** (H1, H2, H3 for the page)
2. **Test the headline story** — read only headlines top-to-bottom. Does it tell a coherent story?
3. **Write supporting copy** for each headline
4. **Apply the 30% cut** to all supporting copy
5. **Now design the layout** around the finalized copy

### Headline Story Test

```markdown
## Headline Story Test: [Page Name]

Read these headlines in order. They should form a coherent narrative:

1. H1: [Hero headline]
2. H2: [First section headline]
3. H2: [Second section headline]
4. H2: [Third section headline]
5. H2: [Final CTA headline]

Coherent narrative? [Yes / No — if No, revise headlines]
```

---

## Hero Copy Contract

The hero section requires exactly this copy:

| Element | Rules | Character Limit |
|---------|-------|----------------|
| **Headline** | 1 sentence. Product language. Specific benefit or capability | 60 characters (6-12 words) |
| **Subline** | Supports headline. Answers "how" or "for whom" | 120 characters (15-25 words) |
| **Primary CTA** | 1 action. Verb-first. Specific outcome | 25 characters (2-4 words) |
| **Secondary CTA** (optional) | Lower commitment alternative | 30 characters (2-5 words) |

### Headline Formula Options

| Formula | Example |
|---------|---------|
| [Action] + [Object] + [Speed/Quality modifier] | "Deploy to production in 60 seconds" |
| [Your/Your team's] + [Desired outcome] | "Your team's best work, organized" |
| [Verb] + [what] + [differentiator] | "Search your entire codebase instantly" |
| [Outcome], [without the pain point] | "Ship daily, without the merge conflicts" |

### CTA Label Rules

| Do | Don't |
|----|-------|
| "Start building" | "Get started" (vague) |
| "Try free for 14 days" | "Learn more" (passive) |
| "See it in action" | "Click here" (mechanical) |
| "Create your first project" | "Sign up" (commitment without context) |

---

## Section Copy Patterns

### Support Section (Trust & Context)

| Pattern | When to Use | Copy Structure |
|---------|-------------|---------------|
| Social proof | When credibility needs establishing | "[Company] uses [Product] to [specific outcome]" |
| Problem statement | When the pain point is relatable | "You know the feeling: [specific frustration]" |
| Key benefit trio | When 3 clear benefits exist | 3 short benefit statements, no icons needed |

### Detail Section (Features)

| Pattern | When to Use | Copy Structure |
|---------|-------------|---------------|
| Feature narrative | For primary features | Headline + 2-sentence explanation + screenshot |
| Comparison | When replacing a competitor/workflow | Before/After or Without/With framing |
| How it works | When the product has a clear workflow | 3-step process with concrete descriptions |

### Final CTA Section

| Element | Rule |
|---------|------|
| Headline | Reinforce the core value proposition (can mirror H1 with variation) |
| Subline | Address remaining objections ("No credit card required", "Free for teams up to 5") |
| CTA | Same as hero CTA (consistency builds confidence) |

---

## Content and Composition Alignment

Copy and layout decisions are interdependent. Use this mapping:

| Copy Decision | Layout Implication |
|--------------|-------------------|
| Short headline (< 6 words) | Can pair with large hero image |
| Long headline (> 10 words) | Needs more text space, image secondary |
| 3 key benefits | Horizontal layout or vertical stack — NOT a card grid |
| Feature with screenshot | Side-by-side layout (text + image alternating) |
| Testimonial quote | Inline within narrative, not a card |
| Process/workflow | Numbered steps or visual flow, not bullet list |

---

## Content Strategy Checklist

```markdown
## Content Strategy Review

### Product Language
- [ ] No design commentary (beautiful, elegant, modern, sleek)
- [ ] Copy describes what the product does, not how it looks
- [ ] Benefits are specific and verifiable

### 30% Cut
- [ ] All copy has been through the 30% cut
- [ ] No adverbs (easily, quickly, simply)
- [ ] No filler phrases or hedging
- [ ] Adjectives are specific and earned

### Copy-First
- [ ] Headlines written before layout decisions
- [ ] Headline story test passes (coherent narrative)
- [ ] Layout serves the copy, not the reverse

### Hero
- [ ] 1 headline (≤ 60 characters)
- [ ] 1 subline (≤ 120 characters)
- [ ] 1 primary CTA (verb-first, specific)
- [ ] Optional secondary CTA (lower commitment)

### Page Flow
- [ ] Each section headline advances the narrative
- [ ] Final CTA reinforces core value
- [ ] No filler sections (every section serves a named purpose)
```
