package service

import (
	"net/http"

	"github.com/go-playground/log"
)

/*RequestInterceptor preprocess a request before it is handled
 * by a listener.
 * It can prevent the request from being further processed by returning
 * an error
 */
type RequestInterceptor func(req *http.Request, respWriter http.ResponseWriter) (*http.Request, bool, error)

var interceptors = []RequestInterceptor{}

/*AddRequestInterceptor add a new request preprocessor (interceptor).
 * Interceptors are called before each request is handled by a listener.
 * This function should be called before start serving requests.
 */
func AddRequestInterceptor(handler RequestInterceptor) {
	interceptors = append(interceptors, handler)
}

/*PreprocessRequest applies each interceptor to the request before it reaches
 * specific request listeners
 * It returns true if no interceptor is complaining error.
 * When the return value is false, the respWriter will be modifed with the error status
 */
func PreprocessRequest(w http.ResponseWriter, req *http.Request) (*http.Request, bool) {
	var handled bool 
	var err error 
	for _, handle := range interceptors {
		defer log.WithTrace().Info("applying preprocessor:", handle)
		req, handled, err = handle(req, w)
		if err != nil {
			log.Warn("error pre-processing", err)
			if !handled {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return req, false
		}
	}
	defer log.WithTrace().Info("all pre-processors done.")
	return req, true
}
