// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"time"
	"unsafe"

	"storj.io/uplink"
)

//export uplink_begin_upload
// uplink_begin_upload begins a new multipart upload to bucket and key.
func uplink_begin_upload(project *C.UplinkProject, bucket_name, object_key *C.uplink_const_char, options *C.UplinkUploadOptions) C.UplinkUploadInfoResult { //nolint:golint
	if project == nil {
		return C.UplinkUploadInfoResult{
			error: mallocError(ErrNull.New("project")),
		}
	}
	if bucket_name == nil {
		return C.UplinkUploadInfoResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}
	if object_key == nil {
		return C.UplinkUploadInfoResult{
			error: mallocError(ErrNull.New("object_key")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.UplinkUploadInfoResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}

	opts := &uplink.UploadOptions{}
	if options != nil {
		if options.expires > 0 {
			opts.Expires = time.Unix(int64(options.expires), 0)
		}
	}

	info, err := proj.BeginUpload(proj.scope.ctx, C.GoString(bucket_name), C.GoString(object_key), opts)
	return C.UplinkUploadInfoResult{
		error: mallocError(err),
		info:  mallocUploadInfo(&info),
	}
}

func mallocUploadInfo(info *uplink.UploadInfo) *C.UplinkUploadInfo {
	if info == nil {
		return nil
	}

	cinfo := (*C.UplinkUploadInfo)(C.calloc(C.sizeof_UplinkUploadInfo, 0))
	cinfo.upload_id = C.CString(info.UploadID)

	return cinfo
}

//export uplink_free_upload_info_result
// uplink_free_upload_info_result frees any resources associated with upload info result.
func uplink_free_upload_info_result(result C.UplinkUploadInfoResult) {
	uplink_free_error(result.error)
	uplink_free_upload_info(result.info)
}

//export uplink_free_upload_info
// uplink_free_upload_info frees memory associated with upload info.
func uplink_free_upload_info(info *C.UplinkUploadInfo) {
	if info == nil {
		return
	}
	defer C.free(unsafe.Pointer(info))

	if info.upload_id != nil {
		C.free(unsafe.Pointer(info.upload_id))
	}
}

//export uplink_commit_upload
// uplink_commit_upload commits a multipart upload to bucket and key started with uplink_begin_upload.
func uplink_commit_upload(project *C.UplinkProject, bucket_name, object_key, upload_id *C.uplink_const_char, options *C.UplinkCommitUploadOptions) C.UplinkCommitUploadResult { //nolint:golint
	if project == nil {
		return C.UplinkCommitUploadResult{
			error: mallocError(ErrNull.New("project")),
		}
	}
	if bucket_name == nil {
		return C.UplinkCommitUploadResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}
	if object_key == nil {
		return C.UplinkCommitUploadResult{
			error: mallocError(ErrNull.New("object_key")),
		}
	}
	if upload_id == nil {
		return C.UplinkCommitUploadResult{
			error: mallocError(ErrNull.New("upload_id")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.UplinkCommitUploadResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}

	opts := &uplink.CommitUploadOptions{}
	if options != nil {
		opts.CustomMetadata = customMetadataFromC(options.custom_metadata)
	}

	object, err := proj.CommitUpload(proj.scope.ctx, C.GoString(bucket_name), C.GoString(object_key), C.GoString(upload_id), opts)
	return C.UplinkCommitUploadResult{
		error:  mallocError(err),
		object: mallocObject(object),
	}
}

//export uplink_free_commit_upload_result
// uplink_free_commit_upload_result frees any resources associated with commit upload result.
func uplink_free_commit_upload_result(result C.UplinkCommitUploadResult) {
	uplink_free_error(result.error)
	uplink_free_object(result.object)
}

//export uplink_abort_upload
// uplink_abort_upload aborts a multipart upload started with uplink_begin_upload.
func uplink_abort_upload(project *C.UplinkProject, bucket_name, object_key, upload_id *C.uplink_const_char) *C.UplinkError { //nolint:golint
	if project == nil {
		return mallocError(ErrNull.New("project"))
	}
	if bucket_name == nil {
		return mallocError(ErrNull.New("bucket_name"))
	}
	if object_key == nil {
		return mallocError(ErrNull.New("object_key"))
	}
	if upload_id == nil {
		return mallocError(ErrNull.New("upload_id"))
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return mallocError(ErrInvalidHandle.New("project"))
	}

	err := proj.AbortUpload(proj.scope.ctx, C.GoString(bucket_name), C.GoString(object_key), C.GoString(upload_id))
	return mallocError(err)
}
