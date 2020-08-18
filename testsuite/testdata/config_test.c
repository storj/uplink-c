// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "helpers.h"
#include "require.h"
#include "uplink.h"

int main()
{
    const char *satellite_addr = getenv("SATELLITE_0_ADDR");
    const char *api_key = getenv("UPLINK_0_APIKEY");
    const char *access_string = getenv("UPLINK_0_ACCESS");

    UplinkConfig config = {
        .user_agent = (const char *)"Test/1.0",
        .dial_timeout_milliseconds = 10000,
    };

    {
        UplinkAccessResult access_result =
            uplink_config_request_access_with_passphrase(config, NULL, api_key, "mypassphrase");
        require_error(access_result.error, UPLINK_ERROR_INTERNAL);
        require(access_result.access == NULL);
        uplink_free_access_result(access_result);

        access_result = uplink_config_request_access_with_passphrase(config, satellite_addr, NULL, "mypassphrase");
        require_error(access_result.error, UPLINK_ERROR_INTERNAL);
        require(access_result.access == NULL);
        uplink_free_access_result(access_result);

        access_result = uplink_config_request_access_with_passphrase(config, satellite_addr, api_key, NULL);
        require_error(access_result.error, UPLINK_ERROR_INTERNAL);
        require(access_result.access == NULL);
        uplink_free_access_result(access_result);

        access_result = uplink_config_request_access_with_passphrase(config, satellite_addr, api_key, "mypassphrase");
        require_noerror(access_result.error);
        require(access_result.access != NULL);

        UplinkProjectResult project_result = uplink_config_open_project(config, access_result.access);
        require_noerror(project_result.error);
        require(project_result.project != NULL);
        // check if project can be used to call satellite
        UplinkBucketResult bucket_result = uplink_stat_bucket(project_result.project, "not-existing-bucket");
        require_error(bucket_result.error, UPLINK_ERROR_BUCKET_NOT_FOUND);
        require(bucket_result.bucket == NULL);
        uplink_free_bucket_result(bucket_result);
        uplink_free_project_result(project_result);

        uplink_free_access_result(access_result);
    }

    {
        UplinkAccessResult access_result = uplink_parse_access(access_string);
        require_noerror(access_result.error);
        require(access_result.access != NULL);

        UplinkProjectResult project_result = uplink_config_open_project(config, access_result.access);
        require_noerror(project_result.error);
        require(project_result.project != NULL);
        uplink_free_project_result(project_result);

        project_result = uplink_config_open_project(config, NULL);
        require_error(project_result.error, UPLINK_ERROR_INTERNAL);
        require(project_result.project == NULL);
        uplink_free_project_result(project_result);

        uplink_free_access_result(access_result);
    }
    return 0;
}
