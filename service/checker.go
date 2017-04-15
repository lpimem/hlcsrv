package service

import "net/http"

// ReqCookieCheckerBuilder helps building a interceptor function to check
// the presents of cookies.
type ReqCookieCheckerBuilder struct {
	headers []string
}

// Tell ReqCookieCheckerBuilder a cookie with name key is required.
func (builder *ReqCookieCheckerBuilder) Require(key string) {
	builder.headers = append(builder.headers, key)
}

// Returns a interceptor function to check required cookies
func (builder ReqCookieCheckerBuilder) Build() RequestInterceptor {
	return func(req *http.Request) (*http.Request, error) {
		for _, expect := range builder.headers {
			if _, err := req.Cookie(expect); err != nil {
				return req, err
			}
		}
		return req, nil
	}
}
