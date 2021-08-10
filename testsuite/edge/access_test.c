// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

#include <stdio.h>
#include <string.h>

#include "../require.h"
#define UPLINK_DISABLE_NAMESPACE_COMPAT
#include "uplink.h"

const char *minimal_access =
    "13J4Upun87ATb3T5T5sDXVeQaCzWFZeF9Ly4ELfxS5hUwTL8APEkwahTEJ1wxZjyErimiDs3kgid33kDLuYPYtwaY7Toy32mCTapfrUB814X13RiA8"
    "44HPWK3QLKZb9cAoVceTowmNZXWbcUMKNbkMHCURE4hn8ZrdHPE3S86yngjvDxwKmarfGx";

int main()
{
    const char *auth_service_addr = getenv("AUTH_SERVICE_ADDR");
    const char *auth_service_cert = getenv("AUTH_SERVICE_CERT");

    fprintf(stdout, "Auth service address is: %s\n", auth_service_addr);

    UplinkAccessResult access_result = uplink_parse_access(minimal_access);
    require_noerror(access_result.error);
    UplinkAccess *access = access_result.access;

    {
        // Happy flow
        EdgeConfig config = {
            .auth_service_address = auth_service_addr,
            .certificate_pem = auth_service_cert,
        };

        EdgeCredentialsResult credentials_result = edge_register_access(config, access, NULL);
        require_noerror(credentials_result.error);

        EdgeCredentials credentials = *credentials_result.credentials;
        require(strcmp("l5pucy3dmvzxgs3fpfewix27l5pq", credentials.access_key_id) == 0);
        require(strcmp("l5pvgzldojsxis3fpfpv6x27l5pv6x27l5pv6x27l5pv6", credentials.secret_key) == 0);
        require(strcmp("https://gateway.example", credentials.endpoint) == 0);

        edge_free_credentials_result(credentials_result);
    }

    {
        // TLS certificate error
        EdgeConfig config = {
            .auth_service_address = auth_service_addr,
        };

        EdgeCredentialsResult credentials_result = edge_register_access(config, access, NULL);
        require_error(credentials_result.error, EDGE_ERROR_AUTH_DIAL_FAILED);
        fprintf(stdout, "TLS error is: %s\n", credentials_result.error->message);

        edge_free_credentials_result(credentials_result);
    }

    {
        // DNS error host does not exist
        EdgeConfig config = {
            .auth_service_address = "doesnotexist.example:1234",
            .certificate_pem = auth_service_cert,
        };

        EdgeCredentialsResult credentials_result = edge_register_access(config, access, NULL);
        require_error(credentials_result.error, EDGE_ERROR_AUTH_DIAL_FAILED);
        fprintf(stdout, "DNS error is: %s\n", credentials_result.error->message);

        edge_free_credentials_result(credentials_result);
    }

    {
        // IP error no server running at this address
        EdgeConfig config = {
            .auth_service_address = "127.0.0.2:864",
            .certificate_pem = auth_service_cert,
        };

        EdgeCredentialsResult credentials_result = edge_register_access(config, access, NULL);
        require_error(credentials_result.error, EDGE_ERROR_AUTH_DIAL_FAILED);
        fprintf(stdout, "IP error is: %s\n", credentials_result.error->message);

        edge_free_credentials_result(credentials_result);
    }

    uplink_free_access_result(access_result);

    return 0;
}
