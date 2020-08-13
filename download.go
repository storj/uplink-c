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

//export download_object
// download_object starts  download to the specified key.
func download_object(project *C.Project, bucket_name, object_key *C.const_char, options *C.DownloadOptions) C.DownloadResult { //nolint:golint
	if project == nil {
		return C.DownloadResult{
			error: mallocError(ErrNull.New("project")),
		}
	}
	if bucket_name == nil {
		return C.DownloadResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}
	if object_key == nil {
		return C.DownloadResult{
			error: mallocError(ErrNull.New("object_key")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.DownloadResult{
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
		return C.DownloadResult{
			error: mallocError(err),
		}
	}

	return C.DownloadResult{
		download: (*C.Download)(mallocHandle(universe.Add(&Download{scope, download}))),
	}
}

//export download_read
// download_read downloads from object's data stream into bytes up to length amount.
// It returns the number of bytes read (0 <= bytes_read <= length) and
// any error encountered that caused the read to stop early.
func download_read(download *C.Download, bytes unsafe.Pointer, length C.size_t) C.ReadResult {
	down, ok := universe.Get(download._handle).(*Download)
	if !ok {
		return C.ReadResult{
			error: mallocError(ErrInvalidHandle.New("download")),
		}
	}

	ilength, ok := safeConvertToInt(length)
	if !ok {
		return C.ReadResult{
			error: mallocError(ErrInvalidArg.New("length too large")),
		}
	}

	var buf []byte
	*(*reflect.SliceHeader)(unsafe.Pointer(&buf)) = reflect.SliceHeader{
		Data: uintptr(bytes),
		Len:  ilength,
		Cap:  ilength,
	}

	n, err := down.download.Read(buf)
	return C.ReadResult{
		bytes_read: C.size_t(n),
		error:      mallocError(err),
	}
}

//export download_info
// download_info returns information about the downloaded object.
func download_info(download *C.Download) C.ObjectResult {
	down, ok := universe.Get(download._handle).(*Download)
	if !ok {
		return C.ObjectResult{
			error: mallocError(ErrInvalidHandle.New("download")),
		}
	}

	info := down.download.Info()
	return C.ObjectResult{
		object: mallocObject(info),
	}
}

//export free_read_result
// free_read_result frees any resources associated with read result.
func free_read_result(result C.ReadResult) {
	free_error(result.error)
}

//export close_download
// close_download closes the download.
func close_download(download *C.Download) *C.Error {
	if download == nil {
		return nil
	}

	down, ok := universe.Get(download._handle).(*Download)
	if !ok {
		return mallocError(ErrInvalidHandle.New("download"))
	}

	return mallocError(down.download.Close())
}

//export free_download_result
// free_download_result frees any associated resources.
func free_download_result(result C.DownloadResult) {
	free_error(result.error)
	freeDownload(result.download)
}

// freeDownload closes the download and frees any associated resources.
func freeDownload(download *C.Download) {
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
