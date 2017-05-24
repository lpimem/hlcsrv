package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"strings"

	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/conf"
	"github.com/lpimem/hlcsrv/storage"
)

// Authenticate implements the interceptor interface.
// It adds a flag to the request context to indicate if the
// request is authenticated. If authenticated, it also checks
// if the request is trying to access an admin URI and returns
// error if the request is not from the admin user.
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
	uid, sid, err = extractUIDSid(req)
	if err != nil {
		log.Info("cannot extract uid/sid:", sid, uid, err)
		ctx = context.WithValue(ctx, REASON, err.Error())
		req = req.WithContext(ctx)
		return req, authorizeAdmin(ctx, req)
	}
	if err = VerifySession(sid, uid, nil); err != nil {
		log.Info("invalid session", sid, uid, err)
		ctx = context.WithValue(ctx, REASON, err.Error())
		req = req.WithContext(ctx)
		return req, authorizeAdmin(ctx, req)
	}
	ctx = context.WithValue(ctx, USER_ID, uid)
	ctx = context.WithValue(ctx, SESSION_ID, sid)
	ctx = context.WithValue(ctx, AUTHENTICATED, true)
	req = req.WithContext(ctx)
	log.Info("request from", uid, "is authorized.")
	return req, authorizeAdmin(ctx, req)
}

// IsSessionTimeout returns if duration since lastAccess exceeds the max session lifetime
// Max session lifetime is defined by func conf.SessionValidHours()
func IsSessionTimeout(lastAccess time.Time) bool {
	return time.Since(lastAccess).Hours() >= conf.SessionValidHours()
}

// VerifySession verifies a session with claimed session id sid for user id uid, previously accessed
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
			return errors.New("no session found")
		}
	}
	if IsSessionTimeout(*lastAccess) {
		err = errors.New("session time out for" + sid)
		return err
	}
	return nil
}

// IsAuthenticated : Verify if a http request r is already validated
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

func authorizeAdmin(ctx context.Context, r *http.Request) error {
	var err error
	const admin string = "admin"
	const adminUserID uint32 = 1
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) > 1 {
		base := parts[1]
		var arg string
		if len(parts) > 2 {
			arg = parts[2]
		}
		if base == admin || base == "static" && arg == admin {
			uid := ctx.Value(USER_ID)
			if uid != adminUserID {
				log.Warn("User [", uid, "] is unauthorized to access ", r.URL.Path)
				err = errors.New("unauthorized")
			}
		}
	}
	return err
}
