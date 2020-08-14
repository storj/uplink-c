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

func mallocError(err error) *C.UplinkError {
	if err == nil {
		return nil
	}

	cerror := (*C.UplinkError)(C.calloc(C.sizeof_UplinkError, 1))

	switch {
	case errors.Is(err, io.EOF):
		cerror.code = C.EOF
		return cerror
	case errors.Is(err, context.Canceled):
		cerror.code = C.UPLINK_ERROR_CANCELED
	case ErrInvalidHandle.Has(err):
		cerror.code = C.UPLINK_ERROR_INVALID_HANDLE

	case errors.Is(err, uplink.ErrTooManyRequests):
		cerror.code = C.UPLINK_ERROR_TOO_MANY_REQUESTS
	case errors.Is(err, uplink.ErrBandwidthLimitExceeded):
		cerror.code = C.UPLINK_ERROR_BANDWIDTH_LIMIT_EXCEEDED

	case errors.Is(err, uplink.ErrBucketNameInvalid):
		cerror.code = C.UPLINK_ERROR_BUCKET_NAME_INVALID
	case errors.Is(err, uplink.ErrBucketAlreadyExists):
		cerror.code = C.UPLINK_ERROR_BUCKET_ALREADY_EXISTS
	case errors.Is(err, uplink.ErrBucketNotEmpty):
		cerror.code = C.UPLINK_ERROR_BUCKET_NOT_EMPTY
	case errors.Is(err, uplink.ErrBucketNotFound):
		cerror.code = C.UPLINK_ERROR_BUCKET_NOT_FOUND

	case errors.Is(err, uplink.ErrObjectKeyInvalid):
		cerror.code = C.UPLINK_ERROR_OBJECT_KEY_INVALID
	case errors.Is(err, uplink.ErrObjectNotFound):
		cerror.code = C.UPLINK_ERROR_OBJECT_NOT_FOUND
	case errors.Is(err, uplink.ErrUploadDone):
		cerror.code = C.UPLINK_ERROR_UPLOAD_DONE

	default:
		cerror.code = C.UPLINK_ERROR_INTERNAL
	}

	cerror.message = C.CString(fmt.Sprintf("%+v", err))
	return cerror
}

//export uplink_free_error
// uplink_free_error frees error data.
func uplink_free_error(err *C.UplinkError) {
	if err == nil {
		return
	}
	defer C.free(unsafe.Pointer(err))

	if err.message != nil {
		C.free(unsafe.Pointer(err.message))
	}
}
