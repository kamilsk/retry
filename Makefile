OPEN_BROWSER       =
SUPPORTED_VERSIONS = 1.5 1.6 1.7 1.8 1.9 latest

include makes/env.mk
include makes/local.mk
include makes/docker.mk

.PHONY: check-code-quality
check-code-quality: ARGS = \
	--exclude='.*_test\.go:.*error return value not checked.*\(errcheck\)$' \
	--exclude='duplicate of.*_test.go.*\(dupl\)$' \
	--vendor --deadline=2m ./...
check-code-quality: docker-tool-gometalinter



.PHONY: cmd-test
cmd-test:
	docker run --rm \
	           -v '$(GOPATH)/src/$(GO_PACKAGE)':'/go/src/$(GO_PACKAGE)' \
	           -w '/go/src/$(GO_PACKAGE)' \
	           golang:1.7 \
	           /bin/sh -c 'go install -ldflags "-X 'main.Version=test'" ./cmd/retry \
	                       && retry -limit=3 -backoff=lin{10ms} -- /bin/sh -c "echo 'trying...'; exit 1"'



.PHONY: pull-github-tpl
pull-github-tpl:
	rm -rf .github
	(git clone git@github.com:kamilsk/shared.git .github && cd .github && git checkout github-tpl-go-v1 \
	  && echo 'github templates at revision' $$(git rev-parse HEAD) && rm -rf .git)
	rm .github/README.md

.PHONY: pull-makes
pull-makes:
	rm -rf makes
	(git clone git@github.com:kamilsk/shared.git makes && cd makes && git checkout makefile-go-v1 \
	  && echo 'makes at revision' $$(git rev-parse HEAD) && rm -rf .git)
