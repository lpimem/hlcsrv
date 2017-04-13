package controller

import "net/http"
import (
	"bytes"

	"fmt"

	"github.com/lpimem/hlcsrv/conf"
	"github.com/lpimem/hlcsrv/util"
)

func Index(w http.ResponseWriter, req *http.Request) {
	w.Write(bytes.NewBufferString("sorry, ").Bytes())
	http.NotFound(w, req)
}

func GetPagenote(w http.ResponseWriter, req *http.Request) {
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
	if !RequirePost(w, req) {
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
	if !RequirePost(w, req) {
		return
	}
	idList := parseRemoveNotesRequest(req)
	deleted := rmNotes(idList)
	writeRespMessage(w, nil, deleted)

}
