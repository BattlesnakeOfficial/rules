GOPATH := $(shell go env GOPATH)

GOLANGCI_LINT_PATH		:= ${GOPATH}/bin/golangci-lint
GOLANGCI_LINT_VERSION	:= 1.55.2


${GOLANGCI_LINT_PATH}:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOPATH}/bin v${GOLANGCI_LINT_VERSION}

install-cli:
	go install ./cli/battlesnake
.PHONY: install-cli

test-format:
	test -z $$(gofmt -l .) || (gofmt -l . && exit 1)
.PHONY: test-format

test-lint: ${GOLANGCI_LINT_PATH}
	golangci-lint run -v ./...
.PHONY: test-lint

test-unit:
	go test -race ./...
.PHONY: test-unit

test: test-format test-lint test-unit
.PHONY: test
