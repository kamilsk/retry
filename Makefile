include makes/env.mk
include makes/local.mk
include makes/docker.mk

.PHONY: docker-bench
docker-bench: ARGS := -benchmem $(ARGS)
docker-bench: docker-bench-1.5
docker-bench: docker-bench-1.6
docker-bench: docker-bench-1.7
docker-bench: docker-bench-1.8
docker-bench: docker-bench-latest

.PHONY: docker-gometalinter
docker-gometalinter: ARGS := --deadline=1m $(ARGS)
docker-gometalinter: docker-tool-gometalinter

.PHONY: docker-pull
docker-pull: docker-pull-1.5
docker-pull: docker-pull-1.6
docker-pull: docker-pull-1.7
docker-pull: docker-pull-1.8
docker-pull: docker-pull-latest
docker-pull: docker-pull-tools
docker-pull: docker-clean

.PHONY: docker-test
docker-test: ARGS := -v $(ARGS)
docker-test: docker-test-1.5
docker-test: docker-test-1.6
docker-test: docker-test-1.7
docker-test: docker-test-1.8
docker-test: docker-test-latest

.PHONY: docker-test-with-coverage
docker-test-with-coverage: ARGS := -v $(ARGS)
docker-test-with-coverage: OPEN_BROWSER := true
docker-test-with-coverage: docker-test-with-coverage-1.5
docker-test-with-coverage: docker-test-with-coverage-1.6
docker-test-with-coverage: docker-test-with-coverage-1.7
docker-test-with-coverage: docker-test-with-coverage-1.8
docker-test-with-coverage: docker-test-with-coverage-latest

.PHONY: cmd-test
cmd-test:
	go install -ldflags "-X 'main.Timeout=100ms'" ./cmd/retry
	retry -limit=3 -backoff=lin[10ms] -timeout=200ms -- curl http://unknown.host
