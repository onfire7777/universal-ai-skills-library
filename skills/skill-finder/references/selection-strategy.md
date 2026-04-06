# Skill Selection Strategy

## When to Use Multiple Skills

Use **1 skill** when the task is narrow and domain-specific (e.g., "create a greeting card").
Use **2-3 skills** when the task spans multiple domains (e.g., "research a company and create a presentation").
Use **4+ skills** rarely — only for complex multi-phase projects. Read them sequentially, not all at once.

## Skill Combination Patterns

| Pattern | Example | How to Combine |
|---|---|---|
| Pipeline | Research → Write → Present | Run skill 1 fully, then skill 2, then skill 3 |
| Parallel expertise | Code audit + security review | Read both, merge their checklists |
| Primary + support | Main task skill + quality skill | Follow primary workflow, apply support skill's standards |

## Handling Conflicts Between Skills

When two skills give contradictory instructions:
1. Prefer the skill with higher relevance score
2. Prefer the skill with scripts (more concrete)
3. Prefer the more specific skill over the more general one

## Score Thresholds

- **5.0+**: Strong match — use the skill confidently
- **3.0-5.0**: Partial match — the skill may be useful but consider rephrasing the query first
- **Below 3.0**: Weak match — unlikely to be relevant; proceed with base Manus capabilities

## When No Skills Match

If `find_skills.py` returns no results or all scores are below 3.0 after retries:
1. Try rephrasing the query with different keywords
2. Try breaking the task into sub-tasks and searching each
3. Consider using `internet-skill-finder` to search GitHub for new skills to install
4. Proceed without a skill — Manus's base capabilities handle most tasks
