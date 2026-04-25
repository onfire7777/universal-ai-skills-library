# Domain Patterns

Purpose: Use this file when the domain is known and you need default axes, common constraints, scoring hints, and downstream routing.

## Contents

- Test
- Load
- Deploy
- UX
- Risk
- Experiment
- Compatibility
- Custom

## Test

Typical axes:

| Axis | Common values | Priority |
|---|---|---|
| browser | Chrome, Firefox, Safari, Edge | high |
| os | Windows, macOS, Linux, iOS, Android | high |
| viewport | desktop, tablet, mobile | medium |
| auth_state | logged_in, anonymous, expired_session | high |
| data_state | empty, populated, edge_case | medium |
| network | wifi, cellular, slow_3g, offline | low |
| locale | ja, en, zh-TW, ko | low |

Common constraints:

```yaml
constraints:
  exclude:
    - {browser: Safari, os: Windows}
    - {browser: Chrome, os: iOS}
  conditional:
    - if: {network: offline}
      then: {auth_state: logged_in}
```

Suggested next agent: `Voyager`, `Radar`, or `Siege`.

## Load

Typical axes:

| Axis | Common values | Priority |
|---|---|---|
| concurrent_users | 10, 100, 500, 1000, 5000 | high |
| data_volume | small, medium, large | high |
| endpoint | /api/users, /api/search, /api/checkout | high |
| duration | 1min, 5min, 30min, 1hour | medium |
| ramp_pattern | constant, gradual, spike, wave | medium |
| region | ap-northeast-1, us-east-1, eu-west-1 | low |

Common constraints:

```yaml
constraints:
  conditional:
    - if: {concurrent_users: 5000}
      then: {duration: "1min"}
  exclude:
    - {data_volume: large, concurrent_users: 5000}
```

Suggested next agent: `Siege`, `Beacon`, or `Bolt`.

## Deploy

Typical axes:

| Axis | Common values | Priority |
|---|---|---|
| environment | dev, staging, production | high |
| region | ap-northeast-1, us-east-1, eu-west-1 | high |
| version | current, next, rollback | high |
| traffic_split | 0%, 1%, 10%, 50%, 100% | medium |
| rollout_strategy | blue-green, canary, rolling | medium |
| feature_flags | enabled, disabled | low |

Default ordering:

1. `dev -> staging -> production`
2. nearest region first
3. `1% -> 10% -> 50% -> 100%`

Suggested next agent: `Scaffold`, `Gear`, or `Beacon`.

## UX

Typical axes:

| Axis | Common values | Priority |
|---|---|---|
| persona | beginner, intermediate, expert, senior, accessibility_user | high |
| device | desktop, tablet, mobile | high |
| scenario | first_visit, return_visit, task_completion, error_recovery | high |
| locale | ja, en, zh, ko, ar | medium |
| accessibility | none, screen_reader, keyboard_only, high_contrast | medium |
| connection | fast, slow, offline | low |

Suggested next agent: `Cast`, `Echo`, or `Researcher`.

## Risk

Typical axes:

| Axis | Common values | Priority |
|---|---|---|
| threat | XSS, SQLi, CSRF, SSRF, XXE, RCE, PathTraversal | high |
| attack_surface | Web_UI, REST_API, GraphQL, WebSocket, File_Upload | high |
| auth_level | anonymous, authenticated, privileged, admin | high |
| data_sensitivity | public, internal, confidential, restricted, PII | high |
| impact | low, medium, high, critical | medium |

Default risk tiers:

| Score | Priority |
|---|---|
| `7.0-10.0` | `P0 / Critical` |
| `4.0-6.9` | `P1 / High` |
| `2.0-3.9` | `P2 / Medium` |
| `0.0-1.9` | `P3 / Low` |

Suggested next agent: `Triage`, `Sentinel`, `Probe`, or `Scout`.

## Experiment

Typical axes:

| Axis | Common values | Priority |
|---|---|---|
| variable | button_color, cta_text, layout, price_display | high |
| user_segment | new_users, returning_users, premium, free | high |
| exposure_rate | 1%, 5%, 10%, 50% | medium |
| duration | 1week, 2weeks, 1month | medium |
| metric | CTR, CVR, retention, ARPU | high |

Suggested next agent: `Experiment` or `Pulse`.

## Compatibility

Typical axes:

| Axis | Common values | Priority |
|---|---|---|
| runtime_version | Node.js 18, 20, 22 / Python 3.9, 3.10, 3.12 | high |
| dependency_version | react@17, react@18, react@19 | high |
| os | ubuntu-22.04, ubuntu-24.04, macos-13, macos-14 | medium |
| architecture | x86_64, arm64 | medium |
| feature | core, experimental, deprecated | medium |

Suggested next agent: `Horizon` or `Builder`.

## Custom

Use `custom` when none of the built-in domains fit. In that case:

- preserve the generic matrix model
- do not invent a downstream specialist
- ask only if domain ambiguity changes the outcome materially
