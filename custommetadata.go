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

func customMetadataToC(customMetadata uplink.CustomMetadata) C.UplinkCustomMetadata {
	if customMetadata == nil {
		return C.UplinkCustomMetadata{}
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

	entries := (*C.UplinkCustomMetadataEntry)(C.calloc(C.size_t(len(sorted)), C.sizeof_UplinkCustomMetadataEntry))
	custom := C.UplinkCustomMetadata{
		entries: entries,
		count:   C.size_t(len(sorted)),
	}

	var array []C.UplinkCustomMetadataEntry
	harray := (*reflect.SliceHeader)(unsafe.Pointer(&array))
	harray.Data = uintptr(unsafe.Pointer(entries))
	harray.Len = len(sorted)
	harray.Cap = len(sorted)

	for i, kv := range sorted {
		ckey := C.CString(kv.key)

		array[i] = C.UplinkCustomMetadataEntry{
			key:        ckey,
			key_length: C.size_t(len(kv.key)),

			value:        C.CString(kv.value),
			value_length: C.size_t(len(kv.value)),
		}
	}

	return custom
}

func customMetadataFromC(custom C.UplinkCustomMetadata) uplink.CustomMetadata {
	if custom.count == 0 {
		return uplink.CustomMetadata{}
	}

	customMetadata := uplink.CustomMetadata{}

	var array []C.UplinkCustomMetadataEntry
	harray := (*reflect.SliceHeader)(unsafe.Pointer(&array))
	harray.Data = uintptr(unsafe.Pointer(custom.entries))
	harray.Len = int(custom.count)
	harray.Cap = int(custom.count)

	for _, e := range array {
		key := C.GoStringN(e.key, C.int(e.key_length))
		value := C.GoStringN(e.value, C.int(e.value_length))
		customMetadata[key] = value
	}

	return customMetadata
}

func freeCustomMetadataData(custom *C.UplinkCustomMetadata) {
	if custom.entries == nil {
		return
	}
	defer func() {
		C.free(unsafe.Pointer(custom.entries))
	}()

	var array []C.UplinkCustomMetadataEntry
	harray := (*reflect.SliceHeader)(unsafe.Pointer(&array))
	harray.Data = uintptr(unsafe.Pointer(custom.entries))
	harray.Len = int(custom.count)
	harray.Cap = int(custom.count)

	for i := range array {
		e := &array[i]
		C.free(unsafe.Pointer(e.key))
		e.key = nil
		C.free(unsafe.Pointer(e.value))
		e.value = nil
	}
}
