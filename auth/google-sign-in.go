package auth

import (
	"context"
	"errors"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/storage"
)

/**AuthenticateGoogleUser authenticates a google user with rawToken
The token will be parsed and validated.

See also:
1. https://developers.google.com/identity/sign-in/web/sign-in
2. https://developers.google.com/identity/sign-in/web/backend-auth#verify-the-integrity-of-the-id-token
3. https://github.com/coreos/go-oidc/blob/c3a2c79e8008bc1b1b0509ae6bf1483642c976f4/example/idtoken/app.go#L66
4. OAuth 2.0 Bearer Token Usage https://tools.ietf.org/html/rfc6750
5. OAuth 2.0 https://tools.ietf.org/html/rfc6749
*/
func AuthenticateGoogleUser(ctx context.Context, rawToken string) (*SessionInfo, error) {
	var (
		idToken *oidc.IDToken
		err     error
	)
	idToken, err = VerifyGoogleAuthIdToken(ctx, rawToken)
	if err != nil {
		return nil, err
	}
	var profile = storage.GoogleTokenClaim{}
	idToken.Claims(&profile)
	if !profile.EmailVerified {
		return nil, errors.New("email not verified")
	}
	return updateGoogleUserSession(&profile)
}

func updateGoogleUserSession(profile *storage.GoogleTokenClaim) (*SessionInfo, error) {
	var (
		uid        uint32
		sid        string
		lastAccess *time.Time
		err        error
	)
	uid, err = storage.GetOrCreateUidForGoogleUser(profile)
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	sInfo, err := storage.QuerySessionByUid(uid)
	if err != nil {
		return nil, err
	}
	if sInfo == nil || sInfo.LastAccess == nil || sInfo.Sid == "" {
		// no existing session, create new.
		sid = computeRandomSessionId(profile.Sub)
		err = storage.UpdateSession(sid, uid)
		if err != nil {
			return nil, err
		}
	} else {
		// verify if existing session timed out
		sid = sInfo.Sid
		lastAccess = sInfo.LastAccess
		if err = VerifySession(sid, uid, lastAccess); err != nil {
			return nil, err
		}
		// refresh session.
		go storage.UpdateSession(sid, uid)
	}
	return &SessionInfo{uid, sid}, nil
}
