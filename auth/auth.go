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

const admin string = "admin"
const adminUserID storage.UserID = 1

// Authenticate implements the interceptor interface.
// It adds a flag to the request context to indicate if the
// request is authenticated. If authenticated, it also checks
// if the request is trying to access an admin URI and returns
// error if the request is not from the admin user.
func Authenticate(req *http.Request) (*http.Request, error) {
	var (
		uid storage.UserID
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
		return req, authorize(ctx, req)
	}
	if err = VerifySession(sid, uid, nil); err != nil {
		log.Info("invalid session", sid, uid, err)
		ctx = context.WithValue(ctx, REASON, err.Error())
		req = req.WithContext(ctx)
		return req, authorize(ctx, req)
	}
	ctx = context.WithValue(ctx, USER_ID, uid)
	ctx = context.WithValue(ctx, SESSION_ID, sid)
	ctx = context.WithValue(ctx, AUTHENTICATED, true)
	req = req.WithContext(ctx)
	log.Info("request from", uid, "is authorized.")
	return req, authorize(ctx, req)
}

// IsSessionTimeout returns if duration since lastAccess exceeds the max session lifetime
// Max session lifetime is defined by func conf.SessionValidHours()
func IsSessionTimeout(lastAccess time.Time) bool {
	return time.Since(lastAccess).Hours() >= conf.SessionValidHours()
}

// VerifySession verifies a session with claimed session id sid for user id uid, previously accessed
// at lastAccess is still valid.
func VerifySession(sid string, uid storage.UserID, lastAccess *time.Time) error {
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

func authorizeUser(ctx context.Context, r *http.Request) error {
	uid := ctx.Value(USER_ID)
	if uid == adminUserID {
		return nil
	}
	uri := r.URL.Path
	unauthorized := errors.New("unauthorized")
	restricted, err := storage.Restriction.Has(uri)
	if err != nil {
		log.Errorf("authorizeUser: cannot query restriction: %s", err)
		return errors.New("bad gateway")
	}
	if restricted {
		if uid == nil {
			return unauthorized
		}
		if acc, err := storage.Permission.HasAccess(uid.(storage.UserID), uri); !acc {
			if err != nil {
				log.Errorf("error checking permission for user %d to %s : %s", uid, uri, err)
			}
			return unauthorized
		}
	} else {
		log.Debugf("URI %s is not restricted", uri)
	}
	return nil
}

func authorizeAdmin(ctx context.Context, r *http.Request) error {
	var err error
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
				log.Warn("User [", uid, "] is unauthorized to access admin URI ", r.URL.Path)
				err = errors.New("unauthorized")
			}
		}
	}
	return err
}

func authorize(ctx context.Context, r *http.Request) error {
	if err := authorizeAdmin(ctx, r); err != nil {
		return err
	}
	if err := authorizeUser(ctx, r); err != nil {
		return err
	}
	return nil
}
