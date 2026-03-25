---
name: internet-skill-finder
description: Search and recommend Agent Skills from verified GitHub repositories. Use when users ask to find, discover, search for, or recommend skills/plugins for specific tasks, domains, or workflows.
---

# Internet Skill Finder

Search 7 verified GitHub repositories for Agent Skills.

## Workflow

### 1. Fetch Skill List

```bash
python3 /home/ubuntu/skills/internet-skill-finder/scripts/fetch_skills.py --search "keyword"
```

Options: `--list` (all skills), `--online` (real-time fetch), `--json` (structured output)

### 2. Deep Dive (if needed)

```bash
python3 /home/ubuntu/skills/internet-skill-finder/scripts/fetch_skills.py --deep-dive REPO SKILL
```

### 3. Present Results

When using cached data, prepend:

> ℹ️ Using cached data. Enable GitHub Connector for real-time results.

Format each match:

```markdown
### [Skill Name]
**Source**: [Repository] | ⭐ [Stars]
**Description**: [From SKILL.md]
👉 **[Import](import_url)**
```

### 4. No Matches

Suggest `/skill-creator` for custom skill creation.

## Data Access

Script auto-detects and uses best method:

| Priority | Method | Rate Limit | Behavior |
|----------|--------|------------|----------|
| 1 | GitHub Connector | 15000/hr | Auto real-time |
| 2 | Offline cache | Unlimited | Fallback |
| 3 | `GITHUB_TOKEN` env | 5000/hr | With `--online` |

JSON output includes `"using_cache": true/false` to indicate data source.

When cache is used, inform user: Enable GitHub Connector for real-time results.

## Sources

7 repositories: anthropics/skills, obra/superpowers, vercel-labs/agent-skills, K-Dense-AI/claude-scientific-skills, ComposioHQ/awesome-claude-skills, travisvn/awesome-claude-skills, BehiSecc/awesome-claude-skills
