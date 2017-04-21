define docker_base_tpl

.PHONY: docker-in-$(1)
docker-in-$(1):
	docker run --rm -it \
	           -v '$${GOPATH}/src/$${GO_PACKAGE}':'/go/src/$${GO_PACKAGE}' \
	           -w '/go/src/$${GO_PACKAGE}' \
	           golang:$(1) \
	           /bin/sh

.PHONY: docker-bench-$(1)
docker-bench-$(1):
	docker run --rm \
	           -v '$${GOPATH}/src/$${GO_PACKAGE}':'/go/src/$${GO_PACKAGE}' \
	           -w '/go/src/$${GO_PACKAGE}' \
	           golang:$(1) \
	           /bin/sh -c '$$(PACKAGES) | xargs go get -d -t && \
	                       $$(PACKAGES) | xargs go test -bench=. $$(strip $$(ARGS))'

.PHONY: docker-pull-$(1)
docker-pull-$(1):
	docker pull golang:$(1)

.PHONY: docker-test-$(1)
docker-test-$(1):
	docker run --rm \
	           -v '$${GOPATH}/src/$${GO_PACKAGE}':'/go/src/$${GO_PACKAGE}' \
	           -w '/go/src/$${GO_PACKAGE}' \
	           golang:$(1) \
	           /bin/sh -c '$$(PACKAGES) | xargs go get -d -t && \
	                       $$(PACKAGES) | xargs go test -race $$(strip $$(ARGS))'

.PHONY: docker-test-check-$(1)
docker-test-check-$(1):
	docker run --rm \
	           -v '$${GOPATH}/src/$${GO_PACKAGE}':'/go/src/$${GO_PACKAGE}' \
	           -w '/go/src/$${GO_PACKAGE}' \
	           golang:$(1) \
	           /bin/sh -c '$$(PACKAGES) | xargs go get -d -t && \
	                       $$(PACKAGES) | xargs go test -run=^hack $$(strip $$(ARGS))'

.PHONY: docker-test-with-coverage-$(1)
docker-test-with-coverage-$(1):
	docker run --rm \
	           -v '$${GOPATH}/src/$${GO_PACKAGE}':'/go/src/$${GO_PACKAGE}' \
	           -w '/go/src/$${GO_PACKAGE}' \
	           golang:$(1) \
	           /bin/sh -c '$$(PACKAGES) | xargs go get -d -t; \
	                       echo "mode: $${GO_TEST_COVERAGE_MODE}" > '$$@.out'; \
	                       for package in $$$$($$(PACKAGES)); do \
	                           go test -covermode '$${GO_TEST_COVERAGE_MODE}' \
	                                   -coverprofile "coverage_$$$${package##*/}.out" \
	                                   $$(strip $$(ARGS)) "$$$${package}"; \
	                           sed '1d' "coverage_$$$${package##*/}.out" >> '$$@.out'; \
	                           rm "coverage_$$$${package##*/}.out"; \
	                       done'
	if [ '$$(OPEN_BROWSER)' != '' ]; then go tool cover -html='$$@.out'; fi

endef

$(foreach v,$(SUPPORTED_VERSIONS),$(eval $(call docker_base_tpl,$(v))))
