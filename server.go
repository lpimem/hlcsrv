package main

import (
	"fmt"
	"net/http"

	"github.com/lpimem/hlcsrv/service"
	"github.com/lpimem/hlcsrv/storage"
	"github.com/lpimem/hlcsrv/util"
)

func main() {
	_init()
	startServ()
}

func startServ() {
	serve("127.0.0.1", 23333)
}

func _init() {
	storage.InitStorage(util.GetAbsRunDirPath() + "/db/dev.db")
}

func serve(ip string, port int64) {
	mux := service.MakeRoutes()
	var server = http.Server{
		Addr:    fmt.Sprintf("%s:%d", ip, port),
		Handler: mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		util.Log("Error: ", err)
	}
}
