// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"unsafe"

	"storj.io/uplink"
)

//export uplink_stat_bucket
// uplink_stat_bucket returns information about a bucket.
func uplink_stat_bucket(project *C.UplinkProject, bucket_name *C.uplink_const_char) C.UplinkBucketResult { //nolint:golint
	if project == nil {
		return C.UplinkBucketResult{
			error: mallocError(ErrNull.New("project")),
		}
	}
	if bucket_name == nil {
		return C.UplinkBucketResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.UplinkBucketResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}

	bucket, err := proj.StatBucket(proj.scope.ctx, C.GoString(bucket_name))

	return C.UplinkBucketResult{
		error:  mallocError(err),
		bucket: mallocBucket(bucket),
	}
}

//export uplink_create_bucket
// uplink_create_bucket creates a new bucket.
//
// When bucket already exists it returns a valid Bucket and ErrBucketExists.
func uplink_create_bucket(project *C.UplinkProject, bucket_name *C.uplink_const_char) C.UplinkBucketResult { //nolint:golint
	if project == nil {
		return C.UplinkBucketResult{
			error: mallocError(ErrNull.New("project")),
		}
	}
	if bucket_name == nil {
		return C.UplinkBucketResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.UplinkBucketResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}

	bucket, err := proj.CreateBucket(proj.scope.ctx, C.GoString(bucket_name))

	return C.UplinkBucketResult{
		error:  mallocError(err),
		bucket: mallocBucket(bucket),
	}
}

//export uplink_ensure_bucket
// uplink_ensure_bucket creates a new bucket and ignores the error when it already exists.
//
// When bucket already exists it returns a valid Bucket and ErrBucketExists.
func uplink_ensure_bucket(project *C.UplinkProject, bucket_name *C.uplink_const_char) C.UplinkBucketResult { //nolint:golint
	if project == nil {
		return C.UplinkBucketResult{
			error: mallocError(ErrNull.New("project")),
		}
	}
	if bucket_name == nil {
		return C.UplinkBucketResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.UplinkBucketResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}

	bucket, err := proj.EnsureBucket(proj.scope.ctx, C.GoString(bucket_name))

	return C.UplinkBucketResult{
		error:  mallocError(err),
		bucket: mallocBucket(bucket),
	}
}

//export uplink_delete_bucket
// uplink_delete_bucket deletes a bucket.
//
// When bucket is not empty it returns ErrBucketNotEmpty.
func uplink_delete_bucket(project *C.UplinkProject, bucket_name *C.uplink_const_char) C.UplinkBucketResult { //nolint:golint
	if project == nil {
		return C.UplinkBucketResult{
			error: mallocError(ErrNull.New("project")),
		}
	}
	if bucket_name == nil {
		return C.UplinkBucketResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.UplinkBucketResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}

	deleted, err := proj.DeleteBucket(proj.scope.ctx, C.GoString(bucket_name))
	return C.UplinkBucketResult{
		error:  mallocError(err),
		bucket: mallocBucket(deleted),
	}
}

//export uplink_delete_bucket_with_objects
// uplink_delete_bucket_with_objects deletes a bucket and all objects within that bucket.
//
// When there are concurrent writes to the bucket it returns ErrBucketNotEmpty.
func uplink_delete_bucket_with_objects(project *C.UplinkProject, bucket_name *C.uplink_const_char) C.UplinkBucketResult { //nolint:golint
	if project == nil {
		return C.UplinkBucketResult{
			error: mallocError(ErrNull.New("project")),
		}
	}
	if bucket_name == nil {
		return C.UplinkBucketResult{
			error: mallocError(ErrNull.New("bucket_name")),
		}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return C.UplinkBucketResult{
			error: mallocError(ErrInvalidHandle.New("project")),
		}
	}

	deleted, err := proj.DeleteBucketWithObjects(proj.scope.ctx, C.GoString(bucket_name))
	return C.UplinkBucketResult{
		error:  mallocError(err),
		bucket: mallocBucket(deleted),
	}
}

func mallocBucket(bucket *uplink.Bucket) *C.UplinkBucket {
	if bucket == nil {
		return nil
	}

	cbucket := (*C.UplinkBucket)(calloc(1, C.sizeof_UplinkBucket))
	cbucket.name = C.CString(bucket.Name)
	cbucket.created = timeToUnix(bucket.Created)

	return cbucket
}

//export uplink_free_bucket_result
// uplink_free_bucket_result frees memory associated with the BucketResult.
func uplink_free_bucket_result(result C.UplinkBucketResult) {
	uplink_free_error(result.error)
	uplink_free_bucket(result.bucket)
}

//export uplink_free_bucket
// uplink_free_bucket frees memory associated with the bucket.
func uplink_free_bucket(bucket *C.UplinkBucket) {
	if bucket == nil {
		return
	}
	defer C.free(unsafe.Pointer(bucket))

	if bucket.name != nil {
		C.free(unsafe.Pointer(bucket.name))
	}
}
