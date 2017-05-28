package controller

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/conf"
	"github.com/lpimem/hlcsrv/storage"
	"net/http"
	"strconv"
	"strings"
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
	if !requireAuth(w, req) {
		return
	}
	var (
		resp []byte
		err  error
	)
	users, err := storage.User.All(100, 0)
	if err != nil {
		log.Errorf("Error querying users %s", err)
		http.Error(w, "Server error", http.StatusBadGateway)
		return
	}
	if resp, err = json.Marshal(users); err != nil {
		log.Errorf("Error JSON-encoding users %s", err)
		http.Error(w, "Server error", http.StatusBadGateway)
		return
	}
	w.Write(resp)
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
	uid := req.PostFormValue("uid")
	uri := req.PostFormValue("uri")
	log.Debugf("Req body: %s", req.PostForm.Encode())
	log.Debugf("Grant: %s -> %s", uid, uri)
	for _, v := range []string{uid, uri} {
		if strings.TrimSpace(v) == "" {
			http.Error(w, "Missing required parameter", http.StatusBadRequest)
			return
		}
	}
	var (
		user  storage.UserID
		uid64 uint64
		err   error
	)
	if uid64, err = strconv.ParseUint(uid, 10, 32); err != nil {
		var msg string
		var debugMsg = fmt.Sprintf("Cannot parse user id from %s : %s", uid, err)
		if conf.IsDebug() {
			msg = debugMsg
		} else {
			msg = "Invalid parameter"
		}
		log.Debug(debugMsg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	user = (storage.UserID)(uid64)
	if err := storage.Permission.Grant(user, uri); err != nil {
		log.Errorf("Error granting permission for %d to %s: %s", user, uri, err)
		http.Error(w, "server error", http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// revoke all access from user
// POST
//     - uid : user id
func (*admin) RevokeAll(w http.ResponseWriter, req *http.Request) {
	// TODO
}
