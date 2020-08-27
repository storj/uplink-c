// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"reflect"
	"unsafe"

	"storj.io/uplink"
)

// Download is a partial download to Storj Network.
type Download struct {
	scope
	download *uplink.Download
}

//export uplink_download_object
// uplink_download_object starts  download to the specified key.
func uplink_download_object(project *C.UplinkProject, bucket_name, object_key *C.uplink_const_char, options *C.UplinkDownloadOptions) C.UplinkDownloadResult { //nolint:golint
	if project == nil {
		return C.UplinkDownloadResult{
			error: mallocError(ErrNull.New("project")),
		}
	}
	if bucket_name == nil {
		return C.UplinkDownloadResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}
	if object_key == nil {
		return C.UplinkDownloadResult{
			error: mallocError(ErrNull.New("object_key")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.UplinkDownloadResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}
	scope := proj.scope.child()

	opts := &uplink.DownloadOptions{
		Offset: 0,
		Length: -1,
	}
	if options != nil {
		opts.Offset = int64(options.offset)
		opts.Length = int64(options.length)
	}

	download, err := proj.DownloadObject(scope.ctx, C.GoString(bucket_name), C.GoString(object_key), opts)
	if err != nil {
		return C.UplinkDownloadResult{
			error: mallocError(err),
		}
	}

	return C.UplinkDownloadResult{
		download: (*C.UplinkDownload)(mallocHandle(universe.Add(&Download{scope, download}))),
	}
}

//export uplink_download_read
// uplink_download_read downloads from object's data stream into bytes up to length amount.
// It returns the number of bytes read (0 <= bytes_read <= length) and
// any error encountered that caused the read to stop early.
func uplink_download_read(download *C.UplinkDownload, bytes unsafe.Pointer, length C.size_t) C.UplinkReadResult {
	down, ok := universe.Get(download._handle).(*Download)
	if !ok {
		return C.UplinkReadResult{
			error: mallocError(ErrInvalidHandle.New("download")),
		}
	}

	ilength, ok := safeConvertToInt(length)
	if !ok {
		return C.UplinkReadResult{
			error: mallocError(ErrInvalidArg.New("length too large")),
		}
	}

	var buf []byte
	hbuf := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	hbuf.Data = uintptr(bytes)
	hbuf.Len = ilength
	hbuf.Cap = ilength

	n, err := down.download.Read(buf)
	return C.UplinkReadResult{
		bytes_read: C.size_t(n),
		error:      mallocError(err),
	}
}

//export uplink_download_info
// uplink_download_info returns information about the downloaded object.
func uplink_download_info(download *C.UplinkDownload) C.UplinkObjectResult {
	down, ok := universe.Get(download._handle).(*Download)
	if !ok {
		return C.UplinkObjectResult{
			error: mallocError(ErrInvalidHandle.New("download")),
		}
	}

	info := down.download.Info()
	return C.UplinkObjectResult{
		object: mallocObject(info),
	}
}

//export uplink_free_read_result
// uplink_free_read_result frees any resources associated with read result.
func uplink_free_read_result(result C.UplinkReadResult) {
	uplink_free_error(result.error)
}

//export uplink_close_download
// uplink_close_download closes the download.
func uplink_close_download(download *C.UplinkDownload) *C.UplinkError {
	if download == nil {
		return nil
	}

	down, ok := universe.Get(download._handle).(*Download)
	if !ok {
		return mallocError(ErrInvalidHandle.New("download"))
	}

	return mallocError(down.download.Close())
}

//export uplink_free_download_result
// uplink_free_download_result frees any associated resources.
func uplink_free_download_result(result C.UplinkDownloadResult) {
	uplink_free_error(result.error)
	freeDownload(result.download)
}

// freeDownload closes the download and frees any associated resources.
func freeDownload(download *C.UplinkDownload) {
	if download == nil {
		return
	}
	defer C.free(unsafe.Pointer(download))
	defer universe.Del(download._handle)

	down, ok := universe.Get(download._handle).(*Download)
	if !ok {
		return
	}

	down.cancel()
	// in case we haven't already closed the download
	_ = down.download.Close()
	// TODO: log error when we didn't close manually and the close returns an error
}
