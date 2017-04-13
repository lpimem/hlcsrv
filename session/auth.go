package session

import (
	"context"
	"net/http"
	"time"

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
)

/** Authenticate always returns nil.
It add a flag to the request context to indicate if the request
is authenticated.
*/
func Authenticate(req *http.Request) error {
	var (
		uid        uint32
		c          *http.Cookie
		sid        string
		lastAccess *time.Time
		idleTime   time.Duration
		err        error
		ctx        context.Context
	)
	ctx = req.Context()
	ctx = context.WithValue(ctx, AUTHENTICATED, false)
	defer req.WithContext(ctx)
	if uid, err = hlccookie.GetRequestUID(req); err != nil {
		util.Log("cannot extract uid from cookie", err)
		return nil
	}
	if c, err = req.Cookie(conf.SessionKeySID()); err != nil {
		return nil
	}
	sid = c.Value
	lastAccess, err = storage.QuerySession(sid, uid)
	if err != nil {
		util.Log("cannot find session for ", sid, uid, err)
		return nil
	}
	idleTime = time.Now().Sub(*lastAccess)
	if idleTime.Hours() > conf.SessionValidHours() {
		util.Log("session time out for", sid, uid)
		return nil
	}
	ctx = context.WithValue(ctx, AUTHENTICATED, true)
	ctx = context.WithValue(ctx, USER_ID, uid)
	ctx = context.WithValue(ctx, SESSION_ID, sid)
	return nil
}

func IsAuthenticated(r *http.Request) bool {
	ctx := r.Context()
	return ctx.Value(AUTHENTICATED).(bool)
}
