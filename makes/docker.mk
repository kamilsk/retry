ifndef CID
$(error Please include env.mk before)
endif

DOCKER_VERSION := $(shell docker version | grep Version | head -1 | awk '{print $$2}')

OPEN_BROWSER       ?= true
SUPPORTED_VERSIONS ?= 1.5 1.6 1.7 1.8 latest

include $(CID)/docker/alpine.mk
include $(CID)/docker/base.mk
include $(CID)/docker/clean.mk
include $(CID)/docker/tools.mk
