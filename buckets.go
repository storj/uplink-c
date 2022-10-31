// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"unsafe"

	"storj.io/uplink"
)

// BucketIterator is an iterator over buckets.
type BucketIterator struct {
	scope
	iterator *uplink.BucketIterator

	initialError error
}

// uplink_list_buckets lists buckets.
//
//export uplink_list_buckets
func uplink_list_buckets(project *C.UplinkProject, options *C.UplinkListBucketsOptions) *C.UplinkBucketIterator {
	if project == nil {
		return (*C.UplinkBucketIterator)(mallocHandle(universe.Add(&BucketIterator{
			initialError: ErrNull.New("project"),
		})))
	}
	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return (*C.UplinkBucketIterator)(mallocHandle(universe.Add(&BucketIterator{
			initialError: ErrInvalidHandle.New("project"),
		})))
	}

	opts := &uplink.ListBucketsOptions{}
	if options != nil {
		opts.Cursor = C.GoString(options.cursor)
	}

	scope := proj.scope.child()
	iterator := proj.ListBuckets(scope.ctx, opts)
	return (*C.UplinkBucketIterator)(mallocHandle(universe.Add(&BucketIterator{
		scope:    scope,
		iterator: iterator,
	})))
}

// uplink_bucket_iterator_next prepares next Bucket for reading.
//
// It returns false if the end of the iteration is reached and there are no more buckets, or if there is an error.
//
//export uplink_bucket_iterator_next
func uplink_bucket_iterator_next(iterator *C.UplinkBucketIterator) C.bool {
	if iterator == nil {
		return C.bool(false)
	}

	iter, ok := universe.Get(iterator._handle).(*BucketIterator)
	if !ok {
		return C.bool(false)
	}
	if iter.initialError != nil {
		return C.bool(false)
	}

	return C.bool(iter.iterator.Next())
}

// uplink_bucket_iterator_err returns error, if one happened during iteration.
//
//export uplink_bucket_iterator_err
func uplink_bucket_iterator_err(iterator *C.UplinkBucketIterator) *C.UplinkError {
	if iterator == nil {
		return mallocError(ErrNull.New("iterator"))
	}

	iter, ok := universe.Get(iterator._handle).(*BucketIterator)
	if !ok {
		return mallocError(ErrInvalidHandle.New("iterator"))
	}
	if iter.initialError != nil {
		return mallocError(iter.initialError)
	}

	return mallocError(iter.iterator.Err())
}

// uplink_bucket_iterator_item returns the current bucket in the iterator.
//
//export uplink_bucket_iterator_item
func uplink_bucket_iterator_item(iterator *C.UplinkBucketIterator) *C.UplinkBucket {
	if iterator == nil {
		return nil
	}

	iter, ok := universe.Get(iterator._handle).(*BucketIterator)
	if !ok {
		return nil
	}

	return mallocBucket(iter.iterator.Item())
}

// uplink_free_bucket_iterator frees memory associated with the BucketIterator.
//
//export uplink_free_bucket_iterator
func uplink_free_bucket_iterator(iterator *C.UplinkBucketIterator) {
	if iterator == nil {
		return
	}
	defer C.free(unsafe.Pointer(iterator))
	defer universe.Del(iterator._handle)

	iter, ok := universe.Get(iterator._handle).(*BucketIterator)
	if ok {
		if iter.scope.cancel != nil {
			iter.scope.cancel()
		}
	}
}
