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
)

//export SUCCESS
const SUCCESS C.uint32_t = 0

//export EOF
const EOF C.uint32_t = 1

//export ERROR_INTERNAL
const ERROR_INTERNAL C.uint32_t = 2

//export ERROR_CANCELED
const ERROR_CANCELED C.uint32_t = 3

//export ERROR_INVALID_HANDLE
const ERROR_INVALID_HANDLE C.uint32_t = 4

var ErrInvalidHandle = errs.Class("invalid handle")

func convertToError(err error) *C.Error {
	if err == nil {
		return nil
	}

	cerror := (*C.Error)(C.malloc(C.sizeof_Error))

	switch {
	case errors.Is(err, context.Canceled):
		cerror.code = ERROR_CANCELED
	case errors.Is(err, io.EOF):
		cerror.code = EOF
		return cerror
	case ErrInvalidHandle.Has(err):
		cerror.code = ERROR_INVALID_HANDLE
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
