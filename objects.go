// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"unsafe"

	"storj.io/uplink"
)

//export list_objects
// list_objects lists objects
func list_objects(project *C.Project, bucket_name *C.char, options *C.ListObjectsOptions) *C.ObjectIterator {
	if project == nil {
		// TODO: should we return an error here?
		return nil
	}
	if bucket_name == nil {
		// TODO: should we return an error here?
		return nil
	}
	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		// TODO: should we return an error here?
		return nil
	}

	opts := &uplink.ListObjectsOptions{}
	if options != nil {
		opts.Prefix = C.GoString(options.prefix)
		opts.Cursor = C.GoString(options.cursor)
		opts.Recursive = bool(options.recursive)

		opts.Info = bool(options.info)
		opts.Standard = bool(options.standard)
		opts.Custom = bool(options.custom)
	}

	// TODO: should we pass in a separate ctx?
	iterator := proj.ListObjects(proj.scope.ctx, C.GoString(bucket_name), opts)

	return (*C.ObjectIterator)(mallocHandle(universe.Add(iterator)))
}

//export object_iterator_next
// object_iterator_next prepares next Object for reading.
//
// It returns false if the end of the iteration is reached and there are no more objects, or if there is an error.
func object_iterator_next(iterator *C.ObjectIterator) C.bool {
	if iterator == nil {
		return false
	}

	iter, ok := universe.Get(iterator._handle).(*uplink.ObjectIterator)
	if !ok {
		return false
	}

	return C.bool(iter.Next())
}

//export object_iterator_err
// object_iterator_err returns error, if one happened during iteration.
func object_iterator_err(iterator *C.ObjectIterator) *C.Error {
	if iterator == nil {
		return mallocError(ErrNull.New("iterator"))
	}

	iter, ok := universe.Get(iterator._handle).(*uplink.ObjectIterator)
	if !ok {
		return mallocError(ErrInvalidHandle.New("iterator"))
	}

	return mallocError(iter.Err())
}

//export object_iterator_item
// object_iterator_item returns the current bucket in the iterator.
func object_iterator_item(iterator *C.ObjectIterator) *C.Object {
	if iterator == nil {
		return nil
	}

	iter, ok := universe.Get(iterator._handle).(*uplink.ObjectIterator)
	if !ok {
		return nil
	}

	return mallocObject(iter.Item())
}

//export free_object_iterator
// free_object_iterator frees memory associated with the ObjectIterator.
func free_object_iterator(iterator *C.ObjectIterator) {
	if iterator == nil {
		return
	}
	defer C.free(unsafe.Pointer(iterator))

	universe.Del(iterator._handle)
}
