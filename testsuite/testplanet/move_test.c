// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

#include <stdlib.h>
#include <string.h>

#include "../require.h"
#include "helpers.h"
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

    size_t data_len = 5 * 1024; // 5KiB;
    uint8_t *data = malloc(data_len);
    fill_random_data(data, data_len);

    {
        UplinkUploadResult upload_result = uplink_upload_object(project, "alpha", "data.txt", NULL);
        require_noerror(upload_result.error);
        require(upload_result.upload->_handle != 0);

        UplinkUpload *upload = upload_result.upload;

        size_t uploaded_total = 0;
        while (uploaded_total < data_len) {
            UplinkWriteResult result = uplink_upload_write(upload, data + uploaded_total, data_len - uploaded_total);
            uploaded_total += result.bytes_written;
            require_noerror(result.error);
            require(result.bytes_written > 0);
            uplink_free_write_result(result);
        }

        UplinkError *commit_err = uplink_upload_commit(upload);
        require_noerror(commit_err);

        UplinkError *move_err = uplink_move_object(project, "alpha", "data.txt", "alpha", "a/prefix/data.txt", NULL);
        require_noerror(move_err);

        UplinkObjectResult object_result = uplink_stat_object(project, "alpha", "a/prefix/data.txt");
        require_noerror(object_result.error);
        require(object_result.object != NULL);

        uplink_free_object_result(object_result);

        uplink_free_upload_result(upload_result);
    }
}
