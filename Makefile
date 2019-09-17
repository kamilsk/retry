GO111MODULE = on
GOFLAGS     = -mod=vendor
PKGS        = $(shell go list ./... | grep -v vendor)
SHELL       = /bin/bash -euo pipefail
TIMEOUT     = 1s


.DEFAULT_GOAL = test-with-coverage


.PHONY: deps
deps:
	@go mod tidy && go mod vendor && go mod verify

.PHONY: update
update:
	@go get -mod= -u


.PHONY: format
format:
	@goimports -local $(dirname $(go list -m)) -ungroup -w $(PKGS)

.PHONY: generate
generate:
	@go generate $(PKGS)

.PHONY: refresh
refresh: generate format


.PHONY: test
test:
	@go test -race -timeout $(TIMEOUT) $(PKGS)

.PHONY: test-with-coverage
test-with-coverage:
	@go test -cover -timeout $(TIMEOUT) $(PKGS) | column -t | sort -r

.PHONY: test-with-coverage-profile
test-with-coverage-profile:
	@go test -cover -covermode count -coverprofile c.out -timeout $(TIMEOUT) $(PKGS)
