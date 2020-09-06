package storage

import (
	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/hlcmsg"
)

var storage *SqliteStorage

// InitStorage initialize a default storage instance
func InitStorage(path string) {
	if storage != nil {
		log.Warn("WARN:", path, "is not inited as storage as storage was already registered at ", storage.path)
		return
	}
	log.Debug("using database @", path)
	storage = NewSqliteStorage(path)
}

// QueryPagenoteByUser select all pagenotes of a user
func QueryPagenoteByUser(uid uint32) []*hlcmsg.Pagenote {
	notes, err := storage.QueryPagenote(uid, 0)
	if err != nil {
		log.Alert("error QueryPagenoteByUser, uid:", uid, err)
		return []*hlcmsg.Pagenote{}
	}
	return notes[uid]
}

//QueryPagenoteByURI selects all pagenotes on a URI
func QueryPagenoteByURI(url string) PagenoteDict {
	pid := storage.QueryPageID(url)
	notes, err := storage.QueryPagenote(0, pid)
	if err != nil {
		log.Alert("error QueryPagenoteByURI, url:", url, "pid", pid, err)
	}
	return notes
}

// QueryPagenote selects pagenotes for a specified user an page
func QueryPagenote(uid uint32, pid uint32) *hlcmsg.Pagenote {
	notes, err := storage.QueryPagenote(uid, pid)
	if err != nil {
		log.Alert("error QueryPagenote, uid:", uid, "pid:", pid, err)
		return nil
	}
	return notes.GetPagenote(uid, pid)
}

// FilterRangeByUID returns a sublist of idList created by user.
func FilterRangeByUID(idList []uint32, user UserID) []uint32 {
	return storage.FilterRangeMeta(idList, user)
}

// DeleteRangeMetas delete list of meta from db
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

// QueryPageID returns the ID of page
func QueryPageID(url string) uint32 {
	defer log.WithTrace().Info("Querying id for page ", url)
	id := storage.QueryPageID(url)
	if id < 1 {
		id = storage.NewPage("", url)
		defer log.WithTrace().Info("new page id", id, url)
	}
	if id < 1 {
		id = storage.QueryPageID(url)
	}
	defer log.WithTrace().Info("page id for ", url, " is ", id)
	return id
}

// QueryPage selects the URI of a page id
func QueryPage(pid uint32) (string, error) {
	return storage.QueryPage(pid)
}

// SavePagenote insert or update pagenotes
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

// QueryUserPagenoteByURI select pagenotes on similar URIs for user
func QueryUserPagenoteByURI(uid UserID, uri string) (PagenoteDict, PagenoteAddon, error,
) {
	return storage.QueryPagenoteFuzzy(uid, uri)
}
