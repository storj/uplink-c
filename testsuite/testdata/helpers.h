// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

#pragma once

#define UPLINK_DISABLE_NAMESPACE_COMPAT

#include <stdlib.h>
#include <time.h>

#include "require.h"
#include "uplink.h"

// with_test_project opens default test project and calls handleProject callback.
void with_test_project(void (*handleProject)(UplinkProject *))
{
    // disable buffering
    setvbuf(stdout, NULL, _IONBF, 0);

    const char *satellite_addr = getenv("SATELLITE_0_ADDR");
    const char *api_key = getenv("UPLINK_0_APIKEY");
    const char *access_string = getenv("UPLINK_0_ACCESS");
    // char *tmp_dir = getenv("TMP_DIR");

    printf("using SATELLITE_0_ADDR: %s\n", satellite_addr);
    printf("using UPLINK_0_ACCESS: %s\n", access_string);

    UplinkAccessResult access_result = uplink_request_access_with_passphrase(satellite_addr, api_key, "mypassphrase");
    require_noerror(access_result.error);

    UplinkProjectResult project_result = uplink_open_project(access_result.access);
    require_noerror(project_result.error);
    requiref(project_result.project->_handle != 0, "got empty project\n");

    uplink_free_access_result(access_result);

    {
        handleProject(project_result.project);
    }

    UplinkError *close_err = uplink_close_project(project_result.project);
    require_noerror(close_err);

    uplink_free_project_result(project_result);

    requiref(uplink_internal_UniverseIsEmpty(), "universe is not empty\n");
}

void fill_random_data(uint8_t *buffer, size_t length)
{
    for (size_t i = 0; i < length; i++) {
        buffer[i] = i * 31;
    }
}

bool array_contains(const char *item, const char *array[], int array_size)
{
    for (int i = 0; i < array_size; i++) {
        if (strcmp(array[i], item) == 0) {
            return true;
        }
    }

    return false;
}
