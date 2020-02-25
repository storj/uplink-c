#pragma once

// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#include <stdint.h>
#include <stdbool.h>
#include <stdlib.h>

typedef struct Handle { long _handle; } Handle;

typedef struct Access   { long _handle; } Access;
typedef struct Project  { long _handle; } Project;
typedef struct Download { long _handle; } Download;
typedef struct Upload   { long _handle; } Upload;

typedef struct Bytes {
    void *data;
    uint64_t length;
} Bytes;

typedef struct Config {
    char *user_agent;
    bool skip_whitelist;
    int32_t dial_timeout_milliseconds;
} Config;

typedef struct Bucket {
    char *name;
    int64_t created;
} Bucket;

typedef struct SystemMetadata {
    int64_t created;
    int64_t expires;
} SystemMetadata;

typedef struct CustomMetadata {
    bool todo; //TODO: remove, here to avoid issues with empty struct
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

typedef struct ObjectIterator { long _handle; } ObjectIterator;
typedef struct BucketIterator { long _handle; } BucketIterator;


typedef struct Error {
    uint32_t code;
    char *message;
} Error;

#define ERROR_EOF            1
#define ERROR_INTERNAL       2
#define ERROR_CANCELED       3
#define ERROR_INVALID_HANDLE 4
#define ERROR_ALREADY_EXISTS 5
#define ERROR_NOT_FOUND      6

typedef struct AccessResult {
    Access *access;
    Error  *error;
} AccessResult;

typedef struct ProjectResult {
    Project *project;
    Error   *error;
} ProjectResult;

typedef struct BucketResult {
    Bucket *bucket;
    Error  *error;
} BucketResult;

typedef struct ObjectResult {
    Object *object;
    Error  *error;
} ObjectResult;

typedef struct UploadResult {
    Upload *upload;
    Error  *error;
} UploadResult;

typedef struct DownloadResult {
    Download *download;
    Error *error;
} DownloadResult;

typedef struct WriteResult {
    size_t bytes_written;
    Error  *error;
} WriteResult;

typedef struct ReadResult {
    size_t bytes_read;
    Error  *error;
} ReadResult;

typedef struct StringResult {
    char   *string;
    Error  *error;
} StringResult;
