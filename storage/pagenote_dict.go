package storage

import (
	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/hlcmsg"
)

// PagenoteDict is a map of user id and list of page notes
type PagenoteDict map[uint32][]*hlcmsg.Pagenote

func (d PagenoteDict) getOrCreatePagenoteList(uid uint32) (notes []*hlcmsg.Pagenote) {
	if _, ok := d[uid]; !ok {
		notes = []*hlcmsg.Pagenote{}
		d[uid] = notes
	} else {
		notes = d[uid]
	}
	return
}

// GetPagenoteList returns pagenotes associated with a user name
// return nil if not found
func (d PagenoteDict) GetPagenoteList(uid uint32) (notes []*hlcmsg.Pagenote) {
	if _, ok := d[uid]; ok {
		notes = d[uid]
	} else {
		notes = nil
	}
	return
}

// AddPagenote adds a pagenote to a user's list
func (d *PagenoteDict) AddPagenote(uid uint32, note *hlcmsg.Pagenote) {
	notes := d.getOrCreatePagenoteList(uid)
	notes = append(notes, note)
	(*d)[uid] = notes
}

// GetPagenote gets pagenote with pid from user uid's list
func (d *PagenoteDict) GetPagenote(uid uint32, pid uint32) *hlcmsg.Pagenote {
	notes := d.getOrCreatePagenoteList(uid)
	for _, n := range notes {
		if pid == 0 || n.Pageid == pid {
			return n
		}
	}
	log.Debug("no pagenote found for ", uid, pid)
	return nil
}

// NewPagenote creates a new pagenote with id pid and add to user uid's list
func (d *PagenoteDict) NewPagenote(uid uint32, pid uint32) *hlcmsg.Pagenote {
	note := &hlcmsg.Pagenote{
		Pageid:     pid,
		Uid:        uid,
		Highlights: []*hlcmsg.RangeMeta{},
	}
	d.AddPagenote(uid, note)
	return note
}
