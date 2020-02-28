// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"reflect"
	"sort"
	"unsafe"

	"storj.io/uplink"
)

// note: while there are restrictions on what can be stored in custom metadata.
// the following functions should work with arbitrary byte strings as keys and values.

func customMetadataToC(customMetadata uplink.CustomMetadata) C.CustomMetadata {
	if customMetadata == nil {
		return C.CustomMetadata{}
	}

	type entry struct {
		key   string
		value string
	}

	var sorted []entry
	for k, v := range customMetadata {
		sorted = append(sorted, entry{key: k, value: v})
	}
	sort.Slice(sorted, func(i, k int) bool { return sorted[i].key < sorted[k].key })

	entries := (*C.CustomMetadataEntry)(C.calloc(C.sizeof_CustomMetadataEntry, C.size_t(len(sorted))))
	custom := C.CustomMetadata{
		entries: entries,
		count:   C.uint64_t(len(sorted)),
	}

	var array []C.CustomMetadataEntry
	*(*reflect.SliceHeader)(unsafe.Pointer(&array)) = reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(entries)),
		Len:  len(sorted),
		Cap:  len(sorted),
	}

	for i, kv := range sorted {
		ckey := C.CString(kv.key)

		array[i] = C.CustomMetadataEntry{
			key:        ckey,
			key_length: C.uint64_t(len(kv.key)),

			value:        C.CString(kv.value),
			value_length: C.uint64_t(len(kv.value)),
		}
	}

	return custom
}

func customMetadataFromC(custom C.CustomMetadata) uplink.CustomMetadata {
	if custom.count == 0 {
		return uplink.CustomMetadata{}
	}

	customMetadata := uplink.CustomMetadata{}

	var array []C.CustomMetadataEntry
	*(*reflect.SliceHeader)(unsafe.Pointer(&array)) = reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(custom.entries)),
		Len:  int(custom.count),
		Cap:  int(custom.count),
	}

	for _, e := range array {
		key := C.GoStringN(e.key, C.int(e.key_length))
		value := C.GoStringN(e.value, C.int(e.value_length))
		customMetadata[key] = value
	}

	return customMetadata
}

func free_custom_metadata_data(custom *C.CustomMetadata) {
	if custom.entries == nil {
		return
	}
	defer func() {
		C.free(unsafe.Pointer(custom.entries))
	}()


	var array []C.CustomMetadataEntry
	*(*reflect.SliceHeader)(unsafe.Pointer(&array)) = reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(custom.entries)),
		Len:  int(custom.count),
		Cap:  int(custom.count),
	}

	for i := range array {
		e := &array[i]
		C.free(unsafe.Pointer(e.key))
		e.key = nil
		C.free(unsafe.Pointer(e.value))
		e.value = nil
	}
}
