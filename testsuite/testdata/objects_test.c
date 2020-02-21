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

void handle_project(Project *project) {
    BucketResult bucket_result = ensure_bucket(project, "test");
    require_noerror(bucket_result.error);
    free_bucket_result(bucket_result);

    char *object_names[] = {"alpha", "beta", "gamma", "delta", "iota", "kappa", "lambda"};
    int object_names_count = 7;

    {
        for(int i = 0; i < object_names_count; i++) {
            UploadResult upload_result = upload_object(project, "test", object_names[i], NULL);
            require_noerror(upload_result.error);
            Upload *upload = upload_result.upload;
            require(upload != NULL);

            uint8_t hello[] = "hello";
            WriteResult result = upload_write(upload, hello, 5);
            require_noerror(result.error);
            free_write_result(result);

            Error *commit_error = upload_commit(upload);
            require_noerror(commit_error);
            free_upload_result(upload_result);
        }
    }

    {
        ObjectIterator* it = list_objects(project, "test", NULL);
        require(it != NULL);

        int count = 0;
        while(object_iterator_next(it)){
            Object *object = object_iterator_item(it);
            require(object != NULL);
            printf("%s\n", object->key);
            free_object(object);
            count++;
        }
        Error *err = object_iterator_err(it);
        require_noerror(err);

        // TODO: verify names returned.
        require(object_names_count = count);

        free_object_iterator(it);
    }

    // TODO: add tests for metadata verification
    // TODO: test options fields
}