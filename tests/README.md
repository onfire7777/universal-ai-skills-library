# tests/

Repository-level test assets for the Universal AI Skills Library.

The router's own Go unit tests live next to the code in `skill-router-cli/`.
This directory holds **cross-cutting** tests that treat the router as a built
binary and the registry as data — so they survive the Goal 1 router decoupling
and the Goal 2/3 consolidation.

| Path | Purpose |
|------|---------|
| `characterization/` | Behavioural baseline pinning router + registry + legacy `manus` alias **before** the refactor, so regressions are caught. See its [README](characterization/README.md). |
| `routing-eval/` | Labeled routing eval set + harness reporting **P@1 / MRR / Recall@5** over the live corpus, with a committed baseline and a no-regression CI gate. Phase 0 of `docs/ARCHITECTURE_IMPROVEMENT_PLAN.md`. See its [README](routing-eval/README.md). |
| `fixtures/skills-lib/` | A **frozen** miniature skills library (manifest + `skills/`) the characterization tests run against, instead of the live `skills/` tree. Shared with `skill-router-cli`'s Go routing tests. |

## Running

```bash
# Fast Python characterization suite (skips the slow Go suite)
CHAR_SKIP_GO_TESTS=1 python3 -m unittest discover -s tests/characterization -p 'test_*.py'

# Full suite (also runs `go test ./...` and checks for NEW failures)
python3 -m unittest discover -s tests/characterization -p 'test_*.py'
```

Set `SKILL_ROUTER_BIN=/path/to/skill-router` to reuse a pre-built binary and
skip the per-run `go build` (CI does this).
