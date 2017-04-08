package service

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/lpimem/hlcsrv/hlcmsg"
	"github.com/lpimem/hlcsrv/storage"
	"github.com/lpimem/hlcsrv/util"
)

const URL_NEW_PAGENOTE = "/pagenote/new"

func TestNewPagenotePostOnly(t *testing.T) {
	// post only
	var w *httptest.ResponseRecorder
	get := httptest.NewRequest("GET", URL_NEW_PAGENOTE, nil)
	w = httptest.NewRecorder()
	savePagenote(w, get)
	if w.Code != http.StatusBadRequest {
		t.Error("Expecting", http.StatusBadRequest, "got", w.Code, w.Body.String())
		t.Fail()
	}
}

func TestNewPagenoteNoEmptyReq(t *testing.T) {
	var w *httptest.ResponseRecorder

	// post empty
	var post *http.Request
	post = httptest.NewRequest("POST", URL_NEW_PAGENOTE, nil)
	w = httptest.NewRecorder()
	savePagenote(w, post)
	if w.Code == http.StatusOK {
		t.Error("null request body shouldn't get accepted")
		t.Fail()
	}
}

func TestNewPagenoteNormal(t *testing.T) {
	var w *httptest.ResponseRecorder
	reqPn := &hlcmsg.Pagenote{
		Uid:    1,
		Pageid: 1,
		Url:    "http://example.com",
		Highlights: []*hlcmsg.RangeMeta{
			&hlcmsg.RangeMeta{
				Anchor:      "/",
				Start:       "/1",
				StartOffset: 0,
				End:         "/2",
				EndOffset:   3,
			},
		},
	}
	buf, _ := proto.Marshal(reqPn)
	reader := bytes.NewReader(buf)
	post := httptest.NewRequest("POST", URL_NEW_PAGENOTE, reader)
	w = httptest.NewRecorder()
	savePagenote(w, post)
	if w.Code != http.StatusOK {
		t.Error("valid new pn request failed", w.Body.String())
		t.Fail()
		return
	}
	var err error
	//respBody := w.Body.Bytes()
	util.Log("encoded resp body:", w.Body.String())
	decoder := base64.NewDecoder(base64.StdEncoding, w.Body)
	respBody, err := ioutil.ReadAll(decoder)
	if err != nil {
		t.Error("cannot decode base64 encoded resp body", err)
		t.Fail()
		return
	}
	util.Log("base64 decoded body:", respBody)
	pnResp := &hlcmsg.HlcResp{}
	err = proto.Unmarshal(respBody, pnResp)
	if err != nil {
		t.Error("cannot parse response body")
		t.Fail()
	}
	if pnResp.Code != hlcmsg.HlcResp_SUC {
		t.Error("valid new pn request failed", pnResp.Msg)
		t.Fail()
	}
	if pnResp.IdList == nil || len(pnResp.IdList.Arr) != 1 {
		t.Error("Failed to get ", 1, "created id")
		t.Fail()
	}
	if pnResp.IdList.Arr[0] < 1 {
		t.Error("ID of created RangeMeta should be > 0")
		t.Fail()
	}
}

func TestGetPageNote(t *testing.T) {
	req := httptest.NewRequest("GET", "/pagenote?uid=1", nil)
	recorder := httptest.NewRecorder()
	getPagenote(recorder, req)
	httpResp := recorder.Result()
	if httpResp.StatusCode != http.StatusOK {
		t.Error("response code should be ", http.StatusOK, "got", httpResp.StatusCode, httpResp.Status)
		t.Fail()
	}
	if recorder.Body == nil {
		t.Error("response body shouldn't be nil ")
		t.Fail()
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
	}
	buf, err := ioutil.ReadAll(
		base64.NewDecoder(
			base64.StdEncoding, bytes.NewBuffer(b64buf)))
	if err != nil {
		t.Error("response body should be base 64 encoded", err)
		t.Fail()
	}
	resp := &hlcmsg.HlcResp{}
	proto.Unmarshal(buf, resp)
	if resp.Code != hlcmsg.HlcResp_SUC {
		t.Error("response code should be ", hlcmsg.HlcResp_SUC, "got", resp.Code)
		t.Fail()
	}
	if len(resp.PagenoteList) < 1 || len(resp.PagenoteList[0].Highlights) < 1 {
		t.Error("parsed response page list should contain 1 range meta")
		t.Fail()
	}
}

func init() {
	storage.ResetTestDb()
}
