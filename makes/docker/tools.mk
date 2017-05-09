.PHONY: docker-in-tools
docker-in-tools:
	docker run --rm -it \
	           -e GITHUB_TOKEN='${GITHUB_TOKEN}' \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
	           kamilsk/go-tools:latest \
	           /bin/sh

.PHONY: docker-pull-tools
docker-pull-tools:
	docker pull kamilsk/go-tools:latest

.PHONY: docker-tool-depth
docker-tool-depth:
	docker run --rm \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
	           kamilsk/go-tools:latest \
	           depth $(strip $(ARGS))

.PHONY: docker-tool-apicompat
docker-tool-apicompat:
	docker run --rm \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
	           kamilsk/go-tools:latest \
	           apicompat $(strip $(ARGS))

.PHONY: docker-tool-benchcmp
docker-tool-benchcmp:
	docker run --rm \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
	           kamilsk/go-tools:latest \
	           benchcmp $(strip $(ARGS))

.PHONY: docker-tool-easyjson
docker-tool-easyjson:
	docker run --rm \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
	           kamilsk/go-tools:latest \
	           easyjson $(strip $(ARGS))

.PHONY: docker-tool-glide
docker-tool-glide:
	docker run --rm \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
	           kamilsk/go-tools:latest \
	           glide $(COMMAND) $(strip $(ARGS))

.PHONY: docker-tool-godepq
docker-tool-godepq:
	docker run --rm \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
	           kamilsk/go-tools:latest \
	           godepq $(strip $(ARGS))

.PHONY: docker-tool-gometalinter
docker-tool-gometalinter:
	docker run --rm \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
	           kamilsk/go-tools:latest \
	           /bin/sh -c '$(PACKAGES) | xargs go test -i && \
	                       gometalinter $(strip $(ARGS))'

.PHONY: docker-tool-goreleaser
docker-tool-goreleaser:
	docker run --rm \
	           -e GITHUB_TOKEN='${GITHUB_TOKEN}' \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
	           kamilsk/go-tools:latest \
	           goreleaser $(strip $(ARGS))

.PHONY: docker-tool-goreporter
docker-tool-goreporter:
	docker run --rm \
	           -e GITHUB_TOKEN='${GITHUB_TOKEN}' \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
	           kamilsk/go-tools:latest \
	           goreporter $(strip $(ARGS))

.PHONY: docker-tool-zb
docker-tool-zb:
	docker run --rm \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
	           kamilsk/go-tools:latest \
	           zb $(COMMAND) $(strip $(ARGS))
