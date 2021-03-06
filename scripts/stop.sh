#! /bin/bash

DIR=$(dirname `readlink -f "$0"`)
if [[ -f "$DIR/_run.pid" ]] ; then
    pid=`cat $DIR/_run.pid`
    kill $pid
    rm $DIR/_run.pid
else
    echo "Not running."
fi
