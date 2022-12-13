#pragma once

// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

#include "uplink_compat.h"

typedef const char uplink_const_char;

typedef struct UplinkHandle {
    size_t _handle;
} UplinkHandle;

typedef struct UplinkAccess {
    size_t _handle;
} UplinkAccess;

typedef struct UplinkProject {
    size_t _handle;
} UplinkProject;

typedef struct UplinkDownload {
    size_t _handle;
} UplinkDownload;

typedef struct UplinkUpload {
    size_t _handle;
} UplinkUpload;

typedef struct UplinkEncryptionKey {
    size_t _handle;
} UplinkEncryptionKey;

typedef struct UplinkPartUpload {
    size_t _handle;
} UplinkPartUpload;

typedef struct UplinkConfig {
    const char *user_agent;

    int32_t dial_timeout_milliseconds;

    // temp_directory specifies where to save data during downloads to use less memory.
    const char *temp_directory;
} UplinkConfig;

typedef struct UplinkBucket {
    char *name;
    int64_t created;
} UplinkBucket;

typedef struct UplinkSystemMetadata {
    int64_t created;
    int64_t expires;
    int64_t content_length;
} UplinkSystemMetadata;

typedef struct UplinkCustomMetadataEntry {
    char *key;
    size_t key_length;

    char *value;
    size_t value_length;
} UplinkCustomMetadataEntry;

typedef struct UplinkCustomMetadata {
    UplinkCustomMetadataEntry *entries;
    size_t count;
} UplinkCustomMetadata;

typedef struct UplinkObject {
    char *key;
    bool is_prefix;
    UplinkSystemMetadata system;
    UplinkCustomMetadata custom;
} UplinkObject;

typedef struct UplinkUploadOptions {
    // When expires is 0 or negative, it means no expiration.
    int64_t expires;
} UplinkUploadOptions;

typedef struct UplinkDownloadOptions {
    // When offset is negative it will read the suffix of the blob.
    // Combining negative offset and positive length is not supported.
    int64_t offset;
    // When length is negative, it will read until the end of the blob.
    int64_t length;
} UplinkDownloadOptions;

typedef struct UplinkListObjectsOptions {
    const char *prefix;
    const char *cursor;
    bool recursive;

    bool system;
    bool custom;
} UplinkListObjectsOptions;

typedef struct UplinkListUploadsOptions {
    const char *prefix;
    const char *cursor;
    bool recursive;

    bool system;
    bool custom;
} UplinkListUploadsOptions;

typedef struct UplinkListBucketsOptions {
    const char *cursor;
} UplinkListBucketsOptions;

typedef struct UplinkObjectIterator {
    size_t _handle;
} UplinkObjectIterator;

typedef struct UplinkBucketIterator {
    size_t _handle;
} UplinkBucketIterator;

typedef struct UplinkUploadIterator {
    size_t _handle;
} UplinkUploadIterator;

typedef struct UplinkPartIterator {
    size_t _handle;
} UplinkPartIterator;

typedef struct UplinkPermission {
    bool allow_download;
    bool allow_upload;
    bool allow_list;
    bool allow_delete;

    // unix time in seconds when the permission becomes valid.
    // disabled when 0.
    int64_t not_before;
    // unix time in seconds when the permission becomes invalid.
    // disabled when 0.
    int64_t not_after;
} UplinkPermission;

typedef struct UplinkPart {
    uint32_t part_number;
    size_t size; // plain size of a part.
    int64_t modified;
    char *etag;
    size_t etag_length;
} UplinkPart;

typedef struct UplinkSharePrefix {
    const char *bucket;
    // prefix is the prefix of the shared object keys.
    const char *prefix;
} UplinkSharePrefix;

typedef struct UplinkError {
    int32_t code;
    char *message;
} UplinkError;

#define UPLINK_ERROR_INTERNAL 0x02
#define UPLINK_ERROR_CANCELED 0x03
#define UPLINK_ERROR_INVALID_HANDLE 0x04
#define UPLINK_ERROR_TOO_MANY_REQUESTS 0x05
#define UPLINK_ERROR_BANDWIDTH_LIMIT_EXCEEDED 0x06
#define UPLINK_ERROR_STORAGE_LIMIT_EXCEEDED 0x07
#define UPLINK_ERROR_SEGMENTS_LIMIT_EXCEEDED 0x08
#define UPLINK_ERROR_PERMISSION_DENIED 0x09

#define UPLINK_ERROR_BUCKET_NAME_INVALID 0x10
#define UPLINK_ERROR_BUCKET_ALREADY_EXISTS 0x11
#define UPLINK_ERROR_BUCKET_NOT_EMPTY 0x12
#define UPLINK_ERROR_BUCKET_NOT_FOUND 0x13

#define UPLINK_ERROR_OBJECT_KEY_INVALID 0x20
#define UPLINK_ERROR_OBJECT_NOT_FOUND 0x21
#define UPLINK_ERROR_UPLOAD_DONE 0x22

#define EDGE_ERROR_AUTH_DIAL_FAILED 0x30
#define EDGE_ERROR_REGISTER_ACCESS_FAILED 0x31

typedef struct UplinkAccessResult {
    UplinkAccess *access;
    UplinkError *error;
} UplinkAccessResult;

typedef struct UplinkProjectResult {
    UplinkProject *project;
    UplinkError *error;
} UplinkProjectResult;

typedef struct UplinkBucketResult {
    UplinkBucket *bucket;
    UplinkError *error;
} UplinkBucketResult;

typedef struct UplinkObjectResult {
    UplinkObject *object;
    UplinkError *error;
} UplinkObjectResult;

typedef struct UplinkUploadResult {
    UplinkUpload *upload;
    UplinkError *error;
} UplinkUploadResult;

typedef struct UplinkPartUploadResult {
    UplinkPartUpload *part_upload;
    UplinkError *error;
} UplinkPartUploadResult;

typedef struct UplinkDownloadResult {
    UplinkDownload *download;
    UplinkError *error;
} UplinkDownloadResult;

typedef struct UplinkWriteResult {
    size_t bytes_written;
    UplinkError *error;
} UplinkWriteResult;

typedef struct UplinkReadResult {
    size_t bytes_read;
    UplinkError *error;
} UplinkReadResult;

typedef struct UplinkStringResult {
    char *string;
    UplinkError *error;
} UplinkStringResult;

typedef struct UplinkEncryptionKeyResult {
    UplinkEncryptionKey *encryption_key;
    UplinkError *error;
} UplinkEncryptionKeyResult;

typedef struct UplinkUploadInfo {
    char *upload_id;

    char *key;
    bool is_prefix;
    UplinkSystemMetadata system;
    UplinkCustomMetadata custom;
} UplinkUploadInfo;

typedef struct UplinkUploadInfoResult {
    UplinkUploadInfo *info;
    UplinkError *error;
} UplinkUploadInfoResult;

typedef struct UplinkCommitUploadOptions {
    UplinkCustomMetadata custom_metadata;
} UplinkCommitUploadOptions;

typedef struct UplinkCommitUploadResult {
    UplinkObject *object;
    UplinkError *error;
} UplinkCommitUploadResult;

typedef struct UplinkPartResult {
    UplinkPart *part;
    UplinkError *error;
} UplinkPartResult;

typedef struct UplinkListUploadPartsOptions {
    uint32_t cursor;
} UplinkListUploadPartsOptions;

// Parameters when connecting to edge services
typedef struct EdgeConfig {
    // DRPC server e.g. auth.[eu|ap|us]1.storjshare.io:7777
    // Mandatory for now because this is no agreement on how to derive this
    const char *auth_service_address;

    // Root certificate(s) or chain(s) against which Uplink checks
    // the auth service.
    // In PEM format.
    // Intended to test against a self-hosted auth service
    // or to improve security.
    const char *certificate_pem;

    // Controls whether a client uses unencrypted connection.
    bool insecure_unencrypted_connection;
} EdgeConfig;

typedef struct EdgeRegisterAccessOptions {
    // Wether objects can be read using only the access_key_id.
    bool is_public;
} EdgeRegisterAccessOptions;

// Gateway credentials in S3 format
typedef struct EdgeCredentials {
    // Is also used in the linkshare url path
    const char *access_key_id;
    const char *secret_key;
    // Base HTTP(S) URL to the gateway.
    // The gateway and linkshare service are different endpoints.
    const char *endpoint;
} EdgeCredentials;

typedef struct EdgeCredentialsResult {
    EdgeCredentials *credentials;
    UplinkError *error;
} EdgeCredentialsResult;

typedef struct EdgeShareURLOptions {
    // Serve the file directly rather than through a landing page.
    bool raw;
} EdgeShareURLOptions;

// we need to suppress 'pedantic' validation because struct is empty for now
#pragma GCC diagnostic push
#pragma GCC diagnostic ignored "-Wpedantic"
typedef struct UplinkMoveObjectOptions {
} UplinkMoveObjectOptions;
#pragma GCC diagnostic pop

#pragma GCC diagnostic push
#pragma GCC diagnostic ignored "-Wpedantic"
typedef struct UplinkUploadObjectMetadataOptions {
} UplinkUploadObjectMetadataOptions;
#pragma GCC diagnostic pop

#pragma GCC diagnostic push
#pragma GCC diagnostic ignored "-Wpedantic"
typedef struct UplinkCopyObjectOptions {
} UplinkCopyObjectOptions;
#pragma GCC diagnostic pop
