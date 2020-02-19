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
const uint32_t ERROR_EXISTS = 5;
const uint32_t ERROR_NOT_EXISTS = 6;

void handle_project(Project *project) {
    {
        BucketResult bucket_result = ensure_bucket(project, "alpha");
        xrequire_noerror(bucket_result.error);
        free_bucket_result(bucket_result);
    }

    size_t  data_len = 5 * 1024; // 5KiB;
    uint8_t *data = malloc(data_len);
    fill_random_data(data, data_len);

    { // basic upload
        UploadResult upload_result = upload_object(project, "alpha", "data.txt");
        xrequire_noerror(upload_result.error);
        require(upload_result.upload->_handle != 0);

        Upload *upload = upload_result.upload;

        size_t uploaded_total = 0;
        while(uploaded_total < data_len) {
            WriteResult result = upload_write(upload, (uint8_t*)data+uploaded_total, data_len-uploaded_total);
            uploaded_total += result.bytes_written;
            xrequire_noerror(result.error);
            require(result.bytes_written > 0);
        }

        Error *commit_err = upload_commit(upload);
        xrequire_noerror(commit_err);

        Error *free_err = free_upload_result(upload_result);
        xrequire_noerror(free_err);
    }

    { // basic download
        size_t downloaded_len = data_len * 2;
        uint8_t *downloaded_data = malloc(downloaded_len);

        DownloadResult download_result = download_object(project, "alpha", "data.txt");
        xrequire_noerror(download_result.error);
        require(download_result.download->_handle != 0);

        Download *download = download_result.download;

        size_t downloaded_total = 0;
        while(true) {
            ReadResult result = download_read(download, (uint8_t*)downloaded_data+downloaded_total, downloaded_len-downloaded_total);
            downloaded_total += result.bytes_read;

            if(result.error != NULL) {
                if(result.error->code == ERROR_EOF) {
                    free_error(result.error);
                    break;
                }
                xrequire_noerror(result.error);
            }
        }

        Error *free_err = free_download_result(download_result);
        xrequire_noerror(free_err);

        require(downloaded_total == data_len);
        require(memcmp(data, downloaded_data, data_len) == 0);
    }

    { // stat object
        ObjectResult object_result = stat_object(project, "alpha", "data.txt");
        xrequire_noerror(object_result.error);
        require(object_result.object != NULL);

        Object *object = object_result.object;
        require(strcmp("data.txt", object->key) == 0);
        require(object->info.created != 0);
        require(object->info.expires == 0);
        // TODO: verify other metadata fields

        free_object_result(object_result);
    }

    { // deleting an existing object
        Error *err = delete_object(project, "alpha", "data.txt");
        xrequire_noerror(err);
    }

    { // deleting a missing object
        Error *err = delete_object(project, "alpha", "data.txt");
        xrequire_error(err, ERROR_NOT_EXISTS);
        free_error(err);
    }
}