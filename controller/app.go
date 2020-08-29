package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/auth"
	"github.com/lpimem/hlcsrv/security"
	"github.com/lpimem/hlcsrv/conf"
	"github.com/lpimem/hlcsrv/hlccookie"
	"github.com/lpimem/hlcsrv/storage"
)

// Index renders the default page of the website.
func Index(w http.ResponseWriter, req *http.Request) {
	http.NotFound(w, req)
}

// AuthenticateGoogleUser handles post request to authenticate a google
// user. Expecting a newly generated google token in the request body.
// The token will be parsed and validated. See also:
// 1. https://developers.google.com/identity/sign-in/web/sign-in
// 2. https://developers.google.com/identity/sign-in/web/backend-auth#verify-the-integrity-of-the-id-token
// 3. https://github.com/coreos/go-oidc/blob/c3a2c79e8008bc1b1b0509ae6bf1483642c976f4/example/idtoken/app.go#L66
// 4. OAuth 2.0 Bearer Token Usage https://tools.ietf.org/html/rfc6750
// 5. OAuth 2.0 https://tools.ietf.org/html/rfc6749
func AuthenticateGoogleUser(w http.ResponseWriter, req *http.Request) {
	if !requirePost(w, req) {
		return
	}
	var (
		rawToken    string
		err         error
		sessionInfo *auth.SessionInfo
	)
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	rawToken = string(reqBody)
	if sessionInfo, err = auth.AuthenticateGoogleUser(req.Context(), rawToken); err != nil {
		log.Warn("Cannot authenticate google user ", err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	respJSON, err := json.Marshal(sessionInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	hlccookie.SetAuthCookies(w, sessionInfo.Sid, sessionInfo.Uid)
	_, err = w.Write(respJSON)
	if err != nil {
		log.Error(err)
	}
}

// Logout handles the logout request by removing the session.
func Logout(w http.ResponseWriter, req *http.Request) {
	if !requirePost(w, req) {
		log.Warn("Logout request should use POST method")
		return
	}
	if !requireAuth(w, req) {
		log.Warn("Logout request should be authenticated")
		return
	}
	sid := req.Context().Value(auth.SESSION_ID).(string)
	err := storage.DeleteSession(sid)
	if err != nil {
		log.Error(err)
	}
	conf.RedirectTo("/", "", w, req)
}

// GetPagenote handles get request to fetch notes for a user and a url
//
// Query parameters:
//     1. uid user identifier, must match the request session
//     2. pid [optional] page id for the url
//     3. url [optional] url string of the page
// Parameter 2 & 3 cannot be both empty.
//
// Response:
// 	Serialized hlcmsg.HlcResp message encoded in base64.
//
// Response Errors:
// 	http.StatusUnauthorized : client must be authenticated with a valid session token
// 	http.StatusBadRequest : client's request is in malformat
//
// See also:
//   1. hlc_resp.proto https://github.com/lpimem/hlcproto/blob/e7787d65aea33d1eb97b3f1f208394ee6a59f187/hlc_resp.proto
func GetPagenote(w http.ResponseWriter, req *http.Request) {
	defer  log.WithTrace().Info("GetPagenote...")
	security.EnableCORS(w)
	if (isHTTPOption(req)) {
		return
	}
	if !requireAuth(w, req) {
		log.Warn("GetPagenote: not authorized...")
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

// SavePagenote handles save pagenote post request.
// Expecting the body of the request to be a serialized hlcmsg/Pagenote message.
// Response:
//   Serialized @hlcmsg.HlcResp message encoded in base64.
func SavePagenote(w http.ResponseWriter, req *http.Request) {
	defer  log.WithTrace().Info("SavePagenote...")
	security.EnableCORS(w)
	if (isHTTPOption(req)) {
		return
	}
	if !requirePost(w, req) {
		log.Warn("SavePagenote: invalid http method...")
		return
	}
	if !requireAuth(w, req) {
		log.Warn("SavePagenote: not authorized...")
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

// DeletePagenote handles the post request to delete an array of notes.
func DeletePagenote(w http.ResponseWriter, req *http.Request) {
	defer  log.WithTrace().Info("DeletePagenote...")
	security.EnableCORS(w)
	if (isHTTPOption(req)) {
		return
	}
	if !requirePost(w, req) {
		log.Warn("DeletePagenote: invalid http method...")
		return
	}
	if !requireAuth(w, req) {
		log.Warn("SavePagenote: not authorized...")
		return
	}
	idList := parseRemoveNotesRequest(req)
	deleted := rmNotes(idList)
	writeRespMessage(w, nil, deleted)
}
