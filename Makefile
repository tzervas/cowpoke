.PHONY: build test docker lint clean help

# Build variables
BINARY_NAME=cowpoke
VERSION?=dev
GOFLAGS=-ldflags "-s -w -X main.version=$(VERSION)"
DOCKER_IMAGE=cowpoke
DOCKER_TAG?=latest

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@go build $(GOFLAGS) -o bin/$(BINARY_NAME) cmd/cowpoke/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

lint: ## Run linters
	@echo "Running linters..."
	@go vet ./...
	@go fmt ./...

docker: ## Build Docker image
	@echo "Building Docker image $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean

.DEFAULT_GOAL := help
