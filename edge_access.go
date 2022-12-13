// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

// in uplink "edge" is a separate package but it's unclear if that is possible in uplink-c
package main

// #include "uplink_definitions.h"
import "C"
import (
	"context"
	"unsafe"

	"storj.io/uplink/edge"
)

// edge_register_access gets credentials for the Storj-hosted Gateway-mt and linkshare service.
// All files uploaded under the Access are then accessible via those services.
//
//export edge_register_access
func edge_register_access(
	config C.EdgeConfig,
	access *C.UplinkAccess,
	options *C.EdgeRegisterAccessOptions,
) C.EdgeCredentialsResult {
	goConfig := edge.Config{
		AuthServiceAddress:            C.GoString(config.auth_service_address),
		CertificatePEM:                []byte(C.GoString(config.certificate_pem)),
		InsecureUnencryptedConnection: bool(config.insecure_unencrypted_connection),
	}
	if options == nil {
		options = &C.EdgeRegisterAccessOptions{}
	}
	goOptions := edge.RegisterAccessOptions{
		Public: bool(options.is_public),
	}

	goAccess, ok := universe.Get(access._handle).(*Access)
	if !ok {
		return C.EdgeCredentialsResult{
			error: mallocError(ErrInvalidHandle.New("access")),
		}
	}

	ctx := context.Background()
	goCredentials, err := goConfig.RegisterAccess(
		ctx,
		goAccess.Access,
		&goOptions,
	)

	return C.EdgeCredentialsResult{
		error:       mallocError(err),
		credentials: mallocEdgeCredentials(goCredentials),
	}
}

//export edge_free_credentials_result
func edge_free_credentials_result(result C.EdgeCredentialsResult) {
	uplink_free_error(result.error)
	edge_free_credentials(result.credentials)
}

//export edge_free_credentials
func edge_free_credentials(credentials *C.EdgeCredentials) {
	if credentials == nil {
		return
	}

	defer C.free(unsafe.Pointer(credentials))

	if credentials.access_key_id != nil {
		C.free(unsafe.Pointer(credentials.access_key_id))
	}
	if credentials.secret_key != nil {
		C.free(unsafe.Pointer(credentials.secret_key))
	}
	if credentials.endpoint != nil {
		C.free(unsafe.Pointer(credentials.endpoint))
	}
}

func mallocEdgeCredentials(credentials *edge.Credentials) *C.EdgeCredentials {
	if credentials == nil {
		return nil
	}

	cCredentials := (*C.EdgeCredentials)(calloc(1, C.sizeof_EdgeCredentials))
	*cCredentials = C.EdgeCredentials{
		access_key_id: C.CString(credentials.AccessKeyID),
		secret_key:    C.CString(credentials.SecretKey),
		endpoint:      C.CString(credentials.Endpoint),
	}
	return cCredentials
}
