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

//export list_objects
// list_objects lists objects
func list_objects(project *C.Project, bucket_name *C.char, options *C.ListObjectsOptions) *C.ObjectIterator { //nolint:golint
	if project == nil {
		return (*C.ObjectIterator)(mallocHandle(universe.Add(&ObjectIterator{
			initialError: ErrNull.New("project"),
		})))
	}
	if bucket_name == nil {
		return (*C.ObjectIterator)(mallocHandle(universe.Add(&ObjectIterator{
			initialError: ErrNull.New("bucket_name"),
		})))
	}
	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return (*C.ObjectIterator)(mallocHandle(universe.Add(&ObjectIterator{
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

	return (*C.ObjectIterator)(mallocHandle(universe.Add(&ObjectIterator{
		scope:    scope,
		iterator: iterator,
	})))
}

//export object_iterator_next
// object_iterator_next prepares next Object for reading.
//
// It returns false if the end of the iteration is reached and there are no more objects, or if there is an error.
func object_iterator_next(iterator *C.ObjectIterator) C.bool {
	if iterator == nil {
		return false
	}

	iter, ok := universe.Get(iterator._handle).(*ObjectIterator)
	if !ok {
		return false
	}
	if iter.initialError != nil {
		return C.bool(false)
	}

	return C.bool(iter.iterator.Next())
}

//export object_iterator_err
// object_iterator_err returns error, if one happened during iteration.
func object_iterator_err(iterator *C.ObjectIterator) *C.Error {
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

//export object_iterator_item
// object_iterator_item returns the current bucket in the iterator.
func object_iterator_item(iterator *C.ObjectIterator) *C.Object {
	if iterator == nil {
		return nil
	}

	iter, ok := universe.Get(iterator._handle).(*ObjectIterator)
	if !ok {
		return nil
	}

	return mallocObject(iter.iterator.Item())
}

//export free_object_iterator
// free_object_iterator frees memory associated with the ObjectIterator.
func free_object_iterator(iterator *C.ObjectIterator) {
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
