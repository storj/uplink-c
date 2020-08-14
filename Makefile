.DEFAULT_GOAL := help

.PHONY: help
help:
	@echo "Usage: make [target]"
	@cat Makefile | awk -F ":.*##"  '/##/ { printf "    %-17s %s\n", $$1, $$2 }' | grep -v  grep

.PHONY: format-c
format-c: ## formats all the C code
	cd testsuite/testdata && clang-format --style=file -i *.c *.h
	clang-format --style=file -i *.h

.PHONY: format-c-check
format-c-check: ## checks C code formatting
	./scripts/format-c-check

.PHONY: build
build: ## builds the Linux dynamic libraries and leave them and a copy of the definitions in .build directory
	go build -ldflags="-s -w" -buildmode c-shared -o .build/uplink.so .
	mv .build/uplink.so .build/libuplink.so
	cp uplink_definitions.h .build/

.PHONY: build-gpl2
build-gpl2: ## builds the Linux dynamic libraries GPL2 license compatible and leave them and a copy of the definitions in .build directory
	./check-licenses.sh
	go build -modfile=go-gpl2.mod -ldflags="-s -w" -buildmode c-shared -tags stdsha256 -o .build/uplink.so .
	mv .build/uplink.so .build/libuplink.so
	cp uplink_definitions.h .build/

.PHONY: bump-dependencies
bump-dependencies: ## bumps the dependencies
	go get storj.io/common@master storj.io/uplink@master
	go mod tidy
	cd testsuite;\
		go get storj.io/common@master storj.io/storj@master storj.io/uplink@master;\
		go mod tidy
