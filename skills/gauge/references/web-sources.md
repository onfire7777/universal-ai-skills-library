# Web Sources

**Purpose:** Information source tier classification, search query templates, and freshness management rules.
**Read when:** Entering the EVOLVE phase or when web research is required.

---

## Source Tier Classification

| Tier | Description | Trust level | Examples | Usage |
|------|-------------|-------------|---------|-------|
| **T1** | Official primary sources | Highest | Anthropic documentation, MCP specification, Claude Code official docs | Accept directly; cite as authoritative |
| **T2** | Official secondary + major community | High | Anthropic blog posts, official MCP examples, major GitHub repos (10k+ stars), Cursor/Windsurf docs | Accept with verification; cross-reference with T1 when possible |
| **T3** | Community knowledge | Medium | Reddit r/ClaudeAI, developer blogs, X/Twitter posts from verified practitioners, YouTube tutorials from known experts | Use as signals only; require T1/T2 corroboration for checklist changes |
| **T4** | Academic and research | Variable | Multi-agent system papers, prompt engineering research, AI safety publications | Use for theoretical grounding; never as sole basis for checklist changes |

### Tier Usage Rules

1. **Checklist changes (Safety Level C):** Require at least 1 T1 or 2 T2 sources.
2. **Detection pattern updates (Safety Level B):** Require at least 1 T2 source.
3. **Reference additions (Safety Level A):** T3 acceptable if clearly attributed.
4. **Never** accept T4 sources alone without T1/T2 cross-reference for any operational change.

---

## Search Query Templates

### Anthropic / Claude Code Best Practices

```
"Claude Code" "SKILL.md" best practices
"Claude Code" agent design patterns 2026
Anthropic "multi-agent" skill design
"claude code" "skill file" structure
site:docs.anthropic.com agent skill
```

### MCP Specification and Patterns

```
"Model Context Protocol" agent design
MCP "skill definition" patterns
"MCP server" agent orchestration best practices
```

### Multi-Agent System Design

```
multi-agent system "skill template" design
AI agent "self-evolution" safety patterns
"agent collaboration" protocol design
LLM agent "normalization" checklist
```

### Prompt Engineering for Agent Skills

```
"prompt engineering" agent specialization 2026
"system prompt" structure best practices
AI agent "role definition" patterns
"agent boundaries" design patterns
```

### Community Practices

```
site:reddit.com/r/ClaudeAI skill design
"claude code" agent tips
"AI agent" skill "quality assurance"
```

---

## Freshness Management

### Staleness Thresholds

| Content type | Stale after | Action |
|--------------|-------------|--------|
| T1 sources (official docs) | 90 days | Re-verify; update citations |
| T2 sources (blogs, repos) | 60 days | Re-verify; flag if outdated |
| T3 sources (community) | 30 days | Re-verify; consider removal |
| T4 sources (academic) | 180 days | Re-verify relevance |
| Detection patterns | 60 days without recalibration | Trigger GT-04 |
| This file (web-sources.md) | 90 days | Full refresh of query templates |

### Freshness Check Protocol

1. Record last-verified date for each source citation.
2. During EVOLVE phase, check all cited sources against staleness thresholds.
3. Flag stale sources for re-verification.
4. Remove sources that are no longer accessible (404, removed).
5. Update query templates if search landscape has changed.

---

## Collection Checklist

When performing web research during EVOLVE phase:

- [ ] Define specific research question before searching.
- [ ] Use at least 2 different query templates.
- [ ] Record tier classification for every source found.
- [ ] Cross-reference T3/T4 findings with T1/T2 sources.
- [ ] Check publication date (reject sources older than staleness threshold without explicit justification).
- [ ] Summarize findings with source attribution.
- [ ] Classify proposed changes by Safety Level before applying.
- [ ] Verify change budget allows proposed modifications.

---

## Currently Tracked Sources

### T1: Official Primary

| Source | URL pattern | Last verified | Notes |
|--------|-------------|---------------|-------|
| Anthropic Docs | docs.anthropic.com | Initial | Claude Code documentation |
| MCP Specification | modelcontextprotocol.io | Initial | Protocol definition |

### T2: Official Secondary

| Source | URL pattern | Last verified | Notes |
|--------|-------------|---------------|-------|
| Anthropic Blog | anthropic.com/blog | Initial | Feature announcements, best practices |
| Claude Code GitHub | github.com/anthropics/claude-code | Initial | Examples, issues, discussions |

### T3: Community

| Source | URL pattern | Last verified | Notes |
|--------|-------------|---------------|-------|
| r/ClaudeAI | reddit.com/r/ClaudeAI | Initial | Community practices, tips |
| Developer blogs | Various | Initial | Practitioner experiences |

### T4: Academic

| Source | URL pattern | Last verified | Notes |
|--------|-------------|---------------|-------|
| Multi-agent papers | arxiv.org | Initial | Theoretical foundations |
| Prompt engineering research | Various | Initial | Design pattern research |
