// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

#include <stdlib.h>
#include <time.h>

// with_test_project opens default test project and calls handleProject callback.
void with_test_project(void (*handleProject)(Project*)) {
    // disable buffering
    setvbuf(stdout, NULL, _IONBF, 0);

    char *satellite_addr = getenv("SATELLITE_0_ADDR");
    char *api_key = getenv("UPLINK_0_APIKEY");
    char *access_string = getenv("UPLINK_0_ACCESS");
    //char *tmp_dir = getenv("TMP_DIR");

    printf("using SATELLITE_0_ADDR: %s\n", satellite_addr);
    printf("using UPLINK_0_ACCESS: %s\n", access_string);

    AccessResult access_result = request_access_with_passphrase(satellite_addr, api_key, "mypassphrase");
    require_noerror(access_result.error);

    ProjectResult project_result = open_project(access_result.access);
    require_noerror(project_result.error);
    requiref(project_result.project->_handle != 0, "got empty project\n");

    free_access_result(access_result);

    {
        handleProject(project_result.project);
    }

    Error *err = free_project_result(project_result);
    require_noerror(err);

    requiref(internal_UniverseIsEmpty(), "universe is not empty\n");
}

void fill_random_data(uint8_t *buffer, size_t length) {
     for(size_t i = 0; i < length; i++) {
          buffer[i] = (uint8_t)i*31;
     }
}

bool array_contains(char *item, char *array[], int array_size) {
    for (int i = 0; i < array_size; i++) {
        if(strcmp(array[i], item) == 0) {
            return true;
        }
    }

    return false;
}
