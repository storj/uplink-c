// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"unsafe"

	"storj.io/uplink"
)

//export stat_object
// stat_object returns information about an object at the specific key.
func stat_object(project *C.Project, bucket_name, object_key *C.char) C.ObjectResult {
	if project == nil {
		return C.ObjectResult{
			error: mallocError(ErrNull.New("project")),
		}
	}
	if bucket_name == nil {
		return C.ObjectResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}
	if object_key == nil {
		return C.ObjectResult{
			error: mallocError(ErrNull.New("object_key")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.ObjectResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}

	object, err := proj.StatObject(proj.scope.ctx, C.GoString(bucket_name), C.GoString(object_key))
	return C.ObjectResult{
		error:  mallocError(err),
		object: mallocObject(object),
	}
}

//export delete_object
// delete_object deletes an object.
func delete_object(project *C.Project, bucket_name, object_key *C.char) C.ObjectResult {
	if project == nil {
		return C.ObjectResult{
			error: mallocError(ErrNull.New("project")),
		}
	}
	if bucket_name == nil {
		return C.ObjectResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}
	if object_key == nil {
		return C.ObjectResult{
			error: mallocError(ErrNull.New("object_key")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.ObjectResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}

	deleted, err := proj.DeleteObject(proj.scope.ctx, C.GoString(bucket_name), C.GoString(object_key))
	return C.ObjectResult{
		error:  mallocError(err),
		object: mallocObject(deleted),
	}
}

func mallocObject(object *uplink.Object) *C.Object {
	if object == nil {
		return nil
	}

	cobject := (*C.Object)(C.calloc(C.sizeof_Object, 1))
	*cobject = objectToC(object)
	return cobject
}

func objectToC(object *uplink.Object) C.Object {
	if object == nil {
		return C.Object{}
	}
	return C.Object{
		key:       C.CString(object.Key),
		is_prefix: C.bool(object.IsPrefix),
		system: C.SystemMetadata{
			created: timeToUnix(object.System.Created),
			expires: timeToUnix(object.System.Expires),
		},
		custom: customMetadataToC(object.Custom),
	}
}

//export free_object_result
// free_object_result frees memory associated with the ObjectResult.
func free_object_result(obj C.ObjectResult) {
	free_error(obj.error)
	free_object(obj.object)
}

//export free_object
// free_object frees memory associated with the Object.
func free_object(obj *C.Object) {
	if obj == nil {
		return
	}
	defer C.free(unsafe.Pointer(obj))

	if obj.key != nil {
		C.free(unsafe.Pointer(obj.key))
	}

	free_system_metadata(&obj.system)
	free_custom_metadata_data(&obj.custom)
}

func free_system_metadata(system *C.SystemMetadata) {
}

func bytesToC(data []byte) C.Bytes {
	if len(data) == 0 {
		return C.Bytes{}
	}

	return C.Bytes{
		data:   C.CBytes(data),
		length: C.uint64_t(len(data)),
	}
}

func free_bytes_data(bytes *C.Bytes) {
	if bytes.data != nil {
		C.free(bytes.data)
		bytes.data = nil
	}
	bytes.length = 0
}
