package controller

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-playground/log"
	"github.com/golang/protobuf/proto"
	"github.com/lpimem/hlcsrv/hlcmsg"
	"github.com/lpimem/hlcsrv/storage"
)

func newNotes(pn *hlcmsg.Pagenote) (*hlcmsg.IdList, error) {
	errList := storage.SavePagenote(pn)
	var err error = nil
	if len(errList) > 0 {
		log.Error("Errors saving pagenotes:")
		for _, e := range errList {
			log.Error("    ", e)
		}
		err = errors.New(fmt.Sprintf("%d errors happed saving %d pagenotes", len(errList), len(pn.Highlights)))
	}
	return getPagenoteMetaIds(pn), err
}

func getNotes(pn *hlcmsg.Pagenote) *hlcmsg.Pagenote {
	if pn != nil && pn.Uid > 0 {
		return storage.QueryPagenote(pn.Uid, pn.Pageid)
	}
	log.Warn("getNotes: empty uid not supported")
	return nil
}

func rmNotes(toRemove *hlcmsg.IdList) *hlcmsg.IdList {
	if toRemove != nil && len(toRemove.Arr) > 0 {
		deleted := storage.DeleteRangeMetas(toRemove.Arr)
		return &hlcmsg.IdList{
			Arr: deleted,
		}
	}
	return nil
}

func parseNewNotesRequest(r *http.Request) (*hlcmsg.Pagenote, error) {
	payload, err := readRequestPayload(r)
	if err != nil {
		return nil, err
	}
	if payload == nil || len(payload) == 0 {
		log.Debug("empty payload")
		return nil, errors.New("request payload is empty")
	}
	log.Trace("received request raw:", payload)
	pn := &hlcmsg.Pagenote{}
	if err = proto.Unmarshal(payload, pn); err != nil {
		log.Debug("Cannot parse Pagenote", err)
		return nil, err
	}
	log.Trace("parsed request:", pn.Pageid, pn.Uid, pn.Url, len(pn.Highlights))
	patchPageId(pn)
	return pn, nil
}

func patchPageId(pn *hlcmsg.Pagenote) error {
	if pn.Pageid < 1 {
		log.Trace("Cleaing url:", pn.Url)
		var err error
		pn.Url, err = cleanUrl(pn.Url)
		if err != nil {
			return err
		}
		log.Trace("Cleaned url:", pn.Url)
		pn.Pageid = storage.QueryPageId(pn.Url)
	}
	return nil
}

func verifyPid(pid uint32) error {
	_, err := storage.QueryPage(pid)
	return err
}

func cleanUrl(urlstr string) (string, error) {
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Warn("Error parsing url", urlstr, err)
		return "", err
	}
	log.Trace("parsed url:", u.String())
	if u.String() == "" {
		return "", errors.New("url shouldn't be empty")
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	// ignore fragment
	u.Fragment = ""
	return u.String(), nil
}

func getPagenoteMetaIds(pn *hlcmsg.Pagenote) *hlcmsg.IdList {
	arr := []uint32{}
	for _, m := range pn.Highlights {
		arr = append(arr, m.Id)
	}
	idl := &hlcmsg.IdList{
		Arr: arr,
	}
	return idl
}
