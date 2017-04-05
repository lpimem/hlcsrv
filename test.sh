#! /bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
TESTDIR=$DIR/build

echo "TEST DIR: $TESTDIR"

mkdir -p $TESTDIR
rm -rf $TESTDIR/db
cp -r $DIR/db $TESTDIR/

RED=`tput setaf 1`
GREEN=`tput setaf 2`
YELLOW=`tput setaf 3`
RESET=`tput sgr0`

TEST_RESULT=$TESTDIR/test_result

for d in */ ; do
  tc=${d%/}
  pushd $tc > /dev/null
  go test --ldflags -s -o $TESTDIR/$tc.test ../$tc > $TEST_RESULT 2>&1
  ret=$?
  popd > /dev/null
  msg=`cat $TEST_RESULT`
  if [[ $msg != *"can't load package"* ]]; then
    if [[ $ret == 1 ]]; then 
      color=$RED
    else
      if [[ $msg == "?"* ]]; then 
        color=$YELLOW
      else
        color=$GREEN
      fi 
    fi 
    echo -e "${color} $msg ${RESET}"
  fi 
done 

rm $TEST_RESULT