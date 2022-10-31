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

// uplink_config_request_access_with_passphrase requests satellite for a new access grant using a passhprase.
//
//export uplink_config_request_access_with_passphrase
func uplink_config_request_access_with_passphrase(config C.UplinkConfig, satellite_address, api_key, passphrase *C.uplink_const_char) C.UplinkAccessResult { //nolint:golint
	if satellite_address == nil {
		return C.UplinkAccessResult{
			error: mallocError(ErrNull.New("satellite_address")),
		}
	}
	if api_key == nil {
		return C.UplinkAccessResult{
			error: mallocError(ErrNull.New("api_key")),
		}
	}
	if passphrase == nil {
		return C.UplinkAccessResult{
			error: mallocError(ErrNull.New("passphrase")),
		}
	}

	ctx := context.Background()

	cfg := uplinkConfig(config)

	access, err := cfg.RequestAccessWithPassphrase(ctx, C.GoString(satellite_address), C.GoString(api_key), C.GoString(passphrase))
	if err != nil {
		return C.UplinkAccessResult{
			error: mallocError(err),
		}
	}

	return C.UplinkAccessResult{
		access: (*C.UplinkAccess)(mallocHandle(universe.Add(&Access{access}))),
	}
}

// uplink_config_open_project opens project using access grant.
//
//export uplink_config_open_project
func uplink_config_open_project(config C.UplinkConfig, access *C.UplinkAccess) C.UplinkProjectResult {
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

	scope := rootScope(C.GoString(config.temp_directory))

	cfg := uplinkConfig(config)
	proj, err := cfg.OpenProject(scope.ctx, acc.Access)
	if err != nil {
		return C.UplinkProjectResult{
			error: mallocError(err),
		}
	}

	return C.UplinkProjectResult{
		project: (*C.UplinkProject)(mallocHandle(universe.Add(&Project{scope, proj}))),
	}
}

func uplinkConfig(config C.UplinkConfig) uplink.Config {
	return uplink.Config{
		UserAgent:   C.GoString(config.user_agent),
		DialTimeout: time.Duration(config.dial_timeout_milliseconds) * time.Millisecond,
	}
}
