package service

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"encoding/base64"

	"github.com/golang/protobuf/proto"
	"github.com/lpimem/hlcsrv/hlcmsg"
	"github.com/lpimem/hlcsrv/storage"
	"github.com/lpimem/hlcsrv/util"
)

func newNotes(pn *hlcmsg.Pagenote) *hlcmsg.IdList {
	storage.SavePagenote(pn)
	return getPagenoteMetaIds(pn)
}

func getNotes(pn *hlcmsg.Pagenote) *hlcmsg.Pagenote {
	if pn.Uid > 0 {
		return storage.QueryPagenote(pn.Uid, pn.Pageid)
	}
	util.Log("error, empty uid not supported")
	return nil
}

func parseNewNotesRequest(r *http.Request) *hlcmsg.Pagenote {
	pn := &hlcmsg.Pagenote{}
	payload, err := readRequestPayload(r)
	if err != nil {
		util.Debug("error loading payload", err)
		return nil
	}
	if payload == nil || len(payload) == 0 {
		util.Debug("empty payload", err)
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
	params := r.URL.Query()
	uid, err := strconv.ParseUint(params.Get("uid"), 10, 32)
	if err != nil {
		util.Log("error cannot extract uid from request", err)
		return nil
	}
	uri := params.Get("uri")
	pn := &hlcmsg.Pagenote{}
	pn.Uid = uint32(uid)
	pn.Url = uri
	patchPageId(pn)
	return pn
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
		pn.Url = cleanUrl(pn.Url)
		pn.Pageid = storage.QueryPageId(pn.Url)
	}
}

func cleanUrl(urlstr string) string {
	u, err := url.Parse(urlstr)
	if err != nil {
		util.Log("Error parsing url", urlstr, err)
		return ""
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	// ignore fragment
	u.Fragment = ""
	return u.String()
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
