package hlccookie

import (
	"net/http"

	"strconv"

	"github.com/lpimem/hlcsrv/conf"
)

func GetRequestUID(r *http.Request) (uint32, error) {
	return getCookieAsUInt32(r, conf.SessionKeyUser())
}

func GetPageId(r *http.Request) (uint32, error) {
	return getCookieAsUInt32(r, conf.SessionKeyPage())
}

func getCookieAsUInt32(r *http.Request, key string) (uint32, error) {
	var (
		err    error
		cookie *http.Cookie
	)
	if cookie, err = r.Cookie(key); err == nil {
		value := cookie.Value
		if valInt, converr := strconv.ParseUint(value, 10, 32); converr == nil {
			return uint32(valInt), nil
		} else {
			return 0, converr
		}
	}
	return 0, err
}
