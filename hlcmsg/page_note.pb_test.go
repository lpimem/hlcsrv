package hlcmsg

import (
	"testing"

	"github.com/golang/protobuf/proto"
)

func TestPagenoteMarshal(t *testing.T) {

	pn := &Pagenote{
		Uid:    1,
		Pageid: 1,
		Url:    "http://example.com",
		Highlights: []*RangeMeta{
			&RangeMeta{
				Anchor:      "/",
				Start:       "/1",
				StartOffset: 0,
				End:         "/2",
				EndOffset:   3,
			},
		},
	}

	buf, err := proto.Marshal(pn)
	if err != nil {
		t.Error("marshal valid pagenote should success, err:", err)
		t.Fail()
		return
	}
	pnNew := &Pagenote{}
	err = proto.UnmarshalMerge(buf, pnNew)
	if err != nil {
		t.Error("error: cannot unmarshal pagenote: ", err)
		t.Fail()
		return
	}
	if pnNew.Uid != pn.Uid || pnNew.Pageid != pn.Pageid || pnNew.Url != pn.Url || len(pnNew.Highlights) != len(pn.Highlights) {
		t.Error("unmarshaled pagenote value changed")
		t.Fail()
	}

}
