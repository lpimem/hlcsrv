package controller

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/lpimem/hlcsrv/auth"

	"github.com/go-playground/log"
	"github.com/golang/protobuf/proto"
	"github.com/lpimem/hlcsrv/hlcmsg"
	"github.com/lpimem/hlcsrv/storage"
)

func newNotes(pn *hlcmsg.Pagenote) (*hlcmsg.IdList, error) {
	errList := storage.SavePagenote(pn)
	var err error
	if len(errList) > 0 {
		log.Error("Errors saving pagenotes:")
		for _, e := range errList {
			log.Error("    ", e)
		}
		err = fmt.Errorf("%d errors happed saving %d pagenotes", len(errList), len(pn.Highlights))
	}
	return getPagenoteMetaIDs(pn), err
}

func getNotes(pn *hlcmsg.Pagenote) *hlcmsg.Pagenote {
	if pn != nil && pn.Uid > 0 {
		return storage.QueryPagenote(pn.Uid, pn.Pageid)
	}
	log.Warn("getNotes: empty uid not supported")
	return nil
}

func rmNotes(user storage.UserID, toRemove *hlcmsg.IdList) *hlcmsg.IdList {
	if toRemove != nil && len(toRemove.Arr) > 0 {
		idList := storage.FilterRangeByUID(toRemove.Arr, user)
		deleted := make([]uint32, 0)
		if len(idList) > 0 {
			deleted = storage.DeleteRangeMetas(toRemove.Arr)
		}
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
	defer log.WithTrace().Info("received request raw:", payload)
	pn := &hlcmsg.Pagenote{}
	if err = proto.Unmarshal(payload, pn); err != nil {
		log.Debug("Cannot parse Pagenote", err)
		return nil, err
	}
	currentUID := uint32(auth.RequestUID(r))
	if currentUID != pn.Uid {
		log.Warnf("User %d is trying to create notes as User %d", currentUID, pn.Uid)
		pn.Uid = currentUID
	}
	defer log.WithTrace().Info("parsed request:", pn.Pageid, pn.Uid, pn.Url, len(pn.Highlights))
	patchPageID(pn)
	return pn, nil
}

func patchPageID(pn *hlcmsg.Pagenote) error {
	if pn.Pageid < 1 {
		defer log.WithTrace().Info("Cleaing url:", pn.Url)
		var err error
		pn.Url, err = cleanURL(pn.Url)
		if err != nil {
			return err
		}
		defer log.WithTrace().Info("Cleaned url:", pn.Url)
		pn.Pageid = storage.QueryPageID(pn.Url)
	}
	return nil
}

func verifyPid(pid uint32) error {
	_, err := storage.QueryPage(pid)
	return err
}

func cleanURL(urlstr string) (string, error) {
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Warn("Error parsing url", urlstr, err)
		return "", err
	}
	defer log.WithTrace().Info("parsed url:", u.String())
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

func getPagenoteMetaIDs(pn *hlcmsg.Pagenote) *hlcmsg.IdList {
	arr := []uint32{}
	for _, m := range pn.Highlights {
		arr = append(arr, m.Id)
	}
	idl := &hlcmsg.IdList{
		Arr: arr,
	}
	return idl
}
