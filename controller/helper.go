package controller

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-playground/log"
	"github.com/golang/protobuf/proto"
	"github.com/lpimem/hlcsrv/auth"
	"github.com/lpimem/hlcsrv/hlcmsg"
)

func parseRemoveNotesRequest(r *http.Request) *hlcmsg.IdList {
	ids := &hlcmsg.IdList{}
	payload, err := readRequestPayload(r)
	if err != nil {
		return nil
	}
	if err = proto.Unmarshal(payload, ids); err != nil {
		log.Debug("Cannot parse IdList", err)
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
		log.Debug("error cannot extract uid from request", err)
		return nil, err
	}
	if pid, err = strconv.ParseUint(
		params.Get("pid"), 10, 32); err != nil {
		log.Debug("warn cannot extract pid from request", err)
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
		log.Debug("parseGetNotesRequest: cleaned url is empty")
		return nil, err
	}
	return pn, nil
}

func readRequestPayload(r *http.Request) ([]byte, error) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Debug("Cannot read request body", err)
	}
	return payload, err
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
		log.Warn("Error: cannot encode message ", resp, err)
		return false
	}

	encoder := base64.NewEncoder(base64.StdEncoding, w)
	defer encoder.Close()
	_, err = encoder.Write(buf)
	if err != nil {
		log.Warn("Error: cannot write message to response", err)
		return false
	}
	return true
}

func requirePost(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodPost {
		return true
	}
	http.Error(w, "bad request", http.StatusBadRequest)
	return false
}

func requireAuth(w http.ResponseWriter, r *http.Request) bool {
	var authorized = true
	if !auth.IsAuthenticated(r) {
		ctx := r.Context()
		reason := ctx.Value(auth.REASON)
		var errMsg string = "not authenticated"
		if reason != nil {
			errMsg = errMsg + ": " + reason.(string)
		}
		log.Warn(errMsg)
		http.Error(w, errMsg, http.StatusUnauthorized)
		authorized = false
	}
	return authorized
}