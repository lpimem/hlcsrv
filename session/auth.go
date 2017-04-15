package session

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/conf"
	"github.com/lpimem/hlcsrv/hlccookie"
	"github.com/lpimem/hlcsrv/storage"
)

type context_key int

const (
	AUTHENTICATED context_key = 1
	USER_ID       context_key = 2
	SESSION_ID    context_key = 3
	REASON        context_key = 4
)

/** Authenticate always returns nil.
It add a flag to the request context to indicate if the request
is authenticated.
*/
func Authenticate(req *http.Request) (*http.Request, error) {
	var (
		uid uint32
		sid string
		err error
		ctx context.Context
	)
	ctx = req.Context()
	ctx = context.WithValue(ctx, AUTHENTICATED, false)
	req = req.WithContext(ctx)
	uid, sid, err = extractUidSid(req)
	if err != nil {
		log.Info("cannot extract uid/sid:", sid, uid, err)
		ctx = context.WithValue(ctx, REASON, err.Error())
		req = req.WithContext(ctx)
		return req, nil
	}
	if err = VerifySession(sid, uid, nil); err != nil {
		log.Info("invalid session", sid, uid, err)
		ctx = context.WithValue(ctx, REASON, err.Error())
		req = req.WithContext(ctx)
		return req, nil
	}
	ctx = context.WithValue(ctx, AUTHENTICATED, true)
	ctx = context.WithValue(ctx, USER_ID, uid)
	ctx = context.WithValue(ctx, SESSION_ID, sid)
	req = req.WithContext(ctx)
	log.Info("request from", uid, "is authorized.")
	return req, nil
}

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

func IsSessionTimeout(lastAccess time.Time) bool {
	return time.Since(lastAccess).Hours() >= conf.SessionValidHours()
}

func VerifySession(sid string, uid uint32, lastAccess *time.Time) error {
	var (
		err error
	)
	if lastAccess == nil {
		lastAccess, err = storage.QuerySession(sid, uid)
		if err != nil {
			log.Warn("cannot find session for ", sid, uid, err)
			return err
		}
		if lastAccess == nil {
			return errors.New("no session found.")
		}
	}
	if IsSessionTimeout(*lastAccess) {
		err = errors.New("session time out for" + sid)
		return err
	}
	return nil
}

func IsAuthenticated(r *http.Request) bool {
	ctx := r.Context()
	v := ctx.Value(AUTHENTICATED)
	if v == nil {
		return false
	}
	return v.(bool)
}
