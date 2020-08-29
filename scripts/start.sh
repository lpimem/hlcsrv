#! /bin/bash

DIR=$(dirname `readlink -f "$0"`)
export HLC_ROOT=$DIR
cd $HLC_ROOT

ENV_PROFILE=$HLC_ROOT/run.env

if [[ ! -f $ENV_PROFILE ]] ; then
  echo "Please create $ENV_PROFILE and add required ENV variables"
  echo "HLC_SESSION_SECRET"
  echo "HLC_SESSION_KEY_USER"
  echo "HLC_SESSION_KEY_SID"
  exit 1
fi

if [[ -f "$HLC_ROOT/_run.pid" ]] ; then
    last_run=`cat $HLC_ROOT/_run.pid`
    kill $last_run
    rm $DIR/_run.pid
fi

source $ENV_PROFILE

$HLC_ROOT/hlcsrv