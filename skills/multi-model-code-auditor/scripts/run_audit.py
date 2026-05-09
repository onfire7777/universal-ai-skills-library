#!/usr/bin/env python3
"""Multi-Model Code Auditor — Elite Edition

Sends a full codebase to the absolute best frontier coding models via
OpenRouter API, using model-specific prompt engineering optimized for
each model's strengths. Dynamically discovers the latest frontier models.

Usage:
    python3 run_audit.py <project_dir> [--output audit_results.json] [--discover]
    python3 run_audit.py <project_dir> --parallel          # Query all models simultaneously
    python3 run_audit.py <project_dir> --models id1,id2    # Use specific models
    python3 run_audit.py --discover                        # List available models

Requires: OPENROUTER_API_KEY environment variable.
"""

import argparse
import json
import os
import re
import sys
import threading
import time
from concurrent.futures import ThreadPoolExecutor, as_completed
from pathlib import Path

try:
    import requests
except ImportError:
    import subprocess
    subprocess.check_call([sys.executable, "-m", "pip", "install", "requests", "-q"])
    import requests

# ─── Constants ────────────────────────────────────────────────────────────────

API_URL = "https://openrouter.ai/api/v1/chat/completions"
MODELS_URL = "https://openrouter.ai/api/v1/models"

# Timeouts: (connect_timeout, read_timeout) in seconds
DEFAULT_CONNECT_TIMEOUT = 30
DEFAULT_READ_TIMEOUT = 300

# Retry configuration
MAX_RETRIES = 3
RETRY_DELAYS = [10, 30, 60]  # Exponential backoff delays in seconds

# Payload limits (in characters, ~4 chars per token)
MAX_PAYLOAD_CHARS = 600_000  # ~150K tokens — safe for most models
WARN_PAYLOAD_CHARS = 400_000

SOURCE_EXTENSIONS = {
    ".py", ".js", ".ts", ".tsx", ".jsx", ".java", ".go", ".rs", ".rb",
    ".php", ".c", ".cpp", ".h", ".hpp", ".cs", ".swift", ".kt", ".scala",
    ".sh", ".bash", ".bat", ".ps1", ".vbs", ".yaml", ".yml", ".toml",
    ".json", ".xml", ".html", ".css", ".sql", ".r", ".lua", ".dart",
    ".vue", ".svelte", ".ex", ".exs", ".zig",
}

SKIP_DIRS = {
    "node_modules", ".git", "__pycache__", ".venv", "venv", "env",
    ".tox", ".mypy_cache", ".pytest_cache", "dist", "build", ".next",
    "target", "bin", "obj", ".idea", ".vscode", ".cache", "coverage",
}

# ─── Frontier Model Tiers ────────────────────────────────────────────────────
# Ordered by audit quality. The script picks the best available from each tier.
# Each tier targets a different model family to maximize diversity of analysis.

MODEL_TIERS = [
    {
        "family": "openai_flagship",
        "role": "Structured Analysis Lead",
        "description": "Best at precise structured output, schema adherence, systematic enumeration",
        "candidates": [
            "openai/gpt-5.4-pro",
            "openai/gpt-5.4",
            "openai/gpt-5.2-pro",
            "openai/gpt-5.2",
            "openai/gpt-5.1",
            "openai/gpt-4.1",
        ],
    },
    {
        "family": "anthropic_flagship",
        "role": "Nuanced Reasoning Specialist",
        "description": "Best at contextual understanding, subtle logic bugs, fewer hallucinations",
        "candidates": [
            "anthropic/claude-opus-4.6",
            "anthropic/claude-sonnet-4.6",
            "anthropic/claude-sonnet-4.5",
        ],
    },
    {
        "family": "google_flagship",
        "role": "Cross-File Architecture Analyst",
        "description": "Best at large context analysis, cross-file dependency tracking, architectural issues",
        "candidates": [
            "google/gemini-3.1-pro-preview",
            "google/gemini-2.5-pro",
            "google/gemini-2.5-pro-preview",
        ],
    },
    {
        "family": "openai_reasoning",
        "role": "Deep Reasoning & Logic Auditor",
        "description": "Best at complex logic chains, race conditions, subtle state machine bugs",
        "candidates": [
            "openai/o4-mini-high",
            "openai/o4-mini",
            "openai/o3-pro",
            "openai/o3-mini-high",
            "openai/o3-mini",
        ],
    },
    {
        "family": "openai_codex",
        "role": "Code-Native Security Scanner",
        "description": "Best at code-level pattern recognition, API misuse, unsafe idioms",
        "candidates": [
            "openai/gpt-5.3-codex",
            "openai/gpt-5.2-codex",
            "openai/gpt-5.1-codex",
        ],
    },
    {
        "family": "deepseek",
        "role": "Algorithmic & Low-Level Bug Hunter",
        "description": "Best at algorithmic correctness, off-by-one errors, data structure misuse",
        "candidates": [
            "deepseek/deepseek-v3.2",
            "deepseek/deepseek-v3.1-terminus",
            "deepseek/deepseek-r1-0528",
            "deepseek/deepseek-r1",
        ],
    },
]

# Known-good fallback models — used when OpenRouter model discovery API fails.
# These are stable, widely-available models that are known to work reliably.
FALLBACK_MODELS = [
    {"id": "anthropic/claude-sonnet-4", "family": "anthropic_flagship",
     "role": "Nuanced Reasoning Specialist", "description": "Fallback — stable Anthropic model"},
    {"id": "google/gemini-2.5-pro-preview", "family": "google_flagship",
     "role": "Cross-File Architecture Analyst", "description": "Fallback — stable Google model"},
    {"id": "openai/gpt-4.1", "family": "openai_flagship",
     "role": "Structured Analysis Lead", "description": "Fallback — stable OpenAI model"},
]

# Family inference from model ID prefix
FAMILY_PREFIXES = {
    "openai/o": "openai_reasoning",
    "openai/gpt-5.3-codex": "openai_codex",
    "openai/gpt-5.2-codex": "openai_codex",
    "openai/gpt-5.1-codex": "openai_codex",
    "openai/codex": "openai_codex",
    "openai/": "openai_flagship",
    "anthropic/": "anthropic_flagship",
    "google/": "google_flagship",
    "deepseek/": "deepseek",
    "meta/": "openai_flagship",  # Use flagship prompts as default
}

# ─── Model-Specific Prompt Engineering ────────────────────────────────────────

# Base system prompt — shared foundation for all models
SYSTEM_PROMPT_BASE = """You are an elite security auditor and principal software engineer with 20+ years of experience in application security, penetration testing, and secure software development. You have deep expertise in OWASP Top 10, CWE/SANS Top 25, MITRE ATT&CK, and real-world exploit development.

Your audit methodology:
1. First, build a mental model of the entire codebase architecture — entry points, data flows, trust boundaries, external interfaces
2. Then systematically analyze each attack surface using adversarial thinking
3. For each finding, trace the full attack chain from input to impact
4. Only report findings where you can cite the EXACT vulnerable code
5. Calibrate severity based on real-world exploitability, not theoretical risk"""

# Few-shot examples — calibrate severity and format for all models
FEW_SHOT_EXAMPLES = """
Here are calibration examples showing the expected severity levels and output quality:

EXAMPLE FINDING (Critical):
{
  "id": "F-01",
  "severity": "critical",
  "category": "injection",
  "file": "utils/shell.py",
  "function": "run_command",
  "line_range": "45-47",
  "title": "OS command injection via unsanitized user input",
  "description": "The run_command function passes user-controlled input directly to subprocess with shell=True. An attacker can inject arbitrary commands via semicolons, pipes, or backticks in the 'cmd' parameter. This is exploitable from the /api/execute endpoint which accepts user input without validation.",
  "attack_chain": "User input → /api/execute → run_command(user_input) → subprocess.Popen(shell=True) → arbitrary command execution",
  "vulnerable_code": "subprocess.Popen(f'process {user_input}', shell=True)",
  "fix_code": "subprocess.run(['process', user_input], shell=False, check=True)",
  "cwe": "CWE-78",
  "cvss_estimate": 9.8
}

EXAMPLE FINDING (Medium):
{
  "id": "F-02",
  "severity": "medium",
  "category": "error_handling",
  "file": "api/handlers.py",
  "function": "handle_upload",
  "line_range": "112-115",
  "title": "Broad exception handler silently swallows file corruption errors",
  "description": "The upload handler catches all exceptions with a bare 'except Exception' and returns a generic 200 OK. If the file is corrupted during write, the caller receives success but the file is truncated. This can lead to data loss that goes undetected.",
  "attack_chain": "File upload → disk full during write → exception caught → 200 OK returned → user believes upload succeeded → data lost",
  "vulnerable_code": "except Exception:\\n    return {'status': 'ok'}",
  "fix_code": "except OSError as e:\\n    logger.error(f'Upload failed: {e}')\\n    return {'status': 'error', 'message': 'Upload failed'}, 500",
  "cwe": "CWE-755",
  "cvss_estimate": 5.3
}

EXAMPLE FALSE POSITIVE (do NOT report these):
- Claiming eval() exists when it doesn't — always verify the function actually exists
- Claiming yaml.safe_load is unsafe — safe_load IS the safe variant
- Citing line numbers that don't match the actual code
- Reporting theoretical issues with no actual attack path in this codebase
"""

# Output schema — identical for all models
OUTPUT_SCHEMA = """
Respond with ONLY valid JSON (no markdown fences, no commentary). Use this exact schema:
{
  "architecture_summary": "2-3 sentence summary of the codebase architecture, entry points, and trust boundaries",
  "attack_surface": ["list of identified attack surfaces and external interfaces"],
  "findings": [
    {
      "id": "F-XX",
      "severity": "critical|high|medium|low",
      "category": "injection|path_traversal|ssrf|auth|crypto|privacy|race_condition|error_handling|logic_bug|resource_leak|config|dependency",
      "file": "relative/path/to/file.py",
      "function": "exact_function_name",
      "line_range": "start-end",
      "title": "Concise title (max 80 chars)",
      "description": "What the vulnerability is, why it matters, and how it's exploitable in this specific codebase",
      "attack_chain": "Step-by-step path from attacker input to impact",
      "vulnerable_code": "The EXACT problematic code (max 3 lines, copy-pasted from source)",
      "fix_code": "Complete working fix code that can be directly applied",
      "cwe": "CWE-XXX",
      "cvss_estimate": 0.0
    }
  ],
  "summary": {
    "total_findings": 0,
    "critical": 0,
    "high": 0,
    "medium": 0,
    "low": 0,
    "overall_assessment": "One paragraph professional assessment of the codebase security posture"
  }
}"""

# Model-specific prompt customizations
MODEL_PROMPTS = {
    "openai_flagship": {
        "system_suffix": """

Your specific strengths for this audit:
- Systematic enumeration: methodically check every function for every vulnerability class
- Schema precision: your output MUST perfectly match the JSON schema with zero deviations
- Completeness: do not stop early — audit every single file and function in the codebase
- Be exhaustive but precise: report everything real, nothing imaginary""",
        "user_prefix": """Perform a complete security audit of this codebase. Be systematic — go file by file, function by function. For each function, check against all vulnerability categories. Your output must be perfectly valid JSON matching the schema exactly.

CRITICAL RULES:
1. Every finding MUST include the exact vulnerable code copied from the source
2. Every fix_code MUST be a complete, working replacement — not pseudocode
3. Do NOT hallucinate functions or line numbers that don't exist in the code
4. Do NOT report yaml.safe_load as unsafe — it IS the safe variant
5. Severity calibration: Critical = remote code execution or data breach. High = privilege escalation or significant data exposure. Medium = denial of service or information disclosure. Low = code quality or defense-in-depth.
""",
        "temperature": 0.05,
    },
    "anthropic_flagship": {
        "system_suffix": """

Your specific strengths for this audit:
- Nuanced contextual reasoning: understand the INTENT of the code, not just the syntax
- Low hallucination rate: only report what you can verify exists in the actual code
- Subtle bug detection: find logic errors that other models miss
- Privacy awareness: identify data leakage paths that are easy to overlook
- Think step-by-step before each finding: trace the data flow from source to sink""",
        "user_prefix": """I need your most careful, thorough security analysis of this codebase. You excel at understanding context and finding subtle issues that other analyzers miss.

For each potential finding, use this mental process:
1. Identify the data flow: where does untrusted input enter? Where does it go?
2. Identify the trust boundary: is there a point where unvalidated data crosses into a privileged context?
3. Verify the code exists: re-read the source to confirm the function and code snippet are real
4. Assess exploitability: is there an actual attack path, or is this theoretical only?
5. Only report if steps 1-4 all confirm a real issue

ANTI-HALLUCINATION RULES:
- If you're not 100% certain a function exists, do NOT report it
- If you can't find the exact vulnerable code in the source, do NOT report it
- Prefer fewer, high-confidence findings over many uncertain ones
""",
        "temperature": 0.1,
    },
    "google_flagship": {
        "system_suffix": """

Your specific strengths for this audit:
- Massive context window: analyze the ENTIRE codebase holistically, not file-by-file
- Cross-file analysis: find vulnerabilities that span multiple files and modules
- Architecture-level issues: identify design flaws, missing security layers, broken trust boundaries
- Dependency chain analysis: trace how data flows across module boundaries
- Import/export analysis: find mismatched APIs between modules""",
        "user_prefix": """Analyze this entire codebase as a unified system. Your key advantage is seeing the big picture — how all the modules interact, where data flows across boundaries, and where architectural security gaps exist.

FOCUS AREAS:
1. Cross-module vulnerabilities: data that's sanitized in one module but used unsafely in another
2. Missing security layers: authentication, authorization, input validation gaps in the overall architecture
3. Trust boundary violations: where does untrusted data cross into trusted contexts across files?
4. API contract mismatches: where do callers pass data that callees don't validate?
5. Configuration security: are defaults secure? Can config values create vulnerabilities?

RULES:
- Cite the EXACT file and function for every finding
- For cross-file issues, trace the full path: file1.func_a() → file2.func_b() → vulnerability
- Do NOT report issues you cannot verify in the actual source code
""",
        "temperature": 0.1,
    },
    "openai_reasoning": {
        "system_suffix": """

Your specific strengths for this audit:
- Deep logical reasoning: find race conditions, state machine bugs, and subtle timing issues
- Multi-step attack chains: trace complex exploitation paths that require multiple steps
- Concurrency analysis: identify thread safety issues, TOCTOU races, deadlocks
- Edge case analysis: find bugs triggered by unusual inputs, boundary values, or error paths
- Think deeply about each issue before reporting — use your reasoning capabilities fully""",
        "user_prefix": """Apply your deep reasoning capabilities to find the most subtle and dangerous vulnerabilities in this codebase. Focus on issues that require multi-step logical analysis to discover.

PRIORITY TARGETS:
1. Race conditions and TOCTOU: check-then-act patterns, file operations without locking, shared mutable state
2. State machine bugs: can the system reach an invalid state through unusual event ordering?
3. Complex attack chains: vulnerabilities that require 2+ steps to exploit
4. Error path exploitation: what happens when operations fail partway through?
5. Integer overflow/underflow, off-by-one errors, boundary conditions

Think through each potential issue step by step. Only report findings where your reasoning confirms a real vulnerability with a concrete attack path.
""",
        "temperature": 0.1,
    },
    "openai_codex": {
        "system_suffix": """

Your specific strengths for this audit:
- Code-native understanding: you think in code, not just about code
- Pattern recognition: identify unsafe coding patterns, anti-patterns, and API misuse
- Fix quality: your fix_code should be production-ready, not pseudocode
- Idiomatic security: know the secure way to do things in each language
- Library expertise: know which stdlib/library functions are safe vs unsafe""",
        "user_prefix": """Audit this codebase with your code-native expertise. For every finding, provide production-ready fix code that can be directly applied.

FOCUS AREAS:
1. Unsafe function usage: shell=True, eval, exec, pickle.loads, yaml.load (without safe_load)
2. Missing input validation: unvalidated user input reaching sensitive operations
3. Insecure defaults: functions called without security-relevant parameters
4. Resource management: unclosed file handles, unbounded allocations, missing timeouts
5. Error handling: bare except, swallowed exceptions, missing error propagation

FOR EACH FINDING:
- The vulnerable_code MUST be copied exactly from the source (not paraphrased)
- The fix_code MUST be a complete, drop-in replacement that compiles/runs
- Include all necessary imports in the fix
""",
        "temperature": 0.05,
    },
    "deepseek": {
        "system_suffix": """

Your specific strengths for this audit:
- Deep code comprehension: understand algorithmic intent and find logic errors
- Low-level analysis: buffer handling, encoding issues, numeric precision
- Algorithmic correctness: verify that algorithms do what they claim
- Data structure misuse: wrong data structure for the access pattern, missing bounds checks
- Performance-security intersection: find DoS vectors through algorithmic complexity""",
        "user_prefix": """Perform a thorough code audit focusing on correctness, security, and robustness. Use your deep code understanding to find issues that surface-level analysis would miss.

FOCUS AREAS:
1. Logic errors: does the code actually do what the developer intended?
2. Algorithmic complexity: are there O(n²) or worse operations on user-controlled input? (DoS vector)
3. Data handling: encoding issues, truncation, precision loss, type confusion
4. Boundary conditions: empty inputs, maximum values, negative numbers, Unicode edge cases
5. Missing validation: assumptions about input that aren't enforced

RULES:
- Only report issues you can verify in the actual source code
- Provide exact code snippets copied from the source
- Fix code must be complete and correct
""",
        "temperature": 0.1,
    },
}


# ─── Core Functions ───────────────────────────────────────────────────────────

def infer_family(model_id: str) -> str:
    """Infer the model family from a model ID for correct prompt selection."""
    model_lower = model_id.lower()
    # Check specific prefixes first (longer matches take priority)
    for prefix in sorted(FAMILY_PREFIXES.keys(), key=len, reverse=True):
        if model_lower.startswith(prefix.lower()):
            return FAMILY_PREFIXES[prefix]
    return "openai_flagship"  # Default fallback


def discover_best_models(api_key: str) -> list[dict]:
    """Dynamically discover the best available frontier model from each family.

    When the OpenRouter API is unreachable, falls back to a curated list of
    known-good stable models instead of blindly selecting the first candidate.
    """
    available = set()
    api_ok = False
    try:
        resp = requests.get(
            MODELS_URL,
            headers={"Authorization": f"Bearer {api_key}"},
            timeout=15,
        )
        resp.raise_for_status()
        available = {m["id"] for m in resp.json().get("data", [])}
        api_ok = True
        print(f"  OpenRouter API: {len(available)} models available")
    except Exception as e:
        print(f"  WARNING: Could not fetch model list ({e})")
        print(f"  Using known-good fallback models")

    if not api_ok:
        return list(FALLBACK_MODELS)

    selected = []
    for tier in MODEL_TIERS:
        for candidate in tier["candidates"]:
            if candidate in available:
                selected.append({
                    "id": candidate,
                    "family": tier["family"],
                    "role": tier["role"],
                    "description": tier["description"],
                })
                break
        else:
            print(f"  WARNING: No model available for {tier['family']} tier")

    # If we found very few models, supplement with fallbacks
    if len(selected) < 2:
        print(f"  WARNING: Only {len(selected)} tier models found. Adding fallbacks.")
        existing_families = {m["family"] for m in selected}
        for fb in FALLBACK_MODELS:
            if fb["family"] not in existing_families:
                selected.append(dict(fb))
                existing_families.add(fb["family"])

    return selected


# Maximum file size to read (1 MB) — skip generated/minified files
MAX_FILE_SIZE = 1_000_000


def collect_source_files(project_dir: Path) -> list[dict]:
    """Collect all source files, sorted by importance (entry points first)."""
    resolved_root = project_dir.resolve()
    files = []
    for path in sorted(project_dir.rglob("*")):
        if any(skip in path.parts for skip in SKIP_DIRS):
            continue
        if path.is_file() and path.suffix.lower() in SOURCE_EXTENSIONS:
            # Guard against symlinks pointing outside the project
            try:
                resolved = path.resolve()
                if not str(resolved).startswith(str(resolved_root)):
                    print(f"  SKIP (symlink escape): {path}", file=sys.stderr)
                    continue
            except (OSError, ValueError):
                continue

            # Guard against excessively large files
            try:
                file_size = path.stat().st_size
                if file_size > MAX_FILE_SIZE:
                    print(f"  SKIP (>{MAX_FILE_SIZE // 1_000_000}MB): {path}", file=sys.stderr)
                    continue
            except OSError:
                continue

            try:
                content = path.read_text(encoding="utf-8", errors="replace")
                if content.strip():
                    rel = str(path.relative_to(project_dir))
                    # Priority scoring: entry points and config first
                    priority = 5
                    name_lower = path.name.lower()
                    if any(k in name_lower for k in ["main", "run", "app", "server", "index"]):
                        priority = 1
                    elif any(k in name_lower for k in ["config", "settings", "auth", "login"]):
                        priority = 2
                    elif any(k in name_lower for k in ["api", "route", "handler", "view"]):
                        priority = 3
                    elif any(k in name_lower for k in ["model", "schema", "database", "db"]):
                        priority = 4
                    files.append({
                        "path": rel,
                        "content": content,
                        "priority": priority,
                        "size": len(content),
                    })
            except Exception as e:
                print(f"  WARN: Could not read {path}: {e}", file=sys.stderr)
                continue

    files.sort(key=lambda f: (f["priority"], f["path"]))
    return files


def build_payload(files: list[dict], max_chars: int = MAX_PAYLOAD_CHARS) -> str:
    """Build the codebase payload with file markers.

    If the total payload exceeds max_chars, truncates low-priority files first,
    preserving high-priority files (entry points, config, auth) in full.
    """
    # First pass: calculate total size
    total = sum(len(f["content"]) + len(f["path"]) + 130 for f in files)  # 130 for markers

    if total <= max_chars:
        # Everything fits — no truncation needed
        parts = []
        for f in files:
            parts.append(f"{'='*60}")
            parts.append(f"FILE: {f['path']}")
            parts.append(f"{'='*60}")
            parts.append(f["content"])
            parts.append("")
        return "\n".join(parts)

    # Need to truncate — keep high-priority files, truncate low-priority
    print(f"  Payload too large ({total:,} chars). Smart-truncating to {max_chars:,} chars...")

    # Reserve space for high-priority files (priority 1-3)
    high_pri = [f for f in files if f["priority"] <= 3]
    low_pri = [f for f in files if f["priority"] > 3]

    high_size = sum(len(f["content"]) + 130 for f in high_pri)
    remaining = max_chars - high_size

    if remaining < 0:
        # Even high-priority files are too large — truncate them too
        print(f"  WARNING: Even high-priority files exceed limit. Truncating all files.")
        per_file = max_chars // len(files)
        for f in files:
            if len(f["content"]) > per_file:
                f["content"] = f["content"][:per_file] + "\n# ... [TRUNCATED — file too large for audit context] ..."
    else:
        # Distribute remaining space among low-priority files
        if low_pri:
            per_file = remaining // len(low_pri)
            included_low = []
            for f in low_pri:
                if per_file > 500:  # Only include if we can show meaningful content
                    if len(f["content"]) > per_file:
                        f["content"] = f["content"][:per_file] + "\n# ... [TRUNCATED] ..."
                    included_low.append(f)
            low_pri = included_low
        files = high_pri + low_pri

    parts = []
    for f in files:
        parts.append(f"{'='*60}")
        parts.append(f"FILE: {f['path']}")
        parts.append(f"{'='*60}")
        parts.append(f["content"])
        parts.append("")
    return "\n".join(parts)


def build_prompt(model_info: dict, payload: str) -> list[dict]:
    """Build the model-specific prompt with elite prompt engineering."""
    family = model_info["family"]
    prompt_config = MODEL_PROMPTS.get(family, MODEL_PROMPTS["openai_flagship"])

    system_msg = SYSTEM_PROMPT_BASE + prompt_config["system_suffix"]
    user_msg = (
        prompt_config["user_prefix"]
        + "\n"
        + FEW_SHOT_EXAMPLES
        + "\n"
        + OUTPUT_SCHEMA
        + "\n\n"
        + "=" * 70
        + "\nCODEBASE TO AUDIT\n"
        + "=" * 70
        + "\n\n"
        + payload
    )

    return [
        {"role": "system", "content": system_msg},
        {"role": "user", "content": user_msg},
    ]


def parse_response(content: str) -> dict:
    """Robustly parse JSON from model response, handling various formats."""
    content = content.strip()

    # Strip markdown fences
    if content.startswith("```"):
        content = re.sub(r"^```(?:json)?\s*\n?", "", content)
        content = re.sub(r"\n?```\s*$", "", content)
        content = content.strip()

    # Try direct parse
    try:
        return json.loads(content)
    except json.JSONDecodeError:
        pass

    # Try to find JSON object in the response
    brace_start = content.find("{")
    brace_end = content.rfind("}")
    if brace_start >= 0 and brace_end > brace_start:
        try:
            return json.loads(content[brace_start : brace_end + 1])
        except json.JSONDecodeError:
            pass

    # Try to find findings array directly
    bracket_start = content.find("[")
    bracket_end = content.rfind("]")
    if bracket_start >= 0 and bracket_end > bracket_start:
        try:
            findings = json.loads(content[bracket_start : bracket_end + 1])
            if isinstance(findings, list):
                return {"findings": findings, "summary": {"overall_assessment": "Parsed from array"}}
        except json.JSONDecodeError:
            pass

    return {"findings": [], "summary": {"overall_assessment": "Failed to parse model response as JSON"}}


def _progress_indicator(model_id: str, stop_event: threading.Event):
    """Print periodic progress dots while waiting for a model response."""
    elapsed = 0
    while not stop_event.is_set():
        stop_event.wait(30)
        if not stop_event.is_set():
            elapsed += 30
            print(f"    ... {model_id} still processing ({elapsed}s elapsed)")


def query_model(model_info: dict, payload: str, api_key: str,
                connect_timeout: int = DEFAULT_CONNECT_TIMEOUT,
                read_timeout: int = DEFAULT_READ_TIMEOUT,
                show_progress: bool = True) -> dict:
    """Send audit request to a single model with retry logic and progress indication.

    Retries up to MAX_RETRIES times with exponential backoff on transient errors
    (timeouts, 429 rate limits, 5xx server errors).
    """
    family = model_info["family"]
    prompt_config = MODEL_PROMPTS.get(family, MODEL_PROMPTS["openai_flagship"])
    messages = build_prompt(model_info, payload)

    headers = {
        "Authorization": f"Bearer {api_key}",
        "Content-Type": "application/json",
        "HTTP-Referer": "https://manus.im",
    }

    body = {
        "model": model_info["id"],
        "messages": messages,
        "temperature": prompt_config.get("temperature", 0.1),
        "max_tokens": 32000,
    }

    # Reasoning models don't support temperature
    if family == "openai_reasoning":
        body.pop("temperature", None)

    last_error = None

    for attempt in range(MAX_RETRIES):
        if attempt > 0:
            delay = RETRY_DELAYS[min(attempt - 1, len(RETRY_DELAYS) - 1)]
            print(f"    Retry {attempt}/{MAX_RETRIES-1} for {model_info['id']} in {delay}s...")
            time.sleep(delay)

        # Start progress indicator thread
        stop_event = threading.Event()
        progress_thread = None
        if show_progress:
            progress_thread = threading.Thread(
                target=_progress_indicator,
                args=(model_info["id"], stop_event),
                daemon=True,
            )
            progress_thread.start()

        try:
            resp = requests.post(
                API_URL,
                headers=headers,
                json=body,
                timeout=(connect_timeout, read_timeout),
            )

            stop_event.set()

            # Check for retryable HTTP errors
            if resp.status_code == 429:
                last_error = f"Rate limited (429)"
                print(f"    Rate limited by {model_info['id']}")
                continue
            if resp.status_code >= 500:
                last_error = f"Server error ({resp.status_code})"
                print(f"    Server error {resp.status_code} from {model_info['id']}")
                continue

            resp.raise_for_status()

            try:
                data = resp.json()
            except json.JSONDecodeError:
                last_error = "Response was not valid JSON"
                print(f"    Invalid JSON response from {model_info['id']}")
                continue

            # Check for API-level errors in the response body
            if "error" in data:
                error_msg = data["error"].get("message", str(data["error"]))
                last_error = f"API error: {error_msg}"
                print(f"    API error from {model_info['id']}: {error_msg}")
                # Don't retry on auth/payment errors
                if any(k in error_msg.lower() for k in ["auth", "key", "payment", "billing", "quota"]):
                    break
                continue

            # Check for empty/missing choices
            choices = data.get("choices", [])
            if not choices:
                last_error = "Empty choices in response"
                print(f"    Empty response from {model_info['id']}")
                continue

            content = choices[0].get("message", {}).get("content", "")
            if not content.strip():
                last_error = "Empty content in response"
                print(f"    Empty content from {model_info['id']}")
                continue

            result = parse_response(content)

            # Dedup findings within this model (same file + same function + similar title)
            seen = set()
            deduped_findings = []
            for f in result.get("findings", []):
                key = (
                    (f.get("file") or "").lower(),
                    (f.get("function") or "").lower(),
                    (f.get("cwe") or "").upper(),
                )
                if key not in seen:
                    seen.add(key)
                    deduped_findings.append(f)
            result["findings"] = deduped_findings

            return {
                "model": model_info["id"],
                "family": family,
                "role": model_info["role"],
                "status": "success",
                "findings_count": len(deduped_findings),
                "result": result,
            }

        except requests.exceptions.Timeout:
            stop_event.set()
            last_error = f"Timeout after {connect_timeout}+{read_timeout}s"
            print(f"    Timeout from {model_info['id']} (attempt {attempt+1}/{MAX_RETRIES})")
            continue

        except requests.exceptions.ConnectionError as e:
            stop_event.set()
            last_error = f"Connection error: {e}"
            print(f"    Connection error for {model_info['id']}: {e}")
            continue

        except Exception as e:
            stop_event.set()
            last_error = str(e)
            print(f"    Unexpected error from {model_info['id']}: {e}")
            # Don't retry on unexpected errors
            break

    # All retries exhausted
    return {
        "model": model_info["id"],
        "family": family,
        "role": model_info["role"],
        "status": "error",
        "findings_count": 0,
        "error": last_error or "Unknown error after all retries",
        "result": {"findings": [], "summary": {"overall_assessment": f"Failed: {last_error}"}},
    }


def run_audit(project_dir: str, output_file: str, models_override: str = None,
              parallel: bool = False, auto_confirm: bool = False) -> dict:
    """Run the full multi-model audit pipeline.

    Args:
        project_dir: Path to the project to audit.
        output_file: Path to write the JSON results.
        models_override: Comma-separated model IDs (overrides auto-discovery).
        parallel: If True, query all models simultaneously using ThreadPoolExecutor.
        auto_confirm: If True, skip confirmation prompt for large payloads.
    """
    api_key = os.environ.get("OPENROUTER_API_KEY", "")
    if not api_key:
        print("ERROR: OPENROUTER_API_KEY environment variable not set.")
        sys.exit(1)

    project_path = Path(project_dir).resolve()
    if not project_path.is_dir():
        print(f"ERROR: {project_dir} is not a directory.")
        sys.exit(1)

    # Step 1: Discover best models
    print("=" * 70)
    print("MULTI-MODEL CODE AUDITOR — ELITE EDITION")
    print("=" * 70)

    if models_override:
        models = []
        for m in models_override.split(","):
            m = m.strip()
            family = infer_family(m)
            role_map = {t["family"]: t["role"] for t in MODEL_TIERS}
            models.append({
                "id": m,
                "family": family,
                "role": role_map.get(family, "Custom Auditor"),
                "description": "User-specified",
            })
    else:
        print("\nDiscovering best available frontier models...")
        models = discover_best_models(api_key)

    print(f"\nSelected {len(models)} models:")
    for m in models:
        print(f"  [{m['family']:20s}] {m['id']:45s} — {m['role']}")

    # Step 2: Collect source files
    print(f"\nCollecting source files from {project_path}...")
    files = collect_source_files(project_path)
    print(f"  Found {len(files)} source files")

    payload = build_payload(files)
    char_count = len(payload)
    token_est = char_count // 4
    print(f"  Payload: {char_count:,} chars (~{token_est:,} tokens)")

    if char_count > WARN_PAYLOAD_CHARS:
        print(f"  NOTE: Large payload. Models with smaller context windows may truncate.")

    # Confirmation for large payloads (source code is being sent to a third-party API)
    if not auto_confirm and char_count > 100_000:
        print(f"\n  WARNING: {char_count:,} chars of source code will be sent to OpenRouter API.")
        try:
            answer = input("  Continue? [Y/n] ").strip().lower()
            if answer and answer not in ("y", "yes"):
                print("  Aborted.")
                sys.exit(0)
        except (EOFError, KeyboardInterrupt):
            # Non-interactive mode (piped stdin) — proceed automatically
            pass

    # Step 3: Query models
    results = []

    if parallel and len(models) > 1:
        print(f"\n  Querying {len(models)} models in PARALLEL...")
        with ThreadPoolExecutor(max_workers=len(models)) as executor:
            futures = {}
            for model in models:
                future = executor.submit(
                    query_model, model, payload, api_key,
                    show_progress=False,  # Avoid interleaved progress output
                )
                futures[future] = model

            for future in as_completed(futures):
                model = futures[future]
                try:
                    result = future.result()
                    print(f"  [{result['status']:12s}] {model['id']:45s} — {result['findings_count']} findings")
                    results.append(result)
                except Exception as e:
                    print(f"  [error       ] {model['id']:45s} — {e}")
                    results.append({
                        "model": model["id"],
                        "family": model["family"],
                        "role": model["role"],
                        "status": "error",
                        "findings_count": 0,
                        "error": str(e),
                        "result": {"findings": [], "summary": {"overall_assessment": f"Thread error: {e}"}},
                    })
    else:
        for i, model in enumerate(models, 1):
            print(f"\n{'─'*50}")
            print(f"[{i}/{len(models)}] {model['id']}")
            print(f"  Role: {model['role']}")
            start = time.time()
            result = query_model(model, payload, api_key)
            elapsed = time.time() - start
            print(f"  Status: {result['status']} | Findings: {result['findings_count']} | Time: {elapsed:.1f}s")
            results.append(result)

    # Step 4: Save results
    output = {
        "project": str(project_path),
        "project_name": project_path.name,
        "files_audited": len(files),
        "payload_chars": char_count,
        "payload_tokens_est": token_est,
        "models_queried": len(models),
        "models": [{"id": m["id"], "family": m["family"], "role": m["role"]} for m in models],
        "timestamp": time.strftime("%Y-%m-%d %H:%M:%S UTC", time.gmtime()),
        "results": results,
    }

    Path(output_file).write_text(json.dumps(output, indent=2))
    print(f"\n{'='*70}")
    print(f"Results saved to {output_file}")

    total = sum(r["findings_count"] for r in results)
    successful = sum(1 for r in results if r["status"] == "success")
    print(f"  {successful}/{len(models)} models responded successfully")
    print(f"  {total} total raw findings (before deduplication)")
    print(f"\nNext step: python3 cross_compare.py {output_file}")

    return output


# ─── CLI ──────────────────────────────────────────────────────────────────────

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Multi-Model Code Auditor — Elite Edition",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  python3 run_audit.py /path/to/project
  python3 run_audit.py /path/to/project --output my_audit.json
  python3 run_audit.py /path/to/project --parallel
  python3 run_audit.py /path/to/project --models anthropic/claude-sonnet-4,google/gemini-2.5-pro
  python3 run_audit.py --discover  # List available frontier models
        """,
    )
    parser.add_argument("project_dir", nargs="?", help="Path to the project directory to audit")
    parser.add_argument("--output", default="audit_results.json", help="Output JSON file")
    parser.add_argument("--models", default=None, help="Comma-separated model IDs (overrides auto-discovery)")
    parser.add_argument("--parallel", action="store_true", help="Query all models simultaneously")
    parser.add_argument("--discover", action="store_true", help="List available frontier models and exit")
    parser.add_argument("-y", "--yes", action="store_true", help="Skip confirmation prompt for large payloads")
    args = parser.parse_args()

    if args.discover:
        api_key = os.environ.get("OPENROUTER_API_KEY", "")
        models = discover_best_models(api_key)
        print("Best available frontier models:")
        for m in models:
            print(f"  [{m['family']:20s}] {m['id']:45s} — {m['role']}")
        sys.exit(0)

    if not args.project_dir:
        parser.error("project_dir is required unless --discover is used")

    run_audit(args.project_dir, args.output, args.models, args.parallel, args.yes)
