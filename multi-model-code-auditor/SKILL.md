---
name: multi-model-code-auditor
description: Run a comprehensive security, privacy, and bug audit on any codebase using the absolute best frontier AI models via OpenRouter API. Dynamically discovers the latest models, uses model-specific elite prompt engineering, merges findings by consensus, implements fixes, and runs a verification re-audit. Use when the user requests a code audit, security review, vulnerability scan, bug audit, or asks to review code with multiple AI models.
---

# Multi-Model Code Auditor — Elite Edition

Audit any codebase by sending it to the best available frontier models (GPT-5.x, Claude Opus/Sonnet 4.x, Gemini 3.x/2.5 Pro, o4-mini, Codex, DeepSeek), each with model-specific prompt engineering optimized for their strengths. Findings are cross-compared, deduplicated by multi-signal similarity, merged by selecting the best description and fix from each cluster, and ranked by consensus. After fixes are applied, a verification re-audit confirms all issues are resolved.

## When to Use

This skill activates when the user requests any of: "security audit", "code review", "vulnerability scan", "bug audit", "review with multiple models", "use AI to check my code", "find all bugs", "privacy audit", or "penetration test the code".

## Architecture

The system uses 6 model tiers, each assigned a specialized role:

| Tier | Family | Role | Strengths |
|------|--------|------|-----------|
| 1 | OpenAI Flagship | Structured Analysis Lead | Schema precision, systematic enumeration, exhaustive coverage |
| 2 | Anthropic Flagship | Nuanced Reasoning Specialist | Contextual understanding, subtle logic bugs, low hallucination |
| 3 | Google Flagship | Cross-File Architecture Analyst | Large context, cross-module dependency tracking, design flaws |
| 4 | OpenAI Reasoning | Deep Reasoning & Logic Auditor | Race conditions, state machines, multi-step attack chains |
| 5 | OpenAI Codex | Code-Native Security Scanner | Pattern recognition, API misuse, production-ready fixes |
| 6 | DeepSeek | Algorithmic & Low-Level Bug Hunter | Algorithmic correctness, boundary conditions, data handling |

Each model receives a tailored prompt with: a specialized persona, model-specific focus areas, few-shot examples calibrating severity levels, negative examples showing what NOT to report, and a strict output schema requiring exact code citations.

When the OpenRouter model discovery API is unavailable, the script falls back to a curated list of known-good stable models (Claude Sonnet 4, Gemini 2.5 Pro, GPT-4.1) to ensure the audit always runs.

## Full Workflow

The audit follows a 7-step pipeline. Steps 1-3 use the scripts; steps 4-7 are performed by the agent.

### Step 1: Discover Best Models

The script dynamically queries OpenRouter to find the best available model from each tier. Run with `--discover` to preview:

```bash
python3 /home/ubuntu/skills/multi-model-code-auditor/scripts/run_audit.py --discover
```

### Step 2: Audit with All Models

```bash
# Sequential (default) — queries models one at a time with retry logic
python3 /home/ubuntu/skills/multi-model-code-auditor/scripts/run_audit.py <project_dir> --output audit_results.json

# Parallel (recommended) — queries all models simultaneously for speed
python3 /home/ubuntu/skills/multi-model-code-auditor/scripts/run_audit.py <project_dir> --output audit_results.json --parallel

# Specific models — override auto-discovery with chosen models
python3 /home/ubuntu/skills/multi-model-code-auditor/scripts/run_audit.py <project_dir> --models anthropic/claude-sonnet-4,google/gemini-2.5-pro --output audit_results.json
```

The script collects source files (prioritizing entry points, config, and auth files first), builds the payload with smart truncation for large codebases, and queries each model with its tailored prompt. Each model gets up to 3 retry attempts with exponential backoff (10s, 30s, 60s) on transient failures.

**Reliability features:**
- Exponential backoff retry (3 attempts per model) for timeouts, rate limits, and server errors
- Smart payload truncation: high-priority files (entry points, config, auth) are preserved in full; low-priority files are truncated to fit within model context limits
- Progress indicators during long API calls
- Per-model intra-finding deduplication
- Family-aware model override: `--models` flag correctly infers prompt engineering from model ID prefix
- Known-good fallback models when API discovery fails

### Step 3: Cross-Compare and Merge

```bash
python3 /home/ubuntu/skills/multi-model-code-auditor/scripts/cross_compare.py audit_results.json --output combined_findings.json
```

The merger performs:

1. **Multi-signal clustering** — groups findings by file + function + title + CWE + code snippet similarity (weighted Jaccard). Cross-file findings can still cluster when they share function names, CWEs, and code patterns.
2. **Best-result selection** — for each cluster, picks: highest severity, most concise title, longest description, most complete fix code, most common CWE, highest CVSS
3. **False positive filtering** — flags single-model findings matching known hallucination patterns (phantom functions, safe function misidentification, suspicious line numbers, theoretical-only risks)
4. **Consensus ranking** — sorts by: model count (desc) then family diversity (desc) then severity (desc) then CVSS (desc)

Output: `combined_findings.json` + `combined_findings_report.md` (generated from `templates/report_template.md`)

### Step 4: Validate Before Implementing

Before implementing ANY fix, the agent MUST:

1. **Verify the code exists** — `grep -n "function_name" file.py` to confirm the function and code snippet are real. Models hallucinate line numbers and function names.
2. **Verify the vulnerability is real** — read the actual code at the cited location. See `references/audit_categories.md` for common false positives.
3. **Check cross-module impact** — will this fix break callers? Check imports and function signatures.
4. **Skip if already fixed** — if a prior audit round already addressed this, document as "already fixed".

Findings with 3+ model consensus are almost certainly real. Findings from only 1 model require manual verification. Findings from 1 model that match false positive patterns should be discarded.

### Step 5: Implement All Verified Fixes

For each confirmed finding, ordered by severity (critical first):

1. Read the actual source file to understand the full context
2. Apply the fix using the best `fix_code` from the merged results
3. Verify the fix doesn't introduce new issues
4. Run existing test suites after each fix

### Step 6: Verification Re-Audit

After ALL fixes are applied, run the full audit pipeline again:

```bash
python3 /home/ubuntu/skills/multi-model-code-auditor/scripts/run_audit.py <project_dir> --output reaudit_results.json --parallel
python3 /home/ubuntu/skills/multi-model-code-auditor/scripts/cross_compare.py reaudit_results.json --output reaudit_combined.json
```

Compare the re-audit results against the original findings. All previously identified issues should now be resolved. Any NEW findings from the re-audit should be evaluated and fixed.

### Step 7: Generate Final Report

The final deliverable includes:

1. **Audit report** — all findings with consensus matrix, severity breakdown, and fix details
2. **Model agreement matrix** — which models found which issues
3. **False positives identified** — with explanation of why they were filtered
4. **Before/after comparison** — original audit vs re-audit results
5. **Test results** — confirmation that all fixes pass tests

The report is automatically generated from `templates/report_template.md` during the cross-compare step.

## Key Principles

**Consensus is king.** The entire system is built around the insight that when multiple independent models agree on a finding, it is almost certainly real. Single-model findings are noise until verified.

**Model diversity matters.** A finding confirmed by GPT + Claude + Gemini (3 different families) is stronger than one confirmed by GPT-5.4 + GPT-5.3 + Codex (same family). The `model_families` count captures this.

**Verify before fixing.** Never implement a fix based solely on a model's claim. Always read the actual code. Models hallucinate function names, line numbers, and even entire vulnerabilities.

**Best-of-breed merging.** Each model has different strengths. The merger takes the best description from the most detailed model, the best fix from the most code-native model, and the highest severity from the most cautious model.

**Preserve functionality.** Every fix must be tested. A security fix that breaks the application is worse than the vulnerability it addresses.

**Re-audit closes the loop.** The verification pass ensures fixes actually work and don't introduce regressions. It also catches issues that only become visible after other fixes are applied.

**Retry and resilience.** API calls fail frequently due to rate limits, timeouts, and server errors. The script retries each model up to 3 times with exponential backoff, and falls back to known-good models when discovery fails.

## Resources

| Resource | Purpose |
|----------|---------|
| `scripts/run_audit.py` | Discovers frontier models, sends codebase with model-specific prompts, collects findings. Supports `--parallel`, `--models`, `--discover` flags. |
| `scripts/cross_compare.py` | Clusters, merges, filters false positives, ranks by consensus, generates template-based report |
| `references/audit_categories.md` | Security checklist, vulnerability taxonomy, and common false positive patterns |
| `templates/report_template.md` | Markdown template for the final audit report with placeholder substitution |
