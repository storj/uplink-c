// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "helpers.h"
#include "require.h"
#include "uplink.h"

void test_access_satellite_address(UplinkAccess *access)
{
    UplinkStringResult satellite_address_result = uplink_access_satellite_address(access);
    require_noerror(satellite_address_result.error);
    require(satellite_address_result.string != NULL);

    uplink_free_string_result(satellite_address_result);
}

void test_access_share(UplinkAccess *access)
{
    {
        UplinkPermission emptyPermission = {0};
        UplinkSharePrefix emptyPrefixes[] = {0};
        UplinkAccessResult shared_access_result = uplink_access_share(NULL, emptyPermission, emptyPrefixes, 0);
        require_error(shared_access_result.error, UPLINK_ERROR_INTERNAL);
        require(shared_access_result.access == NULL);
        uplink_free_access_result(shared_access_result);

        shared_access_result = uplink_access_share(access, emptyPermission, emptyPrefixes, 0);
        require_error(shared_access_result.error, UPLINK_ERROR_INTERNAL);
        require(shared_access_result.access == NULL);
        uplink_free_access_result(shared_access_result);

        UplinkPermission permission = {
            .allow_upload = true,
        };
        shared_access_result = uplink_access_share(access, permission, emptyPrefixes, 0);
        require_noerror(shared_access_result.error);
        requiref(shared_access_result.access->_handle != 0, "got empty access\n");
        uplink_free_access_result(shared_access_result);
    }

    size_t data_len = 1024;

    {
        UplinkPermission uploadPermission = {
            .allow_upload = true,
        };

        UplinkSharePrefix prefixes[] = {
            {"alpha", ""},
        };
        UplinkAccessResult shared_access_result = uplink_access_share(access, uploadPermission, prefixes, 1);
        require_noerror(shared_access_result.error);
        require(shared_access_result.access != NULL);

        UplinkProjectResult project_result = uplink_open_project(shared_access_result.access);
        require_noerror(project_result.error);
        requiref(project_result.project->_handle != 0, "got empty project\n");

        UplinkBucketResult bucket_result = uplink_create_bucket(project_result.project, "alpha");
        require_noerror(bucket_result.error);
        uplink_free_bucket_result(bucket_result);

        uint8_t *data = malloc(data_len);
        fill_random_data(data, data_len);

        UplinkUploadResult upload_result = uplink_upload_object(project_result.project, "alpha", "data.txt", NULL);
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

        UplinkError *commit_err = uplink_upload_commit(upload);
        require_noerror(commit_err);

        uplink_free_upload_result(upload_result);

        uplink_free_access_result(shared_access_result);
        uplink_free_project_result(project_result);
    }

    {
        UplinkPermission uploadPermission = {
            .allow_download = true,
        };

        UplinkSharePrefix prefixes[] = {
            // TODO find out why its not working with {"alpha", "data.txt"},
            // issue is most probably in encryption access
            {"alpha", ""},
        };
        UplinkAccessResult shared_access_result = uplink_access_share(access, uploadPermission, prefixes, 1);
        require_noerror(shared_access_result.error);
        require(shared_access_result.access != NULL);

        UplinkProjectResult project_result = uplink_open_project(shared_access_result.access);
        require_noerror(project_result.error);
        requiref(project_result.project->_handle != 0, "got empty project\n");

        size_t downloaded_len = data_len * 2;
        uint8_t *downloaded_data = malloc(downloaded_len);

        UplinkDownloadResult download_result =
            uplink_download_object(project_result.project, "alpha", "data.txt", NULL);
        require_noerror(download_result.error);
        require(download_result.download->_handle != 0);

        UplinkDownload *download = download_result.download;

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

        uplink_free_download_result(download_result);

        uplink_free_access_result(shared_access_result);
        uplink_free_project_result(project_result);
    }
}

int main()
{
    const char *access_string = getenv("UPLINK_0_ACCESS");

    UplinkAccessResult access_result = uplink_parse_access(access_string);
    require_noerror(access_result.error);
    require(access_result.access != NULL);

    UplinkAccess *access = access_result.access;
    UplinkStringResult serialized = uplink_access_serialize(access);
    require_noerror(serialized.error);
    require(serialized.string != NULL);

    require(strcmp(access_string, serialized.string) == 0);

    // test access satellite node url
    test_access_satellite_address(access);

    // test access share function
    test_access_share(access);

    uplink_free_access_result(access_result);

    requiref(uplink_internal_UniverseIsEmpty(), "universe is not empty\n");

    return 0;
}
