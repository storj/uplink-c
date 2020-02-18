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

// Download is a partial download to Storj Network.
type Download struct {
	scope
	download *uplink.Download
}

//export download_object
// download_object starts  download to the specified key.
func download_object(project C.Project, bucket_name, object_key *C.char, cerr **C.char) C.Download {
	if bucket_name == nil {
		*cerr = C.CString("bucket_name == nil")
		return C.Download{}
	}
	if object_key == nil {
		*cerr = C.CString("object_key == nil")
		return C.Download{}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		*cerr = C.CString("invalid project")
		return C.Download{}
	}
	scope := proj.scope.child()

	download, err := proj.DownloadObject(scope.ctx, C.GoString(bucket_name), C.GoString(object_key))
	if err != nil {
		*cerr = C.CString(fmt.Sprintf("%+v", err))
	}

	return C.Download{universe.Add(&Download{scope, download})}
}

//export download_read
// download_read uploads len(p) bytes from p to the object's data stream.
// It returns the number of bytes written from p (0 <= n <= len(p)) and
// any error encountered that caused the write to stop early.
func download_read(download C.Download, bytes *C.uint8_t, length C.size_t, cerr **C.char) C.size_t {
	down, ok := universe.Get(download._handle).(*Download)
	if !ok {
		*cerr = C.CString("invalid download")
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

	n, err := down.download.Read(buf)
	if err != nil {
		*cerr = C.CString(fmt.Sprintf("%+v", err))
	}
	return C.size_t(n)
}

//export download_info
// download_info returns information about the downloaded object.
func download_info(download C.Download, cerr **C.char) C.Object {
	down, ok := universe.Get(download._handle).(*Download)
	if !ok {
		*cerr = C.CString("invalid download")
		return C.Object{}
	}

	info := down.download.Info()
	return objectToC(info)
}

//export free_download
// free_download closes the download and frees any associated resources.
func free_download(download C.Download, cerr **C.char) {
	down, ok := universe.Get(download._handle).(*Download)
	if !ok {
		*cerr = C.CString("invalid download")
		return
	}

	universe.Del(download._handle)
	defer down.cancel()

	err := down.download.Close()
	if err != nil {
		*cerr = C.CString(fmt.Sprintf("%+v", err))
	}
}
