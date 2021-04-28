// Copyright (C) 2021 Storj Labs, Inc.
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

// PartUpload is an upload for a part.
type PartUpload struct {
	scope
	partUpload *uplink.PartUpload
}

//export uplink_upload_part
// uplink_upload_part starts an part upload to the specified key nad part number.
func uplink_upload_part(project *C.UplinkProject, bucket_name, object_key, upload_id *C.uplink_const_char, part_number C.uint32_t) C.UplinkPartUploadResult { //nolint:golint
	if project == nil {
		return C.UplinkPartUploadResult{
			error: mallocError(ErrNull.New("project")),
		}
	}
	if bucket_name == nil {
		return C.UplinkPartUploadResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}
	if object_key == nil {
		return C.UplinkPartUploadResult{
			error: mallocError(ErrNull.New("object_key")),
		}
	}
	if upload_id == nil {
		return C.UplinkPartUploadResult{
			error: mallocError(ErrNull.New("upload_id")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.UplinkPartUploadResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}

	scope := proj.scope.child()
	partUpload, err := proj.UploadPart(scope.ctx, C.GoString(bucket_name), C.GoString(object_key), C.GoString(upload_id), uint32(part_number))
	return C.UplinkPartUploadResult{
		part_upload: (*C.UplinkPartUpload)(mallocHandle(universe.Add(&PartUpload{scope, partUpload}))),
		error:       mallocError(err),
	}
}

//export uplink_part_upload_write
// uplink_part_upload_write uploads len(p) bytes from p to the object's data stream.
// It returns the number of bytes written from p (0 <= n <= len(p)) and
// any error encountered that caused the write to stop early.
func uplink_part_upload_write(upload *C.UplinkPartUpload, bytes unsafe.Pointer, length C.size_t) C.UplinkWriteResult {
	up, ok := universe.Get(upload._handle).(*PartUpload)
	if !ok {
		return C.UplinkWriteResult{
			error: mallocError(ErrInvalidHandle.New("part upload")),
		}
	}

	ilength, ok := safeConvertToInt(length)
	if !ok {
		return C.UplinkWriteResult{
			error: mallocError(ErrInvalidArg.New("length too large")),
		}
	}

	var buf []byte
	hbuf := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	hbuf.Data = uintptr(bytes)
	hbuf.Len = ilength
	hbuf.Cap = ilength

	n, err := up.partUpload.Write(buf)
	return C.UplinkWriteResult{
		bytes_written: C.size_t(n),
		error:         mallocError(err),
	}
}

//export uplink_part_upload_commit
// uplink_part_upload_commit commits the uploaded part data.
func uplink_part_upload_commit(upload *C.UplinkPartUpload) *C.UplinkError {
	up, ok := universe.Get(upload._handle).(*PartUpload)
	if !ok {
		return mallocError(ErrInvalidHandle.New("part upload"))
	}

	err := up.partUpload.Commit()
	return mallocError(err)
}

//export uplink_part_upload_abort
// uplink_part_upload_abort aborts a part upload.
func uplink_part_upload_abort(upload *C.UplinkPartUpload) *C.UplinkError {
	up, ok := universe.Get(upload._handle).(*PartUpload)
	if !ok {
		return mallocError(ErrInvalidHandle.New("part upload"))
	}

	err := up.partUpload.Abort()
	return mallocError(err)
}

//export uplink_part_upload_set_etag
// uplink_part_upload_set_etag sets part ETag.
func uplink_part_upload_set_etag(upload *C.UplinkPartUpload, etag *C.uplink_const_char) *C.UplinkError {

	up, ok := universe.Get(upload._handle).(*PartUpload)
	if !ok {
		return mallocError(ErrInvalidHandle.New("part upload"))
	}

	err := up.partUpload.SetETag([]byte(C.GoString(etag)))
	return mallocError(err)
}

//export uplink_part_upload_info
// uplink_part_upload_info returns the last information about the uploaded part.
func uplink_part_upload_info(upload *C.UplinkPartUpload) C.UplinkPartResult {
	up, ok := universe.Get(upload._handle).(*PartUpload)
	if !ok {
		return C.UplinkPartResult{
			error: mallocError(ErrInvalidHandle.New("part upload")),
		}
	}

	info := up.partUpload.Info()
	return C.UplinkPartResult{
		part: mallocPart(info),
	}
}

func mallocPart(part *uplink.Part) *C.UplinkPart {
	if part == nil {
		return nil
	}

	cpart := (*C.UplinkPart)(C.calloc(C.sizeof_UplinkPart, 1))
	*cpart = partToC(part)
	return cpart
}

func partToC(part *uplink.Part) C.UplinkPart {
	if part == nil {
		return C.UplinkPart{}
	}
	return C.UplinkPart{
		part_number: C.uint32_t(part.PartNumber),
		size:        C.size_t(part.Size),
		modified:    timeToUnix(part.Modified),
		etag:        C.CString(string(part.ETag)),
		etag_length: C.size_t(len(part.ETag)),
	}
}

//export uplink_free_part_result
// uplink_free_part_result frees memory associated with the part result.
func uplink_free_part_result(result C.UplinkPartResult) {
	uplink_free_error(result.error)
	uplink_free_part(result.part)
}

//export uplink_free_part_upload_result
// uplink_free_part_upload_result frees memory associated with the part upload result.
func uplink_free_part_upload_result(result C.UplinkPartUploadResult) {
	uplink_free_error(result.error)
	freePartUpload(result.part_upload)
}

func freePartUpload(partUpload *C.UplinkPartUpload) {
	if partUpload == nil {
		return
	}
	defer C.free(unsafe.Pointer(partUpload))
	defer universe.Del(partUpload._handle)

	up, ok := universe.Get(partUpload._handle).(*PartUpload)
	if ok {
		up.cancel()
	}
}

//export uplink_free_part
// uplink_free_part frees memory associated with the Part.
func uplink_free_part(part *C.UplinkPart) {
	if part == nil {
		return
	}
	defer C.free(unsafe.Pointer(part))

	if part.etag != nil {
		C.free(unsafe.Pointer(part.etag))
	}
}
