cp db/*.sql build/db/
go build --ldflags -s -o build/hlcsrv
