GIT_ORIGIN?="git@github.com:kamilsk/retrier.git"
GIT_MIRROR?="git@bitbucket.org:kamilsk/retrier.git"
GO_TEST_COVERAGE_MODE?="count"
GO_TEST_COVERAGE_FILE_NAME?="coverage.out"
GOFMT_FLAGS?="-s"
GOLINT_MIN_CONFIDENCE?="0.3"


.PHONY: all build
.PHONY: install install-deps
.PHONY: update-deps
.PHONY: test test-with-coverage test-with-coverage-formatted test-with-coverage-profile
.PHONY: clean vet
.PHONY: dev publish


all: install-deps build install

build:
	go build -v ./...

install:
	go install ./...

install-deps:
	go get -d -t ./...

update-deps:
	go get -d -t -u ./...

test:
	go test -race -v ./...

test-with-coverage:
	go test -cover -race ./...

test-with-coverage-formatted:
	go test -cover -race ./... | column -t | sort -r

test-with-coverage-profile:
	echo "mode: ${GO_TEST_COVERAGE_MODE}" > ${GO_TEST_COVERAGE_FILE_NAME}
	for package in $$(go list ./...); do \
		go test -covermode ${GO_TEST_COVERAGE_MODE} -coverprofile "coverage_$${package##*/}.out" -race "$${package}"; \
		sed '1d' "coverage_$${package##*/}.out" >> ${GO_TEST_COVERAGE_FILE_NAME}; \
	done

clean:
	go clean -i -x ./...

vet:
	go vet ./...

dev:
	git remote set-url origin $(GIT_ORIGIN)
	git remote add mirror $(GIT_MIRROR)

publish:
	git push origin master --tags
	git push mirror master --tags
