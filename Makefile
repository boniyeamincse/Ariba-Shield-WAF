SHELL := /bin/bash

.PHONY: help lint lint-go lint-ts test test-go test-go-all test-ts build build-api build-services build-console gen gen-sdk gen-check schema-check audit fmt fmt-go fmt-ts bootstrap clean check-i18n

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

bootstrap: ## Install local tooling (go, node deps, pre-commit)
	go mod download -C apps/control-api
	go mod download -C services/waf-engine
	go mod download -C services/policy-compiler
	go mod download -C services/event-ingestor
	npm --prefix apps/console-web ci

gen: gen-sdk ## Regenerate SDK types from schema

gen-sdk: ## Generate Go and TS SDK types from JSON Schema (source of truth: packages/*/schema)
	go run ./tools/schema-gen

gen-check: ## Verify generated SDK matches schema (golden test; fails if drift)
	@tmp=$$(mktemp -d); \
	go run ./tools/schema-gen -out-dir "$$tmp" > /dev/null; \
	diff "$$tmp/sdk-go/sdk.go" packages/sdk-go/sdk.go && \
	diff "$$tmp/sdk-typescript/src/types.ts" packages/sdk-typescript/src/types.ts && \
	echo "gen: SDK up to date" || { echo "gen: SDK out of date — run 'make gen'"; exit 1; }

template-check: ## Verify the nginx template copies are identical (P2.17)
	@diff services/policy-compiler/internal/compiler/templates/shield.conf.tmpl \
	  gateways/openresty-gateway/nginx/templates/shield.conf.tmpl \
	  && echo "template: copies in sync" \
	  || { echo "template: DRIFT — run 'make template-sync'"; exit 1; }

template-sync: ## Sync the gateway template copy from the compiler source of truth
	cp services/policy-compiler/internal/compiler/templates/shield.conf.tmpl \
	  gateways/openresty-gateway/nginx/templates/shield.conf.tmpl

schema-check: ## Validate JSON Schema files and golden tests
	python3 -m json.tool packages/policy-schema/schema/policy-v0.json > /dev/null
	python3 -m json.tool packages/event-schema/schema/event-v0.json > /dev/null
	@echo "schema JSON valid"
	cd packages/sdk-typescript && npx tsc --noEmit
	@echo "generated TS SDK typechecks"

fmt: fmt-go fmt-ts ## Format all languages

fmt-go:
	@for d in apps/control-api services/waf-engine services/policy-compiler services/event-ingestor; do \
		(cd $$d && gofmt -l . && gofmt -w .); \
	done

fmt-ts:
	npm --prefix apps/console-web run format -- --write

lint: lint-go lint-ts ## Lint all languages

lint-go:
	cd apps/control-api && golangci-lint run
	cd services/waf-engine && golangci-lint run
	cd services/policy-compiler && golangci-lint run
	cd services/event-ingestor && golangci-lint run

lint-ts:
	npm --prefix apps/console-web run lint

test: test-go-all test-ts ## Run all tests

# P1.8: run tests for ALL Go services, not just control-api.
test-go-all: test-go test-waf test-compiler test-ingestor

test-go:
	cd apps/control-api && go test ./...

test-waf:
	cd services/waf-engine && go test ./...

test-compiler:
	cd services/policy-compiler && go test ./...

test-ingestor:
	cd services/event-ingestor && go test ./...

test-ts:
	npm --prefix apps/console-web run test

build: build-api build-services build-console ## Build all artifacts

build-api:
	cd apps/control-api && go build ./...

build-services: ## Build all Go services
	cd services/waf-engine && go build ./...
	cd services/policy-compiler && go build ./...
	cd services/event-ingestor && go build ./...

build-console:
	npm --prefix apps/console-web run build

check-i18n: ## Verify en/bn message catalogs have matching keys
	python3 tools/check-i18n.py

audit: ## Run dependency and secret scans
	cd apps/control-api && go list -m all | govulncheck ./...
	cd services/waf-engine && go list -m all | govulncheck ./...
	npm --prefix apps/console-web run audit

test-failover: ## Run gateway + API failure drills
	@echo "=== Gateway failure tests ==="
	bash tests/failover/gateway-failure-tests.sh
	@echo "=== API failure tests (requires compose stack up) ==="
	@if [ "$${RUN_API_TESTS:-0}" = "1" ]; then \
		echo "RUN_API_TESTS=1 — running API failure tests"; \
		bash tests/failover/api-failure-tests.sh; \
	else \
		echo "SKIP: set RUN_API_TESTS=1 to run API failure tests (needs: docker compose up -d)"; \
	fi

test-replay: ## Run WAF replay corpus (requires live waf-engine on :8082)
	bash tests/replay/replay.sh --mode detect
	bash tests/replay/replay.sh --mode block

clean:
	rm -rf apps/console-web/node_modules apps/control-api/bin gateways/rust-gateway/target
