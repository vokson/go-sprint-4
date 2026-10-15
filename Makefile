GO            ?= go
GOLANGCI_LINT ?= golangci-lint

.PHONY: test lint fmt

test:
	$(GO) test ./...

lint:
	$(GOLANGCI_LINT) run ./internal/...

fmt:
	$(GOLANGCI_LINT) fmt ./...
