// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"unsafe"

	"storj.io/uplink"
)

// uplink_stat_object returns information about an object at the specific key.
//
//export uplink_stat_object
func uplink_stat_object(project *C.UplinkProject, bucket_name, object_key *C.uplink_const_char) C.UplinkObjectResult { //nolint:golint
	if project == nil {
		return C.UplinkObjectResult{
			error: mallocError(ErrNull.New("project")),
		}
	}
	if bucket_name == nil {
		return C.UplinkObjectResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}
	if object_key == nil {
		return C.UplinkObjectResult{
			error: mallocError(ErrNull.New("object_key")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.UplinkObjectResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}

	object, err := proj.StatObject(proj.scope.ctx, C.GoString(bucket_name), C.GoString(object_key))
	return C.UplinkObjectResult{
		error:  mallocError(err),
		object: mallocObject(object),
	}
}

// uplink_delete_object deletes an object.
//
//export uplink_delete_object
func uplink_delete_object(project *C.UplinkProject, bucket_name, object_key *C.uplink_const_char) C.UplinkObjectResult { //nolint:golint
	if project == nil {
		return C.UplinkObjectResult{
			error: mallocError(ErrNull.New("project")),
		}
	}
	if bucket_name == nil {
		return C.UplinkObjectResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}
	if object_key == nil {
		return C.UplinkObjectResult{
			error: mallocError(ErrNull.New("object_key")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.UplinkObjectResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}

	deleted, err := proj.DeleteObject(proj.scope.ctx, C.GoString(bucket_name), C.GoString(object_key))
	return C.UplinkObjectResult{
		error:  mallocError(err),
		object: mallocObject(deleted),
	}
}

func mallocObject(object *uplink.Object) *C.UplinkObject {
	if object == nil {
		return nil
	}

	cobject := (*C.UplinkObject)(calloc(1, C.sizeof_UplinkObject))
	*cobject = objectToC(object)
	return cobject
}

func objectToC(object *uplink.Object) C.UplinkObject {
	if object == nil {
		return C.UplinkObject{}
	}
	return C.UplinkObject{
		key:       C.CString(object.Key),
		is_prefix: C.bool(object.IsPrefix),
		system: C.UplinkSystemMetadata{
			created:        timeToUnix(object.System.Created),
			expires:        timeToUnix(object.System.Expires),
			content_length: C.int64_t(object.System.ContentLength),
		},
		custom: customMetadataToC(object.Custom),
	}
}

// uplink_free_object_result frees memory associated with the ObjectResult.
//
//export uplink_free_object_result
func uplink_free_object_result(obj C.UplinkObjectResult) {
	uplink_free_error(obj.error)
	uplink_free_object(obj.object)
}

// uplink_free_object frees memory associated with the Object.
//
//export uplink_free_object
func uplink_free_object(obj *C.UplinkObject) {
	if obj == nil {
		return
	}
	defer C.free(unsafe.Pointer(obj))

	if obj.key != nil {
		C.free(unsafe.Pointer(obj.key))
		obj.key = nil
	}

	freeSystemMetadata(&obj.system)
	freeCustomMetadataData(&obj.custom)
}

// uplink_update_object_metadata replaces the custom metadata for the object at the specific key with new_metadata.
// Any existing custom metadata will be deleted.
//
//export uplink_update_object_metadata
func uplink_update_object_metadata(project *C.UplinkProject, bucket_name, object_key *C.uplink_const_char, new_metadata C.UplinkCustomMetadata, options *C.UplinkUploadObjectMetadataOptions) *C.UplinkError { //nolint:golint
	if project == nil {
		return mallocError(ErrNull.New("project"))
	}

	if bucket_name == nil {
		return mallocError(ErrNull.New("bucket_name"))
	}

	if object_key == nil {
		return mallocError(ErrNull.New("object_key"))
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return mallocError(ErrInvalidHandle.New("project"))
	}

	err := proj.UpdateObjectMetadata(proj.scope.ctx, C.GoString(bucket_name), C.GoString(object_key), customMetadataFromC(new_metadata), nil)
	return mallocError(err)
}

func freeSystemMetadata(system *C.UplinkSystemMetadata) {
}
