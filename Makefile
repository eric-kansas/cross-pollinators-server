SHELL := /bin/bash

.PHONY: build run test

build:
	# Building!
	go build

run:
	# Running!
	./cross-pollinators-server

deps:
	dep ensure

compiledaemon:
	CompileDaemon -build='make build' -command='make run' -include=Makefile -graceful-kill=true

test_dependencies:
	go get -u github.com/golang/dep