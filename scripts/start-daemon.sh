#! /bin/bash

DIR=$(dirname `readlink -f "$0"`)
cd ${DIR}
export HLC_ROOT=$DIR
nohup ./start.sh >> $HLC_ROOT/log 2>&1 &
pid=$!
echo $pid > $HLC_ROOT/_run.pid
echo "Started with PID: $pid"
echo "Check logs: $HLC_ROOT/log"