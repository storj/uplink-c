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
	case errors.Is(err, io.EOF):
		cerror.code = C.EOF
		return cerror
	case errors.Is(err, context.Canceled):
		cerror.code = C.ERROR_CANCELED
	case ErrInvalidHandle.Has(err):
		cerror.code = C.ERROR_INVALID_HANDLE

	case errors.Is(err, uplink.ErrTooManyRequests):
		cerror.code = C.ERROR_TOO_MANY_REQUESTS
	case errors.Is(err, uplink.ErrBandwidthLimitExceeded):
		cerror.code = C.ERROR_BANDWIDTH_LIMIT_EXCEEDED

	case errors.Is(err, uplink.ErrBucketNameInvalid):
		cerror.code = C.ERROR_BUCKET_NAME_INVALID
	case errors.Is(err, uplink.ErrBucketAlreadyExists):
		cerror.code = C.ERROR_BUCKET_ALREADY_EXISTS
	case errors.Is(err, uplink.ErrBucketNotEmpty):
		cerror.code = C.ERROR_BUCKET_NOT_EMPTY
	case errors.Is(err, uplink.ErrBucketNotFound):
		cerror.code = C.ERROR_BUCKET_NOT_FOUND

	case errors.Is(err, uplink.ErrObjectKeyInvalid):
		cerror.code = C.ERROR_OBJECT_KEY_INVALID
	case errors.Is(err, uplink.ErrObjectNotFound):
		cerror.code = C.ERROR_OBJECT_NOT_FOUND
	case errors.Is(err, uplink.ErrUploadDone):
		cerror.code = C.ERROR_UPLOAD_DONE

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
