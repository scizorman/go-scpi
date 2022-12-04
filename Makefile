.DEFAULT_GOAL := help

.PHONY: help
help: ## Self documenting help output
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: fmt
fmt: ## Format code
	go fmt ./...

.PHONY: lint
lint: fmt ## Lint code
	golint ./...

.PHONY: vet
vet: fmt ## Vet code
	go vet ./...

.PHONY: test
test: vet ## Run tests
	go test -v ./...

.PHONY: install
install: test ## Install as executable
	go install -v ./...

.PHONY: build
build: test ## Build executable
	go build -v -o bin/ ./...
