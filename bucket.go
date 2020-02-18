// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"fmt"

	"unsafe"

	"storj.io/uplink"
)

//export stat_bucket
// stat_bucket returns information about a bucket.
func stat_bucket(project C.Project, bucketName *C.char, cerr **C.char) C.Bucket {
	if bucketName == nil {
		*cerr = C.CString("bucketName == nil")
		return C.Bucket{}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		*cerr = C.CString("invalid project")
		return C.Bucket{}
	}
	child := proj.scope.child()

	bucket, err := proj.StatBucket(child.ctx, C.GoString(bucketName))
	if err != nil {
		*cerr = C.CString(fmt.Sprintf("%+v", err))
	}

	return bucketToC(bucket)
}

//export create_bucket
// create_bucket creates a new bucket.
//
// When bucket already exists it returns a valid Bucket and ErrBucketExists.
func create_bucket(project C.Project, bucketName *C.char, cerr **C.char) C.Bucket {
	if bucketName == nil {
		*cerr = C.CString("bucketName == nil")
		return C.Bucket{}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		*cerr = C.CString("invalid project")
		return C.Bucket{}
	}
	child := proj.scope.child()

	bucket, err := proj.CreateBucket(child.ctx, C.GoString(bucketName))
	if err != nil {
		*cerr = C.CString(fmt.Sprintf("%+v", err))
	}

	return bucketToC(bucket)
}

//export ensure_bucket
// ensure_bucket creates a new bucket and ignores the error when it already exists.
//
// When bucket already exists it returns a valid Bucket and ErrBucketExists.
func ensure_bucket(project C.Project, bucketName *C.char, cerr **C.char) C.Bucket {
	if bucketName == nil {
		*cerr = C.CString("bucketName == nil")
		return C.Bucket{}
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		*cerr = C.CString("invalid project")
		return C.Bucket{}
	}
	child := proj.scope.child()

	bucket, err := proj.EnsureBucket(child.ctx, C.GoString(bucketName))
	if err != nil {
		*cerr = C.CString(fmt.Sprintf("%+v", err))
	}

	return bucketToC(bucket)
}

//export delete_bucket
// delete_bucket deletes a bucket.
//
// When bucket is not empty it returns ErrBucketNotEmpty.
func delete_bucket(project C.Project, bucketName *C.char, cerr **C.char) {
	if bucketName == nil {
		*cerr = C.CString("bucketName == nil")
		return
	}
	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		*cerr = C.CString("invalid project")
		return
	}
	child := proj.scope.child()

	err := proj.DeleteBucket(child.ctx, C.GoString(bucketName))
	if err != nil {
		*cerr = C.CString(fmt.Sprintf("%+v", err))
	}
}

func bucketToC(bucket *uplink.Bucket) C.Bucket {
	if bucket == nil {
		return C.Bucket{}
	}
	return C.Bucket{
		name:    C.CString(bucket.Name),
		created: C.int64_t(bucket.Created.Unix()),
	}
}

//export free_bucket
// free_bucket frees memory associated with the bucket.
func free_bucket(bucket C.Bucket) {
	if bucket.name != nil {
		C.free(unsafe.Pointer(bucket.name))
		bucket.name = nil
	}
}
