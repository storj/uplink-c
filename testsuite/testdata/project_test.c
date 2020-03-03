// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

#include <stdlib.h>
#include <string.h>

#include "helpers.h"
#include "require.h"
#include "uplink.h"

void handle_project(Project *project);

int main(int argc, char *argv[])
{
    with_test_project(&handle_project);
    return 0;
}

void handle_project(Project *project) { require(project->_handle != 0); };
