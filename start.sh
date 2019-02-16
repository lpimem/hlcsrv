#! /bin/bash

DIR=$(dirname `readlink -f "$0"`)
export HLC_ROOT=$DIR

if [[ -f "$HLC_ROOT/_run.pid" ]] ; then
    last_run=`cat $HLC_ROOT/_run.pid`
    kill $last_run
fi

$HLC_ROOT/hlcsrv >> $HLC_ROOT/log 2>&1 &
pid=$!
echo $pid > $HLC_ROOT/_run.pid
echo "Started with PID: `cat $HLC_ROOT/_run.pid`"
tail -f $HLC_ROOT/log
