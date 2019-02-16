DIR=$(dirname `readlink -f "$0"`)

cd $DIR

if [[ $# -lt 1 ]]; then
  echo "Missing hostname"
  exit 1
fi

HOST=$1
DEPLOY_PATH='$HOME/hlcsrv'

ssh $HOST "cd $DEPLOY_PATH && $DEPLOY_PATH/stop.sh"
rsync -a --progress build/* $HOST:/$DEPLOY_PATH/ --exclude "test" --exclude "*.sqlite"

if [[ -f prod.env ]] ; then
    rsync -a --progress prod.env $HOST:/$DEPLOY_PATH/
    ssh $HOST "ln -s $DEPLOY_PATH/prod.env $DEPLOY_PATH/run.env"
fi

ssh $HOST "$DEPLOY_PATH/start.sh"
