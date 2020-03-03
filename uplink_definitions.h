#pragma once

// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

typedef struct Handle {
    long _handle;
} Handle;

typedef struct Access {
    long _handle;
} Access;
typedef struct Project {
    long _handle;
} Project;
typedef struct Download {
    long _handle;
} Download;
typedef struct Upload {
    long _handle;
} Upload;

typedef struct Config {
    char *user_agent;

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

// TODO: is a structure more convenient:
//   Bytes key;
//   Bytes value;

typedef struct CustomMetadataEntry {
    char *key;           // TODO: should this be void *?
    uint64_t key_length; // TODO: should this be size_t?

    char *value;           // TODO: should this be void *?
    uint64_t value_length; // TODO: should this be size_t?
} CustomMetadataEntry;

typedef struct CustomMetadata {
    CustomMetadataEntry *entries;
    uint64_t count;
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
    long _handle;
} ObjectIterator;
typedef struct BucketIterator {
    long _handle;
} BucketIterator;

typedef struct Permission {
    bool allow_read;
    bool allow_write;
    bool allow_list;
    bool allow_delete;

    // TODO: not before and not after
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
