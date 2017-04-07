package hlcmsg

import "testing"
import "github.com/golang/protobuf/proto"
import "fmt"

func TestResp(t *testing.T) {
	meta := mockupMeta()
	resp := &HlcResp{
		Code: HlcResp_SUC,
		Msg:  "success",
		PagenoteList: []*Pagenote{
			&Pagenote{
				Uid:    1,
				Url:    "example.com",
				Pageid: 1,
				Highlights: []*RangeMeta{
					meta,
				},
			},
		},
	}

	buf, err := proto.Marshal(resp)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	fmt.Printf("%d\r\n", len(buf))
	if len(buf) <= 1 {
		t.Error("buf is empty")
		t.Fail()
	}
	parsed := &HlcResp{}
	err = proto.Unmarshal(buf, parsed)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if len(parsed.PagenoteList) == 0 {
		t.Error("parsed pagenote is emtpy ")
		t.Fail()
	}
	if len(parsed.PagenoteList[0].Highlights) == 0 {
		t.Error("parsed highlights is emtpy ")
		t.Fail()
	}
}
