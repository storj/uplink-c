// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#include <stdlib.h>
#include <string.h>

#include "helpers.h"
#include "require.h"
#include "uplink.h"

void handle_project(UplinkProject *project);
void upload_object(UplinkProject *project, const char *bucket_name, const char *object_key, size_t data_length);

int main()
{
    with_test_project(&handle_project);
    return 0;
}

void handle_project(UplinkProject *project)
{
    {
        // creating a new bucket
        UplinkBucketResult bucket_result = uplink_create_bucket(project, "alpha");
        require_noerror(bucket_result.error);

        UplinkBucket *bucket = bucket_result.bucket;
        require(bucket != NULL);
        require(strcmp("alpha", bucket->name) == 0);
        require(bucket->created != 0);

        uplink_free_bucket_result(bucket_result);
    }

    {
        // creating an existing bucket
        UplinkBucketResult bucket_result = uplink_create_bucket(project, "alpha");
        require_error(bucket_result.error, UPLINK_ERROR_BUCKET_ALREADY_EXISTS);

        UplinkBucket *bucket = bucket_result.bucket;
        require(bucket != NULL);
        require(strcmp("alpha", bucket->name) == 0);
        require(bucket->created != 0);

        uplink_free_bucket_result(bucket_result);
    }

    {
        // ensuring an existing bucket
        UplinkBucketResult bucket_result = uplink_ensure_bucket(project, "alpha");
        require_noerror(bucket_result.error);

        UplinkBucket *bucket = bucket_result.bucket;
        require(bucket != NULL);
        require(strcmp("alpha", bucket->name) == 0);
        require(bucket->created != 0);

        uplink_free_bucket_result(bucket_result);
    }

    {
        // ensuring a new bucket
        UplinkBucketResult bucket_result = uplink_ensure_bucket(project, "beta");
        require_noerror(bucket_result.error);

        UplinkBucket *bucket = bucket_result.bucket;
        require(bucket != NULL);
        require(strcmp("beta", bucket->name) == 0);
        require(bucket->created != 0);

        uplink_free_bucket_result(bucket_result);
    }

    {
        // statting a bucket
        UplinkBucketResult bucket_result = uplink_stat_bucket(project, "alpha");
        require_noerror(bucket_result.error);

        UplinkBucket *bucket = bucket_result.bucket;
        require(bucket != NULL);
        require(strcmp("alpha", bucket->name) == 0);
        require(bucket->created != 0);

        uplink_free_bucket_result(bucket_result);
    }

    {
        // statting a missing bucket
        UplinkBucketResult bucket_result = uplink_stat_bucket(project, "missing");
        require_error(bucket_result.error, UPLINK_ERROR_BUCKET_NOT_FOUND);
        require(bucket_result.bucket == NULL);

        uplink_free_bucket_result(bucket_result);
    }

    {
        // deleting a bucket
        UplinkBucketResult bucket_result = uplink_delete_bucket(project, "alpha");
        require_noerror(bucket_result.error);
        require(bucket_result.bucket != NULL);
        uplink_free_bucket_result(bucket_result);
    }

    {
        // deleting a missing bucket
        UplinkBucketResult bucket_result = uplink_delete_bucket(project, "missing");
        require_error(bucket_result.error, UPLINK_ERROR_BUCKET_NOT_FOUND);
        require(bucket_result.bucket == NULL);
        uplink_free_bucket_result(bucket_result);
    }

    {
        // force deleting a bucket that's not empty
        UplinkBucketResult create_bucket_result = uplink_create_bucket(project, "alpha");
        require_noerror(create_bucket_result.error);

        UplinkBucket *bucket = create_bucket_result.bucket;
        require(bucket != NULL);
        require(strcmp("alpha", bucket->name) == 0);
        require(bucket->created != 0);

        uplink_free_bucket_result(create_bucket_result);

        size_t data_len = 5 * 1024; // 5KiB;
        upload_object(project, "alpha", "data.xtx", data_len);

        UplinkBucketResult delete_bucket_result = uplink_delete_bucket_with_objects(project, "alpha");
        require_noerror(delete_bucket_result.error);
        require(delete_bucket_result.bucket != NULL);
        uplink_free_bucket_result(delete_bucket_result);

        UplinkBucketResult stat_bucket_result = uplink_stat_bucket(project, "alpha");
        require_error(stat_bucket_result.error, UPLINK_ERROR_BUCKET_NOT_FOUND);
        require(stat_bucket_result.bucket == NULL);

        uplink_free_bucket_result(stat_bucket_result);
    }
}

void upload_object(UplinkProject *project, const char *bucket_name, const char *object_key, size_t data_length)
{
    uint8_t *data = malloc(data_length);
    fill_random_data(data, data_length);

    UplinkUploadResult upload_result = uplink_upload_object(project, bucket_name, object_key, NULL);
    require_noerror(upload_result.error);
    require(upload_result.upload->_handle != 0);

    UplinkUpload *upload = upload_result.upload;

    size_t uploaded_total = 0;
    while (uploaded_total < data_length) {
        UplinkWriteResult result = uplink_upload_write(upload, data + uploaded_total, data_length - uploaded_total);
        uploaded_total += result.bytes_written;
        require_noerror(result.error);
        require(result.bytes_written > 0);
        uplink_free_write_result(result);
    }

    UplinkError *commit_err = uplink_upload_commit(upload);
    require_noerror(commit_err);

    UplinkObjectResult object_result = uplink_upload_info(upload);
    require_noerror(object_result.error);
    require(object_result.object != NULL);

    uplink_free_object_result(object_result);

    uplink_free_upload_result(upload_result);
}
