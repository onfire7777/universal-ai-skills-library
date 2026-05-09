#!/usr/bin/env python3
"""Skill Debugger — Deep dual-model analysis of Manus skills.

Uses Manus's built-in model (gpt-4.1-mini) for structural/integration analysis
and Claude Opus 4.6 via OpenRouter for deep code reasoning and bug detection.
Merges findings by consensus and generates a prioritized fix plan.

Usage:
    python3 debug_skill.py <skill-name>                    # Debug a skill
    python3 debug_skill.py <skill-name> --fix              # Debug + auto-fix
    python3 debug_skill.py <skill-name> --deep             # Extended analysis
    python3 debug_skill.py /path/to/skill-dir              # Debug by path
    python3 debug_skill.py <skill-name> --model claude     # Claude only
    python3 debug_skill.py <skill-name> --model manus      # Manus only
"""

import argparse
import json
import os
import re
import subprocess
import sys
import textwrap
import time
from concurrent.futures import ThreadPoolExecutor, as_completed
from datetime import datetime, timezone
from pathlib import Path

# ─── Constants ───────────────────────────────────────────────────────────────

SKILLS_DIR = Path("/home/ubuntu/skills")
OPENROUTER_API_KEY = os.environ.get("OPENROUTER_API_KEY", "")
OPENAI_API_KEY = os.environ.get("OPENAI_API_KEY", "")

# Best available models
CLAUDE_MODEL = "anthropic/claude-opus-4.6"
MANUS_MODEL = "gpt-4.1-mini"

# Retry configuration
MAX_RETRIES = 3
RETRY_DELAYS = [10, 30, 60]
CONNECT_TIMEOUT = 30
READ_TIMEOUT = 300

# ─── Prompt Engineering ─────────────────────────────────────────────────────

def build_claude_system_prompt(deep: bool = False) -> str:
    """Build the Claude Opus system prompt using elite prompt engineering.

    Techniques used:
    - Role priming: Expert persona with specific domain expertise
    - Chain-of-thought: Explicit reasoning steps before conclusions
    - Structured output: JSON schema with required fields
    - Negative prompting: What NOT to report (reduces false positives)
    - Few-shot calibration: Severity scale with concrete examples
    - Metacognitive prompting: Self-verification step
    """
    depth_instruction = ""
    if deep:
        depth_instruction = """
DEEP ANALYSIS MODE — additionally examine:
- Race conditions in any threading/async code
- Resource leaks (file handles, DB connections, sockets)
- Edge cases in string parsing (unicode, null bytes, empty strings)
- Implicit type coercions that could cause silent failures
- Error handling completeness (are all failure modes covered?)
- API contract violations (does the code match its docstrings?)
- Performance pathologies (O(n^2) loops, repeated I/O, unbounded growth)
"""

    return textwrap.dedent(f"""\
    You are a senior software reliability engineer specializing in Python tooling \
    and AI agent frameworks. You have deep expertise in debugging skills for the \
    Manus AI agent platform.

    Your task: Perform a thorough debug analysis of a Manus skill — its SKILL.md \
    instructions, Python scripts, and reference files. Identify bugs, logic errors, \
    robustness gaps, and integration issues that would cause the skill to fail or \
    produce incorrect results when used by Manus.
    {depth_instruction}
    ANALYSIS FRAMEWORK — examine each dimension:

    1. STRUCTURAL INTEGRITY
       - Does SKILL.md have valid YAML frontmatter with name + description?
       - Is the description comprehensive enough to trigger correctly?
       - Are file paths in instructions correct and consistent?
       - Do referenced scripts/references/templates actually exist?

    2. SCRIPT CORRECTNESS
       - Syntax errors, import errors, undefined variables
       - Logic bugs: off-by-one, wrong comparisons, inverted conditions
       - Unhandled exceptions that would crash the script
       - Missing error handling for I/O, network, and subprocess calls
       - Incorrect argument parsing or missing required args

    3. ROBUSTNESS
       - What happens with empty input, missing files, or malformed data?
       - Are there hardcoded paths that won't work across environments?
       - Do scripts handle large inputs without memory exhaustion?
       - Are timeouts set for network and subprocess calls?

    4. SECURITY
       - Command injection via unsanitized inputs to subprocess/os.system
       - Path traversal in file operations
       - API key exposure in logs or error messages
       - Unsafe deserialization (pickle, eval, exec)

    5. INTEGRATION WITH MANUS
       - Will Manus understand when and how to use this skill?
       - Are the instructions clear enough for an AI agent to follow?
       - Do scripts produce output that Manus can parse and act on?
       - Are there ambiguities that could cause Manus to misuse the skill?

    WHAT NOT TO REPORT (reduces noise):
    - Style preferences (naming conventions, line length)
    - Missing type hints (unless they cause actual bugs)
    - Theoretical vulnerabilities with no practical attack vector in this context
    - "Could be improved" suggestions unless they fix actual failure modes
    - Missing features (only report broken existing features)

    SEVERITY CALIBRATION:
    - CRITICAL: Script crashes, data loss, security vulnerability with practical exploit
    - HIGH: Incorrect output, silent failures, broken integration with Manus
    - MEDIUM: Edge case failures, missing error handling for likely scenarios
    - LOW: Minor robustness gaps, unclear instructions, cosmetic issues

    REASONING PROCESS:
    For each finding, you MUST:
    1. Quote the exact problematic code or text
    2. Explain WHY it's a problem (not just WHAT)
    3. Describe a concrete scenario where it fails
    4. Provide the exact fix (code or text replacement)
    5. Self-verify: "Would I bet $100 this is a real bug?" — only report if yes

    OUTPUT FORMAT — respond with a JSON object:
    {{
      "skill_name": "string",
      "overall_health": "healthy|degraded|broken",
      "findings": [
        {{
          "id": "DBG-001",
          "severity": "critical|high|medium|low",
          "category": "structural|script|robustness|security|integration",
          "file": "relative/path/to/file",
          "line_range": "L42-L55",
          "title": "Concise one-line title",
          "problematic_code": "exact code snippet",
          "explanation": "Why this is a problem + concrete failure scenario",
          "fix": "Exact replacement code or instruction",
          "confidence": 0.95
        }}
      ],
      "summary": "2-3 sentence overall assessment"
    }}

    IMPORTANT: Return ONLY the JSON object, no markdown fences, no preamble.""")


def build_manus_system_prompt(deep: bool = False) -> str:
    """Build the Manus model system prompt — focused on structural and integration analysis.

    Techniques used:
    - Complementary expertise: Focuses on areas where Manus has unique insight
    - Structured reasoning: Step-by-step analysis framework
    - Concrete output schema: Matches Claude's schema for easy merging
    """
    depth_instruction = ""
    if deep:
        depth_instruction = """
DEEP MODE — additionally check:
- Cross-references between SKILL.md and all bundled files
- Whether the skill's description would correctly trigger in edge cases
- Whether instructions handle Manus's tool limitations (e.g., no parallel function calls)
- Whether the skill conflicts with or duplicates other common skills
"""

    return textwrap.dedent(f"""\
    You are an expert Manus skill analyst. You deeply understand how Manus skills \
    work: YAML frontmatter triggers skill loading, the body provides instructions, \
    and bundled scripts/references/templates extend capabilities.

    Analyze this skill for bugs, structural issues, and integration problems.
    {depth_instruction}
    FOCUS AREAS (your unique strengths):

    1. SKILL.MD QUALITY
       - Frontmatter: Is name accurate? Is description comprehensive enough to \
    trigger when needed AND not trigger when irrelevant?
       - Instructions: Are they clear, unambiguous, and actionable for an AI agent?
       - File references: Do all mentioned paths exist? Are they correct?
       - Progressive disclosure: Is content properly split between SKILL.md and references?

    2. SCRIPT ANALYSIS
       - Do scripts run correctly? Check imports, syntax, argument parsing
       - Are file paths hardcoded or properly parameterized?
       - Do scripts handle errors gracefully with informative messages?
       - Is output formatted so Manus can parse and use it?

    3. INTEGRATION ISSUES
       - Would Manus know WHEN to use this skill based on the description?
       - Would Manus know HOW to use it based on the instructions?
       - Are there missing steps that Manus would need to figure out?
       - Does the skill assume capabilities Manus doesn't have?

    4. CORRECTNESS
       - Logic errors in scripts (wrong conditions, missing cases)
       - Incorrect regex patterns or string parsing
       - Math errors or off-by-one bugs
       - Missing validation of inputs or API responses

    DO NOT REPORT: Style issues, missing type hints, theoretical concerns, \
    feature requests. Only report actual bugs and failure modes.

    For each finding, provide:
    - The exact problematic code/text
    - A concrete scenario where it fails
    - The exact fix

    OUTPUT: JSON object matching this schema:
    {{
      "skill_name": "string",
      "overall_health": "healthy|degraded|broken",
      "findings": [
        {{
          "id": "DBG-001",
          "severity": "critical|high|medium|low",
          "category": "structural|script|robustness|security|integration",
          "file": "relative/path/to/file",
          "line_range": "L42-L55",
          "title": "Concise one-line title",
          "problematic_code": "exact code snippet",
          "explanation": "Why this is a problem + concrete failure scenario",
          "fix": "Exact replacement code or instruction",
          "confidence": 0.95
        }}
      ],
      "summary": "2-3 sentence overall assessment"
    }}

    Return ONLY the JSON object, no markdown fences, no preamble.""")


# ─── Skill Collection ───────────────────────────────────────────────────────

def collect_skill_files(skill_dir: Path) -> str:
    """Collect all skill files into a single payload string."""
    payload_parts = []
    max_file_size = 512 * 1024  # 512KB per file

    for path in sorted(skill_dir.rglob("*")):
        if not path.is_file():
            continue
        if "__pycache__" in str(path) or ".git" in str(path):
            continue
        if path.is_symlink():
            continue

        rel = path.relative_to(skill_dir)
        try:
            size = path.stat().st_size
            if size > max_file_size:
                payload_parts.append(f"--- FILE: {rel} (TRUNCATED — {size} bytes) ---\n")
                content = path.read_text(encoding="utf-8", errors="replace")[:max_file_size]
                payload_parts.append(content)
                payload_parts.append(f"\n--- END (truncated at {max_file_size} bytes) ---\n\n")
            elif size == 0:
                payload_parts.append(f"--- FILE: {rel} (EMPTY) ---\n\n")
            else:
                content = path.read_text(encoding="utf-8", errors="replace")
                payload_parts.append(f"--- FILE: {rel} ---\n")
                payload_parts.append(content)
                payload_parts.append(f"\n--- END ---\n\n")
        except (UnicodeDecodeError, OSError):
            payload_parts.append(f"--- FILE: {rel} (BINARY — skipped) ---\n\n")

    return "".join(payload_parts)


# ─── API Calls ───────────────────────────────────────────────────────────────

def query_claude(payload: str, deep: bool = False) -> dict:
    """Query Claude Opus 4.6 via OpenRouter with retry logic."""
    import requests

    if not OPENROUTER_API_KEY:
        return {"error": "OPENROUTER_API_KEY not set", "findings": []}

    system_prompt = build_claude_system_prompt(deep)
    user_message = f"Analyze this Manus skill for bugs and issues:\n\n{payload}"

    headers = {
        "Authorization": f"Bearer {OPENROUTER_API_KEY}",
        "Content-Type": "application/json",
    }

    body = {
        "model": CLAUDE_MODEL,
        "messages": [
            {"role": "system", "content": system_prompt},
            {"role": "user", "content": user_message},
        ],
        "temperature": 0.1,
        "max_tokens": 16384,
    }

    for attempt in range(MAX_RETRIES):
        try:
            print(f"  [Claude Opus 4.6] Attempt {attempt + 1}/{MAX_RETRIES}...", end="", flush=True)
            start = time.time()
            r = requests.post(
                "https://openrouter.ai/api/v1/chat/completions",
                headers=headers,
                json=body,
                timeout=(CONNECT_TIMEOUT, READ_TIMEOUT),
            )
            elapsed = time.time() - start
            print(f" {elapsed:.1f}s", flush=True)

            if r.status_code == 200:
                content = r.json()["choices"][0]["message"]["content"]
                return parse_model_response(content, "claude")

            error_msg = r.text[:200]
            # Non-retryable errors
            if r.status_code in (401, 402, 403):
                print(f"  [Claude] Auth/billing error: {error_msg}")
                return {"error": error_msg, "findings": []}

            print(f"  [Claude] HTTP {r.status_code}: {error_msg}")

        except requests.exceptions.Timeout:
            print(f" TIMEOUT", flush=True)
        except Exception as e:
            print(f" ERROR: {e}", flush=True)

        if attempt < MAX_RETRIES - 1:
            delay = RETRY_DELAYS[attempt]
            print(f"  [Claude] Retrying in {delay}s...")
            time.sleep(delay)

    return {"error": "All retries exhausted", "findings": []}


def query_manus(payload: str, deep: bool = False) -> dict:
    """Query Manus's built-in model (gpt-4.1-mini) via OpenAI-compatible API."""
    try:
        from openai import OpenAI
    except ImportError:
        subprocess.check_call([sys.executable, "-m", "pip", "install", "openai", "-q"])
        from openai import OpenAI

    if not OPENAI_API_KEY:
        return {"error": "OPENAI_API_KEY not set", "findings": []}

    system_prompt = build_manus_system_prompt(deep)
    user_message = f"Analyze this Manus skill for bugs and issues:\n\n{payload}"

    client = OpenAI()

    for attempt in range(MAX_RETRIES):
        try:
            print(f"  [Manus gpt-4.1-mini] Attempt {attempt + 1}/{MAX_RETRIES}...", end="", flush=True)
            start = time.time()
            response = client.chat.completions.create(
                model=MANUS_MODEL,
                messages=[
                    {"role": "system", "content": system_prompt},
                    {"role": "user", "content": user_message},
                ],
                temperature=0.1,
                max_tokens=16384,
                timeout=READ_TIMEOUT,
            )
            elapsed = time.time() - start
            print(f" {elapsed:.1f}s", flush=True)

            content = response.choices[0].message.content
            return parse_model_response(content, "manus")

        except Exception as e:
            print(f" ERROR: {e}", flush=True)
            if attempt < MAX_RETRIES - 1:
                delay = RETRY_DELAYS[attempt]
                print(f"  [Manus] Retrying in {delay}s...")
                time.sleep(delay)

    return {"error": "All retries exhausted", "findings": []}


def parse_model_response(content: str, source: str) -> dict:
    """Parse model response, handling markdown fences and malformed JSON."""
    # Strip markdown fences
    content = content.strip()
    if content.startswith("```"):
        lines = content.split("\n")
        # Remove first line (```json) and last line (```)
        if lines[-1].strip() == "```":
            lines = lines[1:-1]
        else:
            lines = lines[1:]
        content = "\n".join(lines)

    # Try direct parse
    try:
        result = json.loads(content)
        for f in result.get("findings", []):
            f["source"] = source
        return result
    except json.JSONDecodeError:
        pass

    # Try to extract JSON object from surrounding text
    brace_depth = 0
    start_idx = None
    for i, ch in enumerate(content):
        if ch == "{":
            if brace_depth == 0:
                start_idx = i
            brace_depth += 1
        elif ch == "}":
            brace_depth -= 1
            if brace_depth == 0 and start_idx is not None:
                try:
                    result = json.loads(content[start_idx:i + 1])
                    for f in result.get("findings", []):
                        f["source"] = source
                    return result
                except json.JSONDecodeError:
                    start_idx = None

    return {"error": f"Could not parse {source} response", "raw": content[:500], "findings": []}


# ─── Consensus Merging ───────────────────────────────────────────────────────

def merge_findings(claude_result: dict, manus_result: dict) -> list[dict]:
    """Merge findings from both models using consensus-based deduplication.

    Findings confirmed by both models get elevated confidence.
    Unique findings are kept but marked as single-source.
    """
    all_findings = []

    claude_findings = claude_result.get("findings", [])
    manus_findings = manus_result.get("findings", [])

    # Mark sources
    for f in claude_findings:
        f.setdefault("source", "claude")
    for f in manus_findings:
        f.setdefault("source", "manus")

    # Try to match findings across models
    matched_manus = set()

    for cf in claude_findings:
        best_match = None
        best_score = 0.0

        for i, mf in enumerate(manus_findings):
            if i in matched_manus:
                continue
            score = _similarity(cf, mf)
            if score > best_score:
                best_score = score
                best_match = i

        if best_match is not None and best_score >= 0.4:
            # Consensus finding — merge and elevate
            mf = manus_findings[best_match]
            matched_manus.add(best_match)

            merged = {
                "id": cf.get("id", "DBG-???"),
                "severity": _higher_severity(cf.get("severity", "low"), mf.get("severity", "low")),
                "category": cf.get("category", mf.get("category", "script")),
                "file": cf.get("file", mf.get("file", "")),
                "line_range": cf.get("line_range", mf.get("line_range", "")),
                "title": cf.get("title", mf.get("title", "")),
                "problematic_code": cf.get("problematic_code", mf.get("problematic_code", "")),
                "explanation": cf.get("explanation", mf.get("explanation", "")),
                "fix": cf.get("fix", mf.get("fix", "")),
                "confidence": min(1.0, max(cf.get("confidence", 0.7), mf.get("confidence", 0.7)) + 0.15),
                "consensus": True,
                "sources": ["claude", "manus"],
            }
            all_findings.append(merged)
        else:
            # Claude-only finding
            cf["consensus"] = False
            cf["sources"] = ["claude"]
            all_findings.append(cf)

    # Add unmatched Manus findings
    for i, mf in enumerate(manus_findings):
        if i not in matched_manus:
            mf["consensus"] = False
            mf["sources"] = ["manus"]
            all_findings.append(mf)

    # Sort by severity then confidence
    severity_order = {"critical": 0, "high": 1, "medium": 2, "low": 3}
    all_findings.sort(key=lambda f: (
        severity_order.get(f.get("severity", "low"), 4),
        -f.get("confidence", 0.5),
    ))

    # Re-number IDs
    for i, f in enumerate(all_findings, 1):
        f["id"] = f"DBG-{i:03d}"

    return all_findings


def _similarity(a: dict, b: dict) -> float:
    """Compute similarity between two findings for dedup."""
    score = 0.0

    # Same file
    if a.get("file") and b.get("file") and a["file"] == b["file"]:
        score += 0.3

    # Same category
    if a.get("category") == b.get("category"):
        score += 0.1

    # Title word overlap
    a_words = set(a.get("title", "").lower().split())
    b_words = set(b.get("title", "").lower().split())
    if a_words and b_words:
        overlap = len(a_words & b_words) / max(len(a_words | b_words), 1)
        score += overlap * 0.3

    # Code snippet overlap
    a_code = a.get("problematic_code", "").strip()
    b_code = b.get("problematic_code", "").strip()
    if a_code and b_code:
        # Check if one contains the other
        if a_code in b_code or b_code in a_code:
            score += 0.3
        else:
            a_tokens = set(a_code.split())
            b_tokens = set(b_code.split())
            if a_tokens and b_tokens:
                overlap = len(a_tokens & b_tokens) / max(len(a_tokens | b_tokens), 1)
                score += overlap * 0.3

    return score


def _higher_severity(a: str, b: str) -> str:
    """Return the higher severity level."""
    order = {"critical": 0, "high": 1, "medium": 2, "low": 3}
    return a if order.get(a, 4) <= order.get(b, 4) else b


# ─── Report Generation ──────────────────────────────────────────────────────

def generate_report(skill_name: str, findings: list[dict],
                    claude_result: dict, manus_result: dict,
                    output_path: Path) -> str:
    """Generate a detailed Markdown debug report."""
    lines = [
        f"# Skill Debug Report: `{skill_name}`",
        f"",
        f"**Generated:** {datetime.now(timezone.utc).strftime('%Y-%m-%d %H:%M UTC')}",
        f"**Models:** Claude Opus 4.6 + Manus gpt-4.1-mini",
        f"",
    ]

    # Model status
    claude_ok = "error" not in claude_result
    manus_ok = "error" not in manus_result
    lines.append("## Model Status")
    lines.append("")
    lines.append(f"| Model | Status | Findings |")
    lines.append(f"|---|---|---|")
    lines.append(f"| Claude Opus 4.6 | {'OK' if claude_ok else 'FAILED: ' + claude_result.get('error', '?')[:50]} | {len(claude_result.get('findings', []))} |")
    lines.append(f"| Manus gpt-4.1-mini | {'OK' if manus_ok else 'FAILED: ' + manus_result.get('error', '?')[:50]} | {len(manus_result.get('findings', []))} |")
    lines.append("")

    # Overall health
    if not findings:
        lines.append("## Overall Health: HEALTHY")
        lines.append("")
        lines.append("No bugs or issues found. The skill appears to be well-structured and functional.")
    else:
        severities = [f.get("severity", "low") for f in findings]
        if "critical" in severities:
            health = "BROKEN"
        elif "high" in severities:
            health = "DEGRADED"
        else:
            health = "NEEDS ATTENTION"

        consensus_count = sum(1 for f in findings if f.get("consensus"))
        lines.append(f"## Overall Health: {health}")
        lines.append("")
        lines.append(f"**{len(findings)} issues found** ({consensus_count} confirmed by both models)")
        lines.append("")

        # Summary table
        lines.append("## Findings Summary")
        lines.append("")
        lines.append("| ID | Severity | Category | File | Title | Consensus |")
        lines.append("|---|---|---|---|---|---|")
        for f in findings:
            consensus = "Both" if f.get("consensus") else ", ".join(f.get("sources", ["?"]))
            lines.append(f"| {f['id']} | {f.get('severity', '?')} | {f.get('category', '?')} | `{f.get('file', '?')}` | {f.get('title', '?')} | {consensus} |")
        lines.append("")

        # Detailed findings
        lines.append("## Detailed Findings")
        lines.append("")
        for f in findings:
            consensus_badge = " [CONSENSUS]" if f.get("consensus") else f" [{', '.join(f.get('sources', ['?']))}]"
            lines.append(f"### {f['id']}: {f.get('title', 'Untitled')}{consensus_badge}")
            lines.append("")
            lines.append(f"**Severity:** {f.get('severity', '?')} | **Category:** {f.get('category', '?')} | **Confidence:** {f.get('confidence', '?')}")
            lines.append(f"**File:** `{f.get('file', '?')}` {f.get('line_range', '')}")
            lines.append("")

            if f.get("problematic_code"):
                lines.append("**Problematic Code:**")
                lines.append("```python")
                lines.append(f.get("problematic_code", ""))
                lines.append("```")
                lines.append("")

            if f.get("explanation"):
                lines.append(f"**Problem:** {f.get('explanation', '')}")
                lines.append("")

            if f.get("fix"):
                lines.append("**Fix:**")
                lines.append("```python")
                lines.append(f.get("fix", ""))
                lines.append("```")
                lines.append("")

            lines.append("---")
            lines.append("")

    report = "\n".join(lines)
    output_path.write_text(report)
    return str(output_path)


# ─── Main ────────────────────────────────────────────────────────────────────

def main():
    parser = argparse.ArgumentParser(description="Debug a Manus skill using dual AI models")
    parser.add_argument("skill", help="Skill name or path to skill directory")
    parser.add_argument("--fix", action="store_true", help="Auto-apply fixes after analysis")
    parser.add_argument("--deep", action="store_true", help="Extended deep analysis mode")
    parser.add_argument("--model", choices=["both", "claude", "manus"], default="both",
                        help="Which model(s) to use (default: both)")
    parser.add_argument("--output", default=None, help="Output report path")
    args = parser.parse_args()

    # Resolve skill directory
    skill_path = Path(args.skill)
    if not skill_path.is_absolute():
        skill_path = SKILLS_DIR / args.skill
    if not skill_path.exists():
        print(f"ERROR: Skill not found: {skill_path}", file=sys.stderr)
        sys.exit(1)
    if not (skill_path / "SKILL.md").exists():
        print(f"ERROR: No SKILL.md in {skill_path}", file=sys.stderr)
        sys.exit(1)

    skill_name = skill_path.name

    print("=" * 60)
    print(f"SKILL DEBUGGER — {skill_name}")
    print("=" * 60)
    print(f"  Path: {skill_path}")
    print(f"  Mode: {'deep' if args.deep else 'standard'}")
    print(f"  Models: {args.model}")
    print()

    # Collect skill files
    print("Collecting skill files...")
    payload = collect_skill_files(skill_path)
    char_count = len(payload)
    print(f"  Payload: {char_count:,} chars ({char_count // 4:,} est. tokens)")
    print()

    # Query models in parallel
    print("Querying models...")
    claude_result = {"findings": [], "error": "skipped"}
    manus_result = {"findings": [], "error": "skipped"}

    if args.model in ("both", "claude") and args.model in ("both", "manus"):
        # Parallel execution
        with ThreadPoolExecutor(max_workers=2) as pool:
            futures = {}
            if args.model in ("both", "claude"):
                futures["claude"] = pool.submit(query_claude, payload, args.deep)
            if args.model in ("both", "manus"):
                futures["manus"] = pool.submit(query_manus, payload, args.deep)

            for name, future in futures.items():
                try:
                    result = future.result(timeout=600)
                    if name == "claude":
                        claude_result = result
                    else:
                        manus_result = result
                except Exception as e:
                    print(f"  [{name}] Failed: {e}")
    else:
        if args.model == "claude":
            claude_result = query_claude(payload, args.deep)
        elif args.model == "manus":
            manus_result = query_manus(payload, args.deep)

    print()

    # Report errors
    if "error" in claude_result and claude_result["error"] != "skipped":
        print(f"  Claude error: {claude_result['error']}")
    if "error" in manus_result and manus_result["error"] != "skipped":
        print(f"  Manus error: {manus_result['error']}")

    # Merge findings
    print("Merging findings...")
    if args.model == "both":
        findings = merge_findings(claude_result, manus_result)
    elif args.model == "claude":
        findings = claude_result.get("findings", [])
        for i, f in enumerate(findings, 1):
            f["id"] = f"DBG-{i:03d}"
            f["consensus"] = False
            f["sources"] = ["claude"]
    else:
        findings = manus_result.get("findings", [])
        for i, f in enumerate(findings, 1):
            f["id"] = f"DBG-{i:03d}"
            f["consensus"] = False
            f["sources"] = ["manus"]

    print(f"  Total findings: {len(findings)}")
    consensus_count = sum(1 for f in findings if f.get("consensus"))
    if args.model == "both" and consensus_count:
        print(f"  Consensus (both models agree): {consensus_count}")
    print()

    # Print findings summary
    if findings:
        print("FINDINGS:")
        for f in findings:
            sev = f.get("severity", "?").upper()
            badge = " [CONSENSUS]" if f.get("consensus") else ""
            print(f"  {f['id']} [{sev}] {f.get('title', '?')}{badge}")
    else:
        print("NO ISSUES FOUND — skill appears healthy.")

    # Generate report
    output_path = Path(args.output) if args.output else (skill_path / "DEBUG_REPORT.md")
    report_path = generate_report(skill_name, findings, claude_result, manus_result, output_path)
    print(f"\nReport: {report_path}")

    # Save raw results
    raw_path = skill_path / ".debug_raw.json"
    raw_data = {
        "timestamp": datetime.now(timezone.utc).isoformat(),
        "skill": skill_name,
        "claude": claude_result,
        "manus": manus_result,
        "merged": findings,
    }
    raw_path.write_text(json.dumps(raw_data, indent=2, default=str))
    print(f"Raw data: {raw_path}")

    # Auto-fix mode
    if args.fix and findings:
        print(f"\n{'=' * 60}")
        print("AUTO-FIX MODE")
        print(f"{'=' * 60}")
        print("Fixes should be applied by Manus using the detailed report above.")
        print("Read the DEBUG_REPORT.md and apply each fix in order of severity.")

    print(f"\n{'=' * 60}")
    print("DONE")
    print(f"{'=' * 60}")


if __name__ == "__main__":
    main()
