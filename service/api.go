package service

import "net/http"

func index(w http.ResponseWriter, req *http.Request) {}

func getPagenote(w http.ResponseWriter, req *http.Request) {}

func savePagenote(w http.ResponseWriter, req *http.Request) {
	if !RequirePost(w, req) {
		return
	}
}

func deletePagenote(w http.ResponseWriter, req *http.Request) {
	if !RequirePost(w, req) {
		return
	}
}
