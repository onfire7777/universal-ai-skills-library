---
name: skill-deployer
description: Deploy and manage skills in Manus project UI via the gRPC-web API. Packages skills as zips, uploads/deletes/syncs per-project skill sets, and curates domain-relevant skills using LLM analysis. Use when needing to make filesystem skills visible in the Manus UI sidebar or synchronize skill sets across projects.
license: Unspecified
---
# Skill Deployer

Deploy skills from the sandbox filesystem into the **Manus project UI** — making them visible in the project sidebar, available as slash commands, and persistent across tasks. This is the only way to make filesystem skills (`/home/ubuntu/skills/`) appear in the Manus interface.

## How It Works

Skills in `/home/ubuntu/skills/` are invisible in the Manus UI by default. This skill bridges that gap:

1. **Packages** each skill directory into a `.zip` (SKILL.md + scripts/references at root)
2. **Authenticates** with the Manus gRPC-web API using a JWT from the browser session
3. **Uploads** zips to specific projects via `UploadProjectSkill` endpoint
4. Skills then appear in the **project sidebar under "Skills"** and as **slash commands** in new tasks

## Prerequisites

Before deploying, two things are needed:

**1. JWT Token** — Extract from the browser while user is logged into manus.im:
```javascript
// Run via browser_console_exec tool:
(document.cookie.split(';').find(c => c.trim().startsWith('session_id=')) || '').split('=').slice(1).join('=') || 'TOKEN_NOT_FOUND'
```
Save to file: `echo "TOKEN" > /tmp/manus_token.txt`

**2. Project UIDs** — Found in the URL when viewing a project:
```
https://manus.im/app/project/<PROJECT_UID>
```
Extract from sidebar via `browser_console_exec` or by clicking each project.

## Commands

### Package Skills (No Upload)

Create zip files from all skills in the skills directory:
```bash
python3 /home/ubuntu/skills/skill-deployer/scripts/deploy_skills.py \
  --token "$(cat /tmp/manus_token.txt)" --package
```

### List Skills in a Project

See what's currently in a project's UI:
```bash
python3 /home/ubuntu/skills/skill-deployer/scripts/deploy_skills.py \
  --token "$(cat /tmp/manus_token.txt)" --project-uid UID --list
```

### Upload a Single Skill to UI

Package and upload one skill by name:
```bash
python3 /home/ubuntu/skills/skill-deployer/scripts/deploy_skills.py \
  --token "$(cat /tmp/manus_token.txt)" --project-uid UID --upload-skill skill-name
```

### Upload All Skills from Zip Directory

```bash
python3 /home/ubuntu/skills/skill-deployer/scripts/deploy_skills.py \
  --token "$(cat /tmp/manus_token.txt)" --project-uid UID --upload-dir /tmp/skill_zips
```

### Delete a Skill from UI

```bash
python3 /home/ubuntu/skills/skill-deployer/scripts/deploy_skills.py \
  --token "$(cat /tmp/manus_token.txt)" --project-uid UID --delete-name skill-name
```

### Sync Project to Exact Skill Set

Deletes skills not in the list, uploads missing ones:
```bash
python3 /home/ubuntu/skills/skill-deployer/scripts/deploy_skills.py \
  --token "$(cat /tmp/manus_token.txt)" --project-uid UID \
  --sync skill-a skill-b skill-c
```

### Multi-Project Deployment with LLM Curation

For deploying curated, domain-relevant skills to multiple projects at once.

**Step A — Create projects file** (`projects.json`):
```json
{
  "projects": [
    {
      "name": "My Web App",
      "uid": "abc123",
      "context": "Full-stack React + Node web app with auth and database",
      "tasks": ["frontend dev", "API design", "database queries"],
      "max_skills": 300
    }
  ]
}
```

**Step B — Run LLM curation** (selects universal + domain-specific skills):
```bash
python3 /home/ubuntu/skills/skill-deployer/scripts/curate_for_project.py \
  --projects projects.json --output deployment_plan.json
```

**Step C — Preview the plan:**
```bash
python3 /home/ubuntu/skills/skill-deployer/scripts/deploy_skills.py \
  --token "$(cat /tmp/manus_token.txt)" --plan deployment_plan.json --dry-run
```

**Step D — Deploy to all projects:**
```bash
python3 /home/ubuntu/skills/skill-deployer/scripts/deploy_skills.py \
  --token "$(cat /tmp/manus_token.txt)" --plan deployment_plan.json
```

## Options

| Flag | Default | Description |
|---|---|---|
| `--skills-dir` | `/home/ubuntu/skills` | Source skills directory |
| `--zip-dir` | `/tmp/skill_zips` | Directory for packaged zips |
| `--rate-limit` | `0.15` | Seconds between API calls |
| `--max-retries` | `3` | Retries on 429 rate limit errors |
| `--dry-run` | off | Preview changes without executing |

## Limits and Constraints

- **500 skills per project** — hard limit enforced by Manus
- **Rate limiting** — 429 errors handled with exponential backoff
- **Token expiry** — JWT tokens last ~24 hours; re-extract if 401 errors occur
- **SKILL.md required** — Zips must contain valid SKILL.md with YAML frontmatter at root

## Troubleshooting

| Error | Fix |
|---|---|
| 401 Unauthorized | Token expired. Re-extract from browser. |
| 429 Rate Limited | Increase `--rate-limit` to 0.3+. |
| 500 SKILL_MD_INVALID_FORMAT | Fix SKILL.md YAML frontmatter in the source skill. |
| Project limit reached | Max 500. Use curation to select only the best skills. |
| Skills not showing in UI | Refresh the project page in browser. Skills appear in sidebar. |

## API Reference

Full endpoint documentation in `references/manus-api.md`.
