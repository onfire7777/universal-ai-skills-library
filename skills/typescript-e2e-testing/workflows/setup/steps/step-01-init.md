---
name: 'step-01-init'
description: 'Initialize setup workflow — set output path, load knowledge, detect continuation'
nextStepFile: './step-02-project-analysis.md'
referenceFiles:
  - 'references/common/knowledge.md'
  - 'references/common/nestjs-setup.md'
  - 'references/common/rules.md'
---

# Step 1: Initialize E2E Setup

## STEP GOAL

Set up the setup session: identify the project, set the output path for the setup checklist, load knowledge base, and check for existing output to resume.

## REFERENCE LOADING

Before starting, load and read:
- `references/common/knowledge.md` - Understand E2E fundamentals and test pyramid
- `references/common/nestjs-setup.md` - Project structure and Jest configuration
- `references/common/rules.md` - Mandatory patterns (GWT, timeouts, isolation)

## EXECUTION

### 1. Ask the User

Ask the user:
- **Project path** to set up E2E testing for
- **Output path** for the setup checklist (suggest a default: `./e2e-setup-checklist-{{date}}.md`)
- **Or provide path to an existing checklist** to resume

### 2. Check for Existing Output

If the user provides a path to an existing checklist file:
- Read the file
- Parse the YAML frontmatter
- If `stepsCompleted` is non-empty → **STOP and load `step-01b-continue.md`**

### 3. Fresh Workflow Setup

If starting fresh:
1. Copy the template from `templates/output-template.md`
2. Fill in the frontmatter:
   - `projectPath`: the project path
   - `outputPath`: the chosen output path
   - `date`: current date
3. Write the initialized checklist to the output path

### 4. Identify Technology References

Based on project, identify which technology-specific references to load:
- **Kafka**: `references/kafka/docker-setup.md`, `references/kafka/knowledge.md`
- **PostgreSQL**: `references/postgres/knowledge.md`
- **MongoDB**: `references/mongodb/docker-setup.md`, `references/mongodb/knowledge.md`
- **Redis**: `references/redis/docker-setup.md`, `references/redis/knowledge.md`
- **API**: `references/api/knowledge.md`

### 5. Append to Output

Append to the output document:

```markdown
## Step 1: Session Setup

**Project**: {{projectPath}}
**References Loaded**: {{list of loaded references}}
```

## FRONTMATTER UPDATE

Update the output document frontmatter:
- Add `1` to `stepsCompleted`

## PRESENT TO USER

Show the user:
- Confirmation of project path and output path
- References loaded

Then ask: **[C] Continue to Step 2: Project Analysis**

## NEXT STEP

After user confirms `[C]`, load `step-02-project-analysis.md`.
