---
name: ultimate-skill-creator
description: >
  The definitive end-to-end pipeline for creating the highest quality portable AI skills. Orchestrates
  five development skills and a local verification stage in sequence to produce
  elite, production-grade skills. Use when the user explicitly requests the full pipeline, says
  "ultimate skill", "ultimate skill creator", "ultimate skill creation", "best skill", "best skill
  creator", "full pipeline", "full skill pipeline", or specifically asks for maximum quality skill
  creation with all optimization stages. For simple skill creation without the full pipeline, use
  skill-creator instead.
---

# Ultimate Skill Creator

Execute a comprehensive six-stage pipeline to create, optimize, triage, connect, debug, and finalize a production-grade portable AI skill. Run each development skill through its complete workflow, then verify the finished local skill.

## Pipeline Overview

| Stage | Skill | Purpose | Exit Criteria |
|---|---|---|---|
| 1 | `skill-creator` | Full 6-step creation process | Passes `quick_validate.py`, all scripts execute |
| 2 | `prompt-engineer` | Deep multi-model SKILL.md optimization | Optimized SKILL.md passes validation |
| 3 | `master-skill-orchestrator` | Skill Triage Protocol for domain analysis | Triage documented in `/tmp/triage_$SKILL_NAME.md` |
| 4 | `skill-connection-map` | Connection Discovery Protocol | Connections documented in `/tmp/connections_$SKILL_NAME.md` |
| 5 | `skill-debugger` | Deep dual-model debug + fix cycle | No critical/high issues after max 3 cycles |
| 6 | Finalize | Verify the local skill and optionally publish | Local skill verified |

## Pre-Flight Check

Run the pre-flight script before starting the pipeline:

```bash
python3 /home/ubuntu/skills/ultimate-skill-creator/scripts/preflight_check.py
```

This verifies the five development skills and checks for an API key. Verify access to a remote only if the user authorizes publishing in Stage 6. If it fails, fix the reported errors before proceeding. If it passes with warnings, the pipeline can proceed with limited functionality.

## Operational Rules

1. **Read each skill's SKILL.md fresh** at the start of each stage. Do not rely on memory.
2. **Complete each stage fully** before advancing. Do not skip steps within a skill's workflow.
3. **Announce each stage transition**: "Starting Stage N: [skill-name] — [purpose]."
4. **Stage transition checkpoint** — Before starting each new stage, confirm the stage number, skill name, and purpose from the Pipeline Overview table above. Do not re-read the entire SKILL.md unless uncertain about requirements.
5. **Gather creation requirements during Stage 1.** Obtain explicit publishing authorization in Stage 6 if the user requests publication and has yet to approve the destination and files.
6. **On fatal errors**, inform the user immediately with error context and recovery options.
7. **Clean up temporary files** after each stage, including on failure.
8. **Sequential execution only** — Complete each stage fully before starting the next. Do not run stages concurrently.

---

## Stage 1: skill-creator — Full Creation

**Read**: `/home/ubuntu/skills/skill-creator/SKILL.md`

Execute the skill-creator's complete 6-step process:

1. **Understand** — Gather concrete examples of how the skill will be used. Ask the user clarifying questions if the skill's purpose is unclear. Do not proceed until you have a clear sense of the functionality.
2. **Plan** — Identify reusable resources (scripts, references, templates). Determine what goes in each directory.
3. **Initialize** — Run `python3 /home/ubuntu/skills/skill-creator/scripts/init_skill.py $SKILL_NAME`. Clean up example files.
4. **Edit** — Read the skill-creator's reference files (`references/workflows.md`, `references/output-patterns.md`, `references/progressive-disclosure-patterns.md`) for design patterns. Write all scripts, references, templates. Write the SKILL.md with proper frontmatter and comprehensive instructions. Test all scripts.
5. **Validate** — Run `python3 /home/ubuntu/skills/skill-creator/scripts/quick_validate.py $SKILL_NAME`. Fix any errors and re-validate.
6. **Confirm** — Inform the user that Stage 1 is complete. Show the skill structure.

**Exit**: Skill passes `quick_validate.py` and all scripts execute without errors.

---

## Stage 2: prompt-engineer — Deep Optimization

**Read**: `/home/ubuntu/skills/prompt-engineer/SKILL.md`

Execute the prompt-engineer's complete workflow. Replace `$SKILL_NAME` with the actual skill name throughout.

1. **Prepare** — Copy the skill's SKILL.md to a temp file:
   ```bash
   cp /home/ubuntu/skills/$SKILL_NAME/SKILL.md /tmp/pe_input_$SKILL_NAME.txt
   ```
2. **Verify flags** — Check the script's supported arguments:
   ```bash
   python3 /home/ubuntu/skills/prompt-engineer/scripts/optimize_prompt.py --help
   ```
   Adjust the flags in step 3 if any are not supported.
3. **Optimize** — Run at deep depth:
   ```bash
   python3 /home/ubuntu/skills/prompt-engineer/scripts/optimize_prompt.py \
     --file /tmp/pe_input_$SKILL_NAME.txt \
     --depth deep --show-analysis --show-comparison \
     -o /tmp/pe_output_$SKILL_NAME.md
   ```
4. **Synthesize** — Read the optimized output. Rewrite the SKILL.md incorporating the best improvements. Preserve all script paths, command references, and technical accuracy. Do not blindly replace — merge intelligently.
5. **Validate** — Re-run `quick_validate.py` to confirm the optimized SKILL.md passes.
6. **Clean up** — Remove `/tmp/pe_input_$SKILL_NAME.txt` and `/tmp/pe_output_$SKILL_NAME.md`, even if earlier steps failed.

**Exit**: SKILL.md is rewritten with prompt-engineer improvements and passes validation.

---

## Stage 3: master-skill-orchestrator — Full Triage

**Read**: `/home/ubuntu/skills/master-skill-orchestrator/SKILL.md`

Execute the orchestrator's complete Skill Triage Protocol:

1. **Domain Classification** — Identify ALL domains the new skill touches using the orchestrator's domain table.
2. **Complexity Assessment** — Classify as simple/moderate/complex/wicked.
3. **Phase Classification** — Determine which phases the skill covers (Understand → Reflect).
4. **Skill Selection** — Using the Priority Stack, Common Task Patterns, Prerequisite Chains, and Emergent Combinations, identify additional installed skills whose knowledge should inform the new skill.
5. **Apply insights** — If relevant skills or patterns are found, read those skills and incorporate their knowledge into the SKILL.md or reference files.
6. **Anti-Pattern Check** — Verify no Skill Overload, Tunnel Vision, Cargo Culting, or Skipping.
7. **Document** — Save the triage analysis to `/tmp/triage_$SKILL_NAME.md` with domain classification, complexity level, phase classification, and any insights applied.

**Exit**: Triage document exists at `/tmp/triage_$SKILL_NAME.md` with all required sections.

---

## Stage 4: skill-connection-map — Connection Discovery

**Read**: `/home/ubuntu/skills/skill-connection-map/SKILL.md`

Execute the Connection Discovery Protocol:

1. **Name the problem domain** — What field is the new skill obviously in?
2. **Name three distant domains** — What fields have nothing to do with this skill?
3. **Find one skill per distant domain** — What installed skill has a structural analogy?
4. **Apply distant frameworks** — Translate the distant skill's concepts into the new skill's domain.
5. **Evaluate all seven connection types** — Structural Analogy, Inverse/Complement, Prerequisite Chain, Cross-Domain Transfer, Emergent Combination, Tension/Dialectic, Recursive/Fractal.
6. **Incorporate discoveries** — Apply any valuable cross-domain insights to the SKILL.md or references.
7. **Document** — Save the connection analysis to `/tmp/connections_$SKILL_NAME.md` listing all seven connection types evaluated and any improvements applied.

**Exit**: Connection document exists at `/tmp/connections_$SKILL_NAME.md` with all seven types evaluated.

---

## Stage 5: skill-debugger — Full Debug Cycle

**Read**: `/home/ubuntu/skills/skill-debugger/SKILL.md`

This is the most time-intensive stage (2-5 minutes per run, potentially multiple cycles). Set expectations accordingly.

Execute the skill-debugger's complete workflow:

1. **Run deep analysis**:
   ```bash
   python3 /home/ubuntu/skills/skill-debugger/scripts/debug_skill.py $SKILL_NAME --deep
   ```
   Wait for completion.
2. **Read the full report** — Read `/home/ubuntu/skills/$SKILL_NAME/DEBUG_REPORT.md`. Understand every finding.
3. **Apply ALL fixes** — Fix in severity order (critical → low). Fully rewrite affected files — no partial edits.
4. **Re-run debugger** — Verify fixes. If new issues emerge, fix those too. **Maximum 3 debug cycles.** If critical/high issues persist after 3 cycles, list all unresolved issues, inform the user, and ask whether to proceed to Stage 6 or abort the pipeline. Do not proceed silently.
5. **Clean up** — Remove `/home/ubuntu/skills/$SKILL_NAME/DEBUG_REPORT.md` and `/home/ubuntu/skills/$SKILL_NAME/.debug_raw.json`.
6. **Final validation** — Run `quick_validate.py` one last time.

**Exit**: No critical/high severity issues (or user confirmed proceed after 3 cycles), passes `quick_validate.py`.

---

## Stage 6: Verify the Local Skill

1. Confirm the finished skill is present in the local skills root and rerun the local validator. The skill created in Stage 1 is already local.
2. If the user requests publication, confirm the exact repository, target branch, destination path, and files to publish. Show the proposed diff or file list before committing or pushing and wait for approval of that scope. Review the selected files for secrets and private material. Obtain fresh approval if the destination or file set changes.
3. After approval, use an isolated checkout and stage only the approved skill files. Preserve existing work in the destination and verify the staged diff before committing and pushing to the approved branch.
4. Remove the temporary triage and connection reports created by this run.

**Exit:** The local skill passes validation. If publication was approved, verify the result at the approved destination and report the commit.

---

## Completion Deliverables

After all six stages, deliver to the user:

1. **Attachment**: Final `SKILL.md` file
2. **Summary table**:

| Stage | Skill | Key Actions | Result |
|---|---|---|---|
| 1 | skill-creator | [what was created] | [validation result] |
| 2 | prompt-engineer | [improvements applied] | [improvement score] |
| 3 | master-skill-orchestrator | [domains/skills identified] | [insights applied] |
| 4 | skill-connection-map | [connections discovered] | [improvements made] |
| 5 | skill-debugger | [issues found/fixed] | [final health] |
| 6 | Finalize | [local validation and optional publication] | [local path or approved commit] |

3. **Key metrics**: Prompt-engineer improvement score, debugger findings resolved, files created, and the commit SHA if published
4. **Reflection**: Brief summary of what was learned during the pipeline — what worked well, what was challenging, and any insights for future skill creation

## Constraints

- NEVER skip or combine stages. Each skill runs fully and independently.
- NEVER partially edit files during debug fixes. Fully rewrite affected files.
- Maximum 3 debug cycles in Stage 5 to prevent infinite loops.
- Sequential execution only — no concurrent stages.
- Clean up all temporary files after each stage, including on failure.
- Preserve all original script paths, commands, and technical details exactly.
- Verify $SKILL_NAME is set correctly before any destructive operations.
