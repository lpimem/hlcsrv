package storage

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/lpimem/hlcsrv/hlcmsg"
)

func TestQueryNotesByUID(t *testing.T) {
	ResetTestDb()
	testQueryNotes(1, 0, "uid", t)
}

func TestQueryNotesByPid(t *testing.T) {
	ResetTestDb()
	testQueryNotes(0, 1, "pid", t)
}

func TestQueryNotesByUIDAndPid(t *testing.T) {
	ResetTestDb()
	testQueryNotes(1, 1, "uid and pid", t)
}

func testQueryNotes(uid, pid uint32, msg string, t *testing.T) (notes []*hlcmsg.Pagenote) {
	var err error
	noteDict, err := storage.QueryPagenote(uid, pid)
	if err != nil {
		t.Error("cannot query Pagenote", msg, err)
		t.Fail()
		return
	}
	if uid > 0 {
		notes = noteDict.GetPagenoteList(uid)
	} else {
		for _, uidNotes := range noteDict {
			notes = uidNotes
			break
		}
	}
	if notes == nil {
		t.Error("queried notes shouldn't be nil for uid", uid)
		t.Fail()
		return
	}
	if len(notes) < 1 {
		t.Error("queried pagenote list is empty", msg, notes)
		t.Fail()
		return
	}
	if len(notes[0].Highlights) < 1 {
		t.Error("queried note Highlights shouldn't be empty", msg, notes[0])
		t.Fail()
	}
	return
}

func TestNewRangeMeta(t *testing.T) {
	ResetTestDb()
	var err error
	metas := storage.QueryMetaList(1, 1)
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
	metas = storage.QueryMetaList(1, 1)
	if len(metas) < 2 {
		t.Error("should be able to query 2 metas, actually got ", metas)
		t.Fail()
		return
	}
}

func TestDeleteRangeMeta(t *testing.T) {
	ResetTestDb()
	metas := storage.QueryMetaList(1, 1)
	count := len(metas)
	if count < 1 {
		t.Error("should be able to query 1 meta")
		t.Fail()
		return
	}
	err := storage.DeleteRangeMeta(metas[0].Id)
	if err != nil {
		t.Error("cannot delete range meta ", metas[0].Id, err)
		t.Fail()
		return
	}
	metas = storage.QueryMetaList(1, 1)
	if len(metas) >= count {
		t.Error("rangemeta is not deleted ", metas[0].Id)
		t.Fail()
		return
	}
}

func TestQueryPageID(t *testing.T) {
	ResetTestDb()
	pid := storage.QueryPageID("http://example.com")
	if pid != 1 {
		t.Error("Should get 1 for page id, but got", pid)
		t.Fail()
	}
	pid = storage.QueryPageID("notexist.example.com")
	if pid != 0 {
		t.Error("Shouldn't get pid != 0 for unknown uri, but got:", pid)
		t.Fail()
	}
}

func TestNewPage(t *testing.T) {
	ResetTestDb()
	url := "http://new.example.com"
	pid := storage.NewPage("test", url)
	if pid < 1 {
		t.Error("created page id should be larger than 0")
		t.Fail()
	}
	queriedID := storage.QueryPageID(url)
	if queriedID != pid {
		t.Error("queried id does not match created id")
		t.Fail()
	}
}

func TestQueryUser(t *testing.T) {
	ResetTestDb()
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
	ResetTestDb()
	uname := "Alice"
	email := "alice@example.com"
	passwd := "unsafe"
	slt := "unsafe"
	uid := storage.NewUser(uname, email, passwd, slt)
	if uid < 1 {
		t.Error("created user id should be > 0", uid)
		t.Fail()
	}
	queriedID := storage.QueryUser(uname, passwd)
	if queriedID < 1 {
		t.Error("cannot query created user")
		t.Fail()
	}
	if queriedID != uid {
		t.Error("quried ID not matching created id", queriedID, uid)
		t.Fail()
	}
}

func init() {
	ResetTestDb()
}

func countPagenoteDictItems(d PagenoteDict) int {
	var sum int
	for _, ls := range d {
		sum += len(ls)
	}
	return sum
}

func TestQueryPagenoteFuzzy(t *testing.T) {
	/*
		func (s *SqliteStorage) QueryPagenoteFuzzy(
			uid uint32, uriPattern string) (PagenoteDict, PagenoteAddon, error)
	*/
	type testcase struct {
		// testcase name
		Name string
		// query user id
		UID uint32
		// query uri pattern
		URI string
		// expected result count
		Count int
		// should execute successfully without error
		Success bool
	}

	tcs := []*testcase{
		&testcase{"Valid Query", 1, "example", 1, true},
		&testcase{"NO UID", 0, "example", 1, true},
		&testcase{"NO URL", 1, "", 1, true},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			ResetTestDb()
			pd, addons, err := storage.QueryPagenoteFuzzy(tc.UID, tc.URI)
			nPdItem := countPagenoteDictItems(pd)
			if nPdItem != len(addons) {
				t.Error("count of pagenotes dismatch with count of addons")
				t.Fail()
				return
			}
			if tc.Success && err != nil {
				t.Error(err)
				t.Fail()
				return
			}
			var fail bool
			if nPdItem != tc.Count {
				t.Error("count of pagenotes dismatch")
				fail = true
			}
			if len(addons) != tc.Count {
				t.Error("count of addons dismatch")
				fail = true
			}
			if fail {
				t.Fail()
			}
		})
	}

}
