package service

import "net/http"

type ReqCookieCheckerBuilder []string

func (builder ReqCookieCheckerBuilder) Require(key string) {
	builder = append(builder, key)
}

func (builder ReqCookieCheckerBuilder) Build() RequestInterceptor {
	return func(req *http.Request) error {
		for _, expect := range builder {
			if _, err := req.Cookie(expect); err != nil {
				return err
			}
		}
		return nil
	}
}
