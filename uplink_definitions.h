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

typedef struct UplinkListBucketsOptions {
    const char *cursor;
} UplinkListBucketsOptions;

typedef struct UplinkObjectIterator {
    size_t _handle;
} UplinkObjectIterator;

typedef struct UplinkBucketIterator {
    size_t _handle;
} UplinkBucketIterator;

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

#define UPLINK_ERROR_BUCKET_NAME_INVALID 0x10
#define UPLINK_ERROR_BUCKET_ALREADY_EXISTS 0x11
#define UPLINK_ERROR_BUCKET_NOT_EMPTY 0x12
#define UPLINK_ERROR_BUCKET_NOT_FOUND 0x13

#define UPLINK_ERROR_OBJECT_KEY_INVALID 0x20
#define UPLINK_ERROR_OBJECT_NOT_FOUND 0x21
#define UPLINK_ERROR_UPLOAD_DONE 0x22

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
