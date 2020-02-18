// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"fmt"
	"unsafe"

	"storj.io/uplink"
)

//export stat_object
// stat_object returns information about an object at the specific key.
func stat_object(project C.Project, bucket_name, object_key *C.char, cerr **C.char) C.Object {
	if bucket_name == nil {
		*cerr = C.CString("bucket_name == nil")
		return C.Object{}
	}
	if object_key == nil {
		*cerr = C.CString("object_key == nil")
		return C.Object{}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		*cerr = C.CString("invalid project")
		return C.Object{}
	}
	child := proj.scope.child()

	object, err := proj.StatObject(child.ctx, C.GoString(bucket_name), C.GoString(object_key))
	if err != nil {
		*cerr = C.CString(fmt.Sprintf("%+v", err))
	}

	return objectToC(object)
}

//export delete_object
// delete_object deletes an object.
func delete_object(project C.Project, bucket_name, object_key *C.char, cerr **C.char) {
	if bucket_name == nil {
		*cerr = C.CString("bucket_name == nil")
		return
	}
	if object_key == nil {
		*cerr = C.CString("object_key == nil")
		return
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		*cerr = C.CString("invalid project")
		return
	}
	child := proj.scope.child()

	err := proj.DeleteObject(child.ctx, C.GoString(bucket_name), C.GoString(object_key))
	if err != nil {
		*cerr = C.CString(fmt.Sprintf("%+v", err))
	}
}

func objectToC(object *uplink.Object) C.Object {
	if object == nil {
		return C.Object{}
	}
	return C.Object{
		key: C.CString(object.Key),
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

//export free_object
// free_object frees memory associated with the object.
func free_object(obj C.Object) {
	// TODO: figure out whether argument should be a pointer

	if obj.key != nil {
		C.free(unsafe.Pointer(obj.key))
		obj.key = nil
	}

	free_object_info(&obj.info)
	free_standard_metadata(&obj.standard)
	free_custom_metadata(&obj.custom)
}

func free_object_info(info *C.ObjectInfo) {
	info.created = 0
	info.expires = 0
}

func free_standard_metadata(standard *C.StandardMetadata) {
	standard.content_length = 0
	if standard.content_type != nil {
		C.free(unsafe.Pointer(standard.content_type))
		standard.content_type = nil
	}

	standard.file_created = 0
	standard.file_modified = 0
	standard.file_permissions = 0

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
