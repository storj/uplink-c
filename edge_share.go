// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import "storj.io/uplink/edge"

// edge_join_share_url concats a linkshare URL
// Example result: https://link.us1.storjshare.io/s/l5pucy3dmvzxgs3fpfewix27l5pq/mybucket/myprefix/myobject
// The existence or accessibility of the object is not checked,
// the object might not exist or be inaccessible.
//
// baseURL: linkshare service, e.g. https://link.us1.storjshare.io
// accessKeyId: as returned from RegisterAccess. Must be associated with public visibility.
// bucket: optional bucket, if empty shares the entire project.
// key: optional object key or prefix, if empty shares the entire bucket. A prefix must end with "/".
//
//export edge_join_share_url
func edge_join_share_url(
	baseURL *C.uplink_const_char,
	accessKeyID *C.uplink_const_char,
	bucket *C.uplink_const_char,
	key *C.uplink_const_char,
	options *C.EdgeShareURLOptions,
) C.UplinkStringResult {
	var goOptions *edge.ShareURLOptions

	if options != nil {
		goOptions = &edge.ShareURLOptions{
			Raw: bool(options.raw),
		}
	}

	url, err := edge.JoinShareURL(
		C.GoString(baseURL),
		C.GoString(accessKeyID),
		C.GoString(bucket),
		C.GoString(key),
		goOptions,
	)

	return C.UplinkStringResult{
		error:  mallocError(err),
		string: C.CString(url),
	}
}
