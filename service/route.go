package service

import (
	"net/http"

	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/controller"
)

// MakeRoutes returns an http.ServeMux instance for handling application http requests
func MakeRoutes() *http.ServeMux {
	mux := routes()
	return wrapProcessors(mux)
}

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	// hlc
	mux.HandleFunc("/pagenote/delete",
		controller.DeletePagenote)

	mux.HandleFunc("/pagenote/new",
		controller.SavePagenote)

	mux.HandleFunc("/pagenote",
		controller.GetPagenote)

	mux.HandleFunc("/google_auth",
		controller.AuthenticateGoogleUser)

	mux.HandleFunc("/q",
		controller.HandleQuery)

	// admin
	mux.HandleFunc("/admin/users",
		controller.Admin.Users)

	mux.HandleFunc("/admin/permissions",
		controller.Admin.Permissions)

	mux.HandleFunc("/admin/restrictions",
		controller.Admin.Restrictions)

	mux.HandleFunc("/admin/restrict",
		controller.Admin.AddRestriction)

	mux.HandleFunc("/admin/unrestrict",
		controller.Admin.RemoveRestriction)

	// static files
	fs := http.FileServer(
		http.Dir("static"))
	mux.Handle("/static/",
		http.StripPrefix("/static/", fs))

	// index
	mux.HandleFunc("/", controller.Index)

	return mux
}

func wrapProcessors(mux *http.ServeMux) *http.ServeMux {
	wrapper := http.NewServeMux()
	wrapper.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			log.Info(r.Method, "\t", r.URL.String())
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
