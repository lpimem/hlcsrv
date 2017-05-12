package auth

import (
	"context"
	"errors"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/lpimem/hlcsrv/storage"
)

/*AuthenticateGoogleUser authenticates a google user with rawToken
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
	idToken, err = VerifyGoogleAuthIDToken(ctx, rawToken)
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

func updateGoogleUserSession(profile *storage.GoogleTokenClaim) (session *SessionInfo, err error) {
	var (
		uid        uint32
		sid        string
		lastAccess *time.Time
	)
	uid, err = storage.GetOrCreateUIDForGoogleUser(profile)
	if err != nil {
		return
	}
	lastSession, err := storage.QuerySessionByUID(uid)
	if err != nil {
		return
	}
	if lastSession == nil || lastSession.LastAccess == nil || lastSession.Sid == "" {
		// no existing session, create new.
		sid, err = createSessionForGoogleUser(profile.Sub, uid)
	} else {
		// verify if existing session timed out
		sid = lastSession.Sid
		lastAccess = lastSession.LastAccess
		if err = VerifySession(sid, uid, lastAccess); err != nil {
			sid, err = createSessionForGoogleUser(profile.Sub, uid)
		} else {
			err = storage.UpdateSession(sid, uid)
		}
	}
	session = &SessionInfo{uid, sid}
	return
}

func createSessionForGoogleUser(sub string, uid uint32) (sid string, err error) {
	sid = computeRandomSessionID(sub)
	err = storage.UpdateSession(sid, uid)
	return
}
