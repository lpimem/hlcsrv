package service

import (
	"net/http"

	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/controller"
)

// returns an http.ServeMux instance for handling application http requests
func MakeRoutes() *http.ServeMux {
	mux := routes()
	return wrapProcessors(mux)
}

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/pagenote/delete", controller.DeletePagenote)
	mux.HandleFunc("/pagenote/new", controller.SavePagenote)
	mux.HandleFunc("/pagenote", controller.GetPagenote)
	mux.HandleFunc("/google_auth", controller.AuthenticateGoogleUser)
	fs := http.FileServer(
		http.Dir("static"))
	// http.Dir(util.GetAbsRunDirPath() + "static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", controller.Index)
	return mux
}

func wrapProcessors(mux *http.ServeMux) *http.ServeMux {
	wrapper := http.NewServeMux()
	wrapper.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			log.Info(r.Method, "\t", r.URL.String())
			//ctx := r.Context()
			//ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
			//defer cancel()
			//r = r.WithContext(ctx)
			var correct bool
			if r, correct = PreprocessRequest(w, r); !correct {
				return
			}
			mux.ServeHTTP(w, r)
		})
	return wrapper
}

func init() {
	buildDefaultInterceptors()
	log.Debug(len(interceptors), "interceptors loaded.")
}
