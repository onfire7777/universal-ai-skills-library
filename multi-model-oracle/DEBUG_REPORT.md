# Skill Debug Report: `multi-model-oracle`

**Generated:** 2026-04-05 06:41 UTC
**Models:** Claude Opus 4.6 + Manus gpt-4.1-mini

## Model Status

| Model | Status | Findings |
|---|---|---|
| Claude Opus 4.6 | OK | 11 |
| Manus gpt-4.1-mini | OK | 8 |

## Overall Health: DEGRADED

**12 issues found** (7 confirmed by both models)

## Findings Summary

| ID | Severity | Category | File | Title | Consensus |
|---|---|---|---|---|---|
| DBG-001 | high | script | `scripts/oracle.py` | requests module imported inside function but not at top level — no install check or graceful error | Both |
| DBG-002 | high | script | `scripts/oracle.py` | OpenRouter response parsing assumes 'choices' key exists — will crash with KeyError on unexpected response format | Both |
| DBG-003 | high | integration | `SKILL.md` | SKILL.md references --stdin flag but doesn't document it in the workflow steps | Both |
| DBG-004 | high | structural | `SKILL.md` | Model names in description reference non-existent models (Claude Opus 4.6, GPT-5.4) | claude |
| DBG-005 | medium | script | `scripts/oracle.py` | RETRY_DELAYS has only 2 elements but is indexed by attempt number — IndexError if MAX_RETRIES is increased | Both |
| DBG-006 | medium | script | `scripts/oracle.py` | run_oracle uses `str | None` type hint syntax requiring Python 3.10+ | Both |
| DBG-007 | medium | integration | `scripts/oracle.py` | Progress/status messages printed to stdout mix with the final output, making it hard for Manus to parse | Both |
| DBG-008 | medium | script | `scripts/oracle.py` | merge_responses passes model names to the merge LLM despite instructions saying not to reference models | Both |
| DBG-009 | medium | integration | `scripts/oracle.py` | OPENAI_API_KEY environment variable required but not checked | manus |
| DBG-010 | medium | robustness | `scripts/oracle.py` | Merge stage has no retry logic — single failure loses the entire merge | claude |
| DBG-011 | medium | robustness | `scripts/oracle.py` | prompt_engineer makes 3 sequential API calls, not parallel — slow and no timeout | claude |
| DBG-012 | low | robustness | `scripts/oracle.py` | Intent detection 'code' pattern is overly broad — matches 'create a story' as code intent | claude |

## Detailed Findings

### DBG-001: requests module imported inside function but not at top level — no install check or graceful error [CONSENSUS]

**Severity:** high | **Category:** script | **Confidence:** 1.0
**File:** `scripts/oracle.py` L1-L10

**Problematic Code:**
```python
import requests
```

**Problem:** The `requests` module is imported inside `query_openrouter()` at runtime. If `requests` is not installed, the function will raise an ImportError that is caught by the generic `except Exception as e` handler, but the error message will be cryptic ('No module named requests'). More importantly, the SKILL.md Requirements section lists `requests` as needed but there's no check or installation step. The `openai` package has a graceful fallback in `prompt_engineer()` but `requests` does not get the same treatment in `query_openrouter()` — it will fail on the first attempt and then retry twice (wasting 20 seconds on sleep) before finally returning an error.

**Fix:**
```python
Add `requests` to the top-level imports with a try/except, or add an early check in `main()` that verifies both `requests` and `openai` are importable before starting the pipeline:

```python
def main():
    # Early dependency check
    missing = []
    try:
        import requests
    except ImportError:
        missing.append('requests')
    try:
        import openai
    except ImportError:
        missing.append('openai')
    if missing:
        print(f"ERROR: Missing required packages: {', '.join(missing)}", file=sys.stderr)
        print("Install with: pip install " + ' '.join(missing), file=sys.stderr)
        sys.exit(1)
    ...
```
```

---

### DBG-002: OpenRouter response parsing assumes 'choices' key exists — will crash with KeyError on unexpected response format [CONSENSUS]

**Severity:** high | **Category:** script | **Confidence:** 1.0
**File:** `scripts/oracle.py` L168-L175

**Problematic Code:**
```python
data = resp.json()
content = data["choices"][0]["message"]["content"].strip()
```

**Problem:** If the OpenRouter API returns a 200 status but with an error body (which some API proxies do), or if the response JSON structure is different than expected (e.g., missing 'choices' key, empty 'choices' array, or 'content' is None), this line will raise a KeyError, IndexError, or AttributeError. The exception is caught by the outer `except Exception` but only after the code has already passed the `resp.status_code == 200` check, so it will be treated as a retryable error and waste time on retries for a structural issue.

**Fix:**
```python
Add defensive parsing:
```python
if resp.status_code == 200:
    data = resp.json()
    choices = data.get("choices")
    if not choices or not isinstance(choices, list) or len(choices) == 0:
        return {"model": model_id, "success": False,
                "error": f"Unexpected response structure: no choices in {str(data)[:200]}"}
    content = (choices[0].get("message", {}).get("content") or "").strip()
    if not content:
        return {"model": model_id, "success": False,
                "error": "Empty content in response"}
    usage = data.get("usage", {})
    return {...}
```
```

---

### DBG-003: SKILL.md references --stdin flag but doesn't document it in the workflow steps [CONSENSUS]

**Severity:** high | **Category:** integration | **Confidence:** 1.0
**File:** `SKILL.md` L38-L42

**Problematic Code:**
```python
echo "Explain quantum entanglement" | python3 oracle.py --stdin
```

**Problem:** The `--stdin` flag is shown in the script's docstring but not in the SKILL.md workflow section. Manus relies on SKILL.md instructions to know how to use the skill. While this won't cause a crash, it means Manus won't know about the `--stdin` option and may not use it when appropriate (e.g., for very long queries that could exceed shell argument limits).

**Fix:**
```python
Add a `--stdin` example to the SKILL.md workflow section:
```
For piped input:
```bash
echo "Your question" | python3 /home/ubuntu/skills/multi-model-oracle/scripts/oracle.py --stdin
```
```

---

### DBG-004: Model names in description reference non-existent models (Claude Opus 4.6, GPT-5.4) [claude]

**Severity:** high | **Category:** structural | **Confidence:** 0.92
**File:** `SKILL.md` L3

**Problematic Code:**
```python
description: Get the ultimate merged answer from the 3 best AI models (Anthropic Claude Opus 4.6, OpenAI GPT-5.4, Manus gpt-4.1-mini).
```

**Problem:** The model identifiers 'anthropic/claude-opus-4.6' and 'openai/gpt-5.4' do not exist on OpenRouter. These are fabricated model names. When the script sends requests to OpenRouter with these model IDs, OpenRouter will return a 404 or 422 error. Since 422 is in FATAL_STATUSES, these requests will fail immediately without retry. This means 2 of the 3 models will always fail, and the 'oracle' will effectively just be a single gpt-4.1-mini response. This is the most fundamental issue with the skill.

**Fix:**
```python
Replace with real model IDs. For example:
  ANTHROPIC_MODEL = "anthropic/claude-sonnet-4" 
  OPENAI_MODEL = "openai/gpt-4.1"
And update SKILL.md descriptions accordingly. Alternatively, make the model IDs configurable via environment variables or arguments so they can be updated without code changes.
```

---

### DBG-005: RETRY_DELAYS has only 2 elements but is indexed by attempt number — IndexError if MAX_RETRIES is increased [CONSENSUS]

**Severity:** medium | **Category:** script | **Confidence:** 1.0
**File:** `scripts/oracle.py` L46-L48

**Problematic Code:**
```python
MAX_RETRIES = 2
RETRY_DELAYS = [5, 15]
```

**Problem:** The retry loop uses `RETRY_DELAYS[attempt]` where attempt ranges from 0 to MAX_RETRIES-1 (i.e., 0 and 1). Currently this works because RETRY_DELAYS has exactly 2 elements and MAX_RETRIES is 2. But if someone changes MAX_RETRIES to 3 without updating RETRY_DELAYS, `RETRY_DELAYS[2]` will raise an IndexError. This is a latent bug / maintenance trap.

**Fix:**
```python
Use `min(attempt, len(RETRY_DELAYS) - 1)` as the index, or generate delays dynamically:
```python
delay = RETRY_DELAYS[min(attempt, len(RETRY_DELAYS) - 1)]
time.sleep(delay)
```
```

---

### DBG-006: run_oracle uses `str | None` type hint syntax requiring Python 3.10+ [CONSENSUS]

**Severity:** medium | **Category:** script | **Confidence:** 1.0
**File:** `scripts/oracle.py` L290-L295

**Problematic Code:**
```python
def run_oracle(
    query: str,
    show_individual: bool = False,
    show_prompts: bool = False,
    output_file: str | None = None,
    skip_engineering: bool = False,
) -> str:
```

**Problem:** The `str | None` union syntax was introduced in Python 3.10. On Python 3.8 or 3.9 (which are still common on Ubuntu systems, especially older Manus environments), this will raise a `TypeError: unsupported operand type(s) for |: 'type' and 'NoneType'` at import time, crashing the entire script before any code runs.

**Fix:**
```python
Use `Optional[str]` from typing, or use `from __future__ import annotations` at the top of the file:
```python
from __future__ import annotations
```
Or change to:
```python
from typing import Optional
...
def run_oracle(
    query: str,
    show_individual: bool = False,
    show_prompts: bool = False,
    output_file: Optional[str] = None,
    skip_engineering: bool = False,
) -> str:
```
```

---

### DBG-007: Progress/status messages printed to stdout mix with the final output, making it hard for Manus to parse [CONSENSUS]

**Severity:** medium | **Category:** integration | **Confidence:** 1.0
**File:** `scripts/oracle.py` L310-L340

**Problematic Code:**
```python
print("=" * 60)
print("MULTI-MODEL ORACLE")
...
print(f"  Detected intent: {intent}")
...
output = run_oracle(...)
print(output)
```

**Problem:** The `run_oracle()` function prints progress messages to stdout (e.g., 'Querying models in parallel...', 'Detected intent: code', etc.) and then `main()` also prints the final formatted output to stdout. This means the actual output is interleaved with progress messages. When Manus captures stdout, it will get a messy mix of progress logs and the actual answer. The progress messages should go to stderr, or the function should return the output without printing progress to stdout.

**Fix:**
```python
Change all progress `print()` calls inside `run_oracle()` and `query_all_models()` to use `file=sys.stderr`:
```python
print("[Stage 1] Prompt Engineering...", file=sys.stderr)
print(f"  Detected intent: {intent}", file=sys.stderr)
# ... etc for all progress messages
```
Some warning prints already use `file=sys.stderr` correctly (e.g., in `prompt_engineer` and `merge_responses`), but the main pipeline progress messages in `run_oracle` and `query_all_models` do not.
```

---

### DBG-008: merge_responses passes model names to the merge LLM despite instructions saying not to reference models [CONSENSUS]

**Severity:** medium | **Category:** script | **Confidence:** 1.0
**File:** `scripts/oracle.py` L247-L260

**Problematic Code:**
```python
model_label = r["model"].split("/")[-1] if "/" in r["model"] else r["model"]
response_sections.append(
    f"=== Response {i} (from {model_label}) ===")
```

**Problem:** The merge prompt includes model names like 'claude-opus-4.6' and 'gpt-5.4' in the response headers sent to the merge LLM. While the MERGE_SYSTEM prompt instructs the LLM not to reference individual models, providing the model names increases the chance the LLM will leak them into the merged output, violating the merge strategy in the reference doc ('Never reference the individual models in the final output'). This is a design inconsistency rather than a crash bug, but it can produce incorrect output.

**Fix:**
```python
Use anonymous labels instead:
```python
response_sections.append(
    f"=== Response {i} ==="
    f"\n{r['content']}"
)
```
```

---

### DBG-009: OPENAI_API_KEY environment variable required but not checked [manus]

**Severity:** medium | **Category:** integration | **Confidence:** 0.9
**File:** `scripts/oracle.py` L90-L110

**Problematic Code:**
```python
client = OpenAI()  # Uses pre-configured API key and base URL
```

**Problem:** The Manus model querying requires OPENAI_API_KEY but the script does not check or warn if it is missing, leading to confusing errors at runtime.

**Fix:**
```python
Add explicit checks for OPENAI_API_KEY environment variable at startup and print clear error if missing.
```

---

### DBG-010: Merge stage has no retry logic — single failure loses the entire merge [claude]

**Severity:** medium | **Category:** robustness | **Confidence:** 0.85
**File:** `scripts/oracle.py` L260-L280

**Problematic Code:**
```python
try:
    from openai import OpenAI
    client = OpenAI()
    response = client.chat.completions.create(
        ...
    )
    return response.choices[0].message.content.strip()
except Exception as e:
    # Fallback: concatenate with headers
```

**Problem:** The merge stage (Stage 3) makes a single API call with no retry logic. If this call fails due to a transient network error or rate limit, the fallback is simple concatenation which defeats the entire purpose of the oracle. Stage 2 has retry logic (MAX_RETRIES=2) but Stage 3 does not, even though the merge is arguably the most important call since it produces the final output.

**Fix:**
```python
Add retry logic similar to Stage 2:
```python
for attempt in range(MAX_RETRIES + 1):
    try:
        response = client.chat.completions.create(...)
        return response.choices[0].message.content.strip()
    except Exception as e:
        if attempt < MAX_RETRIES:
            time.sleep(RETRY_DELAYS[attempt])
            continue
        # Final fallback
        print(f"  [Merge] Warning: AI merge failed ({e}), using concatenation fallback",
              file=sys.stderr)
        ...
```
```

---

### DBG-011: prompt_engineer makes 3 sequential API calls, not parallel — slow and no timeout [claude]

**Severity:** medium | **Category:** robustness | **Confidence:** 0.8
**File:** `scripts/oracle.py` L108-L140

**Problematic Code:**
```python
for model_key, model_id in models_info:
    ...
    try:
        response = client.chat.completions.create(
            model=MANUS_MODEL,
            ...
        )
```

**Problem:** Stage 1 makes 3 sequential API calls to gpt-4.1-mini to generate prompt variants. Each call could take 10-30 seconds. With no explicit timeout on the OpenAI client, if the API is slow, this stage alone could take 90+ seconds before Stage 2 even begins. The SKILL.md says the pipeline uses parallel execution, but that only applies to Stage 2. Stage 1 is sequential and has no timeout protection.

**Fix:**
```python
Either parallelize the 3 prompt engineering calls using ThreadPoolExecutor (similar to Stage 2), or set a timeout on the OpenAI client:
```python
client = OpenAI(timeout=60.0)
```
This ensures the prompt engineering stage doesn't hang indefinitely.
```

---

### DBG-012: Intent detection 'code' pattern is overly broad — matches 'create a story' as code intent [claude]

**Severity:** low | **Category:** robustness | **Confidence:** 0.78
**File:** `scripts/oracle.py` L62-L80

**Problematic Code:**
```python
"code": [
    r"\b(write|create|build|implement|code|script|function|class|program|debug|fix)\b",
```

**Problem:** The first pattern for 'code' intent matches common verbs like 'write' and 'create' without requiring any code-related context. A query like 'Write a poem about nature' would match 'write' in the code patterns AND 'write...poem' in the creative patterns. Since scoring is by match count, the code intent could win if other code patterns also partially match. The creative patterns require a compound match (verb + noun) which is more specific but scores lower. This could cause misclassification of creative queries as code queries.

**Fix:**
```python
Make the first code pattern require code-related context, or increase the weight of compound matches:
```python
"code": [
    r"\b(implement|code|script|function|class|program|debug|fix)\b",
    r"\b(write|create|build)\b.*\b(code|script|function|program|api|app|tool|bot|service)\b",
    ...
]
```
```

---
