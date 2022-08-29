// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#include <stdlib.h>
#include <string.h>

#include "../require.h"
#include "helpers.h"
#include "uplink.h"

void handle_project(UplinkProject *project);
void test_revoke_access(UplinkProject *project);
UplinkAccessResult derive_access(UplinkAccess *access);
void test_access_availability(UplinkAccess *access, bool expect_available);

int main()
{
    with_test_project(&handle_project);
    return 0;
}

void handle_project(UplinkProject *project)
{
    require(project->_handle != 0);
    test_revoke_access(project);
}

void test_revoke_access(UplinkProject *project)
{
    const char *access_string = getenv("UPLINK_0_ACCESS");

    UplinkAccessResult access_result = uplink_parse_access(access_string);
    require_noerror(access_result.error);
    require(access_result.access != NULL);

    UplinkBucketResult bucket_result = uplink_create_bucket(project, "alpha");
    require_noerror(bucket_result.error);
    uplink_free_bucket_result(bucket_result);

    UplinkAccess *access = access_result.access;

    {
        UplinkAccessResult derived_access_result = derive_access(access);
        test_access_availability(derived_access_result.access, true);

        UplinkError *revoke_error = uplink_revoke_access(NULL, derived_access_result.access);
        require_error(revoke_error, UPLINK_ERROR_INTERNAL);
        uplink_free_error(revoke_error);

        test_access_availability(derived_access_result.access, true);
        uplink_free_access_result(derived_access_result);
    }

    {
        UplinkAccessResult derived_access_result = derive_access(access);
        test_access_availability(derived_access_result.access, true);

        UplinkError *revoke_error = uplink_revoke_access(project, NULL);
        require_error(revoke_error, UPLINK_ERROR_INTERNAL);
        uplink_free_error(revoke_error);

        test_access_availability(derived_access_result.access, true);
        uplink_free_access_result(derived_access_result);
    }

    {
        UplinkAccessResult derived_access_result = derive_access(access);

        test_access_availability(derived_access_result.access, true);

        UplinkError *revoke_error = uplink_revoke_access(project, derived_access_result.access);
        require_noerror(revoke_error);
        uplink_free_error(revoke_error);

        test_access_availability(derived_access_result.access, false);
        uplink_free_access_result(derived_access_result);
    }

    uplink_free_access_result(access_result);
}

UplinkAccessResult derive_access(UplinkAccess *access)
{
    UplinkPermission full_permissions = {
        .allow_upload = true,
        .allow_download = true,
        .allow_delete = true,
        .allow_list = true,
    };

    UplinkSharePrefix prefixes[] = {
        {"alpha", ""},
    };
    UplinkAccessResult shared_access_result = uplink_access_share(access, full_permissions, prefixes, 1);
    require_noerror(shared_access_result.error);
    require(shared_access_result.access != NULL);

    return shared_access_result;
}

void test_access_availability(UplinkAccess *access, bool expect_available)
{
    size_t data_len = 1024;

    UplinkProjectResult project_result = uplink_open_project(access);
    require_noerror(project_result.error);
    requiref(project_result.project->_handle != 0, "got empty project\n");

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
    if (expect_available) {
        require_noerror(commit_err);
    } else {
        require_error(commit_err, UPLINK_ERROR_PERMISSION_DENIED);
    }

    uplink_free_upload_result(upload_result);
    uplink_free_project_result(project_result);
}
