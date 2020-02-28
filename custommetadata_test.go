// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/testrand"
	"storj.io/uplink"
)

func TestCustomMetadata_Conversion(t *testing.T) {
	inputs := []uplink.CustomMetadata{
		{},
		{"A": "B"},
		{"A": "", "": "B"},
		{"": ""},
		{"\x00": "\x00", "\xFF": "\xFF"},
	}

	for _, meta := range inputs {
		t.Log(fmt.Sprintf("%+v", meta))
		cmeta := customMetadataToC(meta)
		gometa := customMetadataFromC(cmeta)
		require.Equal(t, meta, gometa)
		freeCustomMetadataData(&cmeta)
	}
}

func TestCustomMetadata_Random(t *testing.T) {
	for i := 0; i < 100; i++ {
		meta := uplink.CustomMetadata{}

		for k := 0; k < testrand.Intn(10); k++ {
			key := string(testrand.BytesInt(20))
			value := string(testrand.BytesInt(20))
			meta[key] = value
		}

		cmeta := customMetadataToC(meta)
		gometa := customMetadataFromC(cmeta)
		require.Equal(t, meta, gometa)
		freeCustomMetadataData(&cmeta)
	}
}
