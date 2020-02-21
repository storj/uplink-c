// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"storj.io/uplink"
)

//export list_buckets
// list_buckets lists buckets
func list_buckets(project *C.Project, options *C.ListBucketsOptions) *C.BucketIterator {
	if project == nil {
		// TODO: should we return an error here?
		return nil
	}
	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		// TODO: should we return an error here?
		return nil
	}

	opts := &uplink.ListBucketsOptions{}
	if options != nil {
		opts.Cursor = C.GoString(options.cursor)
	}

	// TODO: should we pass in a separate ctx?
	iterator := proj.ListBuckets(proj.scope.ctx, opts)

	return (*C.BucketIterator)(mallocHandle(universe.Add(iterator)))
}

//export bucket_iterator_next
// bucket_iterator_next prepares next Bucket for reading.
//
// It returns false if the end of the iteration is reached and there are no more buckets, or if there is an error.
func bucket_iterator_next(iterator *C.BucketIterator) C.bool {
	if iterator == nil {
		return false
	}

	iter, ok := universe.Get(iterator._handle).(*uplink.BucketIterator)
	if !ok {
		return false
	}

	return C.bool(iter.Next())
}

//export bucket_iterator_err
// bucket_iterator_err returns error, if one happened during iteration.
func bucket_iterator_err(iterator *C.BucketIterator) *C.Error {
	if iterator == nil {
		return mallocError(ErrNull.New("iterator"))
	}

	iter, ok := universe.Get(iterator._handle).(*uplink.BucketIterator)
	if !ok {
		return mallocError(ErrInvalidHandle.New("iterator"))
	}

	return mallocError(iter.Err())
}

//export bucket_iterator_item
// bucket_iterator_item returns the current bucket in the iterator.
func bucket_iterator_item(iterator *C.BucketIterator) *C.Bucket {
	if iterator == nil {
		return nil
	}

	iter, ok := universe.Get(iterator._handle).(*uplink.BucketIterator)
	if !ok {
		return nil
	}

	return mallocBucket(iter.Item())
}

//export free_bucket_iterator
// free_bucket_iterator frees memory associated with the BucketIterator.
func free_bucket_iterator(iterator *C.BucketIterator) {
	if iterator == nil {
		return
	}

	universe.Del(iterator._handle)
}
