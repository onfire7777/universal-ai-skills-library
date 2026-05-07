# Skill Debug Report: `context-anchor`

**Generated:** 2026-05-07 20:07 UTC
**Models:** Claude Opus 4.6 + Manus gpt-4.1-mini

## Model Status

| Model | Status | Findings |
|---|---|---|
| Claude Opus 4.6 | OK | 5 |
| Manus gpt-4.1-mini | OK | 8 |

## Overall Health: NEEDS ATTENTION

**8 issues found** (5 confirmed by both models)

## Findings Summary

| ID | Severity | Category | File | Title | Consensus |
|---|---|---|---|---|---|
| DBG-001 | medium | script | `scripts/anchor.py` | _insert_into_section never returns None on section-not-found but callers check for None | Both |
| DBG-002 | medium | script | `scripts/anchor.py` | Fallback keyword relevance scoring produces misleading scores for short messages | Both |
| DBG-003 | medium | integration | `SKILL.md` | SKILL.md description is comprehensive and well-structured for Manus triggering | Both |
| DBG-004 | medium | script | `scripts/anchor.py` | cmd_set with --file and positional topic args: --file silently wins | Both |
| DBG-005 | medium | script | `scripts/anchor.py` | cmd_update --refine skips blank lines between header and content, potentially losing structure | Both |
| DBG-006 | medium | integration | `SKILL.md` | Hardcoded absolute script paths reduce portability | manus |
| DBG-007 | medium | integration | `SKILL.md` | Implicit trigger detection not automated or specified how to implement | manus |
| DBG-008 | low | script | `scripts/anchor.py` | Timestamp update regex may replace unintended content | manus |

## Detailed Findings

### DBG-001: _insert_into_section never returns None on section-not-found but callers check for None [CONSENSUS]

**Severity:** medium | **Category:** script | **Confidence:** 1.0
**File:** `scripts/anchor.py` L108-L140

**Problematic Code:**
```python
    return "\n".join(new_lines)
```

**Problem:** The function _insert_into_section always returns a string (the joined new_lines). If the section_header is NOT found in the content, it simply returns the original content unchanged — it never returns None. However, callers in cmd_update check `if result:` which would be False for an empty string but never None. This isn't a crash bug, but the logic is misleading: if the section header doesn't exist (e.g., typo or format mismatch), the function silently returns the unmodified content and the caller treats it as success. The real issue is that if `## Key Objectives` exists in content but with slightly different formatting (e.g., trailing space), the `line.strip() == section_header` check will fail, and the item won't be inserted, but the code will still report success and set `modified = True`.

**Fix:**
```python
This is a minor logic gap. The function works correctly for the expected case. To make it more robust, add a flag to track whether insertion actually happened:

```python
def _insert_into_section(content, section_header, new_item, numbered=False):
    lines = content.split("\n")
    new_lines = []
    in_section = False
    inserted = False
    section_items = 0

    for line in lines:
        if line.strip() == section_header:
            in_section = True
            new_lines.append(line)
            continue
        # ... rest of logic ...

    # If section was the last one (no subsequent ##)
    if in_section and not inserted:
        if numbered:
            new_lines.append(f"{section_items + 1}. {new_item}")
        else:
            new_lines.append(f"- {new_item}")
        new_lines.append("")
        inserted = True

    if not inserted:
        return None  # Section not found
    return "\n".join(new_lines)
```
```

---

### DBG-002: Fallback keyword relevance scoring produces misleading scores for short messages [CONSENSUS]

**Severity:** medium | **Category:** script | **Confidence:** 1.0
**File:** `scripts/anchor.py` L215-L222

**Problematic Code:**
```python
    if not message_keywords:
        overlap_pct = 0
    else:
        overlap = anchor_keywords & message_keywords
        overlap_pct = len(overlap) / len(message_keywords) * 100

    score = min(10, int(overlap_pct / 10))
```

**Problem:** The keyword-based fallback measures overlap as a percentage of message keywords that appear in the anchor. For a short message like `check "database schema"` where both words appear in the anchor, overlap_pct would be 100%, giving a score of 10/10. Conversely, a highly relevant but differently-worded message like `check "Should I normalize the tables?"` would score 0/10 because none of those specific words appear in the anchor. This makes the fallback unreliable for its stated purpose. However, this is explicitly labeled as a 'fallback' and the primary path uses AI scoring, so the impact is limited to environments without an OpenAI API key.

**Fix:**
```python
This is a known limitation of keyword-based approaches and is already mitigated by the AI-first design. The fallback label in the output makes this clear. No critical fix needed, but could add a disclaimer in the output:
```python
print(f"  Note: Keyword matching is approximate. Set OPENAI_API_KEY for AI-powered scoring.")
```
```

---

### DBG-003: SKILL.md description is comprehensive and well-structured for Manus triggering [CONSENSUS]

**Severity:** medium | **Category:** integration | **Confidence:** 1.0
**File:** `SKILL.md` L1-L5

**Problematic Code:**
```python
N/A
```

**Problem:** This is not a finding — noting that the skill description, activation triggers, and instructions are well-designed for Manus integration. The YAML frontmatter is valid, paths are consistent, and the behavioral instructions are clear and actionable.

**Fix:**
```python
N/A
```

---

### DBG-004: cmd_set with --file and positional topic args: --file silently wins [CONSENSUS]

**Severity:** medium | **Category:** script | **Confidence:** 0.9500000000000001
**File:** `scripts/anchor.py` L155-L160

**Problematic Code:**
```python
    if args.file:
        try:
            with open(args.file, "r") as f:
                topic = f.read().strip()
        except (IOError, FileNotFoundError) as e:
            print(f"Error reading file: {e}", file=sys.stderr)
            sys.exit(1)
    elif args.topic:
        topic = " ".join(args.topic)
```

**Problem:** If a user provides both `--file somefile.txt` and a positional topic string, the file takes precedence silently. This is standard argparse behavior and documented implicitly by the if/elif structure, but could confuse users. This is a very minor UX issue, not a bug.

**Fix:**
```python
No fix needed — this is standard CLI precedence behavior.
```

---

### DBG-005: cmd_update --refine skips blank lines between header and content, potentially losing structure [CONSENSUS]

**Severity:** medium | **Category:** script | **Confidence:** 0.9
**File:** `scripts/anchor.py` L253-L270

**Problematic Code:**
```python
            if in_core and not replaced:
                if line.startswith("##"):
                    # End of Core Topic section — insert new topic
                    new_lines.append("")
                    new_lines.append(new_topic)
                    new_lines.append("")
                    new_lines.append(line)
                    replaced = True
                    in_core = False
                else:
                    # Skip ALL old topic content (blank and non-blank)
                    continue
```

**Problem:** The refine logic correctly replaces the Core Topic section content. However, if the Core Topic section is the last section in the document (no subsequent ## heading), the code enters the `if in_core and not replaced` block at the end and appends the new topic. This works. But there's a subtle issue: the `continue` statement skips ALL lines (including blank lines) after `## Core Topic` until the next `##` heading. This means if there's content between sections that doesn't start with `##` (like the `<!-- anchor-timestamp -->` line if it appears after Core Topic), it would be swallowed. In practice, the timestamp appears before Core Topic in the generated format, so this is unlikely to trigger with the default template, but could be an issue with manually edited anchor files.

**Fix:**
```python
This is acceptable for the default template but could be documented. No code change strictly necessary for the default flow.
```

---

### DBG-006: Hardcoded absolute script paths reduce portability [manus]

**Severity:** medium | **Category:** integration | **Confidence:** 0.9
**File:** `SKILL.md` L20-L30

**Problematic Code:**
```python
python3 /home/ubuntu/skills/context-anchor/scripts/anchor.py set "Building a SaaS analytics dashboard" \
```

**Problem:** The SKILL.md Quick Start examples use absolute paths (/home/ubuntu/skills/context-anchor/scripts/anchor.py) which will fail if Manus is installed or run in a different environment or user home directory. This reduces portability and integration ease.

**Fix:**
```python
Use relative paths or environment variables to locate the script, e.g., python3 ./scripts/anchor.py or parameterize the base path in the skill config.
```

---

### DBG-007: Implicit trigger detection not automated or specified how to implement [manus]

**Severity:** medium | **Category:** integration | **Confidence:** 0.85
**File:** `SKILL.md` L80-L90

**Problematic Code:**
```python
Detect **implicit** triggers: user frames session purpose with "today we're building X", "let's work on Y", "this project is about Z". Summarize the detected intent and confirm with the user before setting.
```

**Problem:** The skill description says to detect implicit triggers by user phrasing but does not specify how Manus should detect these triggers (e.g., regex, NLP intent detection). Without implementation details or a helper, Manus may not know how to reliably detect these triggers.

**Fix:**
```python
Add explicit instructions or helper functions/scripts for detecting implicit triggers, or provide example regex patterns or prompt templates for detection.
```

---

### DBG-008: Timestamp update regex may replace unintended content [manus]

**Severity:** low | **Category:** script | **Confidence:** 0.7
**File:** `scripts/anchor.py` L350-L370

**Problematic Code:**
```python
content = re.sub(
    r'<!-- anchor-timestamp -->.*',
    f'<!-- anchor-timestamp -->{timestamp} (updated)',
    content
)
```

**Problem:** The regex replaces everything after the anchor-timestamp comment on the same line, but if the timestamp comment appears multiple times or with trailing content, it may replace more than intended or miss updates if formatting changes.

**Fix:**
```python
Use a regex that matches only the timestamp comment and its immediate content, e.g., r'(<!-- anchor-timestamp -->)[^\n]*' and replace only that group.
```

---
