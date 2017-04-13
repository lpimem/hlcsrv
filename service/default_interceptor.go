package service

import "github.com/lpimem/hlcsrv/session"

func buildDefaultInterceptors() {
	buildDefaultCookieChecker()
	AddRequestInterceptor(session.Authenticate)
}

func buildDefaultCookieChecker() {
	//var builder ReqCookieCheckerBuilder
	//builder.Require("uid")
	//AddRequestInterceptor(builder.Build())
}
