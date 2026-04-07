#pragma once

#include <assert.h>
#include <execinfo.h>
#include <stdio.h>
#include <stdlib.h>

#ifdef __clang__
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wgnu-zero-variadic-macro-arguments"
#endif

static inline void _require_print_backtrace(void)
{
    void *bt[16];
    int n = backtrace(bt, 16);
    // skip frame 0 (_require_print_backtrace itself)
    if (n > 1) {
        fflush(stdout);
        backtrace_symbols_fd(bt + 1, n - 1, 1);
    }
}

#define require(test)                                                                                                  \
    do {                                                                                                               \
        if (!(test)) {                                                                                                 \
            printf("failed:\n\t%s:%d: %s\n", __FILE__, __LINE__, #test);                                               \
            _require_print_backtrace();                                                                                \
            exit(1);                                                                                                   \
        }                                                                                                              \
    } while (0)

#define requiref(test, msg, ...)                                                                                       \
    do {                                                                                                               \
        if (!(test)) {                                                                                                 \
            printf(msg, ##__VA_ARGS__);                                                                                \
            printf("failed:\n\t%s:%d: %s\n", __FILE__, __LINE__, #test);                                               \
            _require_print_backtrace();                                                                                \
            exit(1);                                                                                                   \
        }                                                                                                              \
    } while (0)

#define require_noerror(err)                                                                                           \
    do {                                                                                                               \
        if (err != NULL) {                                                                                             \
            printf("failed:\n\t%s:%d: [%d] %s\n", __FILE__, __LINE__, err->code, err->message);                        \
            _require_print_backtrace();                                                                                \
            exit(1);                                                                                                   \
        }                                                                                                              \
    } while (0)

#define require_error(err, expected)                                                                                   \
    do {                                                                                                               \
        if (err == NULL) {                                                                                             \
            printf("failed:\n\t%s:%d: NULL\n", __FILE__, __LINE__);                                                    \
            _require_print_backtrace();                                                                                \
            exit(1);                                                                                                   \
        } else if (err->code != expected) {                                                                            \
            printf("failed:\n\t%s:%d: [%d != %d] %s\n", __FILE__, __LINE__, expected, err->code, err->message);        \
            _require_print_backtrace();                                                                                \
            exit(1);                                                                                                   \
        }                                                                                                              \
    } while (0)

#ifdef __clang__
#pragma clang diagnostic pop
#endif
