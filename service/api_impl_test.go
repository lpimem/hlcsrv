package service

import "testing"
import "github.com/lpimem/hlcsrv/hlcmsg"
import "github.com/lpimem/hlcsrv/storage"

func TestNewNotes(t *testing.T) {
	// newNotes(pn *hlcmsg.Pagenote) *hlcmsg.IdList
	pn := mockPageNote()
	idlist := newNotes(pn)
	if idlist == nil {
		t.Error("should received a valid idlist")
		t.Fail()
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
		got := cleanUrl(u)
		if got != ex {
			t.Error("url not cleaned:", got)
			t.Fail()
		}
	}
}

func mockPageNote() *hlcmsg.Pagenote {
	return &hlcmsg.Pagenote{
		Uid:    1,
		Pageid: 1,
		Url:    "http://example.com",
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
