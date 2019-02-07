package storage

import (
	"testing"
)

func TestUserAll(t *testing.T) {
	ResetTestDb()
	all, err := User.All(100, 0)
	if err != nil {
		t.Errorf("Cannot query all users: %s", err)
		t.FailNow()
	}
	expected := 3
	if len(all) != expected {
		t.Errorf("Expecting %d users, got %d", expected, len(all))
	}
	for _, uInfo := range all {
		if uInfo.ID < 1 {
			t.Errorf("User's info shouldn't be < 1, %v", uInfo)
		}
		if len(uInfo.Email) < 1 {
			t.Errorf("Invalid email address, %v", uInfo)
		}
	}
}
