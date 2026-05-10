# Game Research Reference

Purpose: Use this file when conducting competitive game analysis, market research, design reference gathering, or community feedback collection for game planning.

## Contents

- Game source tiers
- Research dimensions
- Collection checklists
- Competing games intelligence template
- Design reference intelligence template
- Market intelligence template
- Community feedback intelligence template
- Game research brief
- Web tool fallback

## Game Source Tiers

| Tier | Source type | Default reliability |
|---|---|---:|
| Tier 1 | official store pages (Steam, PS Store, eShop), official websites, press releases | `0.90` |
| Tier 2 | review aggregators (Metacritic, OpenCritic), industry media (Game Developer, IGN, Eurogamer) | `0.75-0.85` |
| Tier 3 | community sources (Reddit r/gamedesign, Steam Community, Discord, Twitter/X) | `0.60-0.65` |
| Tier 4 | postmortems & academic (GDC Vault, Gamasutra/Game Developer postmortems, game studies journals) | `0.80-0.90` |

Rules:
- start with Tier 1 for factual claims (release dates, platforms, features)
- Tier 4 carries high reliability for design insights but may be dated — note publication year
- use Tier 3 for signal generation and sentiment, not as sole proof
- when sources conflict, report the disagreement explicitly with tier labels

## Research Dimensions

| Dimension | Phase | Purpose |
|---|---|---|
| Competing Games | DISCOVER | Identify direct/indirect competitors, feature matrices, positioning |
| Design References | DESIGN | GDC talks, postmortems, articles on relevant mechanics/systems |
| Market & Audience | DISCOVER, FRAME | Genre market size, growth trends, audience demographics |
| Community Feedback | DISCOVER, VALIDATE | Player sentiment, common complaints, unmet needs |
| Technical References | DESIGN | Engine capabilities, platform constraints, performance benchmarks |

## Collection Checklists

### Competing Games

```markdown
## Competing Games: [Genre / Subgenre]

### Store & Product Pages
- [ ] Steam store page (tags, reviews, concurrent players)
- [ ] Console store pages (PS Store, eShop, Xbox Store)
- [ ] Official website and marketing materials
- [ ] Press releases and announcements

### Reviews & Ratings
- [ ] Metacritic / OpenCritic aggregate score
- [ ] Notable critic reviews (key praise and criticism)
- [ ] Steam review sentiment (recent vs all-time)
- [ ] User review themes

### Player Data
- [ ] SteamDB / SteamCharts (concurrent players, peak, trend)
- [ ] Estimated sales (VG Insights, Game Data Crunch)
- [ ] Twitch / YouTube viewership signals
```

### Design References

```markdown
## Design References: [Mechanic / System]

### GDC & Conferences
- [ ] GDC Vault talks (search by topic, year)
- [ ] CEDEC presentations (Japanese market)
- [ ] DiGRA papers (academic game studies)

### Articles & Postmortems
- [ ] Game Developer (gamedeveloper.com) postmortems
- [ ] Developer blogs and technical write-ups
- [ ] Notable Reddit r/gamedesign discussions

### Books & Frameworks
- [ ] Relevant game design books (Schell, Koster, Adams)
- [ ] Design pattern catalogs
- [ ] Balance methodology papers
```

### Market & Audience

```markdown
## Market Analysis: [Genre / Platform]

### Market Reports
- [ ] Genre market size and growth rate
- [ ] Platform distribution (PC, console, mobile)
- [ ] Regional market breakdown

### Audience Analysis
- [ ] Steam tag analysis (genre overlaps, trending tags)
- [ ] Target demographic profiles
- [ ] Purchasing behavior patterns

### Trend Signals
- [ ] Emerging subgenres or hybrid genres
- [ ] Successful recent releases in genre
- [ ] Steam Next Fest / showcase performance
```

### Community

```markdown
## Community Research: [Game / Genre]

### Reddit
- [ ] r/gamedesign discussions
- [ ] Genre-specific subreddits
- [ ] r/Games and r/gaming threads

### Steam Community
- [ ] Steam discussion forums (competitor titles)
- [ ] Steam review analysis (positive/negative themes)
- [ ] Steam Workshop activity (modding community)

### Discord & Social
- [ ] Genre-focused Discord servers
- [ ] Twitter/X developer discussions
- [ ] YouTube creator coverage
```

## Competing Games Intelligence

```markdown
## Competing Games Intelligence: [Genre]

### Title Comparison Matrix
| Attribute | Our Game | Comp A | Comp B | Comp C |
|---|---|---|---|---|
| Genre/Subgenre | [...] | [...] | [...] | [...] |
| Platform | [...] | [...] | [...] | [...] |
| Price Point | [...] | [...] | [...] | [...] |
| Metacritic | N/A | [...] | [...] | [...] |
| Steam Reviews | N/A | [...] | [...] | [...] |
| Peak Concurrent | N/A | [...] | [...] | [...] |
| Core Loop | [...] | [...] | [...] | [...] |
| Unique Selling Point | [...] | [...] | [...] | [...] |
| Monetization | [...] | [...] | [...] | [...] |
| Release Date | [...] | [...] | [...] | [...] |

### Feature Gap Analysis
| Feature | Our Game | Comp A | Comp B | Priority |
|---|---|---|---|---|
| [...] | ✓/✗/Planned | ✓/✗ | ✓/✗ | H/M/L |

### Competitive Positioning
- Direct competitors: [...]
- Indirect competitors: [...]
- Our differentiation: [...]
- Market gap opportunity: [...]
```

## Design Reference Intelligence

```markdown
## Design Reference Intelligence: [Topic]

### Key Sources
| Source | Type | Year | Tier | Key Insight |
|---|---|---|---|---|
| [GDC talk title] | Talk | [...] | T4 | [...] |
| [Article title] | Article | [...] | T2 | [...] |
| [Postmortem title] | Postmortem | [...] | T4 | [...] |

### Extracted Design Lessons
| Lesson | Source | Applicability | Risk if Ignored |
|---|---|---|---|
| [...] | [...] | High/Med/Low | [...] |

### Design Pattern Summary
- Pattern: [...]
- Proven in: [titles that used it successfully]
- Failed in: [titles where it didn't work, with reasons]
- Adaptation for our game: [...]
```

## Market Intelligence

```markdown
## Market Intelligence: [Genre / Platform]

### Market Overview
| Metric | Value | Source | Tier | Year |
|---|---|---|---|---|
| Genre market size | [...] | [...] | [...] | [...] |
| YoY growth rate | [...] | [...] | [...] | [...] |
| Avg. price point | [...] | [...] | [...] | [...] |
| Top platform share | [...] | [...] | [...] | [...] |

### Audience Profile
| Segment | Size | Engagement | Spending |
|---|---|---|---|
| [...] | [...] | [...] | [...] |

### Trend Analysis
| Trend | Direction | Confidence | Design Implication |
|---|---|---|---|
| [...] | ↑/↓/→ | High/Med/Low | [...] |
```

## Community Feedback Intelligence

```markdown
## Community Feedback Intelligence: [Genre / Topic]

### Sentiment Themes
| Theme | Frequency | Sentiment | Source | Design Implication |
|---|---|---|---|---|
| [...] | High/Med/Low | +/−/Mixed | [...] | [...] |

### Unmet Player Needs
| Need | Evidence | Frequency | Opportunity Score |
|---|---|---|---|
| [...] | [...] | High/Med/Low | H/M/L |

### Common Complaints (Competitor Titles)
| Complaint | Affected Titles | Frequency | Our Advantage |
|---|---|---|---|
| [...] | [...] | [...] | [...] |

### Community Wishlist
| Feature/Improvement | Community Support | Feasibility | Priority |
|---|---|---|---|
| [...] | High/Med/Low | High/Med/Low | H/M/L |
```

## Game Research Brief

Integrated output template combining all research dimensions.

```markdown
## Game Research Brief: [Project Name]

### Research Scope
- Genre: [...]
- Target platforms: [...]
- Research dimensions covered: [Competing Games | Design References | Market | Community | Technical]
- Date: [YYYY-MM-DD]

### Executive Summary
[2-3 sentences summarizing key findings and design implications]

### Competing Games Summary
[Top 3-5 competitors with key differentiators — reference Competing Games Intelligence]

### Design References Summary
[Top 3-5 design insights from GDC/postmortems/articles — reference Design Reference Intelligence]

### Market Summary
[Market size, trends, opportunity — reference Market Intelligence]

### Community Insights Summary
[Key player sentiments, unmet needs — reference Community Feedback Intelligence]

### Design Implications
| Finding | Implication | Recommended Action | Priority |
|---|---|---|---|
| [...] | [...] | [...] | H/M/L |

### Source List
| # | Source | Tier | URL/Reference | Date Accessed |
|---|---|---|---|---|
| 1 | [...] | T[1-4] | [...] | [...] |
```

## Web Tool Fallback

When `WebSearch` or `WebFetch` tools are not available:

1. **Acknowledge the limitation** — Inform the user that web research cannot be performed automatically.
2. **Provide a manual collection guide** — Generate the relevant collection checklist(s) from this reference, pre-filled with specific search queries and URLs to visit.
3. **Request user-supplied data** — Ask the user to paste or upload:
   - Store page URLs or screenshots
   - Review scores and excerpts
   - Market report excerpts
   - Community discussion links or summaries
4. **Process provided data** — Apply source tiers and templates to any data the user supplies.
5. **Document gaps** — Clearly mark which sections of the research brief could not be completed and why.

Fallback output format:

```markdown
## Research Collection Request

Web tools are not available. Please provide data for the following:

### Required (for minimum viable research)
- [ ] [Specific data point 1 with suggested source]
- [ ] [Specific data point 2 with suggested source]

### Optional (for comprehensive research)
- [ ] [Additional data point with suggested source]

### Search Queries to Try
1. [Suggested search query for competing games]
2. [Suggested search query for market data]
3. [Suggested search query for design references]
```
