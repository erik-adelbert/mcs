# Makefile for mcs

SHELL := /bin/bash

check :
	cd cmd/gomer-uct; go run ./main.go -f $$GOPATH/src/mcs/assets/www.js-games.de/problem01.txt

bench :
	cd cmd/benchmarks; go test -v -timeout 0

clean :
	cd cmd/benchmarks; rm *log
