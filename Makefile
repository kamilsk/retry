include makes/env.mk
include makes/bench.mk
include makes/deps.mk
include makes/docker.mk
include makes/tests.mk
include makes/tools.mk

.PHONY: all
all: install-deps build install

.PHONY: docker-bench
docker-bench: ARGS := -benchmem $(ARGS)
# blocked by https://github.com/kamilsk/shared/issues/65
# docker-bench: docker-bench-1.5
docker-bench: docker-bench-1.6
docker-bench: docker-bench-1.7
docker-bench: docker-bench-latest

.PHONY: docker-gometalinter
docker-gometalinter: ARGS := --deadline=30s $(ARGS)
docker-gometalinter: docker-tool-gometalinter

.PHONY: docker-pull
docker-pull: docker-pull-1.5
docker-pull: docker-pull-1.6
docker-pull: docker-pull-1.7
docker-pull: docker-pull-latest
docker-pull: docker-pull-tools
docker-pull: docker-clean

.PHONY: docker-test
docker-test: ARGS := -v $(ARGS)
# blocked by https://github.com/kamilsk/shared/issues/65
# docker-test: docker-test-1.5
docker-test: docker-test-1.6
docker-test: docker-test-1.7
docker-test: docker-test-latest

.PHONY: docker-test-with-coverage
docker-test-with-coverage: ARGS := -v $(ARGS)
docker-test-with-coverage: OPEN_BROWSER := true
docker-test-with-coverage: docker-test-1.5-with-coverage
docker-test-with-coverage: docker-test-1.6-with-coverage
docker-test-with-coverage: docker-test-1.7-with-coverage
docker-test-with-coverage: docker-test-latest-with-coverage
