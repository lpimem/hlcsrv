package storage

import (
	"github.com/lpimem/hlcsrv/hlcmsg"
	"github.com/lpimem/hlcsrv/util"
)

var storage *SqliteStorage = nil

func InitStorage(path string) {
	storage = NewSqliteStorage(path)
}

func QueryPagenoteByUser(uid uint32) []*hlcmsg.Pagenote {
	notes, err := storage.QueryPagenote(uid, 0)
	if err != nil {
		util.Log("error QueryPagenoteByUser, uid:", uid, err)
		return []*hlcmsg.Pagenote{}
	}
	return notes[uid]
}

func QueryPagenoteByUrl(url string) PagenoteDict {
	pid := storage.QueryPageId(url)
	notes, err := storage.QueryPagenote(0, pid)
	if err != nil {
		util.Log("error QueryPagenoteByUrl, url:", url, "pid", pid, err)
	}
	return notes
}

func QueryPagenote(uid uint32, url string) *hlcmsg.Pagenote {
	pid := storage.QueryPageId(url)
	if pid <= 0 {
		pid = storage.NewPage("unknown", url)
	}
	notes, err := storage.QueryPagenote(uid, pid)
	if err != nil {
		util.Log("error QueryPagenote, uid:", uid, "pid:", pid, err)
		return nil
	}
	return notes.GetPagenote(uid, pid)
}

func DeleteRangeMetas(idList []uint32) {
	for _, id := range idList {
		err := storage.DeleteRangeMeta(id)
		if err != nil {
			util.Log("error cannot delete RangeMeta", id, err)
		}
	}
}

func QueryPageId(url string) uint32 {
	id := storage.QueryPageId(url)
	if id < 0 {
		id = storage.NewPage("", url)
	}
	if id < 0 {
		id = storage.QueryPageId(url)
	}
	return id
}

func SavePagenote(pn *hlcmsg.Pagenote) uint32 {
	// storage.SavePagenote()
	for _, hlt := range pn.Highlights {
		storage.NewRangeMeta(pn.Uid, pn.Pageid, hlt)
	}
	return pn.Pageid
}
