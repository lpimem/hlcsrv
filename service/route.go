package service

import (
	"net/http"

	"github.com/lpimem/hlcsrv/util"
)

func MakeRoutes() *http.ServeMux {
	mux := routes()
	return wrapProcessors(mux)
}

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/pagenote/delete", deletePagenote)
	mux.HandleFunc("/pagenote/new", savePagenote)
	mux.HandleFunc("/pagenote", getPagenote)
	fs := http.FileServer(
		http.Dir("static"))
	// http.Dir(util.GetAbsRunDirPath() + "static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", index)
	return mux
}

func wrapProcessors(mux *http.ServeMux) *http.ServeMux {
	wrapper := http.NewServeMux()
	wrapper.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			util.Debug(r.Method, "\t", r.URL.String())
			if !PreprocessRequest(w, r) {
				return
			}
			mux.ServeHTTP(w, r)
		})
	return wrapper
}

func init() {
	buildDefaultInterceptors()
}
