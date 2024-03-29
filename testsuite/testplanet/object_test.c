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

void handle_project(UplinkProject *project)
{
    {
        UplinkBucketResult bucket_result = uplink_ensure_bucket(project, "alpha");
        require_noerror(bucket_result.error);
        uplink_free_bucket_result(bucket_result);
    }

    time_t current_time = time(NULL);

    size_t data_len = 5 * 1024; // 5KiB;
    uint8_t *data = malloc(data_len);
    fill_random_data(data, data_len);

    { // basic upload
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

        UplinkObjectResult object_result = uplink_upload_info(upload);
        require_noerror(object_result.error);
        require(object_result.object != NULL);

        UplinkObject *object = object_result.object;
        require(strcmp("data.txt", object->key) == 0);
        require(object->system.created >= current_time);
        require(object->system.expires == 0);
        require(object->system.content_length == (int64_t)data_len);
        require(object->custom.count == 2);
        require(strcmp(object->custom.entries[0].key, "key1") == 0);
        require(strcmp(object->custom.entries[0].value, "value1") == 0);
        require(strcmp(object->custom.entries[1].key, "key2") == 0);
        require(strcmp(object->custom.entries[1].value, "value2") == 0);

        uplink_free_object_result(object_result);

        uplink_free_upload_result(upload_result);
    }

    { // basic download
        size_t downloaded_len = data_len * 2;
        uint8_t *downloaded_data = malloc(downloaded_len);

        UplinkDownloadResult download_result = uplink_download_object(project, "alpha", "data.txt", NULL);
        require_noerror(download_result.error);
        require(download_result.download->_handle != 0);

        UplinkDownload *download = download_result.download;

        UplinkObjectResult object_result = uplink_download_info(download);
        require_noerror(object_result.error);
        require(object_result.object != NULL);

        UplinkObject *object = object_result.object;
        require(strcmp("data.txt", object->key) == 0);
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

    { // stat object
        UplinkObjectResult object_result = uplink_stat_object(project, "alpha", "data.txt");
        require_noerror(object_result.error);
        require(object_result.object != NULL);

        UplinkObject *object = object_result.object;
        require(strcmp("data.txt", object->key) == 0);
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

    { // update object matadata.
        UplinkCustomMetadataEntry entries[] = {
            {.key = "key1", .key_length = 4, .value = "value11", .value_length = 7},
            {.key = "key2", .key_length = 4, .value = "value22", .value_length = 7},
            {.key = "key3", .key_length = 4, .value = "value33", .value_length = 7},
        };
        UplinkCustomMetadata customMetadata = {
            .entries = entries,
            .count = 3,
        };

        UplinkError *update_error = uplink_update_object_metadata(project, "alpha", "data.txt", customMetadata, NULL);
        require_noerror(update_error);

        UplinkObjectResult object_result = uplink_stat_object(project, "alpha", "data.txt");
        require_noerror(object_result.error);
        require(object_result.object != NULL);

        UplinkObject *object = object_result.object;
        require(strcmp("data.txt", object->key) == 0);
        require(object->system.created >= current_time);
        require(object->system.expires == 0);
        require(object->system.content_length == (int64_t)data_len);
        require(object->custom.count == 3);
        require(strcmp(object->custom.entries[0].key, "key1") == 0);
        require(strcmp(object->custom.entries[0].value, "value11") == 0);
        require(strcmp(object->custom.entries[1].key, "key2") == 0);
        require(strcmp(object->custom.entries[1].value, "value22") == 0);
        require(strcmp(object->custom.entries[2].key, "key3") == 0);
        require(strcmp(object->custom.entries[2].value, "value33") == 0);

        uplink_free_object_result(object_result);
    }

    { // deleting an existing object
        UplinkObjectResult object_result = uplink_delete_object(project, "alpha", "data.txt");
        require_noerror(object_result.error);
        require(object_result.object != NULL);
        uplink_free_object_result(object_result);
    }

    { // deleting a missing object
        UplinkObjectResult object_result = uplink_delete_object(project, "alpha", "data.txt");
        require_noerror(object_result.error);
        require(object_result.object == NULL);
        uplink_free_object_result(object_result);
    }
}
