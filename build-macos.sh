DIR=$(dirname `realpath -P $0`)
cd $DIR
CGO_ENABLED=1 go build --ldflags -s "$@" -v -o build/hlcsrv
ret=$!
if [[ ${ret} -ne 0 ]]; then
    echo "If go-sqlite3 is complaining, try the following commands. 
    brew install FiloSottile/musl-cross/musl-cross  
    brew install mingw-w64
If you are cross compiling from macos to linux, try this:
    GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc ./build-macos.sh
Reference:
* https://github.com/mattn/go-sqlite3/issues/372#issuecomment-396863368
* https://github.com/mattn/go-sqlite3/issues/372#issuecomment-398001083"
    exit ${ret}
fi
rsync -a db/*.sql build/db/
rsync -a view build/
rsync -a static build/
rsync -a scripts/*.sh build/
exit $ret
