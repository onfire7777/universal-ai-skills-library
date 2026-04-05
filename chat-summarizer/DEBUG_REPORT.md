# Skill Debug Report: `chat-summarizer`

**Generated:** 2026-04-05 06:22 UTC
**Models:** Claude Opus 4.6 + Manus gpt-4.1-mini

## Model Status

| Model | Status | Findings |
|---|---|---|
| Claude Opus 4.6 | OK | 7 |
| Manus gpt-4.1-mini | OK | 4 |

## Overall Health: DEGRADED

**8 issues found** (3 confirmed by both models)

## Findings Summary

| ID | Severity | Category | File | Title | Consensus |
|---|---|---|---|---|---|
| DBG-001 | high | script | `scripts/format_summary.py` | Handoff prompt header uses string multiplication on a string literal instead of a repeated character | Both |
| DBG-002 | high | integration | `SKILL.md` | SKILL.md instructs Manus to run with --all which triggers --enrich, but --enrich will almost always fail | Both |
| DBG-003 | medium | script | `scripts/format_summary.py` | Section heading regex fails to match sections with mixed case or single-word names | Both |
| DBG-004 | medium | script | `scripts/format_summary.py` | enrich_with_ai uses non-existent model name 'gpt-4.1-mini' | claude |
| DBG-005 | medium | integration | `SKILL.md` | Output file paths in Step 4 assume output-dir is /home/ubuntu but script defaults to input file's parent directory | claude |
| DBG-006 | low | script | `scripts/format_summary.py` | extract_file_paths regex captures trailing punctuation and markdown table characters | claude |
| DBG-007 | low | integration | `scripts/format_summary.py` | AI enrichment requires OPENAI_API_KEY and openai package but no fallback or clear error propagation | manus |
| DBG-008 | low | script | `scripts/format_summary.py` | parse_summary discards content before the first ## heading (including the H1 title and Generated timestamp) | claude |

## Detailed Findings

### DBG-001: Handoff prompt header uses string multiplication on a string literal instead of a repeated character [CONSENSUS]

**Severity:** high | **Category:** script | **Confidence:** 1.0
**File:** `scripts/format_summary.py` L195-L197

**Problematic Code:**
```python
    header = (
        "=== AI HANDOFF PROMPT ===\n"
        "Paste this into a new chat to give the AI full context:\n"
        "─" * 60 + "\n"
    )
```

**Problem:** Due to Python operator precedence, `"─" * 60 + "\n"` is evaluated as `("─" * 60) + "\n"`, which produces the 60-character line plus newline. However, this expression is being implicitly concatenated with the previous string literal `"Paste this into a new chat to give the AI full context:\n"`. Python's implicit string concatenation only works between string literals, but `"─" * 60` is an expression, not a literal. So the actual behavior is: the implicit concatenation joins `"Paste this into a new chat to give the AI full context:\n"` with `"─"`, producing `"...context:\n─"`, and then `* 60` is applied to the ENTIRE concatenated string, and then `+ "\n"` is appended. This means the header line `Paste this into a new chat...` gets repeated 60 times, producing a massively bloated header instead of a clean separator line. This will make the handoff_prompt.txt file enormous and unusable.

**Fix:**
```python
    header = (
        "=== AI HANDOFF PROMPT ===\n"
        "Paste this into a new chat to give the AI full context:\n"
        + "─" * 60 + "\n"
    )
```

---

### DBG-002: SKILL.md instructs Manus to run with --all which triggers --enrich, but --enrich will almost always fail [CONSENSUS]

**Severity:** high | **Category:** integration | **Confidence:** 1.0
**File:** `SKILL.md` L93-L95

**Problematic Code:**
```python
python3 /home/ubuntu/skills/chat-summarizer/scripts/format_summary.py /home/ubuntu/chat_summary.md --all
```

**Problem:** The --all flag sets args.enrich = True, which calls enrich_with_ai(). This function requires the `openai` Python package AND the OPENAI_API_KEY environment variable. In most Manus sandbox environments, one or both of these will be missing. When the openai package is not installed, the function returns a string 'AI enrichment unavailable: openai package not installed.' — this is handled gracefully and won't crash, but it writes a useless ai_recommendations.md file. More importantly, the SKILL.md Step 4 tells Manus to deliver ai_recommendations.md as if it contains useful content ('AI-generated next steps (if --enrich was used)'), which will confuse the user. The parenthetical '(if --enrich was used)' is misleading since --all always enables --enrich. This is a silent failure that produces a misleading artifact.

**Fix:**
```python
Change the SKILL.md Step 3 command to not use --all by default, and instead use explicit flags:

```bash
python3 /home/ubuntu/skills/chat-summarizer/scripts/format_summary.py /home/ubuntu/chat_summary.md --json --handoff
```

And update Step 4 to only mention ai_recommendations.md conditionally:
4. `/home/ubuntu/ai_recommendations.md` — (only if `openai` package and OPENAI_API_KEY are available; run with `--enrich` to generate)
```

---

### DBG-003: Section heading regex fails to match sections with mixed case or single-word names [CONSENSUS]

**Severity:** medium | **Category:** script | **Confidence:** 1.0
**File:** `scripts/format_summary.py` L62-L72

**Problematic Code:**
```python
        heading_match = re.match(
            r"^##\s+(?:\d+\.\s+)?([A-Z][A-Z_]+(?:\s*[A-Z_]+)*)", line
        )
```

**Problem:** The regex `[A-Z][A-Z_]+` requires at least 2 uppercase characters. A section heading like `## DISCOVERIES` works, but the regex character class `[A-Z_]` does not match digits or other characters. More critically, the regex uses `[A-Z_]+` which won't match spaces between words properly — `(?:\s*[A-Z_]+)*` uses `\s*` (zero or more) which means it can greedily match across word boundaries. However, the real issue is that the summary template in SKILL.md uses headings like `## SESSION_METADATA` (with underscores), but if Manus writes `## Session Metadata` or `## Session_Metadata` (mixed case), the regex will fail to match because it requires all uppercase `[A-Z]`. Since Manus is an AI generating the markdown, case variation is a realistic scenario.

**Fix:**
```python
Make the regex case-insensitive or normalize:
        heading_match = re.match(
            r"^##\s+(?:\d+\.\s+)?([A-Za-z][A-Za-z_]+(?:\s*[A-Za-z_]+)*)", line
        )
        if heading_match:
            if current_section:
                sections[current_section] = "\n".join(current_content).strip()
            current_section = heading_match.group(1).strip().replace(" ", "_").upper()
            current_content = []
```

---

### DBG-004: enrich_with_ai uses non-existent model name 'gpt-4.1-mini' [claude]

**Severity:** medium | **Category:** script | **Confidence:** 0.72
**File:** `scripts/format_summary.py` L230-L232

**Problematic Code:**
```python
        response = client.chat.completions.create(
            model="gpt-4.1-mini",
```

**Problem:** The model name 'gpt-4.1-mini' is not a standard OpenAI model identifier. The likely intended model is 'gpt-4o-mini' or 'gpt-4.1-mini' (if this is a future/internal model). If this model doesn't exist in the API, the call will raise an exception like 'The model `gpt-4.1-mini` does not exist'. While the exception is caught by the broad `except Exception as e` block, it means AI enrichment will always fail even when the API key and package are available, producing a useless error message in the output file.

**Fix:**
```python
Change to a known model name:
            model="gpt-4o-mini",
```

---

### DBG-005: Output file paths in Step 4 assume output-dir is /home/ubuntu but script defaults to input file's parent directory [claude]

**Severity:** medium | **Category:** integration | **Confidence:** 0.7
**File:** `SKILL.md` L99-L103

**Problematic Code:**
```python
1. `/home/ubuntu/chat_summary.md` — the full structured summary
2. `/home/ubuntu/handoff_prompt.txt` — dense single-block for pasting into new chats
3. `/home/ubuntu/chat_summary.json` — machine-parseable version
4. `/home/ubuntu/ai_recommendations.md` — AI-generated next steps
```

**Problem:** The script's --output-dir defaults to `summary_path.parent`, which is the parent directory of the input file. Since the input is `/home/ubuntu/chat_summary.md`, the parent is `/home/ubuntu/`, so the paths happen to be correct in this case. However, if Manus writes the summary to a different location (e.g., a project directory), the output files would be written there instead, but Manus would still look for them at `/home/ubuntu/`. This is fragile but works for the documented happy path.

**Fix:**
```python
No immediate fix needed for the happy path, but for robustness, add --output-dir /home/ubuntu to the command in SKILL.md:
```bash
python3 /home/ubuntu/skills/chat-summarizer/scripts/format_summary.py /home/ubuntu/chat_summary.md --all --output-dir /home/ubuntu
```
```

---

### DBG-006: extract_file_paths regex captures trailing punctuation and markdown table characters [claude]

**Severity:** low | **Category:** script | **Confidence:** 0.88
**File:** `scripts/format_summary.py` L84-L90

**Problematic Code:**
```python
        r"`(/[^\s`]+)`",           # Backtick-wrapped absolute paths
        r"(/home/\S+)",            # /home/ paths
        r"(/mnt/\S+)",             # /mnt/ paths
```

**Problem:** The non-backtick patterns `/home/\S+` and `/mnt/\S+` will greedily match trailing markdown table pipes `|`, commas, parentheses, and other punctuation. While there's a `rstrip(".,;:)")` call, it doesn't strip `|` which is common in the ARTIFACTS_REGISTRY table format. For example, from a table row like `| /home/ubuntu/app.py | Main app |`, the regex would capture `/home/ubuntu/app.py` followed by trailing content. The backtick pattern is fine, but the bare path patterns will produce false path references that then generate spurious 'path does not exist' warnings.

**Fix:**
```python
Add `|` to the rstrip characters:
            path = match.group(1).rstrip(".,;:)|")
```

---

### DBG-007: AI enrichment requires OPENAI_API_KEY and openai package but no fallback or clear error propagation [manus]

**Severity:** low | **Category:** integration | **Confidence:** 0.8
**File:** `scripts/format_summary.py` L230-L240

**Problematic Code:**
```python
try:
    from openai import OpenAI
except ImportError:
    return "AI enrichment unavailable: openai package not installed."

api_key = os.environ.get("OPENAI_API_KEY")
if not api_key:
    return "AI enrichment unavailable: OPENAI_API_KEY not set."
```

**Problem:** If the openai package is missing or the API key is not set, the enrichment silently returns a string message instead of raising an error or logging. Manus may not detect this as a failure and may assume enrichment succeeded.

**Fix:**
```python
Raise exceptions or log errors explicitly so Manus can detect enrichment failure. Alternatively, document clearly that enrichment is optional and may silently skip if dependencies are missing.
```

---

### DBG-008: parse_summary discards content before the first ## heading (including the H1 title and Generated timestamp) [claude]

**Severity:** low | **Category:** script | **Confidence:** 0.75
**File:** `scripts/format_summary.py` L55-L76

**Problematic Code:**
```python
    for line in text.splitlines():
        # Match ## SECTION_NAME or ## 1. SECTION_NAME or ## SECTION_METADATA
        heading_match = re.match(
            r"^##\s+(?:\d+\.\s+)?([A-Z][A-Z_]+(?:\s*[A-Z_]+)*)", line
        )
        if heading_match:
            ...
        elif current_section:
            current_content.append(line)
```

**Problem:** The `# Chat Session Summary` H1 heading and `Generated: [timestamp]` line appear before any `##` section. Since `current_section` is None at that point, these lines are silently discarded. The timestamp is potentially useful metadata that could be included in the JSON output. This is a minor data loss issue — the information isn't used anywhere downstream, but it means the JSON export lacks the generation timestamp from the source document.

**Fix:**
```python
This is a minor issue. If desired, add pre-section content capture:
    preamble_lines = []
    ...
        elif current_section:
            current_content.append(line)
        else:
            preamble_lines.append(line)
    sections['_preamble'] = '\n'.join(preamble_lines).strip()
```

---
