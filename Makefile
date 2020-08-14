.DEFAULT_GOAL := help

SHELL ?= bash

DESTDIR ?= /usr/local

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
	go build -ldflags="-s -w" -buildmode c-archive -o .build/uplink.a .
	mv .build/uplink.so .build/libuplink.so
	mv .build/uplink.a .build/libuplink.a
	mkdir -p .build/uplink
	mv .build/*.h .build/uplink
	cp uplink_definitions.h .build/uplink
	./scripts/gen-pkg-config > .build/libuplink.pc

.PHONY: build-gpl2
build-gpl2: ## builds the Linux dynamic libraries GPL2 license compatible and leave them and a copy of the definitions in .build directory
	./check-licenses.sh
	go build -modfile=go-gpl2.mod -ldflags="-s -w" -buildmode c-shared -tags stdsha256 -o .build/uplink.so .
	go build -modfile=go-gpl2.mod -ldflags="-s -w" -buildmode c-shared -tags stdsha256 -o .build/uplink.a .
	mv .build/uplink.so .build/libuplink.so
	mv .build/uplink.a .build/libuplink.a
	mkdir -p .build/uplink
	mv .build/*.h .build/uplink
	cp uplink_definitions.h .build/uplink
	./scripts/gen-pkg-config > .build/libuplink.pc

.PHONY: bump-dependencies
bump-dependencies: ## bumps the dependencies
	go get storj.io/common@master storj.io/uplink@master
	go mod tidy
	cd testsuite;\
		go get storj.io/common@master storj.io/storj@master storj.io/uplink@master;\
		go mod tidy

.PHONY: test
test: ## run test suite
	cd testsuite && go test

.PHONY: test-install
test-install: ## test install process
	./scripts/test-install

.PHONY: install
install: build ## install library and headers
	install -d \
		${DESTDIR}/include/uplink \
		${DESTDIR}/lib \
		${DESTDIR}/lib/pkgconfig
	install .build/libuplink.so ${DESTDIR}/lib
	install .build/libuplink.a ${DESTDIR}/lib
	install --mode 644 .build/uplink/* ${DESTDIR}/include/uplink
	install --mode 644 .build/libuplink.pc ${DESTDIR}/lib/pkgconfig
