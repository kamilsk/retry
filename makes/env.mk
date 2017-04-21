ifndef GOPATH
$(error $$GOPATH not set)
endif


ARGS =

CWD        := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))
CID        := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
DATE       := $(shell date -u "+%Y-%m-%d %H:%M:%S")
GO_VERSION := $(shell go version | awk '{print $$3}' | tr -d 'go')

GIT_REV    := $(shell git rev-parse --short HEAD)
GO_PACKAGE := $(patsubst %/,%,$(subst $(GOPATH)/src/,,$(CWD)))
PACKAGES   := go list ./... | grep -v vendor | grep -v ^_
