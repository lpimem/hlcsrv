package storage

import (
	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/hlcmsg"
)

var storage *SqliteStorage = nil

func InitStorage(path string) {
	if storage != nil {
		log.Warn("WARN:", path, "is not inited as storage as storage was already registered at ", storage.path)
		return
	}
	log.Debug("using database @", path)
	storage = NewSqliteStorage(path)
}

func QueryPagenoteByUser(uid uint32) []*hlcmsg.Pagenote {
	notes, err := storage.QueryPagenote(uid, 0)
	if err != nil {
		log.Alert("error QueryPagenoteByUser, uid:", uid, err)
		return []*hlcmsg.Pagenote{}
	}
	return notes[uid]
}

func QueryPagenoteByUrl(url string) PagenoteDict {
	pid := storage.QueryPageId(url)
	notes, err := storage.QueryPagenote(0, pid)
	if err != nil {
		log.Alert("error QueryPagenoteByUrl, url:", url, "pid", pid, err)
	}
	return notes
}

func QueryPagenote(uid uint32, pid uint32) *hlcmsg.Pagenote {
	notes, err := storage.QueryPagenote(uid, pid)
	if err != nil {
		log.Alert("error QueryPagenote, uid:", uid, "pid:", pid, err)
		return nil
	}
	return notes.GetPagenote(uid, pid)
}

func DeleteRangeMetas(idList []uint32) []uint32 {
	deleted := []uint32{}
	for _, id := range idList {
		err := storage.DeleteRangeMeta(id)
		if err != nil {
			log.Alert("error cannot delete RangeMeta", id, err)
		} else {
			deleted = append(deleted, id)
		}
	}
	return deleted
}

func QueryPageId(url string) uint32 {
	log.Trace("Querying id for page ", url)
	id := storage.QueryPageId(url)
	if id < 1 {
		id = storage.NewPage("", url)
		log.Trace("new page id", id, url)
	}
	if id < 1 {
		id = storage.QueryPageId(url)
	}
	log.Trace("page id for ", url, " is ", id)
	return id
}

func QueryPage(pid uint32) (string, error) {
	return storage.QueryPage(pid)
}

func SavePagenote(pn *hlcmsg.Pagenote) []error {
	// storage.SavePagenote()
	errs := []error{}
	log.Debug(len(pn.Highlights), " blocks to save")
	for _, hlt := range pn.Highlights {
		id, err := storage.NewRangeMeta(pn.Uid, pn.Pageid, hlt)
		if err != nil {
			log.Error("Error saving new range meta", err)
			errs = append(errs, err)
		} else {
			hlt.Id = id
		}
	}
	return errs
}
