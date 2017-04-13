package controller

import "testing"
import "github.com/lpimem/hlcsrv/hlcmsg"
import "github.com/lpimem/hlcsrv/storage"

func TestNewNotes(t *testing.T) {
	pn := mockPageNote(1, 2, "http://example.com/index.html")
	idlist, err := newNotes(pn)
	if err != nil {
		t.Error("error: ", err)
		t.Fail()
		return
	}
	if idlist == nil {
		t.Error("should received a valid idlist")
		t.Fail()
		return
	}
	if len(idlist.Arr) != 1 {
		t.Error("idlist should contain 1 id")
		t.Fail()
	}
}

func testGetPagenote(t *testing.T, uid uint32, pid uint32, count uint32) {
	storage.ResetTestDb()
	req := &hlcmsg.Pagenote{
		Uid:    uid,
		Pageid: pid,
	}
	note := getNotes(req)
	if count == 0 {
		if note != nil {
			t.Error("expecting nil, got", note)
			t.Fail()
		}
		return
	}
	if note == nil {
		t.Error("note should be not nil")
		t.Fail()
		return
	}
	if len(note.Highlights) != int(count) {
		t.Error("note should contain", count, "highlights")
		t.Fail()
		return
	}
}

func TestGetNote(t *testing.T) {
	for _, tc := range [][]uint32{
		[]uint32{1, 0, 1},
		[]uint32{1, 1, 1},
		[]uint32{2, 1, 0},
	} {
		testGetPagenote(t, tc[0], tc[1], tc[2])
	}
}

func TestCleanUrl(t *testing.T) {
	tcs := [][]string{
		[]string{
			"http://example.com/hello?a=1&b=2#anchor&hiho=3",
			"http://example.com/hello?a=1&b=2",
		},
		[]string{
			"example.com/hello?a=1&b=2#anchor&hiho=3",
			"http://example.com/hello?a=1&b=2",
		},
		[]string{
			"https://example.com/hello?a=1&b=2#anchor&hiho=3",
			"https://example.com/hello?a=1&b=2",
		},
		[]string{
			"https://example.com/hello",
			"https://example.com/hello",
		},
		[]string{
			"https://example.com",
			"https://example.com",
		},
		[]string{
			"https://example.com/",
			"https://example.com/",
		},
		// []string{
		// 	"https://example/",
		// 	"",
		// },
	}
	for _, tc := range tcs {
		u := tc[0]
		ex := tc[1]
		got, err := cleanUrl(u)
		if err != nil {
			t.Error("error should'nt happen", err)
			t.Fail()
			continue
		}
		if got != ex {
			t.Error("url should be properly cleaned. expect: ", ex, "| got:", got)
			t.Fail()
		}
	}
}

func mockPageNote(uid uint32, pid uint32, url string) *hlcmsg.Pagenote {
	return &hlcmsg.Pagenote{
		Uid:    uid,
		Pageid: pid,
		Url:    url,
		Highlights: []*hlcmsg.RangeMeta{
			&hlcmsg.RangeMeta{
				Anchor:      "/",
				Start:       "/1",
				StartOffset: 0,
				End:         "/2",
				EndOffset:   3,
			},
		},
	}
}
