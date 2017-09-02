SHELL := /bin/bash

.PHONY: build run test

build:
	# Building!
	go build

run:
	# Running!
	./beeswax-log-ingestion

deps:
	govendor sync -v

compiledaemon:
	CompileDaemon -build='make build' -command='make run' -include=Makefile -graceful-kill=true

test_dependencies:
	go get -u github.com/kardianos/govendor