package storage

import (
	"github.com/lpimem/hlcsrv/hlcmsg"
)

var storage *SqliteStorage = nil

func InitStorage(path string) {
	storage = NewSqliteStorage(path)
}

func QueryPageNotesByUser(uid string) []*hlcmsg.PageNotes {
	return []*hlcmsg.PageNotes{}
}

func QueryPageNotesByUrl(url string) []*hlcmsg.PageNotes {
	return []*hlcmsg.PageNotes{}
}

func QueryPageNote(uid uint32, url string) *hlcmsg.PageNotes {
	return nil
}

func (s *SqliteStorage) UpdatePageNotes(notes *hlcmsg.PageNotes) {

}
