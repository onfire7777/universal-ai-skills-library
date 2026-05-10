# Security Audit Report — {{PROJECT_NAME}}

**Date:** {{DATE}}
**Auditor:** Multi-Model Code Auditor (Elite Edition)
**Scope:** Full codebase — {{FILES_AUDITED}} files, ~{{TOKEN_COUNT}} tokens

## Executive Summary

{{EXECUTIVE_SUMMARY}}

## Models Used

| Model | Family | Role | Findings |
|-------|--------|------|----------|
{{MODEL_TABLE_ROWS}}

## Results Overview

| Metric | Value |
|--------|-------|
| Total raw findings (pre-dedup) | {{RAW_FINDINGS}} |
| Clusters formed | {{CLUSTERS}} |
| Confirmed findings | {{CONFIRMED}} |
| False positives filtered | {{FALSE_POSITIVES}} |
| Fixes implemented | {{FIXES_IMPLEMENTED}} |
| Re-audit clean | {{REAUDIT_CLEAN}} |

## Severity Breakdown

| Severity | Count | Percentage |
|----------|-------|------------|
| Critical | {{CRITICAL}} | {{CRITICAL_PCT}} |
| High | {{HIGH}} | {{HIGH_PCT}} |
| Medium | {{MEDIUM}} | {{MEDIUM_PCT}} |
| Low | {{LOW}} | {{LOW_PCT}} |

## Consensus Matrix

This table shows which models independently identified each finding. Higher consensus indicates higher confidence.

| ID | Severity | Consensus | Families | File | Title |
|----|----------|-----------|----------|------|-------|
{{CONSENSUS_MATRIX_ROWS}}

## False Positives Filtered

The following findings were identified as likely false positives and excluded from the fix list.

{{FALSE_POSITIVE_DETAILS}}

## Detailed Findings

{{DETAILED_FINDINGS}}

## Verification Re-Audit

After all fixes were applied, the full audit was re-run to verify resolution.

| Metric | Original | Re-Audit |
|--------|----------|----------|
| Critical | {{ORIG_CRITICAL}} | {{REAUDIT_CRITICAL}} |
| High | {{ORIG_HIGH}} | {{REAUDIT_HIGH}} |
| Medium | {{ORIG_MEDIUM}} | {{REAUDIT_MEDIUM}} |
| Low | {{ORIG_LOW}} | {{REAUDIT_LOW}} |
| Total | {{ORIG_TOTAL}} | {{REAUDIT_TOTAL}} |

{{REAUDIT_NOTES}}

## Recommendations

{{RECOMMENDATIONS}}
