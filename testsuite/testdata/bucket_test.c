// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#include <string.h>
#include <stdlib.h>

#include "require.h"
#include "uplink.h"
#include "helpers.h"

void handle_project(Project project);

int main(int argc, char *argv[]) {
    with_test_project(&handle_project);
}

void handle_project(Project project) {
    char *_err = "";
    char **err = &_err;

    {
        // creating a new bucket
        Bucket bucket = create_bucket(project, "alpha", err);
        require_noerror(*err);

        require(strcmp("alpha", bucket.name) == 0);
        require(bucket.created != 0);

        free_bucket(bucket);
    }

    {
        // creating an existing bucket
        Bucket bucket = create_bucket(project, "alpha", err);
        require_error(*err);
        free(err);
        // TODO: verify exact error

        require(strcmp("alpha", bucket.name) == 0);
        require(bucket.created != 0);

        free_bucket(bucket);
    }

    {
        // ensuring an existing bucket
        Bucket bucket = ensure_bucket(project, "alpha", err);
        require_noerror(*err);

        require(strcmp("alpha", bucket.name) == 0);
        require(bucket.created != 0);

        free_bucket(bucket);
    }

    {
        // ensuring a new bucket
        Bucket bucket = ensure_bucket(project, "beta", err);
        require_noerror(*err);

        require(strcmp("beta", bucket.name) == 0);
        require(bucket.created != 0);

        free_bucket(bucket);
    }

    {
        // statting a bucket
        Bucket bucket = stat_bucket(project, "alpha", err);
        require_noerror(*err);

        require(strcmp("alpha", bucket.name) == 0);
        require(bucket.created != 0);

        free_bucket(bucket);
    }

    {
        // statting a missing bucket
        Bucket bucket = stat_bucket(project, "missing", err);
        require_error(*err);
        free(err);

        require(strcmp("", bucket.name) == 0);
        require(bucket.created == 0);

        free_bucket(bucket);
    }

    {
        // deleting a bucket
        delete_bucket(project, "alpha", err);
        require_noerror(*err);
    }

    {
        // deleting a missing bucket
        delete_bucket(project, "missing", err);
        require_error(*err);
        free(err);
    }

}