// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"context"
	"errors"
	"fmt"
	"io"
	"unsafe"

	"github.com/zeebo/errs"

	"storj.io/uplink"
)

//export SUCCESS
const SUCCESS C.uint32_t = 0

//export ERROR_EOF
const ERROR_EOF C.uint32_t = 1

//export ERROR_INTERNAL
const ERROR_INTERNAL C.uint32_t = 2

//export ERROR_CANCELED
const ERROR_CANCELED C.uint32_t = 3

//export ERROR_INVALID_HANDLE
const ERROR_INVALID_HANDLE C.uint32_t = 4

//export ERROR_ALREADY_EXISTS
const ERROR_ALREADY_EXISTS C.uint32_t = 5

//export ERROR_NOT_FOUND
const ERROR_NOT_FOUND C.uint32_t = 6

var ErrInvalidHandle = errs.Class("invalid handle")
var ErrNull = errs.Class("NULL")
var ErrInvalidArg = errs.Class("invalid argument")

func mallocError(err error) *C.Error {
	if err == nil {
		return nil
	}

	cerror := (*C.Error)(C.malloc(C.sizeof_Error))

	switch {
	case errors.Is(err, context.Canceled):
		cerror.code = ERROR_CANCELED
	case errors.Is(err, io.EOF):
		cerror.code = ERROR_EOF
		return cerror

	case ErrInvalidHandle.Has(err):
		cerror.code = ERROR_INVALID_HANDLE
	case uplink.ErrBucketAlreadyExists.Has(err):
		cerror.code = ERROR_ALREADY_EXISTS
	case uplink.ErrBucketNotFound.Has(err):
		cerror.code = ERROR_NOT_FOUND
	case uplink.ErrObjectNotFound.Has(err):
		cerror.code = ERROR_NOT_FOUND
	default:
		cerror.code = ERROR_INTERNAL
	}

	cerror.message = C.CString(fmt.Sprintf("%+v", err))
	return cerror
}

//export free_error
// free_error frees error data.
func free_error(err *C.Error) {
	if err == nil {
		return
	}

	C.free(unsafe.Pointer(err.message))
	C.free(unsafe.Pointer(err))
}
