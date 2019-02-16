DIR=$(dirname `readlink -f "$0"`)

HOST=io3g
PATH=$HOME/hlcsrv

ssh $HOST "cd $PATH && $PATH/stop.sh"
rsync -a --progress build/* $HOST/$PATH/ --exclude "test"
ssh $HOST "$PATH/start.sh"
