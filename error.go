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

var (
	// ErrInvalidHandle is used when the handle passed as an argument is invalid.
	ErrInvalidHandle = errs.Class("invalid handle")
	// ErrNull is returned when an argument is NULL, however it should not be.
	ErrNull = errs.Class("NULL")
	// ErrInvalidArg is returned when the argument is not valid.
	ErrInvalidArg = errs.Class("invalid argument")
)

func mallocError(err error) *C.Error {
	if err == nil {
		return nil
	}

	cerror := (*C.Error)(C.calloc(C.sizeof_Error, 1))

	switch {
	case errors.Is(err, context.Canceled):
		cerror.code = C.ERROR_CANCELED
	case errors.Is(err, io.EOF):
		cerror.code = C.ERROR_EOF
		return cerror

	case ErrInvalidHandle.Has(err):
		cerror.code = C.ERROR_INVALID_HANDLE
	case uplink.ErrBucketAlreadyExists.Has(err):
		cerror.code = C.ERROR_ALREADY_EXISTS
	case uplink.ErrBucketNotFound.Has(err):
		cerror.code = C.ERROR_NOT_FOUND
	case uplink.ErrObjectNotFound.Has(err):
		cerror.code = C.ERROR_NOT_FOUND
	default:
		cerror.code = C.ERROR_INTERNAL
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
	defer C.free(unsafe.Pointer(err))

	if err.message != nil {
		C.free(unsafe.Pointer(err.message))
	}
}
