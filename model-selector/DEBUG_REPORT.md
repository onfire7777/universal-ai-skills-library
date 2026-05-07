# Skill Debug Report: `model-selector`

**Generated:** 2026-05-07 20:08 UTC
**Models:** Claude Opus 4.6 + Manus gpt-4.1-mini

## Model Status

| Model | Status | Findings |
|---|---|---|
| Claude Opus 4.6 | OK | 10 |
| Manus gpt-4.1-mini | OK | 8 |

## Overall Health: DEGRADED

**14 issues found** (4 confirmed by both models)

## Findings Summary

| ID | Severity | Category | File | Title | Consensus |
|---|---|---|---|---|---|
| DBG-001 | high | script | `scripts/model_selector.py` | refresh_cache rate limit check uses 'force' parameter incorrectly — force=True skips rate limit but cmd_refresh always passes force=True | claude |
| DBG-002 | medium | robustness | `scripts/model_selector.py` | atomic_write_json writes to temp file in same directory as target — fails if directory doesn't exist | Both |
| DBG-003 | medium | script | `scripts/model_selector.py` | cmd_recommend triggers network refresh via select_best_model despite being read-only | claude |
| DBG-004 | medium | robustness | `scripts/model_selector.py` | Ambiguous model name handling in cmd_set may cause abrupt exit | manus |
| DBG-005 | medium | script | `scripts/model_selector.py` | detect_category gives 'general' keywords equal weight, making general win over specialized categories with same score | Both |
| DBG-006 | medium | script | `scripts/model_selector.py` | cmd_recommend calls select_best_model which may return None model_id, then prints None without handling | claude |
| DBG-007 | medium | script | `scripts/model_selector.py` | get_rankings returns empty list when cache exists but no ranked models validate, causing silent fallback loss | claude |
| DBG-008 | medium | script | `scripts/model_selector.py` | cmd_auto calls select_best_model which refreshes cache, then calls get_rankings again redundantly | claude |
| DBG-009 | low | integration | `scripts/model_selector.py` | refresh_cache imports requests at function level — unclear error if package missing | Both |
| DBG-010 | low | script | `scripts/model_selector.py` | detect_category scoring gives multi-word keywords 1.5 points per occurrence but single-word keywords 1.5 * count — inconsistent weighting | Both |
| DBG-011 | low | robustness | `scripts/model_selector.py` | Broad exception swallowing in frontier model identification may hide errors | manus |
| DBG-012 | low | robustness | `scripts/model_selector.py` | Skipping malformed model entries without logging | manus |
| DBG-013 | low | robustness | `scripts/model_selector.py` | Pricing conversion may fail silently on invalid data | manus |
| DBG-014 | low | integration | `SKILL.md` | SKILL.md trigger phrases don't include 'recommend model' as two separate words matching the command | claude |

## Detailed Findings

### DBG-001: refresh_cache rate limit check uses 'force' parameter incorrectly — force=True skips rate limit but cmd_refresh always passes force=True [claude]

**Severity:** high | **Category:** script | **Confidence:** 0.88
**File:** `scripts/model_selector.py` L155-L160

**Problematic Code:**
```python
def refresh_cache(force=False):
    ...
    # Rate limit check
    if not force:
        allowed, remaining = check_rate_limit()
        if not allowed:
            print(f"  Rate limited: wait {remaining}s before refreshing again.")
            # Try to return stale cache
            stale = safe_read_json(CACHE_PATH)
            return stale
```

**Problem:** The `cmd_refresh` function calls `refresh_cache(force=True)`, which bypasses the rate limit entirely. This means a user can spam `/model-selector refresh` and hit the OpenRouter API unlimited times, violating the documented 'max 1 per minute' rate limit. The rate limit is only enforced for auto-refresh (when `get_or_refresh_cache` calls `refresh_cache()` without force). The SKILL.md explicitly states 'Rate limit: max 1 per minute' for the refresh command.

**Fix:**
```python
Always check rate limit, but allow force to override only the staleness check, not the rate limit:

def refresh_cache(force=False):
    import requests

    # Rate limit check — always enforced
    allowed, remaining = check_rate_limit()
    if not allowed:
        print(f"  Rate limited: wait {remaining}s before refreshing again.")
        stale = safe_read_json(CACHE_PATH)
        return stale

    # If not force, also check if cache is still fresh (handled by caller)
    print("  Refreshing model cache from OpenRouter...")
    ...
```

---

### DBG-002: atomic_write_json writes to temp file in same directory as target — fails if directory doesn't exist [CONSENSUS]

**Severity:** medium | **Category:** robustness | **Confidence:** 1.0
**File:** `scripts/model_selector.py` L56-L60

**Problematic Code:**
```python
def atomic_write_json(path, data):
    """Write JSON atomically using temp file + rename to prevent corruption."""
    dir_name = os.path.dirname(path) or "."
    tmp_path = None
    try:
        fd, tmp_path = tempfile.mkstemp(dir=dir_name, suffix=".tmp")
```

**Problem:** CONFIG_PATH and CACHE_PATH are in the home directory (`~/.model_selector_config.json`). The `os.path.expanduser` is called at module level, so `dir_name` will be `/home/ubuntu`. This directory should always exist. However, if for some reason the home directory doesn't exist or isn't writable, the error message from the IOError handler is clear. This is not a practical bug in the Manus sandbox environment.

**Fix:**
```python
N/A - home directory always exists in Manus sandbox
```

---

### DBG-003: cmd_recommend triggers network refresh via select_best_model despite being read-only [claude]

**Severity:** medium | **Category:** script | **Confidence:** 0.92
**File:** `scripts/model_selector.py` L283-L285

**Problematic Code:**
```python
model_id, category, reason = select_best_model(task_text=task_text)
```

**Problem:** The `recommend` command is documented as a read-only command that should complete under 1 second with no network calls. However, `select_best_model()` calls `get_rankings(category, use_refresh=True)`, which calls `get_or_refresh_cache()`, which may trigger a network request to OpenRouter if the cache is stale or missing. This violates the performance target of <1s for read-only commands and could cause unexpected delays or failures when the user just wants a recommendation without changing config.

**Fix:**
```python
Change `select_best_model` to accept a `use_refresh` parameter, or create a separate read-only path for recommend:

def select_best_model(task_text=None, category=None, use_refresh=True):
    """Select the best model for a task."""
    if category is None:
        category = detect_category(task_text or "")

    rankings = get_rankings(category, use_refresh=use_refresh)
    if not rankings:
        return None, category, "No models available"

    best = rankings[0]
    return best["id"], category, best["reason"]

Then in cmd_recommend:
    model_id, category, reason = select_best_model(task_text=task_text, use_refresh=False)
```

---

### DBG-004: Ambiguous model name handling in cmd_set may cause abrupt exit [manus]

**Severity:** medium | **Category:** robustness | **Confidence:** 0.9
**File:** `scripts/model_selector.py` L277-L295

**Problematic Code:**
```python
if len(matches) > 1:
    # DBG-010: Show descriptions to help disambiguate
    print(f"  Ambiguous model name. Did you mean one of:")
    lookup = cache["model_lookup"]
    for m in matches[:10]:
        desc = lookup.get(m, {}).get("description", "")[:60]
        ctx = lookup.get(m, {}).get("context_length", 0)
        print(f"    {m}  (ctx={ctx:,}) {desc}")
    print(f"\n  Provide a more specific model ID to disambiguate.")
    sys.exit(1)
```

**Problem:** Abruptly exiting with sys.exit(1) on ambiguous model names may be unfriendly for integration or scripting. Manus might prefer a structured error or suggestion instead of exit.

**Fix:**
```python
Replace sys.exit(1) with raising a custom exception or returning an error message so Manus can handle it gracefully.
```

---

### DBG-005: detect_category gives 'general' keywords equal weight, making general win over specialized categories with same score [CONSENSUS]

**Severity:** medium | **Category:** script | **Confidence:** 0.85
**File:** `scripts/model_selector.py` L218-L230

**Problematic Code:**
```python
best = max(scores, key=scores.get)
if scores[best] == 0:
    return "general"

# DBG-013: If 'general' tied with a specialized category, prefer specialized
if best == "general":
    specialized = {k: v for k, v in scores.items() if k != "general" and v == scores[best]}
    if specialized:
        best = max(specialized, key=specialized.get)
```

**Problem:** The tie-breaking logic only activates when `best == "general"`. But `max()` with a dict returns the first key with the maximum value in insertion order. If 'general' is the last category in CATEGORY_KEYWORDS (which it is), and another category has the same score, `max()` will return the other category first (since it appears earlier in iteration). So the tie-break code for `best == "general"` would only trigger if Python's `max()` happened to pick 'general' — which it would only do if 'general' appeared first in iteration. Since 'general' is last in CATEGORY_KEYWORDS, `max()` over `scores` (a dict preserving insertion order) would return the first category with the max score, not 'general'. So the tie-break code is actually dead code in practice. This is not a bug per se — it's defensive code that works correctly because of dict ordering. No real failure scenario.

**Fix:**
```python
N/A - behavior is correct due to dict insertion order, the defensive check is harmless
```

---

### DBG-006: cmd_recommend calls select_best_model which may return None model_id, then prints None without handling [claude]

**Severity:** medium | **Category:** script | **Confidence:** 0.75
**File:** `scripts/model_selector.py` L283-L298

**Problematic Code:**
```python
model_id, category, reason = select_best_model(task_text=task_text)

print(f"\n  Task: \"{task_text}\"")
print(f"  Category: {category}")
print(f"  Recommended: {model_id}")
print(f"  Reason: {reason}")
```

**Problem:** If `select_best_model` returns `(None, category, 'No models available')` — which happens when `rankings` is empty (e.g., `get_rankings` returns an empty list because cache validation filtered out all models AND default rankings were somehow empty) — the output will show 'Recommended: None'. While this is unlikely with the hardcoded DEFAULT_RANKINGS, it's inconsistent with `cmd_auto` which explicitly checks for None and exits. The `cmd_recommend` should handle this case too.

**Fix:**
```python
Add a None check after select_best_model:

    model_id, category, reason = select_best_model(task_text=task_text)

    if model_id is None:
        print(f"\n  Task: \"{task_text}\"")
        print(f"  Category: {category}")
        print(f"  ERROR: Could not determine a recommendation. Try running 'refresh' first.")
        sys.exit(1)
```

---

### DBG-007: get_rankings returns empty list when cache exists but no ranked models validate, causing silent fallback loss [claude]

**Severity:** medium | **Category:** script | **Confidence:** 0.3
**File:** `scripts/model_selector.py` L247-L252

**Problematic Code:**
```python
validated = []
for r in rankings:
    if r["id"] in lookup:
        validated.append(r)
if validated:
    rankings = validated
```

**Problem:** If the cache is populated but none of the default-ranked models exist in the OpenRouter catalog (e.g., model IDs changed, API returned different format), `validated` will be empty and the code falls through to return the original `rankings`. This is actually correct behavior (graceful fallback). However, there's a subtle issue: if the cache has a `model_lookup` that is an empty dict `{}`, the condition `cache.get('model_lookup')` is falsy for empty dict, so validation is skipped entirely. This is fine. No actual bug here upon closer inspection.

**Fix:**
```python
N/A - this is actually correct on re-examination
```

---

### DBG-008: cmd_auto calls select_best_model which refreshes cache, then calls get_rankings again redundantly [claude]

**Severity:** medium | **Category:** script | **Confidence:** 0.25
**File:** `scripts/model_selector.py` L270-L278

**Problematic Code:**
```python
model_id, category, reason = select_best_model(task_text=task_text)
...
# Show alternatives — read-only, no refresh needed
rankings = get_rankings(category, use_refresh=False)
```

**Problem:** This is not a bug but a minor inefficiency. `select_best_model` already calls `get_rankings(category, use_refresh=True)` which loads/refreshes the cache. Then `cmd_auto` calls `get_rankings` again with `use_refresh=False`. The second call re-reads the cache from disk. This is harmless but wasteful. Not a real bug.

**Fix:**
```python
N/A - not a real bug, just minor inefficiency
```

---

### DBG-009: refresh_cache imports requests at function level — unclear error if package missing [CONSENSUS]

**Severity:** low | **Category:** integration | **Confidence:** 1.0
**File:** `scripts/model_selector.py` L155-L162

**Problematic Code:**
```python
def refresh_cache(force=False):
    """Fetch latest models from OpenRouter and rebuild the cache."""
    import requests
```

**Problem:** If the `requests` package is not installed, the import will raise `ImportError` inside the function. This error will be caught by the top-level exception handler in `main()` and printed as a generic error message. While SKILL.md states requests is pre-installed in the Manus sandbox, the error message won't clearly tell the user to install requests. However, this is documented as a requirement and the sandbox has it pre-installed, so this is a very minor concern.

**Fix:**
```python
Add a more specific error message by catching ImportError:

def refresh_cache(force=False):
    try:
        import requests
    except ImportError:
        print("  ERROR: 'requests' package not installed. Run: pip install requests")
        return safe_read_json(CACHE_PATH)
```

---

### DBG-010: detect_category scoring gives multi-word keywords 1.5 points per occurrence but single-word keywords 1.5 * count — inconsistent weighting [CONSENSUS]

**Severity:** low | **Category:** script | **Confidence:** 0.85
**File:** `scripts/model_selector.py` L209-L215

**Problematic Code:**
```python
if " " in kw or "-" in kw or "/" in kw:
    # Multi-word/special keywords: check substring and give full bonus
    if kw in task_lower:
        score += 1.5
else:
    # Single-word keywords: use word boundary regex
    pattern = r"\b" + re.escape(kw) + r"\b"
    matches = re.findall(pattern, task_lower)
    if matches:
        score += len(matches) * 1.5
```

**Problem:** Single-word keywords get `count * 1.5` points (so if 'code' appears 3 times, that's 4.5 points), while multi-word keywords like 'ci/cd' or 'pros and cons' only get 1.5 points regardless of how many times they appear. This means repeating a single word in a task description can disproportionately skew category detection. For example, 'help help help help code' would score general=6.0 vs coding=1.5, even though the task is about code. However, this is a design choice rather than a bug — repeated words in natural language are rare, and the weighting generally works for typical inputs.

**Fix:**
```python
Consider capping single-word matches at 1 occurrence:
    matches = re.findall(pattern, task_lower)
    if matches:
        score += 1.5  # Count presence, not frequency
```

---

### DBG-011: Broad exception swallowing in frontier model identification may hide errors [manus]

**Severity:** low | **Category:** robustness | **Confidence:** 0.85
**File:** `scripts/model_selector.py` L209-L234

**Problematic Code:**
```python
for m in models:
    try:
        model_id = m.get("id")
        if not model_id:
            continue  # DBG-012: skip entries without an ID
        provider = model_id.split("/")[0]
        if provider in frontier_providers:
            try:
                price_out = float(m.get("pricing", {}).get("completion", "0") or "0")
            except (ValueError, TypeError):
                price_out = 0.0  # DBG-012: handle non-numeric pricing
            ctx = m.get("context_length", 0) or 0
            frontier.append({
                "id": model_id,
                "name": m.get("name", ""),
                "provider": provider,
                "context_length": ctx,
                "price_per_m_output": price_out * 1_000_000,
                "description": (m.get("description", "") or "")[:200],
            })
    except Exception:
        continue  # DBG-012: skip any malformed entry
```

**Problem:** Catching all exceptions and continuing silently may hide unexpected bugs or data issues.

**Fix:**
```python
Log or print a warning when skipping malformed entries due to exceptions to aid debugging.
```

---

### DBG-012: Skipping malformed model entries without logging [manus]

**Severity:** low | **Category:** robustness | **Confidence:** 0.8
**File:** `scripts/model_selector.py` L196-L204

**Problematic Code:**
```python
for m in models:
    model_id = m.get("id")
    if not model_id:
        continue  # DBG-011: skip entries without an ID
    model_lookup[model_id] = {
        "id": model_id,
        "name": m.get("name", model_id),
        "context_length": m.get("context_length", 0) or 0,
        "pricing_prompt": m.get("pricing", {}).get("prompt", "0") or "0",
        "pricing_completion": m.get("pricing", {}).get("completion", "0") or "0",
        "description": m.get("description", "") or "",
        "modality": m.get("architecture", {}).get("modality", "text->text"),
    }
```

**Problem:** Malformed entries without 'id' are silently skipped without any warning or logging, which may hide data issues.

**Fix:**
```python
Add a debug or warning print statement when skipping entries without 'id' to aid troubleshooting.
```

---

### DBG-013: Pricing conversion may fail silently on invalid data [manus]

**Severity:** low | **Category:** robustness | **Confidence:** 0.8
**File:** `scripts/model_selector.py` L345-L355

**Problematic Code:**
```python
try:
    pi = float(m.get("pricing_prompt") or "0") * 1_000_000
    po = float(m.get("pricing_completion") or "0") * 1_000_000
except (ValueError, TypeError):
    pi, po = 0.0, 0.0
```

**Problem:** If pricing data is malformed, it defaults silently to zero without notifying the user, which may mislead about actual costs.

**Fix:**
```python
Consider logging a warning when pricing data is invalid to alert maintainers or users.
```

---

### DBG-014: SKILL.md trigger phrases don't include 'recommend model' as two separate words matching the command [claude]

**Severity:** low | **Category:** integration | **Confidence:** 0.15
**File:** `SKILL.md` L8-L12

**Problematic Code:**
```python
"recommend model", "refresh models"
```

**Problem:** The trigger phrase 'recommend model' is listed, which matches the `/model-selector recommend` command. This is actually fine. No issue here.

**Fix:**
```python
N/A
```

---
