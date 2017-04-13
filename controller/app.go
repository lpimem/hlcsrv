package controller

import "net/http"
import (
	"fmt"

	"encoding/json"

	"github.com/lpimem/hlcsrv/conf"
	"github.com/lpimem/hlcsrv/util"
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
	rawToken = req.FormValue("google_token")
	if sessionInfo, err = doAuthenticateGoogleUser(req.Context(), rawToken); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	payload, err := json.Marshal(sessionInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	_, err = w.Write(payload)
	if err != nil {
		util.Log(err)
	}
}

func GetPagenote(w http.ResponseWriter, req *http.Request) {
	if !requireAuth(w, req) {
		return
	}
	pn, err := parseGetNotesRequest(req)
	if err != nil {
		util.Log("Cannot parse request, error:", err)
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
	if !requirePost(w, req) {
		return
	}
	if !requireAuth(w, req) {
		return
	}
	pn, err := parseNewNotesRequest(req)
	if err != nil {
		util.Error("cannot parse request, error: ", err)
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
		util.Error("savePagenote: ", err)
	}
	writeRespMessage(w, nil, idList)
}

func DeletePagenote(w http.ResponseWriter, req *http.Request) {
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
