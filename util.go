// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import "time"

// safeConvertToInt converts the C.size_t to an int, and returns a boolean
// indicating if the conversion was lossless and semantically equivalent.
func safeConvertToInt(n C.size_t) (int, bool) {
	return int(n), C.size_t(int(n)) == n && int(n) >= 0
}

// timeToUnix converts to C.int64_t.
func timeToUnix(t time.Time) C.int64_t {
	if t.IsZero() {
		return 0
	}
	return C.int64_t(t.Unix())
}
