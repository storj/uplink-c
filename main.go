// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"

func main() {}

var universe = newHandles()

//export uplink_internal_UniverseIsEmpty
// uplink_internal_UniverseIsEmpty returns true if nothing is stored in the global map.
func uplink_internal_UniverseIsEmpty() bool {
	return universe.Empty()
}
