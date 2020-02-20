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
func download_object(project *C.Project, bucket_name, object_key *C.char) C.DownloadResult {
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

	download, err := proj.DownloadObject(scope.ctx, C.GoString(bucket_name), C.GoString(object_key))
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
// download_read uploads len(p) bytes from p to the object's data stream.
// It returns the number of bytes written from p (0 <= n <= len(p)) and
// any error encountered that caused the write to stop early.
func download_read(download *C.Download, bytes *C.uint8_t, length C.size_t) C.ReadResult {
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
		Data: uintptr(unsafe.Pointer(bytes)),
		Len:  ilength,
		Cap:  ilength,
	}

	n, err := down.download.Read(buf)
	return C.ReadResult{
		bytes_read: C.size_t(n),
		error:      mallocError(err),
	}
}

// TODO: should we have free_read_result?

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

//export free_download_result
// free_download_result closes the download and frees any associated resources.
func free_download_result(result C.DownloadResult) *C.Error {
	free_error(result.error)
	return free_download(result.download)
}

//export free_download
// free_download closes the download and frees any associated resources.
func free_download(download *C.Download) *C.Error {
	if download == nil {
		return nil
	}

	down, ok := universe.Get(download._handle).(*Download)
	if !ok {
		return mallocError(ErrInvalidHandle.New("download"))
	}

	universe.Del(download._handle)
	defer down.cancel()

	return mallocError(down.download.Close())
}