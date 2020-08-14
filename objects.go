// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"unsafe"

	"storj.io/uplink"
)

// ObjectIterator is an iterator over objects.
type ObjectIterator struct {
	scope
	iterator *uplink.ObjectIterator

	initialError error
}

//export uplink_list_objects
// uplink_list_objects lists objects.
func uplink_list_objects(project *C.UplinkProject, bucket_name *C.uplink_const_char, options *C.UplinkListObjectsOptions) *C.UplinkObjectIterator { //nolint:golint
	if project == nil {
		return (*C.UplinkObjectIterator)(mallocHandle(universe.Add(&ObjectIterator{
			initialError: ErrNull.New("project"),
		})))
	}
	if bucket_name == nil {
		return (*C.UplinkObjectIterator)(mallocHandle(universe.Add(&ObjectIterator{
			initialError: ErrNull.New("bucket_name"),
		})))
	}
	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return (*C.UplinkObjectIterator)(mallocHandle(universe.Add(&ObjectIterator{
			initialError: ErrInvalidHandle.New("project"),
		})))
	}

	opts := &uplink.ListObjectsOptions{}
	if options != nil {
		opts.Prefix = C.GoString(options.prefix)
		opts.Cursor = C.GoString(options.cursor)
		opts.Recursive = bool(options.recursive)

		opts.System = bool(options.system)
		opts.Custom = bool(options.custom)
	}

	scope := proj.scope.child()
	iterator := proj.ListObjects(scope.ctx, C.GoString(bucket_name), opts)

	return (*C.UplinkObjectIterator)(mallocHandle(universe.Add(&ObjectIterator{
		scope:    scope,
		iterator: iterator,
	})))
}

//export uplink_object_iterator_next
// uplink_object_iterator_next prepares next Object for reading.
//
// It returns false if the end of the iteration is reached and there are no more objects, or if there is an error.
func uplink_object_iterator_next(iterator *C.UplinkObjectIterator) C.bool {
	if iterator == nil {
		return C.bool(false)
	}

	iter, ok := universe.Get(iterator._handle).(*ObjectIterator)
	if !ok {
		return C.bool(false)
	}
	if iter.initialError != nil {
		return C.bool(false)
	}

	return C.bool(iter.iterator.Next())
}

//export uplink_object_iterator_err
// uplink_object_iterator_err returns error, if one happened during iteration.
func uplink_object_iterator_err(iterator *C.UplinkObjectIterator) *C.UplinkError {
	if iterator == nil {
		return mallocError(ErrNull.New("iterator"))
	}

	iter, ok := universe.Get(iterator._handle).(*ObjectIterator)
	if !ok {
		return mallocError(ErrInvalidHandle.New("iterator"))
	}
	if iter.initialError != nil {
		return mallocError(iter.initialError)
	}

	return mallocError(iter.iterator.Err())
}

//export uplink_object_iterator_item
// uplink_object_iterator_item returns the current object in the iterator.
func uplink_object_iterator_item(iterator *C.UplinkObjectIterator) *C.UplinkObject {
	if iterator == nil {
		return nil
	}

	iter, ok := universe.Get(iterator._handle).(*ObjectIterator)
	if !ok {
		return nil
	}

	return mallocObject(iter.iterator.Item())
}

//export uplink_free_object_iterator
// uplink_free_object_iterator frees memory associated with the ObjectIterator.
func uplink_free_object_iterator(iterator *C.UplinkObjectIterator) {
	if iterator == nil {
		return
	}
	defer C.free(unsafe.Pointer(iterator))
	defer universe.Del(iterator._handle)

	iter, ok := universe.Get(iterator._handle).(*ObjectIterator)
	if ok {
		if iter.scope.cancel != nil {
			iter.scope.cancel()
		}
	}
}
