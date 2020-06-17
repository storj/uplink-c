// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "uplink.h"

void handle_project(Project *project)
{
    {
        printf("# creating buckets\n");

        char *bucket_names[] = {"alpha", "beta", "gamma", "delta"};
        int bucket_names_count = 4;

        for (int i = 0; i < bucket_names_count; i++) {
            BucketResult bucket_result = ensure_bucket(project, bucket_names[i]);
            if (bucket_result.error) {
                fprintf(stderr, "failed to create bucket %s: %s\n", bucket_names[i], bucket_result.error->message);
                free_bucket_result(bucket_result);
                return;
            }

            Bucket *bucket = bucket_result.bucket;
            fprintf(stdout, "created bucket %s\n", bucket->name);

            free_bucket_result(bucket_result);
        }
    }

    {
        printf("# listing buckets\n");

        BucketIterator *it = list_buckets(project, NULL);

        int count = 0;
        while (bucket_iterator_next(it)) {
            Bucket *bucket = bucket_iterator_item(it);
            printf("bucket %s\n", bucket->name);
            free_bucket(bucket);
            count++;
        }
        Error *err = bucket_iterator_err(it);
        if (err) {
            fprintf(stderr, "bucket listing failed: %s\n", err->message);
            free_error(err);
            free_bucket_iterator(it);
            return;
        }
        free_bucket_iterator(it);
    }

    {
        printf("# uploading objects\n");

        char *object_names[] = {"a.txt", "b/1.blob", "b/2.blob", "c.txt"};
        int object_names_count = 4;

        for (int i = 0; i < object_names_count; i++) {
            UploadResult upload_result = upload_object(project, "alpha", object_names[i], NULL);
            if (upload_result.error) {
                fprintf(stderr, "upload starting failed: %s\n", upload_result.error->message);
                free_upload_result(upload_result);
                return;
            }

            char *data = "testing data";
            size_t data_written = 0;
            size_t data_length = 12;

            Upload *upload = upload_result.upload;
            while (data_written < data_length) {
                WriteResult result = upload_write(upload, data + data_written, data_length - data_written);
                data_written += result.bytes_written;
                if (result.error) {
                    fprintf(stderr, "upload failed to write: %s\n", result.error->message);

                    Error *abort_error = upload_abort(upload);
                    if (abort_error) {
                        fprintf(stderr, "upload failed to abort: %s\n", abort_error->message);
                        free_error(abort_error);
                    }

                    free_write_result(result);
                    free_upload_result(upload_result);
                    return;
                }
                free_write_result(result);
            }

            Error *commit_error = upload_commit(upload);
            if (commit_error) {
                fprintf(stderr, "upload committing failed: %s\n", commit_error->message);
                free_error(commit_error);
                free_upload_result(upload_result);
                return;
            }
            free_upload_result(upload_result);
        }
    }

    {
        printf("# listing objects\n");

        ObjectIterator *it = list_objects(project, "alpha", NULL);

        int count = 0;
        while (object_iterator_next(it)) {
            Object *object = object_iterator_item(it);
            printf("object %s\n", object->key);
            free_object(object);
            count++;
        }
        Error *err = object_iterator_err(it);
        if (err) {
            fprintf(stderr, "object listing failed: %s\n", err->message);
            free_error(err);
            free_object_iterator(it);
            return;
        }
        free_object_iterator(it);
    }

    {
        printf("# downloading an object\n");

        DownloadResult download_result = download_object(project, "alpha", "a.txt", NULL);
        if (download_result.error) {
            fprintf(stderr, "download starting failed: %s\n", download_result.error->message);
            free_download_result(download_result);
            return;
        }

        size_t buffer_size = 1024;
        char *buffer = malloc(buffer_size);

        Download *download = download_result.download;
        while (true) {
            ReadResult result = download_read(download, buffer, buffer_size);

            // TODO: is there a nicer way to output a blob of binary data
            for (size_t p = 0; p < result.bytes_read; p++) {
                putchar(buffer[p]);
            }

            if (result.error) {
                if (result.error->code == EOF) {
                    free_read_result(result);
                    break;
                }
                fprintf(stderr, "download failed to read: %s\n", result.error->message);
                free_read_result(result);
                return;
            }
            free_read_result(result);
        }

        Error *close_error = close_download(download);
        if (close_error) {
            fprintf(stderr, "download failed to close: %s\n", close_error->message);
            free_error(close_error);
        }

        free_download_result(download_result);
    }

    {
        printf("# downloading an object range\n");

        DownloadOptions options = {};
        options.offset = 6;
        options.length = 3;

        DownloadResult download_result = download_object(project, "alpha", "a.txt", &options);
        if (download_result.error) {
            fprintf(stderr, "download starting failed: %s\n", download_result.error->message);
            free_download_result(download_result);
            return;
        }

        size_t buffer_size = 1024;
        char *buffer = malloc(buffer_size);

        Download *download = download_result.download;
        while (true) {
            ReadResult result = download_read(download, buffer, buffer_size);

            // TODO: is there a nicer way to output a blob of binary data
            for (size_t p = 0; p < result.bytes_read; p++) {
                putchar(buffer[p]);
            }

            if (result.error) {
                if (result.error->code == EOF) {
                    free_read_result(result);
                    break;
                }
                fprintf(stderr, "download failed to read: %s\n", result.error->message);
                free_read_result(result);
                return;
            }
            free_read_result(result);
        }

        Error *close_error = close_download(download);
        if (close_error) {
            fprintf(stderr, "download failed to close: %s\n", close_error->message);
            free_error(close_error);
        }

        free_download_result(download_result);
    }
}

int main(int argc, char *argv[])
{
    char *access_string = getenv("UPLINK_0_ACCESS");

    AccessResult access_result = parse_access(access_string);
    if (access_result.error) {
        fprintf(stderr, "failed to parse access: %s\n", access_result.error->message);
        goto done_access_result;
    }

    ProjectResult project_result = open_project(access_result.access);
    if (project_result.error) {
        fprintf(stderr, "failed to open project: %s\n", project_result.error->message);
        goto done_project_result;
    }

    handle_project(project_result.project);

    Error *close_error = close_project(project_result.project);
    if (close_error) {
        fprintf(stderr, "failed to close project: %s\n", close_error->message);
        free_error(close_error);
    }

done_project_result:
    free_project_result(project_result);
done_access_result:
    free_access_result(access_result);

    return 0;
}