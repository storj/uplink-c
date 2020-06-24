#!/bin/bash

set -e
set -o pipefail

if ! hash go-licenses; then
	echo "Please put the go-licenses tool in your path"
	echo "You can download it from https://github.com/google/go-licenses"
	exit 1
fi

# we build this way for gpl v2 compatibility. github.com/minio/sha256-simd
# is better but apache v2
export GOFLAGS='-tags=stdsha256'

# make sure we get licenses at all (we're ignoring go-licenses' exit code)
if [ "x$(go-licenses csv storj.io/uplink-c --stderrthreshold=FATAL | grep -E ',MIT$')" == "x" ]; then
	echo "unexpected error"
	exit 1
fi

# TODO: get go-licenses working without silencing so many non-fatal errors
# TODO: fix go-licenses for storj.io/common
# MIT, ISC, and BSD-3-Clause are okay
# storj.io/common is okay
# github.com/spacemonkeygo/monkit/v3 is explicitly replaced
go-licenses csv storj.io/uplink-c --stderrthreshold=FATAL \
	| grep -vE ',(MIT|ISC|BSD-3-Clause)$' \
	| grep -vE '^storj.io/common/' \
	| grep -vE '^github.com/spacemonkeygo/monkit/v3,' \
	> .license-exceptions || true

if [ "x$(cat .license-exceptions)" != "x" ]; then
	echo "License exceptions:"
	cat .license-exceptions
	rm .license-exceptions
	exit 1
fi

rm .license-exceptions

exit 0
