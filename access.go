// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import "storj.io/uplink"

// Access contains everything to access a project
// and specific buckets.
type Access struct {
	*uplink.Access
}

//export parse_access
// parse_access parses access string.
//
// For convenience with using other arguments,
// parse does not return an error. But, instead
// delays the calls.
func parse_access(accessString *C.char) C.Access {
	access := uplink.ParseAccess(C.GoString(accessString))
	return C.Access{universe.Add(&Access{access})}
}

//export free_access
// free_access frees the resources associated with Access.
func free_access(access C.Access) {
	universe.Del(access._handle)
}
