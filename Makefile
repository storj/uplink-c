# format-c formats all C-code
format-c:
	cd testsuite/testdata && clang-format --style=file -i *.c *.h
	clang-format --style=file -i *.h


build:
	go build -ldflags="-s -w" -buildmode c-shared -o .build/uplink.so .
	mv .build/uplink.so .build/libuplink.so
	cp uplink_definitions.h .build/