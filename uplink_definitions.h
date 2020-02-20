#pragma once

#ifndef UPLINK_DEFINITIONS
#define UPLINK_DEFINITIONS

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

typedef struct ObjectInfo {
    int64_t created;
    int64_t expires;
} ObjectInfo;

typedef struct StandardMetadata {
    int64_t content_length;
    char *content_type;
    char *etag;

    int64_t file_created;
    int64_t file_modified;
    uint32_t file_permissions;

    Bytes unknown;
} StandardMetadata;

typedef struct CustomMetadata {
    bool todo; //TODO: remove, here to avoid issues with empty struct
} CustomMetadata;

typedef struct Object {
    char *key;
    ObjectInfo info;
    StandardMetadata standard;
    CustomMetadata custom;
} Object;

typedef struct Objects { long _handle; } Objects;
typedef struct ObjectsOptions { 
    char *prefix;
    char *cursor;
    bool recursive;

    bool info;
    bool standard;
    bool custom;
} ObjectsOptions;

typedef struct Buckets { long _handle; } Buckets;
typedef struct BucketsOptions { 
    char *cursor;
} BucketsOptions;

typedef struct Error {
    uint32_t code;
    char *message;
} Error;

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

#endif