package auth

import (
	"errors"
	"os"

	oidc "github.com/coreos/go-oidc"
	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/conf"
	"golang.org/x/net/context"
)

var (
	clientID     = os.Getenv("GOOGLE_OAUTH2_CLIENT_ID")
	clientSecret = os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET")
	verifier     *oidc.IDTokenVerifier
	//config       oauth2.Config
)

// VerifyGoogleAuthIDToken parses and validate the raw google sign-in token string.
func VerifyGoogleAuthIDToken(ctx context.Context, rawToken string) (*oidc.IDToken, error) {
	if err := ensureGoogleAppConfig(); err != nil {
		panic(err)
	}
	idToken, err := verifier.Verify(ctx, rawToken)
	if err != nil {
		log.Warn("Cannot verify google token: ", err)
		log.Info(rawToken)
		return nil, err
	}
	if err := verifyAud(ctx, idToken); err != nil {
		log.Warn(" Cannot verify google token: ", err)
		log.Info(rawToken)
		return nil, err
	}
	return idToken, nil
}

func verifyAud(ctx context.Context, idToken *oidc.IDToken) error {
	for _, aud := range idToken.Audience {
		if aud == conf.GoogleSignInAppID() {
			return nil
		}
	}
	return errors.New("audience dismatch")
}

func configure(ctx context.Context) error {
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		log.Info("error setting up google authentication service", err)
		return err
	}
	oidcConfig := &oidc.Config{
		ClientID:       clientID,
		SkipNonceCheck: true,
	}
	verifier = provider.Verifier(oidcConfig)
	//config = oauth2.Config{
	//	ClientID:     clientID,
	//	ClientSecret: clientSecret,
	//	Endpoint:     provider.Endpoint(),
	//	RedirectURL:  conf.GoogleOAuthRedirectURL(),
	//	Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	//}
	return nil
}

func ensureGoogleAppConfig() error {
	if clientID == "" || clientSecret == "" {
		msg := "Cannot Find Google App Config"
		log.Warn(msg)
		return errors.New(msg)
	}
	return nil
}

func init() {
	_ = ensureGoogleAppConfig()
	log.Info("google client id:", clientID)
	log.Info("google client secret:", len(clientSecret))
	if err := configure(context.Background()); err != nil {
		panic(err)
	}
}
