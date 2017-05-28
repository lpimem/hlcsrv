package auth

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/lpimem/hlcsrv/conf"
	"github.com/lpimem/hlcsrv/hlccookie"
	"github.com/lpimem/hlcsrv/security"
	"github.com/lpimem/hlcsrv/storage"
)

func extractUIDSid(req *http.Request) (uid storage.UserID, sid string, err error) {
	uid, sid, err = extractUIDSidFromCookies(req)
	if err != nil {
		var err2 error
		uid, sid, err2 = extractUIDSidFromRequestHeader(req)
		if err2 != nil {
			err = errors.New(err.Error() + " & " + err2.Error())
			return
		}
		err = nil
	}
	return
}

func extractUIDSidFromCookies(req *http.Request) (uid storage.UserID, sid string, err error) {
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

func extractUIDSidFromRequestHeader(req *http.Request) (uid storage.UserID, sid string, err error) {
	var uid64 uint64
	uid64, err = strconv.ParseUint(req.Header.Get(HUSER_ID), 10, 32)
	if err != nil {
		return
	}
	uid = storage.UserID(uint32(uid64))
	sid = req.Header.Get(HSESSION_ID)
	if sid == "" {
		err = errors.New("missing session id header")
	}
	return
}

func computeRandomSessionID(seed string) string {
	once := security.RandStringBytesMaskImprSrc(32)
	return security.HashWithSlt(once, seed)
}
