// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

func main() {}

var universe = newHandles()

//export internal_UniverseIsEmpty
// internal_UniverseIsEmpty returns true if nothing is stored in the global map.
func internal_UniverseIsEmpty() bool {
	return universe.Empty()
}
