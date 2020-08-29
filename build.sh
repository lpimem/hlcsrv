DIR=$(dirname `readlink -f $0`)
cd $DIR

go build --ldflags -s -o build/hlcsrv
ret=$!
rsync -a db/*.sql build/db/
rsync -a view build/
rsync -a static build/
rsync -a scripts/*.sh build/
exit $ret
