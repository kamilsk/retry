ifndef PACKAGES
$(error Please include env.mk before)
endif


.PHONY: install-deps
install-deps:
	$(PACKAGES) | xargs go get -d -t

.PHONY: update-deps
update-deps:
	$(PACKAGES) | xargs go get -d -t -u


.PHONY: build
build:
	$(PACKAGES) | xargs go build -v $(strip $(ARGS))

.PHONY: clean
clean:
	$(PACKAGES) | xargs go clean -i -x

.PHONY: install
install:
	$(PACKAGES) | xargs go install

.PHONY: vet
vet:
	$(PACKAGES) | xargs go vet


OPEN_BROWSER              ?= true
GO_TEST_COVERAGE_MODE     ?= count
GO_TEST_COVERAGE_FILENAME ?= coverage.out

.PHONY: test
test:
	$(PACKAGES) | xargs go test -race $(strip $(ARGS))

.PHONY: test-check
test-check:
	$(PACKAGES) | xargs go test -run=^hack $(strip $(ARGS))

.PHONY: test-with-coverage
test-with-coverage:
	$(PACKAGES) | xargs go test -cover $(strip $(ARGS))

.PHONY: test-with-coverage-formatted
test-with-coverage-formatted:
	$(PACKAGES)| xargs go test -cover $(strip $(ARGS)) | column -t | sort -r

.PHONY: test-with-coverage-profile
test-with-coverage-profile:
	echo 'mode: ${GO_TEST_COVERAGE_MODE}' > '${GO_TEST_COVERAGE_FILENAME}'
	for package in $$($(PACKAGES)); do \
	    go test -covermode '${GO_TEST_COVERAGE_MODE}' \
	            -coverprofile "coverage_$${package##*/}.out" \
	            $(strip $(ARGS)) "$${package}"; \
	    sed '1d' "coverage_$${package##*/}.out" >> '${GO_TEST_COVERAGE_FILENAME}'; \
	    rm "coverage_$${package##*/}.out"; \
	done
	if [ '$(OPEN_BROWSER)' != '' ]; then go tool cover -html='${GO_TEST_COVERAGE_FILENAME}'; fi

.PHONY: test-example
test-example: GO_TEST_COVERAGE_FILENAME = coverage_example.out
test-example:
	echo 'mode: ${GO_TEST_COVERAGE_MODE}' > '${GO_TEST_COVERAGE_FILENAME}'
	for package in $$($(PACKAGES)); do \
	    go test -v -run=Example \
	            -covermode '${GO_TEST_COVERAGE_MODE}' \
	            -coverprofile "coverage_example_$${package##*/}.out" \
	            $(strip $(ARGS)) "$${package}"; \
	    sed '1d' "coverage_example_$${package##*/}.out" >> '${GO_TEST_COVERAGE_FILENAME}'; \
	    rm "coverage_example_$${package##*/}.out"; \
	done
	if [ '$(OPEN_BROWSER)' != '' ]; then go tool cover -html='${GO_TEST_COVERAGE_FILENAME}'; fi


.PHONY: bench
bench:
	$(PACKAGES) | xargs go test -bench=. $(strip $(ARGS))
