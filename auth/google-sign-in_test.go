package auth

import (
	"fmt"
	"testing"

	"github.com/lpimem/hlcsrv/storage"
)

func TestUpdateGoogleUserSession(t *testing.T) {
	storage.ResetTestDb()
	testcases := []struct {
		name    string
		profile *storage.GoogleTokenClaim
		pass    bool
	}{
		{
			"existing gid",
			&storage.GoogleTokenClaim{
				Email: "example@google.com",
				Sub:   "100000",
				Name:  "Test User",
			},
			true,
		},
		{
			"new gid",
			&storage.GoogleTokenClaim{
				Email: "example_2@google.com",
				Sub:   "1000011",
				Name:  "Test User 11",
			},
			true,
		},
		{
			"duplicate email",
			&storage.GoogleTokenClaim{
				Email: "example_2@google.com",
				Sub:   "1000112",
				Name:  "Test User 3",
			},
			false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var suc bool
			sInfo, err := updateGoogleUserSession(tc.profile)
			suc = sInfo != nil && err == nil
			if suc != tc.pass {
				fmt.Println(err)
				t.Fail()
			}
			if sInfo == nil {
				return
			}
			lastAccess, err := storage.QuerySession(sInfo.Sid, sInfo.Uid)
			suc = err == nil && !IsSessionTimeout(*lastAccess)
			if suc != tc.pass {
				fmt.Println(err)
				t.Fail()
			}
		})
	}
}
