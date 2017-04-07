package storage

import (
	"github.com/lpimem/hlcsrv/hlcmsg"
	"github.com/lpimem/hlcsrv/util"
)

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

func (d PagenoteDict) GetPagenoteList(uid uint32) (notes []*hlcmsg.Pagenote) {
	if _, ok := d[uid]; ok {
		notes = d[uid]
	} else {
		notes = nil
	}
	return
}

func (d *PagenoteDict) AddPagenote(uid uint32, note *hlcmsg.Pagenote) {
	notes := d.getOrCreatePagenoteList(uid)
	notes = append(notes, note)
	(*d)[uid] = notes
}

func (d *PagenoteDict) GetPagenote(uid uint32, pid uint32) *hlcmsg.Pagenote {
	notes := d.getOrCreatePagenoteList(uid)
	for _, n := range notes {
		if pid == 0 || n.Pageid == pid {
			return n
		}
	}
	util.Log("no pagenote found for", uid, pid)
	return nil
}

func (d *PagenoteDict) NewPagenote(uid uint32, pid uint32) *hlcmsg.Pagenote {
	note := &hlcmsg.Pagenote{
		Pageid:     pid,
		Uid:        uid,
		Highlights: []*hlcmsg.RangeMeta{},
	}
	d.AddPagenote(uid, note)
	return note
}
