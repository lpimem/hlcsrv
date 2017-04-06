package service

import "net/http"

func index(w http.ResponseWriter, req *http.Request) {
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
	writeRespMessage(w, pn)
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
	writeRespMessage(w, idList)
}

func deletePagenote(w http.ResponseWriter, req *http.Request) {
	if !RequirePost(w, req) {
		return
	}
	http.NotFound(w, req)
}
