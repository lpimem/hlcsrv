package storage

import (
	"github.com/lpimem/hlc/util"
	"github.com/lpimem/hlcsrv/hlcmsg"
)

var storage *SqliteStorage = nil

func InitStorage(path string) {
	storage = NewSqliteStorage(path)
}

func QueryPageNotesByUser(uid uint32) []*hlcmsg.PageNotes {
	notes, err := storage.QueryPageNotes(uid, 0)
	if err != nil {
		util.Log("error QueryPageNotesByUser, uid:", uid, err)
		return []*hlcmsg.PageNotes{}
	}
	return notes[uid]
}

func QueryPageNotesByUrl(url string) PageNoteDict {
	pid := storage.QueryPageId(url)
	notes, err := storage.QueryPageNotes(0, pid)
	if err != nil {
		util.Log("error QueryPageNotesByUrl, url:", url, "pid", pid, err)
	}
	return notes
}

func QueryPageNote(uid uint32, url string) *hlcmsg.PageNotes {
	pid := storage.QueryPageId(url)
	if pid <= 0 {
		pid = storage.NewPage("unknown", url)
	}
	notes, err := storage.QueryPageNotes(uid, pid)
	if err != nil {
		util.Log("error QueryPageNote, uid:", uid, "pid:", pid, err)
		return nil
	}
	return notes.GetPageNote(uid, pid)
}

func DeleteRangeMetas(idList []uint32) {
	for _, id := range idList {
		err := storage.DeleteRangeMeta(id)
		if err != nil {
			util.Log("error cannot delete RangeMeta", id, err)
		}
	}
}
