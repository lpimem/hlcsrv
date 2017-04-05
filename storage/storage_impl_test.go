package storage

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/lpimem/hlcsrv/hlcmsg"
	"github.com/lpimem/hlcsrv/util"
)

func TestNewRangeMeta(t *testing.T) {
	resetDb()
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
	newMeta.Id, err = storage.NewRangeMeta(1, 1, newMeta)
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

func TestQueryPageId(t *testing.T) {
	resetDb()
	pid := storage.QueryPageId("example.com")
	if pid != 1 {
		t.Error("Should get 1 for page id, but got", pid)
		t.Fail()
	}
	pid = storage.QueryPageId("notexist.example.com")
	if pid != 0 {
		t.Error("Shouldn't get pid != 0 for unknown uri, but got:", pid)
		t.Fail()
	}
}

func TestNewPage(t *testing.T) {
	resetDb()
	url := "http://new.example.com"
	pid := storage.NewPage("test", url)
	if pid < 1 {
		t.Error("created page id should be larger than 0")
		t.Fail()
	}
	queriedId := storage.QueryPageId(url)
	if queriedId != pid {
		t.Error("queried id does not match created id")
		t.Fail()
	}
}

func TestQueryUser(t *testing.T) {
	resetDb()
	uname := "Bob"
	uemail := "bob@example.com"
	passwd := "unsafe"
	uid := storage.QueryUser(uname, passwd)
	if uid < 1 {
		t.Error("Cannot get user id using uname + passwd", uid)
		t.Fail()
	}
	uid = storage.QueryUser(uemail, passwd)
	if uid < 1 {
		t.Error("Cannot get user id using email + passwd", uid)
		t.Fail()
	}
}

func TestNewUser(t *testing.T) {
	resetDb()
	uname := "Alice"
	email := "alice@example.com"
	passwd := "unsafe"
	slt := "unsafe"
	uid := storage.NewUser(uname, email, passwd, slt)
	if uid < 1 {
		t.Error("created user id should be > 0", uid)
		t.Fail()
	}
	queriedId := storage.QueryUser(uname, passwd)
	if queriedId < 1 {
		t.Error("cannot query created user")
		t.Fail()
	}
	if queriedId != uid {
		t.Error("quried ID not matching created id", queriedId, uid)
		t.Fail()
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

func resetDb() {
	CleanDb()
	SeedTestDb()
}

func init() {
	InitTestDb()
}
