// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"reflect"
	"time"
	"unsafe"

	"storj.io/uplink"
)

// Upload is a partial upload to Storj Network.
type Upload struct {
	scope
	upload *uplink.Upload
}

//export uplink_upload_object
// uplink_upload_object starts an upload to the specified key.
func uplink_upload_object(project *C.UplinkProject, bucket_name, object_key *C.uplink_const_char, options *C.UplinkUploadOptions) C.UplinkUploadResult { //nolint:golint
	if project == nil {
		return C.UplinkUploadResult{
			error: mallocError(ErrNull.New("project")),
		}
	}
	if bucket_name == nil {
		return C.UplinkUploadResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}
	if object_key == nil {
		return C.UplinkUploadResult{
			error: mallocError(ErrNull.New("object_key")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.UplinkUploadResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}
	scope := proj.scope.child()

	opts := &uplink.UploadOptions{}
	if options != nil {
		if options.expires > 0 {
			opts.Expires = time.Unix(int64(options.expires), 0)
		}
	}

	upload, err := proj.UploadObject(scope.ctx, C.GoString(bucket_name), C.GoString(object_key), opts)
	if err != nil {
		return C.UplinkUploadResult{
			error: mallocError(err),
		}
	}

	return C.UplinkUploadResult{
		upload: (*C.UplinkUpload)(mallocHandle(universe.Add(&Upload{scope, upload}))),
	}
}

//export uplink_upload_write
// uplink_upload_write uploads len(p) bytes from p to the object's data stream.
// It returns the number of bytes written from p (0 <= n <= len(p)) and
// any error encountered that caused the write to stop early.
func uplink_upload_write(upload *C.UplinkUpload, bytes unsafe.Pointer, length C.size_t) C.UplinkWriteResult {
	up, ok := universe.Get(upload._handle).(*Upload)
	if !ok {
		return C.UplinkWriteResult{
			error: mallocError(ErrInvalidHandle.New("upload")),
		}
	}

	ilength, ok := safeConvertToInt(length)
	if !ok {
		return C.UplinkWriteResult{
			error: mallocError(ErrInvalidArg.New("length too large")),
		}
	}

	var buf []byte
	*(*reflect.SliceHeader)(unsafe.Pointer(&buf)) = reflect.SliceHeader{
		Data: uintptr(bytes),
		Len:  ilength,
		Cap:  ilength,
	}

	n, err := up.upload.Write(buf)
	return C.UplinkWriteResult{
		bytes_written: C.size_t(n),
		error:         mallocError(err),
	}
}

//export uplink_upload_commit
// uplink_upload_commit commits the uploaded data.
func uplink_upload_commit(upload *C.UplinkUpload) *C.UplinkError {
	up, ok := universe.Get(upload._handle).(*Upload)
	if !ok {
		return mallocError(ErrInvalidHandle.New("upload"))
	}

	err := up.upload.Commit()
	return mallocError(err)
}

//export uplink_upload_abort
// uplink_upload_abort aborts an upload.
func uplink_upload_abort(upload *C.UplinkUpload) *C.UplinkError {
	up, ok := universe.Get(upload._handle).(*Upload)
	if !ok {
		return mallocError(ErrInvalidHandle.New("upload"))
	}

	err := up.upload.Abort()
	return mallocError(err)
}

//export uplink_upload_info
// uplink_upload_info returns the last information about the uploaded object.
func uplink_upload_info(upload *C.UplinkUpload) C.UplinkObjectResult {
	up, ok := universe.Get(upload._handle).(*Upload)
	if !ok {
		return C.UplinkObjectResult{
			error: mallocError(ErrInvalidHandle.New("upload")),
		}
	}

	info := up.upload.Info()
	return C.UplinkObjectResult{
		object: mallocObject(info),
	}
}

//export uplink_upload_set_custom_metadata
// uplink_upload_set_custom_metadata returns the last information about the uploaded object.
func uplink_upload_set_custom_metadata(upload *C.UplinkUpload, custom C.UplinkCustomMetadata) *C.UplinkError {
	up, ok := universe.Get(upload._handle).(*Upload)
	if !ok {
		return mallocError(ErrInvalidHandle.New("upload"))
	}

	customMetadata := customMetadataFromC(custom)
	err := up.upload.SetCustomMetadata(up.scope.ctx, customMetadata)

	return mallocError(err)
}

//export uplink_free_write_result
// uplink_free_write_result frees any resources associated with write result.
func uplink_free_write_result(result C.UplinkWriteResult) {
	uplink_free_error(result.error)
}

//export uplink_free_upload_result
// uplink_free_upload_result closes the upload and frees any associated resources.
func uplink_free_upload_result(result C.UplinkUploadResult) {
	uplink_free_error(result.error)
	freeUpload(result.upload)
}

func freeUpload(upload *C.UplinkUpload) {
	if upload == nil {
		return
	}
	defer C.free(unsafe.Pointer(upload))
	defer universe.Del(upload._handle)

	up, ok := universe.Get(upload._handle).(*Upload)
	if ok {
		up.cancel()
	}
}
