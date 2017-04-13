package controller

import (
	"context"

	"errors"

	"time"

	"log"

	"github.com/coreos/go-oidc"
	"github.com/lpimem/hlcsrv/google-auth"
	"github.com/lpimem/hlcsrv/security"
	"github.com/lpimem/hlcsrv/session"
	"github.com/lpimem/hlcsrv/storage"
)

type SessionInfo struct {
	Uid uint32
	Sid string
}

type GoogleTokenClaim struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func computeRandomSessionId(seed string) string {
	once := security.RandStringBytesMaskImprSrc(32)
	return security.HashWithSlt(once, seed)
}

func doAuthenticateGoogleUser(ctx context.Context, rawToken string) (*SessionInfo, error) {
	var (
		idToken *oidc.IDToken
		err     error
	)
	idToken, err = google_auth.VerifyGoogleAuthIdToken(ctx, rawToken)
	if err != nil {
		return nil, err
	}
	var profile = GoogleTokenClaim{}
	idToken.Claims(&profile)
	if !profile.EmailVerified {
		return nil, errors.New("email not verified")
	}
	return updateGoogleUserSession(profile.Email, profile.Email)
}

func updateGoogleUserSession(
	gid, email string,
) (*SessionInfo, error) {
	var (
		uid        uint32
		sid        string
		lastAccess *time.Time
		err        error
	)
	uid, err = storage.GetOrCreateUidForGoogleUser(gid, email)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	sInfo, err := storage.QuerySessionByUid(uid)
	if err != nil {
		return nil, err
	}
	if sInfo == nil || sInfo.LastAccess == nil || sInfo.Sid == "" {
		// no existing session, create new.
		sid = computeRandomSessionId(gid)
		err = storage.UpdateSession(sid, uid)
		if err != nil {
			return nil, err
		}
	} else {
		// verify if existing session timed out
		sid = sInfo.Sid
		lastAccess = sInfo.LastAccess
		if err = session.VerifySession(sid, uid, lastAccess); err != nil {
			return nil, err
		}
		// refresh session.
		go storage.UpdateSession(sid, uid)
	}
	return &SessionInfo{uid, sid}, nil
}
