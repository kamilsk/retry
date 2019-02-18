.PHONY: deps
deps:
	@(go mod tidy && go mod vendor && go mod verify)

.PHONY: goimports
goimports:
	@(goimports -ungroup -w .)

.PHONY: test
test:
	@(go test -cover -race -timeout 1s -v ./...)
