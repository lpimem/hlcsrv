cp db/*.sql build/db/
go build --ldflags -s -o build/hlcsrv
ret=$!
rsync -a --progress view build/
rsync -a --progress static build/
exit $ret