GOLANGCI_LINT_VERSION ?= 1.51.1

/go/bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOPATH}/bin v${GOLANGCI_LINT_VERSION}

build-cli:
	go install ./cli/battlesnake
.PHONY: cli

test: test-format test-lint test-unit
.PHONY: test

test-format:
	test -z $(gofmt -l .) || (gofmt -l . && exit 1)
.PHONY: test-format

test-lint: /go/bin/golangci-lint
	/go/bin/golangci-lint run -v ./...
.PHONY: test-lint

test-unit:
	go test -race ./...
.PHONY: test-unit
