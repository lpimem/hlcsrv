package google_auth

import (
	"log"
	"os"

	"errors"

	oidc "github.com/coreos/go-oidc"
	"github.com/lpimem/hlcsrv/conf"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

var (
	clientID     = os.Getenv("GOOGLE_OAUTH2_CLIENT_ID")
	clientSecret = os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET")
	verifier     *oidc.IDTokenVerifier
	config       oauth2.Config
)

const redirectURL = "http://127.0.0.1:5556/auth/google/callback"

func VerifyGoogleAuthIdToken(ctx context.Context, rawToken string) (*oidc.IDToken, error) {
	idToken, err := verifier.Verify(ctx, rawToken)
	if err != nil {
		log.Println("WARN: Cannot verify token: ", err)
		log.Println(rawToken)
		return nil, err
	}
	if err := verifyAud(ctx, idToken); err != nil {
		log.Println("WARN: Cannot verify token: ", err)
		log.Println(rawToken)
		return nil, err
	}
	return idToken, nil
}

func verifyAud(ctx context.Context, idToken *oidc.IDToken) error {
	for _, aud := range idToken.Audience {
		if aud == conf.GoogleSignInAppId() {
			return nil
		}
	}
	return errors.New("audience dismatch")
}

func configure(ctx context.Context) error {
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		log.Fatal("error setting up google authentication service", err)
		return err
	}
	oidcConfig := &oidc.Config{
		ClientID:       clientID,
		SkipNonceCheck: true,
	}
	verifier = provider.Verifier(oidcConfig)
	config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  redirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
	return nil
}

func init() {
	log.Println("google client id:", clientID)
	log.Println("google client secret:", clientSecret[:4])
	configure(context.Background())
}
