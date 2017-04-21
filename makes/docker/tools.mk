.PHONY: docker-in-tools
docker-in-tools:
	docker run --rm -it \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
	           kamilsk/go-tools:latest \
	           /bin/sh

.PHONY: docker-pull-tools
docker-pull-tools:
	docker pull kamilsk/go-tools:latest

.PHONY: docker-tool-gometalinter
docker-tool-gometalinter:
	docker run --rm \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
	           kamilsk/go-tools:latest \
	           /bin/sh -c '$(PACKAGES) | xargs go test -i && \
	                       gometalinter --vendor $(strip $(ARGS)) ./...'

.PHONY: docker-tool-glide
docker-tool-glide:
	docker run --rm \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
	           kamilsk/go-tools:latest \
	           /bin/sh -c 'glide install --strip-vendor $(strip $(ARGS)) && \
	                       rm -rf /go/src/.glide'
