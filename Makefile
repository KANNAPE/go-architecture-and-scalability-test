# Variables
BINARY_NAME=upfluence-analysis
MAIN_PATH=./cmd/server/main.go

# Default target
.DEFAULT_GOAL := help

.PHONY: help deps check build run test coverage clean check

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

deps: ## Download and tidy dependencies
	go mod download
	go mod tidy

check: ## Run gofumpt & go vet 
	gofumpt -l -w .
	go vet ./...

build: deps ## Build the binary for production
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

run: ## Run the application locally
	go run $(MAIN_PATH)

test: ## Run all tests with race detector
	go test -v -race ./internal/... # running all tests in verbose mode

coverage: ## Run tests and generate coverage report
	go test -cover ./internal/...

clean: ## Remove binary files
	rm -rf bin/

check: deps build test # Run full CI pipeline