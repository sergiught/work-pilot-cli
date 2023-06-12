GO_BIN ?= $(shell go env GOPATH)/bin

$(GO_BIN)/golangci-lint:
	@echo "==> Installing linter"
	@go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: lint
lint: $(GO_BIN)/golangci-lint ## Run go linter checks
	@echo "==> Running linter"
	@golangci-lint run -v --fix -c .golangci.yml ./...
