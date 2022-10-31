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

// uplink_parse_access parses serialized access grant string.
//
//export uplink_parse_access
func uplink_parse_access(accessString *C.uplink_const_char) C.UplinkAccessResult { //nolint:golint
	access, err := uplink.ParseAccess(C.GoString(accessString))
	if err != nil {
		return C.UplinkAccessResult{
			error: mallocError(err),
		}
	}

	return C.UplinkAccessResult{
		access: (*C.UplinkAccess)(mallocHandle(universe.Add(&Access{access}))),
	}
}

// uplink_request_access_with_passphrase requests satellite for a new access grant using a passhprase.
//
//export uplink_request_access_with_passphrase
func uplink_request_access_with_passphrase(satellite_address, api_key, passphrase *C.uplink_const_char) C.UplinkAccessResult { //nolint:golint
	if satellite_address == nil {
		return C.UplinkAccessResult{
			error: mallocError(ErrNull.New("satellite_address")),
		}
	}
	if api_key == nil {
		return C.UplinkAccessResult{
			error: mallocError(ErrNull.New("api_key")),
		}
	}
	if passphrase == nil {
		return C.UplinkAccessResult{
			error: mallocError(ErrNull.New("passphrase")),
		}
	}

	ctx := context.Background()
	access, err := uplink.RequestAccessWithPassphrase(ctx, C.GoString(satellite_address), C.GoString(api_key), C.GoString(passphrase))
	if err != nil {
		return C.UplinkAccessResult{
			error: mallocError(err),
		}
	}

	return C.UplinkAccessResult{
		access: (*C.UplinkAccess)(mallocHandle(universe.Add(&Access{access}))),
	}
}

// uplink_access_satellite_address returns the satellite node URL for this access grant.
//
//export uplink_access_satellite_address
func uplink_access_satellite_address(access *C.UplinkAccess) C.UplinkStringResult {
	if access == nil {
		return C.UplinkStringResult{
			error: mallocError(ErrNull.New("access")),
		}
	}

	acc, ok := universe.Get(access._handle).(*Access)
	if !ok {
		return C.UplinkStringResult{
			error: mallocError(ErrInvalidHandle.New("access")),
		}
	}

	return C.UplinkStringResult{
		string: C.CString(acc.Access.SatelliteAddress()),
	}
}

// uplink_access_serialize serializes access grant into a string.
//
//export uplink_access_serialize
func uplink_access_serialize(access *C.UplinkAccess) C.UplinkStringResult {
	if access == nil {
		return C.UplinkStringResult{
			error: mallocError(ErrNull.New("access")),
		}
	}

	acc, ok := universe.Get(access._handle).(*Access)
	if !ok {
		return C.UplinkStringResult{
			error: mallocError(ErrInvalidHandle.New("access")),
		}
	}

	str, err := acc.Serialize()
	if err != nil {
		return C.UplinkStringResult{
			error: mallocError(err),
		}
	}
	return C.UplinkStringResult{
		string: C.CString(str),
	}
}

// uplink_access_share creates new access grant with specific permission. Permission will be applied to prefixes when defined.
//
//export uplink_access_share
func uplink_access_share(access *C.UplinkAccess, permission C.UplinkPermission, prefixes *C.UplinkSharePrefix, prefixes_count int) C.UplinkAccessResult { //nolint:golint
	if access == nil {
		return C.UplinkAccessResult{
			error: mallocError(ErrNull.New("access")),
		}
	}

	acc, ok := universe.Get(access._handle).(*Access)
	if !ok {
		return C.UplinkAccessResult{
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
		var array []C.UplinkSharePrefix
		harray := (*reflect.SliceHeader)(unsafe.Pointer(&array))
		harray.Data = uintptr(unsafe.Pointer(prefixes))
		harray.Len = prefixes_count
		harray.Cap = prefixes_count

		for _, p := range array {
			goprefixes = append(goprefixes, uplink.SharePrefix{
				Bucket: C.GoString(p.bucket),
				Prefix: C.GoString(p.prefix),
			})
		}
	}

	newAccess, err := acc.Share(perm, goprefixes...)
	if err != nil {
		return C.UplinkAccessResult{
			error: mallocError(err),
		}
	}
	return C.UplinkAccessResult{
		access: (*C.UplinkAccess)(mallocHandle(universe.Add(&Access{newAccess}))),
	}
}

// uplink_access_override_encryption_key overrides the root encryption key for the prefix in
// bucket with encryptionKey.
//
// This function is useful for overriding the encryption key in user-specific
// access grants when implementing multitenancy in a single app bucket.
//
//export uplink_access_override_encryption_key
func uplink_access_override_encryption_key(access *C.UplinkAccess, bucket, prefix *C.uplink_const_char, encryptionKey *C.UplinkEncryptionKey) *C.UplinkError { //nolint:golint
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

// uplink_free_string_result frees the resources associated with string result.
//
//export uplink_free_string_result
func uplink_free_string_result(result C.UplinkStringResult) {
	uplink_free_error(result.error)
	C.free(unsafe.Pointer(result.string))
}

// uplink_free_access_result frees the resources associated with access grant.
//
//export uplink_free_access_result
func uplink_free_access_result(result C.UplinkAccessResult) {
	uplink_free_error(result.error)
	freeAccess(result.access)
}

func freeAccess(access *C.UplinkAccess) {
	if access == nil {
		return
	}
	defer C.free(unsafe.Pointer(access))
	defer universe.Del(access._handle)
}
