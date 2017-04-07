package service

import (
	"net/http/httptest"
	"testing"

	"net/http"

	"io/ioutil"

	"bytes"

	"github.com/golang/protobuf/proto"
	"github.com/lpimem/hlcsrv/hlcmsg"
	"github.com/lpimem/hlcsrv/storage"
)

func TestNewPagenote(t *testing.T) {
	const url = "/pagenote/new"
	var w *httptest.ResponseRecorder
	// post only
	get := httptest.NewRequest("GET", url, nil)
	w = httptest.NewRecorder()
	savePagenote(w, get)
	if w.Code != http.StatusBadRequest {
		t.Error("Expecting", http.StatusBadRequest, "got", w.Code, w.Body.String())
	}

	// post empty
	var post *http.Request
	post = httptest.NewRequest("POST", url, nil)
	w = httptest.NewRecorder()
	savePagenote(w, post)
	if w.Code == http.StatusOK {
		t.Error("null request body shouldn't get accepted")
		t.Fail()
	}
	// post normal
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
	post = httptest.NewRequest("POST", url, reader)
	w = httptest.NewRecorder()
	savePagenote(w, post)
	if w.Code != http.StatusOK {
		t.Error("valid new pn request failed", w.Body.String())
		t.Fail()
		return
	}
	respBody := w.Body.Bytes()
	pnResp := &hlcmsg.HlcResp{}
	err := proto.Unmarshal(respBody, pnResp)
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
	buf, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		t.Error("cannot read resp body", err)
		t.Fail()
		return
	}
	if len(buf) < 1 {
		t.Error("response body buf shouldn't be empty ")
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
