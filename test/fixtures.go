package test

import (
	"github.com/lpimem/hlcsrv/hlcmsg"
)

func MockupMeta() *hlcmsg.RangeMeta {
	return &hlcmsg.RangeMeta{
		Id:          123,
		Anchor:      "#123/1/2/3",
		Start:       "#123/1/2/3/1",
		StartOffset: 0,
		End:         "#123/1/2/3/5",
		EndOffset:   19,
		Text:        "A Sample Text message",
	}
}
