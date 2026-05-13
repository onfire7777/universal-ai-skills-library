## Summary

-

## Validation

- [ ] `git diff --check`
- [ ] `skill-router skills validate-manifest`
- [ ] `powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\validate-universal-ai-stack.ps1`
- [ ] `powershell -NoProfile -ExecutionPolicy Bypass -File .\ai-setup\scripts\public-release-audit.ps1`
- [ ] `go test ./...` from `skill-router-cli`

## Safety

- [ ] No secrets, OAuth/session files, logs, generated state, or local model files are committed.
- [ ] New commands are not redundant with existing `skill-router` commands.
- [ ] New adapters keep the full skill corpus router-first and on demand.
