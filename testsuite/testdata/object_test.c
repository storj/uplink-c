// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#include <stdlib.h>
#include <string.h>

#include "helpers.h"
#include "require.h"
#include "uplink.h"

void handle_project(Project *project);

int main(int argc, char *argv[])
{
    with_test_project(&handle_project);
    return 0;
}

void handle_project(Project *project)
{
    {
        BucketResult bucket_result = ensure_bucket(project, "alpha");
        require_noerror(bucket_result.error);
        free_bucket_result(bucket_result);
    }

    size_t data_len = 5 * 1024; // 5KiB;
    uint8_t *data = malloc(data_len);
    fill_random_data(data, data_len);

    { // basic upload
        UploadResult upload_result = upload_object(project, "alpha", "data.txt", NULL);
        require_noerror(upload_result.error);
        require(upload_result.upload->_handle != 0);

        Upload *upload = upload_result.upload;

        size_t uploaded_total = 0;
        while (uploaded_total < data_len) {
            WriteResult result = upload_write(upload, data + uploaded_total, data_len - uploaded_total);
            uploaded_total += result.bytes_written;
            require_noerror(result.error);
            require(result.bytes_written > 0);
            free_write_result(result);
        }

        Error *commit_err = upload_commit(upload);
        require_noerror(commit_err);

        free_upload_result(upload_result);
    }

    { // basic download
        size_t downloaded_len = data_len * 2;
        uint8_t *downloaded_data = malloc(downloaded_len);

        DownloadResult download_result = download_object(project, "alpha", "data.txt", NULL);
        require_noerror(download_result.error);
        require(download_result.download->_handle != 0);

        Download *download = download_result.download;

        size_t downloaded_total = 0;
        while (true) {
            ReadResult result =
                download_read(download, downloaded_data + downloaded_total, downloaded_len - downloaded_total);
            downloaded_total += result.bytes_read;

            if (result.error) {
                if (result.error->code == ERROR_EOF) {
                    free_read_result(result);
                    break;
                }
                require_noerror(result.error);
            }
            free_read_result(result);
        }

        Error *close_err = close_download(download);
        require_noerror(close_err);

        free_download_result(download_result);

        require(downloaded_total == data_len);
        require(memcmp(data, downloaded_data, data_len) == 0);
    }

    { // stat object
        ObjectResult object_result = stat_object(project, "alpha", "data.txt");
        require_noerror(object_result.error);
        require(object_result.object != NULL);

        Object *object = object_result.object;
        require(strcmp("data.txt", object->key) == 0);
        require(object->system.created != 0);
        require(object->system.expires == 0);
        // TODO: verify other metadata fields

        free_object_result(object_result);
    }

    { // deleting an existing object
        ObjectResult object_result = delete_object(project, "alpha", "data.txt");
        require_noerror(object_result.error);
        require(object_result.object != NULL);
        free_object_result(object_result);
    }

    { // deleting a missing object
        ObjectResult object_result = delete_object(project, "alpha", "data.txt");
        require_error(object_result.error, ERROR_NOT_FOUND);
        require(object_result.object == NULL);
        free_object_result(object_result);
    }
}