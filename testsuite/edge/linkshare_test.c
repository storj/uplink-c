// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

#include <stdio.h>
#include <string.h>

#include "../require.h"
#define UPLINK_DISABLE_NAMESPACE_COMPAT
#include "uplink.h"

int main()
{
    {
        // Happy
        UplinkStringResult result = edge_join_share_url("https://storjshare.example", "l5pucy3dmvzxgs3fpfewix27l5pq",
                                                        "mybucket", "myprefix/mykey", NULL);

        require_noerror(result.error);

        require(strcmp("https://storjshare.example/s/l5pucy3dmvzxgs3fpfewix27l5pq/mybucket/myprefix/mykey",
                       result.string) == 0);

        uplink_free_string_result(result);
    }

    {
        // Raw
        EdgeShareURLOptions options = {
            .raw = true,
        };

        UplinkStringResult result = edge_join_share_url("https://storjshare.example", "l5pucy3dmvzxgs3fpfewix27l5pq",
                                                        "mybucket", "myprefix/mykey", &options);

        require_noerror(result.error);

        require(strcmp("https://storjshare.example/raw/l5pucy3dmvzxgs3fpfewix27l5pq/mybucket/myprefix/mykey",
                       result.string) == 0);

        uplink_free_string_result(result);
    }

    {
        // Error prefix without bucket
        UplinkStringResult result = edge_join_share_url("https://storjshare.example", "l5pucy3dmvzxgs3fpfewix27l5pq",
                                                        "", "myprefix/mykey", NULL);

        require_error(result.error, UPLINK_ERROR_INTERNAL);

        uplink_free_string_result(result);
    }

    return 0;
}