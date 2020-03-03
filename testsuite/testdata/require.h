#pragma once

#include <assert.h>
#include <stdio.h>

#define require(test)                                                                                                  \
    do {                                                                                                               \
        if (!(test)) {                                                                                                 \
            printf("failed:\n\t%s:%d: %s\n", __FILE__, __LINE__, #test);                                               \
            exit(1);                                                                                                   \
        }                                                                                                              \
    } while (0)

#define requiref(test, msg, ...)                                                                                       \
    do {                                                                                                               \
        if (!(test)) {                                                                                                 \
            printf(msg, ##__VA_ARGS__);                                                                                \
            printf("failed:\n\t%s:%d: %s\n", __FILE__, __LINE__, #test);                                               \
            exit(1);                                                                                                   \
        }                                                                                                              \
    } while (0)

#define require_noerror(err)                                                                                           \
    do {                                                                                                               \
        if (err != NULL) {                                                                                             \
            printf("failed:\n\t%s:%d: [%d] %s\n", __FILE__, __LINE__, err->code, err->message);                        \
            exit(1);                                                                                                   \
        }                                                                                                              \
    } while (0)

#define require_error(err, expected)                                                                                   \
    do {                                                                                                               \
        if (err == NULL) {                                                                                             \
            printf("failed:\n\t%s:%d: NULL\n", __FILE__, __LINE__);                                                    \
            exit(1);                                                                                                   \
        } else if (err->code != expected) {                                                                            \
            printf("failed:\n\t%s:%d: [%d != %d] %s\n", __FILE__, __LINE__, expected, err->code, err->message);        \
            exit(1);                                                                                                   \
        }                                                                                                              \
    } while (0)
