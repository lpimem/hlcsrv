
DIR=$(dirname `readlink -f "$0"`)
SERVICE_DIR=/etc/systemd/system

if [[ ! -d $SERVICE_DIR ]] ; then
  echo "ERROR: service dir [$SERVICE_DIR] does not exist! "
  exit 1
fi

sudo echo "[Unit]
Description=HLC Server
After=nginx.service

[Service]
Type=forking
ExecStart=$DIR/start.sh
ExecStop=$DIR/stop.sh
PIDFile=$DIR/_run.pid
Restart=always" > $SERVICE_DIR/hlcsrv.service

sudo systemctl daemon-reload
