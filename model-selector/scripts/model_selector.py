#!/usr/bin/env python3
"""Model Selector — Set, auto-select, or toggle the default AI model for Manus.

Commands:
    set <model_id>       Set a specific model as the default backbone LLM
    auto                 Auto-select the best model for the current task
    auto <task_desc>     Auto-select the best model for a described task
    on                   Enable model-selector (Manus reads config at task start)
    off                  Disable model-selector (Manus uses its default model)
    status               Show current configuration and selected model
    refresh              Refresh the cached model leaderboard from OpenRouter
    list [category]      List top models, optionally filtered by category
    recommend <task>     Recommend the best model for a task without setting it

Config: ~/.model_selector_config.json
Cache:  ~/.model_selector_cache.json
"""

import argparse
import json
import os
import re
import sys
import tempfile
import time
from datetime import datetime, timezone

# ---------------------------------------------------------------------------
# Paths
# ---------------------------------------------------------------------------
CONFIG_PATH = os.path.expanduser("~/.model_selector_config.json")
CACHE_PATH = os.path.expanduser("~/.model_selector_cache.json")
RATE_LIMIT_PATH = os.path.expanduser("~/.model_selector_last_refresh")
CACHE_MAX_AGE_HOURS = 24
RATE_LIMIT_SECONDS = 60  # Min seconds between refreshes

# ---------------------------------------------------------------------------
# Default config
# ---------------------------------------------------------------------------
DEFAULT_CONFIG = {
    "enabled": True,
    "mode": "auto",
    "manual_model": None,
    "auto_category": None,
    "last_selected_model": None,
    "last_selected_reason": None,
    "updated_at": None,
}

# ---------------------------------------------------------------------------
# Task category definitions — dynamically matched via keywords
# ---------------------------------------------------------------------------
CATEGORY_KEYWORDS = {
    "coding": [
        "code", "program", "script", "debug", "function", "class", "api",
        "build", "develop", "software", "app", "website", "deploy", "refactor",
        "test", "bug", "compile", "syntax", "algorithm", "database", "sql",
        "python", "javascript", "typescript", "react", "node", "html", "css",
        "git", "docker", "kubernetes", "ci/cd", "backend", "frontend",
    ],
    "reasoning": [
        "analyze", "reason", "think", "logic", "strategy", "plan", "evaluate",
        "compare", "decide", "assess", "framework", "architecture", "design",
        "system", "complex", "tradeoff", "pros and cons", "optimize",
        "problem-solving", "critical thinking", "decision",
    ],
    "research": [
        "research", "investigate", "find", "search", "discover", "study",
        "survey", "literature", "paper", "academic", "source", "citation",
        "evidence", "data", "statistics", "benchmark", "report", "fact-check",
    ],
    "creative": [
        "write", "story", "creative", "poem", "essay", "blog", "article",
        "content", "copy", "narrative", "fiction", "marketing", "brand",
        "slogan", "tagline", "tone", "voice", "style",
    ],
    "math": [
        "math", "calculate", "equation", "formula", "proof", "theorem",
        "algebra", "calculus", "statistics", "probability", "geometry",
        "numerical", "computation", "solve",
    ],
    "general": [
        "help", "explain", "summarize", "translate", "convert", "format",
        "organize", "define", "overview",
    ],
}

# ---------------------------------------------------------------------------
# Model tier rankings per category (curated from April 2026 benchmarks)
# ---------------------------------------------------------------------------
DEFAULT_RANKINGS = {
    "coding": [
        {"id": "anthropic/claude-opus-4.6", "reason": "80.8% SWE-bench, best for complex codebases"},
        {"id": "openai/gpt-5.4", "reason": "57.7% SWE-bench Pro, strong terminal/tool use"},
        {"id": "google/gemini-3.1-pro-preview", "reason": "Top coding benchmarks, 1M context"},
        {"id": "x-ai/grok-4.20", "reason": "2M context, strong multi-agent coding"},
        {"id": "deepseek/deepseek-v3.2", "reason": "Excellent value, strong coding at low cost"},
    ],
    "reasoning": [
        {"id": "anthropic/claude-opus-4.6", "reason": "Strongest reasoning and analysis"},
        {"id": "openai/gpt-5.4", "reason": "Excellent structured reasoning"},
        {"id": "google/gemini-3.1-pro-preview", "reason": "Strong analytical capabilities"},
        {"id": "x-ai/grok-4.20", "reason": "Deep reasoning with massive context"},
        {"id": "deepseek/deepseek-r1-0528", "reason": "Dedicated reasoning model"},
    ],
    "research": [
        {"id": "google/gemini-3.1-pro-preview", "reason": "Best for research with 1M context"},
        {"id": "anthropic/claude-opus-4.6", "reason": "Excellent synthesis and analysis"},
        {"id": "openai/gpt-5.4", "reason": "Strong factual accuracy"},
        {"id": "x-ai/grok-4.20", "reason": "2M context for massive document analysis"},
        {"id": "deepseek/deepseek-v3.2", "reason": "Cost-effective research at scale"},
    ],
    "creative": [
        {"id": "anthropic/claude-opus-4.6", "reason": "Best creative writing quality"},
        {"id": "anthropic/claude-sonnet-4.6", "reason": "Excellent creative at lower cost"},
        {"id": "openai/gpt-5.4", "reason": "Strong narrative and style control"},
        {"id": "google/gemini-3.1-pro-preview", "reason": "Good creative with multimodal"},
        {"id": "x-ai/grok-4.20", "reason": "Creative with fewer restrictions"},
    ],
    "math": [
        {"id": "openai/gpt-5.4", "reason": "Strongest math benchmarks"},
        {"id": "anthropic/claude-opus-4.6", "reason": "Excellent mathematical reasoning"},
        {"id": "deepseek/deepseek-r1-0528", "reason": "Dedicated reasoning, strong math"},
        {"id": "google/gemini-3.1-pro-preview", "reason": "Top math benchmark scores"},
        {"id": "x-ai/grok-4.20", "reason": "Strong computational abilities"},
    ],
    "general": [
        {"id": "anthropic/claude-opus-4.6", "reason": "Best overall general capability"},
        {"id": "openai/gpt-5.4", "reason": "Excellent all-around performance"},
        {"id": "google/gemini-3.1-pro-preview", "reason": "Strong general with best value"},
        {"id": "anthropic/claude-sonnet-4.6", "reason": "Great balance of quality and speed"},
        {"id": "x-ai/grok-4.20", "reason": "Versatile with massive context"},
    ],
}


# ===========================================================================
# Atomic file I/O helpers
# ===========================================================================

def atomic_write_json(path, data):
    """Write JSON atomically using temp file + rename to prevent corruption."""
    dir_name = os.path.dirname(path) or "."
    tmp_path = None  # DBG-007: init before try to avoid UnboundLocalError
    try:
        fd, tmp_path = tempfile.mkstemp(dir=dir_name, suffix=".tmp")
        with os.fdopen(fd, "w") as f:
            json.dump(data, f, indent=2)
        os.replace(tmp_path, path)
    except OSError as e:
        # Clean up temp file on failure
        if tmp_path:  # DBG-007: only unlink if tmp_path was assigned
            try:
                os.unlink(tmp_path)
            except OSError:
                pass
        raise IOError(f"Failed to write {path}: {e}")


def safe_read_json(path, default=None):
    """Read JSON with error handling, returning default on failure."""
    if not os.path.isfile(path):
        return default
    try:
        with open(path) as f:
            return json.load(f)
    except (json.JSONDecodeError, IOError, OSError):
        return default


# ===========================================================================
# Config helpers
# ===========================================================================

def load_config():
    """Load config from disk, or return defaults."""
    cfg = safe_read_json(CONFIG_PATH)
    if cfg is None:
        return dict(DEFAULT_CONFIG)
    # Merge with defaults for any missing keys
    for k, v in DEFAULT_CONFIG.items():
        cfg.setdefault(k, v)
    return cfg


def save_config(cfg):
    """Persist config to disk atomically."""
    cfg["updated_at"] = datetime.now(timezone.utc).isoformat()
    try:
        atomic_write_json(CONFIG_PATH, cfg)
    except IOError as e:
        print(f"  ERROR: Could not save config: {e}")
        print(f"  Check file permissions for {CONFIG_PATH}")
        sys.exit(1)


# ===========================================================================
# Rate limiting
# ===========================================================================

def check_rate_limit():
    """Check if refresh is rate-limited. Returns (allowed, seconds_remaining)."""
    if not os.path.isfile(RATE_LIMIT_PATH):
        return True, 0
    try:
        with open(RATE_LIMIT_PATH) as f:
            last_refresh = float(f.read().strip())
        elapsed = time.time() - last_refresh
        if elapsed < RATE_LIMIT_SECONDS:
            return False, int(RATE_LIMIT_SECONDS - elapsed)
        return True, 0
    except (ValueError, IOError):
        return True, 0


def update_rate_limit():
    """Record the current time as the last refresh time."""
    try:
        with open(RATE_LIMIT_PATH, "w") as f:
            f.write(str(time.time()))
    except IOError:
        pass  # Non-critical


# ===========================================================================
# Cache helpers
# ===========================================================================

def load_cache():
    """Load cached model data, or return None if stale/missing."""
    cache = safe_read_json(CACHE_PATH)
    if cache is None:
        return None
    # Check age
    cached_at = cache.get("cached_at", 0)
    age_hours = (time.time() - cached_at) / 3600
    if age_hours > CACHE_MAX_AGE_HOURS:
        return None  # Stale
    return cache


def refresh_cache(force=False):
    """Fetch latest models from OpenRouter and rebuild the cache."""
    import requests

    # Rate limit check
    if not force:
        allowed, remaining = check_rate_limit()
        if not allowed:
            print(f"  Rate limited: wait {remaining}s before refreshing again.")
            # Try to return stale cache
            stale = safe_read_json(CACHE_PATH)
            return stale

    print("  Refreshing model cache from OpenRouter...")
    api_key = os.environ.get("OPENROUTER_API_KEY", "")
    headers = {}
    if api_key:
        headers["Authorization"] = f"Bearer {api_key}"
    else:
        print("  WARNING: OPENROUTER_API_KEY not set. Fetching public model list.")

    try:
        resp = requests.get(
            "https://openrouter.ai/api/v1/models",
            headers=headers,
            timeout=30,
        )
        resp.raise_for_status()
        data = resp.json()
    except requests.exceptions.Timeout:
        print("  ERROR: OpenRouter API timed out after 30s.")
        print("  Keeping existing cache. Try again later.")
        return safe_read_json(CACHE_PATH)
    except requests.exceptions.ConnectionError:
        print("  ERROR: Could not connect to OpenRouter API.")
        print("  Keeping existing cache. Check internet connection.")
        return safe_read_json(CACHE_PATH)
    except Exception as e:
        # Sanitize error message to avoid leaking API key
        err_msg = str(e)
        if api_key and api_key in err_msg:
            err_msg = err_msg.replace(api_key, "[REDACTED]")
        print(f"  ERROR: Failed to fetch models: {err_msg}")
        print("  Keeping existing cache.")
        return safe_read_json(CACHE_PATH)

    models = data.get("data", [])
    print(f"  Fetched {len(models)} models from OpenRouter.")

    # Build a lookup for quick access — DBG-011: use .get() to skip malformed entries
    model_lookup = {}
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

    # Identify frontier models — DBG-012: safe access with try/except
    frontier_providers = ["openai", "anthropic", "google", "x-ai", "deepseek", "meta-llama", "mistralai"]
    frontier = []
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

    cache = {
        "cached_at": time.time(),
        "cached_at_human": datetime.now(timezone.utc).isoformat(),
        "total_models": len(models),
        "frontier_models": frontier,
        "model_lookup": model_lookup,
    }

    try:
        atomic_write_json(CACHE_PATH, cache)
        update_rate_limit()
        print(f"  Cache saved: {len(frontier)} frontier models indexed.")
    except IOError as e:
        print(f"  WARNING: Could not save cache: {e}")
        print(f"  Check file permissions for {CACHE_PATH}")

    return cache


def get_cache():
    """Get cache if available. Does NOT trigger network refresh.
    Use for read-only commands (status, list, recommend)."""
    # DBG-005: read-only — never triggers network call
    return load_cache()


def get_or_refresh_cache():
    """Get cache, refreshing if stale/missing. Use for commands that need fresh data (auto, set)."""
    cache = load_cache()
    if cache is None:
        cache = refresh_cache()
    # DBG-001: ensure we never return None
    if cache is None:
        cache = {}
    return cache


# ===========================================================================
# Category detection
# ===========================================================================

def detect_category(task_text):
    """Dynamically classify a task description into a category."""
    if not task_text:
        return "general"

    task_lower = task_text.lower()
    scores = {}

    for category, keywords in CATEGORY_KEYWORDS.items():
        score = 0
        for kw in keywords:
            # DBG-003: Use word boundary regex instead of substring matching
            # DBG-016: Handle multi-word keywords properly
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
        scores[category] = score

    best = max(scores, key=scores.get)
    if scores[best] == 0:
        return "general"

    # DBG-013: If 'general' tied with a specialized category, prefer specialized
    if best == "general":
        specialized = {k: v for k, v in scores.items() if k != "general" and v == scores[best]}
        if specialized:
            best = max(specialized, key=specialized.get)

    return best


# ===========================================================================
# Model selection
# ===========================================================================

def get_rankings(category, use_refresh=False):
    """Get model rankings for a category, using cache if available.
    
    Args:
        category: Task category to get rankings for.
        use_refresh: If True, use get_or_refresh_cache() for fresh data.
                     If False, use get_cache() for read-only access.
    """
    # DBG-005: callers choose whether to trigger network refresh
    cache = get_or_refresh_cache() if use_refresh else get_cache()

    # Start with default rankings
    rankings = DEFAULT_RANKINGS.get(category, DEFAULT_RANKINGS["general"])

    if cache and cache.get("model_lookup"):
        # Validate that ranked models actually exist in OpenRouter
        lookup = cache["model_lookup"]
        validated = []
        for r in rankings:
            if r["id"] in lookup:
                validated.append(r)
        if validated:
            rankings = validated

    return rankings


def select_best_model(task_text=None, category=None):
    """Select the best model for a task."""
    if category is None:
        category = detect_category(task_text or "")

    rankings = get_rankings(category, use_refresh=True)
    if not rankings:
        return None, category, "No models available"

    best = rankings[0]
    return best["id"], category, best["reason"]


# ===========================================================================
# Commands
# ===========================================================================

def cmd_set(args):
    """Set a specific model as the default."""
    model_id = args.model_id
    cfg = load_config()

    # Validate model exists in cache — uses refresh for validation
    cache = get_or_refresh_cache()
    if cache and cache.get("model_lookup"):
        if model_id not in cache["model_lookup"]:
            # Try partial match
            matches = [m for m in cache["model_lookup"] if model_id.lower() in m.lower()]
            if len(matches) == 1:
                model_id = matches[0]
                print(f"  Matched to: {model_id}")
            elif len(matches) > 1:
                # DBG-010: Show descriptions to help disambiguate
                print(f"  Ambiguous model name. Did you mean one of:")
                lookup = cache["model_lookup"]
                for m in matches[:10]:
                    desc = lookup.get(m, {}).get("description", "")[:60]
                    ctx = lookup.get(m, {}).get("context_length", 0)
                    print(f"    {m}  (ctx={ctx:,}) {desc}")
                print(f"\n  Provide a more specific model ID to disambiguate.")
                sys.exit(1)
            else:
                print(f"  WARNING: Model '{model_id}' not found in OpenRouter catalog.")
                print(f"  Setting anyway — it may be a valid model not in the public list.")

    cfg["mode"] = "manual"
    cfg["manual_model"] = model_id
    cfg["last_selected_model"] = model_id
    cfg["last_selected_reason"] = "Manually set by user"
    cfg["enabled"] = True
    save_config(cfg)

    print(f"\n  Model set: {model_id}")
    print(f"  Mode: manual")
    print(f"  Status: enabled")
    print(f"\n  Manus will use this model as its backbone LLM for all tasks.")


def cmd_auto(args):
    """Auto-select the best model for a task."""
    task_text = " ".join(args.task_description) if args.task_description else None
    cfg = load_config()

    if task_text:
        model_id, category, reason = select_best_model(task_text=task_text)
        print(f"\n  Task analysis: \"{task_text}\"")
    else:
        model_id, category, reason = select_best_model()
        print(f"\n  Auto-selecting best general-purpose model...")

    if model_id is None:
        print("  ERROR: Could not determine best model.")
        sys.exit(1)

    cfg["mode"] = "auto"
    cfg["auto_category"] = category
    cfg["last_selected_model"] = model_id
    cfg["last_selected_reason"] = reason
    cfg["enabled"] = True
    save_config(cfg)

    print(f"  Category detected: {category}")
    print(f"  Selected model: {model_id}")
    print(f"  Reason: {reason}")
    print(f"  Status: enabled")

    # Show alternatives — read-only, no refresh needed
    rankings = get_rankings(category, use_refresh=False)
    if len(rankings) > 1:
        print(f"\n  Alternatives for {category}:")
        for i, r in enumerate(rankings[1:5], 2):
            print(f"    #{i}: {r['id']} — {r['reason']}")


def cmd_on(args):
    """Enable model-selector."""
    cfg = load_config()
    cfg["enabled"] = True
    save_config(cfg)
    model = cfg.get("last_selected_model", "not set")
    mode = cfg.get("mode", "auto")
    print(f"\n  Model selector: ENABLED")
    print(f"  Mode: {mode}")
    print(f"  Current model: {model}")


def cmd_off(args):
    """Disable model-selector."""
    cfg = load_config()
    cfg["enabled"] = False
    save_config(cfg)
    print(f"\n  Model selector: DISABLED")
    print(f"  Manus will use its default backbone model.")


def cmd_status(args):
    """Show current configuration."""
    cfg = load_config()
    cache = load_cache()  # DBG-005: read-only, never triggers network

    print("\n" + "=" * 60)
    print("  MODEL SELECTOR — Status")
    print("=" * 60)
    print(f"  Enabled:    {'YES' if cfg.get('enabled') else 'NO'}")
    print(f"  Mode:       {cfg.get('mode', 'auto')}")

    if cfg.get("mode") == "manual":
        print(f"  Model:      {cfg.get('manual_model', 'not set')}")
    else:
        print(f"  Model:      {cfg.get('last_selected_model', 'not set')}")
        print(f"  Category:   {cfg.get('auto_category', 'not set')}")

    print(f"  Reason:     {cfg.get('last_selected_reason', 'N/A')}")
    print(f"  Updated:    {cfg.get('updated_at', 'never')}")

    if cache:
        age_hours = (time.time() - cache.get("cached_at", 0)) / 3600
        print(f"\n  Cache:      {cache.get('total_models', 0)} models indexed")
        print(f"  Cache age:  {age_hours:.1f} hours")
        print(f"  Cached at:  {cache.get('cached_at_human', 'unknown')}")
    else:
        print(f"\n  Cache:      not loaded (run 'refresh' to populate)")

    print(f"\n  Config:     {CONFIG_PATH}")
    print(f"  Cache:      {CACHE_PATH}")
    print("-" * 60)


def cmd_refresh(args):
    """Refresh the model cache."""
    cache = refresh_cache(force=True)
    if cache and cache.get("total_models"):
        print(f"\n  Cache refreshed successfully.")
        print(f"  Total models: {cache['total_models']}")
        print(f"  Frontier models: {len(cache.get('frontier_models', []))}")
    else:
        print(f"\n  Cache refresh failed or returned empty. Using default rankings.")


def cmd_list(args):
    """List top models by category."""
    category = args.category or "general"
    if category not in CATEGORY_KEYWORDS and category != "all":
        print(f"  Unknown category: {category}")
        print(f"  Available: {', '.join(CATEGORY_KEYWORDS.keys())}, all")
        sys.exit(1)

    if category == "all":
        for cat in CATEGORY_KEYWORDS:
            print(f"\n  === {cat.upper()} ===")
            rankings = get_rankings(cat, use_refresh=False)  # DBG-005: read-only
            for i, r in enumerate(rankings[:3], 1):
                print(f"    #{i}: {r['id']:50s} — {r['reason']}")
    else:
        print(f"\n  === {category.upper()} — Top Models ===")
        rankings = get_rankings(category, use_refresh=False)  # DBG-005: read-only
        for i, r in enumerate(rankings, 1):
            print(f"    #{i}: {r['id']:50s} — {r['reason']}")

        # Show pricing if cache available — DBG-005: read-only
        cache = get_cache()
        if cache and cache.get("model_lookup"):
            print(f"\n  Pricing (per 1M tokens):")
            for r in rankings:
                m = cache["model_lookup"].get(r["id"])
                if m:
                    # DBG-004: safe pricing conversion with fallback
                    try:
                        pi = float(m.get("pricing_prompt") or "0") * 1_000_000
                        po = float(m.get("pricing_completion") or "0") * 1_000_000
                    except (ValueError, TypeError):
                        pi, po = 0.0, 0.0
                    ctx = m.get("context_length", 0) or 0
                    print(f"    {r['id']:50s} ${pi:.2f}/${po:.2f}  ctx={ctx:,}")


def cmd_recommend(args):
    """Recommend the best model for a task without setting it."""
    task_text = " ".join(args.task_description)
    model_id, category, reason = select_best_model(task_text=task_text)

    print(f"\n  Task: \"{task_text}\"")
    print(f"  Category: {category}")
    print(f"  Recommended: {model_id}")
    print(f"  Reason: {reason}")

    rankings = get_rankings(category, use_refresh=False)  # Read-only for display
    if len(rankings) > 1:
        print(f"\n  Full ranking for {category}:")
        for i, r in enumerate(rankings, 1):
            marker = " <-- recommended" if r["id"] == model_id else ""
            print(f"    #{i}: {r['id']:50s} — {r['reason']}{marker}")

    print(f"\n  To set this model: python3 model_selector.py set {model_id}")
    print(f"  To auto-select:    python3 model_selector.py auto {task_text}")


# ===========================================================================
# Main
# ===========================================================================

def main():
    parser = argparse.ArgumentParser(
        description="Model Selector — Set or auto-select the best AI model for Manus",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=__doc__,
    )
    subparsers = parser.add_subparsers(dest="command", help="Available commands")

    # set
    p_set = subparsers.add_parser("set", help="Set a specific model as default")
    p_set.add_argument("model_id", help="Model ID (e.g., anthropic/claude-opus-4.6)")

    # auto
    p_auto = subparsers.add_parser("auto", help="Auto-select the best model")
    p_auto.add_argument("task_description", nargs="*", help="Optional task description")

    # on/off
    subparsers.add_parser("on", help="Enable model-selector")
    subparsers.add_parser("off", help="Disable model-selector")

    # status
    subparsers.add_parser("status", help="Show current configuration")

    # refresh
    subparsers.add_parser("refresh", help="Refresh model cache from OpenRouter")

    # list
    p_list = subparsers.add_parser("list", help="List top models by category")
    p_list.add_argument("category", nargs="?", help="Category (coding/reasoning/research/creative/math/general/all)")

    # recommend
    p_rec = subparsers.add_parser("recommend", help="Recommend best model for a task")
    p_rec.add_argument("task_description", nargs="+", help="Task description")

    args = parser.parse_args()

    if not args.command:
        parser.print_help()
        sys.exit(0)

    commands = {
        "set": cmd_set,
        "auto": cmd_auto,
        "on": cmd_on,
        "off": cmd_off,
        "status": cmd_status,
        "refresh": cmd_refresh,
        "list": cmd_list,
        "recommend": cmd_recommend,
    }

    cmd_func = commands.get(args.command)
    if cmd_func is None:
        print(f"  Unknown command: {args.command}")
        parser.print_help()
        sys.exit(1)

    try:
        cmd_func(args)
    except KeyboardInterrupt:
        print("\n  Interrupted.")
        sys.exit(130)
    except Exception as e:
        # Sanitize any potential API key leaks
        err_msg = str(e)
        api_key = os.environ.get("OPENROUTER_API_KEY", "")
        if api_key and api_key in err_msg:
            err_msg = err_msg.replace(api_key, "[REDACTED]")
        print(f"  ERROR: {err_msg}")
        sys.exit(1)


if __name__ == "__main__":
    main()
