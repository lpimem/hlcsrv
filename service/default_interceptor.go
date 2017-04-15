package service

import "github.com/lpimem/hlcsrv/auth"

func buildDefaultInterceptors() {
	buildDefaultCookieChecker()
	AddRequestInterceptor(auth.Authenticate)
}

func buildDefaultCookieChecker() {
	//var builder ReqCookieCheckerBuilder
	//builder.Require("uid")
	//AddRequestInterceptor(builder.Build())
}
