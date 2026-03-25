---
name: multi-repo-git-ops
description: Manages git operations (branching, committing, pushing) across multi-repo systems with git submodules. Use this skill whenever the user wants to create a feature branch, commit changes, push code, sync submodules, or manage branches across the parent repo and service repos. Also use when implementing BMAD stories that touch one or more services, or when the user asks about git workflow in a multi-repo project. Triggers on phrases like "create branch", "commit changes", "push to remote", "sync submodules", "start working on story", "update service branch", or any git operation that spans parent and submodule repos.
---

# Multi-Repo Git Operations

This skill handles git operations in **multi-repo systems** that use git submodules. The parent repo orchestrates multiple service repos, each an independent git repository with its own branches, history, and CI/CD.

Understanding this structure is essential — git operations here always involve deciding **which repos** are affected and operating in each one correctly.

## Discovery Protocol

Before performing any git operations, discover the project structure dynamically. Never assume service names, branch mappings, or directory layout — always read from `.gitmodules`.

### Step 1: Identify submodules and their paths

```bash
# List all submodules with their paths
git config -f .gitmodules --get-regexp '\.path$'
```

### Step 2: Identify each submodule's tracked branch

```bash
# List all submodule branch mappings
git config -f .gitmodules --get-regexp '\.branch$'

# Get a specific service's default branch
git config -f .gitmodules submodule."{submodule-path}".branch
```

If a submodule has no branch configured in `.gitmodules`, check the remote:
```bash
cd {submodule-path}
git remote show origin | grep 'HEAD branch'
```

### Step 3: Detect project methodology

- **BMAD**: Check if `_bmad-output/implementation-artifacts/` directory exists
- If BMAD is detected, the BMAD Integration section applies

### Step 4: Cache discovered info for the session

After discovery, the agent knows:
- All submodule names and paths
- Each submodule's default branch
- Whether BMAD is in use
- The submodule directory prefix (commonly `services/`)

Use this cached information for all subsequent operations in the session.

## Repository Structure

A typical multi-repo project looks like this:

```
{parent-repo}/                       (parent repo — planning & orchestration only)
├── .gitmodules                      (submodule definitions with tracked branches)
├── services/                        (submodule repos — directory name may vary)
│   ├── auth-service/                (branch: develop)
│   ├── booking-service/             (branch: dev)
│   ├── config-repo/                 (branch: main)
│   └── ...
└── [project methodology dirs]       (e.g., _bmad-output/ if using BMAD)
```

**Key principle**: The parent repo tracks **commit pointers** to each submodule, not the submodule code itself. When you change code in a service, you commit in the service repo, then update the parent's pointer.

## Branch Conventions

### Default Branches by Service

Services may track different default branches — always look up the correct one before branching:

```bash
# Get a service's default branch (always use this — never assume)
git config -f .gitmodules submodule."{submodule-path}".branch
```

Different services in the same project often use different default branches (e.g., `develop`, `dev`, `main`). The only reliable source of truth is `.gitmodules`.

### Feature Branch Naming

**For ticket-based work** (JIRA/project tracker):
```
{type}/{TICKET-ID}-{short-description}
```
Examples: `feat/PROJ-1234-add-search-by-date`, `bugfix/PROJ-5678-fix-null-pointer`

**For BMAD story work** (derived from story file name):
```
feat/{story-key}
```
The story key comes from the story file name. A story file named `1-2-user-authentication.md` (representing Epic 1, Story 2) produces branch name `feat/1-2-user-authentication`.

**Branch type prefixes**: `feat/`, `fix/`, `bugfix/`, `chore/`, `base/`

## BMAD Integration (If Applicable)

This section applies only if the project uses the BMAD methodology. Detect this by checking for the `_bmad-output/implementation-artifacts/` directory.

BMAD manages work through story files and sprint status tracking. This section explains how git operations map to the BMAD lifecycle.

### How Stories Reference Services

BMAD story files live at `_bmad-output/implementation-artifacts/{story-key}.md`. Stories reference affected services in several ways — check all of these:

1. **Service tags in story title or body**: `[auth-service]`, `[scheduler-service]`
2. **Dev Notes section**: Lists "Source tree components to touch" with service paths
3. **Project Structure Notes**: References paths like `services/{service-name}/src/...`
4. **Tasks/Subtasks**: Individual tasks may reference different services

When a story doesn't explicitly tag services, look at the file paths mentioned in Dev Notes and Tasks to determine which submodule directories are involved.

### BMAD Status and Git Operations Mapping

BMAD tracks story status in `_bmad-output/implementation-artifacts/sprint-status.yaml`. Each status transition has corresponding git operations:

| Story Status Change | Git Operations Required |
|---|---|
| `backlog` → `ready-for-dev` | No git ops (story file creation only) |
| `ready-for-dev` → `in-progress` | Create feature branches in affected services |
| `in-progress` (ongoing work) | Commit changes in service repos |
| `in-progress` → `review` | Push all service branches, ensure clean state |
| `review` → `in-progress` (fixes) | Continue on same branches, commit fixes |
| `review` → `done` | PRs merged, update parent submodule pointers |

### Starting a BMAD Story (Git Setup)

When beginning work on a BMAD story, perform these git operations:

**Step 1: Read the story and identify services**
```bash
# Story files follow pattern: {epic_num}-{story_num}-{slug}.md
# Example: _bmad-output/implementation-artifacts/1-2-user-authentication.md
```

Read the story file. Extract affected services from:
- Service tags like `[auth-service]`
- File paths in Dev Notes (e.g., `services/{service-name}/src/...`)
- Task descriptions referencing specific services

**Step 2: Derive the branch name from the story key**
```bash
# Story file: 1-2-user-authentication.md → branch: feat/1-2-user-authentication
STORY_KEY="1-2-user-authentication"  # from the story file name (without .md)
BRANCH_NAME="feat/${STORY_KEY}"
```

**Step 3: Create branches in all affected services**
```bash
# For each affected service (replace with actual discovered service names):
for SERVICE in auth-service scheduler-service; do
  DEFAULT_BRANCH=$(git config -f .gitmodules submodule."services/$SERVICE".branch)
  cd services/$SERVICE
  git checkout $DEFAULT_BRANCH
  git pull origin $DEFAULT_BRANCH
  git checkout -b $BRANCH_NAME
  git push -u origin $BRANCH_NAME
  cd ../..
done
```

Use the **same branch name** across all affected services for traceability. This makes it easy to find all changes related to a story across the system.

**Step 4: Verify setup**
```bash
# Confirm all affected services are on the correct branch
git submodule foreach --quiet 'BRANCH=$(git branch --show-current); echo "$(basename $(pwd)): $BRANCH"'
```

### During Story Implementation

**Committing changes** — always commit inside the service directory:
```bash
cd services/{service-name}
git add {specific-files}
git commit -m "feat({scope}): add validation for booking dates

Story: {story-key}"
```

**Working across multiple services** — commit in each separately:
```bash
# Service A
cd services/auth-service
git add src/auth/jwt.service.ts src/auth/jwt.service.spec.ts
git commit -m "feat(auth): add JWT validation endpoint

Story: 1-2-user-authentication"

# Service B
cd ../bff-service
git add src/middleware/auth.middleware.ts
git commit -m "feat(middleware): add JWT auth middleware

Story: 1-2-user-authentication"
```

All commits must follow the **Commit Message Format** section below — this is critical for release-please to generate correct changelogs and version bumps.

### Completing a Story (Push & PR)

When story is ready for review:

**Step 1: Run quality checks in each affected service**
```bash
cd services/{service-name}
npm run lint
npm run format
npm run typecheck
npm test
```

**Step 2: Push each service's feature branch**
```bash
# Push services first — always before parent
cd services/{service-name}
git push origin feat/{story-key}
```

**Step 3: Create PR per service**
Each service gets its own PR: `feat/{story-key}` → service's default branch.

PR title format: `feat({story-key}): {story title}`
PR body should reference the story: `Story: {epic_num}.{story_num} - {title}`

**Step 4: After PRs are merged, update parent**
```bash
# For each merged service:
cd services/{service-name}
git checkout {default-branch}
git pull origin {default-branch}
cd ../..

# Update parent pointers (use actual discovered service paths)
git add services/auth-service services/bff-service
git commit -m "chore: update submodule pointers after story {story-key}

Services updated:
- auth-service
- bff-service"
```

## Commit Message Format (Release-Please Compatible)

If the project uses [release-please](https://github.com/googleapis/release-please) for automated versioning and changelog generation, every commit message must follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification so release-please can correctly determine version bumps and produce meaningful changelogs.

Even without release-please, Conventional Commits is the recommended format for multi-repo systems because it produces clean, parseable git history across many repositories.

### No AI Co-Author Trailers

**NEVER add `Co-Authored-By` trailers for AI agents in commit messages.** This means no lines like:
```
Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
```
Commit messages must only contain the conventional commit structure described below. AI attribution in commits pollutes changelogs, adds noise to git history, and provides no value. This rule applies to all commits — feature, fix, chore, or otherwise.

### Structure

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

**Header rules:**
- `type` — required (see type table below)
- `scope` — optional, a noun describing the affected code section in parentheses
- `description` — required, immediately after `: `
  - Use imperative, present tense: "add" not "added" or "adds"
  - Do not capitalize the first letter
  - Do not end with a period
  - Keep under ~100 characters

**Body** — optional, separated from header by a blank line. Free-form, explains the "what" and "why".

**Footer(s)** — optional, separated from body by a blank line. Used for breaking changes, story references, and other metadata.

### Commit Types and Version Bumps

Release-please only creates a release when it detects **releasable units** — commits with types that map to a version bump. Choosing the right type matters because it directly controls what version gets released and what appears in the changelog.

#### Releasable types (trigger a release)

| Type | SemVer Bump | When to use |
|------|-------------|-------------|
| `feat` | **Minor** (0.x.0) | New capability, endpoint, feature, UI component |
| `fix` | **Patch** (0.0.x) | Bug fix for existing functionality |
| `deps` | **Patch** (0.0.x) | Dependency updates (package.json, lock files) |

#### Non-releasable types (won't trigger a release alone)

| Type | When to use |
|------|-------------|
| `refactor` | Code restructuring without behavior change |
| `perf` | Performance improvement (special refactoring) |
| `test` | Adding or correcting tests |
| `docs` | Documentation changes only |
| `style` | Code style (formatting, whitespace, semicolons) |
| `chore` | Maintenance tasks (.gitignore, configs, tooling) |
| `build` | Build system, dependencies, project version |
| `ci` | CI/CD pipeline and workflow changes |

If you only commit non-releasable types since the last release, no release PR will be created. Be intentional — don't use `chore:` when the change is actually a `feat:` or `fix:`.

#### Choosing the right type

```
New capability or feature?           → feat
Fixing a bug?                        → fix
Dependency update?                   → deps
Performance improvement?             → perf
Code restructuring, no behavior Δ?   → refactor
Documentation only?                  → docs
Tests only?                          → test
CI/CD or build changes?              → ci / build
Everything else (configs, tooling)?  → chore
```

### Breaking Changes (Major Bump)

Any commit type with a breaking change triggers a **major** version bump. Two formats:

**Option A — `!` in the header:**
```
feat(api)!: redesign authentication flow

BREAKING CHANGE: the /auth/login endpoint now requires OAuth2 tokens
instead of API keys. All existing integrations must migrate.
```

**Option B — footer only:**
```
refactor(database): migrate from MongoDB to PostgreSQL

BREAKING CHANGE: all database connection strings must be updated.
The MongoDB driver is no longer included.
```

Both `BREAKING CHANGE` (space) and `BREAKING-CHANGE` (hyphen) are recognized in footers.

When marking a breaking change, always include a `BREAKING CHANGE:` footer explaining the migration impact — even when using the `!` format. The footer is what consumers read to understand what they need to change.

### Scopes

Scopes describe which **code section** is affected — they help organize changelogs and, in monorepo setups, route changes to the correct package.

**Good scopes** (describe code areas):
```
feat(auth): add JWT validation
fix(scheduler): correct cron expression parsing
refactor(booking): extract date utility
```

**Bad scopes** (these are anti-patterns):
```
fix(PROJ-1234): ...          ← ticket IDs are not scopes
feat(john): ...               ← people are not scopes
```

Use the **domain or module name** as the scope (e.g., `auth`, `booking`, `scheduler`, `middleware`, `config`, `migration`).

### Story References

If using BMAD, include the story key in the commit **body or footer**, not in the scope:

```
feat(auth): add JWT validation endpoint

Implement token validation with RS256 signing.

Story: 1-2-user-authentication
```

### Special Release-Please Footers

| Footer | Purpose |
|--------|---------|
| `BREAKING CHANGE: <desc>` | Triggers major version bump |
| `Release-As: x.x.x` | Force a specific version number |

**`Release-As` example** — useful for initial releases or coordinated version jumps:
```
chore: prepare v3.0.0 release

Release-As: 3.0.0
```

### Complete Examples

**Simple feature** (minor bump):
```
feat(booking): add search by date range endpoint
```

**Bug fix with context** (patch bump):
```
fix(auth): resolve token refresh race condition

The refresh token was being invalidated before the new access token
was issued, causing a brief window where requests would fail.

Story: 2-3-auth-improvements
```

**Breaking change** (major bump):
```
feat(api)!: remove deprecated v1 booking endpoints

BREAKING CHANGE: all /api/v1/bookings/* endpoints are removed.
Consumers must migrate to /api/v2/bookings/* which uses the new
pagination format.

Story: 3-1-api-v2-migration
```

**Dependency update** (patch bump):
```
deps: upgrade @nestjs/core to v11.0.0
```

**Test-only change** (no release):
```
test(scheduler): add unit tests for cron expression parser
```

**Chore** (no release):
```
chore: update .gitignore to exclude coverage reports
```

### Parent Repo Commit Messages

When updating submodule pointers in the parent repo, use `chore:` since these don't represent feature changes in the parent itself:

```
chore: update submodule pointers after story 1-2-user-authentication

Services updated:
- auth-service
- bff-service
```

## General Operations Reference

### Create Feature Branch (Non-BMAD)

For ad-hoc work not tied to a BMAD story:

```bash
DEFAULT_BRANCH=$(git config -f .gitmodules submodule."{submodule-path}".branch)
cd {submodule-path}
git checkout $DEFAULT_BRANCH
git pull origin $DEFAULT_BRANCH
git checkout -b feat/{TICKET-ID}-{description}
git push -u origin feat/{TICKET-ID}-{description}
```

### Sync All Submodules

```bash
# Pull parent and update all submodule pointers
git pull origin main
git submodule update --recursive

# OR: Pull latest from each submodule's tracked branch
git submodule foreach 'git pull origin $(git config -f $toplevel/.gitmodules submodule.$name.branch)'
```

### Check Status Across All Repos

```bash
# Quick status: show only services with changes or non-default branches
git submodule foreach --quiet \
  'STATUS=$(git status --porcelain); BRANCH=$(git branch --show-current); DEFAULT=$(git config -f $toplevel/.gitmodules submodule.$name.branch); if [ -n "$STATUS" ] || [ "$BRANCH" != "$DEFAULT" ]; then echo "$(basename $(pwd)) [$BRANCH]: $([ -n "$STATUS" ] && echo "has changes" || echo "clean, non-default branch")"; fi'
```

### Switch Service Back to Default Branch

```bash
DEFAULT_BRANCH=$(git config -f .gitmodules submodule."{submodule-path}".branch)
cd {submodule-path}
git checkout $DEFAULT_BRANCH
git pull origin $DEFAULT_BRANCH
```

## CI/CD Awareness

Understanding what happens when you push is critical in multi-repo systems. Typical branch-to-deployment mappings:

| Branch | Typical Push Effect |
|---|---|
| `develop` / `dev` | **May auto-deploy to dev environment** — confirm with user first |
| `main` | Stable reference — deployment workflows available |
| `test`, `stage`, `release` | Manual deployment via `workflow_dispatch` |
| `feat/*`, `base/*` | CI runs (lint, test, build) but **no deployment** |

Actual deployment behavior depends on the project's CI/CD configuration. Check if the project has a shared CI/CD workflow repository among its submodules, and confirm with the user which branches trigger deployments.

## Safety Rules

These rules prevent common mistakes in multi-repo systems:

1. **Always push service repos before parent** — pushing parent first creates broken submodule pointers that break other developers' checkouts
2. **Never force-push** to shared default branches (`develop`, `dev`, `main`) — these are shared branches
3. **Confirm with user** before pushing to deployment-triggering branches
4. **Check service's default branch** before creating feature branches — they differ across services, and branching from the wrong base causes merge conflicts
5. **Never use `git add .` or `git add -A`** in service repos — always add specific files to avoid committing secrets or build artifacts
6. **Never commit** `.env`, credentials, secrets, or `node_modules` in any repo
7. **Same branch name across services** for multi-service stories — enables traceability
8. **Update sprint-status.yaml** when git operations change story state (branch created → `in-progress`, pushed for review → `review`) — only if using BMAD
9. **Never add `Co-Authored-By` trailers for AI agents** in commit messages — no Claude, Copilot, or any AI attribution. These pollute changelogs and git history. Commit messages must contain only the conventional commit structure (type, scope, description, body, footers for breaking changes/story refs)
10. **All commit messages must be Conventional Commits format** — release-please depends on this for automated versioning and changelog generation. A malformed commit message can cause missed releases or incorrect version bumps
