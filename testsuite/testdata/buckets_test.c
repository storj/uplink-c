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
    // creating a buckets
    const char *bucket_names[] = {"alpha", "beta", "delta", "gamma", "iota", "kappa", "lambda"};
    int bucket_names_count = 7;

    time_t current_time = time(NULL);

    for (int i = 0; i < bucket_names_count; i++) {
        BucketResult bucket_result = ensure_bucket(project, bucket_names[i]);
        require_noerror(bucket_result.error);
        free_bucket_result(bucket_result);
    }

    {
        BucketIterator *it = list_buckets(project, NULL);
        require(it != NULL);

        int count = 0;
        while (bucket_iterator_next(it)) {
            Bucket *bucket = bucket_iterator_item(it);
            require(bucket != NULL);
            require(strcmp(bucket->name, bucket_names[count]) == 0);
            require(bucket->created >= current_time);
            free_bucket(bucket);
            count++;
        }

        Error *err = bucket_iterator_err(it);
        require_noerror(err);
        require(bucket_names_count = count);

        free_bucket_iterator(it);
    }

    {
        ListBucketsOptions options = {
            .cursor = "gamma",
        };

        BucketIterator *it = list_buckets(project, &options);
        require(it != NULL);

        int iota_index = 4;
        int count = 0;
        while (bucket_iterator_next(it)) {
            Bucket *bucket = bucket_iterator_item(it);
            require(bucket != NULL);
            require(strcmp(bucket->name, bucket_names[count + iota_index]) == 0);
            require(bucket->created != 0);
            free_bucket(bucket);
            count++;
        }

        Error *err = bucket_iterator_err(it);
        require_noerror(err);
        require(bucket_names_count = count);

        free_bucket_iterator(it);
    }
}
