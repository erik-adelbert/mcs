# Makefile for mcs

SHELL := /bin/bash

check :
	cd examples/gomer-uct; go run ./main.go -f ../../assets/www.js-game.de/problem01.txt

bench :
	cd examples/benchmarks; go test -v -timeout 0

clean :
	cd examples/benchmarks; rm *log
