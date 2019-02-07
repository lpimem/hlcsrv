package auth

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/conf"
	"github.com/lpimem/hlcsrv/storage"
)

type authCase struct {
	Name string
	Sid  string
	UID  uint32
	Suc  bool
}

func TestVerifySession(t *testing.T) {
	now := time.Now()
	invalidDate := now.Add(-1 * (time.Duration(conf.SessionValidHours()) + 1) * time.Hour)
	testcases := []struct {
		name string
		sid  string
		UID  int
		la   *time.Time
		pass bool
	}{
		{"valid session without last access time", "fake_session_id", 10, nil, true},
		{"valid session with last access time", "fake_session_id", 10, &now, true},
		{"session timeout", "fake_session_id", 10, &invalidDate, false},
		{"missing UID", "fake_session_id", 0, &invalidDate, false},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := VerifySession(tc.sid, storage.UserID(uint32(tc.UID)), tc.la)
			if (err == nil) != tc.pass {
				t.Error("ERROR", err)
				t.Fail()
			}
		})
	}
}

func TestAuthenticate(t *testing.T) {

	testCases := []*authCase{
		&authCase{"valid request", "fake_session_id", uint32(10), true},
		&authCase{"bare request is not authenticated", "", uint32(0), false},
		&authCase{"UID and sid dismatch", "fake_session_id", uint32(1), false},
		&authCase{"UID missing", "fake_session_id", uint32(0), false},
		&authCase{"sid missing", "", uint32(10), false},
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
					req, err = setByCookie(req, tc.UID, tc.Sid)
					if err != nil {
						t.Error(err)
						t.Fail()
						return
					}
				} else {
					req, err = setByHeader(req, tc.UID, tc.Sid)
				}
				req, err = Authenticate(req)
				if IsAuthenticated(req) != tc.Suc {
					log.Debug(req.Context().Value(AUTHENTICATED))
					log.Debug(req.Context().Value(USER_ID))
					log.Debug(req.Context().Value(SESSION_ID))
					log.Debug(req.Context().Value(REASON))
					if err != nil {
						t.Error(err)
					}
					t.Fail()
				}
			})
		}
	}
}

func TestAuthorizeAdmin(t *testing.T) {
	type testcase struct {
		Name string
		Sid  string
		UID  uint32
		URI  string
		Suc  bool
	}

	testCases := []*testcase{
		&testcase{"valid admin request from admin", "fake_session_id_for_1", uint32(1), "/admin", true},
		&testcase{"valid admin request from admin", "fake_session_id_for_1", uint32(1), "/static/admin", true},
		&testcase{"valid admin request from admin", "fake_session_id_for_1", uint32(1), "/static/admin/anything", true},
		&testcase{"invalid admin request from user", "fake_session_id", uint32(10), "/admin", false},
		&testcase{"invalid admin request from user", "fake_session_id", uint32(10), "/static/admin", false},
		&testcase{"invalid admin request from user", "fake_session_id", uint32(10), "/static/admin/anything", false},
		&testcase{"admin request from non-user", "", 0, "/admin", false},
		&testcase{"admin request from non-user", "", 0, "/static/admin", false},
		&testcase{"admin request from non-user", "", 0, "/static/admin/anything", false},
		&testcase{"invalid admin request from user", "fake_session_id", uint32(10), "/admin", false},
		&testcase{"admin request from non-user", "", 0, "/admin", false},
		&testcase{"valid normal request from admin", "fake_session_id_for_1", uint32(1), "/random", true},
		&testcase{"valid normal request from user", "fake_session_id", uint32(10), "/random", true},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var (
				req *http.Request
				err error
			)
			req = httptest.NewRequest("GET", tc.URI, nil)
			req, err = setByCookie(req, tc.UID, tc.Sid)
			if err != nil {
				t.Error(err)
				t.Fail()
				return
			}
			req, err = Authenticate(req)
			pass := IsAuthenticated(req) && err == nil
			if pass != tc.Suc {
				if err != nil {
					t.Error("unexpected error :", err)
				}
				t.Fail()
			}
		})
	}
}

func setByCookie(req *http.Request, UID uint32, sid string) (*http.Request, error) {
	if sid != "" {
		req.AddCookie(&http.Cookie{
			Name:  conf.SessionKeySID(),
			Value: sid,
		})
		log.Info("req add cookie:", conf.SessionKeySID(), sid)
		if _, err := req.Cookie(conf.SessionKeySID()); err != nil {
			return nil, err
		}
	}
	if UID > 0 {
		req.AddCookie(&http.Cookie{
			Name:  conf.SessionKeyUser(),
			Value: strconv.FormatUint(uint64(UID), 10),
		})
		log.Info("req add cookie:", conf.SessionKeyUser(), UID)
	}
	return req, nil
}

func setByHeader(req *http.Request, UID uint32, sid string) (*http.Request, error) {
	if UID > 0 {
		req.Header.Set(HUSER_ID, strconv.FormatUint(uint64(UID), 10))
	}
	if sid != "" {
		req.Header.Set(HSESSION_ID, sid)
	}
	return req, nil
}

func init() {
	storage.ResetTestDb()
}
