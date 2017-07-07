package controller

import (
	"bytes"
	"context"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"strconv"

	"github.com/go-playground/log"
	"github.com/golang/protobuf/proto"
	"github.com/lpimem/hlcsrv/auth"
	"github.com/lpimem/hlcsrv/hlcmsg"
	"github.com/lpimem/hlcsrv/storage"
)

const URLNewPagenote = "/pagenote/new"

const fakeUID = 10
const fakeSID = "fake_session_id"

func TestNewPagenotePostOnly(t *testing.T) {
	// post only
	var w *httptest.ResponseRecorder
	get := httptest.NewRequest("GET", URLNewPagenote, nil)
	get = fakeAuthenticateion(get)
	w = httptest.NewRecorder()
	SavePagenote(w, get)
	if w.Code != http.StatusBadRequest {
		t.Error("Expecting", http.StatusBadRequest, "got", w.Code, w.Body.String())
		t.Fail()
	}
}

func TestNewPagenoteNoEmptyReq(t *testing.T) {
	var w *httptest.ResponseRecorder

	// post empty
	var post *http.Request
	post = httptest.NewRequest("POST", URLNewPagenote, nil)
	post = fakeAuthenticateion(post)
	w = httptest.NewRecorder()
	SavePagenote(w, post)
	if w.Code == http.StatusOK {
		t.Error("null request body shouldn't get accepted")
		t.Fail()
	}
}

func TestNewPagenoteNormal(t *testing.T) {
	var w *httptest.ResponseRecorder
	reqPn := mockPageNote(1, 1, "http://example.com/index.html")
	buf, _ := proto.Marshal(reqPn)
	reader := bytes.NewReader(buf)
	post := httptest.NewRequest("POST", URLNewPagenote, reader)
	post = fakeAuthenticateion(post)
	w = httptest.NewRecorder()
	SavePagenote(w, post)
	if w.Code != http.StatusOK {
		t.Error("valid new pn request failed", w.Body.String())
		t.Fail()
		return
	}
	var err error
	log.Trace("encoded resp body:", w.Body.String())
	decoder := base64.NewDecoder(base64.StdEncoding, w.Body)
	respBody, err := ioutil.ReadAll(decoder)
	if err != nil {
		t.Error("cannot decode base64 encoded resp body", err)
		t.Fail()
		return
	}
	log.Trace("base64 decoded body:", respBody)
	pnResp := &hlcmsg.HlcResp{}
	err = proto.Unmarshal(respBody, pnResp)
	if err != nil {
		t.Error("cannot parse response body")
		t.Fail()
		return
	}
	if pnResp.Code != hlcmsg.HlcResp_SUC {
		t.Error("valid new pn request failed", pnResp.Msg)
		t.Fail()
		return
	}
	if pnResp.IdList == nil || len(pnResp.IdList.Arr) != 1 {
		t.Error("Failed to get ", 1, "created id")
		t.Fail()
		return
	}
	if pnResp.IdList.Arr[0] < 1 {
		t.Error("ID of created RangeMeta should be > 0")
		t.Fail()
		return
	}
}

func BenchmarkSavePagenote(b *testing.B) {
	var w *httptest.ResponseRecorder
	reqPn := mockPageNote(1, 1, "http://example.com/index.html")
	buf, _ := proto.Marshal(reqPn)
	//requests := []*http.Request{}
	//for n := 0; n < b.N; n++ {
	//	reader := bytes.NewReader(buf)
	//	post := httptest.NewRequest("POST", URLNewPagenote, reader)
	//	post = fakeAuthenticateion(post)
	//	requests = append(requests, post)
	//}
	for n := 0; n < b.N; n++ {
		reader := bytes.NewReader(buf)
		post := httptest.NewRequest("POST", URLNewPagenote, reader)
		post = fakeAuthenticateion(post)
		w = httptest.NewRecorder()
		SavePagenote(w, post)
	}
}

func BenchmarkSavePagenoteP(b *testing.B) {
	var w *httptest.ResponseRecorder
	reqPn := mockPageNote(1, 1, "http://example.com/index.html")
	buf, _ := proto.Marshal(reqPn)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			reader := bytes.NewReader(buf)
			post := httptest.NewRequest("POST", URLNewPagenote, reader)
			post = fakeAuthenticateion(post)
			w = httptest.NewRecorder()
			SavePagenote(w, post)
		}
	})
}

func BenchmarkGetPagenote(b *testing.B) {
	req := httptest.NewRequest("GET", "/pagenote?uid=1&url=example.com", nil)
	req = fakeAuthenticateion(req)
	tcs := []int{10, 100, 1000}
	for _, tc := range tcs {
		storage.ResetTestDb()
		storage.SeedDbForBench(tc)
		b.Run(strconv.FormatInt(int64(tc), 10), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				recorder := httptest.NewRecorder()
				GetPagenote(recorder, req)
			}
		})
	}
}

func BenchmarkGetPagenoteP(b *testing.B) {
	storage.ResetTestDb()
	storage.SeedDbForBench(100)
	req := httptest.NewRequest("GET", "/pagenote?uid=1&url=example.com", nil)
	req = fakeAuthenticateion(req)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			recorder := httptest.NewRecorder()
			GetPagenote(recorder, req)
		}
	})
}

func TestGetPageNote(t *testing.T) {
	req := httptest.NewRequest("GET", "/pagenote?uid=1&url=example.com", nil)
	req = fakeAuthenticateion(req)
	recorder := httptest.NewRecorder()
	GetPagenote(recorder, req)
	httpResp := recorder.Result()
	if httpResp.StatusCode != http.StatusOK {
		t.Error("response code should be ", http.StatusOK, "got", httpResp.StatusCode, recorder.Body.String())
		t.Fail()
		return
	}
	if recorder.Body == nil {
		t.Error("response body shouldn't be nil ")
		t.Fail()
		return
	}
	b64buf, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		t.Error("cannot read resp body", err)
		t.Fail()
		return
	}
	if len(b64buf) < 1 {
		t.Error("response body buf shouldn't be empty ")
		t.Fail()
		return
	}
	buf, err := ioutil.ReadAll(
		base64.NewDecoder(
			base64.StdEncoding, bytes.NewBuffer(b64buf)))
	if err != nil {
		t.Error("response body should be base 64 encoded", err)
		t.Fail()
		return
	}
	resp := &hlcmsg.HlcResp{}
	proto.Unmarshal(buf, resp)
	if resp.Code != hlcmsg.HlcResp_SUC {
		t.Error("response code should be ", hlcmsg.HlcResp_SUC, "got", resp.Code)
		t.Fail()
		return
	}
	if len(resp.PagenoteList) < 1 || len(resp.PagenoteList[0].Highlights) < 1 {
		t.Error("parsed response page list should contain 1 range meta")
		t.Fail()
		return
	}
}

func TestFakeAuthentication(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r = fakeAuthenticateion(r)
	if !auth.IsAuthenticated(r) {
		t.Error("fakeAuthenticateion should mark request authenticated.")
		t.Fail()
	}
}

func TestLogout(t *testing.T) {
	r := httptest.NewRequest("POST", "/logout", nil)
	fakeAuthenticateion(r)
	if !auth.IsAuthenticated(r) {
		t.Error("faked session not working")
		t.FailNow()
	}
	resp := httptest.NewRecorder()
	Logout(resp, r)
	if resp.Code != http.StatusOK {
		t.Errorf("response code is not OK: %d", resp.Code)
		t.Fail()
	}

	lastAccess, err := storage.QuerySession(fakeSID, fakeUID)
	if err != nil {
		t.Error("cannot query user session after logout: ", err)
		t.Fail()
	}
	if lastAccess != nil {
		t.Error("session still exists after logout: ", lastAccess)
		t.Fail()
	}
	if auth.IsAuthenticated(r) {
		t.Error("Logout failed: request is still authenticated")
		t.Fail()
	}
}

func fakeAuthenticateion(r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, auth.AUTHENTICATED, true)
	ctx = context.WithValue(ctx, auth.USER_ID, fakeUID)
	ctx = context.WithValue(ctx, auth.SESSION_ID, fakeSID)
	storage.UpdateSession(fakeSID, fakeUID)
	return r.WithContext(ctx)
}

func fakeAdminAuthenticateion(r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, auth.AUTHENTICATED, true)
	ctx = context.WithValue(ctx, auth.USER_ID, 1)
	ctx = context.WithValue(ctx, auth.SESSION_ID, "fake_session_id_for_1")
	return r.WithContext(ctx)
}

func init() {
	storage.ResetTestDb()
}
