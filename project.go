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

//export uplink_open_project
// uplink_open_project opens project using access grant.
func uplink_open_project(access *C.UplinkAccess) C.UplinkProjectResult {
	if access == nil {
		return C.UplinkProjectResult{
			error: mallocError(ErrNull.New("access")),
		}
	}

	acc, ok := universe.Get(access._handle).(*Access)
	if !ok {
		return C.UplinkProjectResult{
			error: mallocError(ErrInvalidHandle.New("Access")),
		}
	}

	scope := rootScope("")
	config := uplink.Config{}

	proj, err := config.OpenProject(scope.ctx, acc.Access)
	if err != nil {
		return C.UplinkProjectResult{
			error: mallocError(err),
		}
	}

	return C.UplinkProjectResult{
		project: (*C.UplinkProject)(mallocHandle(universe.Add(&Project{scope, proj}))),
	}
}

//export uplink_close_project
// uplink_close_project closes the project.
func uplink_close_project(project *C.UplinkProject) *C.UplinkError {
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

//export uplink_revoke_access
// uplink_revoke_access revokes the API key embedded in the provided access grant.
func uplink_revoke_access(project *C.UplinkProject, access *C.UplinkAccess) *C.UplinkError {
	if project == nil {
		return mallocError(ErrNull.New("project"))
	}

	if access == nil {
		return mallocError(ErrNull.New("access"))
	}

	proj, ok := universe.Get(project._handle).(*Project)
	if !ok {
		return mallocError(ErrInvalidHandle.New("project"))
	}

	acc, ok := universe.Get(access._handle).(*Access)
	if !ok {
		return mallocError(ErrInvalidHandle.New("access"))
	}

	scope := rootScope("")

	return mallocError(proj.RevokeAccess(scope.ctx, acc.Access))
}

//export uplink_free_project_result
// uplink_free_project_result frees any associated resources.
func uplink_free_project_result(result C.UplinkProjectResult) {
	uplink_free_error(result.error)
	freeProject(result.project)
}

// freeProject closes the project and frees any associated resources.
func freeProject(project *C.UplinkProject) {
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
