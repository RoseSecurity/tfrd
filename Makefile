BINARY_NAME=tfrd
VERSION=v1
GO=go

default: help

help: ## List Makefile targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

all: build

fmt: ## Format Go files
	gofumpt -w .

build: ## Build tfrd
	env $(if $(GOOS),GOOS=$(GOOS)) $(if $(GOARCH),GOARCH=$(GOARCH)) $(GO) build -o build/$(BINARY_NAME) -ldflags "-X github.com/RoseSecurity/tfrd/cmd.Version=local"

install: ## Download dependencies
	$(GO) mod download

generate: ## Run tooling
	if [ -d tools ]; then cd tools; go generate ./...; fi

clean: ## Clean up build artifacts
	$(GO) clean
	rm -rf ./build

run: build ## Run tfrd
	./build/$(BINARY_NAME)

docs: build ## Generate documentation
	./build/$(BINARY_NAME) docs

prep: generate build docs ## Prepare for release
	go mod tidy
	go vet

.PHONY: default all build install clean run fmt help generate docs prep
