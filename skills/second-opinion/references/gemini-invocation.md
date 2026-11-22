# Gemini CLI Invocation

Use the user’s configured Gemini CLI account and preserve any model they requested. The default model is `gemini-3.1-pro-preview`.

## Review command

Prepare the same scoped diff, context, and review focus used for Codex. Include untracked files for an uncommitted review. Pass that prompt through stdin in headless mode.

```bash
gemini --model gemini-3.1-pro-preview \
  --approval-mode default \
  --output-format text \
  --prompt "Review the supplied changes and return findings." \
  < "$prompt_file" > "$output_file" 2> "$stderr_log"
```

The prompt supports general, security, performance, and error-handling reviews with the same command. Keep Gemini in default approval mode. Report any approval or workspace-trust requirement and let the user decide through their normal setup. Treat a nonzero exit, empty output, or denied inspection as an incomplete review and preserve any partial findings.

| Failure | Response |
|---|---|
| Gemini CLI missing | Report that the executable is missing and provide its official installation documentation |
| Workspace trust or approval required | Report the blocked action and the diagnostic |
| Authentication, quota, model, or timeout error | Report the failure and any partial results |
| Account directs the user to another CLI | Report the account limitation and ask before changing CLI |

The invocation does not install extensions or change workspace trust. Existing extensions and tool policies follow the user’s Gemini CLI configuration.

[Gemini CLI headless mode](https://geminicli.com/docs/cli/headless/) · [Upstream second-opinion reference](https://github.com/trailofbits/skills/blob/f44eeba661acb45b0efecdc89a8aba7f0b7835ff/plugins/second-opinion/skills/second-opinion/references/gemini-invocation.md)
