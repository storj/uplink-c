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

	case uplink.ErrRequestsLimitExceeded.Has(err):
		cerror.code = C.ERROR_TOO_MANY_REQUESTS
	case uplink.ErrBandwidthLimitExceeded.Has(err):
		cerror.code = C.ERROR_BANDWIDTH_LIMIT_EXCEEDED

	case uplink.ErrBucketNameInvalid.Has(err):
		cerror.code = C.ERROR_BUCKET_NAME_INVALID
	case uplink.ErrBucketAlreadyExists.Has(err):
		cerror.code = C.ERROR_BUCKET_ALREADY_EXISTS
	case uplink.ErrBucketNotEmpty.Has(err):
		cerror.code = C.ERROR_BUCKET_NOT_EMPTY
	case uplink.ErrBucketNotFound.Has(err):
		cerror.code = C.ERROR_BUCKET_NOT_FOUND

	case uplink.ErrObjectKeyInvalid.Has(err):
		cerror.code = C.ERROR_OBJECT_KEY_INVALID
	case uplink.ErrObjectNotFound.Has(err):
		cerror.code = C.ERROR_OBJECT_NOT_FOUND
	case uplink.ErrUploadDone.Has(err):
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
