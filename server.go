package main

import (
	"github.com/lpimem/hlcsrv/service"
	"github.com/lpimem/hlcsrv/storage"
	"github.com/lpimem/hlcsrv/util"
)

func main() {
	_init()
	startServ()
}

func startServ() {
	service.Serve("127.0.0.1", 23333)
}

func _init() {
	storage.InitStorage(util.GetAbsRunDirPath() + "/db/dev.db")
}
