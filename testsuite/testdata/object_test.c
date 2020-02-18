// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#include <string.h>
#include <stdlib.h>

#include "require.h"
#include "uplink.h"
#include "helpers.h"

void handle_project(Project project);

int main(int argc, char *argv[]) {
    with_test_project(&handle_project);
}

void handle_project(Project project) {
    char *_err = "";
    char **err = &_err;

    {
        Bucket bucket = ensure_bucket(project, "alpha", err);
        require_noerror(*err);
        free_bucket(bucket);
    }

    size_t  data_len = 5 * 1024; // 5KiB;
    uint8_t *data = malloc(data_len);
    fill_random_data(data, data_len);

    { // basic upload
        Upload upload = upload_object(project, "alpha", "data.txt", err);
        require_noerror(*err);
        require(upload._handle != 0);

        size_t uploaded_total = 0;
        while(uploaded_total < data_len) {
            size_t data_written = upload_write(upload, (uint8_t*)data+uploaded_total, data_len-uploaded_total, err);
            require_noerror(*err);
            uploaded_total += data_written;
            require(data_written > 0);
        }

        upload_commit(upload, err);
        require_noerror(*err);

        free_upload(upload, err);
        require_noerror(*err);
    }

    { // basic download
        size_t downloaded_len = data_len * 2;
        uint8_t *downloaded_data = malloc(downloaded_len);

        Download download = download_object(project, "alpha", "data.txt", err);
        require_noerror(*err);
        require(download._handle != 0);

        size_t downloaded_total = 0;
        while(true) {
            size_t data_read = download_read(download, (uint8_t*)downloaded_data+downloaded_total, downloaded_len-downloaded_total, err);
            downloaded_total += data_read;
            // TODO: check for io.EOF
            if(err != NULL) {
                free(err);
                break;
            }
        }

        free_download(download, err);
        require_noerror(*err);

        require(downloaded_total == data_len);
        require(memcmp(data, downloaded_data, data_len) == 0);
    }

    { // stat object
        Object object = stat_object(project, "alpha", "data.txt", err);
        require_noerror(*err);

        require(strcmp("data.txt", object.key) == 0);
        require(object.info.created != 0);
        require(object.info.expires == 0);

        free_object(object);
        // TODO: verify other status
    }

    { // deleting an existing object
        delete_object(project, "alpha", "data.txt", err);
        require_noerror(*err);
    }

    { // deleting a missing object
        delete_object(project, "alpha", "data.txt", err);
        require_error(*err);
        free(err);
    }
}