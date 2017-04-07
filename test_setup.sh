#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

if [[ $1 == 1 ]]; then
    ## setup: rename file so that compile will include this file.
    mv $DIR/storage/storage_utils_test.go $DIR/storage/storage_test_utils.go
    ##
elif [[ $1 == 2 ]]; then
    ## reset filename so that compiler will ignore it for normal builds.
    mv $DIR/storage/storage_test_utils.go $DIR/storage/storage_utils_test.go
    ##
fi