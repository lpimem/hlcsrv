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
White='\033[0;37m'
GREY='\033[1;30m'

TEST_RESULT=$TESTDIR/test_result

./test_setup.sh 1

for d in */ ; do
  tc=${d%/}
  pushd $tc > /dev/null
  go test -v --ldflags -s -o $TESTDIR/$tc.test ../$tc > $TEST_RESULT 2>&1
  ret=$?
  suc=0
  popd > /dev/null

  input="$TEST_RESULT"
  while IFS= read -r msg
  do
    if [[ $msg != *"can't load package"* ]]; then
        color=$GREY
        if [[ $msg == *"---"* ]]; then
          if [[ $msg != *"--- PASS"* ]]; then
            color=$RED
          else 
            color="${RESET}$GREY$GREEN"
          fi
        fi
        if [[ $msg == "?"* ]]; then
          color=$YELLOW
        elif [[ $msg == "ok"* ]]; then
          color="${RESET}$GREEN"
        elif [[ $msg == *"FAIL"* || $msg == *"should"* ]]; then
          color=$RED
          suc=1
        fi
      echo -e "${color}$msg"
    else
      break
    fi
  done < "$input"
done 

./test_setup.sh 2

rm $TEST_RESULT
exit $suc