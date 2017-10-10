#
# Temporary file for experiments.
#
# make -pnrR

SUPPORTED_VERSIONS ?= 1.5 1.6 1.7 1.8 latest

include env.mk
include local.mk
include docker.mk
include docker/hugo.mk

ARGS = a

test1: ARGS = 1
test1:
	echo $(ARGS)

test2: override ARGS += 2
test2:
	echo $(ARGS)
