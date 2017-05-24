package storage

import (
	"fmt"
	"testing"
)

func TestQuerySession(t *testing.T) {
	type testcase struct {
		Name string
		Sid  string
		Uid  uint32
		Pass bool
	}
	var tcs = []*testcase{
		&testcase{"valid query", "fake_session_id", 10, true},
		&testcase{"missed query 1", "fake_session_id", 11, false},
		&testcase{"missed query 2", "fake_session_id 2", 10, false},
		&testcase{"invalid uid", "fake_session_id 2", 0, false},
		&testcase{"invalid sid", "", 1, false},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			lastAccess, err := QuerySession(tc.Sid, tc.Uid)
			suc := lastAccess != nil && err == nil
			if suc != tc.Pass {
				fmt.Println(err)
				t.Fail()
			}
		})
	}
}

func TestQuerySessionByUID(t *testing.T) {
	var tcs = []struct {
		name string
		uid  int
		pass bool
	}{
		{"valid query", 10, true},
		{"no session", 2, false},
		{"invalid uid", 0, false},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r, err := QuerySessionByUID(uint32(tc.uid))
			suc := r != nil && err == nil && r.LastAccess != nil && r.Sid != ""
			if suc != tc.pass {
				fmt.Println(err)
				t.Fail()
			}
		})
	}
}

func TestUpdateSession(t *testing.T) {
	var tcs = []struct {
		name string
		sid  string
		uid  int
		pass bool
	}{
		{"update existing session", "fake_session_id", 10, true},
		{"update session id", "fake_session_id 2", 10, true},
		{"create session id", "fake_session_id 3", 1, true},
		{"missing uid", "fake_session_id 3", 0, false},
		{"missing session id", "", 10, false},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := UpdateSession(tc.sid, uint32(tc.uid))
			if (err == nil) != tc.pass {
				fmt.Println(err)
				t.Fail()
				return
			}
			lastAccess, err := QuerySession(tc.sid, uint32(tc.uid))
			suc := lastAccess != nil && err == nil
			if suc != tc.pass {
				fmt.Println(err)
				t.Fail()
			}
		})
	}
}
