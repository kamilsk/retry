SHELL := /bin/bash -euo pipefail


.PHONY: deps
deps:
	@(go mod tidy && go mod vendor && go mod verify)


.PHONY: test
test:
	@(go test -cover -race -timeout 1s -v ./...)
