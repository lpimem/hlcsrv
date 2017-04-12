package storage

import "testing"

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
