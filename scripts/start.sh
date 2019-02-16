#! /bin/bash

DIR=$(dirname `readlink -f "$0"`)
export HLC_ROOT=$DIR

ENV_PROFILE=$HLC_ROOT/run.env

if [[ -f "$HLC_ROOT/_run.pid" ]] ; then
    last_run=`cat $HLC_ROOT/_run.pid`
    kill $last_run
    rm $DIR/_run.pid
fi

if [[ ! -f $ENV_PROFILE ]] ; then
  echo "Please create $ENV_PROFILE and add required ENV variables"
  echo "HLC_SESSION_SECRET"
  echo "HLC_SESSION_KEY_USER"
  echo "HLC_SESSION_KEY_SID"
  exit 1
fi

source $ENV_PROFILE

cd $HLC_ROOT
nohup $HLC_ROOT/hlcsrv >> $HLC_ROOT/log 2>&1 &
pid=$!
echo $pid > $HLC_ROOT/_run.pid
echo "Started with PID: $pid"
echo "Check logs: $HLC_ROOT/log"
