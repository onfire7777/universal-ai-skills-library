#!/usr/bin/env python3
"""Cross-Compare & Merge — Elite Edition

Takes raw findings from multiple frontier models, deduplicates via multi-signal
similarity, ranks by consensus strength, selects the best description and fix
for each unique issue, and produces the ultimate merged audit result.

Pipeline:
  1. Normalize all findings into a common schema
  2. Cluster similar findings using multi-signal similarity (file + function + title + CWE + code)
  3. For each cluster, pick the best description, best fix, and highest-confidence severity
  4. Rank by consensus (models agreeing) then severity then CVSS
  5. Flag likely false positives (single-model findings with hallucination markers)
  6. Output: merged JSON + Markdown report

Usage:
    python3 cross_compare.py audit_results.json [--output combined.json] [--threshold 0.4]
"""

import argparse
import json
import re
import sys
from collections import Counter
from pathlib import Path

# ─── Constants ────────────────────────────────────────────────────────────────

SEVERITY_SCORE = {"critical": 4, "high": 3, "medium": 2, "low": 1}

# Known false positive patterns — findings matching these are flagged.
# Each pattern has: a regex to match, a reason explaining why it's a false positive,
# and optionally a verify command to confirm.
# Patterns that match against the TITLE + DESCRIPTION fields only
TEXT_FP_PATTERNS = [
    {
        "pattern": r"yaml\.safe_load.*unsafe",
        "reason": "yaml.safe_load IS the safe variant — this is a model hallucination",
    },
    {
        "pattern": r"if.*attacker.*root|if.*root.*access",
        "reason": "Theoretical-only risk requiring root access — not a real vulnerability",
    },
]

# Patterns that match against the VULNERABLE_CODE field only
CODE_FP_PATTERNS = [
    {
        "pattern": r"\beval\s*\(",
        "verify": "grep -rn 'eval(' to confirm it actually exists",
        "reason": "Model may be hallucinating eval() usage — verify with grep",
    },
    {
        "pattern": r"\bexec\s*\(",
        "verify": "grep -rn 'exec(' to confirm it actually exists",
        "reason": "Model may be hallucinating exec() usage — verify with grep",
    },
    {
        "pattern": r"pickle\.loads",
        "verify": "grep -rn 'pickle.loads' to confirm it actually exists",
        "reason": "Model may be hallucinating pickle.loads usage — verify with grep",
    },
]

# ─── Similarity Engine ───────────────────────────────────────────────────────

def normalize_text(text: str) -> str:
    """Normalize text for comparison: lowercase, strip non-alphanumeric."""
    return re.sub(r"[^a-z0-9\s]", "", (text or "").lower()).strip()


def word_set(text: str) -> set:
    """Extract word set from text for Jaccard similarity."""
    return set(normalize_text(text).split())


def jaccard(a: set, b: set) -> float:
    """Jaccard similarity between two sets."""
    if not a and not b:
        return 0.0
    return len(a & b) / len(a | b)


def code_similarity(a: str, b: str) -> float:
    """Compare code snippets by normalized token overlap."""
    a_tokens = set(re.findall(r"[a-zA-Z_]\w*", a or ""))
    b_tokens = set(re.findall(r"[a-zA-Z_]\w*", b or ""))
    if not a_tokens or not b_tokens:
        return 0.0
    return len(a_tokens & b_tokens) / len(a_tokens | b_tokens)


def compute_similarity(a: dict, b: dict) -> float:
    """Multi-signal similarity between two findings.

    Signals (weighted):
      - File match:     0.25 (exact match on normalized filename)
      - Function match: 0.20 (exact match on function name)
      - Title overlap:  0.20 (Jaccard on title words)
      - CWE match:      0.15 (exact CWE ID match)
      - Code overlap:   0.10 (token overlap in vulnerable_code)
      - Category match:  0.10 (same category)

    When files differ, file_score is 0 but other signals can still drive
    clustering — this allows cross-file architectural findings to be grouped
    when they share the same function, CWE, title, and code patterns.
    """
    # File match (normalize path separators)
    file_a = normalize_text(a.get("file", "")).replace(" ", "")
    file_b = normalize_text(b.get("file", "")).replace(" ", "")
    file_score = 1.0 if file_a and file_b and file_a == file_b else 0.0

    # Function match
    func_a = normalize_text(a.get("function", ""))
    func_b = normalize_text(b.get("function", ""))
    func_score = 1.0 if func_a and func_b and func_a == func_b else 0.0

    # Title Jaccard
    title_score = jaccard(word_set(a.get("title", "")), word_set(b.get("title", "")))

    # CWE match
    cwe_a = (a.get("cwe") or "").upper()
    cwe_b = (b.get("cwe") or "").upper()
    cwe_score = 1.0 if cwe_a and cwe_b and cwe_a == cwe_b else 0.0

    # Code similarity
    code_score = code_similarity(a.get("vulnerable_code", ""), b.get("vulnerable_code", ""))

    # Category match
    cat_a = normalize_text(a.get("category", ""))
    cat_b = normalize_text(b.get("category", ""))
    cat_score = 1.0 if cat_a and cat_b and cat_a == cat_b else 0.0

    total = (
        0.25 * file_score
        + 0.20 * func_score
        + 0.20 * title_score
        + 0.15 * cwe_score
        + 0.10 * code_score
        + 0.10 * cat_score
    )

    # When files are completely different AND no other strong signals match,
    # the score will naturally be low. But if function + title + CWE all match
    # across different files, the score can still exceed the threshold — which
    # is correct for cross-file architectural findings.
    return total


# ─── Clustering ──────────────────────────────────────────────────────────────

def cluster_findings(all_findings: list[dict], threshold: float = 0.4) -> list[dict]:
    """Cluster similar findings across models using greedy agglomerative clustering.

    Each finding is compared against existing cluster representatives. If the
    similarity exceeds the threshold, it joins that cluster. Otherwise, it
    starts a new cluster.
    """
    clusters = []

    for finding in all_findings:
        best_cluster = None
        best_score = 0.0

        for cluster in clusters:
            # Compare against the cluster representative
            score = compute_similarity(finding, cluster["representative"])
            if score > best_score:
                best_score = score
                best_cluster = cluster

        if best_cluster and best_score >= threshold:
            best_cluster["members"].append(finding)
            best_cluster["models"].add(finding["_source_model"])
            best_cluster["families"].add(finding["_source_family"])
        else:
            clusters.append({
                "representative": finding,
                "members": [finding],
                "models": {finding["_source_model"]},
                "families": {finding["_source_family"]},
            })

    return clusters


# ─── Best-Result Selection ───────────────────────────────────────────────────

def select_best_from_cluster(cluster: dict) -> dict:
    """From a cluster of similar findings, select the best version of each field.

    Strategy:
      - Severity: take the HIGHEST severity (most cautious)
      - Title: take the most concise one
      - Description: take the LONGEST (most detailed)
      - Attack chain: take the LONGEST
      - Vulnerable code: take the LONGEST (most complete snippet)
      - Fix code: take the LONGEST (most complete fix)
      - CWE: take the most common one across members
      - CVSS: take the HIGHEST
    """
    members = cluster["members"]

    # Severity: highest across all members
    severity_scores = [SEVERITY_SCORE.get(m.get("severity", "low"), 1) for m in members]
    max_sev_idx = severity_scores.index(max(severity_scores))
    best_severity = members[max_sev_idx].get("severity", "medium")

    # Title: shortest non-empty (most concise)
    titles = [(len(m.get("title", "")), m.get("title", "")) for m in members if m.get("title")]
    best_title = min(titles, key=lambda x: x[0])[1] if titles else "Untitled"

    # Description: longest (most detailed)
    best_description = max(
        (m.get("description", "") for m in members), key=len, default=""
    )

    # Attack chain: longest
    best_attack_chain = max(
        (m.get("attack_chain", "") for m in members), key=len, default=""
    )

    # Vulnerable code: longest (most complete)
    best_vuln_code = max(
        (m.get("vulnerable_code", "") for m in members), key=len, default=""
    )

    # Fix code: longest (most complete)
    best_fix_code = max(
        (m.get("fix_code", "") for m in members), key=len, default=""
    )

    # CWE: most common
    cwes = [m.get("cwe", "") for m in members if m.get("cwe")]
    best_cwe = Counter(cwes).most_common(1)[0][0] if cwes else ""

    # CVSS: highest
    cvss_values = [m.get("cvss_estimate", 0) for m in members if m.get("cvss_estimate")]
    best_cvss = max(cvss_values) if cvss_values else 0.0

    # File and function: from the member with the longest description (most detailed analysis)
    detail_member = max(members, key=lambda m: len(m.get("description", "")))
    best_file = detail_member.get("file", "unknown")
    best_function = detail_member.get("function", "unknown")
    best_line_range = detail_member.get("line_range", "")

    # Category: most common
    cats = [m.get("category", "") for m in members if m.get("category")]
    best_category = Counter(cats).most_common(1)[0][0] if cats else "security"

    return {
        "severity": best_severity,
        "category": best_category,
        "file": best_file,
        "function": best_function,
        "line_range": best_line_range,
        "title": best_title,
        "description": best_description,
        "attack_chain": best_attack_chain,
        "vulnerable_code": best_vuln_code,
        "fix_code": best_fix_code,
        "cwe": best_cwe,
        "cvss_estimate": best_cvss,
        "consensus": len(cluster["models"]),
        "model_families": len(cluster["families"]),
        "reported_by": sorted(cluster["models"]),
        "reported_by_families": sorted(cluster["families"]),
    }


# ─── False Positive Detection ────────────────────────────────────────────────

def check_false_positive(finding: dict) -> dict | None:
    """Check if a finding matches known false positive patterns.

    Returns a dict with reason/verify fields if it matches, or None if clean.
    Uses field-specific matching to avoid false triggers from unrelated fields.
    """
    # Check title + description for text-level FP patterns
    text_fields = " ".join([
        finding.get("title", ""),
        finding.get("description", ""),
    ]).lower()

    for fp in TEXT_FP_PATTERNS:
        if re.search(fp["pattern"], text_fields, re.IGNORECASE):
            return {
                "reason": fp.get("reason", "Matches known false positive pattern"),
                "verify": fp.get("verify", "Manual verification required"),
            }

    # Check vulnerable_code for code-level FP patterns
    vuln_code = (finding.get("vulnerable_code", "") or "").lower()
    for fp in CODE_FP_PATTERNS:
        if re.search(fp["pattern"], vuln_code, re.IGNORECASE):
            return {
                "reason": fp.get("reason", "Matches known false positive pattern"),
                "verify": fp.get("verify", "Manual verification required"),
            }

    # Check line_range specifically for suspiciously high line numbers
    line_range = finding.get("line_range", "")
    if line_range:
        numbers = re.findall(r"\d+", str(line_range))
        if any(int(n) > 9999 for n in numbers):
            return {
                "reason": "Suspiciously high line number (>9999) — model may be hallucinating code location",
                "verify": "Verify the file actually has this many lines",
            }

    # Additional heuristic: if vulnerable_code is suspiciously short (< 5 chars)
    # Real code snippets like 'except:' (7 chars) or 'eval(x)' (7 chars) are short
    # but legitimate. Only flag truly trivial snippets like single words.
    vuln_code = finding.get("vulnerable_code", "")
    if vuln_code and len(vuln_code.strip()) < 5:
        return {
            "reason": "Vulnerable code snippet is suspiciously short (< 5 chars)",
            "verify": "Check if the cited code actually exists in the source",
        }

    return None


# ─── Ranking ─────────────────────────────────────────────────────────────────

def rank_findings(findings: list[dict]) -> list[dict]:
    """Rank findings by: consensus (desc) -> family diversity (desc) -> severity (desc) -> CVSS (desc)."""
    return sorted(
        findings,
        key=lambda f: (
            -f.get("consensus", 1),
            -f.get("model_families", 1),
            -SEVERITY_SCORE.get(f.get("severity", "low"), 1),
            -f.get("cvss_estimate", 0),
        ),
    )


# ─── Report Generation ──────────────────────────────────────────────────────

def generate_report(findings: list[dict], false_positives: list[dict], meta: dict) -> str:
    """Generate a comprehensive Markdown audit report using the template.

    Loads the report template from templates/report_template.md and fills in
    all placeholders. Falls back to inline generation if the template is missing.
    """
    template_path = Path(__file__).parent.parent / "templates" / "report_template.md"

    if template_path.exists():
        return _generate_from_template(findings, false_positives, meta, template_path)
    else:
        print(f"  WARNING: Template not found at {template_path}. Using inline report.")
        return _generate_inline(findings, false_positives, meta)


def _generate_from_template(findings: list[dict], false_positives: list[dict],
                            meta: dict, template_path: Path) -> str:
    """Fill the report template with actual data."""
    template = template_path.read_text(encoding="utf-8")

    total = len(findings)
    sev_counts = Counter(f.get("severity", "low") for f in findings)

    # Model table rows
    model_rows = []
    for m in meta.get("models", []):
        # Count findings per model
        count = sum(1 for f in findings if m["id"] in f.get("reported_by", []))
        model_rows.append(f"| `{m['id']}` | {m['family']} | {m['role']} | {count} |")

    # Consensus matrix rows
    consensus_rows = []
    for f in findings:
        consensus_str = f"{f['consensus']}/{meta.get('models_queried', '?')}"
        families_str = ", ".join(f.get("reported_by_families", []))
        consensus_rows.append(
            f"| {f['id']} | **{f.get('severity', '?').upper()}** | "
            f"{consensus_str} | {families_str} | `{f.get('file', '?')}` | {f.get('title', '?')} |"
        )

    # False positive details
    if false_positives:
        fp_lines = []
        for fp in false_positives:
            fp_lines.append(
                f"- **{fp.get('title', '?')}** in `{fp.get('file', '?')}` — "
                f"{fp.get('_fp_reason', 'Pattern match')}"
            )
        fp_text = "\n".join(fp_lines)
    else:
        fp_text = "No false positives were detected."

    # Detailed findings
    detail_lines = []
    for f in findings:
        detail_lines.append(f"### {f['id']}: {f.get('title', 'Untitled')}\n")
        detail_lines.append(f"| Field | Value |")
        detail_lines.append(f"|---|---|")
        detail_lines.append(f"| Severity | **{f.get('severity', '?').upper()}** |")
        detail_lines.append(f"| Category | {f.get('category', '?')} |")
        detail_lines.append(f"| CWE | {f.get('cwe', 'N/A')} |")
        detail_lines.append(f"| CVSS Estimate | {f.get('cvss_estimate', 'N/A')} |")
        detail_lines.append(f"| File | `{f.get('file', '?')}` |")
        detail_lines.append(f"| Function | `{f.get('function', '?')}` |")
        detail_lines.append(f"| Consensus | {f['consensus']}/{meta.get('models_queried', '?')} models |")
        detail_lines.append(f"| Reported By | {', '.join(f.get('reported_by_families', []))} |")
        detail_lines.append(f"\n{f.get('description', 'No description.')}\n")
        if f.get("attack_chain"):
            detail_lines.append(f"**Attack Chain:** {f['attack_chain']}\n")
        if f.get("vulnerable_code"):
            detail_lines.append(f"**Vulnerable Code:**\n```\n{f['vulnerable_code']}\n```\n")
        if f.get("fix_code"):
            detail_lines.append(f"**Fix:**\n```\n{f['fix_code']}\n```\n")
        detail_lines.append("---\n")

    # Executive summary
    if total == 0:
        exec_summary = "The audit found no significant security issues in the codebase."
    else:
        sev_parts = []
        for sev in ["critical", "high", "medium", "low"]:
            c = sev_counts.get(sev, 0)
            if c > 0:
                sev_parts.append(f"{c} {sev}")
        exec_summary = (
            f"The multi-model audit identified {total} confirmed findings "
            f"({', '.join(sev_parts)}). "
            f"{len(false_positives)} potential false positives were filtered. "
            f"Findings with multi-model consensus are high-confidence issues that "
            f"should be prioritized for remediation."
        )

    # Calculate percentages
    def pct(count):
        return f"{count/total*100:.0f}%" if total > 0 else "0%"

    # Perform substitutions
    replacements = {
        "{{PROJECT_NAME}}": meta.get("project_name", "Unknown"),
        "{{DATE}}": meta.get("timestamp", "Unknown"),
        "{{FILES_AUDITED}}": str(meta.get("files_audited", "?")),
        "{{TOKEN_COUNT}}": f"{meta.get('payload_tokens_est', 0):,}",
        "{{EXECUTIVE_SUMMARY}}": exec_summary,
        "{{MODEL_TABLE_ROWS}}": "\n".join(model_rows),
        "{{RAW_FINDINGS}}": str(meta.get("total_raw_findings", "?")),
        "{{CLUSTERS}}": str(meta.get("clusters_formed", "?")),
        "{{CONFIRMED}}": str(total),
        "{{FALSE_POSITIVES}}": str(len(false_positives)),
        "{{FIXES_IMPLEMENTED}}": "Pending",
        "{{REAUDIT_CLEAN}}": "Pending",
        "{{CRITICAL}}": str(sev_counts.get("critical", 0)),
        "{{CRITICAL_PCT}}": pct(sev_counts.get("critical", 0)),
        "{{HIGH}}": str(sev_counts.get("high", 0)),
        "{{HIGH_PCT}}": pct(sev_counts.get("high", 0)),
        "{{MEDIUM}}": str(sev_counts.get("medium", 0)),
        "{{MEDIUM_PCT}}": pct(sev_counts.get("medium", 0)),
        "{{LOW}}": str(sev_counts.get("low", 0)),
        "{{LOW_PCT}}": pct(sev_counts.get("low", 0)),
        "{{CONSENSUS_MATRIX_ROWS}}": "\n".join(consensus_rows),
        "{{FALSE_POSITIVE_DETAILS}}": fp_text,
        "{{DETAILED_FINDINGS}}": "\n".join(detail_lines),
        "{{ORIG_CRITICAL}}": str(sev_counts.get("critical", 0)),
        "{{ORIG_HIGH}}": str(sev_counts.get("high", 0)),
        "{{ORIG_MEDIUM}}": str(sev_counts.get("medium", 0)),
        "{{ORIG_LOW}}": str(sev_counts.get("low", 0)),
        "{{ORIG_TOTAL}}": str(total),
        "{{REAUDIT_CRITICAL}}": "—",
        "{{REAUDIT_HIGH}}": "—",
        "{{REAUDIT_MEDIUM}}": "—",
        "{{REAUDIT_LOW}}": "—",
        "{{REAUDIT_TOTAL}}": "—",
        "{{REAUDIT_NOTES}}": "Re-audit has not been performed yet.",
        "{{RECOMMENDATIONS}}": _generate_recommendations(findings),
    }

    report = template
    for placeholder, value in replacements.items():
        report = report.replace(placeholder, value)

    return report


def _generate_recommendations(findings: list[dict]) -> str:
    """Generate actionable recommendations based on findings."""
    if not findings:
        return "No specific recommendations — the codebase appears secure."

    recs = []
    sev_counts = Counter(f.get("severity", "low") for f in findings)
    categories = Counter(f.get("category", "unknown") for f in findings)

    if sev_counts.get("critical", 0) > 0:
        recs.append(
            "1. **Immediate action required:** Address all critical findings before the next release. "
            "These represent exploitable vulnerabilities that could lead to system compromise."
        )

    if categories.get("injection", 0) > 0:
        recs.append(
            f"2. **Input sanitization:** {categories['injection']} injection-related findings were detected. "
            "Implement a centralized input validation layer and avoid string interpolation in commands/queries."
        )

    if categories.get("error_handling", 0) > 0:
        recs.append(
            f"3. **Error handling:** {categories['error_handling']} error handling issues found. "
            "Replace broad exception handlers with specific catches and ensure errors are properly logged."
        )

    high_consensus = [f for f in findings if f.get("consensus", 1) >= 2]
    if high_consensus:
        recs.append(
            f"4. **High-confidence issues:** {len(high_consensus)} findings were confirmed by multiple "
            "independent AI models. These should be prioritized as they have the highest confidence of being real."
        )

    if not recs:
        recs.append("1. Review all findings and implement fixes in order of severity.")
        recs.append("2. Run the verification re-audit after all fixes are applied.")

    return "\n\n".join(recs)


def _generate_inline(findings: list[dict], false_positives: list[dict], meta: dict) -> str:
    """Fallback inline report generation when template is unavailable."""
    lines = [
        "# Multi-Model Code Audit — Merged Results\n",
        f"**Project:** {meta.get('project_name', 'Unknown')}",
        f"**Date:** {meta.get('timestamp', 'Unknown')}",
        f"**Files audited:** {meta.get('files_audited', '?')}",
        f"**Payload:** ~{meta.get('payload_tokens_est', 0):,} tokens",
        f"**Models queried:** {meta.get('models_queried', '?')}",
        "",
    ]

    # Model table
    lines.append("## Models Used\n")
    lines.append("| Model | Family | Role |")
    lines.append("|---|---|---|")
    for m in meta.get("models", []):
        lines.append(f"| `{m['id']}` | {m['family']} | {m['role']} |")

    # Summary stats
    confirmed = [f for f in findings if not f.get("_false_positive")]
    lines.append(f"\n## Summary\n")
    lines.append(f"**{len(confirmed)} confirmed findings** ({len(false_positives)} false positives filtered)\n")

    sev_counts = Counter(f.get("severity", "low") for f in confirmed)
    lines.append(f"| Severity | Count |")
    lines.append(f"|---|---|")
    for sev in ["critical", "high", "medium", "low"]:
        if sev_counts.get(sev, 0) > 0:
            lines.append(f"| {sev.upper()} | {sev_counts[sev]} |")

    # Consensus matrix
    lines.append("\n## Consensus Matrix\n")
    lines.append("| ID | Severity | Consensus | File | Title |")
    lines.append("|---|---|---|---|---|")
    for f in confirmed:
        consensus_str = f"{f['consensus']}/{meta.get('models_queried', '?')}"
        lines.append(
            f"| {f['id']} | **{f.get('severity', '?').upper()}** | "
            f"{consensus_str} | `{f.get('file', '?')}` | {f.get('title', '?')} |"
        )

    # False positives
    if false_positives:
        lines.append("\n## False Positives Filtered\n")
        for fp in false_positives:
            lines.append(f"- **{fp.get('title', '?')}** in `{fp.get('file', '?')}` — {fp.get('_fp_reason', 'Pattern match')}")

    # Detailed findings
    lines.append("\n## Detailed Findings\n")
    for f in confirmed:
        lines.append(f"### {f['id']}: {f.get('title', 'Untitled')}\n")
        lines.append(f"| Field | Value |")
        lines.append(f"|---|---|")
        lines.append(f"| Severity | **{f.get('severity', '?').upper()}** |")
        lines.append(f"| Category | {f.get('category', '?')} |")
        lines.append(f"| CWE | {f.get('cwe', 'N/A')} |")
        lines.append(f"| CVSS Estimate | {f.get('cvss_estimate', 'N/A')} |")
        lines.append(f"| File | `{f.get('file', '?')}` |")
        lines.append(f"| Function | `{f.get('function', '?')}` |")
        lines.append(f"| Consensus | {f['consensus']}/{meta.get('models_queried', '?')} models |")
        lines.append(f"| Reported By | {', '.join(f.get('reported_by_families', []))} |")
        lines.append(f"\n{f.get('description', 'No description.')}\n")
        if f.get("attack_chain"):
            lines.append(f"**Attack Chain:** {f['attack_chain']}\n")
        if f.get("vulnerable_code"):
            lines.append(f"**Vulnerable Code:**\n```\n{f['vulnerable_code']}\n```\n")
        if f.get("fix_code"):
            lines.append(f"**Fix:**\n```\n{f['fix_code']}\n```\n")
        lines.append("---\n")

    return "\n".join(lines)


# ─── Main Pipeline ───────────────────────────────────────────────────────────

def cross_compare(input_file: str, output_file: str, threshold: float = 0.4) -> dict:
    """Full cross-comparison and merging pipeline."""
    data = json.loads(Path(input_file).read_text())

    print("=" * 70)
    print("CROSS-COMPARE & MERGE — ELITE EDITION")
    print("=" * 70)

    # Step 1: Collect all findings with source attribution
    all_findings = []
    for result in data.get("results", []):
        model = result.get("model", "unknown")
        family = result.get("family", "unknown")
        status = result.get("status", "error")
        if status != "success":
            print(f"  Skipping {model} (status: {status})")
            continue
        findings = result.get("result", {}).get("findings", [])
        print(f"  {model}: {len(findings)} findings")
        for f in findings:
            f["_source_model"] = model
            f["_source_family"] = family
            all_findings.append(f)

    print(f"\nTotal raw findings: {len(all_findings)}")

    if not all_findings:
        print("  No findings to process. All models may have failed or returned empty results.")
        output = {
            "project": data.get("project"),
            "project_name": data.get("project_name"),
            "files_audited": data.get("files_audited"),
            "payload_tokens_est": data.get("payload_tokens_est"),
            "models_queried": data.get("models_queried"),
            "models": data.get("models", []),
            "timestamp": data.get("timestamp"),
            "total_raw_findings": 0,
            "clusters_formed": 0,
            "confirmed_findings": 0,
            "false_positives_filtered": 0,
            "findings": [],
            "false_positives": [],
        }
        Path(output_file).write_text(json.dumps(output, indent=2))
        print(f"\nEmpty results saved to {output_file}")
        return output

    # Step 2: Cluster similar findings
    clusters = cluster_findings(all_findings, threshold)
    print(f"Clusters formed: {len(clusters)}")

    # Step 3: Select best from each cluster
    merged = []
    for i, cluster in enumerate(clusters, 1):
        best = select_best_from_cluster(cluster)
        best["id"] = f"AUDIT-{i:03d}"
        merged.append(best)

    # Step 4: Check for false positives
    confirmed = []
    false_positives = []
    for f in merged:
        fp_check = check_false_positive(f)
        if fp_check and f.get("consensus", 1) <= 1:
            # Only flag as FP if low consensus AND matches FP pattern
            f["_false_positive"] = True
            f["_fp_reason"] = fp_check["reason"]
            false_positives.append(f)
        else:
            confirmed.append(f)

    # Step 5: Rank by consensus then severity
    ranked = rank_findings(confirmed)

    # Re-number after filtering
    for i, f in enumerate(ranked, 1):
        f["id"] = f"AUDIT-{i:03d}"

    print(f"\nConfirmed findings: {len(ranked)}")
    print(f"False positives filtered: {len(false_positives)}")

    # Consensus breakdown
    consensus_counts = Counter(f.get("consensus", 1) for f in ranked)
    for c in sorted(consensus_counts.keys(), reverse=True):
        print(f"  {consensus_counts[c]} findings with {c}-model consensus")

    # Step 6: Save outputs
    output = {
        "project": data.get("project"),
        "project_name": data.get("project_name"),
        "files_audited": data.get("files_audited"),
        "payload_tokens_est": data.get("payload_tokens_est"),
        "models_queried": data.get("models_queried"),
        "models": data.get("models", []),
        "timestamp": data.get("timestamp"),
        "total_raw_findings": len(all_findings),
        "clusters_formed": len(clusters),
        "confirmed_findings": len(ranked),
        "false_positives_filtered": len(false_positives),
        "findings": ranked,
        "false_positives": false_positives,
    }

    Path(output_file).write_text(json.dumps(output, indent=2))
    print(f"\nMerged findings saved to {output_file}")

    # Generate Markdown report
    report_path = output_file.replace(".json", "_report.md")
    report = generate_report(ranked, false_positives, output)
    Path(report_path).write_text(report)
    print(f"Markdown report saved to {report_path}")

    # Print top findings
    if ranked:
        print(f"\n{'─'*50}")
        print("TOP FINDINGS (by consensus):")
        for f in ranked[:10]:
            print(f"  [{f['severity'].upper():8s}] {f['consensus']}/{data.get('models_queried','?')} "
                  f"| {f['file']:30s} | {f['title']}")

    return output


# ─── CLI ──────────────────────────────────────────────────────────────────────

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Cross-compare and merge multi-model audit findings")
    parser.add_argument("input_file", help="Path to audit_results.json from run_audit.py")
    parser.add_argument("--output", default="combined_findings.json", help="Output JSON file")
    parser.add_argument("--threshold", type=float, default=0.4,
                        help="Similarity threshold for clustering (0.0-1.0, default 0.4)")
    args = parser.parse_args()

    cross_compare(args.input_file, args.output, args.threshold)
