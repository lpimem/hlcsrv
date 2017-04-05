package service

import (
	"fmt"
	"net/http"

	"github.com/lpimem/hlcsrv/util"
)

func Serve(ip string, port int64) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		util.Log("Request: ", req.URL.Path)
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "welcome.")
	})
	var server = http.Server{
		Addr:    fmt.Sprintf("%s:%d", ip, port),
		Handler: mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Errorf("Error starting server", err)
	}
}
