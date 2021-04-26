// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

#include <stdlib.h>
#include <string.h>

#include "helpers.h"
#include "require.h"
#include "uplink.h"

void handle_project(UplinkProject *project);

int main()
{
    with_test_project(&handle_project);
    return 0;
}

void handle_project(UplinkProject *project)
{
    {
        UplinkBucketResult bucket_result = uplink_ensure_bucket(project, "alpha");
        require_noerror(bucket_result.error);
        uplink_free_bucket_result(bucket_result);
    }

    { // test upload
        UplinkUploadInfoResult info_result = uplink_begin_upload(project, "alpha", "data.txt", NULL);
        require_noerror(info_result.error);
        require(info_result.info != NULL);
        require(strlen(info_result.info->upload_id) > 0);

        UplinkCommitUploadResult commit_result =
            uplink_commit_upload(project, "alpha", "data.txt", info_result.info->upload_id, NULL);
        require_noerror(commit_result.error);
        require(commit_result.object != NULL);
        require(strcmp(commit_result.object->key, "data.txt") == 0);

        UplinkObjectResult object_result = uplink_stat_object(project, "alpha", "data.txt");
        require_noerror(object_result.error);
        require(object_result.object != NULL);
        require(strcmp(object_result.object->key, "data.txt") == 0);

        uplink_free_upload_info_result(info_result);
        uplink_free_commit_upload_result(commit_result);
        uplink_free_object_result(object_result);
    }

    { // test abort pending object
        UplinkUploadInfoResult info_result = uplink_begin_upload(project, "alpha", "pending-data.txt", NULL);
        require_noerror(info_result.error);
        require(info_result.info != NULL);
        require(strlen(info_result.info->upload_id) > 0);

        UplinkError *abort_error =
            uplink_abort_upload(project, "alpha", "pending-data.txt", info_result.info->upload_id);
        require_noerror(abort_error);
        uplink_free_error(abort_error);

        // TODO I expect this to pass, we need to fix uplink
        // UplinkCommitUploadResult commit_result =
        //     uplink_commit_upload(project, "alpha", "pending-data.txt", info_result.info->upload_id, NULL);
        // require_error(commit_result.error, UPLINK_ERROR_OBJECT_NOT_FOUND);

        UplinkObjectResult object_result = uplink_stat_object(project, "alpha", "pending-data.txt");
        require_error(object_result.error, UPLINK_ERROR_OBJECT_NOT_FOUND);

        uplink_free_upload_info_result(info_result);
        //  uplink_free_commit_upload_result(commit_result);
        uplink_free_object_result(object_result);
    }
}