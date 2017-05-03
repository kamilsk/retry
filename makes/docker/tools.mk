GITHUB_TOKEN ?=

.PHONY: docker-in-tools
docker-in-tools:
	docker run --rm -it \
	           -e GITHUB_TOKEN='$(GITHUB_TOKEN)' \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
	           kamilsk/go-tools:latest \
	           /bin/sh

.PHONY: docker-pull-tools
docker-pull-tools:
	docker pull kamilsk/go-tools:latest

.PHONY: docker-tool-glide
docker-tool-glide: COMMAND = install
docker-tool-glide:
	docker run --rm \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
	           kamilsk/go-tools:latest \
	           /bin/sh -c 'glide $(COMMAND) --strip-vendor $(strip $(ARGS)) && \
	                       rm -rf /go/src/.glide'

.PHONY: docker-tool-gometalinter
docker-tool-gometalinter:
	docker run --rm \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
	           kamilsk/go-tools:latest \
	           /bin/sh -c '$(PACKAGES) | xargs go test -i && \
	                       gometalinter --vendor $(strip $(ARGS)) ./...'

.PHONY: docker-tool-goreleaser
docker-tool-goreleaser:
	docker run --rm \
	           -e GITHUB_TOKEN='${GITHUB_TOKEN}' \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
	           kamilsk/go-tools:latest \
	           goreleaser $(strip $(ARGS))
