# Foundation Makefile — build/test entrypoints for the canonical tree.
# Owner: FOUNDATION task. See STRUCTURE.md + docs/ARCHITECTURE-decoupling.md.
#
# NOTE: `make test` is RED at HEAD by design — 4 known pre-existing failures
# (2 real routing-drift, 2 macOS-only). "Green" for the refactor = no NEW
# failures beyond those 4. See STRUCTURE.md §4.

ROUTER_DIR := skill-router-cli

.PHONY: help build vet test baseline

help: ## Show available targets
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) \
		| awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-10s\033[0m %s\n",$$1,$$2}'

build: ## go build ./... (router)
	cd $(ROUTER_DIR) && go build ./...

vet: ## go vet ./... (router)
	cd $(ROUTER_DIR) && go vet ./...

test: ## go test ./... (router) — RED at HEAD: 4 known baseline failures
	cd $(ROUTER_DIR) && go test ./...

baseline: ## Record the baseline: build + vet + test (test failures are expected/known)
	@echo "== go build ==" ; cd $(ROUTER_DIR) && go build ./... && echo "BUILD_PASS"
	@echo "== go vet ==" ; cd $(ROUTER_DIR) && go vet ./... && echo "VET_PASS"
	@echo "== go test (expect 4 known failures; see STRUCTURE.md) ==" ; \
		cd $(ROUTER_DIR) && go test ./... ; echo "test rc=$$?"
