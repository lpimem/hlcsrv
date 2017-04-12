package service

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"errors"

	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/lpimem/hlcsrv/hlcmsg"
	"github.com/lpimem/hlcsrv/storage"
	"github.com/lpimem/hlcsrv/util"
)

func newNotes(pn *hlcmsg.Pagenote) (*hlcmsg.IdList, error) {
	errList := storage.SavePagenote(pn)
	var err error = nil
	if len(errList) > 0 {
		util.Error("Errors saving pagenotes:")
		for _, e := range errList {
			util.Error("    ", e)
		}
		err = errors.New(fmt.Sprintf("%d errors happed saving %d pagenotes", len(errList), len(pn.Highlights)))
	}
	return getPagenoteMetaIds(pn), err
}

func getNotes(pn *hlcmsg.Pagenote) *hlcmsg.Pagenote {
	if pn != nil && pn.Uid > 0 {
		return storage.QueryPagenote(pn.Uid, pn.Pageid)
	}
	util.Log("error, empty uid not supported")
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
		util.Debug("error loading payload", err)
		return nil, err
	}
	if payload == nil || len(payload) == 0 {
		util.Debug("empty payload")
		return nil, errors.New("request payload is empty")
	}
	util.Debug("received request raw:", payload)
	pn := &hlcmsg.Pagenote{}
	if err = proto.Unmarshal(payload, pn); err != nil {
		util.Log("Cannot parse Pagenote", err)
		return nil, err
	}
	util.Debug("parsed request:", pn.Pageid, pn.Uid, pn.Url, len(pn.Highlights))
	patchPageId(pn)
	return pn, nil
}

func parseRemoveNotesRequest(r *http.Request) *hlcmsg.IdList {
	ids := &hlcmsg.IdList{}
	payload, err := readRequestPayload(r)
	if err != nil {
		return nil
	}
	if err = proto.Unmarshal(payload, ids); err != nil {
		util.Log("Cannot parse IdList", err)
		return nil
	}
	return ids
}

func parseGetNotesRequest(r *http.Request) (*hlcmsg.Pagenote, error) {
	var (
		uid, pid uint64
		err      error
	)
	params := r.URL.Query()
	if uid, err = strconv.ParseUint(
		params.Get("uid"), 10, 32); err != nil {
		util.Log("error cannot extract uid from request", err)
		return nil, err
	}
	if pid, err = strconv.ParseUint(
		params.Get("pid"), 10, 32); err != nil {
		util.Log("warn cannot extract pid from request", err)
		pid = 0
	}
	pn := &hlcmsg.Pagenote{}
	pn.Uid = uint32(uid)
	if pid <= 0 {
		pn.Url = params.Get("url")
		err = patchPageId(pn)
	} else {
		if err = verifyPid(uint32(pid)); err == nil {
			pn.Pageid = uint32(pid)
		}
	}
	if err != nil {
		util.Error("parseGetNotesRequest: cleaned url is empty")
		return nil, err
	}
	return pn, nil
}

func readRequestPayload(r *http.Request) ([]byte, error) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.Log("Cannot read request body", err)
	}
	return payload, err
}

func patchPageId(pn *hlcmsg.Pagenote) error {
	if pn.Pageid < 1 {
		util.Debug("Cleaing url:", pn.Url)
		var err error
		pn.Url, err = cleanUrl(pn.Url)
		if err != nil {
			return err
		}
		util.Debug("Cleaned url:", pn.Url)
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
		util.Log("Error parsing url", urlstr, err)
		return "", err
	}
	util.Log("parsed url:", u.String())
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

func writeRespMessage(w http.ResponseWriter, pn *hlcmsg.Pagenote, idlist *hlcmsg.IdList) bool {
	resp := &hlcmsg.HlcResp{
		Code:         hlcmsg.HlcResp_SUC,
		Msg:          "sucess",
		PagenoteList: []*hlcmsg.Pagenote{},
		IdList:       idlist,
	}
	if pn != nil {
		resp.PagenoteList = append(resp.PagenoteList, pn)
	}
	buf, err := proto.Marshal(resp)
	if err != nil {
		util.Log("Error: cannot encode message ", resp, err)
		return false
	}

	encoder := base64.NewEncoder(base64.StdEncoding, w)
	defer encoder.Close()
	_, err = encoder.Write(buf)
	if err != nil {
		util.Log("Error: cannot write message to response", err)
		return false
	}
	//if _, err = w.Write(buf); err != nil {
	//	util.Log("Error: cannot write message to response", err)
	//	return false
	//}
	return true
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
