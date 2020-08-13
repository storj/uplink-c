// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#include <stdlib.h>
#include <string.h>

#include "helpers.h"
#include "require.h"
#include "uplink.h"

void handle_project(Project *project);

int main()
{
    with_test_project(&handle_project);
    return 0;
}

void handle_project(Project *project)
{
    {
        // creating a new bucket
        BucketResult bucket_result = create_bucket(project, "alpha");
        require_noerror(bucket_result.error);

        Bucket *bucket = bucket_result.bucket;
        require(bucket != NULL);
        require(strcmp("alpha", bucket->name) == 0);
        require(bucket->created != 0);

        free_bucket_result(bucket_result);
    }

    {
        // creating an existing bucket
        BucketResult bucket_result = create_bucket(project, "alpha");
        require_error(bucket_result.error, ERROR_BUCKET_ALREADY_EXISTS);

        Bucket *bucket = bucket_result.bucket;
        require(bucket != NULL);
        require(strcmp("alpha", bucket->name) == 0);
        require(bucket->created != 0);

        free_bucket_result(bucket_result);
    }

    {
        // ensuring an existing bucket
        BucketResult bucket_result = ensure_bucket(project, "alpha");
        require_noerror(bucket_result.error);

        Bucket *bucket = bucket_result.bucket;
        require(bucket != NULL);
        require(strcmp("alpha", bucket->name) == 0);
        require(bucket->created != 0);

        free_bucket_result(bucket_result);
    }

    {
        // ensuring a new bucket
        BucketResult bucket_result = ensure_bucket(project, "beta");
        require_noerror(bucket_result.error);

        Bucket *bucket = bucket_result.bucket;
        require(bucket != NULL);
        require(strcmp("beta", bucket->name) == 0);
        require(bucket->created != 0);

        free_bucket_result(bucket_result);
    }

    {
        // statting a bucket
        BucketResult bucket_result = stat_bucket(project, "alpha");
        require_noerror(bucket_result.error);

        Bucket *bucket = bucket_result.bucket;
        require(bucket != NULL);
        require(strcmp("alpha", bucket->name) == 0);
        require(bucket->created != 0);

        free_bucket_result(bucket_result);
    }

    {
        // statting a missing bucket
        BucketResult bucket_result = stat_bucket(project, "missing");
        require_error(bucket_result.error, ERROR_BUCKET_NOT_FOUND);
        require(bucket_result.bucket == NULL);

        free_bucket_result(bucket_result);
    }

    {
        // deleting a bucket
        BucketResult bucket_result = delete_bucket(project, "alpha");
        require_noerror(bucket_result.error);
        require(bucket_result.bucket != NULL);
        free_bucket_result(bucket_result);
    }

    {
        // deleting a missing bucket
        BucketResult bucket_result = delete_bucket(project, "missing");
        require_error(bucket_result.error, ERROR_BUCKET_NOT_FOUND);
        require(bucket_result.bucket == NULL);
        free_bucket_result(bucket_result);
    }
}
