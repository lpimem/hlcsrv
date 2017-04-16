package storage

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/hlcmsg"
	"github.com/lpimem/hlcsrv/util"
	_ "github.com/mattn/go-sqlite3"
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

/**Upsert performs `update on insert duplicated` query.
Note, this method uses unchecked string format to build the query.
You should never use user's input as the value for table, fields, and keys.
*/
func (s *SqliteStorage) Upsert(
	table string,
	fields []string,
	keys []string,
	parameters []interface{},
	keyValues []interface{},
) (r sql.Result, err error) {
	if len(fields) != len(parameters) {
		return nil, errors.New("Numbers of fields and parameters do not match.")
	}
	if len(keys) != len(keyValues) {
		return nil, errors.New("Numbers of keys and keyValues do not match.")
	}
	if len(keys) < 1 {
		return nil, errors.New("Numbers of keys cannot be 0")
	}

	var (
		snippetUpdateFields string
		snippetUpdateCond   string
		snippetInsertFields string
		snippetInsertValues string
		query               string
		queryParameters     []interface{}
	)
	const template = `
		UPDATE %s
		SET %s
		WHERE %s;
		INSERT INTO %s (%s)
		SELECT %s
		WHERE (Select Changes() = 0);`

	snippetInsertFields = strings.Join(fields, ",")
	snippetInsertValues = strings.Repeat("?,", len(fields)-1) + "?"

	var updateFieldsBuilder bytes.Buffer
	var fieldsNumber = len(fields) - 1
	for i, f := range fields {
		updateFieldsBuilder.WriteString(f)
		updateFieldsBuilder.WriteString(" = ?")
		if i < fieldsNumber {
			updateFieldsBuilder.WriteString((","))
		}
	}
	snippetUpdateFields = updateFieldsBuilder.String()

	var updateCondBuilder bytes.Buffer
	var keysNumber = len(keys) - 1
	for i, k := range keys {
		updateCondBuilder.WriteString(k)
		updateCondBuilder.WriteString(" = ?")
		if i < keysNumber {
			updateCondBuilder.WriteString((","))
		}
	}
	snippetUpdateCond = updateCondBuilder.String()

	query = fmt.Sprintf(template, table, snippetUpdateFields, snippetUpdateCond,
		table, snippetInsertFields, snippetInsertValues)
	log.Debug(query)

	queryParameters = append(parameters, keyValues...)
	queryParameters = append(queryParameters, parameters...)

	r, err = s.DB.Exec(query, queryParameters...)
	return
}

func (s *SqliteStorage) QueryMetaList(uid uint32, pid uint32) []*hlcmsg.RangeMeta {
	dict, err := s.QueryPagenote(uid, pid)
	if err != nil {
		log.Warn("Error:", err)
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
	queryBuilder.WriteString(
		`select id, anchor, start, startOffset, end, endOffset, page, author, option from hlc_range where 1=1 `)
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
		var anchor, start, end, option string
		err := rows.Scan(&id, &anchor, &start, &startOffset, &end, &endOffset, &page, &author, &option)
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
			Option:      option,
		}
		note.Highlights = append(note.Highlights, meta)
		return nil
	})
	return result, err
}

func (s *SqliteStorage) NewRangeMeta(uid uint32, pid uint32, m *hlcmsg.RangeMeta) (uint32, error) {
	r, err := s.DB.Exec(
		`insert into hlc_range(anchor, start, startOffset, end, endOffset, text, page, author, option)
values (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		m.Anchor, m.Start, m.StartOffset, m.End, m.EndOffset, m.Text, pid, uid, m.Option)
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
		log.Warn("ignored error: ", err)
	}
	return id
}

func (s *SqliteStorage) QueryPage(pid uint32) (string, error) {
	var url string
	err := util.QueryDb(s.DB,
		"select url from hlc_page where id=?",
		[]interface{}{pid},
		func(idx int, rows *sql.Rows) error {
			return rows.Scan(&url)
		})
	return url, err

}

func (s *SqliteStorage) NewPage(title, url string) (id uint32) {
	rst, err := s.DB.Exec(
		"insert into hlc_page (title, url) values (?, ?)",
		title, url,
	)
	if err != nil {
		log.Warn("ignored error:", err)
		return
	}
	lastId, err := rst.LastInsertId()
	if err != nil {
		log.Warn("ignored error:", err)
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
		log.Warn("ignored error: ", err)
		return
	}
	lastId, err := r.LastInsertId()
	if err != nil {
		log.Warn("ignored error: ", err)
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
		log.Warn("error querying user id :", err)
	}
	return
}

func initDb(db *sql.DB) error {
	fpath := util.GetAbsRunDirPath() + "/db/tables.sql"
	createTables, err := ioutil.ReadFile(fpath)
	if err != nil {
		log.Error("Cannot init db: ", err)
		log.Error("    Current dir: ", util.GetAbsRunDirPath())
		log.Error("    file path:", fpath)
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
