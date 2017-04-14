package google_auth

import "context"

func validToken() string {
	return ""
}

//func TestVerifyGoogleAuthToken(t *testing.T) {
//	return
//	var token = validToken() // fetch/generate a valid token
//	ctx := context.Background()
//	if id, err := VerifyGoogleAuthIdToken(ctx, token); err != nil {
//		t.Error("token should be valid", err)
//		t.Fail()
//	} else {
//		util.Log(id.Audience, id.Subject, id.Expiry, id.Issuer, id.IssuedAt)
//	}
//}

func init() {
	// setup test client id and secret
	configure(context.Background())
}
