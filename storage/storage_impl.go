package storage

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"os"

	// ignore warning
	_ "github.com/mattn/go-sqlite3"

	"bytes"

	"github.com/lpimem/hlcsrv/hlcmsg"
	"github.com/lpimem/hlcsrv/util"
)

type SqliteStorage struct {
	path string
	DB   *sql.DB
}

func NewSqliteStorage(path string) *SqliteStorage {
	db, err := prepareSQLDb(path)
	if err != nil {
		panic(err)
	}
	return &SqliteStorage{path, db}
}

func (s *SqliteStorage) Close() {
	s.DB.Close()
}

func (s *SqliteStorage) QueryMetaList(uid uint32, pid uint32) []*hlcmsg.RangeMeta {
	dict, err := s.QueryPagenote(uid, pid)
	if err != nil {
		util.Log("Error:", err)
		return []*hlcmsg.RangeMeta{}
	} else {
		return dict.GetPagenote(uid, pid).GetHighlights()
	}
}

func (s *SqliteStorage) QueryPagenote(uid uint32, pid uint32) (PagenoteDict, error) {
	if uid == 0 && pid == 0 {
		return PagenoteDict{}, errors.New("uid and url cannot both be 0")
	}
	var queryBuilder bytes.Buffer
	queryBuilder.WriteString(`select id, anchor, start, startOffset, end, endOffset, page, author from hlc_range where 1=1 `)
	var parameters = []interface{}{}
	if uid > 0 {
		queryBuilder.WriteString(" and author = ?")
		parameters = append(parameters, uid)
	}
	if pid > 0 {
		queryBuilder.WriteString(" and page = ?")
		parameters = append(parameters, pid)
	}
	var query = queryBuilder.String()
	result := PagenoteDict{}
	err := util.QueryDb(s.DB, query, parameters, func(rowno int, rows *sql.Rows) error {
		var id, startOffset, endOffset, page, author uint32
		var anchor, start, end string
		err := rows.Scan(&id, &anchor, &start, &startOffset, &end, &endOffset, &page, &author)
		if err != nil {
			return err
		}
		note := result.GetPagenote(author, page)
		if note == nil {
			note = result.NewPagenote(author, page)
		}
		meta := &hlcmsg.RangeMeta{
			Id:          id,
			Anchor:      anchor,
			Start:       start,
			StartOffset: startOffset,
			End:         end,
			EndOffset:   endOffset,
			Text:        "",
		}
		note.Highlights = append(note.Highlights, meta)
		return nil
	})
	return result, err
}

func (s *SqliteStorage) NewRangeMeta(uid uint32, pid uint32, m *hlcmsg.RangeMeta) (uint32, error) {
	r, err := s.DB.Exec(`insert into hlc_range(anchor, start, startOffset, end, endOffset, text, page, author) values(?, ?, ?, ?, ?, ?, ?, ?)`, m.Anchor, m.Start, m.StartOffset, m.End, m.EndOffset, m.Text, pid, uid)
	if err != nil {
		return 0, err
	}
	lastId, err := r.LastInsertId()
	return uint32(lastId), err
}

func (s *SqliteStorage) DeleteRangeMeta(id uint32) error {
	if id < 1 {
		return errors.New("invalid range meta id, should be > 0")
	}
	_, err := s.DB.Exec(`delete from hlc_range where id=?`, id)
	return err
}

func (s *SqliteStorage) QueryPageId(url string) uint32 {
	var id uint32
	err := util.QueryDb(s.DB,
		"select id from hlc_page where url = ?",
		[]interface{}{url},
		func(rowno int, rows *sql.Rows) error {
			return rows.Scan(&id)
		})
	if err != nil {
		util.Log("ignored error: ", err)
	}
	return id
}

func (s *SqliteStorage) NewPage(title, url string) (id uint32) {
	rst, err := s.DB.Exec(
		"insert into hlc_page (title, url) values (?, ?)",
		title, url,
	)
	if err != nil {
		util.Log("ignored error:", err)
		return
	}
	lastId, err := rst.LastInsertId()
	if err != nil {
		util.Log("ignored error:", err)
		id = 0
	} else {
		id = uint32(lastId)
	}
	return
}

func (s *SqliteStorage) NewUser(name, email, password, slt string) (id uint32) {
	r, err := s.DB.Exec(
		"insert into hlc_user (name, email, password, _slt) values (?, ?, ?, ?)",
		name, email, password, slt,
	)
	if err != nil {
		util.Log("ignored error: ", err)
		return
	}
	lastId, err := r.LastInsertId()
	if err != nil {
		util.Log("ignored error: ", err)
		return
	}
	id = uint32(lastId)
	return
}

func (s *SqliteStorage) QueryUser(handle, password string) (id uint32) {
	const active = 1
	query := `select id from hlc_user where _status = ? and ((name=? and password=?) or (email=? and password=?)) `
	err := util.QueryDb(s.DB, query,
		[]interface{}{active, handle, password, handle, password},
		func(idx int, rows *sql.Rows) error {
			return rows.Scan(&id)
		})
	if err != nil {
		util.Log("error querying user id :", err)
	}
	return
}

func initDb(db *sql.DB) error {
	fpath := util.GetAbsRunDirPath() + "/db/tables.sql"
	createTables, err := ioutil.ReadFile(fpath)
	if err != nil {
		util.Log("Current dir: ", util.GetAbsRunDirPath())
		util.Log("file path:", fpath)
		return err
	}
	_, err = db.Exec(string(createTables))
	return err
}

func prepareSQLDb(path string) (*sql.DB, error) {
	var isNew = false
	if _, err := os.Stat(path); os.IsNotExist(err) {
		isNew = true
	}
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if isNew {
		if err := initDb(db); err != nil {
			return nil, err
		}
	}
	return db, nil
}
