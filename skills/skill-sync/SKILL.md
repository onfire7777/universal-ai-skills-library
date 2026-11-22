---
name: skill-sync
description: Sync and install skills from a GitHub skills repository into a local or global AI skill root. Use when the user asks to install skills, update skills, sync skills from GitHub, load a skill library, or ensure cross-AI skill roots are up to date.
---

# Skill Sync

Sync all skills from a GitHub repository into `/home/ubuntu/skills/` with version tracking, validation, and change detection.

## Quick Start

Run the script to preview the default repository, `onfire7777/universal-ai-skills-library`.

```bash
python3 /home/ubuntu/skills/skill-sync/scripts/sync_skills.py --dry-run
```

Show the user the repository, full commit SHA, target directory, and skills that would be installed or updated. Let the user inspect the selected changes and wait for approval of that exact sync. Then rerun with the same repository, target, and filters, replacing `--dry-run` with `--apply <commit-sha>`. Review a new preview if the repository has moved on.

## Commands

```bash
# Preview selected skills from a different repository
python3 /home/ubuntu/skills/skill-sync/scripts/sync_skills.py --repo user/repo-name --filter "security,audit"

# Preview replacing locally modified skills
python3 /home/ubuntu/skills/skill-sync/scripts/sync_skills.py --force

# Validate currently installed skills
python3 /home/ubuntu/skills/skill-sync/scripts/sync_skills.py --validate-only
```

## How It Works

1. **Fetch** the repository's default branch with Git and preview its exact commit
2. **Discover** all directories containing a valid `SKILL.md` (supports both `repo/skill-name/` and `repo/skills/skill-name/` layouts)
3. **Compare** each skill's content hash against the last sync state
4. **Apply** only when the supplied commit matches the previewed source. Install new skills, update changed skills, and skip unchanged ones.
5. **Validate** that each skill has SKILL.md frontmatter (`name` + `description`)
6. **Save state** after applying, with the commit hash and per-skill content hashes

## Change Detection

The script tracks content hashes to avoid unnecessary overwrites:

- **Source unchanged, dest unchanged** - skipped (already up-to-date)
- **Source changed, dest unchanged** - updated from repo
- **Source unchanged, dest changed** - skipped (local changes preserved)
- **Both changed** - skipped unless `--force` is used

## State File

Sync state is stored at `~/.skill-sync-cache/sync_state.json`. Contains last sync timestamp, commit hash, and per-skill content hashes. Delete this file to force a full re-sync.

## Configuration

Use Git and Python 3.12 or newer. The preview may populate the Git cache. It leaves the skills directory and installation state alone.

Paths can be overridden via environment variables or CLI flags:

| Setting | Default | Override |
|---|---|---|
| Skills directory | `/home/ubuntu/skills` | `SKILL_ROUTER_SKILLS_DIR` or `--skills-dir` |
| Cache directory | `~/.skill-sync-cache` | `SKILL_ROUTER_SYNC_CACHE` or `--cache-dir` |

## After Syncing

Skills are immediately available to compatible AI clients after syncing. Agents route skills from the target skill root using the `name` and `description` in each `SKILL.md` frontmatter.
