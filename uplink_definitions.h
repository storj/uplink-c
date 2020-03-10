#pragma once

// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

typedef struct Handle {
    size_t _handle;
} Handle;

typedef struct Access {
    size_t _handle;
} Access;
typedef struct Project {
    size_t _handle;
} Project;
typedef struct Download {
    size_t _handle;
} Download;
typedef struct Upload {
    size_t _handle;
} Upload;

typedef struct Config {
    char *user_agent;

    int32_t dial_timeout_milliseconds;

    // temp_directory specifies where to save data during downloads to use less memory.
    char *temp_directory;
} Config;

typedef struct Bucket {
    char *name;
    int64_t created;
} Bucket;

typedef struct SystemMetadata {
    int64_t created;
    int64_t expires;
    int64_t content_length;
} SystemMetadata;

typedef struct CustomMetadataEntry {
    char *key;
    size_t key_length;

    char *value;
    size_t value_length;
} CustomMetadataEntry;

typedef struct CustomMetadata {
    CustomMetadataEntry *entries;
    size_t count;
} CustomMetadata;

typedef struct Object {
    char *key;
    bool is_prefix;
    SystemMetadata system;
    CustomMetadata custom;
} Object;

typedef struct UploadOptions {
    // When expires is 0 or negative, it means no expiration.
    int64_t expires;
} UploadOptions;

typedef struct DownloadOptions {
    int64_t offset;
    // When length is negative, it will read until the end of the blob.
    int64_t length;
} DownloadOptions;

typedef struct ListObjectsOptions {
    char *prefix;
    char *cursor;
    bool recursive;

    bool system;
    bool custom;
} ListObjectsOptions;

typedef struct ListBucketsOptions {
    char *cursor;
} ListBucketsOptions;

typedef struct ObjectIterator {
    size_t _handle;
} ObjectIterator;
typedef struct BucketIterator {
    size_t _handle;
} BucketIterator;

typedef struct Permission {
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
} Permission;

typedef struct SharePrefix {
    char *bucket;
    // prefix is the prefix of the shared object keys.
    char *prefix;
} SharePrefix;

typedef struct Error {
    uint32_t code;
    char *message;
} Error;

#define ERROR_EOF 1
#define ERROR_INTERNAL 2
#define ERROR_CANCELED 3
#define ERROR_INVALID_HANDLE 4
#define ERROR_ALREADY_EXISTS 5
#define ERROR_NOT_FOUND 6

typedef struct AccessResult {
    Access *access;
    Error *error;
} AccessResult;

typedef struct ProjectResult {
    Project *project;
    Error *error;
} ProjectResult;

typedef struct BucketResult {
    Bucket *bucket;
    Error *error;
} BucketResult;

typedef struct ObjectResult {
    Object *object;
    Error *error;
} ObjectResult;

typedef struct UploadResult {
    Upload *upload;
    Error *error;
} UploadResult;

typedef struct DownloadResult {
    Download *download;
    Error *error;
} DownloadResult;

typedef struct WriteResult {
    size_t bytes_written;
    Error *error;
} WriteResult;

typedef struct ReadResult {
    size_t bytes_read;
    Error *error;
} ReadResult;

typedef struct StringResult {
    char *string;
    Error *error;
} StringResult;
