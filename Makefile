include makes/env.mk
include makes/local.mk
include makes/docker.mk

OPEN_BROWSER        =
SUPPORTED_VERSIONS ?= 1.5 1.6 1.7 1.8 latest

.PHONY: check-code-quality
check-code-quality: ARGS = --vendor --deadline=1m ./...
check-code-quality: docker-tool-gometalinter

.PHONY: complex-bench
complex-bench: ARGS = -benchmem
complex-bench: docker-bench-1.5
complex-bench: docker-bench-1.6
complex-bench: docker-bench-1.7
complex-bench: docker-bench-1.8
complex-bench: docker-bench-latest

.PHONY: complex-tests
complex-tests: ARGS = -timeout=1s
complex-tests: docker-test-1.5
complex-tests: docker-test-1.6
complex-tests: docker-test-1.7
complex-tests: docker-test-1.8
complex-tests: docker-test-latest

.PHONY: complex-tests-with-coverage
complex-tests-with-coverage: ARGS = -timeout=1s
complex-tests-with-coverage: docker-test-with-coverage-1.5
complex-tests-with-coverage: docker-test-with-coverage-1.6
complex-tests-with-coverage: docker-test-with-coverage-1.7
complex-tests-with-coverage: docker-test-with-coverage-1.8
complex-tests-with-coverage: docker-test-with-coverage-latest

.PHONY: docker-pull
docker-pull: docker-pull-1.5
docker-pull: docker-pull-1.6
docker-pull: docker-pull-1.7
docker-pull: docker-pull-1.8
docker-pull: docker-pull-latest
docker-pull: docker-pull-tools
docker-pull: PRUNE = --force
docker-pull: docker-clean

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

.PHONY: research
research: COMMAND = -y research.yml install
research: ARGS    = --strip-vendor
research: docker-tool-glide
research:
	rm -rf .glide

.PHONY: test-cmd
test-cmd:
	docker run --rm \
	           -v '$(GOPATH)/src/$(GO_PACKAGE)':'/go/src/$(GO_PACKAGE)' \
	           -w '/go/src/$(GO_PACKAGE)' \
	           golang:1.7 \
	           /bin/sh -c 'go install -ldflags "-X 'main.Version=test'" ./cmd/retry \
	                       && retry -limit=3 -backoff=lin{10ms} -- /bin/sh -c "echo 'trying...'; exit 1"'
