# Failure Taxonomy

Purpose: load this when classifying loop anomalies. It maps observed evidence to `CONTRACT_MISSING`, `STATE_DRIFT`, `VERIFY_GAP`, `COMMIT_SCOPE_RISK`, or `TOOL_FAILURE`, then defines severity and the smallest reversible recovery step.

## Contents

1. Class definitions
2. Severity priority
3. Reporting schema
4. Contract-check mapping
5. Multiple-failure rule

## Class Definitions

### `CONTRACT_MISSING`

Missing or invalid loop artifacts prevent deterministic execution.

| Signal | First response | Recovery |
|--------|----------------|----------|
| `goal.md` missing | classify `P1`, halt execution | `bash bootstrap.sh ${LOOP_DIR}` |
| AC section empty | classify `P1`, trigger `ON_GOAL_CONTRACT_WEAK` | rewrite measurable ACs and verify commands |
| `progress.md` missing or empty | classify `P1`, initialize iteration `0` | append `## Iteration 0 — Bootstrap` |
| `state.env` missing | classify `P1`, rebuild defaults | `bash recover.sh` or create `NEXT_ITERATION=1`, `LAST_STATUS=READY` |
| footer malformed | classify `P1`, re-apply | emit valid `NEXUS_LOOP_STATUS` and `NEXUS_LOOP_SUMMARY` |
| partial template only | classify `P1`, preserve what exists | rebuild missing artifacts only |
| AC not measurable | classify `P1` | apply `vague-goal-handling.md` |

### `STATE_DRIFT`

Iteration counters or statuses disagree across artifacts.

| Signal | First response | Recovery |
|--------|----------------|----------|
| `NEXT_ITERATION` too high | re-sync from `progress.md` | `bash recover.sh` |
| `NEXT_ITERATION` too low | re-sync from `progress.md` | `bash recover.sh` |
| `LAST_STATUS` contradicts `progress.md` | update from evidence | `bash recover.sh` |
| `state.env.sha256` mismatch | treat as corrupted state | auto-run `recover.sh` |
| drift across `2+` consecutive iterations | classify `P0`, pause loop | manual investigation |

### `VERIFY_GAP`

`DONE` is claimed without successful verification evidence.

| Signal | First response | Recovery |
|--------|----------------|----------|
| `verify.sh` missing | classify `P0`, downgrade to `CONTINUE` | generate `verify.sh` |
| `verify.sh` exits non-zero | classify `P0`, downgrade to `CONTINUE` | rerun `bash verify.sh`, then fix failures |
| `done.md` exists but no verify evidence | classify `P0`, downgrade to `CONTINUE` | run verification independently |
| partial verification pass | classify `P1`, log partial state | continue with failing ACs as focus |
| `VERIFY_RESULT=SKIP` with `done.md` present | classify `P1`, warn | generate verify coverage or re-evaluate the claim |

### `COMMIT_SCOPE_RISK`

Auto-commit may include unrelated changes.

| Signal | First response | Recovery |
|--------|----------------|----------|
| `dirty-start-paths.txt` missing while `AUTOCOMMIT=true` | classify `P1`, regenerate baseline | resnapshot dirty baseline |
| staged files include baseline paths | classify `P0`, unstage baseline paths | `git reset HEAD -- <baseline_paths>` |
| single-loop candidate scope exceeds goal | classify `P1`, restrict staging | verify each candidate path |
| cross-loop path overlap | classify `P1`, suspend affected loop | delegate via `ORBIT_TO_GUARDIAN_HANDOFF` |
| `git add -A` used with dirty baseline | classify `P0`, revert commit | reset, fix baseline, recommit with scoped staging |

### `TOOL_FAILURE`

The runner or executor cannot make progress.

| Signal | First response | Recovery |
|--------|----------------|----------|
| network timeout / DNS failure | retry with backoff | network-focused retry |
| OOM or disk full | classify `P1`, pause | free resources, recover, resume |
| permission or token error | classify `P1`, pause | refresh credentials, resume |
| `EXEC_TIMEOUT` exceeded | retry once, then `P1` | increase timeout or enable adaptive timeout |
| deterministic non-zero exit | apply retry policy once | escalate to Builder if repeatable |
| `RETRY_LIMIT` exhausted | classify `P1`, record `BLOCKED` | investigate via `ORBIT_TO_SCOUT_HANDOFF` |
| circuit breaker `OPEN` | classify `P1`, record `BLOCKED` | cooldown or manual reset via `recover.sh --reset-circuit` |

### Retry Policy by Error Type

Errors are classified into two categories for retry policy:

| Category | Examples | Retry behavior | Backoff |
|----------|----------|----------------|---------|
| `TRANSIENT` | network timeout, DNS failure, API rate limit, OOM (after cleanup), `EXEC_TIMEOUT` exceeded | retry up to `RETRY_LIMIT` | exponential with jitter: `base^attempt + random(0, base^attempt * 0.5)` |
| `PERMANENT` | permission/token error, invalid input, deterministic non-zero exit, schema violation | no retry; escalate immediately | — |

Detection heuristics:

| Signal | Category | Rationale |
|--------|----------|-----------|
| exit code `124` (timeout) | `TRANSIENT` | timeout is environment-dependent |
| exit code `1` on consecutive identical input | `PERMANENT` | deterministic failure |
| exit code `137` (OOM killed) | `TRANSIENT` | memory pressure may be temporary |
| HTTP `429`, `502`, `503` in stderr | `TRANSIENT` | server-side rate limit or outage |
| HTTP `401`, `403` in stderr | `PERMANENT` | credentials or permissions issue |
| identical error message `3` consecutive times | `PERMANENT` | upgrade from transient to permanent |

When a `TRANSIENT` error is retried `3` consecutive times with the same error signature, it is reclassified as `PERMANENT` and the circuit breaker is engaged (see SKILL.md § Circuit Breaker).

## Severity Priority

| Severity | Classes | Action |
|----------|---------|--------|
| `P0` | false `DONE`, destructive commit-scope risk, escalated drift | stop and confirm |
| `P1` | drift, contract gaps, retry exhaustion | recover then continue |
| `P2` | intermittent tool failure within retry budget | retry policy |

## Reporting Schema

```yaml
ORBIT_FAILURE_REPORT:
  class: "<taxonomy class>"
  severity: "P0|P1|P2"
  sub_class: "<specific signal>"
  retry_category: "TRANSIENT|PERMANENT"
  circuit_state: "CLOSED|HALF_OPEN|OPEN"
  evidence:
    - "<file:line>"
  impact: "<what breaks>"
  recommendation: "<smallest reversible action>"
  recovery_cmd: "<command>"
```

## Contract-Check Mapping

| Contract score | Finding | Failure class | Severity | Next action |
|----------------|---------|---------------|----------|-------------|
| Missing | `goal.md` missing | `CONTRACT_MISSING` | `P1` | rebuild in bootstrap mode |
| Missing | `progress.md` empty | `CONTRACT_MISSING` | `P1` | initialize iteration `0` |
| Missing | `state.env` missing | `CONTRACT_MISSING` | `P1` | generate defaults |
| Partial | footer missing or malformed | `CONTRACT_MISSING` | `P1` | repair footer |
| Partial | state and timeline mismatch | `STATE_DRIFT` | `P1` | resync from progress evidence |
| Partial | `done.md` present and verify evidence unclear | `VERIFY_GAP` | `P0` | downgrade to `CONTINUE` |
| Partial | candidate scope may include baseline paths | `COMMIT_SCOPE_RISK` | `P0` | hand off to Guardian |
| Partial | `runner.log` contains failure entries | `TOOL_FAILURE` | `P2` | retry policy and environment check |
| Complete | retry budget exhausted | `TOOL_FAILURE` | `P1` | investigate via Scout |
| Complete | artifacts consistent | none | — | continue safely |

## Rules for Multiple Failure Classes

- Process the highest severity first.
- `P0` always preempts `P1` and `P2`.
- Use one handoff at a time.
- Record remaining issues as pending issues for the next audit cycle.

## Circuit Breaker Integration

When the same failure signature (class + sub_class + error message hash) occurs `3` times:

| State | Condition | Behavior |
|-------|-----------|----------|
| `CLOSED` | `< 3` consecutive same-signature failures | normal retry policy applies |
| `HALF_OPEN` | exactly `3` consecutive same-signature failures | allow one probe execution; if it fails, transition to `OPEN` |
| `OPEN` | `> 3` consecutive same-signature failures or probe failed | stop execution, emit `BLOCKED` status, require manual reset or cooldown |

Cooldown: `OPEN` state auto-transitions to `HALF_OPEN` after `CIRCUIT_COOLDOWN` seconds (default: `300`).

Reset: manual reset via `recover.sh --reset-circuit` or by clearing `${LOOP_DIR}/.circuit-state`.

Reporting: circuit breaker state changes are logged in the structured log (see `script-template-runner.md § Structured Logging`).
