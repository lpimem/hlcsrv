package service

import (
	"fmt"
	"net/http"

	"github.com/lpimem/hlcsrv/util"
)

/*Serve start listening http requests to the given ip and port
 */
func Serve(ip string, port int64) {
	mux := routes()
	var server = http.Server{
		Addr:    fmt.Sprintf("%s:%d", ip, port),
		Handler: mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		util.Log("Error: ", err)
	}
}

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/pagenote/delete", deletePageNotes)
	mux.HandleFunc("/pagenote", getPageNotes)
	mux.HandleFunc("/", index)
	return mux
}

func wrap_processors(mux *http.ServeMux) *http.ServeMux {
	wrapper := http.NewServeMux()
	wrapper.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			if !PreprocessRequest(w, r) {
				return
			}
			mux.ServeHTTP(w, r)
		})
	return wrapper
}

func make_routes() *http.ServeMux {
	mux := routes()
	return wrap_processors(mux)
}

func init() {
	buildDefaultInterceptors()
}
