SHELL := /bin/bash -euo pipefail


.PHONY: generate
generate:
	go generate ./cmd/generate
	mv ./cmd/generate/parser_gen.go ./cmd/retry/parser_gen.go


.PHONY: test
test:                         #| Runs tests with race.
	go test -race ./...

.PHONY: test-check
test-check:                   #| Fast runs tests to check their compilation errors.
	go test -run=^hack ./...

.PHONY: test-with-coverage
test-with-coverage:           #| Runs tests with coverage.
	go test -cover ./...

.PHONY: test-with-coverage-formatted
test-with-coverage-formatted: #| Runs tests with coverage and formats the result.
	go test -cover ./... | column -t | sort -r

.PHONY: test-with-coverage-profile
test-with-coverage-profile:   #| Runs tests with coverage and collects the result.
	go test -covermode count -coverprofile cover.out ./...

.PHONY: test-example
test-example:                 #| Runs example tests with coverage and collects the result.
	go test -covermode count -coverprofile -run=Example -v example.out ./...
