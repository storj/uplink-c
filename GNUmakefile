.DEFAULT_GOAL := help


DESTDIR ?= /usr/local
GPL2 ?= false

ifeq (${GPL2},true)
GOFLAGS += -modfile=go-gpl2.mod -tags=stdsha256
export GOFLAGS
endif

GO111MODULE=on
export GO111MODULE

.PHONY: help
help:
	@echo "Usage: make [target]"
	@cat Makefile | awk -F ":.*##"  '/##/ { printf "    %-17s %s\n", $$1, $$2 }' | grep -v  grep

.PHONY: format-c
format-c: ## formats all the C code
	find . -type f -iname \*.h -o -iname \*.c | xargs clang-format --style=file -i

.PHONY: format-c-check
format-c-check: ## checks C code formatting
	./scripts/format-c-check

.PHONY: build
build: ## builds the Linux dynamic libraries and leave them and a copy of the definitions in .build directory
ifeq (${GPL2},true)
	cp go.mod go-gpl2.mod
	cp go.sum go-gpl2.sum
	go mod edit -replace github.com/spacemonkeygo/monkit/v3=./internal/replacements/monkit
	./scripts/check-licenses-gpl2
endif
	go build -ldflags="-s -w" -buildmode c-shared -o .build/uplink.so .
	go build -ldflags="-s -w" -buildmode c-archive -o .build/uplink.a .
	mv .build/uplink.so .build/libuplink.so
	mv .build/uplink.a .build/libuplink.a
	mkdir -p .build/uplink
	mv .build/*.h .build/uplink
	cp uplink_definitions.h .build/uplink
	cp uplink_compat.h .build/uplink
	./scripts/gen-pkg-config > .build/libuplink.pc

.PHONY: bump-dependencies
bump-dependencies: ## bumps the dependencies
	go get storj.io/common@main storj.io/uplink@main
	go mod tidy
	cd testsuite;\
		go get storj.io/common@main storj.io/storj@main storj.io/uplink@main;\
		go mod tidy

.PHONY: test
test: ## run test suite
	cd testsuite && go test

.PHONY: test-install
test-install: ## test install process
	./scripts/test-install

.PHONY: test-namespace
test-namespace: ## test namespacing
	./scripts/test-namespace

.PHONY: install
install: build ## install library and headers
	install -d \
		${DESTDIR}/include/uplink \
		${DESTDIR}/lib \
		${DESTDIR}/lib/pkgconfig
	install .build/libuplink.so ${DESTDIR}/lib
	install .build/libuplink.a ${DESTDIR}/lib
	install -m 644 .build/uplink/* ${DESTDIR}/include/uplink
	install -m 644 .build/libuplink.pc ${DESTDIR}/lib/pkgconfig
