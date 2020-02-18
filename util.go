// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import "C"

// safeConvertToInt converts the C.size_t to an int, and returns a boolean
// indicating if the conversion was lossless and semantically equivalent.
func safeConvertToInt(n C.size_t) (int, bool) {
	return int(n), C.size_t(int(n)) == n && int(n) >= 0
}
