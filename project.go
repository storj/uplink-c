// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
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

	proj, err := config.Open(scope.ctx, acc.Access)
	if err != nil {
		return C.ProjectResult{
			error: mallocError(err),
		}
	}

	return C.ProjectResult{
		project: (*C.Project)(mallocHandle(universe.Add(&Project{scope, proj}))),
	}
}

//export free_project_result
// free_project_result closes the ProjectResult and frees any associated resources.
func free_project_result(result C.ProjectResult) *C.Error {
	free_error(result.error)
	return free_project(result.project)
}

//export free_project
// free_project closes the project and frees any associated resources.
func free_project(project *C.Project) *C.Error {
	if project == nil {
		return nil
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return mallocError(ErrInvalidHandle.New("project"))
	}

	universe.Del(project._handle)
	defer proj.cancel()

	return mallocError(proj.Close())
}
