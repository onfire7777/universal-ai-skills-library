#!/usr/bin/env python3
"""Multi-Model Oracle — Get the ultimate merged answer from the best AI models.

4-Stage Pipeline:
  1. Prompt Engineering: Analyze intent, optimize the query, create model-specific variants
  2. Parallel Query: Send to Anthropic's best, OpenAI's best, and Manus simultaneously
  3. Intelligent Merge: Extract strongest elements, resolve contradictions, unify
  4. Output: Deliver the ultimate merged answer

Usage:
    python3 oracle.py "What is the best approach to distributed consensus?"
    python3 oracle.py "Write a Python async web scraper" --show-individual
    python3 oracle.py --file prompt.txt --output result.md
    echo "Explain quantum entanglement" | python3 oracle.py --stdin
"""

import argparse
import json
import os
import re
import sys
import time
from concurrent.futures import ThreadPoolExecutor, as_completed
from datetime import datetime, timezone
from pathlib import Path
from typing import Optional

# ─── Dependency Check ─────────────────────────────────────────────────────

try:
    import requests
except ImportError:
    print(
        "ERROR: 'requests' package not installed.\n"
        "Install with: pip3 install requests",
        file=sys.stderr,
    )
    sys.exit(1)

try:
    from openai import OpenAI
except ImportError:
    OpenAI = None  # Will be checked when needed


def _check_openai():
    """Check if the openai package is available."""
    if OpenAI is None:
        print(
            "WARNING: 'openai' package not installed. Manus model and prompt "
            "engineering will be unavailable.\n"
            "Install with: pip3 install openai",
            file=sys.stderr,
        )
        return False
    return True


def _check_env_keys():
    """Check required environment variables and warn if missing."""
    warnings = []
    if not os.environ.get("OPENROUTER_API_KEY"):
        warnings.append(
            "OPENROUTER_API_KEY not set — Anthropic and OpenAI models will be unavailable"
        )
    if not os.environ.get("OPENAI_API_KEY"):
        warnings.append(
            "OPENAI_API_KEY not set — Manus model and prompt engineering will be unavailable"
        )
    for w in warnings:
        _log(f"WARNING: {w}")
    return len(warnings) == 0


def _log(msg: str):
    """Print progress/status to stderr (never stdout)."""
    print(msg, file=sys.stderr)


# ─── Configuration ──────────────────────────────────────────────────────────

# Model IDs — the latest best models as of April 2026
ANTHROPIC_MODEL = "anthropic/claude-opus-4.6"
OPENAI_MODEL = "openai/gpt-5.4"
MANUS_MODEL = "gpt-4.1-mini"

# API endpoints
OPENROUTER_URL = "https://openrouter.ai/api/v1/chat/completions"

# Timeouts and retries
REQUEST_TIMEOUT = (15, 180)  # (connect, read) in seconds
MAX_RETRIES = 2
RETRY_DELAYS = [5, 15, 30]  # Extra element for safety if MAX_RETRIES is increased

# Non-retryable HTTP status codes
FATAL_STATUSES = {401, 402, 403, 422}


# ─── Intent Detection ──────────────────────────────────────────────────────

INTENT_PATTERNS = {
    "creative": [
        # Check creative FIRST to avoid "create a story" matching code
        r"\b(write|compose|draft|create|generate)\b.*\b(story|poem|essay|article|blog|letter|speech|song|novel|screenplay)\b",
        r"\b(creative|imaginative|artistic|literary|narrative|fiction|poetic)\b",
    ],
    "code": [
        r"\b(implement|code|script|function|class|program|debug|refactor|compile)\b",
        r"\b(python|javascript|typescript|rust|go|java|c\+\+|sql|html|css|bash|react|django|flask)\b",
        r"\b(api|endpoint|database|server|frontend|backend|algorithm|data structure|repository|git)\b",
        r"\b(write|create|build)\b.*\b(code|script|function|program|app|application|tool|cli)\b",
    ],
    "analysis": [
        r"\b(analyze|compare|evaluate|assess|review|critique|examine|investigate)\b",
        r"\b(pros?\s+and\s+cons?|trade.?offs?|advantages|disadvantages|strengths|weaknesses)\b",
        r"\b(difference|similarities|versus|vs\.?)\b",
    ],
    "research": [
        r"\b(explain|describe|what\s+is|how\s+does|why\s+does|define|overview)\b",
        r"\b(history|background|context|theory|concept|principle|mechanism)\b",
    ],
    "reasoning": [
        r"\b(solve|prove|derive|calculate|logic|reason|deduce|infer)\b",
        r"\b(mathematical|philosophical|ethical|moral|paradox|dilemma)\b",
        r"\b(step.by.step|think through|work through|break down)\b",
    ],
    "practical": [
        r"\b(how\s+to|guide|tutorial|steps|instructions|setup|configure|install)\b",
        r"\b(best\s+practice|recommendation|suggest|advice|tip)\b",
    ],
}


def detect_intent(query: str) -> str:
    """Detect the primary intent of the user's query.

    Creative intent is checked first to prevent 'create a story' from
    matching the code intent's 'create' keyword.
    """
    query_lower = query.lower()
    scores = {}

    # Check creative first — if it matches, return immediately to avoid
    # false positives from broad code patterns
    for intent in ["creative", "code", "analysis", "research", "reasoning", "practical"]:
        patterns = INTENT_PATTERNS.get(intent, [])
        score = 0
        for pattern in patterns:
            matches = re.findall(pattern, query_lower)
            score += len(matches)
        scores[intent] = score

    if max(scores.values()) == 0:
        return "general"
    return max(scores, key=scores.get)


# ─── Stage 1: Prompt Engineering ────────────────────────────────────────────

PROMPT_ENGINEER_SYSTEM = """You are an elite prompt engineer. Your job is to take a user's raw query and transform it into the most effective possible prompt for an AI model.

You will receive:
1. The user's original query
2. The detected intent category
3. The target model's strengths

Your output MUST be valid JSON with exactly these fields:
{
  "enhanced_prompt": "The improved prompt text",
  "reasoning": "Brief explanation of what you improved and why (1-2 sentences)"
}

Prompt engineering principles to apply:
- ADD SPECIFICITY: If the query is vague, add concrete parameters (format, length, depth, audience)
- ADD STRUCTURE: Request organized output (sections, numbered lists, tables where appropriate)
- ADD CONTEXT: Infer and state implicit requirements the user likely wants
- ROLE PRIMING: Add an expert persona appropriate to the domain
- CHAIN OF THOUGHT: For reasoning tasks, request step-by-step thinking
- FEW-SHOT: For format-sensitive tasks, include a brief example of desired output
- NEGATIVE CONSTRAINTS: Specify what to avoid (generic advice, filler, hedging)
- PRESERVE INTENT: Never change what the user is actually asking for

Do NOT:
- Add unnecessary complexity to simple questions
- Change the core question or topic
- Add requirements the user clearly didn't want
- Make the prompt longer than necessary — concise improvement is best"""


def build_model_strength_description(model_name: str) -> str:
    """Describe each model's strengths for the prompt engineer."""
    strengths = {
        "anthropic": (
            "Anthropic Claude Opus 4.6 — excels at nuanced reasoning, creative writing, "
            "philosophical depth, careful analysis, safety-aware responses, and long-form "
            "structured output. Particularly strong at considering multiple perspectives "
            "and acknowledging uncertainty."
        ),
        "openai": (
            "OpenAI GPT-5.4 — excels at breadth of knowledge, technical precision, "
            "code generation, structured data output, mathematical reasoning, and "
            "following complex instructions. Particularly strong at systematic "
            "problem-solving and comprehensive coverage."
        ),
        "manus": (
            "GPT-4.1-mini — fast and efficient, good at structured tasks, concise "
            "answers, and practical advice. Best used for clear, direct responses "
            "that prioritize actionability over depth."
        ),
    }
    for key, desc in strengths.items():
        if key in model_name.lower():
            return desc
    return strengths["manus"]


def _engineer_single_variant(client, model_key: str, model_id: str,
                              query: str, intent: str) -> dict:
    """Engineer a single prompt variant for one model. Used in parallel."""
    strengths = build_model_strength_description(model_id)
    engineer_prompt = (
        f"User's original query:\n\"{query}\"\n\n"
        f"Detected intent: {intent}\n\n"
        f"Target model strengths:\n{strengths}\n\n"
        f"Transform this query into the best possible prompt for this specific model. "
        f"Leverage the model's strengths. Output valid JSON only."
    )

    try:
        response = client.chat.completions.create(
            model=MANUS_MODEL,
            messages=[
                {"role": "system", "content": PROMPT_ENGINEER_SYSTEM},
                {"role": "user", "content": engineer_prompt},
            ],
            max_tokens=1500,
            temperature=0.3,
            response_format={"type": "json_object"},
            timeout=60,
        )
        raw = response.choices[0].message.content.strip()
        parsed = json.loads(raw)
        return {
            "key": model_key,
            "prompt": parsed.get("enhanced_prompt", query),
            "reasoning": parsed.get("reasoning", ""),
        }
    except Exception as e:
        _log(f"  [Prompt Engineer] Warning: {model_key} variant failed: {e}")
        return {"key": model_key, "prompt": query, "reasoning": ""}


def prompt_engineer(query: str, intent: str) -> dict:
    """Stage 1: Use Manus model to prompt-engineer the user's query IN PARALLEL.

    Returns dict with keys: anthropic_prompt, openai_prompt, manus_prompt, reasoning
    """
    if not _check_openai() or not os.environ.get("OPENAI_API_KEY"):
        return {
            "anthropic_prompt": query,
            "openai_prompt": query,
            "manus_prompt": query,
            "reasoning": "Prompt engineering skipped: openai package or API key not available.",
        }

    client = OpenAI()
    results = {}

    models_info = [
        ("anthropic", ANTHROPIC_MODEL),
        ("openai", OPENAI_MODEL),
        ("manus", MANUS_MODEL),
    ]

    # Parallelize the 3 prompt engineering calls
    with ThreadPoolExecutor(max_workers=3) as executor:
        futures = {
            executor.submit(
                _engineer_single_variant, client, model_key, model_id, query, intent
            ): model_key
            for model_key, model_id in models_info
        }

        reasoning_parts = []
        for future in as_completed(futures):
            try:
                result = future.result()
                results[f"{result['key']}_prompt"] = result["prompt"]
                if result["reasoning"]:
                    reasoning_parts.append(result["reasoning"])
            except Exception as e:
                model_key = futures[future]
                _log(f"  [Prompt Engineer] Warning: {model_key} variant failed: {e}")
                results[f"{model_key}_prompt"] = query

    # Ensure all keys exist
    for key in ["anthropic_prompt", "openai_prompt", "manus_prompt"]:
        if key not in results:
            results[key] = query

    results["reasoning"] = reasoning_parts[0] if reasoning_parts else "Prompt engineering completed."

    return results


# ─── Stage 2: Parallel Model Query ─────────────────────────────────────────

def query_openrouter(model_id: str, prompt: str, intent: str) -> dict:
    """Query a model via OpenRouter API with retry logic."""
    api_key = os.environ.get("OPENROUTER_API_KEY")
    if not api_key:
        return {"model": model_id, "success": False,
                "error": "OPENROUTER_API_KEY not set"}

    # Build system message based on intent
    system_messages = {
        "code": "You are an expert software engineer. Provide production-quality, well-documented code with clear explanations.",
        "creative": "You are a masterful writer with deep literary sensibility. Create vivid, original, emotionally resonant content.",
        "analysis": "You are a rigorous analytical thinker. Provide balanced, evidence-based analysis with clear structure.",
        "research": "You are a knowledgeable researcher. Explain concepts clearly with accurate details and helpful examples.",
        "reasoning": "You are a precise logical thinker. Show your reasoning step by step, verify each step, and state your confidence.",
        "practical": "You are a practical expert advisor. Provide clear, actionable guidance with concrete steps.",
        "general": "You are a knowledgeable, thoughtful assistant. Provide comprehensive, well-structured responses.",
    }
    system_msg = system_messages.get(intent, system_messages["general"])

    headers = {
        "Authorization": f"Bearer {api_key}",
        "Content-Type": "application/json",
        "HTTP-Referer": "https://manus.im",
        "X-Title": "Multi-Model Oracle",
    }
    payload = {
        "model": model_id,
        "messages": [
            {"role": "system", "content": system_msg},
            {"role": "user", "content": prompt},
        ],
        "max_tokens": 4096,
        "temperature": 0.4,
    }

    for attempt in range(MAX_RETRIES + 1):
        try:
            t0 = time.time()
            resp = requests.post(
                OPENROUTER_URL, headers=headers, json=payload,
                timeout=REQUEST_TIMEOUT,
            )
            elapsed = time.time() - t0

            if resp.status_code == 200:
                data = resp.json()
                # Safe access: validate response structure
                choices = data.get("choices")
                if not choices or not isinstance(choices, list):
                    return {"model": model_id, "success": False,
                            "error": f"Unexpected response: no 'choices' in response"}
                message = choices[0].get("message")
                if not message or "content" not in message:
                    return {"model": model_id, "success": False,
                            "error": "Unexpected response: no 'content' in message"}
                content = message["content"].strip()
                usage = data.get("usage", {})
                return {
                    "model": model_id,
                    "success": True,
                    "content": content,
                    "elapsed": round(elapsed, 1),
                    "tokens_in": usage.get("prompt_tokens", 0),
                    "tokens_out": usage.get("completion_tokens", 0),
                }

            if resp.status_code in FATAL_STATUSES:
                error_detail = resp.text[:200]
                return {"model": model_id, "success": False,
                        "error": f"HTTP {resp.status_code}: {error_detail}"}

            # Retryable error
            if attempt < MAX_RETRIES:
                delay = RETRY_DELAYS[min(attempt, len(RETRY_DELAYS) - 1)]
                time.sleep(delay)
                continue

            return {"model": model_id, "success": False,
                    "error": f"HTTP {resp.status_code} after {MAX_RETRIES + 1} attempts"}

        except requests.exceptions.Timeout:
            if attempt < MAX_RETRIES:
                delay = RETRY_DELAYS[min(attempt, len(RETRY_DELAYS) - 1)]
                time.sleep(delay)
                continue
            return {"model": model_id, "success": False,
                    "error": f"Timeout after {REQUEST_TIMEOUT[1]}s"}
        except Exception as e:
            if attempt < MAX_RETRIES:
                delay = RETRY_DELAYS[min(attempt, len(RETRY_DELAYS) - 1)]
                time.sleep(delay)
                continue
            return {"model": model_id, "success": False, "error": str(e)}


def query_manus(prompt: str, intent: str) -> dict:
    """Query the Manus built-in model via OpenAI-compatible API."""
    if not _check_openai():
        return {"model": MANUS_MODEL, "success": False,
                "error": "openai package not installed"}

    if not os.environ.get("OPENAI_API_KEY"):
        return {"model": MANUS_MODEL, "success": False,
                "error": "OPENAI_API_KEY not set"}

    system_messages = {
        "code": "You are an expert software engineer. Provide production-quality, well-documented code with clear explanations.",
        "creative": "You are a masterful writer. Create vivid, original content.",
        "analysis": "You are a rigorous analytical thinker. Provide balanced, evidence-based analysis.",
        "research": "You are a knowledgeable researcher. Explain concepts clearly with accurate details.",
        "reasoning": "You are a precise logical thinker. Show your reasoning step by step.",
        "practical": "You are a practical expert advisor. Provide clear, actionable guidance.",
        "general": "You are a knowledgeable, thoughtful assistant. Provide comprehensive responses.",
    }
    system_msg = system_messages.get(intent, system_messages["general"])

    client = OpenAI()
    try:
        t0 = time.time()
        response = client.chat.completions.create(
            model=MANUS_MODEL,
            messages=[
                {"role": "system", "content": system_msg},
                {"role": "user", "content": prompt},
            ],
            max_tokens=4096,
            temperature=0.4,
        )
        elapsed = time.time() - t0
        content = response.choices[0].message.content.strip()
        usage = response.usage
        return {
            "model": MANUS_MODEL,
            "success": True,
            "content": content,
            "elapsed": round(elapsed, 1),
            "tokens_in": usage.prompt_tokens if usage else 0,
            "tokens_out": usage.completion_tokens if usage else 0,
        }
    except Exception as e:
        return {"model": MANUS_MODEL, "success": False, "error": str(e)}


def query_all_models(prompts: dict, intent: str) -> list:
    """Stage 2: Query all 3 models in parallel."""
    results = []

    def _query_anthropic():
        return query_openrouter(
            ANTHROPIC_MODEL, prompts["anthropic_prompt"], intent
        )

    def _query_openai():
        return query_openrouter(
            OPENAI_MODEL, prompts["openai_prompt"], intent
        )

    def _query_manus():
        return query_manus(prompts["manus_prompt"], intent)

    tasks = {
        "Anthropic Claude Opus 4.6": _query_anthropic,
        "OpenAI GPT-5.4": _query_openai,
        "Manus gpt-4.1-mini": _query_manus,
    }

    _log("  Querying models in parallel...")
    with ThreadPoolExecutor(max_workers=3) as executor:
        futures = {}
        for name, fn in tasks.items():
            futures[executor.submit(fn)] = name

        for future in as_completed(futures):
            name = futures[future]
            try:
                result = future.result()
                status = "OK" if result["success"] else f"FAIL: {result.get('error', '?')}"
                elapsed = f" ({result.get('elapsed', '?')}s)" if result["success"] else ""
                _log(f"    [{name}] {status}{elapsed}")
                results.append(result)
            except Exception as e:
                _log(f"    [{name}] EXCEPTION: {e}")
                results.append({"model": name, "success": False, "error": str(e)})

    return results


# ─── Stage 3: Intelligent Merge ────────────────────────────────────────────

MERGE_SYSTEM = """You are an expert synthesis engine. You receive responses from multiple AI sources to the same question. Your job is to create the ULTIMATE merged answer that is better than any individual response.

Merge strategy:
1. IDENTIFY the strongest elements from each response (unique insights, best explanations, most accurate details)
2. RESOLVE contradictions: if sources disagree, use majority opinion OR present both perspectives with clear labeling
3. ELIMINATE redundancy: don't repeat the same point from multiple sources
4. PRESERVE unique value: if only one source mentioned something important, keep it
5. STRUCTURE coherently: the final answer should flow naturally, not feel like a patchwork
6. MAINTAIN the appropriate tone and depth for the query type

Output rules:
- Write the merged answer directly — do NOT mention the individual sources or the merging process
- Do NOT say "Source A said X while Source B said Y" — just present the best unified answer
- The output should read as if it came from a single, exceptionally knowledgeable source
- If the query was about code, include the best code with the best explanations
- If the query was creative, produce the most compelling creative output
- Match the format the user would expect (paragraphs for explanations, code blocks for code, etc.)"""

MERGE_MAX_RETRIES = 2
MERGE_RETRY_DELAYS = [3, 10]


def merge_responses(query: str, results: list, intent: str) -> str:
    """Stage 3: Intelligently merge all successful model responses with retry logic."""
    successful = [r for r in results if r["success"]]

    if not successful:
        return "All models failed to respond. Please check your API keys and try again."

    if len(successful) == 1:
        return successful[0]["content"]

    # Build the merge prompt — use generic labels, not model names
    response_sections = []
    for i, r in enumerate(successful, 1):
        response_sections.append(
            f"=== Response {i} ===\n{r['content']}"
        )

    merge_prompt = (
        f"Original user query:\n\"{query}\"\n\n"
        f"Query type: {intent}\n\n"
        f"{'=' * 60}\n"
        + "\n\n".join(response_sections)
        + f"\n\n{'=' * 60}\n\n"
        f"Create the ultimate merged answer. Write it directly — do not reference the individual responses."
    )

    if not _check_openai() or not os.environ.get("OPENAI_API_KEY"):
        # Fallback: concatenate
        _log("  [Merge] Warning: openai not available, using concatenation fallback")
        parts = []
        for i, r in enumerate(successful, 1):
            parts.append(f"## Response {i}\n\n{r['content']}")
        return "\n\n---\n\n".join(parts)

    client = OpenAI()

    for attempt in range(MERGE_MAX_RETRIES + 1):
        try:
            response = client.chat.completions.create(
                model=MANUS_MODEL,
                messages=[
                    {"role": "system", "content": MERGE_SYSTEM},
                    {"role": "user", "content": merge_prompt},
                ],
                max_tokens=8000,
                temperature=0.3,
            )
            return response.choices[0].message.content.strip()
        except Exception as e:
            if attempt < MERGE_MAX_RETRIES:
                delay = MERGE_RETRY_DELAYS[min(attempt, len(MERGE_RETRY_DELAYS) - 1)]
                _log(f"  [Merge] Attempt {attempt + 1} failed ({e}), retrying in {delay}s...")
                time.sleep(delay)
                continue
            # Final fallback: concatenate
            _log(f"  [Merge] Warning: AI merge failed after {MERGE_MAX_RETRIES + 1} attempts ({e}), using concatenation fallback")
            parts = []
            for i, r in enumerate(successful, 1):
                parts.append(f"## Response {i}\n\n{r['content']}")
            return "\n\n---\n\n".join(parts)


# ─── Stage 4: Output ───────────────────────────────────────────────────────

def format_output(
    query: str,
    intent: str,
    prompts: dict,
    results: list,
    merged: str,
    show_individual: bool = False,
    show_prompts: bool = False,
) -> str:
    """Format the final output."""
    lines = []

    # Header
    lines.append("=" * 70)
    lines.append("MULTI-MODEL ORACLE — ULTIMATE ANSWER")
    lines.append("=" * 70)
    lines.append(f"Query: {query[:100]}{'...' if len(query) > 100 else ''}")
    lines.append(f"Intent: {intent}")
    lines.append(f"Timestamp: {datetime.now(timezone.utc).strftime('%Y-%m-%d %H:%M:%S UTC')}")

    # Model status
    lines.append("")
    lines.append("Models:")
    for r in results:
        model = r["model"].split("/")[-1] if "/" in r["model"] else r["model"]
        if r["success"]:
            tokens = f"  ({r.get('tokens_in', 0)}+{r.get('tokens_out', 0)} tokens)"
            lines.append(f"  + {model:25s} {r.get('elapsed', '?')}s{tokens}")
        else:
            lines.append(f"  x {model:25s} FAILED: {r.get('error', 'unknown')}")

    successful_count = sum(1 for r in results if r["success"])
    lines.append(f"\nMerged from {successful_count}/3 models")

    # Prompt engineering info
    if show_prompts:
        lines.append("")
        lines.append("-" * 70)
        lines.append("PROMPT ENGINEERING")
        lines.append("-" * 70)
        lines.append(f"Reasoning: {prompts.get('reasoning', 'N/A')}")
        for key in ["anthropic_prompt", "openai_prompt", "manus_prompt"]:
            label = key.replace("_prompt", "").upper()
            lines.append(f"\n[{label}]:\n{prompts.get(key, 'N/A')}")

    # The merged answer
    lines.append("")
    lines.append("=" * 70)
    lines.append("ANSWER")
    lines.append("=" * 70)
    lines.append("")
    lines.append(merged)

    # Individual responses
    if show_individual:
        lines.append("")
        lines.append("=" * 70)
        lines.append("INDIVIDUAL MODEL RESPONSES")
        lines.append("=" * 70)
        for r in results:
            model = r["model"].split("/")[-1] if "/" in r["model"] else r["model"]
            lines.append(f"\n--- {model} ---")
            if r["success"]:
                lines.append(r["content"])
            else:
                lines.append(f"[FAILED: {r.get('error', 'unknown')}]")

    lines.append("")
    lines.append("=" * 70)

    return "\n".join(lines)


# ─── Main ───────────────────────────────────────────────────────────────────

def run_oracle(
    query: str,
    show_individual: bool = False,
    show_prompts: bool = False,
    output_file: Optional[str] = None,
    skip_engineering: bool = False,
) -> str:
    """Run the full 4-stage oracle pipeline."""
    _log("=" * 60)
    _log("MULTI-MODEL ORACLE")
    _log("=" * 60)

    # Check environment
    _check_env_keys()

    # Stage 1: Prompt Engineering
    _log("\n[Stage 1] Prompt Engineering...")
    intent = detect_intent(query)
    _log(f"  Detected intent: {intent}")

    if skip_engineering:
        prompts = {
            "anthropic_prompt": query,
            "openai_prompt": query,
            "manus_prompt": query,
            "reasoning": "Prompt engineering skipped (--raw flag).",
        }
        _log("  Skipped (--raw flag)")
    else:
        prompts = prompt_engineer(query, intent)
        _log(f"  Reasoning: {prompts.get('reasoning', 'N/A')}")

    # Stage 2: Parallel Query
    _log(f"\n[Stage 2] Querying 3 models...")
    results = query_all_models(prompts, intent)

    successful = sum(1 for r in results if r["success"])
    _log(f"\n  {successful}/3 models responded successfully")

    # Stage 3: Merge
    _log(f"\n[Stage 3] Merging responses...")
    merged = merge_responses(query, results, intent)
    _log(f"  Merged answer: {len(merged)} chars")

    # Stage 4: Output
    _log(f"\n[Stage 4] Formatting output...")
    output = format_output(
        query, intent, prompts, results, merged,
        show_individual=show_individual,
        show_prompts=show_prompts,
    )

    if output_file:
        out_path = Path(output_file)
        out_path.parent.mkdir(parents=True, exist_ok=True)
        out_path.write_text(output, encoding="utf-8")
        _log(f"  Saved to: {output_file}")

    return output


def main():
    parser = argparse.ArgumentParser(
        description="Multi-Model Oracle — Get the ultimate merged answer from the best AI models"
    )
    parser.add_argument(
        "query", nargs="?", default=None,
        help="The query or prompt to send to all models",
    )
    parser.add_argument(
        "--file", "-f", default=None,
        help="Read query from a file instead of command line",
    )
    parser.add_argument(
        "--stdin", action="store_true",
        help="Read query from stdin (pipe input)",
    )
    parser.add_argument(
        "--output", "-o", default=None,
        help="Save full output to a file",
    )
    parser.add_argument(
        "--show-individual", "-i", action="store_true",
        help="Show individual model responses alongside the merged answer",
    )
    parser.add_argument(
        "--show-prompts", "-p", action="store_true",
        help="Show the prompt-engineered variants sent to each model",
    )
    parser.add_argument(
        "--raw", action="store_true",
        help="Skip prompt engineering — send the query as-is to all models",
    )
    args = parser.parse_args()

    # Determine the query source
    if args.file:
        file_path = Path(args.file)
        if not file_path.exists():
            print(f"ERROR: File not found: {args.file}", file=sys.stderr)
            sys.exit(1)
        query = file_path.read_text(encoding="utf-8").strip()
    elif args.stdin:
        query = sys.stdin.read().strip()
    elif args.query:
        query = args.query
    else:
        parser.print_help()
        sys.exit(1)

    if not query:
        print("ERROR: Empty query.", file=sys.stderr)
        sys.exit(1)

    output = run_oracle(
        query=query,
        show_individual=args.show_individual,
        show_prompts=args.show_prompts,
        output_file=args.output,
        skip_engineering=args.raw,
    )

    # Only the final output goes to stdout
    print(output)


if __name__ == "__main__":
    main()
