// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "helpers.h"
#include "require.h"
#include "uplink.h"

void test_access_share(Access *access)
{
    {
        Permission emptyPermission = {};
        SharePrefix emptyPrefixes[] = {};
        AccessResult shared_access_result = access_share(NULL, emptyPermission, emptyPrefixes, 0);
        require_error(shared_access_result.error, ERROR_INTERNAL);
        require(shared_access_result.access == NULL);
        free_access_result(shared_access_result);

        shared_access_result = access_share(access, emptyPermission, emptyPrefixes, 0);
        require_error(shared_access_result.error, ERROR_INTERNAL);
        require(shared_access_result.access == NULL);
        free_access_result(shared_access_result);

        Permission permission = {
            allow_upload : true,
        };
        shared_access_result = access_share(access, permission, emptyPrefixes, 0);
        require_noerror(shared_access_result.error);
        requiref(shared_access_result.access->_handle != 0, "got empty access\n");
        free_access_result(shared_access_result);
    }

    size_t data_len = 1024;

    {
        Permission uploadPermission = {
            allow_upload : true,
        };

        SharePrefix prefixes[] = {
            {"alpha", ""},
        };
        AccessResult shared_access_result = access_share(access, uploadPermission, prefixes, 1);
        require_noerror(shared_access_result.error);
        require(shared_access_result.access != NULL);

        ProjectResult project_result = open_project(shared_access_result.access);
        require_noerror(project_result.error);
        requiref(project_result.project->_handle != 0, "got empty project\n");

        BucketResult bucket_result = create_bucket(project_result.project, "alpha");
        require_noerror(bucket_result.error);
        free_bucket_result(bucket_result);

        uint8_t *data = malloc(data_len);
        fill_random_data(data, data_len);

        UploadResult upload_result = upload_object(project_result.project, "alpha", "data.txt", NULL);
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

        free_access_result(shared_access_result);
        free_project_result(project_result);
    }

    {
        Permission uploadPermission = {
            allow_download : true,
        };

        SharePrefix prefixes[] = {
            // TODO find out why its not working with {"alpha", "data.txt"},
            // issue is most probably in encryption access
            {"alpha", ""},
        };
        AccessResult shared_access_result = access_share(access, uploadPermission, prefixes, 1);
        require_noerror(shared_access_result.error);
        require(shared_access_result.access != NULL);

        ProjectResult project_result = open_project(shared_access_result.access);
        require_noerror(project_result.error);
        requiref(project_result.project->_handle != 0, "got empty project\n");

        size_t downloaded_len = data_len * 2;
        uint8_t *downloaded_data = malloc(downloaded_len);

        DownloadResult download_result = download_object(project_result.project, "alpha", "data.txt", NULL);
        require_noerror(download_result.error);
        require(download_result.download->_handle != 0);

        Download *download = download_result.download;

        size_t downloaded_total = 0;
        while (true) {
            ReadResult result =
                download_read(download, downloaded_data + downloaded_total, downloaded_len - downloaded_total);
            downloaded_total += result.bytes_read;

            if (result.error) {
                if (result.error->code == EOF) {
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

        free_access_result(shared_access_result);
        free_project_result(project_result);
    }
}

int main(int argc, char *argv[])
{
    char *access_string = getenv("UPLINK_0_ACCESS");

    AccessResult access_result = parse_access(access_string);
    require_noerror(access_result.error);
    require(access_result.access != NULL);

    Access *access = access_result.access;
    StringResult serialized = access_serialize(access);
    require_noerror(serialized.error);
    require(serialized.string != NULL);

    require(strcmp(access_string, serialized.string) == 0);

    // test access share function
    test_access_share(access);

    free_access_result(access_result);

    requiref(internal_UniverseIsEmpty(), "universe is not empty\n");

    return 0;
}