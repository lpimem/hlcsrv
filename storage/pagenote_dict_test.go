package storage

import (
	"testing"

	"github.com/lpimem/hlcsrv/hlcmsg"
)

func TestGetPagenoteList(t *testing.T) {
	ins := mock()
	notes := ins.GetPagenoteList(1)
	if len(notes) < 1 {
		t.Error("cannot query existing note by uid")
		t.Fail()
	}
	if notes[0].Url != "http://example.com" {
		t.Error("url dismatch", notes[0].Url)
		t.Fail()
	}
	if len(notes[0].Highlights) < 1 {
		t.Error("highlight inside pagenote is not retrived ")
		t.Fail()
	}
	if notes[0].Highlights[0].End != "/2/1" {
		t.Error("highlight inside pagenote is not retrived ")
		t.Fail()
	}
}

func TestNewPagenoteList(t *testing.T) {
	ins := mock()

	const newUID = uint32(2)
	const pid = 2

	nilNote := ins.GetPagenoteList(newUID)

	if nilNote != nil {
		t.Error("expecting nil but got ", nilNote)
		t.Fail()
	}

	ins.AddPagenote(newUID, &hlcmsg.Pagenote{
		Uid:        newUID,
		Pageid:     pid,
		Highlights: []*hlcmsg.RangeMeta{},
	})

	query := ins.GetPagenoteList(newUID)

	if query == nil {
		t.Error("query shouldn't be nil")
		t.Fail()
	}

	if len(query) < 1 {
		t.Error("query shouldn't be empty")
		t.Fail()
	}

	query[0].Highlights = append(query[0].Highlights, &hlcmsg.RangeMeta{
		Id:          1,
		Anchor:      "/1",
		Start:       "/12",
		StartOffset: 0,
		End:         "/22/1",
		EndOffset:   10,
	})

	if len(query[0].Highlights) < 1 || query[0].Highlights[0].End != "/22/1" {
		t.Error("4. query not reflecting new note change ")
		t.Fail()
	}

	query = ins.GetPagenoteList(newUID)

	if query == nil {
		t.Error("2 query shouldn't be nil")
		t.Fail()
	}

	if len(query) < 1 {
		t.Error("2 query shouldn't be empty")
		t.Fail()
	}

	if len(query[0].Highlights) < 1 || query[0].Highlights[0].End != "/22/1" {
		t.Error("re-query not reflecting new note change ")
		t.Fail()
	}

}

func mock() *PagenoteDict {
	return &PagenoteDict{
		1: []*hlcmsg.Pagenote{
			&hlcmsg.Pagenote{
				Uid:    1,
				Pageid: 2,
				Url:    "http://example.com",
				Highlights: []*hlcmsg.RangeMeta{
					&hlcmsg.RangeMeta{
						Id:          1,
						Anchor:      "/",
						Start:       "/1",
						StartOffset: 0,
						End:         "/2/1",
						EndOffset:   10,
					},
				},
			},
		},
	}
}
