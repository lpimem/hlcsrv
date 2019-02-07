package storage

import (
	"github.com/lpimem/hlcsrv/conf"
	"testing"
)

func TestHasAccess(t *testing.T) {
	ResetTestDb()
	type testcase struct {
		Name string
		User UserID
		URI  string
		Suc  bool
	}

	testcases := []*testcase{
		&testcase{"NormalBaseCase", 1, "/example", true},
		&testcase{"Should have access to sub URI", 1, "/example/C", true},
		&testcase{"Should have access to sub URI 2", 1, "/example/Boy/C", true},
		&testcase{"Should have NO access to unlisted URI 1 ", 1, "/exampleA", false},
		&testcase{"Should have NO access to unlisted URI 2 ", 2, "/example/A", false},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			access, err := Permission.HasAccess(tc.User, tc.URI)
			if err != nil {
				t.Error("unexpected error", err)
				t.Fail()
			}
			if access != tc.Suc {
				t.Error("Expecting ", tc.Suc, " got ", access)
				t.Fail()
			}
		})
	}
}

func TestGrant(t *testing.T) {
	ResetTestDb()
	type testcase struct {
		Name       string
		User       UserID
		URI        string
		TestURIPos []string
		TestURINeg []string
	}
	testcases := []*testcase{
		&testcase{"Should grant new access", 1, "/localhost/Can",
			[]string{"/localhost/Can", "/localhost/Can/D"},
			[]string{"/localhost/D", "/localhost/Cannon", "/"},
		},
		&testcase{"Duplicate granting should be fine", 1, "/localhost/Can",
			[]string{"/localhost/Can", "/localhost/Can/D"},
			[]string{"/localhost/D", "/localhost/Cannon", "/"},
		},
		&testcase{"Grant higher level", 2, "/example",
			[]string{"/example", "/example/A", "/example/Boy", "/example/Can"},
			[]string{"/example1", "/exampl", "/"},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			err := Permission.Grant(tc.User, tc.URI)
			if err != nil {
				t.Error("unexpected error: ", err)
				t.Fail()
			}
			for _, testURI := range tc.TestURIPos {
				if acc, err := Permission.HasAccess(tc.User, testURI); err != nil || !acc {
					if err != nil {
						t.Error("unexpected error: ", err)
					}
					t.Errorf("User %d should have access to %s", tc.User, testURI)
					t.Fail()
				}
			}
			for _, testURI := range tc.TestURINeg {
				if acc, err := Permission.HasAccess(tc.User, testURI); err != nil || acc {
					if err != nil {
						t.Error("unexpected error: ", err)
					}
					t.Errorf("User %d shouldn't have access to %s", tc.User, testURI)
					t.Fail()
				}
			}
		})
	}
}

func init() {
	conf.SetDebug(true)
}
