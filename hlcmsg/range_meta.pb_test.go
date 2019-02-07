package hlcmsg

import "testing"
import "github.com/golang/protobuf/proto"

func TestRangeMeta(t *testing.T) {
	rmeta := mockupMeta()
	data, err := proto.Marshal(rmeta)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	newMeta := &RangeMeta{}
	err = proto.Unmarshal(data, newMeta)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	switch {
	case rmeta.Id != newMeta.Id:
		t.Error("id not match")
		t.Fail()
	case rmeta.Start != newMeta.Start:
		t.Error("Start not match")
		t.Fail()
	case rmeta.End != newMeta.End:
		t.Error("End not match")
		t.Fail()
	case rmeta.StartOffset != newMeta.StartOffset:
		t.Error("StartOffset not match")
		t.Fail()
	case rmeta.EndOffset != newMeta.EndOffset:
		t.Error("EndOffset not match")
		t.Fail()
	case rmeta.Anchor != newMeta.Anchor:
		t.Error("Anchor not match")
		t.Fail()
	case rmeta.Text != newMeta.Text:
		t.Error("Text not match")
		t.Fail()
	}
}

func mockupMeta() *RangeMeta {
	return &RangeMeta{
		Id:          123,
		Anchor:      "#123/1/2/3",
		Start:       "#123/1/2/3/1",
		StartOffset: 0,
		End:         "#123/1/2/3/5",
		EndOffset:   19,
		Text:        "A Sample Text message",
	}
}
