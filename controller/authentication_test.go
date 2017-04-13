package controller

import (
	"fmt"
	"testing"

	"github.com/lpimem/hlcsrv/session"
	"github.com/lpimem/hlcsrv/storage"
)

func TestUpdateGoogleUserSession(t *testing.T) {
	testcases := []struct {
		name  string
		gid   string
		email string
		pass  bool
	}{
		{"existing gid", "example@google.com", "example@google.com", true},
		{"new gid", "example_2@google.com", "example_2@google.com", true},
		{"duplicate email", "example_3@google.com", "example@google.com", false},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var suc bool
			sInfo, err := updateGoogleUserSession(tc.gid, tc.email)
			suc = sInfo != nil && err == nil
			if suc != tc.pass {
				fmt.Println(err)
				t.Fail()
			}
			if sInfo == nil {
				return
			}

			lastAccess, err := storage.QuerySession(sInfo.Sid, sInfo.Uid)
			suc = !session.IsSessionTimeout(*lastAccess) && err == nil
			if suc != tc.pass {
				fmt.Println(err)
				t.Fail()
			}
		})
	}
}
