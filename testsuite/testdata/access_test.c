// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#include <string.h>
#include <stdlib.h>
#include <stdio.h>

#include "uplink.h"
#include "require.h"
#include "helpers.h"

int main(int argc, char *argv[]) {
    char *access_string = getenv("UPLINK_0_ACCESS");

    AccessResult access_result = parse_access(access_string);
    require_noerror(access_result.error);
    require(access_result.access != NULL);

    Access *access = access_result.access;
    StringResult serialized = access_serialize(access);
    require_noerror(serialized.error);
    require(serialized.string != NULL);

    require(strcmp(access_string, serialized.string) == 0);

    free_access_result(access_result);

    return 0;
}