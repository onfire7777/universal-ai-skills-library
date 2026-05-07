# Skill Debug Report: `prompt-engineer`

**Generated:** 2026-05-07 20:09 UTC
**Models:** Claude Opus 4.6 + Manus gpt-4.1-mini

## Model Status

| Model | Status | Findings |
|---|---|---|
| Claude Opus 4.6 | OK | 5 |
| Manus gpt-4.1-mini | OK | 10 |

## Overall Health: DEGRADED

**12 issues found** (3 confirmed by both models)

## Findings Summary

| ID | Severity | Category | File | Title | Consensus |
|---|---|---|---|---|---|
| DBG-001 | high | integration | `SKILL.md` | Hardcoded script path in Quick Start section | manus |
| DBG-002 | high | integration | `scripts/optimize_prompt.py` | Stdin input mode hangs if no piped input is provided | manus |
| DBG-003 | medium | robustness | `scripts/optimize_prompt.py` | strip_code_fences only extracts first code block, misses non-anchored blocks | Both |
| DBG-004 | medium | script | `scripts/optimize_prompt.py` | OpenAI client initialized with potentially empty API key | Both |
| DBG-005 | medium | integration | `SKILL.md` | SKILL.md mentions --stdin will hang interactively but doesn't warn about isatty behavior | Both |
| DBG-006 | medium | script | `scripts/optimize_prompt.py` | Fallback to quick mode when only one model is available is only warned, not enforced | manus |
| DBG-007 | medium | script | `scripts/optimize_prompt.py` | call_openrouter does not handle HTTP errors with detailed messages | manus |
| DBG-008 | medium | script | `scripts/optimize_prompt.py` | Intent detection uses simple keyword matching without stemming or context | manus |
| DBG-009 | medium | script | `scripts/optimize_prompt.py` | Mode auto-detection may misclassify prompts due to simplistic keyword scoring | manus |
| DBG-010 | medium | script | `scripts/optimize_prompt.py` | merge_optimizations fallback returns longest valid result without quality check | manus |
| DBG-011 | medium | robustness | `scripts/optimize_prompt.py` | analyze_prompt passes error string to json.loads causing silent fallback | claude |
| DBG-012 | low | robustness | `scripts/optimize_prompt.py` | optimize_prompt defaults available_models to all models regardless of API key availability | claude |

## Detailed Findings

### DBG-001: Hardcoded script path in Quick Start section [manus]

**Severity:** high | **Category:** integration | **Confidence:** 0.99
**File:** `SKILL.md` L18-L22

**Problematic Code:**
```python
python3 /home/ubuntu/skills/prompt-engineer/scripts/optimize_prompt.py "your raw prompt here"
```

**Problem:** The Quick Start instructions hardcode an absolute path (/home/ubuntu/skills/prompt-engineer/scripts/optimize_prompt.py) which will not exist on most users' systems. This will cause confusion or failure when users try to run the script as instructed.

**Fix:**
```python
Replace the hardcoded absolute path with a relative or generic path, e.g., `python3 scripts/optimize_prompt.py "your raw prompt here"` or instruct users to run from the skill directory.
```

---

### DBG-002: Stdin input mode hangs if no piped input is provided [manus]

**Severity:** high | **Category:** integration | **Confidence:** 0.95
**File:** `scripts/optimize_prompt.py` L292-L298

**Problematic Code:**
```python
if args.stdin:
    if sys.stdin.isatty():
        print("Reading from stdin (press Ctrl+D when done):", file=sys.stderr)
    raw_prompt = sys.stdin.read().strip()
```

**Problem:** The code reads from stdin without timeout or fallback. If the user runs the script with --stdin but does not pipe input, the script will hang waiting for input indefinitely, which is noted in SKILL.md but not handled programmatically.

**Fix:**
```python
Add a timeout or detect empty input after a short wait and exit with an informative error message to prevent hanging.
```

---

### DBG-003: strip_code_fences only extracts first code block, misses non-anchored blocks [CONSENSUS]

**Severity:** medium | **Category:** robustness | **Confidence:** 1.0
**File:** `scripts/optimize_prompt.py` L139-L142

**Problematic Code:**
```python
match = re.match(r'^```(?:json)?\s*\n?(.*?)\n?```', text, re.DOTALL)
if match:
    return match.group(1).strip()
return text
```

**Problem:** The regex uses `re.match` which anchors at the start of the string. If the model returns any preamble text before the code fence (e.g., 'Here is the JSON:\n```json\n...```'), the match fails and the raw text (including the preamble and fences) is returned. This then causes `json.loads()` to fail in `analyze_prompt` and `compare_prompts`, falling back to the default analysis dict. While the fallback exists, it degrades quality silently. Using `re.search` would be more robust.

**Fix:**
```python
match = re.search(r'```(?:json)?\s*\n?(.*?)\n?```', text, re.DOTALL)
if match:
    return match.group(1).strip()
return text
```

---

### DBG-004: OpenAI client initialized with potentially empty API key [CONSENSUS]

**Severity:** medium | **Category:** script | **Confidence:** 1.0
**File:** `scripts/optimize_prompt.py` L87-L91

**Problematic Code:**
```python
def _get_openai_client():
    global _openai_client
    if _openai_client is None:
        _openai_client = OpenAI(api_key=OPENAI_API_KEY, timeout=180)
    return _openai_client
```

**Problem:** If OPENAI_API_KEY is empty string (the default from `os.environ.get('OPENAI_API_KEY', '')`), the OpenAI client is created with an empty key. This won't crash at creation time but will produce confusing auth errors at call time. The `available_models` logic in `main()` prevents 'mini' from being added when the key is empty, but `_get_openai_client()` could still be called if code paths change. More critically, if only OPENROUTER_API_KEY is set and the merge step picks 'mini' via `pick_model`, it would fail — but `pick_model` only picks from `available_models` so this is protected. The risk is low but the empty-string check is fragile.

**Fix:**
```python
def _get_openai_client():
    global _openai_client
    if _openai_client is None:
        if not OPENAI_API_KEY:
            raise RuntimeError("OPENAI_API_KEY not set but OpenAI client requested")
        _openai_client = OpenAI(api_key=OPENAI_API_KEY, timeout=180)
    return _openai_client
```

---

### DBG-005: SKILL.md mentions --stdin will hang interactively but doesn't warn about isatty behavior [CONSENSUS]

**Severity:** medium | **Category:** integration | **Confidence:** 1.0
**File:** `SKILL.md` L30-L30

**Problematic Code:**
```python
**Stdin Input**: Pipe input with `echo "prompt" | python3 scripts/optimize_prompt.py --stdin`. Note: `--stdin` requires piped input; it will hang if run interactively without piping.
```

**Problem:** The script actually handles this case with `if sys.stdin.isatty(): print('Reading from stdin (press Ctrl+D when done):', file=sys.stderr)`, so it won't truly 'hang' — it will prompt the user. The SKILL.md instruction is slightly misleading but not critically wrong. Manus might avoid using --stdin unnecessarily based on this warning, which is fine behavior.

**Fix:**
```python
**Stdin Input**: Pipe input with `echo "prompt" | python3 scripts/optimize_prompt.py --stdin`. If run without piped input, it will prompt for interactive input (end with Ctrl+D).
```

---

### DBG-006: Fallback to quick mode when only one model is available is only warned, not enforced [manus]

**Severity:** medium | **Category:** script | **Confidence:** 0.9
**File:** `scripts/optimize_prompt.py` L354-L361

**Problematic Code:**
```python
if args.depth != "quick" and len(available_models) < 2:
    print(f"Warning: Only {len(available_models)} model(s) available ({', '.join(available_models)}). Falling back to quick mode.", file=sys.stderr)
    args.depth = "quick"
```

**Problem:** The code prints a warning and sets args.depth to quick mode if fewer than two models are available, but this happens after parsing arguments. If the user explicitly requested a deeper mode, this silently overrides it with only a warning, which may confuse users or cause unexpected behavior.

**Fix:**
```python
Either enforce the fallback with an explicit message and exit or document clearly that the mode will be overridden. Alternatively, fail early if requested mode is incompatible with available models.
```

---

### DBG-007: call_openrouter does not handle HTTP errors with detailed messages [manus]

**Severity:** medium | **Category:** script | **Confidence:** 0.85
**File:** `scripts/optimize_prompt.py` L237-L252

**Problematic Code:**
```python
resp = requests.post(
    "https://openrouter.ai/api/v1/chat/completions",
    headers=headers, json=payload, timeout=180,
)
resp.raise_for_status()
data = resp.json()
if "error" in data:
    raise ValueError(f"API error: {str(data['error'])[:200]}")
```

**Problem:** If the HTTP response is an error (e.g., 401 Unauthorized, 429 Rate Limit), resp.raise_for_status() raises an HTTPError but the except block only returns a sanitized error string without logging or differentiating error types. This may obscure the root cause and complicate debugging.

**Fix:**
```python
Catch HTTPError separately, log status codes and messages, and provide more informative error messages to the user. Consider retrying on rate limits with backoff.
```

---

### DBG-008: Intent detection uses simple keyword matching without stemming or context [manus]

**Severity:** medium | **Category:** script | **Confidence:** 0.8
**File:** `scripts/optimize_prompt.py` L172-L179

**Problematic Code:**
```python
def detect_intent(prompt: str) -> str:
    prompt_lower = prompt.lower()
    scores = {}
    for category, info in INTENT_CATEGORIES.items():
        score = sum(1 for kw in info["keywords"] if re.search(r'\b' + re.escape(kw) + r'\b', prompt_lower))
        scores[category] = score
    best = max(scores, key=scores.get)
    return best if scores[best] > 0 else "general"
```

**Problem:** The intent detection relies on exact keyword matches with word boundaries, which may miss relevant intents if the user uses synonyms, plural forms, or related words. This can cause misclassification and suboptimal prompt optimization.

**Fix:**
```python
Enhance intent detection with stemming, lemmatization, or use a small ML classifier to improve accuracy. Alternatively, expand keyword lists or add fuzzy matching.
```

---

### DBG-009: Mode auto-detection may misclassify prompts due to simplistic keyword scoring [manus]

**Severity:** medium | **Category:** script | **Confidence:** 0.8
**File:** `scripts/optimize_prompt.py` L187-L198

**Problematic Code:**
```python
def detect_mode(prompt: str) -> str:
    prompt_lower = prompt.lower()
    agent_signals = ["build", "create", "deploy", "research", "analyze", "file", "website",
                     "app", "automate", "organize", "manage", "monitor", "download",
                     "scrape", "install", "configure", "set up", "generate"]
    project_signals = ["recurring", "weekly", "daily", "every time", "template",
                       "standard", "always", "consistent", "reusable", "team"]
    agent_score = sum(1 for s in agent_signals if re.search(r'\b' + re.escape(s) + r'\b', prompt_lower))
    project_score = sum(1 for s in project_signals if re.search(r'\b' + re.escape(s) + r'\b', prompt_lower))
    if project_score >= 2:
        return "project"
    if agent_score >= 1:
        return "agent"
    return "chat"
```

**Problem:** The mode detection uses simple keyword counts and thresholds that may misclassify prompts, especially if keywords appear in different contexts or multiple modes are relevant. For example, a prompt mentioning 'create' and 'weekly' might be ambiguous but forced into 'project' mode only if two project signals appear.

**Fix:**
```python
Improve mode detection with more nuanced NLP techniques or allow user override with clearer feedback. Consider weighting keywords or combining with intent detection.
```

---

### DBG-010: merge_optimizations fallback returns longest valid result without quality check [manus]

**Severity:** medium | **Category:** script | **Confidence:** 0.8
**File:** `scripts/optimize_prompt.py` L404-L408

**Problematic Code:**
```python
if is_error_result(merged):
    merged = max(valid_results.values(), key=len)
```

**Problem:** If merging multiple optimized prompts fails, the fallback returns the longest valid result by length, which may not be the best quality or most coherent prompt. This heuristic can produce suboptimal output.

**Fix:**
```python
Implement a better fallback strategy, such as selecting the prompt with the highest clarity or specificity score if available, or returning a default optimized prompt with a warning.
```

---

### DBG-011: analyze_prompt passes error string to json.loads causing silent fallback [claude]

**Severity:** medium | **Category:** robustness | **Confidence:** 0.75
**File:** `scripts/optimize_prompt.py` L280-L283

**Problematic Code:**
```python
result = call_model(model, ANALYSIS_SYSTEM_PROMPT, raw_prompt, temperature=0.3)
if is_error_result(result):
    print(f"  Warning: Analysis model failed: {result}", file=sys.stderr)
try:
    return json.loads(strip_code_fences(result))
```

**Problem:** When `is_error_result(result)` is True, the code prints a warning but then still attempts `json.loads(strip_code_fences(result))` on the error string. This will always fail and fall through to the except block returning the default dict. While this works functionally, it's wasteful and confusing — the warning is printed but then the code proceeds as if it might succeed. More importantly, if the error message happens to contain valid JSON (unlikely but possible with some API error formats), it could return garbage analysis data.

**Fix:**
```python
result = call_model(model, ANALYSIS_SYSTEM_PROMPT, raw_prompt, temperature=0.3)
if is_error_result(result):
    print(f"  Warning: Analysis model failed: {result}", file=sys.stderr)
    return {
        "intent": detect_intent(raw_prompt),
        "mode": detect_mode(raw_prompt),
        "clarity_score": 5,
        "specificity_score": 5,
        "actionability_score": 5,
        "missing_elements": ["Analysis model failed"],
        "ambiguities": [],
        "strengths": [],
        "manus_tools_needed": [],
        "suggested_structure": "standard",
        "complexity_level": "moderate",
    }
try:
    return json.loads(strip_code_fences(result))
```

---

### DBG-012: optimize_prompt defaults available_models to all models regardless of API key availability [claude]

**Severity:** low | **Category:** robustness | **Confidence:** 0.7
**File:** `scripts/optimize_prompt.py` L310-L312

**Problematic Code:**
```python
if available_models is None:
    available_models = list(MODELS.keys())
```

**Problem:** If `optimize_prompt` is called without the `available_models` parameter (e.g., from external code or future refactoring), it defaults to all models including those whose API keys may not be set. This would cause API call failures. In the current `main()` flow, `available_models` is always passed explicitly, so this is only a latent issue for future callers.

**Fix:**
```python
if available_models is None:
    available_models = []
    if OPENROUTER_API_KEY:
        available_models.extend(["opus", "gpt"])
    if OPENAI_API_KEY:
        available_models.append("mini")
```

---
