package service

import "net/http"
import "bytes"

func index(w http.ResponseWriter, req *http.Request) {
	w.Write(bytes.NewBufferString("sorry, ").Bytes())
	http.NotFound(w, req)
}

func getPagenote(w http.ResponseWriter, req *http.Request) {
	pn := parseGetNotesRequest(req)
	if pn == nil {
		http.Error(w, "invalid get pagenote request", http.StatusBadRequest)
		return
	}
	pn = getNotes(pn)
	if pn == nil {
		http.Error(w, "cannot get notes", http.StatusBadGateway)
		return
	}
	writeRespMessage(w, pn, nil)
}

func savePagenote(w http.ResponseWriter, req *http.Request) {
	if !RequirePost(w, req) {
		return
	}
	pn := parseNewNotesRequest(req)
	if pn == nil {
		http.Error(w, "cannot parse request", http.StatusBadRequest)
		return
	}
	idList := newNotes(pn)
	writeRespMessage(w, nil, idList)
}

func deletePagenote(w http.ResponseWriter, req *http.Request) {
	if !RequirePost(w, req) {
		return
	}
	idList := parseRemoveNotesRequest(req)
	deleted := rmNotes(idList)
	writeRespMessage(w, nil, deleted)

}
