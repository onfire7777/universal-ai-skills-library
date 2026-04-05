# Prompt Engineering Techniques Used

This reference documents the elite prompt engineering techniques applied in the skill-debugger's dual-model prompts.

## Techniques by Category

### 1. Role Priming
- Claude: "senior software reliability engineer specializing in Python tooling and AI agent frameworks"
- Manus: "expert Manus skill analyst" — leverages platform-specific knowledge

### 2. Complementary Expertise Assignment
- Claude focuses on: deep code reasoning, security, race conditions, resource leaks
- Manus focuses on: structural integrity, integration quality, trigger accuracy
- This avoids redundant analysis and maximizes coverage

### 3. Chain-of-Thought Enforcement
- "For each finding, you MUST: 1) Quote exact code, 2) Explain WHY, 3) Describe failure scenario, 4) Provide exact fix"
- Forces structured reasoning instead of surface-level pattern matching

### 4. Metacognitive Verification
- "Self-verify: Would I bet $100 this is a real bug?"
- Reduces false positives by forcing the model to evaluate its own confidence

### 5. Negative Prompting (What NOT to Report)
- Explicit exclusion list: style preferences, missing type hints, theoretical vulnerabilities
- Dramatically reduces noise in findings

### 6. Severity Calibration with Concrete Examples
- Each severity level has specific criteria (not vague "important" vs "minor")
- CRITICAL = "Script crashes, data loss, security vulnerability with practical exploit"

### 7. Structured Output Schema
- Both models use identical JSON schema for easy consensus merging
- Required fields prevent incomplete findings

### 8. Temperature Control
- temperature=0.1 for deterministic, reproducible analysis
- High temperature would cause inconsistent findings across runs

### 9. Deep Mode Extension
- Optional additional analysis dimensions activated by --deep flag
- Keeps standard mode fast while allowing thorough analysis when needed

## Model-Specific Optimizations

### Claude Opus 4.6
- Stronger at: multi-file reasoning, security analysis, subtle logic bugs
- Prompt emphasizes: code-level analysis, security, robustness edge cases
- Uses: longer, more detailed system prompt (Claude handles long context well)

### Manus gpt-4.1-mini
- Stronger at: understanding Manus platform conventions, structural analysis
- Prompt emphasizes: SKILL.md quality, integration issues, trigger accuracy
- Uses: more focused prompt (smaller context window budget)
