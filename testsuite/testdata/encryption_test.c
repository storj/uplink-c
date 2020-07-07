// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "helpers.h"
#include "require.h"
#include "uplink.h"

int main(int argc, char *argv[])
{
    char *access_string = getenv("UPLINK_0_ACCESS");

    AccessResult access_result = parse_access(access_string);
    require_noerror(access_result.error);
    require(access_result.access != NULL);

    Access *access = access_result.access;

    char salt[] = { 4, 5, 6 };
    EncryptionKeyResult key_result = derive_encryption_key("my-password", salt, 3);
    Error *error = access_override_encryption_key(access, "bucket", "prefix/", key_result.encryption_key);
    require_noerror(error);

    free_access_result(access_result);
    free_encryption_key_result(key_result);

    requiref(internal_UniverseIsEmpty(), "universe is not empty\n");

    return 0;
}