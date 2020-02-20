// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#include <string.h>
#include <stdlib.h>
#include <stdio.h>

#include "uplink.h"

void handle_project(Project *project) {
    {
        printf("# creating buckets\n");

        char *bucket_names[] = {"alpha", "beta", "gamma", "delta"};
        int bucket_names_count = 4;

        for(int i = 0; i < bucket_names_count; i++) {
            BucketResult bucket_result = ensure_bucket(project, bucket_names[i]);
            if(bucket_result.error){
                fprintf(stderr, "failed to create bucket %s: %s\n", bucket_names[i], bucket_result.error->message);
                free_bucket_result(bucket_result);
                return;
            }

            Bucket *bucket = bucket_result.bucket;
            fprintf(stdout, "created bucket %s\n", bucket->name);

            free_bucket_result(bucket_result);
        }
    }

    {
        printf("# listing buckets\n");

        BucketIterator* it = list_buckets(project, NULL);

        int count = 0;
        while(bucket_iterator_next(it)){
            Bucket *bucket = bucket_iterator_item(it);
            printf("bucket %s\n", bucket->name);
            free_bucket(bucket);
            count++;
        }
        Error *err = bucket_iterator_err(it);
        if(err){
            fprintf(stderr, "bucket listing failed: %s\n", err->message);
            free_error(err);
            free_bucket_iterator(it);
            return;
        }
        free_bucket_iterator(it);
    }
}

int main(int argc, char *argv[]) {
    char *access_string = getenv("UPLINK_0_ACCESS");

    AccessResult access_result = parse_access(access_string);
    if(access_result.error){
        fprintf(stderr, "failed to parse access: %s\n", access_result.error->message);
        goto done_access_result;
    }

    ProjectResult project_result = open_project(access_result.access);
    if(project_result.error){
        fprintf(stderr, "failed to open project: %s\n", project_result.error->message);
        goto done_project_result;
    }

    handle_project(project_result.project);

done_project_result:
    free_project_result(project_result);
done_access_result:
    free_access_result(access_result);

    return 0;
}