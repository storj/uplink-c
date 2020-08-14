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
    const char *access_string = getenv("UPLINK_0_ACCESS");

    UplinkAccessResult access_result = uplink_parse_access(access_string);
    require_noerror(access_result.error);
    require(access_result.access != NULL);

    UplinkAccess *access = access_result.access;

    char salt[] = {4, 5, 6};
    UplinkEncryptionKeyResult key_result = uplink_derive_encryption_key("my-password", salt, 3);
    UplinkError *error = uplink_access_override_encryption_key(access, "bucket", "prefix/", key_result.encryption_key);
    require_noerror(error);

    uplink_free_access_result(access_result);
    uplink_free_encryption_key_result(key_result);

    requiref(uplink_internal_UniverseIsEmpty(), "universe is not empty\n");

    return 0;
}
