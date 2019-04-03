SHELL := /bin/bash -euo pipefail
PKGS  := go list ./... | grep -v vendor | grep -v ^_


.PHONY: deps
deps:
	@(go mod tidy && go mod vendor && go mod verify)

.PHONY: update
update:
	@(go get -mod= -u)


.PHONY: format
format:
	@(goimports -ungroup -w .)


.PHONY: test
test:                         #| Runs tests with race.
	@(go test -race -timeout 1s ./...)

.PHONY: test-check
test-check:                   #| Fast runs tests to check their compilation errors.
	@(go test -run=^hack ./...)

.PHONY: test-with-coverage
test-with-coverage:           #| Runs tests with coverage.
	@(go test -cover -timeout 1s  ./...)

.PHONY: test-with-coverage-formatted
test-with-coverage-formatted: #| Runs tests with coverage and formats the result.
	@(go test -cover -timeout 1s  ./... | column -t | sort -r)

.PHONY: test-with-coverage-profile
test-with-coverage-profile:   #| Runs tests with coverage and collects the result.
	@(go test -covermode count -coverprofile c.out -timeout 1s ./...)

.PHONY: test-with-coverage-profile-old
test-with-coverage-profile-old:
	@(echo 'mode: count' > 'cover.out')
	@(set -e; for package in $$($(PKGS)); do \
	    go test -covermode count \
	            -coverprofile "coverage_$${package##*/}.out" \
	            -timeout 1s "$${package}"; \
	    if [ -f "coverage_$${package##*/}.out" ]; then \
	        sed '1d' "coverage_$${package##*/}.out" >> cover.out; \
	        rm "coverage_$${package##*/}.out"; \
	    fi \
	done)

.PHONY: test-example
test-example:                 #| Runs example tests with coverage and collects the result.
	@(go test -covermode count -coverprofile -run=Example -timeout 1s -v example.out ./...)


.PHONY: sync
sync:
	@(git stash && git pull --rebase && git stash pop || true)

.PHONY: upgrade
upgrade: sync update deps format test-with-coverage-formatted
