#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

$DIR/test_setup.sh 1

cd $DIR/controller

go test -bench=. -benchtime 10s -benchmem -blockprofile block.out -mutexprofile mutex.out --outputdir $DIR/build --ldflags -s -o ../build/controller.test ../controller


$DIR/test_setup.sh 0