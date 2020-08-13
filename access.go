// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"context"
	"reflect"
	"time"
	"unsafe"

	"storj.io/uplink"
)

// Access grant contains everything to access a project and specific buckets.
type Access struct {
	*uplink.Access
}

//export parse_access
// parse_access parses serialized access grant string.
func parse_access(accessString *C.const_char) C.AccessResult { //nolint:golint
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
// request_access_with_passphrase requests satellite for a new access grant using a passhprase.
func request_access_with_passphrase(satellite_address, api_key, passphrase *C.const_char) C.AccessResult { //nolint:golint
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
// access_serialize serializes access grant into a string.
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
// access_share creates new access grant with specific permission. Permission will be applied to prefixes when defined.
func access_share(access *C.Access, permission C.Permission, prefixes *C.SharePrefix, prefixes_count int) C.AccessResult { //nolint:golint
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
		AllowDownload: bool(permission.allow_download),
		AllowUpload:   bool(permission.allow_upload),
		AllowList:     bool(permission.allow_list),
		AllowDelete:   bool(permission.allow_delete),
	}

	if permission.not_before != 0 {
		perm.NotBefore = time.Unix(int64(permission.not_before), 0)
	}
	if permission.not_after != 0 {
		perm.NotAfter = time.Unix(int64(permission.not_after), 0)
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

//export access_override_encryption_key
// access_override_encryption_key overrides the root encryption key for the prefix in
// bucket with encryptionKey.
//
// This function is useful for overriding the encryption key in user-specific
// access grants when implementing multitenancy in a single app bucket.
func access_override_encryption_key(access *C.Access, bucket, prefix *C.const_char, encryptionKey *C.EncryptionKey) *C.Error { //nolint:golint
	if access == nil {
		return mallocError(ErrNull.New("access"))
	}

	acc, ok := universe.Get(access._handle).(*Access)
	if !ok {
		return mallocError(ErrInvalidHandle.New("access"))
	}

	if encryptionKey == nil {
		return mallocError(ErrNull.New("encryption key"))
	}

	encKey, ok := universe.Get(encryptionKey._handle).(*EncryptionKey)
	if !ok {
		return mallocError(ErrInvalidHandle.New("encryption key"))
	}

	err := acc.OverrideEncryptionKey(C.GoString(bucket), C.GoString(prefix), encKey.EncryptionKey)
	return mallocError(err)
}

//export free_string_result
// free_string_result frees the resources associated with string result.
func free_string_result(result C.StringResult) {
	free_error(result.error)
	C.free(unsafe.Pointer(result.string))
}

//export free_access_result
// free_access_result frees the resources associated with access grant.
func free_access_result(result C.AccessResult) {
	free_error(result.error)
	freeAccess(result.access)
}

func freeAccess(access *C.Access) {
	if access == nil {
		return
	}
	defer C.free(unsafe.Pointer(access))
	defer universe.Del(access._handle)
}
