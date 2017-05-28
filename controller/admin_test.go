package controller

import (
	"bytes"
	"encoding/json"
	"github.com/lpimem/hlcsrv/storage"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func TestAdminUsers(t *testing.T) {
	var w *httptest.ResponseRecorder
	get := httptest.NewRequest("GET", "/admin/users", nil)
	get = fakeAuthenticateion(get)
	w = httptest.NewRecorder()
	Admin.Users(w, get)
	if w.Code != http.StatusOK {
		t.Errorf("Error retriving users %s", w.Body.String())
		t.Fail()
	}
	users := []storage.UserInfo{}
	if err := json.Unmarshal(w.Body.Bytes(), &users); err != nil {
		t.Error("error parsing response", err)
		t.Fail()
	}
	const expected = 3
	if len(users) != expected {
		t.Errorf("Should have %d user info, got %d", expected, len(users))
	}
}

func TestAdminGrant(t *testing.T) {
	var w *httptest.ResponseRecorder
	var uid storage.UserID = 10
	var uri = "/localhost/new/permission"
	if acc, err := storage.Permission.HasAccess(uid, uri); err != nil || acc {
		if acc {
			t.Errorf("Pre-condition dismatch, %d should have no access to %s", uid, uri)
		}
		if err != nil {
			t.Errorf("Unexpected error querying permission: %s", err)
		}
		t.FailNow()
	}
	var form url.Values = make(url.Values)
	form.Set("uid", strconv.Itoa(int(uid)))
	form.Set("uri", uri)
	request := httptest.NewRequest("POST", "/admin/grant", bytes.NewBufferString(form.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request = fakeAdminAuthenticateion(request)
	w = httptest.NewRecorder()
	Admin.Grant(w, request)
	if w.Code != http.StatusOK {
		t.Error("Error : ", w.Code, " body: ", w.Body.String())
		t.FailNow()
	}
	if acc, err := storage.Permission.HasAccess(uid, uri); err != nil || !acc {
		if !acc {
			t.Error("should have access after granting")
		}
		if err != nil {
			t.Error("error querying permission ", err)
		}
		t.Fail()
	}
}
