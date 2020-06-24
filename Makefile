# format-c formats all C-code
format-c:
	cd testsuite/testdata && clang-format --style=file -i *.c *.h
	clang-format --style=file -i *.h

build:
	go build -ldflags="-s -w" -buildmode c-shared -o .build/uplink.so .
	mv .build/uplink.so .build/libuplink.so
	cp uplink_definitions.h .build/

build-gpl2:
	./check-licenses.sh
	go build -modfile=go-gpl2.mod -ldflags="-s -w" -buildmode c-shared -tags stdsha256 -o .build/uplink.so .
	mv .build/uplink.so .build/libuplink.so
	cp uplink_definitions.h .build/

.PHONY: bump-dependencies
bump-dependencies:
	go get storj.io/common@master storj.io/uplink@master
	go mod tidy
	cd testsuite;\
		go get storj.io/common@master storj.io/storj@master storj.io/uplink@master;\
		go mod tidy
