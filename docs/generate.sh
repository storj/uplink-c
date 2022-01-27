#!/bin/bash
##
## This script clones the gh-pages branch and generates documentation for the currently checked out version.
## Commit and push the documentation branch manually after changes.
## It will be hosted at https://storj.github.io/uplink-c/
##

set -e

## go to script dir
cd $(dirname "$0")

PROJECT_ROOT=$(realpath ..)

## checkout gh-pages branch here
if [ ! -d gh-pages ] ; then
    git clone git@github.com:storj/uplink-c.git gh-pages --branch gh-pages
else
    cd gh-pages
    git pull
fi

cd $PROJECT_ROOT

## make sure the header files are not tampered with
make build

cd .build/uplink

## remove copyright header because it will attach it to the next declaration
sed --in-place --regexp-extended 's+// Copyright \(C\) 2020 Storj Labs, Inc\.++g' uplink_definitions.h
sed --in-place 's+// See LICENSE for copying information\.++g' uplink_definitions.h

## remove cgo barf
## https://stackoverflow.com/questions/6287755/using-sed-to-delete-all-lines-between-two-matching-patterns/6287940
sed --in-place -n -e "1,/Start of boilerplate cgo prologue/ p" -e "/\End of boilerplate cgo prologue/,$ p" uplink.h
sed --in-place 's/#define GO_CGO_EXPORT_PROLOGUE_H//g' uplink.h
sed --in-place -e 's/typedef struct { const char \*p; ptrdiff_t n; } _GoString_;//g' uplink.h

for HEADER_FILE in *.h; do
    ## change comment style to 2 lines of "///" so that it is picked up by doxygen
    sed --in-place 's+//+///+g' $HEADER_FILE
done

cd $PROJECT_ROOT

## rewrite local links to github code browser
cp README.md .build/README.md
sed --in-place -r 's~\[(.+)\]\(([^h].+)\)~[\1](https://github.com/storj/uplink-c/tree/main/\2)~g' .build/README.md

## generate docs
doxygen docs/Doxyfile
