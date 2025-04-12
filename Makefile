# Get the module path from go.mod
MODULE := $(shell go list -m)

# Define standard Go targets
.PHONY: test lint fmt tidy clean bench help

test: ## Run tests with coverage
	@echo "==> Running tests..."
	@go test ./... -v -coverprofile=coverage.out -coverpkg=$(MODULE)

lint: ## Run linter (requires golangci-lint)
	@echo "==> Running linter..."
	@golangci-lint run ./... || echo "Linting failed or golangci-lint not found. Please install it: https://golangci-lint.run/usage/install/"

fmt: ## Format code
	@echo "==> Formatting code..."
	@go fmt ./...

tidy: ## Tidy go.mod
	@echo "==> Tidying go modules..."
	@go mod tidy

clean: ## Clean up generated files
	@echo "==> Cleaning up..."
	@rm -f coverage.out

bench: ## Run benchmarks
	@echo "==> Running benchmarks..."
	@go test ./benchmarks -bench=. -benchmem

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\\033[36m%-20s\\033[0m %s\\n", $$1, $$2}'

.DEFAULT_GOAL := help
