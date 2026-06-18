# Foundation Makefile — build/test entrypoints for the canonical tree.
# Owner: FOUNDATION task. See STRUCTURE.md + docs/ARCHITECTURE-decoupling.md.
#
# NOTE: `make test` is expected to pass at HEAD. Use `make eval-check` for
# routing-quality regression gating in addition to unit tests.

ROUTER_DIR := skill-router-cli

.PHONY: help build vet test baseline schema eval eval-check eval-baseline \
        registry-check registry-write parity dist release-dry install-local

help: ## Show available targets
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) \
		| awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-10s\033[0m %s\n",$$1,$$2}'

build: ## go build ./... (router)
	cd $(ROUTER_DIR) && go build ./...

vet: ## go vet ./... (router)
	cd $(ROUTER_DIR) && go vet ./...

test: ## go test ./... (router)
	cd $(ROUTER_DIR) && go test ./...

baseline: ## Record the baseline: build + vet + test
	@echo "== go build ==" ; cd $(ROUTER_DIR) && go build ./... && echo "BUILD_PASS"
	@echo "== go vet ==" ; cd $(ROUTER_DIR) && go vet ./... && echo "VET_PASS"
	@echo "== go test ==" ; \
		cd $(ROUTER_DIR) && go test ./... ; echo "test rc=$$?"

# --- Phase 0 (docs/ARCHITECTURE_IMPROVEMENT_PLAN.md) — measurable routing ---

schema: ## Validate SKILL.md frontmatter vs schemas/skill.schema.json (strict)
	node scripts/registry/validate-schema.mjs --error

eval: ## Run the routing eval (P@1 / MRR / Recall@5) over tests/routing-eval/cases.jsonl
	python3 tests/routing-eval/run_eval.py

eval-check: ## Routing metrics gate — fail if metrics regress vs the committed baseline
	python3 tests/routing-eval/run_eval.py --check

eval-baseline: ## Re-record the routing-eval baseline snapshot (only after an INTENDED change)
	python3 tests/routing-eval/run_eval.py --baseline

# --- Phase 0 (cont.) — Go owns the build; Node retired behind the parity gate ---

registry-check: ## Verify registry artifacts via the Go owner (skill-router registry build --check)
	cd $(ROUTER_DIR) && go run . registry build --check

registry-write: ## Regenerate the registry artifacts with the Go owner
	cd $(ROUTER_DIR) && go run . registry build --write

parity: ## Cut-over gate: prove Go registry build is byte-identical to the Node generator
	bash scripts/registry/parity-check.sh

# --- Phase 5 (cont.) — binary distribution / packaging ---

dist: ## Cross-compile a local snapshot of release binaries (requires goreleaser)
	goreleaser release --clean --snapshot --skip=publish

release-dry: ## Validate the goreleaser config (requires goreleaser)
	goreleaser check

install-local: ## Build skill-router and install compatibility aliases into ~/.local/bin
	cd $(ROUTER_DIR) && go build -o skill-router . \
		&& mkdir -p $$HOME/.local/bin \
		&& install -m 0755 skill-router $$HOME/.local/bin/skill-router \
		&& ln -sf skill-router $$HOME/.local/bin/manus \
		&& rm -f skill-router \
	&& echo "installed skill-router + compatibility alias -> $$HOME/.local/bin"
