package service

import "github.com/lpimem/hlcsrv/auth"
import "github.com/lpimem/hlcsrv/conf"

func buildDefaultInterceptors() {
	buildDefaultCookieChecker()
	AddRequestInterceptor(auth.Authenticate)
	if conf.IsDebug() {
		AddRequestInterceptor(logAccess)
	}
}

func buildDefaultCookieChecker() {
	//var builder ReqCookieCheckerBuilder
	//builder.Require("uid")
	//AddRequestInterceptor(builder.Build())
}
