package service

import (
	"net/http"

	"io/ioutil"

	"github.com/golang/protobuf/proto"
	"github.com/lpimem/hlcsrv/hlcmsg"
	"github.com/lpimem/hlcsrv/util"
)

func parseNewNotesRequest(r *http.Request) *hlcmsg.Pagenote {
	pn := &hlcmsg.Pagenote{}
	payload, err := readRequestPayload(r)
	if err != nil {
		return nil
	}
	if err = proto.Unmarshal(payload, pn); err != nil {
		util.Log("Cannot parse Pagenote", err)
		return nil
	}
	return pn
}

func parseRemoveNotesRequest(r *http.Request) []uint32 {
	return []uint32{}
}

func parseGetNotesRequest(r *http.Request) *hlcmsg.Pagenote {
	return parseNewNotesRequest(r)
}

func readRequestPayload(r *http.Request) ([]byte, error) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.Log("Cannot read request body", err)
	}
	return payload, err
}
