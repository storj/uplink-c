// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"context"
	"reflect"
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

//export request_access_with_passphrase
// request_access_with_passphrase requests satellite for a new access using a passhprase.
func request_access_with_passphrase(satellite_address, api_key, passphrase *C.char) C.AccessResult {
	if satellite_address == nil {
		return C.AccessResult{
			error: mallocError(ErrNull.New("satellite_address")),
		}
	}
	if api_key == nil {
		return C.AccessResult{
			error: mallocError(ErrNull.New("api_key")),
		}
	}
	if passphrase == nil {
		return C.AccessResult{
			error: mallocError(ErrNull.New("passphrase")),
		}
	}

	ctx := context.Background()
	access, err := uplink.RequestAccessWithPassphrase(ctx, C.GoString(satellite_address), C.GoString(api_key), C.GoString(passphrase))
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

//export access_share
// access_share creates new Access with specific permission. Permission will be applied to prefixes when defined.
func access_share(access *C.Access, permission C.Permission, prefixes *C.SharePrefix, prefixes_count int) C.AccessResult {
	if access == nil {
		return C.AccessResult{
			error: mallocError(ErrNull.New("access")),
		}
	}

	acc, ok := universe.Get(access._handle).(*Access)
	if !ok {
		return C.AccessResult{
			error: mallocError(ErrInvalidHandle.New("access")),
		}
	}

	perm := uplink.Permission{
		AllowRead:   bool(permission.allow_read),
		AllowWrite:  bool(permission.allow_write),
		AllowList:   bool(permission.allow_list),
		AllowDelete: bool(permission.allow_delete),

		// TODO: not before and not after
	}

	var goprefixes []uplink.SharePrefix
	if prefixes != nil && prefixes_count > 0 {
		var array []C.SharePrefix
		*(*reflect.SliceHeader)(unsafe.Pointer(&array)) = reflect.SliceHeader{
			Data: uintptr(unsafe.Pointer(prefixes)),
			Len:  prefixes_count,
			Cap:  prefixes_count,
		}

		for _, p := range array {
			goprefixes = append(goprefixes, uplink.SharePrefix{
				Bucket: C.GoString(p.bucket),
				Prefix: C.GoString(p.prefix),
			})
		}
	}

	newAccess, err := acc.Share(perm, goprefixes...)
	if err != nil {
		return C.AccessResult{
			error: mallocError(err),
		}
	}
	return C.AccessResult{
		access: (*C.Access)(mallocHandle(universe.Add(&Access{newAccess}))),
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
	defer C.free(unsafe.Pointer(access))
	defer universe.Del(access._handle)

	// TODO: should this return an error?
}
