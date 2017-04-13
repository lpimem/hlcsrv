package service

import "net/http"

/*RequestInterceptor preprocess a request before it is handled
 * by a listener.
 * It can prevent the request from being further processed by returning
 * an error
 */
type RequestInterceptor func(req *http.Request) (*http.Request, error)

var interceptors []RequestInterceptor

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
func PreprocessRequest(respWriter http.ResponseWriter, req *http.Request) bool {
	var err error
	for _, handle := range interceptors {
		req, err = handle(req)
		if err != nil {
			http.Error(respWriter, err.Error(), http.StatusBadRequest)
			return false
		}
	}
	return true
}
