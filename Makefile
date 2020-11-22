# sourced by https://github.com/octomation/makefiles

.DEFAULT_GOAL = test-with-coverage
GIT_HOOKS     = post-merge pre-commit pre-push
GO_VERSIONS   = 1.11 1.12 1.13 1.14 1.15
GO111MODULE   = on
SHELL         = /bin/bash -euo pipefail

AT    := @
OS    := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH  := $(shell uname -m | tr '[:upper:]' '[:lower:]')

SHELL ?= /bin/bash -euo pipefail

verbose:
	$(eval AT :=)
	@echo > /dev/null
.PHONY: verbose

todo:
	@grep \
		--exclude=Makefile \
		--exclude-dir={bin,components,node_modules,vendor} \
		--color \
		--text \
		-nRo -E ' TODO:.*|SkipNow' . || true
.PHONY: todo

rmdir:
	$(AT) for dir in `git ls-files --others --exclude-standard --directory`; do \
		find $${dir%%/} -depth -type d -empty | xargs rmdir; \
	done
.PHONY: rmdir

GO111MODULE ?= on
GOFLAGS     ?= -mod=
GOPRIVATE   ?= go.octolab.net
GOPROXY     ?= direct
LOCAL       ?= $(MODULE)
MODULE      ?= `GO111MODULE=on go list -m $(GOFLAGS)`
PACKAGES    ?= `GO111MODULE=on go list $(GOFLAGS) ./...`
PATHS       ?= $(shell echo $(PACKAGES) | sed -e "s|$(MODULE)/||g" | sed -e "s|$(MODULE)|$(PWD)/*.go|g")
TIMEOUT     ?= 1s

ifeq (, $(PACKAGES))
	PACKAGES = $(MODULE)
endif

ifeq (, $(PATHS))
	PATHS = .
endif

export GO111MODULE := $(GO111MODULE)
export GOFLAGS     := $(GOFLAGS)
export GOPRIVATE   := $(GOPRIVATE)
export GOPROXY     := $(GOPROXY)

go-env:
	@echo "GO111MODULE: `go env GO111MODULE`"
	@echo "GOFLAGS:     $(strip `go env GOFLAGS`)"
	@echo "GOPRIVATE:   $(strip `go env GOPRIVATE`)"
	@echo "GOPROXY:     $(strip `go env GOPROXY`)"
	@echo "LOCAL:       $(LOCAL)"
	@echo "MODULE:      $(MODULE)"
	@echo "PACKAGES:    $(PACKAGES)"
	@echo "PATHS:       $(strip $(PATHS))"
	@echo "TIMEOUT:     $(TIMEOUT)"
.PHONY: go-env

export GOBIN := $(PWD)/bin/$(OS)/$(ARCH)

deps-check:
	@go mod verify
	@if command -v egg > /dev/null; then \
		egg deps check license; \
		egg deps check version; \
	fi
.PHONY: deps-check

deps-clean:
	@go clean -modcache
.PHONY: deps-clean

deps-fetch:
	@go mod download
	@if [[ "`go env GOFLAGS`" =~ -mod=vendor ]]; then go mod vendor; fi
.PHONY: deps-fetch

deps-tidy:
	@go mod tidy
	@if [[ "`go env GOFLAGS`" =~ -mod=vendor ]]; then go mod vendor; fi
.PHONY: deps-tidy

deps-update: selector = '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}'
deps-update:
	$(AT) if command -v egg > /dev/null; then \
		packages="`egg deps list`"; \
	else \
		packages="`go list -f $(selector) -m -mod=readonly all`"; \
	fi; \
	if [[ "`go version`" == *1.1[1-3]* ]]; then \
		go get -d -mod= -u $$packages; \
	else \
		go get -d -u $$packages; \
	fi; \
	if [[ "`go env GOFLAGS`" =~ -mod=vendor ]]; then go mod vendor; fi
.PHONY: deps-update

deps-update-all:
	$(AT) if [[ "`go version`" == *1.1[1-3]* ]]; then \
		go get -d -mod= -u ./...; \
	else \
		go get -d -u ./...; \
	fi; \
	if [[ "`go env GOFLAGS`" =~ -mod=vendor ]]; then go mod vendor; fi
.PHONY: deps-update-all

go-fmt:
	@if command -v goimports > /dev/null; then \
		goimports -local $(LOCAL) -ungroup -w $(PATHS); \
	else \
		gofmt -s -w $(PATHS); \
	fi
.PHONY: go-fmt

go-generate:
	@go generate $(PACKAGES)
.PHONY: go-generate

lint:
	@golangci-lint run ./...
	@looppointer ./...
.PHONY: lint

GODOC_HOST ?= localhost:6060

docs:
	@(sleep 2 && open http://$(GODOC_HOST)/pkg/$(LOCAL)/) &
	@godoc -http=$(GODOC_HOST)
.PHONY: docs

test:
	@go test -race -timeout $(TIMEOUT) $(PACKAGES)
.PHONY: test

test-clean:
	@go clean -testcache
.PHONY: test-clean

test-quick: GOTAGS = integration,tools
test-quick:
	@go test -run ^Fake$$ -tags $(GOTAGS) ./... | { grep -v 'no tests to run' || true; }
	@go test -timeout $(TIMEOUT) $(PACKAGES)
.PHONY: test-quick

test-verbose:
	@go test -race -timeout $(TIMEOUT) -v $(PACKAGES)
.PHONY: test-verbose

test-with-coverage:
	@go test \
		-cover \
		-covermode atomic \
		-coverprofile c.out \
		-race \
		-timeout $(TIMEOUT) \
		$(PACKAGES) | column -t | sort -r
.PHONY: test-with-coverage

test-with-coverage-report: test-with-coverage
	@go tool cover -html c.out
.PHONY: test-with-coverage-report

test-integration: GOTAGS = integration
test-integration:
	@go test \
		-cover \
		-covermode atomic \
		-coverprofile integration.out \
		-race \
		-tags $(GOTAGS) \
		./... | column -t | sort -r
.PHONY: test-integration

test-integration-quick: GOTAGS = integration
test-integration-quick:
	@go test -tags $(GOTAGS) ./...
.PHONY: test-integration-quick

test-integration-report: test-integration
	@go tool cover -html integration.out
.PHONY: test-integration-report

TOOLFLAGS ?= -mod=

tools-env:
	@echo "GOBIN:       `go env GOBIN`"
	@echo "TOOLFLAGS:   $(TOOLFLAGS)"
.PHONY: tools-env

toolset: GOTAGS = tools
toolset:
	$(AT) ( \
		GOFLAGS=$(TOOLFLAGS); \
		cd tools; \
		go mod tidy; \
		if [[ "`go env GOFLAGS`" =~ -mod=vendor ]]; then go mod vendor; fi; \
		go generate -tags $(GOTAGS) tools.go; \
	)
.PHONY: toolset

ifdef GIT_HOOKS

hooks: unhook
	$(AT) for hook in $(GIT_HOOKS); do cp githooks/$$hook .git/hooks/; done
.PHONY: hooks

unhook:
	@ls .git/hooks | grep -v .sample | sed 's|.*|.git/hooks/&|' | xargs rm -f || true
.PHONY: unhook

define hook_tpl
$(1):
	@githooks/$(1)
.PHONY: $(1)
endef

render_hook_tpl = $(eval $(call hook_tpl,$(hook)))
$(foreach hook,$(GIT_HOOKS),$(render_hook_tpl))

endif

git-check:
	$(AT) git diff --exit-code >/dev/null
	$(AT) git diff --cached --exit-code >/dev/null
	$(AT) ! git ls-files --others --exclude-standard | grep -q ^
.PHONY: git-check

ifdef GO_VERSIONS

define go_tpl
go$(1):
	@docker run \
		--rm -it \
		-v $(PWD):/src \
		-w /src \
		golang:$(1) bash
.PHONY: go$(1)
endef

render_go_tpl = $(eval $(call go_tpl,$(version)))
$(foreach version,$(GO_VERSIONS),$(render_go_tpl))

endif


export PATH := $(GOBIN):$(PATH)


init: deps test lint hooks
	@git config core.autocrlf input
.PHONY: init

clean: deps-clean test-clean
.PHONY: clean

deps: deps-fetch toolset
.PHONY: deps

env: go-env tools-env
env:
	@echo "PATH:        $(PATH)"
.PHONY: env

format: go-fmt
.PHONY: format

generate: go-generate format
.PHONY: generate

refresh: deps-tidy update deps generate test
.PHONY: refresh

update: deps-update
.PHONY: update

verify: deps-check generate test lint git-check
.PHONY: verify
