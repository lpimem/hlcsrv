package google_auth

import (
	"log"
	"os"

	oidc "github.com/coreos/go-oidc"
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

func verifyGoogleAuthIdToken(ctx context.Context, rawToken string) (*oidc.IDToken, error) {
	idToken, err := verifier.Verify(ctx, rawToken)
	if err != nil {
		log.Println("WARN: Cannot verify token", err)
		return nil, err
	}
	return idToken, nil
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
	configure(context.Background())
}
