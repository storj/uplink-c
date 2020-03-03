# format-c formats all C-code
format-c:
	cd testsuite/testdata && clang-format --style=file -i *.c *.h
	clang-format --style=file -i *.h