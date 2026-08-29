SHELL := /bin/bash

.PHONY: help lint lint-go lint-ts test test-go test-ts build build-api build-console gen gen-sdk schema-check audit fmt fmt-go fmt-ts bootstrap clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

bootstrap: ## Install local tooling (go, node deps, pre-commit)
	go mod download -C apps/control-api
	npm --prefix apps/console-web ci

gen: gen-sdk ## Regenerate SDK types from schema

gen-sdk: ## Generate Go and TS SDK types from JSON Schema (source of truth: packages/*/schema)
	@echo "TODO(Sprint 2): schema-codegen wired to packages/policy-schema + packages/event-schema"

schema-check: ## Validate JSON Schema files and golden tests
	@echo "TODO(Sprint 2): schema validator"

fmt: fmt-go fmt-ts ## Format all languages

fmt-go:
	cd apps/control-api && gofmt -l . && gofmt -w .

fmt-ts:
	npm --prefix apps/console-web run format -- --write

lint: lint-go lint-ts ## Lint all languages

lint-go:
	cd apps/control-api && golangci-lint run

lint-ts:
	npm --prefix apps/console-web run lint

test: test-go test-ts ## Run all tests

test-go:
	cd apps/control-api && go test ./...

test-ts:
	npm --prefix apps/console-web run test

build: build-api build-console ## Build all artifacts

build-api:
	cd apps/control-api && go build ./...

build-console:
	npm --prefix apps/console-web run build

audit: ## Run dependency and secret scans
	cd apps/control-api && go list -m all | govulncheck ./...
	npm --prefix apps/console-web run audit

test-failover: ## Run gateway + API failure drills
	@echo "=== Gateway failure tests ==="
	bash tests/failover/gateway-failure-tests.sh
	@echo "=== API failure tests (requires compose stack) ==="
	@echo "Skipping: run manually with docker compose up -d first"
	@echo "  bash tests/failover/api-failure-tests.sh"

test-replay: ## Run WAF replay corpus (requires live waf-engine on :8082)
	bash tests/replay/replay.sh --mode detect
	bash tests/replay/replay.sh --mode block

clean:
	rm -rf apps/console-web/node_modules apps/control-api/bin gateways/rust-gateway/target
