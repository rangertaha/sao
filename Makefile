SHELL := /bin/sh

.PHONY: help fmt vet test test-race test-cover lint tidy verify ci ui-install ui-build run-server run-ui

help: ## Show available targets
	@awk 'BEGIN {FS = ":.*##"; printf "\nTargets:\n"} /^[a-zA-Z0-9_-]+:.*##/ {printf "  %-12s %s\n", $$1, $$2} END {printf "\n"}' $(MAKEFILE_LIST)

fmt: ## Format all Go packages
	go fmt ./...

vet: ## Run go vet static checks
	go vet ./...

test: ## Run unit tests
	go test ./...

test-race: ## Run tests with the race detector
	go test -race ./...

test-cover: ## Generate test coverage profile
	go test -coverprofile=coverage.out ./...

lint: ## Run golangci-lint checks
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint is not installed. Install from https://golangci-lint.run/"; \
		exit 1; \
	}
	golangci-lint run

tidy: ## Tidy go module dependencies
	go mod tidy

verify: fmt vet test ## Run local verification checks

ci: fmt vet test-race lint ## Run CI-equivalent checks locally

ui-install: ## Install UI dependencies
	cd web/ui && npm install

ui-build: ## Build React UI assets into web/ui/dist
	cd web/ui && npm run build

run-server: ## Run SAO server command
	go run . server

run-ui: ## Run embedded UI command
	go run . ui
