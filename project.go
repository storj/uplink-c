// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"fmt"

	"storj.io/uplink"
)

// Project provides access to managing buckets.
type Project struct {
	scope
	*uplink.Project
}

//export open_project
// open_project opens project using access.
func open_project(access C.Access, cerr **C.char) C.Project {
	acc, ok := universe.Get(access._handle).(*Access)
	if !ok {
		*cerr = C.CString("invalid access")
		return C.Project{}
	}

	scope := rootScope("") // TODO: should we provide this as an argument here as well?

	// TODO: remove testcode
	config := uplink.Config{
		Whitelist: uplink.InsecureSkipConnectionVerify(),
	}

	proj, err := config.Open(scope.ctx, acc.Access)
	if err != nil {
		*cerr = C.CString(fmt.Sprintf("%+v", err))
		return C.Project{}
	}

	return C.Project{universe.Add(&Project{scope, proj})}
}

//export free_project
// free_project closes the project and frees any associated resources.
func free_project(project C.Project, cerr **C.char) {
	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		*cerr = C.CString("invalid project")
		return
	}

	universe.Del(project._handle)
	defer proj.cancel()

	if err := proj.Close(); err != nil {
		*cerr = C.CString(fmt.Sprintf("%+v", err))
		return
	}
}
