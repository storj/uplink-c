// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"unsafe"

	"storj.io/uplink"
)

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
func parse_access(accessString *C.char) C.AccessResult {
	access, err := uplink.ParseAccess(C.GoString(accessString))
	if err != nil {
		return C.AccessResult{
			error: mallocError(err),
		}
	}

	return C.AccessResult{
		access: (*C.Access)(mallocHandle(universe.Add(&Access{access}))),
	}
}

//export access_serialize
// access_serialize serializes access into a string.
func access_serialize(access *C.Access) C.StringResult {
	if access == nil {
		return C.StringResult{
			error: mallocError(ErrNull.New("access")),
		}
	}

	acc, ok := universe.Get(access._handle).(*Access)
	if !ok {
		return C.StringResult{
			error: mallocError(ErrInvalidHandle.New("access")),
		}
	}

	str, err := acc.Serialize()
	if err != nil {
		return C.StringResult{
			error: mallocError(err),
		}
	}
	return C.StringResult{
		string: C.CString(str),
	}
}

//export free_string_result
// free_string_result frees the resources associated with Access.
func free_string_result(result C.StringResult) {
	free_error(result.error)
	C.free(unsafe.Pointer(result.string))
}

//export free_access_result
// free_access_result frees the resources associated with Access.
func free_access_result(result C.AccessResult) {
	free_error(result.error)
	free_access(result.access)
}

//export free_access
// free_access frees the resources associated with Access.
func free_access(access *C.Access) {
	if access == nil {
		return
	}

	// TODO: should this return an error?
	universe.Del(access._handle)
}
