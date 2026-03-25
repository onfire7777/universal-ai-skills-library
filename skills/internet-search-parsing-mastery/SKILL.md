---
name: Internet Search & Parsing Mastery
description: Advanced internet search techniques, web content extraction, OSINT methodology, fact-checking, and information synthesis for finding and parsing the best information from the web.
version: 1.0.0
author: Custom Meta-Skill
tags: [search, internet, parsing, web-scraping, OSINT, fact-checking, information-retrieval]
---

# Internet Search & Parsing Mastery

## Purpose
Master the art and science of finding, extracting, validating, and synthesizing information from the internet with maximum efficiency and accuracy.

## Search Strategy Framework

### The 4-Phase Search Protocol

**Phase 1: Query Formulation**
Before searching, define:
- What exactly do I need to know?
- What type of source would have this information? (academic, official, community, commercial)
- What are the key terms and their synonyms?
- What date range is relevant?

**Phase 2: Multi-Source Search**
Never rely on a single search. Execute searches across:
- General web search (for broad coverage)
- Specialized databases (for depth)
- Community sources (Reddit, Stack Overflow, HN for practitioner knowledge)
- Official sources (documentation, government sites, company blogs)
- Academic sources (Google Scholar, arXiv for research)

**Phase 3: Content Extraction**
- Extract the actual content, not just snippets
- Navigate to primary sources, not summaries
- Save key findings immediately to prevent loss
- Note the date, author, and context of each source

**Phase 4: Synthesis & Validation**
- Cross-reference findings across sources
- Identify consensus vs. disagreement
- Weight by source quality and recency
- Produce a synthesized answer with citations

### Advanced Search Techniques

**Query Expansion**: For any search, generate 3 query variants:
1. The literal question
2. The question rephrased with synonyms
3. The question from a different angle or domain

**Specificity Ladder**:
- Start specific → if too few results, broaden
- Start broad → if too many results, narrow with additional terms

**Domain-Specific Search**:
- `site:github.com` for code and repos
- `site:stackoverflow.com` for programming Q&A
- `site:reddit.com` for community discussions
- `site:arxiv.org` for research papers
- `site:gov` for government data
- `filetype:pdf` for documents and reports

**Temporal Search**:
- Use date filters for time-sensitive topics
- Search for "[topic] 2025" or "[topic] latest" for recency
- Check if older results have been superseded

### Web Content Extraction Best Practices

**Extraction Priority Order**:
1. **APIs first**: Check if the site has a public API (REST, GraphQL)
2. **Structured data**: Look for JSON-LD, microdata, RSS/Atom feeds
3. **Semantic HTML**: Extract from `<article>`, `<main>`, `<h1>`-`<h6>`, `<p>` tags
4. **Markdown conversion**: Convert HTML to clean Markdown for processing
5. **Full page rendering**: For SPAs and JS-heavy sites, render first then extract
6. **Screenshot analysis**: For visual content, charts, infographics

**Content Quality Signals**:
- Author credentials and expertise
- Publication date and update history
- Peer review or editorial process
- Citation count and backlinks
- Domain authority and reputation
- Presence of references and sources

### Fact-Checking Protocol

For every important claim:
1. **Source Check**: Who is making this claim? What's their credibility?
2. **Evidence Check**: What evidence supports it? Is it primary or secondary?
3. **Corroboration Check**: Do 2+ independent sources confirm it?
4. **Contradiction Check**: Are there credible sources that disagree?
5. **Recency Check**: Is this information still current?
6. **Context Check**: Is the claim being taken out of context?

**Red Flags for Misinformation**:
- No author or anonymous source
- No dates or very old content
- Emotional language designed to provoke
- No citations or references
- Claims that seem too good/bad to be true
- Single-source claims with no corroboration
- Sites with excessive ads or clickbait

### Information Synthesis Protocol

**When combining information from multiple sources**:
1. Create a source matrix: Source × Key Finding × Quality Rating
2. Identify areas of convergence (high confidence)
3. Identify areas of divergence (investigate further)
4. Weight by source quality, recency, and relevance
5. Produce a synthesized narrative with inline citations
6. Explicitly state confidence level for each conclusion
7. Note gaps — what questions remain unanswered?

### Rate Limiting & Ethical Search
- Respect robots.txt and terms of service
- Don't hammer servers with rapid requests
- Cache results to avoid redundant searches
- Use official APIs when available
- Attribute sources properly
- Don't scrape personal/private information
