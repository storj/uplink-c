// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#include <limits.h>
#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "helpers.h"
#include "require.h"
#include "uplink.h"

#define THREADS 7

void handle_project(UplinkProject *project);
void *handle_thread(void *arg);

typedef struct {
    int id;
    pthread_t tid;
    UplinkProject *project;
} cfg_t;

cfg_t cfgs[THREADS];

int main()
{
    with_test_project(&handle_project);
    return 0;
}

void handle_project(UplinkProject *project)
{
    fprintf(stderr, "Handling project...\n");

    {
        UplinkBucketResult bucket_result = uplink_ensure_bucket(project, "alpha");
        require_noerror(bucket_result.error);
        uplink_free_bucket_result(bucket_result);
    }

    for (int i = 0; i < THREADS; i++) {
        cfgs[i].id = i;
        cfgs[i].project = project;

        pthread_create(&cfgs[i].tid, NULL, handle_thread, (void *)&cfgs[i]);
    }

    fprintf(stderr, "Waiting for threads...\n");

    for (int i = 0; i < THREADS; i++) {
        pthread_join(cfgs[i].tid, NULL);
    }

    fprintf(stderr, "Leaving project...\n");
}

void *handle_thread(void *arg)
{
    cfg_t *cfg = (cfg_t *)arg;

    fprintf(stderr, "[%d] Handling thread...\n", cfg->id);

    UplinkProject *project = cfg->project;

    time_t current_time = time(NULL);
    srandom(current_time + cfg->id);

    size_t data_len = 5 * 1024; // 5KiB;
    uint8_t *data = malloc(data_len);
    fill_random_data(data, data_len);

    char object_key[PATH_MAX];
    snprintf(object_key, PATH_MAX, "thread-%d/data.txt", cfg->id);

    fprintf(stderr, "[%d] Basic Upload...\n", cfg->id);

    { // basic upload
        UplinkUploadResult upload_result = uplink_upload_object(project, "alpha", object_key, NULL);
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

        UplinkCustomMetadataEntry entries[] = {
            {.key = "key1", .key_length = 4, .value = "value1", .value_length = 6},
            {.key = "key2", .key_length = 4, .value = "value2", .value_length = 6},
        };
        UplinkCustomMetadata customMetadata = {
            .entries = entries,
            .count = 2,
        };
        UplinkError *error = uplink_upload_set_custom_metadata(upload, customMetadata);
        require_noerror(error);

        UplinkError *commit_err = uplink_upload_commit(upload);
        require_noerror(commit_err);

        uplink_free_upload_result(upload_result);
    }

    fprintf(stderr, "[%d] Basic Download...\n", cfg->id);

    { // basic download
        size_t downloaded_len = (int64_t)data_len * 2;
        uint8_t *downloaded_data = malloc(downloaded_len);

        UplinkDownloadResult download_result = uplink_download_object(project, "alpha", object_key, NULL);
        require_noerror(download_result.error);
        require(download_result.download->_handle != 0);

        UplinkDownload *download = download_result.download;

        UplinkObjectResult object_result = uplink_download_info(download);
        require_noerror(object_result.error);
        require(object_result.object != NULL);

        UplinkObject *object = object_result.object;
        require(strcmp(object_key, object->key) == 0);
        require(object->system.created >= current_time);
        require(object->system.expires == 0);
        require(object->system.content_length == (int64_t)data_len);
        require(object->custom.count == 2);
        require(strcmp(object->custom.entries[0].key, "key1") == 0);
        require(strcmp(object->custom.entries[0].value, "value1") == 0);
        require(strcmp(object->custom.entries[1].key, "key2") == 0);
        require(strcmp(object->custom.entries[1].value, "value2") == 0);

        size_t downloaded_total = 0;
        while (true) {
            UplinkReadResult result =
                uplink_download_read(download, downloaded_data + downloaded_total, downloaded_len - downloaded_total);
            downloaded_total += result.bytes_read;

            if (result.error) {
                if (result.error->code == EOF) {
                    uplink_free_read_result(result);
                    break;
                }
                require_noerror(result.error);
            }
            uplink_free_read_result(result);
        }

        UplinkError *close_err = uplink_close_download(download);
        require_noerror(close_err);

        uplink_free_object_result(object_result);

        uplink_free_download_result(download_result);

        require(downloaded_total == data_len);
        require(memcmp(data, downloaded_data, data_len) == 0);
    }

    fprintf(stderr, "[%d] Stat Object...\n", cfg->id);

    { // stat object
        UplinkObjectResult object_result = uplink_stat_object(project, "alpha", object_key);
        require_noerror(object_result.error);
        require(object_result.object != NULL);

        UplinkObject *object = object_result.object;
        require(strcmp(object_key, object->key) == 0);
        require(object->system.created >= current_time);
        require(object->system.expires == 0);
        require(object->system.content_length == (int64_t)data_len);
        require(object->custom.count == 2);
        require(strcmp(object->custom.entries[0].key, "key1") == 0);
        require(strcmp(object->custom.entries[0].value, "value1") == 0);
        require(strcmp(object->custom.entries[1].key, "key2") == 0);
        require(strcmp(object->custom.entries[1].value, "value2") == 0);

        uplink_free_object_result(object_result);
    }

    fprintf(stderr, "[%d] Delete Object...\n", cfg->id);

    { // deleting an existing object
        UplinkObjectResult object_result = uplink_delete_object(project, "alpha", object_key);
        require_noerror(object_result.error);
        require(object_result.object != NULL);
        uplink_free_object_result(object_result);
    }

    { // deleting a missing object
        UplinkObjectResult object_result = uplink_delete_object(project, "alpha", object_key);
        require_error(object_result.error, UPLINK_ERROR_OBJECT_NOT_FOUND);
        require(object_result.object == NULL);
        uplink_free_object_result(object_result);
    }

    fprintf(stderr, "[%d] Leaving thread.\n", cfg->id);

    return NULL;
}
