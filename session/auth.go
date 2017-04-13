package session

import (
	"context"
	"net/http"
	"time"

	"errors"

	"github.com/lpimem/hlcsrv/conf"
	"github.com/lpimem/hlcsrv/hlccookie"
	"github.com/lpimem/hlcsrv/storage"
	"github.com/lpimem/hlcsrv/util"
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
		c   *http.Cookie
		sid string
		err error
		ctx context.Context
	)
	ctx = req.Context()
	ctx = context.WithValue(ctx, AUTHENTICATED, false)
	req = req.WithContext(ctx)
	if uid, err = hlccookie.GetRequestUID(req); err != nil {
		util.Log("cannot extract uid from cookie", err)
		ctx = context.WithValue(ctx, REASON, "COOKIE UID "+err.Error())
		req = req.WithContext(ctx)
		return req, nil
	}
	if c, err = req.Cookie(conf.SessionKeySID()); err != nil {
		util.Log("cannot extract sid from cookie", err)
		ctx = context.WithValue(ctx, REASON, "COOKIE SID "+err.Error())
		req = req.WithContext(ctx)
		return req, nil
	}
	sid = c.Value
	if err = VerifySession(sid, uid, nil); err != nil {
		util.Log("invalid session", sid, uid, err)
		ctx = context.WithValue(ctx, REASON, err.Error())
		req = req.WithContext(ctx)
		return req, nil
	}
	ctx = context.WithValue(ctx, AUTHENTICATED, true)
	ctx = context.WithValue(ctx, USER_ID, uid)
	ctx = context.WithValue(ctx, SESSION_ID, sid)
	req = req.WithContext(ctx)
	return req, nil
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
			util.Log("cannot find session for ", sid, uid, err)
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
