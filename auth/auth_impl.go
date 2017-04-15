package auth

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/lpimem/hlcsrv/conf"
	"github.com/lpimem/hlcsrv/hlccookie"
	"github.com/lpimem/hlcsrv/security"
)

func extractUidSid(req *http.Request) (uid uint32, sid string, err error) {
	uid, sid, err = extractUidSidFromCookies(req)
	if err != nil {
		var err2 error
		uid, sid, err2 = extractUidSidFromRequestHeader(req)
		if err2 != nil {
			err = errors.New(err.Error() + " & " + err2.Error())
			return
		} else {
			err = nil
		}
	}
	return
}

func extractUidSidFromCookies(req *http.Request) (uid uint32, sid string, err error) {
	var c *http.Cookie
	if uid, err = hlccookie.GetRequestUID(req); err != nil {
		return
	}
	if c, err = req.Cookie(conf.SessionKeySID()); err != nil {
		return
	}
	sid = c.Value
	return
}

func extractUidSidFromRequestHeader(req *http.Request) (uid uint32, sid string, err error) {
	var uid64 uint64
	uid64, err = strconv.ParseUint(req.Header.Get(HUSER_ID), 10, 32)
	if err != nil {
		return
	}
	uid = uint32(uid64)
	sid = req.Header.Get(HSESSION_ID)
	if sid == "" {
		err = errors.New("missing session id header")
	}
	return
}

func computeRandomSessionId(seed string) string {
	once := security.RandStringBytesMaskImprSrc(32)
	return security.HashWithSlt(once, seed)
}
