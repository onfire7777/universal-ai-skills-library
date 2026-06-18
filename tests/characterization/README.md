# Characterization tests

These tests pin the **current behaviour** of the skill-router and the registry
*before* the Goal 1–3 refactor (router decoupling, repository consolidation, registry
unification), so any unintended behaviour change is caught. They are the
regression safety net for a "refactor / consolidate without breaking anything"
effort.

## What is pinned

| File | Pins |
|------|------|
| `test_router_characterization.py` | The router's `decision` + chosen skill (`best.name`, `best.source`) for a battery of prompts, and that `skill <name>` still loads a known skill. |
| `test_registry_characterization.py` | Fixture `validate-manifest` is clean; **live** registry keeps structural integrity (no dup names/dirs, no missing SKILL.md), never loses or relocates an existing skill, and only ever grows. |
| `test_legacy_alias_characterization.py` | Compatibility binary aliases still resolve skills and validate the manifest identically, and the manifest records them in generic compatibility fields. |
| `test_go_unit_baseline.py` | `go test ./...` introduces **no new** failures (subset of `baseline/go_known_failures.json`, now empty = stay green). |

## The fixture (why behaviour is pinned against it, not `skills/`)

Routing/lookup/alias assertions run against a **frozen** fixture skills library,
not the live `skills/` tree, so that Goal 2/3 edits to `skills/` and
`manifest.json` cannot break them. The fixture is the shared, routing-verified
corpus authored by Builder 2's Go routing tests at:

```
skill-router-cli/cmd/skills/testdata/route-fixture/
```

The router is pointed at it via `SKILL_ROUTER_REPO_DIR` (with `HOME`/config
redirected to throwaway dirs so no live external skill roots leak in). Override
the location with `CHAR_FIXTURE_DIR=/abs/path`.

## Baselines

* `baseline/routing_golden.json` — pinned router outcomes (decision + skill;
  scores are diagnostic only and **not** asserted, since scoring internals may
  drift).
* `baseline/registry_baseline.json` — fingerprint of the live 1,813-skill corpus.
  The slug set and its SHA-256 were **cross-validated** against Scout 1's
  independent baseline (`~/.bridgespace-scout2-work/baseline_uasl_slugs.txt`,
  sha256 `f1a5c648…`) — identical sets, 1,813 registered slugs. The gate counts
  *registered* slugs (manifest entries), deliberately excluding the one nested
  vendored `SKILL.md` under `infographic/` that a raw file scan over-counts.
* `baseline/go_known_failures.json` — hand-curated list of accepted Go failures
  (empty; the 4 pre-refactor reds were fixed by Builder 2).

## Running

```bash
# fast (skip the slow Go suite)
CHAR_SKIP_GO_TESTS=1 python3 -m unittest discover -s tests/characterization -p 'test_*.py'

# full
python3 -m unittest discover -s tests/characterization -p 'test_*.py'

# reuse a prebuilt router (skip per-run `go build`)
SKILL_ROUTER_BIN=/path/to/skill-router python3 -m unittest discover -s tests/characterization -p 'test_*.py'
```

## Updating a baseline (deliberate, reviewed changes only)

```bash
python3 tests/characterization/update_baseline.py   # refreshes routing + registry baselines
```

`go_known_failures.json` is intentionally **not** auto-generated — edit it by
hand so an accepted Go failure is always a conscious, reviewed decision.
