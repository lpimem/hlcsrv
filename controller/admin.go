package controller

import (
	"encoding/json"
	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/storage"
	"net/http"
)

type admin struct{}

// Admin provide admin request handlers
var Admin admin

// GET a list of restricted URI in json array
func (*admin) Restrictions(w http.ResponseWriter, req *http.Request) {
	if !requireAuth(w, req) {
		return
	}
	restrictedURI, err := storage.Restriction.All()
	if err != nil {
		log.Errorf("Error querying restrictions: %s", err)
		http.Error(w, "Server error", http.StatusBadGateway)
		return
	}
	resp, err := json.Marshal(restrictedURI)
	if err != nil {
		log.Errorf("Error encoding JSON: %s", err)
		http.Error(w, "Server error", http.StatusBadGateway)
		return
	}
	w.Write(resp)
}

// POST: add a uri to to restriction list
// ARG: uri = encoded uri
func (*admin) AddRestriction(w http.ResponseWriter, req *http.Request) {
	if !requireAuth(w, req) || !requirePost(w, req) {
		return
	}
	uri := req.Form.Get("uri")
	err := storage.Restriction.Add(uri)
	if err != nil {
		log.Errorf("Cannot add restriction: %s", err)
		http.Error(w, "Cannot add restriction: "+err.Error(), http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (*admin) RemoveRestriction(w http.ResponseWriter, req *http.Request) {
	if !requireAuth(w, req) || !requirePost(w, req) {
		return
	}
	uri := req.Form.Get("uri")
	err := storage.Restriction.Remove(uri)
	if err != nil {
		log.Errorf("Cannot add restriction: %s", err)
		http.Error(w, "Cannot add restriction: "+err.Error(), http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// List users
func (*admin) Users(w http.ResponseWriter, req *http.Request) {

}

// List permissions
func (*admin) Permissions(w http.ResponseWriter, req *http.Request) {
	// TODO
}

// grant access for user to uri
// POST
//     - uid : user id
//     - uri : uri prefix
func (*admin) Grant(w http.ResponseWriter, req *http.Request) {
	if !requireAuth(w, req) || !requirePost(w, req) {
		return
	}
}

// revoke all access from user
// POST
//     - uid : user id
func (*admin) RevokeAll(w http.ResponseWriter, req *http.Request) {
	// TODO
}
