package main

import (
	"github.com/lpimem/hlcsrv/service"
	"github.com/lpimem/hlcsrv/storage"
	"github.com/lpimem/hlcsrv/util"
)

func startServ() {
	util.Log("start ...")
	service.Serve("127.0.0.1", 23333)
	util.Log("end ...")
}

func main() {
	util.Log("init db ...")
	storage.InitStorage("db/test.db")

	util.Log("exit 0")
}
