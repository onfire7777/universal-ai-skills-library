# Testing & CI

UASL is gated by five CI workflows and three complementary test surfaces (Go unit, routing-eval, characterization), plus a byte-parity oracle. Together they ensure routing quality never regresses and the registry stays deterministic.

Source: [`.github/workflows/`](https://github.com/onfire7777/universal-ai-skills-library/tree/main/.github/workflows) · [`tests/`](https://github.com/onfire7777/universal-ai-skills-library/tree/main/tests) · [`Makefile`](https://github.com/onfire7777/universal-ai-skills-library/blob/main/Makefile)

---

## CI workflows

| Workflow | Trigger | Platform | Enforces |
|---|---|---|---|
| **ci.yml** | push to `main`, all PRs | Windows | Go tests + build, stack-template validation, public-release audit |
| **security.yml** | push to `main`, all PRs | Ubuntu | Gitleaks `v8.30.1` secret scan |
| **registry-parity.yml** | changes to corpus/router/generator/artifacts | Ubuntu | Go == committed artifacts **and** Go == Node, byte-for-byte |
| **characterization.yml** | push to `main`, all PRs | Ubuntu | Routing-eval regression + behavioral baseline |
| **release.yml** | tags `v*`, manual dispatch | Ubuntu | goreleaser + cosign signing, parity gate |

### ci.yml (Windows)
Sets up Go from `skill-router-cli/go.mod`, runs `go test ./...`, builds `skill-router.exe`, then runs the PowerShell validators `validate-universal-ai-stack.ps1` and `public-release-audit.ps1 -SkipGoTests`. All steps must pass to merge.

### security.yml (Gitleaks)
Installs gitleaks `v8.30.1` and scans the **working tree** (`--no-git`) with `.gitleaks.toml`:
```bash
gitleaks detect --no-git --config .gitleaks.toml --redact --verbose --exit-code 1
```
Real secrets **fail** CI. The config allowlists ~45 known-illustrative findings (example JWTs, RFC samples, security-skill detection patterns, vendored third-party references) — each documented as a non-secret with the narrowest possible scope. Teaching examples in history are exempted via `--no-git`; real commits are still checked.

### registry-parity.yml
Runs on changes to `skills/**`, `scripts/registry/**`, `skill-router-cli/**`, `manifest.json`, or `docs/build_manifest.json`:
```bash
cd skill-router-cli && go run . registry build --check     # Go matches committed artifacts
bash scripts/registry/parity-check.sh                       # Go output == Node output, byte-identical
cd skill-router-cli && go test ./internal/registry/...
```
This is the safety net for the [Node → Go migration](Node-to-Go-Migration.md): Go is authoritative, Node is the non-blocking oracle.

### characterization.yml
Runs the routing-eval regression gate and the behavioral baseline on Ubuntu:
```bash
python3 tests/routing-eval/run_eval.py --check
```
It gates on **no NEW Go failures** (tracked against `tests/characterization/baseline/go_known_failures.json`), complementing the all-must-pass Windows CI.

### release.yml
On a `v*` tag: re-runs the parity gate (refuses to release on divergence), installs cosign, then runs **goreleaser** to build cross-platform binaries, create the GitHub Release, and **cosign-sign** the checksums (keyless OIDC). Optionally publishes a Homebrew tap when `HOMEBREW_TAP_GITHUB_TOKEN` is set.

---

## Test surfaces

### Go unit tests
```bash
cd skill-router-cli
go test ./...
go vet ./...
```
Cover the router CLI, the config-driven resolver, manifest parsing, and the registry builder (`internal/registry`).

### Routing-eval harness
Lives in [`tests/routing-eval/`](https://github.com/onfire7777/universal-ai-skills-library/tree/main/tests/routing-eval):

| File | Role |
|---|---|
| `cases.jsonl` | Curated, reviewed eval set — 55 cases (47 positive, 8 negative). The CI gate runs here. |
| `run_eval.py` | Routes each case; computes P@1, MRR, Recall@5, abstention accuracy. |
| `baseline/metrics.json` | Committed baseline; `--check` fails on regression. |
| `seed_corpus.py` | Generates a broad regression net (2k+ cases) from corpus descriptions + author `triggers:` (output is gitignored). |

Case format:
```json
{"id": "para-pytest", "prompt": "write pytest fixtures …", "expect": ["python-testing-patterns"], "reject": [], "tags": ["paraphrase"]}
```
- `expect` non-empty → positive case (any listed id is an acceptable top route).
- `expect` empty → abstention case (router should refuse to route).
- `reject` → explicitly-bad destinations.

```bash
make eval          # human report over cases.jsonl
make eval-check    # CI gate: fail if metrics regress vs baseline
make eval-baseline # re-record baseline (only after an intended change)
python3 tests/routing-eval/run_eval.py --show-misses
```

### Characterization tests
[`tests/characterization/`](https://github.com/onfire7777/universal-ai-skills-library/tree/main/tests/characterization) — 13+ Python tests pinning behavioral contracts: router, registry, legacy aliases, provider boundary, installer contracts, MCP bridge naming, model/platform compatibility, migration end-state, and more. Baselines (`routing_golden.json`, `registry_baseline.json`, `go_known_failures.json`) record ground truth. The suite runs on Linux and gates on no *new* Go regressions.

---

## Makefile cheat sheet

```bash
make build            # go build ./...
make vet              # go vet ./...
make test             # go test ./...
make baseline         # build + vet + test
make schema           # validate SKILL.md frontmatter (strict)
make registry-check   # verify manifest matches committed artifacts
make registry-write   # regenerate registry (Go owner)
make parity           # byte-parity oracle: Go == Node
make eval             # routing-eval human report
make eval-check       # routing-eval regression gate
make eval-baseline    # re-record eval baseline
make install-local    # install skill-router locally
make release-dry      # dry-run goreleaser
```

---

## Related pages

- [Manifest & Registry](Manifest-and-Registry.md) — the artifacts parity protects
- [Performance & Benchmarks](Performance-and-Benchmarks.md) — the numbers eval guards
- [Contributing](Contributing.md) — the pre-PR validation checklist
- [Security](Security.md) — the gitleaks policy in full
