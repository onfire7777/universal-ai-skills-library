# Universal AI Skills Codex Plugin Rule

Use `skill-router` as the router-first skill loader:

```bash
skill-router preflight --hook-event UserPromptSubmit --json "<latest user prompt>"
skill-router skill <name>
skill-router skill search <query> --limit 10
skill-router route --explain "<user prompt>"
skill-router doctor
```

Run preflight only for real user prompts. Do not route from lifecycle hooks, assistant output, tool output, status checks, or background jobs.

For a request with multiple deliverables, constraints, named skill passes or verification needs, identify each requirement, then call `workflow_plan` or `skill-router skills workflow plan`. Plan even when the preflight candidate covers only part of the task. Keep the user’s original task in the prompt. For each requirement, identify the specific object, action, platform, output format and verification details that matter. Express these as source facets, give each acceptance criterion a check ID and evidence kind, and link the observed result to that ID during review. For uncovered requirements or mismatched context, use `workflow_search` to inspect alternatives and pin the relevant file, lines and digest in a revised plan. Use `workflow_load` for each stage, read its full skills and supporting files, and carry forward decisions. Record each applied element with its load receipt, any path-only files consulted, and the observed evidence. For large or truncated loads, request the file manifest, then load one contribution at a time. Save the full CLI response to a task-local file and read it in ranges if needed. Give a reason for each gap. Replan when requirements change and use `workflow_review` before claiming coverage.

For a single requirement, load a clearly relevant hook candidate (the workflow_candidate field in fallback JSON) with `skill-router skill <name>`.

Keep preflight internal and quiet. Do not use host-native skill tools for router-selected universal skills. Do not preload or duplicate the full canonical skill corpus.
