// Copyright (C) 2020 Storj Labs, Inc.
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

int cstring_cmp(const void *a, const void *b)
{
    const char **ia = (const char **)a;
    const char **ib = (const char **)b;
    return strcmp(*ia, *ib);
}

void handle_project(UplinkProject *project)
{
    UplinkBucketResult bucket_result = uplink_ensure_bucket(project, "test");
    require_noerror(bucket_result.error);
    uplink_free_bucket_result(bucket_result);

    time_t current_time = time(NULL);

    char *object_names[] = {"alpha/one", "beta", "delta", "gamma", "iota", "kappa", "lambda", "alpha/two"};
    const int object_names_count = 8;

    {
        for (int i = 0; i < object_names_count; i++) {
            UplinkUploadResult upload_result = uplink_upload_object(project, "test", object_names[i], NULL);
            require_noerror(upload_result.error);
            UplinkUpload *upload = upload_result.upload;
            require(upload != NULL);

            uint8_t hello[] = "hello";
            UplinkWriteResult result = uplink_upload_write(upload, hello, 5);
            require_noerror(result.error);
            uplink_free_write_result(result);

            UplinkCustomMetadataEntry entries[] = {
                {.key = "object_key",
                 .key_length = 10,
                 .value = object_names[i],
                 .value_length = strlen(object_names[i])},
            };
            UplinkCustomMetadata customMetadata = {.entries = entries, .count = 1};
            UplinkError *error = uplink_upload_set_custom_metadata(upload, customMetadata);
            require_noerror(error);

            UplinkError *commit_error = uplink_upload_commit(upload);
            require_noerror(commit_error);
            uplink_free_upload_result(upload_result);
        }
    }

    {
#define EXPECTED_RESULTS_COUNT 7
        UplinkObjectIterator *it = uplink_list_objects(project, "test", NULL);
        require(it != NULL);

        char *expected_results[] = {"alpha/", "beta", "delta", "gamma", "iota", "kappa", "lambda"};
        char *results[EXPECTED_RESULTS_COUNT];
        int count = 0;
        while (uplink_object_iterator_next(it)) {
            UplinkObject *object = uplink_object_iterator_item(it);
            require(object != NULL);
            bool is_prefix = object->key[strlen(object->key) - 1] == '/';
            require(object->is_prefix == is_prefix);
            require(object->system.created == 0);
            require(object->system.expires == 0);
            require(object->custom.count == 0);

            results[count] = strdup(object->key);

            uplink_free_object(object);
            count++;
        }
        UplinkError *err = uplink_object_iterator_err(it);
        require_noerror(err);

        require(EXPECTED_RESULTS_COUNT == count);

        qsort(results, count, sizeof(char *), cstring_cmp);
        for (int i = 0; i < EXPECTED_RESULTS_COUNT; i++) {
            require(strcmp(expected_results[i], results[i]) == 0);
        }

        uplink_free_object_iterator(it);
#undef EXPECTED_RESULTS_COUNT
    }

    {
#define EXPECTED_RESULTS_COUNT 2
        UplinkListObjectsOptions options = {
            .prefix = "alpha/",
            .system = true,
            .custom = true,
        };

        UplinkObjectIterator *it = uplink_list_objects(project, "test", &options);
        require(it != NULL);

        char *expected_results[] = {"alpha/one", "alpha/two"};
        char *results[EXPECTED_RESULTS_COUNT];

        int count = 0;
        while (uplink_object_iterator_next(it)) {
            UplinkObject *object = uplink_object_iterator_item(it);
            require(object != NULL);

            bool is_prefix = object->key[strlen(object->key) - 1] == '/';
            require(object->is_prefix == is_prefix);
            require(object->system.created >= current_time);
            require(object->system.expires == 0);
            require(object->custom.count == 1);
            require(strcmp(object->custom.entries[0].key, "object_key") == 0);
            require(strcmp(object->custom.entries[0].value, object->key) == 0);

            results[count] = strdup(object->key);

            uplink_free_object(object);
            count++;
        }
        UplinkError *err = uplink_object_iterator_err(it);
        require_noerror(err);

        require(EXPECTED_RESULTS_COUNT == count);

        qsort(results, count, sizeof(char *), cstring_cmp);
        for (int i = 0; i < EXPECTED_RESULTS_COUNT; i++) {
            require(strcmp(expected_results[i], results[i]) == 0);
        }

        uplink_free_object_iterator(it);
#undef EXPECTED_RESULTS_COUNT
    }
}
