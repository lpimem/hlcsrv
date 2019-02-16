DIR=$(dirname `readlink -f "$0"`)

if [[ -d $DIR/build/ ]] ; then
  rm -rf $DIR/build/*
fi
