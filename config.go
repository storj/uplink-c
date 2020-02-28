// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	"context"
	"time"

	"storj.io/uplink"
)

//export config_request_access_with_passphrase
// config_request_access_with_passphrase requests satellite for a new access using a passhprase.
func config_request_access_with_passphrase(config C.Config, satellite_address, api_key, passphrase *C.char) C.AccessResult {
	if satellite_address != nil {
		return C.AccessResult{
			error: mallocError(ErrNull.New("satellite_address")),
		}
	}
	if api_key != nil {
		return C.AccessResult{
			error: mallocError(ErrNull.New("api_key")),
		}
	}
	if passphrase != nil {
		return C.AccessResult{
			error: mallocError(ErrNull.New("passphrase")),
		}
	}

	ctx := context.Background()

	cfg := uplinkConfig(config)

	access, err := cfg.RequestAccessWithPassphrase(ctx, C.GoString(satellite_address), C.GoString(api_key), C.GoString(passphrase))
	if err != nil {
		return C.AccessResult{
			error: mallocError(err),
		}
	}

	return C.AccessResult{
		access: (*C.Access)(mallocHandle(universe.Add(&Access{access}))),
	}
}

//export config_open_project
// config_open_project opens project using access.
func config_open_project(config C.Config, access *C.Access) C.ProjectResult {
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

	cfg := uplinkConfig(config)
	proj, err := cfg.OpenProject(scope.ctx, acc.Access)
	if err != nil {
		return C.ProjectResult{
			error: mallocError(err),
		}
	}

	return C.ProjectResult{
		project: (*C.Project)(mallocHandle(universe.Add(&Project{scope, proj}))),
	}
}

func uplinkConfig(config C.Config) uplink.Config {
	return uplink.Config{
		UserAgent:   C.GoString(config.user_agent),
		DialTimeout: time.Duration(config.dial_timeout_milliseconds) * time.Millisecond,
	}
}
