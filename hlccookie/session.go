package hlccookie

import (
	"net/http"

	"strconv"

	"github.com/lpimem/hlcsrv/conf"
	"github.com/lpimem/hlcsrv/storage"
)

// GetRequestUID extracts user id from request cookie
func GetRequestUID(r *http.Request) (storage.UserID, error) {
	uid32, err := getCookieAsUInt32(r, conf.SessionKeyUser())
	if err != nil {
		return 0, err
	}
	return storage.UserID(uid32), err
}

func getCookieAsUInt32(r *http.Request, key string) (uint32, error) {
	var (
		err    error
		cookie *http.Cookie
	)
	if cookie, err = r.Cookie(key); err == nil {
		value := cookie.Value
		valInt, converr := strconv.ParseUint(value, 10, 32)
		if converr == nil {
			return uint32(valInt), nil
		}
		return 0, converr
	}
	return 0, err
}

// SetAuthCookies set authentication cookie for response.
func SetAuthCookies(w http.ResponseWriter, sid string, uid storage.UserID) {
	http.SetCookie(w, &http.Cookie{Name: conf.SessionKeySID(), Value: sid})
	http.SetCookie(w, &http.Cookie{Name: conf.SessionKeyUser(), Value: strconv.FormatUint(uint64(uid), 10)})
}
