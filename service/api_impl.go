package service

import (
	"net/http"

	"io/ioutil"

	"github.com/golang/protobuf/proto"
	"github.com/lpimem/hlcsrv/hlcmsg"
	"github.com/lpimem/hlcsrv/storage"
	"github.com/lpimem/hlcsrv/util"
)

func newNotes(pn *hlcmsg.Pagenote) uint32 {
	return storage.SavePagenote(pn)
}

func getNotes(pn *hlcmsg.Pagenote) {
	storage.QueryPagenote(pn.Uid, pn.Url)
}

func parseNewNotesRequest(r *http.Request) *hlcmsg.Pagenote {
	pn := &hlcmsg.Pagenote{}
	payload, err := readRequestPayload(r)
	if err != nil {
		return nil
	}
	if err = proto.Unmarshal(payload, pn); err != nil {
		util.Log("Cannot parse Pagenote", err)
		return nil
	}
	patchPageId(pn)
	return pn
}

func parseRemoveNotesRequest(r *http.Request) []uint32 {
	ids := &hlcmsg.IdList{}
	payload, err := readRequestPayload(r)
	if err != nil {
		return nil
	}
	if err = proto.Unmarshal(payload, ids); err != nil {
		util.Log("Cannot parse IdList", err)
		return nil
	}
	return ids.Arr
}

func parseGetNotesRequest(r *http.Request) *hlcmsg.Pagenote {
	return parseNewNotesRequest(r)
}

func readRequestPayload(r *http.Request) ([]byte, error) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.Log("Cannot read request body", err)
	}
	return payload, err
}

func patchPageId(pn *hlcmsg.Pagenote) {
	if pn.Pageid < 1 {
		pn.Pageid = storage.QueryPageId(pn.Url)
	}
}
