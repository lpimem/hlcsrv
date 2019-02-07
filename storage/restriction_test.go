package storage

import (
	"testing"
)

func TestRestrictionHas(t *testing.T) {
	type testcase struct {
		Name       string
		URI        string
		Restricted bool
	}

	testcases := []*testcase{
		&testcase{"Exact match should be restricted", "/example/", true},
		&testcase{"Match without trailing slash be restricted", "/example", true},
		&testcase{"Prefix match should be restricted ", "/example/a", true},
		&testcase{"Nested prefix match should be restricted", "/example/a/b/c", true},
		&testcase{"Unmatch should not be restricted", "/localhost/a/b/c", false},
	}
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			restricted, err := Restriction.Has(tc.URI)
			if restricted != tc.Restricted || err != nil {
				if err != nil {
					t.Errorf("Unexpected error: %s", err)
				}
				msg := "unrestricted"
				if tc.Restricted {
					msg = "restricted"
				}
				t.Errorf("%s should be %s", tc.URI, msg)
				t.Fail()
			}
		})
	}
}

func TestRestrictionAdd(t *testing.T) {
	const newURI = "/localhost/ban"
	var (
		res bool
		err error
	)
	if res, err = Restriction.Has(newURI); res || err != nil {
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		}
		t.Errorf("Pre-condition dismatch. Expecting no restrictions on %s", newURI)
		t.FailNow()
	}
	if err = Restriction.Add(newURI); err != nil {
		t.Errorf("Unexpected error %s", err)
		t.FailNow()
	}
	if res, err = Restriction.Has(newURI); !res || err != nil {
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		}
		t.Errorf("Post-condition dismatch. Expecting new restrictions on %s", newURI)
		t.FailNow()
	}
}
