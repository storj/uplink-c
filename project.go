// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"storj.io/uplink"
)

// Project is a scoped uplink.Project
type Project struct {
	scope
	*uplink.Project
}

//export open_project
// open_project opens project using uplink
func open_project(access C.Access, cerr **C.char) C.Project {
	return C.Project{}
}
