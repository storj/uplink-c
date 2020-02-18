// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#include <string.h>
#include <stdlib.h>

#include "require.h"
#include "uplink.h"
#include "helpers.h"

void handle_project(Project project)
{};

int main(int argc, char *argv[]) {
    with_test_project(&handle_project);
}
