.DEFAULT_GOAL := help


DESTDIR ?= /usr/local
GPL2 ?= false
NO_QUIC ?= false

ifeq (${GPL2},true)
GOFLAGS += -modfile=go-gpl2.mod -tags=stdsha256
export GOFLAGS
endif

ifeq (${NO_QUIC},true)
GOFLAGS += -tags=noquic
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


OS ?= $(shell uname)
ZIG ?= $(shell which zig)
TAG ?= $(shell git tag -l | grep -E 'v[0-9]+\.[0-9]+\.[0-9]+' | sort -V | tail -n 1)

.PHONY: release-files
release-files: ## builds and copies release files to Github

## MacOS is required by license to build MacOS libraries
## Forcing it here just prevents things from being half-done
ifneq ($(OS),Darwin)
	echo This tool must be run from MacOS
	exit 1
endif
ifndef GITHUB_TOKEN
	echo GITHUB_TOKEN is undefined
	exit 1
endif

	echo "Uploading binaries to release draft"

	git checkout $(TAG)

##	linux_amd64
	CGO_ENABLED=1 CC="${ZIG} cc -target x86_64-linux-gnu" GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -buildmode c-shared -o .build/uplink_linux_amd64.so .
	github-release upload --user storj --repo uplink-c --tag "$(TAG)" --name "uplink_linux_amd64.so" --file ".build/uplink_linux_amd64.so"
##	linux_arm64
	CGO_ENABLED=1 CC="${ZIG} cc -target aarch64-linux-gnu" GOOS=linux go build -ldflags="-s -w" -buildmode c-shared -o .build/uplink_linux_arm64.so .
	github-release upload --user storj --repo uplink-c --tag "$(TAG)" --name "uplink_linux_arm64.so" --file ".build/uplink_linux_arm64.so"
##	windows_amd64
	CGO_ENABLED=1 CC="${ZIG} cc -target x86_64-windows" GOARCH=amd64 GOOS=windows go build -ldflags="-s -w" -buildmode c-shared -o .build/uplink_windows_amd64.dll .
	github-release upload --user storj --repo uplink-c --tag "$(TAG)" --name "uplink_windows_amd64.dll" --file ".build/uplink_windows_amd64.dll"
##	darwin_amd64
	CGO_ENABLED=1 GOARCH=amd64 go build -ldflags="-s -w" -buildmode c-shared -o .build/uplink_darwin_amd64.dylib .
	github-release upload --user storj --repo uplink-c --tag "$(TAG)" --name "uplink_darwin_amd64.dylib" --file ".build/uplink_darwin_amd64.dylib"
##	darwin_arm64
	CGO_ENABLED=1 GOARCH=arm64 go build -ldflags="-s -w" -buildmode c-shared -o .build/uplink_darwin_arm64.dylib .
	github-release upload --user storj --repo uplink-c --tag "$(TAG)" --name "uplink_darwin_arm64.dylib" --file ".build/uplink_darwin_arm64.dylib"

	echo "Uploading release binaries done"