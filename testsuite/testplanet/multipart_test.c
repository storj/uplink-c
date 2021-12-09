// Copyright (C) 2021 Storj Labs, Inc.
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
    {
        UplinkBucketResult bucket_result = uplink_ensure_bucket(project, "alpha");
        require_noerror(bucket_result.error);
        uplink_free_bucket_result(bucket_result);
    }
    {
        UplinkBucketResult bucket_result = uplink_ensure_bucket(project, "alpha-second");
        require_noerror(bucket_result.error);
        uplink_free_bucket_result(bucket_result);
    }
    {
        UplinkBucketResult bucket_result = uplink_ensure_bucket(project, "alpha-third");
        require_noerror(bucket_result.error);
        uplink_free_bucket_result(bucket_result);
    }

    time_t current_time = time(NULL);

    { // test upload
        size_t data_len = 50 * 1024;
        uint8_t *data = malloc(data_len);
        fill_random_data(data, data_len);

        UplinkUploadInfoResult info_result = uplink_begin_upload(project, "alpha", "data.txt", NULL);
        require_noerror(info_result.error);
        require(info_result.info != NULL);
        require(strlen(info_result.info->upload_id) > 0);

        size_t uploaded_total = 0;
        for (size_t i = 0; i < 5; i++) {
            UplinkPartUploadResult upload_result =
                uplink_upload_part(project, "alpha", "data.txt", info_result.info->upload_id, i + 1);
            require_noerror(upload_result.error);
            require(upload_result.part_upload->_handle != 0);

            UplinkPartUpload *part_upload = upload_result.part_upload;

            size_t part_size = (i + 1) * 10240;
            while (uploaded_total < part_size) {
                UplinkWriteResult result =
                    uplink_part_upload_write(part_upload, data + uploaded_total, part_size - uploaded_total);
                uploaded_total += result.bytes_written;
                require_noerror(result.error);
                require(result.bytes_written > 0);
                uplink_free_write_result(result);
            }

            UplinkError *etag_err = uplink_part_upload_set_etag(part_upload, "test");
            require_noerror(etag_err);

            UplinkError *commit_err = uplink_part_upload_commit(part_upload);
            require_noerror(commit_err);

            UplinkPartResult info_result = uplink_part_upload_info(part_upload);
            require_noerror(info_result.error);
            require(info_result.part->part_number == (i + 1));
            require(info_result.part->size == 10240);
            // require(info_result.part->modified != 0); TODO enable when it will be fixed
            require(strcmp(info_result.part->etag, "test") == 0);
            require(info_result.part->etag_length == 4);

            uplink_free_part_result(info_result);
            uplink_free_part_upload_result(upload_result);
        }

        UplinkCommitUploadResult commit_result =
            uplink_commit_upload(project, "alpha", "data.txt", info_result.info->upload_id, NULL);
        require_noerror(commit_result.error);
        require(commit_result.object != NULL);
        require(strcmp(commit_result.object->key, "data.txt") == 0);

        UplinkObjectResult object_result = uplink_stat_object(project, "alpha", "data.txt");
        require_noerror(object_result.error);
        require(object_result.object != NULL);
        require(strcmp(object_result.object->key, "data.txt") == 0);
        require(object_result.object->system.content_length == (int64_t)data_len);

        uplink_free_upload_info_result(info_result);
        uplink_free_commit_upload_result(commit_result);
        uplink_free_object_result(object_result);
    }

    { // test abort pending object
        UplinkUploadInfoResult info_result = uplink_begin_upload(project, "alpha", "pending-data.txt", NULL);
        require_noerror(info_result.error);
        require(info_result.info != NULL);
        require(strlen(info_result.info->upload_id) > 0);

        UplinkError *abort_error =
            uplink_abort_upload(project, "alpha", "pending-data.txt", info_result.info->upload_id);
        require_noerror(abort_error);
        uplink_free_error(abort_error);

        // TODO I expect this to pass, we need to fix uplink
        // UplinkCommitUploadResult commit_result =
        //     uplink_commit_upload(project, "alpha", "pending-data.txt", info_result.info->upload_id, NULL);
        // require_error(commit_result.error, UPLINK_ERROR_OBJECT_NOT_FOUND);

        UplinkObjectResult object_result = uplink_stat_object(project, "alpha", "pending-data.txt");
        require_error(object_result.error, UPLINK_ERROR_OBJECT_NOT_FOUND);

        uplink_free_upload_info_result(info_result);
        //  uplink_free_commit_upload_result(commit_result);
        uplink_free_object_result(object_result);
    }

    { // test listing pending objects

        char *uploads[] = {"alpha/alpha_object", "beta", "delta", "gamma", "iota", "kappa", "lambda"};
        for (int i = 0; i < 7; i++) {
            UplinkUploadInfoResult info_result = uplink_begin_upload(project, "alpha-second", uploads[i], NULL);
            require_noerror(info_result.error);
            require(info_result.info != NULL);
            require(strlen(info_result.info->upload_id) > 0);
            uplink_free_upload_info_result(info_result);
        }

        {
#define EXPECTED_RESULTS_COUNT 7
            char *results[EXPECTED_RESULTS_COUNT];

            UplinkUploadIterator *it = uplink_list_uploads(project, "alpha-second", NULL);
            require(it != NULL);

            int count = 0;
            while (uplink_upload_iterator_next(it)) {
                UplinkUploadInfo *upload_info = uplink_upload_iterator_item(it);
                require(upload_info != NULL);

                bool is_prefix = upload_info->key[strlen(upload_info->key) - 1] == '/';
                require(upload_info->is_prefix == is_prefix);
                require(upload_info->system.created == 0);
                require(upload_info->system.expires == 0);
                require(upload_info->custom.count == 0);

                results[count] = strdup(upload_info->key);

                uplink_free_upload_info(upload_info);
                count++;
            }

            require(EXPECTED_RESULTS_COUNT == count);

            char *expected_results[] = {"alpha/", "beta", "delta", "gamma", "iota", "kappa", "lambda"};
            qsort(results, count, sizeof(char *), cstring_cmp);
            for (int i = 0; i < EXPECTED_RESULTS_COUNT; i++) {
                require(strcmp(expected_results[i], results[i]) == 0);
            }

            uplink_free_upload_iterator(it);
#undef EXPECTED_RESULTS_COUNT
        }

        {
#define EXPECTED_RESULTS_COUNT 1
            UplinkListUploadsOptions options = {
                .prefix = "alpha/",
                .system = true,
                .custom = true,
            };

            UplinkUploadIterator *it = uplink_list_uploads(project, "alpha-second", &options);
            require(it != NULL);

            char *expected_results[] = {"alpha/alpha_object"};
            char *results[EXPECTED_RESULTS_COUNT];

            int count = 0;
            while (uplink_upload_iterator_next(it)) {
                UplinkUploadInfo *upload_info = uplink_upload_iterator_item(it);
                require(upload_info != NULL);

                bool is_prefix = upload_info->key[strlen(upload_info->key) - 1] == '/';
                require(upload_info->is_prefix == is_prefix);
                require(upload_info->system.created >= current_time);
                require(upload_info->system.expires == 0);
                require(upload_info->custom.count == 0);

                results[count] = strdup(upload_info->key);

                uplink_free_upload_info(upload_info);
                count++;
            }
            UplinkError *err = uplink_upload_iterator_err(it);
            require_noerror(err);

            require(EXPECTED_RESULTS_COUNT == count);
            require(strcmp(expected_results[0], results[0]) == 0);

            uplink_free_upload_iterator(it);
#undef EXPECTED_RESULTS_COUNT
        }
    }

    { // test listing upload parts
        size_t data_len = 70 * 1024;
        uint8_t *data = malloc(data_len);
        fill_random_data(data, data_len);

        UplinkUploadInfoResult info_result = uplink_begin_upload(project, "alpha-third", "data.txt", NULL);
        require_noerror(info_result.error);
        require(info_result.info != NULL);
        require(strlen(info_result.info->upload_id) > 0);

        size_t uploaded_total = 0;
        for (size_t i = 0; i < 7; i++) {
            UplinkPartUploadResult upload_result =
                uplink_upload_part(project, "alpha-third", "data.txt", info_result.info->upload_id, i + 1);
            require_noerror(upload_result.error);
            require(upload_result.part_upload->_handle != 0);

            UplinkPartUpload *part_upload = upload_result.part_upload;

            size_t part_size = (i + 1) * 10240;
            while (uploaded_total < part_size) {
                UplinkWriteResult result =
                    uplink_part_upload_write(part_upload, data + uploaded_total, part_size - uploaded_total);
                uploaded_total += result.bytes_written;
                require_noerror(result.error);
                require(result.bytes_written > 0);
                uplink_free_write_result(result);
            }

            UplinkError *etag_err = uplink_part_upload_set_etag(part_upload, "my-etag");
            require_noerror(etag_err);

            UplinkError *commit_err = uplink_part_upload_commit(part_upload);
            require_noerror(commit_err);

            uplink_free_part_upload_result(upload_result);
        }

        { // no options
            UplinkPartIterator *it =
                uplink_list_upload_parts(project, "alpha-third", "data.txt", info_result.info->upload_id, NULL);
            require(it != NULL);

            int count = 0;
            while (uplink_part_iterator_next(it)) {
                UplinkPart *part = uplink_part_iterator_item(it);
                require(part != NULL);

                require(part->part_number == (uint32_t)count + 1);
                require(part->size == 10240);
                require(part->modified >= current_time);
                require(strcmp(part->etag, "my-etag") == 0);

                uplink_free_part(part);
                count++;
            }

            require(count == 7);

            uplink_free_part_iterator(it);
        }

        { // with options
            uint32_t cursor = (uint32_t)3;
            UplinkListUploadPartsOptions options = {
                .cursor = cursor,
            };

            UplinkPartIterator *it =
                uplink_list_upload_parts(project, "alpha-third", "data.txt", info_result.info->upload_id, &options);
            require(it != NULL);

            int count = 0;
            while (uplink_part_iterator_next(it)) {
                UplinkPart *part = uplink_part_iterator_item(it);
                require(part != NULL);

                require(part->part_number == (uint32_t)(count + 1 + cursor));
                require(part->size == 10240);
                require(part->modified >= current_time);
                require(strcmp(part->etag, "my-etag") == 0);

                uplink_free_part(part);
                count++;
            }

            require(count == 4);

            uplink_free_part_iterator(it);
        }

        uplink_free_upload_info_result(info_result);
    }
}
