// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"unsafe"

	"storj.io/uplink"
)

// Project provides access to managing buckets.
type Project struct {
	scope
	*uplink.Project
}

//export open_project
// open_project opens project using access.
func open_project(access *C.Access) C.ProjectResult {
	if access == nil {
		return C.ProjectResult{
			error: mallocError(ErrNull.New("access")),
		}
	}

	acc, ok := universe.Get(access._handle).(*Access)
	if !ok {
		return C.ProjectResult{
			error: mallocError(ErrInvalidHandle.New("Access")),
		}
	}

	scope := rootScope("") // TODO: should we provide this as an argument here as well?

	config := uplink.Config{}

	proj, err := config.OpenProject(scope.ctx, acc.Access)
	if err != nil {
		return C.ProjectResult{
			error: mallocError(err),
		}
	}

	return C.ProjectResult{
		project: (*C.Project)(mallocHandle(universe.Add(&Project{scope, proj}))),
	}
}

//export close_project
// close_project closes the project.
func close_project(project *C.Project) *C.Error {
	if project == nil {
		return nil
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return mallocError(ErrInvalidHandle.New("project"))
	}

	proj.cancel()
	return mallocError(proj.Close())
}

//export free_project_result
// free_project_result frees any associated resources.
func free_project_result(result C.ProjectResult) {
	free_error(result.error)
	freeProject(result.project)
}

// freeProject closes the project and frees any associated resources.
func freeProject(project *C.Project) {
	if project == nil {
		return
	}
	defer C.free(unsafe.Pointer(project))
	defer universe.Del(project._handle)

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return
	}

	proj.cancel()
	// in case we haven't already closed the project
	_ = proj.Close()
	// TODO: log error when we didn't close manually and the close returns an error
}
