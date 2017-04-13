package service

import "net/http"

type ReqCookieCheckerBuilder struct {
	headers []string
}

func (builder *ReqCookieCheckerBuilder) Require(key string) {
	builder.headers = append(builder.headers, key)
}

func (builder ReqCookieCheckerBuilder) Build() RequestInterceptor {
	return func(req *http.Request) error {
		for _, expect := range builder.headers {
			if _, err := req.Cookie(expect); err != nil {
				return err
			}
		}
		return nil
	}
}


