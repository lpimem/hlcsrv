package storage

import (
	"testing"
)

func TestSqliteStorage_QueryUIDByGoogleID(t *testing.T) {
	ResetTestDb()
	uid := storage.QueryUIDByGoogleID("100000")
	if uid != 1 {
		t.Error("should return uid 1, got:", uid)
		t.Fail()
	}

	uid = storage.QueryUIDByGoogleID("100001")
	if uid != 0 {
		t.Error("should return uid 0 (not found), got:", uid)
		t.Fail()
	}
}

func TestSqliteStorage_NewUserByGoogleID(t *testing.T) {
	ResetTestDb()
	tcs := []struct {
		name    string
		profile *GoogleTokenClaim
		pass    bool
	}{
		{
			"valid new google user",
			&GoogleTokenClaim{
				"abc@example.com",
				true,
				"ABC",
				"100001",
				"http://example.com/a.png",
				"EN"},
			true,
		},
		{
			"another new google user",
			&GoogleTokenClaim{
				"abcde@example.com",
				true,
				"ABCD",
				"100002",
				"http://example.com/a.png",
				"EN"},
			true,
		},
		{
			"new google user with duplicate gid",
			&GoogleTokenClaim{
				"abcd@example.com",
				true,
				"ABCDE",
				"100001",
				"http://example.com/a.png",
				"EN"},
			false,
		},
		{
			"new google user with duplicate name",
			&GoogleTokenClaim{
				"abcd123@example.com",
				true,
				"ABC",
				"100003",
				"http://example.com/a.png",
				"EN"},
			false,
		},
		{
			"new google user with duplicate email",
			&GoogleTokenClaim{
				"abc@example.com",
				true,
				"ABCDEF",
				"100004",
				"http://example.com/a.png",
				"EN"},
			false,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			uid, err := storage.NewUserByGoogleProfile(tc.profile)
			suc := uid > 0 && err == nil
			if suc != tc.pass {
				t.Error(err)
				t.Fail()
			}
		})
	}
}

func TestGetOrCreateUIDForGoogleUser(t *testing.T) {
	ResetTestDb()
	var (
		uid uint32
		err error
	)
	uid, err = GetOrCreateUIDForGoogleUser(&GoogleTokenClaim{
		Email: "abc@example.com",
		Sub:   "100000",
		Name:  "abc",
	})
	if uid != 1 || err != nil {
		t.Error("should return uid 1, got:", uid, err)
		t.Fail()
	}

	uid, err = GetOrCreateUIDForGoogleUser(&GoogleTokenClaim{
		Email: "abc2@example.com",
		Sub:   "100010",
		Name:  "abc",
	})
	if uid < 1 || err != nil {
		t.Error("should return uid > 0, got:", uid, err)
		t.Fail()
	}

	uid2, err := GetOrCreateUIDForGoogleUser(&GoogleTokenClaim{
		Email: "abc2@example.com",
		Sub:   "100010",
		Name:  "abc",
	})
	if uid2 != uid || err != nil {
		t.Error("should return uid=", uid, "got:", uid2, err)
		t.Fail()
	}

	uid, err = GetOrCreateUIDForGoogleUser(&GoogleTokenClaim{
		Email: "abc2@example.com",
		Sub:   "100012",
		Name:  "abc",
	})
	if err == nil {
		t.Error("should raise error for duplicated email, got", uid, err)
		t.Fail()
	}
}
