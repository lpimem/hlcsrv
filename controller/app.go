package controller

import "net/http"
import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/conf"
	"github.com/lpimem/hlcsrv/hlccookie"
)

func Index(w http.ResponseWriter, req *http.Request) {
	http.NotFound(w, req)
}

func AuthenticateGoogleUser(w http.ResponseWriter, req *http.Request) {
	if !requirePost(w, req) {
		return
	}
	var (
		rawToken    string
		err         error
		sessionInfo *SessionInfo
	)
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	rawToken = string(reqBody)
	if sessionInfo, err = doAuthenticateGoogleUser(req.Context(), rawToken); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	hlccookie.SetAuthCookies(w, sessionInfo.Sid, sessionInfo.Uid)
	respJson, err := json.Marshal(sessionInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	_, err = w.Write(respJson)
	if err != nil {
		log.Error(err)
	}
}

func GetPagenote(w http.ResponseWriter, req *http.Request) {
	defer log.Trace("GetPagenote...")
	if !requireAuth(w, req) {
		log.Warn("Not authorized...")
		return
	}
	pn, err := parseGetNotesRequest(req)
	if err != nil {
		log.Warn("Cannot parse request, error:", err)
		var errMsg string
		if conf.IsDebug() {
			errMsg = fmt.Sprintln("Cannot parse request, error:", err)
		} else {
			errMsg = "invalid get pagenote request"
		}
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	pn = getNotes(pn)
	writeRespMessage(w, pn, nil)
}

func SavePagenote(w http.ResponseWriter, req *http.Request) {
	defer log.Trace("SavePagenote...")
	if !requirePost(w, req) {
		return
	}
	if !requireAuth(w, req) {
		return
	}
	pn, err := parseNewNotesRequest(req)
	if err != nil {
		log.Warn("cannot parse request, error: ", err)
		var errMsg string
		if conf.IsDebug() {
			errMsg = fmt.Sprintln("cannot parse request, error: ", err)
		} else {
			errMsg = "cannot parse request"
		}
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	idList, err := newNotes(pn)
	if err != nil {
		log.Error("savePagenote: ", err)
	}
	writeRespMessage(w, nil, idList)
}

func DeletePagenote(w http.ResponseWriter, req *http.Request) {
	defer log.Trace("DeletePagenote...")
	if !requirePost(w, req) {
		return
	}
	if !requireAuth(w, req) {
		return
	}
	idList := parseRemoveNotesRequest(req)
	deleted := rmNotes(idList)
	writeRespMessage(w, nil, deleted)

}
