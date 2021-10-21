#pragma once

// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#ifndef UPLINK_DISABLE_NAMESPACE_COMPAT
#warning                                                                                                               \
    "Please consider updating your code to use the namespaced symbols. These compatibility defines will be removed in a future release. Define UPLINK_DISABLE_NAMESPACE_COMPAT to remove this warning."

#define const_char uplink_const_char

#define Handle UplinkHandle
#define Access UplinkAccess
#define Project UplinkProject
#define Download UplinkDownload
#define Upload UplinkUpload
#define EncryptionKey UplinkEncryptionKey
#define Config UplinkConfig
#define Bucket UplinkBucket
#define SystemMetadata UplinkSystemMetadata
#define CustomMetadataEntry UplinkCustomMetadataEntry
#define CustomMetadata UplinkCustomMetadata
#define Object UplinkObject
#define UploadOptions UplinkUploadOptions
#define DownloadOptions UplinkDownloadOptions
#define ListObjectsOptions UplinkListObjectsOptions
#define ListBucketsOptions UplinkListBucketsOptions
#define ObjectIterator UplinkObjectIterator
#define BucketIterator UplinkBucketIterator
#define Permission UplinkPermission
#define SharePrefix UplinkSharePrefix
#define Error UplinkError

#define ERROR_INTERNAL UPLINK_ERROR_INTERNAL
#define ERROR_CANCELED UPLINK_ERROR_CANCELED
#define ERROR_INVALID_HANDLE UPLINK_ERROR_INVALID_HANDLE
#define ERROR_TOO_MANY_REQUESTS UPLINK_ERROR_TOO_MANY_REQUESTS
#define ERROR_BANDWIDTH_LIMIT_EXCEEDED UPLINK_ERROR_BANDWIDTH_LIMIT_EXCEEDED
#define ERROR_BUCKET_NAME_INVALID UPLINK_ERROR_BUCKET_NAME_INVALID
#define ERROR_BUCKET_ALREADY_EXISTS UPLINK_ERROR_BUCKET_ALREADY_EXISTS
#define ERROR_BUCKET_NOT_EMPTY UPLINK_ERROR_BUCKET_NOT_EMPTY
#define ERROR_BUCKET_NOT_FOUND UPLINK_ERROR_BUCKET_NOT_FOUND
#define ERROR_OBJECT_KEY_INVALID UPLINK_ERROR_OBJECT_KEY_INVALID
#define ERROR_OBJECT_NOT_FOUND UPLINK_ERROR_OBJECT_NOT_FOUND
#define ERROR_UPLOAD_DONE UPLINK_ERROR_UPLOAD_DONE

#define AccessResult UplinkAccessResult
#define ProjectResult UplinkProjectResult
#define BucketResult UplinkBucketResult
#define ObjectResult UplinkObjectResult
#define UploadResult UplinkUploadResult
#define DownloadResult UplinkDownloadResult
#define WriteResult UplinkWriteResult
#define ReadResult UplinkReadResult
#define StringResult UplinkStringResult
#define EncryptionKeyResult UplinkEncryptionKeyResult

#define parse_access uplink_parse_access
#define request_access_with_passphrase uplink_request_access_with_passphrase
#define access_satellite_address uplink_access_satellite_address
#define access_serialize uplink_access_serialize
#define access_share uplink_access_share
#define access_override_encryption_key uplink_access_override_encryption_key
#define free_string_result uplink_free_string_result
#define free_access_result uplink_free_access_result
#define stat_bucket uplink_stat_bucket
#define create_bucket uplink_create_bucket
#define ensure_bucket uplink_ensure_bucket
#define delete_bucket uplink_delete_bucket
#define free_bucket_result uplink_free_bucket_result
#define free_bucket uplink_free_bucket
#define list_buckets uplink_list_buckets
#define bucket_iterator_next uplink_bucket_iterator_next
#define bucket_iterator_err uplink_bucket_iterator_err
#define bucket_iterator_item uplink_bucket_iterator_item
#define free_bucket_iterator uplink_free_bucket_iterator
#define config_request_access_with_passphrase uplink_config_request_access_with_passphrase
#define config_open_project uplink_config_open_project
#define download_object uplink_download_object
#define download_read uplink_download_read
#define download_info uplink_download_info
#define free_read_result uplink_free_read_result
#define close_download uplink_close_download
#define free_download_result uplink_free_download_result
#define derive_encryption_key uplink_derive_encryption_key
#define free_encryption_key_result uplink_free_encryption_key_result
#define free_error uplink_free_error
#define internal_universeisempty uplink_internal_universeisempty
#define stat_object uplink_stat_object
#define update_object_metadata uplink_update_object_metadata
#define delete_object uplink_delete_object
#define free_object_result uplink_free_object_result
#define free_object uplink_free_object
#define list_objects uplink_list_objects
#define object_iterator_next uplink_object_iterator_next
#define object_iterator_err uplink_object_iterator_err
#define object_iterator_item uplink_object_iterator_item
#define free_object_iterator uplink_free_object_iterator
#define open_project uplink_open_project
#define close_project uplink_close_project
#define free_project_result uplink_free_project_result
#define upload_object uplink_upload_object
#define upload_write uplink_upload_write
#define upload_commit uplink_upload_commit
#define upload_abort uplink_upload_abort
#define upload_info uplink_upload_info
#define upload_set_custom_metadata uplink_upload_set_custom_metadata
#define free_write_result uplink_free_write_result
#define free_upload_result uplink_free_upload_result

#endif // UPLINK_DISABLE_NAMESPACE_COMPAT
