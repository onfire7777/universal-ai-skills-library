#!/usr/bin/env python3
"""
Prompt Engineer — Manus-Optimized Prompt Optimizer

Takes a user's raw prompt/idea and transforms it into the best possible
Manus-optimized prompt using multi-model AI analysis and Manus-specific
context engineering principles.

Usage:
    python3 optimize_prompt.py "your raw prompt here"
    python3 optimize_prompt.py --file /path/to/prompt.txt
    echo "your prompt" | python3 optimize_prompt.py --stdin

Options:
    --mode agent|chat|project   Target prompt mode (default: auto-detect)
    --depth quick|standard|deep Optimization depth (default: standard)
    --show-analysis             Show the analysis breakdown
    --show-comparison           Show before/after comparison
    -o, --output FILE           Save optimized prompt to file

Requirements:
    pip install requests openai
    Environment variables: OPENROUTER_API_KEY and/or OPENAI_API_KEY
"""

import argparse
import json
import os
import re
import sys
import time
from concurrent.futures import ThreadPoolExecutor, as_completed

# ---------------------------------------------------------------------------
# Dependency checks (DBG-014 fix from round 1)
# ---------------------------------------------------------------------------
try:
    import requests
except ImportError:
    print("Error: 'requests' library required. Install with: pip install requests", file=sys.stderr)
    sys.exit(1)

try:
    from openai import OpenAI
except ImportError:
    print("Error: 'openai' library required. Install with: pip install openai", file=sys.stderr)
    sys.exit(1)

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
SKILL_DIR = os.path.dirname(SCRIPT_DIR)

# ---------------------------------------------------------------------------
# Model configuration
# ---------------------------------------------------------------------------
OPENROUTER_API_KEY = os.environ.get("OPENROUTER_API_KEY", "")
OPENAI_API_KEY = os.environ.get("OPENAI_API_KEY", "")

MODELS = {
    "opus": {
        "name": "Anthropic Claude Opus 4",
        "id": "anthropic/claude-opus-4",
        "provider": "openrouter",
        "strengths": "nuanced reasoning, edge cases, multi-perspective analysis",
    },
    "gpt": {
        "name": "OpenAI GPT-4.1",
        "id": "openai/gpt-4.1",
        "provider": "openrouter",
        "strengths": "technical precision, systematic coverage, structured output",
    },
    "mini": {
        "name": "Manus gpt-4.1-mini",
        "id": "gpt-4.1-mini",
        "provider": "openai",
        "strengths": "speed, efficiency, prompt engineering, structured tasks",
    },
}

# Cached OpenAI client (DBG-002 fix: explicit key, single instance)
_openai_client = None


def _get_openai_client():
    global _openai_client
    if _openai_client is None:
        _openai_client = OpenAI(api_key=OPENAI_API_KEY, timeout=180)
    return _openai_client


# ---------------------------------------------------------------------------
# Intent detection (word boundary matching)
# ---------------------------------------------------------------------------
INTENT_CATEGORIES = {
    "agent_task": {
        "description": "Multi-step autonomous task requiring tool use, research, file creation, or web interaction",
        "keywords": ["build", "create", "deploy", "research", "analyze", "automate", "organize", "manage", "monitor", "plan", "develop", "design", "implement"],
        "optimization": "Structure as agent-first prompt with clear end-goal, deliverables, constraints, and success criteria. Leverage Manus tools (shell, browser, file system, search, code execution). Use imperative language.",
    },
    "code_development": {
        "description": "Software development, debugging, code review, or technical implementation",
        "keywords": ["code", "program", "debug", "fix", "implement", "function", "api", "script", "app", "website", "database"],
        "optimization": "Specify language, framework, error handling, testing requirements. Include architecture preferences, file structure, and deployment target. Be explicit about quality standards.",
    },
    "content_creation": {
        "description": "Writing, document creation, presentations, reports, or creative content",
        "keywords": ["write", "draft", "article", "report", "presentation", "blog", "email", "document", "content", "copy"],
        "optimization": "Define tone, audience, style, format, length, and structure. Include examples of desired output quality. Specify research depth and citation requirements.",
    },
    "research_analysis": {
        "description": "Deep research, data analysis, competitive analysis, or market research",
        "keywords": ["research", "analyze", "compare", "investigate", "study", "evaluate", "assess", "benchmark", "survey"],
        "optimization": "Define scope, depth, sources, evaluation criteria, and output format. Request structured findings with evidence. Specify actionable recommendations.",
    },
    "data_processing": {
        "description": "Data manipulation, visualization, spreadsheet work, or statistical analysis",
        "keywords": ["data", "spreadsheet", "chart", "graph", "csv", "excel", "statistics", "visualize", "calculate", "transform"],
        "optimization": "Specify input format, desired transformations, output format, and visualization preferences. Include sample data structure and expected results.",
    },
    "workflow_automation": {
        "description": "Process automation, scheduling, integrations, or recurring task setup",
        "keywords": ["automate", "schedule", "integrate", "workflow", "recurring", "trigger", "pipeline", "process"],
        "optimization": "Define trigger conditions, steps, error handling, and success verification. Specify integrations (MCP, APIs, email). Include rollback procedures.",
    },
    "general": {
        "description": "General questions, advice, brainstorming, or conversational queries",
        "keywords": [],
        "optimization": "Add structure, depth calibration, format preference, and specificity. Transform from chatbot-style to agent-style where possible.",
    },
}


def detect_intent(prompt: str) -> str:
    """Detect the primary intent category of a prompt using word boundary matching."""
    prompt_lower = prompt.lower()
    scores = {}
    for category, info in INTENT_CATEGORIES.items():
        score = sum(1 for kw in info["keywords"] if re.search(r'\b' + re.escape(kw) + r'\b', prompt_lower))
        scores[category] = score
    best = max(scores, key=scores.get)
    return best if scores[best] > 0 else "general"


def detect_mode(prompt: str) -> str:
    """Auto-detect whether the prompt is best suited for agent, chat, or project mode."""
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


# ---------------------------------------------------------------------------
# Utilities
# ---------------------------------------------------------------------------
def strip_code_fences(text: str) -> str:
    """Remove markdown code fences from model output."""
    text = text.strip()
    # Handle multiple code blocks — extract content from the first one
    match = re.match(r'^```(?:json)?\s*\n?(.*?)\n?```', text, re.DOTALL)
    if match:
        return match.group(1).strip()
    return text


def sanitize_error(msg: str) -> str:
    """Remove API keys from error messages."""
    error_msg = str(msg)
    if OPENROUTER_API_KEY and OPENROUTER_API_KEY in error_msg:
        error_msg = error_msg.replace(OPENROUTER_API_KEY, "[REDACTED]")
    if OPENAI_API_KEY and OPENAI_API_KEY in error_msg:
        error_msg = error_msg.replace(OPENAI_API_KEY, "[REDACTED]")
    # Also redact anything that looks like a bearer token
    error_msg = re.sub(r'(sk-[a-zA-Z0-9]{10,})', '[REDACTED]', error_msg)
    error_msg = re.sub(r'(Bearer\s+)[^\s"\']+', r'\1[REDACTED]', error_msg)
    return error_msg


def is_error_result(text) -> bool:
    """Check if a model result is an error string. Handles None safely."""
    if text is None:
        return True
    if not isinstance(text, str):
        return True
    return text.startswith("[ERROR")


def pick_model(preferred: str, available_models: list) -> str:
    """Pick the preferred model if available, otherwise the first available."""
    if preferred in available_models:
        return preferred
    return available_models[0]


# ---------------------------------------------------------------------------
# API calls
# ---------------------------------------------------------------------------
def call_openrouter(model_id: str, system_prompt: str, user_prompt: str,
                    temperature: float = 0.7, max_tokens: int = 4096) -> str:
    """Call a model via OpenRouter API with retry logic and robust parsing."""
    headers = {
        "Authorization": f"Bearer {OPENROUTER_API_KEY}",
        "Content-Type": "application/json",
        "HTTP-Referer": "https://manus.im",
        "X-Title": "Prompt Engineer Skill",
    }
    payload = {
        "model": model_id,
        "messages": [
            {"role": "system", "content": system_prompt},
            {"role": "user", "content": user_prompt},
        ],
        "temperature": temperature,
        "max_tokens": max_tokens,
    }
    for attempt in range(3):
        try:
            resp = requests.post(
                "https://openrouter.ai/api/v1/chat/completions",
                headers=headers, json=payload, timeout=180,
            )
            resp.raise_for_status()
            data = resp.json()
            if "error" in data:
                raise ValueError(f"API error: {str(data['error'])[:200]}")
            content = data.get("choices", [{}])[0].get("message", {}).get("content")
            if content is None:
                raise ValueError(f"Unexpected response structure: {json.dumps(data)[:200]}")
            return content
        except Exception as e:
            if attempt == 2:
                return f"[ERROR from {model_id}]: {sanitize_error(e)}"
            time.sleep(2 ** attempt)
    return "[ERROR]: All retries exhausted"


def call_openai(system_prompt: str, user_prompt: str,
                temperature: float = 0.7, max_tokens: int = 4096) -> str:
    """Call gpt-4.1-mini via the OpenAI-compatible API with retry logic."""
    client = _get_openai_client()
    for attempt in range(3):
        try:
            response = client.chat.completions.create(
                model="gpt-4.1-mini",
                messages=[
                    {"role": "system", "content": system_prompt},
                    {"role": "user", "content": user_prompt},
                ],
                temperature=temperature,
                max_tokens=max_tokens,
            )
            return response.choices[0].message.content
        except Exception as e:
            if attempt == 2:
                return f"[ERROR from gpt-4.1-mini]: {sanitize_error(e)}"
            time.sleep(2 ** attempt)
    return "[ERROR from gpt-4.1-mini]: All retries exhausted"


def call_model(model_key: str, system_prompt: str, user_prompt: str,
               temperature: float = 0.7) -> str:
    """Route to the correct API based on model provider."""
    model = MODELS[model_key]
    if model["provider"] == "openrouter":
        return call_openrouter(model["id"], system_prompt, user_prompt, temperature)
    else:
        return call_openai(system_prompt, user_prompt, temperature)


# ---------------------------------------------------------------------------
# Core prompt engineering pipeline
# ---------------------------------------------------------------------------

ANALYSIS_SYSTEM_PROMPT = """You are an elite prompt engineering analyst specializing in Manus AI.
Your job is to analyze a user's raw prompt and produce a structured analysis.

Manus AI is an autonomous AI agent (not a chatbot) that:
- Operates in a sandboxed Linux VM with internet access
- Has tools: shell, file system, browser, search, code execution, image generation, slides, scheduling
- Works in an agent loop: Analyze → Think → Select tool → Execute → Observe → Iterate
- Can deploy websites, run code, browse the web, create files, and automate workflows
- Uses ~50 tool calls per typical task
- Credits grow multiplicatively with follow-up prompts (one-shot is best)
- Supports Projects (persistent instructions + knowledge base) and Skills (modular capabilities)

Analyze the prompt and output valid JSON with these fields:
{
  "intent": "primary intent category",
  "mode": "agent|chat|project",
  "clarity_score": 1-10,
  "specificity_score": 1-10,
  "actionability_score": 1-10,
  "missing_elements": ["list of what's missing"],
  "ambiguities": ["list of unclear aspects"],
  "strengths": ["what's already good"],
  "manus_tools_needed": ["which Manus tools this task likely needs"],
  "suggested_structure": "brief description of optimal prompt structure",
  "complexity_level": "simple|moderate|complex|expert"
}

Output ONLY the JSON, no other text."""

OPTIMIZER_SYSTEM_PROMPT = """You are the world's best prompt engineer, specializing in crafting perfect prompts for Manus AI.

## Manus AI Context Engineering Principles (from official Manus engineering team):

1. **One-Shot Density**: Credit usage grows multiplicatively with follow-ups. Pack everything into ONE dense, thorough prompt. Even "change the color to green" costs nearly as much as the original prompt.

2. **Agent-First Thinking**: Manus is an autonomous agent, not a chatbot. Prompts should leverage multi-step execution, tool use, file creation, web research, and structured outputs. Use phrases like "Take ownership", "Guide step by step", "Track progress".

3. **Imperative Clarity**: Be imperative and specific. Manus is an "unfocused lens" — with enough adjustment, you get a clear picture. Vague prompts waste credits and produce poor results.

4. **Structured Deliverables**: Always specify the exact output format, file types, structure, and quality standards expected. Manus excels at producing structured, reusable deliverables.

5. **Context as Memory**: Manus uses the file system as externalized memory. For complex tasks, instruct it to create planning files (todo.md), save intermediate results, and organize outputs in clear directory structures.

6. **Error Prevention**: Most errors are preventable through prompting alone. Specify constraints, edge cases, fallback behaviors, and validation criteria upfront.

7. **Leverage Manus Tools**: Explicitly reference capabilities when relevant: shell commands, Python code, web browsing, search, file operations, image generation, slides, scheduling, MCP integrations, GitHub, Google Drive.

## Your Task:
Transform the user's raw prompt into the BEST possible Manus-optimized prompt. The optimized prompt must be:
- Dense and self-contained (one-shot ready)
- Specific with clear deliverables and success criteria
- Structured with logical sections
- Agent-aware (leveraging Manus's multi-step capabilities)
- Credit-efficient (preventing errors and loops)

Output ONLY the optimized prompt text, ready to paste into Manus. No meta-commentary."""

DEEP_OPTIMIZER_SYSTEM_PROMPT = """You are the world's foremost prompt engineer. You have deep expertise in:

1. **Manus AI Architecture**: Agent loop (Analyze → Think → Select tool → Execute → Observe → Iterate), sandbox VM, tool ecosystem (shell, browser, file, search, code, image gen, slides, scheduling, MCP, GitHub, Google Drive), Skills system, Projects system.

2. **Context Engineering** (from Manus's own engineering team):
   - KV-Cache optimization: stable prefixes, append-only context, deterministic serialization
   - File system as context: unlimited, persistent, externalized memory
   - Attention manipulation: todo.md recitation to keep goals in recent attention
   - Error retention: leave failures in context so the model learns and adapts
   - Anti-few-shot: introduce diversity to prevent pattern lock-in

3. **Advanced Prompting Patterns**:
   - Tree of Thoughts: explore multiple solution paths
   - Self-Consistency: generate multiple approaches, pick best
   - Socratic Questioning: clarify before solving
   - Reverse Engineering: start from desired outcome
   - First Principles: break down to fundamentals
   - Red Team / Blue Team: attack and defend strategies
   - Iterative Refinement: evolutionary improvement
   - Prompt Chaining: sequential multi-stage processing
   - Context Injection: gather context before answering
   - Progressive Disclosure: layered depth

4. **Credit Optimization**:
   - One-shot density: pack everything into a single prompt
   - Prevent loops: specify validation and fallback behaviors
   - Reduce tool calls: batch related operations
   - Explicit permissions: tell Manus when to ask vs. decide autonomously

## Your Task:
You will receive a raw prompt AND an analysis of that prompt. Transform it into the absolute best possible Manus-optimized prompt. Apply every relevant technique. The result should be a masterpiece of prompt engineering — dense, clear, structured, agent-aware, and credit-efficient.

Structure the optimized prompt with clear sections:
- **Objective**: What to accomplish (1-2 sentences)
- **Context**: Background information and constraints
- **Requirements**: Specific deliverables with quality standards
- **Process**: How to approach the task (leveraging Manus tools)
- **Output Format**: Exact structure of expected deliverables
- **Success Criteria**: How to verify the task is complete
- **Constraints**: What to avoid, edge cases, error handling

Output ONLY the optimized prompt text. No meta-commentary."""

COMPARISON_SYSTEM_PROMPT = """You are a prompt engineering evaluator. Compare the original and optimized prompts.

Score each on these dimensions (1-10):
1. Clarity: How clear and unambiguous is the prompt?
2. Specificity: How specific are the requirements and deliverables?
3. Actionability: How actionable is the prompt for an AI agent?
4. Completeness: Does it cover all necessary aspects?
5. Efficiency: Will it minimize wasted effort and credits?
6. Manus-Optimization: How well does it leverage Manus's specific capabilities?

Output valid JSON:
{
  "original_scores": {"clarity": N, "specificity": N, "actionability": N, "completeness": N, "efficiency": N, "manus_optimization": N, "total": N},
  "optimized_scores": {"clarity": N, "specificity": N, "actionability": N, "completeness": N, "efficiency": N, "manus_optimization": N, "total": N},
  "improvement_percentage": N,
  "key_improvements": ["list of specific improvements made"],
  "remaining_suggestions": ["any further improvements possible"]
}

Output ONLY the JSON."""


def analyze_prompt(raw_prompt: str, available_models: list) -> dict:
    """Stage 1: Analyze the raw prompt using the best available model."""
    print("  [1/4] Analyzing prompt...", flush=True)
    model = pick_model("mini", available_models)
    result = call_model(model, ANALYSIS_SYSTEM_PROMPT, raw_prompt, temperature=0.3)
    if is_error_result(result):
        print(f"  Warning: Analysis model failed: {result}", file=sys.stderr)
    try:
        return json.loads(strip_code_fences(result))
    except (json.JSONDecodeError, TypeError):
        return {
            "intent": detect_intent(raw_prompt),
            "mode": detect_mode(raw_prompt),
            "clarity_score": 5,
            "specificity_score": 5,
            "actionability_score": 5,
            "missing_elements": ["Could not parse AI analysis"],
            "ambiguities": [],
            "strengths": [],
            "manus_tools_needed": [],
            "suggested_structure": "standard",
            "complexity_level": "moderate",
        }


def optimize_prompt(raw_prompt: str, analysis: dict, depth: str = "standard",
                    available_models: list = None) -> dict:
    """Stage 2: Optimize the prompt. Always returns dict of model_key -> result."""
    print("  [2/4] Optimizing prompt (parallel models)...", flush=True)

    analysis_context = f"""
## Analysis of the Raw Prompt:
- Intent: {analysis.get('intent', 'unknown')}
- Mode: {analysis.get('mode', 'agent')}
- Clarity: {analysis.get('clarity_score', '?')}/10
- Specificity: {analysis.get('specificity_score', '?')}/10
- Actionability: {analysis.get('actionability_score', '?')}/10
- Missing: {', '.join(analysis.get('missing_elements', []))}
- Ambiguities: {', '.join(analysis.get('ambiguities', []))}
- Strengths: {', '.join(analysis.get('strengths', []))}
- Tools Needed: {', '.join(analysis.get('manus_tools_needed', []))}
- Complexity: {analysis.get('complexity_level', 'moderate')}

## Raw Prompt to Optimize:
{raw_prompt}
"""

    system = DEEP_OPTIMIZER_SYSTEM_PROMPT if depth == "deep" else OPTIMIZER_SYSTEM_PROMPT

    if available_models is None:
        available_models = list(MODELS.keys())

    if depth == "quick":
        model_key = pick_model("mini", available_models)
        return {model_key: call_model(model_key, system, analysis_context, temperature=0.7)}

    # Standard/Deep: parallel query across available models
    results = {}
    models_to_query = [k for k in available_models if k in MODELS]
    with ThreadPoolExecutor(max_workers=len(models_to_query)) as executor:
        futures = {}
        for key in models_to_query:
            model_system = system + f"\n\nYou are optimizing this prompt with the strengths of {MODELS[key]['name']}: {MODELS[key]['strengths']}."
            futures[executor.submit(call_model, key, model_system, analysis_context, 0.7)] = key

        for future in as_completed(futures):
            model_key = futures[future]
            try:
                results[model_key] = future.result()
            except Exception as e:
                results[model_key] = f"[ERROR]: {sanitize_error(e)}"

    return results


def merge_optimizations(results: dict, raw_prompt: str, analysis: dict,
                        available_models: list) -> str:
    """Stage 3: Merge multiple optimization results into the best single prompt."""
    print("  [3/4] Merging optimizations...", flush=True)

    # Filter out failed results (DBG-009 fix: None-safe check)
    valid_results = {k: v for k, v in results.items() if not is_error_result(v)}
    if not valid_results:
        return "Error: All optimization models failed. Please check API keys.\n" + \
               "\n".join(f"  {k}: {v}" for k, v in results.items())
    if len(valid_results) == 1:
        return next(iter(valid_results.values()))

    # Build merge prompt with only valid results
    version_texts = []
    for i, (key, text) in enumerate(valid_results.items(), 1):
        model_info = MODELS.get(key, {})
        version_texts.append(
            f"## Optimized Version {i} ({model_info.get('name', key)} — {model_info.get('strengths', 'general')}):\n{text}"
        )

    merge_prompt = f"""You have {len(valid_results)} different optimized versions of the same prompt, each created by a different AI model.
Your job is to merge them into ONE ultimate optimized prompt that takes the best elements from each.

## Original Raw Prompt:
{raw_prompt}

## Analysis:
- Intent: {analysis.get('intent', 'unknown')}
- Mode: {analysis.get('mode', 'agent')}
- Complexity: {analysis.get('complexity_level', 'moderate')}

{chr(10).join(version_texts)}

## Merge Instructions:
1. Extract the strongest elements from each version
2. Resolve any contradictions by choosing the most specific/actionable option
3. Eliminate redundancy without losing important details
4. Ensure the final prompt is dense, one-shot ready, and Manus-optimized
5. Structure with clear sections (Objective, Context, Requirements, Process, Output Format, Success Criteria, Constraints)
6. The result should read as a single, coherent, masterfully crafted prompt

Output ONLY the merged optimized prompt. No meta-commentary."""

    merge_system = """You are the world's best prompt merger. You take multiple AI-generated prompt optimizations and synthesize them into one ultimate prompt that is better than any individual version. Output only the final merged prompt."""

    # Use available models for merging (DBG-001 fix)
    merge_model = pick_model("mini", available_models)
    merged = call_model(merge_model, merge_system, merge_prompt, temperature=0.5)

    # Retry with a different model if merge fails
    if is_error_result(merged):
        fallback = [m for m in available_models if m != merge_model]
        if fallback:
            merged = call_model(fallback[0], merge_system, merge_prompt, temperature=0.5)

    # Final fallback: return the longest valid result
    if is_error_result(merged):
        merged = max(valid_results.values(), key=len)

    return merged


def compare_prompts(raw_prompt: str, optimized_prompt: str, available_models: list) -> dict:
    """Stage 4: Compare original vs optimized prompt."""
    print("  [4/4] Evaluating improvement...", flush=True)

    comparison_input = f"""## Original Prompt:
{raw_prompt}

## Optimized Prompt:
{optimized_prompt}"""

    model = pick_model("mini", available_models)
    result = call_model(model, COMPARISON_SYSTEM_PROMPT, comparison_input, temperature=0.3)
    try:
        return json.loads(strip_code_fences(result))
    except (json.JSONDecodeError, TypeError):
        return {"error": "Could not parse comparison", "raw": result}


# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------
def main():
    parser = argparse.ArgumentParser(
        description="Manus-Optimized Prompt Engineer — Transform any prompt into the best possible Manus prompt"
    )
    parser.add_argument("prompt", nargs="?", help="The raw prompt to optimize")
    parser.add_argument("--file", "-f", help="Read prompt from a file")
    parser.add_argument("--stdin", action="store_true",
                        help="Read prompt from stdin (requires piped input)")
    parser.add_argument("--mode", choices=["agent", "chat", "project", "auto"],
                        default="auto", help="Target prompt mode (default: auto-detect)")
    parser.add_argument("--depth", choices=["quick", "standard", "deep"],
                        default="standard", help="Optimization depth (default: standard)")
    parser.add_argument("--show-analysis", action="store_true",
                        help="Show the analysis breakdown")
    parser.add_argument("--show-comparison", action="store_true",
                        help="Show before/after comparison scores")
    parser.add_argument("--show-individual", action="store_true",
                        help="Show individual model optimizations before merge")
    parser.add_argument("-o", "--output", help="Save optimized prompt to file")
    args = parser.parse_args()

    # Get the raw prompt (with error handling)
    if args.stdin:
        if sys.stdin.isatty():
            print("Reading from stdin (press Ctrl+D when done):", file=sys.stderr)
        raw_prompt = sys.stdin.read().strip()
    elif args.file:
        try:
            with open(args.file, "r", encoding="utf-8") as f:
                raw_prompt = f.read().strip()
        except FileNotFoundError:
            print(f"Error: File not found: {args.file}", file=sys.stderr)
            sys.exit(1)
        except (PermissionError, OSError) as e:
            print(f"Error: Cannot read file {args.file}: {e}", file=sys.stderr)
            sys.exit(1)
    elif args.prompt:
        raw_prompt = args.prompt
    else:
        parser.print_help()
        sys.exit(1)

    if not raw_prompt:
        print("Error: Empty prompt provided.", file=sys.stderr)
        sys.exit(1)

    # Validate API keys and determine available models
    available_models = []
    if OPENROUTER_API_KEY:
        available_models.extend(["opus", "gpt"])
    if OPENAI_API_KEY:
        available_models.append("mini")

    if not available_models:
        print("Error: No API keys configured. Set OPENROUTER_API_KEY and/or OPENAI_API_KEY.", file=sys.stderr)
        sys.exit(1)

    if args.depth != "quick" and len(available_models) < 2:
        print(f"Warning: Only {len(available_models)} model(s) available ({', '.join(available_models)}). Falling back to quick mode.", file=sys.stderr)
        args.depth = "quick"

    start_time = time.time()
    print(f"\n{'='*70}")
    print(f"  MANUS PROMPT ENGINEER — Optimizing Your Prompt")
    print(f"  Depth: {args.depth} | Mode: {args.mode} | Models: {len(available_models)}")
    print(f"{'='*70}\n")

    # Stage 1: Analyze
    analysis = analyze_prompt(raw_prompt, available_models)

    if args.mode != "auto":
        analysis["mode"] = args.mode

    if args.show_analysis:
        print(f"\n{'─'*50}")
        print("  PROMPT ANALYSIS")
        print(f"{'─'*50}")
        for key, value in analysis.items():
            if isinstance(value, list):
                print(f"  {key}: {', '.join(str(v) for v in value) if value else '(none)'}")
            else:
                print(f"  {key}: {value}")
        print()

    # Stage 2: Optimize (always returns dict)
    results = optimize_prompt(raw_prompt, analysis, args.depth, available_models)

    if args.depth == "quick":
        # Quick mode: extract the single result
        optimized = next(iter(results.values()))
        if is_error_result(optimized):
            print(f"Error: Optimization failed: {optimized}", file=sys.stderr)
            sys.exit(1)
    else:
        if args.show_individual:
            for model_key, result in results.items():
                model_info = MODELS.get(model_key, {})
                print(f"\n{'─'*50}")
                print(f"  {model_info.get('name', model_key)} Optimization:")
                print(f"{'─'*50}")
                print(result)
                print()

        # Stage 3: Merge
        optimized = merge_optimizations(results, raw_prompt, analysis, available_models)
        if optimized.startswith("Error: All optimization models failed"):
            print(optimized, file=sys.stderr)
            sys.exit(1)

    # Stage 4: Compare (optional)
    comparison = None
    if args.show_comparison:
        comparison = compare_prompts(raw_prompt, optimized, available_models)

    elapsed = time.time() - start_time

    # Output
    print(f"\n{'='*70}")
    print(f"  OPTIMIZED PROMPT")
    print(f"{'='*70}\n")
    print(optimized)

    if comparison and "error" not in comparison:
        print(f"\n{'─'*50}")
        print("  IMPROVEMENT SCORES")
        print(f"{'─'*50}")
        orig = comparison.get("original_scores", {})
        opt = comparison.get("optimized_scores", {})
        print(f"  {'Dimension':<22} {'Original':>10} {'Optimized':>10} {'Change':>10}")
        print(f"  {'─'*52}")
        for dim in ["clarity", "specificity", "actionability", "completeness", "efficiency", "manus_optimization"]:
            o = orig.get(dim, "?")
            n = opt.get(dim, "?")
            change = ""
            if isinstance(o, (int, float)) and isinstance(n, (int, float)):
                diff = n - o
                change = f"+{diff}" if diff > 0 else str(diff)
            print(f"  {dim:<22} {str(o):>10} {str(n):>10} {change:>10}")
        print(f"\n  Overall improvement: {comparison.get('improvement_percentage', '?')}%")
        if comparison.get("key_improvements"):
            print(f"\n  Key improvements:")
            for imp in comparison["key_improvements"]:
                print(f"    + {imp}")
        print()

    print(f"\n{'─'*50}")
    print(f"  Completed in {elapsed:.1f}s | Depth: {args.depth} | Mode: {analysis.get('mode', 'auto')}")
    print(f"  Intent: {analysis.get('intent', 'unknown')} | Complexity: {analysis.get('complexity_level', 'unknown')}")
    print(f"  Models available: {', '.join(available_models)}")
    print(f"{'─'*50}\n")

    # Save to file if requested (with error handling)
    if args.output:
        try:
            out_dir = os.path.dirname(os.path.abspath(args.output))
            if out_dir:
                os.makedirs(out_dir, exist_ok=True)
            with open(args.output, "w", encoding="utf-8") as f:
                f.write(optimized)
            print(f"  Saved to: {args.output}\n")
        except OSError as e:
            print(f"Error: Could not save to {args.output}: {e}", file=sys.stderr)


if __name__ == "__main__":
    main()
