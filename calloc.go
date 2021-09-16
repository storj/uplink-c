// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include <stdlib.h>
import "C"

import "unsafe"

//go:linkname calloc_runtime_throw runtime.throw
func calloc_runtime_throw(string)

func calloc(nitems C.size_t, size C.size_t) unsafe.Pointer {
	ptr := C.calloc(nitems, size)

	if ptr == nil {
		// if requested a zero-sized allocation and a nil pointer
		// is returned, instead malloc a byte to match the behavior
		// of the Go provided malloc wrapper in that it never
		// returns a nil pointer.
		if nitems == 0 || size == 0 {
			return C.malloc(1)
		}
		calloc_runtime_throw("runtime: C calloc failed")
		panic("unreachable")
	}

	return ptr
}
