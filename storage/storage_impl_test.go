package storage

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/lpimem/hlcsrv/hlcmsg"
	"github.com/lpimem/hlcsrv/util"
)

func TestInsertRangeMeta(t *testing.T) {
	InitTestDb()
	metas, err := storage.QueryMetas(1, 1)
	if err != nil {
		t.Error("cannot query meta", err)
		t.Fail()
		return
	}
	if len(metas) < 1 {
		t.Error("should be able to query 1 meta")
		t.Fail()
		return
	}
	meta := metas[0]
	newMsg := proto.Clone(meta)
	newMeta, ok := newMsg.(*hlcmsg.RangeMeta)
	if !ok {
		t.Error("cannot convert cloned message to RangeMeta")
		t.Fail()
		return
	}
	newMeta.Id = 0
	newMeta.Id, err = storage.InsertRangeMeta(1, 1, newMeta)
	if err != nil {
		t.Error("cannot insert range meta: ", err)
		t.Fail()
		return
	}
	metas, err = storage.QueryMetas(1, 1)
	if len(metas) < 2 {
		t.Error("should be able to query 2 metas, actually got ", metas, err)
		t.Fail()
		return
	}
}

func InitTestDb() {
	InitStorage(util.GetAbsRunDirPath() + "/db/test.db")
	err := CleanDb()
	if err != nil {
		panic(err)
	}
	err = SeedTestDb()
	if err != nil {
		panic(err)
	}
}

func CleanDb() error {
	_, err := storage.DB.Exec(`
	delete from hlc_range;
	delete from hlc_user;
	delete from hlc_page;
	delete from hlc_comments;
	`)
	return err
}

func SeedTestDb() error {
	_, err := storage.DB.Exec(`
	insert into hlc_user(id, name, email, password, _slt) 
		values (1, "Bob", "bob@example.com", "unsafe", "unsafe");

	insert into hlc_page(id, title, url)
		values (1, "example", "example.com");

	insert into hlc_range(id, anchor, start, startOffset, end, endOffset, text, page, author)
		values (1, "#c", "#c/1", 0, "#c/12", 32, "This is the selected text", 1, 1);
	`)
	return err
}
