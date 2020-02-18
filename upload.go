// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"

	"storj.io/uplink"
)

// Upload is a partial upload to Storj Network.
type Upload struct {
	scope
	upload *uplink.Upload
}

//export upload_object
// upload_object starts an upload to the specified key.
func upload_object(project C.Project, bucket_name, object_key *C.char, cerr **C.char) C.Upload {
	if bucket_name == nil {
		*cerr = C.CString("bucket_name == nil")
		return C.Upload{}
	}
	if object_key == nil {
		*cerr = C.CString("object_key == nil")
		return C.Upload{}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		*cerr = C.CString("invalid project")
		return C.Upload{}
	}
	scope := proj.scope.child()

	upload, err := proj.UploadObject(scope.ctx, C.GoString(bucket_name), C.GoString(object_key))
	if err != nil {
		*cerr = C.CString(fmt.Sprintf("%+v", err))
	}

	return C.Upload{universe.Add(&Upload{scope, upload})}
}

//export upload_write
// upload_write uploads len(p) bytes from p to the object's data stream.
// It returns the number of bytes written from p (0 <= n <= len(p)) and
// any error encountered that caused the write to stop early.
func upload_write(upload C.Upload, bytes *C.uint8_t, length C.size_t, cerr **C.char) C.size_t {
	up, ok := universe.Get(upload._handle).(*Upload)
	if !ok {
		*cerr = C.CString("invalid upload")
		return C.size_t(0)
	}

	ilength, ok := safeConvertToInt(length)
	if !ok {
		*cerr = C.CString("invalid length: too large or negative")
		return C.size_t(0)
	}

	var buf []byte
	*(*reflect.SliceHeader)(unsafe.Pointer(&buf)) = reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(bytes)),
		Len:  ilength,
		Cap:  ilength,
	}

	n, err := up.upload.Write(buf)
	if err != nil {
		*cerr = C.CString(fmt.Sprintf("%+v", err))
	}
	return C.size_t(n)
}

//export upload_commit
// upload_commit commits the uploaded data.
func upload_commit(upload C.Upload, cerr **C.char) {
	up, ok := universe.Get(upload._handle).(*Upload)
	if !ok {
		*cerr = C.CString("invalid upload")
		return
	}

	err := up.upload.Commit()
	if err != nil {
		*cerr = C.CString(fmt.Sprintf("%+v", err))
	}
}

//export upload_abort
// upload_abort aborts an upload.
func upload_abort(upload C.Upload, cerr **C.char) {
	up, ok := universe.Get(upload._handle).(*Upload)
	if !ok {
		*cerr = C.CString("invalid upload")
		return
	}

	err := up.upload.Abort()
	if err != nil {
		*cerr = C.CString(fmt.Sprintf("%+v", err))
	}
}

//export upload_info
// upload_info returns the last information about the uploaded object.
func upload_info(upload C.Upload, cerr **C.char) C.Object {
	up, ok := universe.Get(upload._handle).(*Upload)
	if !ok {
		*cerr = C.CString("invalid upload")
		return C.Object{}
	}

	info := up.upload.Info()
	return objectToC(info)
}

//export free_upload
// free_upload closes the upload and frees any associated resources.
func free_upload(upload C.Upload, cerr **C.char) {
	up, ok := universe.Get(upload._handle).(*Upload)
	if !ok {
		*cerr = C.CString("invalid upload")
		return
	}

	universe.Del(upload._handle)
	defer up.cancel()
}
