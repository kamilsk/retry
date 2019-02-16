PACKAGES := go list ./... | grep -v vendor | grep -v ^_
SHELL    ?= /bin/bash -euo pipefail


.PHONY: test
test:                         #| Runs tests with race.
                              #| Accepts: ARGS.
                              #| Uses: PACKAGES.
	$(PACKAGES) | xargs go test -race $(strip $(ARGS))

.PHONY: test-check
test-check:                   #| Fast runs tests to check their compilation errors.
                              #| Accepts: ARGS.
                              #| Uses: PACKAGES.
	$(PACKAGES) | xargs go test -run=^hack $(strip $(ARGS))

.PHONY: test-with-coverage
test-with-coverage:           #| Runs tests with coverage.
                              #| Accepts: ARGS.
                              #| Uses: PACKAGES.
	$(PACKAGES) | xargs go test -cover $(strip $(ARGS))

.PHONY: test-with-coverage-formatted
test-with-coverage-formatted: #| Runs tests with coverage and formats the result.
                              #| Accepts: ARGS.
                              #| Uses: PACKAGES.
	$(PACKAGES) | xargs go test -cover $(strip $(ARGS)) | column -t | sort -r

.PHONY: test-with-coverage-profile
test-with-coverage-profile:   #| Runs tests with coverage and collects the result.
                              #| Accepts: ARGS.
                              #| Uses: PACKAGES.
	echo 'mode: count' > cover.out
	for package in $$($(PACKAGES)); do \
	    go test -covermode count \
	            -coverprofile "coverage_$${package##*/}.out" \
	            $(strip $(ARGS)) "$${package}"; \
	    if [ -f "coverage_$${package##*/}.out" ]; then \
	        sed 1d "coverage_$${package##*/}.out" >> cover.out; \
	        rm "coverage_$${package##*/}.out"; \
	    fi \
	done

.PHONY: test-example
test-example:                 #| Runs example tests with coverage and collects the result.
                              #| Accepts: ARGS.
                              #| Uses: PACKAGES.
	echo 'mode: count' > coverage_example.out
	for package in $$($(PACKAGES)); do \
	    go test -v -run=Example \
	            -covermode count \
	            -coverprofile "coverage_example_$${package##*/}.out" \
	            $(strip $(ARGS)) "$${package}"; \
	    if [ -f "coverage_$${package##*/}.out" ]; then \
	        sed 1d "coverage_example_$${package##*/}.out" >> coverage_example.out; \
	        rm "coverage_example_$${package##*/}.out"; \
	    fi \
	done
