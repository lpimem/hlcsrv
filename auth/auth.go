package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/conf"
	"github.com/lpimem/hlcsrv/storage"
)

/** Authenticate always returns nil. It implements the interceptor
interface.
Authenticate add a flag to the request context to indicate if the request
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

// return if duration since lastAccess exceeds the max session lifetime
// Max session lifetime is defined by func conf.SessionValidHours()
func IsSessionTimeout(lastAccess time.Time) bool {
	return time.Since(lastAccess).Hours() >= conf.SessionValidHours()
}

// verify a session with claimed session id sid for user id uid, previously accessed
// at lastAccess is still valid.
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

// Verify if a http request r is already validated
// this function should be called in the begining of each
// request handler which checks authentication.
func IsAuthenticated(r *http.Request) bool {
	ctx := r.Context()
	v := ctx.Value(AUTHENTICATED)
	if v == nil {
		return false
	}
	return v.(bool)
}
