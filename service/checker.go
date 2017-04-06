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

func RequirePost(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodPost {
		return true
	} else {
		http.Error(w, "only post accepted", http.StatusBadRequest)
		return false
	}
}
