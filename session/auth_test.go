package session

import (
	"testing"

	"net/http/httptest"

	"net/http"

	"strconv"

	"fmt"

	"time"

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

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var (
				req *http.Request
			)
			req = httptest.NewRequest("GET", "/", nil)
			if tc.Sid != "" {
				req.AddCookie(&http.Cookie{
					Name:  conf.SessionKeySID(),
					Value: tc.Sid,
				})
				fmt.Println("req add cookie:", conf.SessionKeySID(), tc.Sid)
				if _, err := req.Cookie(conf.SessionKeySID()); err != nil {
					t.Error("cookie not set", err)
					t.Fail()
				}
			}
			if tc.Uid > 0 {
				req.AddCookie(&http.Cookie{
					Name:  conf.SessionKeyUser(),
					Value: strconv.FormatUint(uint64(tc.Uid), 10),
				})
				fmt.Println("req add cookie:", conf.SessionKeyUser(), tc.Uid)
			}
			req, err := Authenticate(req)
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

func init() {
	storage.ResetTestDb()
}
