// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"

//export uplink_copy_object
// uplink_copy_object copies object to a same/different bucket and key.
func uplink_copy_object(project *C.UplinkProject, old_bucket_name, old_object_key, new_bucket_name, new_object_key *C.uplink_const_char,
	options *C.UplinkMoveObjectOptions) *C.UplinkError { //nolint:golint
	if project == nil {
		return mallocError(ErrNull.New("project"))
	}
	if old_bucket_name == nil {
		return mallocError(ErrNull.New("old_bucket_name"))
	}
	if old_object_key == nil {
		return mallocError(ErrNull.New("old_object_key"))
	}
	if new_bucket_name == nil {
		return mallocError(ErrNull.New("new_bucket_name"))
	}
	if new_object_key == nil {
		return mallocError(ErrNull.New("new_object_key"))
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return mallocError(ErrInvalidHandle.New("project"))
	}

	_, err := proj.CopyObject(proj.scope.ctx,
		C.GoString(old_bucket_name), C.GoString(old_object_key),
		C.GoString(new_bucket_name), C.GoString(new_object_key),
		nil)
	return mallocError(err)
}
