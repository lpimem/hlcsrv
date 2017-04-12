package storage

import (
	"fmt"
	"testing"
)

func TestSqliteStorage_QueryUidByGoogleId(t *testing.T) {
	ResetTestDb()
	uid := storage.QueryUidByGoogleId("100000")
	if uid != 1 {
		t.Error("should return uid 1, got:", uid)
		t.Fail()
	}

	uid = storage.QueryUidByGoogleId("100001")
	if uid != 0 {
		t.Error("should return uid 0 (not found), got:", uid)
		t.Fail()
	}
}

func TestSqliteStorage_NewUserByGoogleId(t *testing.T) {
	ResetTestDb()
	uid, err := storage.NewUserByGoogleId("100001", "abc@example.com")
	if err != nil {
		t.Error("NewUserByGoogleId: should be success", err)
		t.Fail()
		return
	}
	if uid <= 0 {
		t.Error("NewUserByGoogleId: should get uid > 0, got:", uid)
		t.Fail()
	}

	uid, err = storage.NewUserByGoogleId("100001", "abc@example.com")
	if err == nil {
		t.Error("Should report error for duplicated gid")
		t.Fail()
		return
	}
}

func TestGetOrCreateUidForGoogleUser(t *testing.T) {
	ResetTestDb()
	var (
		uid uint32
		err error
	)
	uid, err = GetOrCreateUidForGoogleUser("100000", "abc@example.com")
	if uid != 1 || err != nil {
		t.Error("should return uid 1, got:", uid, err)
		t.Fail()
	}

	uid, err = GetOrCreateUidForGoogleUser("100010", "abc2@example.com")
	if uid < 1 || err != nil {
		t.Error("should return uid > 0, got:", uid, err)
		t.Fail()
	}

	uid_2, err := GetOrCreateUidForGoogleUser("100010", "abc2@example.com")
	if uid_2 != uid || err != nil {
		t.Error("should return uid=", uid, "got:", uid_2, err)
		t.Fail()
	}

	uid, err = GetOrCreateUidForGoogleUser("100012", "abc2@example.com")
	fmt.Println("expected msg for duplicated email:", err)
	if err == nil {
		t.Error("should raise error for duplicated email, got", uid, err)
		t.Fail()
	}
}
