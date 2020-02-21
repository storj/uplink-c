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

	cobject := (*C.Object)(C.malloc(C.sizeof_Object))
	*cobject = objectToC(object)
	return cobject
}

func objectToC(object *uplink.Object) C.Object {
	if object == nil {
		return C.Object{}
	}
	return C.Object{
		key: C.CString(object.Key),
		is_prefix: C.bool(object.IsPrefix),
		info: C.ObjectInfo{
			created: timeToUnix(object.Info.Created),
			expires: timeToUnix(object.Info.Expires),
		},
		standard: C.StandardMetadata{
			content_length: C.int64_t(object.Standard.ContentLength),
			content_type:   C.CString(object.Standard.ContentType),

			file_created:     timeToUnix(object.Standard.FileCreated),
			file_modified:    timeToUnix(object.Standard.FileModified),
			file_permissions: C.uint32_t(object.Standard.FilePermissions),

			unknown: bytesToC(object.Standard.Unknown),
		},
		custom: C.CustomMetadata{
			// TODO:
		},
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

	if obj.key != nil {
		C.free(unsafe.Pointer(obj.key))
	}

	free_object_info(&obj.info)
	free_standard_metadata(&obj.standard)
	free_custom_metadata(&obj.custom)

	C.free(unsafe.Pointer(obj))
}

func free_object_info(info *C.ObjectInfo) {
}

func free_standard_metadata(standard *C.StandardMetadata) {
	standard.content_length = 0
	if standard.content_type != nil {
		C.free(unsafe.Pointer(standard.content_type))
	}

	free_bytes(&standard.unknown)
}

func free_custom_metadata(custom *C.CustomMetadata) {
	// TODO:
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

//export free_bytes
// free_bytes frees memory associated with bytes.
func free_bytes(bytes *C.Bytes) {
	if bytes.data != nil {
		C.free(bytes.data)
		bytes.data = nil
	}
	bytes.length = 0
}
