.PHONY: deps
deps:
	@(go mod tidy && go mod verify)

.PHONY: goimports
goimports:
	@(goimports -ungroup -w .)

.PHONY: test
test:
	@(go test -cover -race -v ./...)
