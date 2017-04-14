package session

import (
	"testing"

	"net/http/httptest"

	"net/http"

	"fmt"

	"time"

	"strconv"

	"github.com/lpimem/hlcsrv/conf"
	"github.com/lpimem/hlcsrv/storage"
)

type auth_case struct {
	Name string
	Sid  string
	Uid  uint32
	Suc  bool
}

func TestVerifySession(t *testing.T) {
	now := time.Now()
	invalidDate := now.Add(-1 * (time.Duration(conf.SessionValidHours()) + 1) * time.Hour)
	testcases := []struct {
		name string
		sid  string
		uid  int
		la   *time.Time
		pass bool
	}{
		{"valid session without last access time", "fake_session_id", 10, nil, true},
		{"valid session with last access time", "fake_session_id", 10, &now, true},
		{"session timeout", "fake_session_id", 10, &invalidDate, false},
		{"missing uid", "fake_session_id", 0, &invalidDate, false},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := VerifySession(tc.sid, uint32(tc.uid), tc.la)
			if (err == nil) != tc.pass {
				fmt.Println("ERROR", err)
				t.Fail()
			}
		})
	}
}

func TestAuthenticate(t *testing.T) {

	testCases := []*auth_case{
		&auth_case{"valid request", "fake_session_id", uint32(10), true},
		&auth_case{"bare request is not authenticated", "", uint32(0), false},
		&auth_case{"uid and sid dismatch", "fake_session_id", uint32(1), false},
		&auth_case{"uid missing", "fake_session_id", uint32(0), false},
		&auth_case{"sid missing", "", uint32(10), false},
	}

	for _, byCookie := range []bool{true, false} {
		for _, tc := range testCases {
			var tcname string
			if byCookie {
				tcname = tc.Name + "[cookie]"
			} else {
				tcname = tc.Name + "[header]"
			}
			t.Run(tcname, func(t *testing.T) {
				var (
					req *http.Request
					err error
				)
				req = httptest.NewRequest("GET", "/", nil)
				if byCookie {
					req, err = setByCookie(req, tc.Uid, tc.Sid)
					if err != nil {
						t.Error(err)
						t.Fail()
						return
					}
				} else {
					req, err = setByHeader(req, tc.Uid, tc.Sid)
				}
				req, err = Authenticate(req)
				if err != nil {
					fmt.Println("Authenticate should return no error")
					t.Fail()
				}
				if IsAuthenticated(req) != tc.Suc {
					fmt.Println(req.Context().Value(AUTHENTICATED))
					fmt.Println(req.Context().Value(USER_ID))
					fmt.Println(req.Context().Value(SESSION_ID))
					fmt.Println(req.Context().Value(REASON))
					t.Fail()
				}
			})
		}
	}

}

func setByCookie(req *http.Request, uid uint32, sid string) (*http.Request, error) {
	if sid != "" {
		req.AddCookie(&http.Cookie{
			Name:  conf.SessionKeySID(),
			Value: sid,
		})
		fmt.Println("req add cookie:", conf.SessionKeySID(), sid)
		if _, err := req.Cookie(conf.SessionKeySID()); err != nil {
			return nil, err
		}
	}
	if uid > 0 {
		req.AddCookie(&http.Cookie{
			Name:  conf.SessionKeyUser(),
			Value: strconv.FormatUint(uint64(uid), 10),
		})
		fmt.Println("req add cookie:", conf.SessionKeyUser(), uid)
	}
	return req, nil
}

func setByHeader(req *http.Request, uid uint32, sid string) (*http.Request, error) {
	if uid > 0 {
		req.Header.Set(HUSER_ID, strconv.FormatUint(uint64(uid), 10))
	}
	if sid != "" {
		req.Header.Set(HSESSION_ID, sid)
	}
	return req, nil
}

func init() {
	storage.ResetTestDb()
}
