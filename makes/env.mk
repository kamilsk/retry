ifndef GOPATH
$(error $$GOPATH not set)
endif

CWD        := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))
CID        := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
DATE       := $(shell date -u "+%Y-%m-%d_%H:%M:%S")
GO_VERSION := $(shell go version | awk '{print $$3}' | tr -d 'go')

GIT_REV    := $(shell git rev-parse --short HEAD)
GO_PACKAGE := $(patsubst %/,%,$(subst $(GOPATH)/src/,,$(CWD)))
PACKAGES   := go list ./... | grep -v vendor | grep -v ^_

SHELL      ?= /bin/bash -euo pipefail

.PHONY: help
help:
	@fgrep -h "#|" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/#| //'
	# TODO make -pnrR
