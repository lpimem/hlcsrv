package conf

import (
	"net/http"
	"net/url"
	"os"

	"github.com/go-playground/log"
)

var debugFlag = false

const (
	_HLC_NEXT = "hlc.next"
)

// IsDebug returns true should the app run in debugging mode.
func IsDebug() bool {
	if os.Getenv("HLC_DEBUG") == "1" {
		return true
	}
	return debugFlag
}

// SetDebug change debug option
func SetDebug(option bool) {
	debugFlag = option
}

/*SessionSecret is the secret used to sign your cookies
 * TODO : Change the value before deployment.
 * On Mac or Linux, you can use the following command to generate one:
 *   ```
 *   env LC_CTYPE=C tr -dc "a-zA-Z0-9-_\$\?" < /dev/urandom | fold -w 64 | head -n 1
 *   ```
 */
func SessionSecret() string {
	return os.Getenv("HLC_SESSION_SECRET")
}

// Page should be a request parameter, not a cookie
//
///*SessionKeyPage is the random seed for key name for Page Id.
// */
//func SessionKeyPage() string {
//	return "PkNMgRN4kx_uxrmduaVK1AyL8L7aCxhVDHmSPWHpp9v6UD-BJGEMPMbRPQaa9Dc1"
//}

/*SessionKeyUser is the random seed for key name for User Id.
 */
func SessionKeyUser() string {
	return os.Getenv("HLC_SESSION_KEY_USER")
}

/*SessionKeySID is the key for session id
 */
func SessionKeySID() string {
	return os.Getenv("HLC_SESSION_KEY_SID")
}

/*SessionValidHours defines how long a session could be idle for.
 */
func SessionValidHours() float64 {
	return 24 * 30
}

/*GoogleSignInAppID extract google client id from $GOOGLE_OAUTH2_CLIENT_ID
environment variable.
*/
func GoogleSignInAppID() string {
	return os.Getenv("GOOGLE_OAUTH2_CLIENT_ID")
}

/*GoogleOAuthRedirectURL returns the OAuth2.0 redirect URL for 3-legged authentication.
(Currently this is feature is not implemented)
*/
func GoogleOAuthRedirectURL() string {
	return "http://127.0.0.1:5556/auth/google/callback"
}

func LoginURL() string {
	return "/static/login.html"
}

func RedirectToLogin(next string, w http.ResponseWriter, r *http.Request) {
	RedirectTo(LoginURL(), next, w, r)
}

func RedirectTo(redirectUrl string, next_url string, w http.ResponseWriter, r *http.Request) {
	log.Debug("Redirecting to ", redirectUrl, " -> ", next_url)
	http.SetCookie(w, &http.Cookie{Name: _HLC_NEXT, Value: next_url, Path: "/"})
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func SetNext(w http.ResponseWriter, u *url.URL) {
	encodedPath := EncodePath(u)
	http.SetCookie(w, &http.Cookie{Name: _HLC_NEXT, Value: encodedPath, Path: "/"})
}

func GetNext(r *http.Request) (string, error) {
	var (
		err error
		c   *http.Cookie
	)
	if c, err = r.Cookie(_HLC_NEXT); err == nil {
		return c.Value, err
	}
	return "", err
}

func EncodePath(u *url.URL) string {
	log.Debug("Encoding:", u)
	path := u.Path
	if u.RawQuery != "" {
		path += "%3F"
		path += u.RawQuery
	}
	return path
}
