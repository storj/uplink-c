// Copyright (C) 2026 Storj Labs, Inc.
// See LICENSE for copying information.

//go:build uplink_coverage

package main

/*
// cgo concatenates the preamble into the exported header without
// wrapping it in include guards. If uplink.h is #included more than
// once in the same translation unit, the static function below would
// be redefined. Guard this block ourselves so the preamble is idempotent.
#ifndef UPLINK_COVERAGE_FLUSH_PREAMBLE
#define UPLINK_COVERAGE_FLUSH_PREAMBLE

#include <stdlib.h>

extern void uplinkFlushCoverage(void);

// __attribute__((unused)) suppresses -Wunused-function in C tests that
// include the generated uplink.h but don't call this function directly.
__attribute__((unused))
static void uplink_register_atexit_flush(void) {
	atexit(uplinkFlushCoverage);
}

#endif
*/
import "C"

import (
	"fmt"
	"os"
	"runtime/coverage"
)

// uplinkFlushCoverage writes the coverage meta and counter data to
// GOCOVERDIR. It is registered as an atexit handler so that host
// processes linking the c-shared library produce coverage data on exit
// without needing to call anything explicitly.
//
//export uplinkFlushCoverage
func uplinkFlushCoverage() {
	dir := os.Getenv("GOCOVERDIR")
	if dir == "" {
		return
	}
	if err := coverage.WriteMetaDir(dir); err != nil {
		fmt.Fprintln(os.Stderr, "uplink coverage: WriteMetaDir:", err)
	}
	if err := coverage.WriteCountersDir(dir); err != nil {
		fmt.Fprintln(os.Stderr, "uplink coverage: WriteCountersDir:", err)
	}
}

func init() {
	C.uplink_register_atexit_flush()
}
