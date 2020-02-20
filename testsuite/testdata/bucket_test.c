// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#include <string.h>
#include <stdlib.h>

#include "require.h"
#include "uplink.h"
#include "helpers.h"

void handle_project(Project *project);

int main(int argc, char *argv[]) {
    with_test_project(&handle_project);
    return 0;
}

const uint32_t SUCCESS = 0;
const uint32_t ERROR_EOF = 1;
const uint32_t ERROR_INTERNAL = 2;
const uint32_t ERROR_CANCELED = 3;
const uint32_t ERROR_INVALID_HANDLE = 4;
const uint32_t ERROR_ALREADY_EXISTS = 5;
const uint32_t ERROR_NOT_FOUND = 6;

void handle_project(Project *project) {
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
        require_error(bucket_result.error, ERROR_ALREADY_EXISTS);
        
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
        require_error(bucket_result.error, ERROR_NOT_FOUND);
        require(bucket_result.bucket == NULL);

        free_bucket_result(bucket_result);
    }

    {
        // deleting a bucket
        Error *err = delete_bucket(project, "alpha");
        require_noerror(err);
    }

    {
        // deleting a missing bucket
        Error *err = delete_bucket(project, "missing");
        require_error(err, ERROR_NOT_FOUND);
        free_error(err);
    }

}