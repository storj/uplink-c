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

//export list_buckets
// list_buckets lists buckets
func list_buckets(project *C.Project, options *C.ListBucketsOptions) *C.BucketIterator {
	if project == nil {
		return (*C.BucketIterator)(mallocHandle(universe.Add(&BucketIterator{
			initialError: ErrNull.New("project"),
		})))
	}
	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return (*C.BucketIterator)(mallocHandle(universe.Add(&BucketIterator{
			initialError: ErrInvalidHandle.New("project"),
		})))
	}

	opts := &uplink.ListBucketsOptions{}
	if options != nil {
		opts.Cursor = C.GoString(options.cursor)
	}

	scope := proj.scope.child()
	iterator := proj.ListBuckets(scope.ctx, opts)
	return (*C.BucketIterator)(mallocHandle(universe.Add(&BucketIterator{
		scope:    scope,
		iterator: iterator,
	})))
}

//export bucket_iterator_next
// bucket_iterator_next prepares next Bucket for reading.
//
// It returns false if the end of the iteration is reached and there are no more buckets, or if there is an error.
func bucket_iterator_next(iterator *C.BucketIterator) C.bool {
	if iterator == nil {
		return false
	}

	iter, ok := universe.Get(iterator._handle).(*BucketIterator)
	if !ok {
		return false
	}
	if iter.initialError != nil {
		return false
	}

	return C.bool(iter.iterator.Next())
}

//export bucket_iterator_err
// bucket_iterator_err returns error, if one happened during iteration.
func bucket_iterator_err(iterator *C.BucketIterator) *C.Error {
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

//export bucket_iterator_item
// bucket_iterator_item returns the current bucket in the iterator.
func bucket_iterator_item(iterator *C.BucketIterator) *C.Bucket {
	if iterator == nil {
		return nil
	}

	iter, ok := universe.Get(iterator._handle).(*BucketIterator)
	if !ok {
		return nil
	}

	return mallocBucket(iter.iterator.Item())
}

//export free_bucket_iterator
// free_bucket_iterator frees memory associated with the BucketIterator.
func free_bucket_iterator(iterator *C.BucketIterator) {
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
