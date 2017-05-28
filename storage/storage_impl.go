package storage

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"

	"strings"

	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/hlcmsg"
	"github.com/lpimem/hlcsrv/util"
	_ "github.com/mattn/go-sqlite3" // sqlite driver
)

// SqliteStorage refers to one database instance.
type SqliteStorage struct {
	path string
	DB   *sql.DB
}

// NewSqliteStorage create a new database at the given path, or load the existing one.
func NewSqliteStorage(path string) *SqliteStorage {
	db, err := prepareSQLDb(path)
	if err != nil {
		panic(err)
	}
	return &SqliteStorage{path, db}
}

// Close connection
func (s *SqliteStorage) Close() {
	s.DB.Close()
}

// QueryMetaList returns all meta of user @p uid for a page @p pid
func (s *SqliteStorage) QueryMetaList(uid uint32, pid uint32) []*hlcmsg.RangeMeta {
	dict, err := s.QueryPagenote(uid, pid)
	if err != nil {
		log.Warn("Error:", err)
		return []*hlcmsg.RangeMeta{}
	}
	return dict.GetPagenote(uid, pid).GetHighlights()
}

//PagenoteAddon key: range id , value: strings for addtional values
type PagenoteAddon map[uint32][]interface{}

// QueryPagenoteFuzzy queries notes for all URIs that match uriPattern
func (s *SqliteStorage) QueryPagenoteFuzzy(
	uid UserID, uriPattern string) (PagenoteDict, PagenoteAddon, error,
) {
	var condBuilder bytes.Buffer
	var parameters = []interface{}{}
	if uid > 0 {
		condBuilder.WriteString(" and a.author = ?")
		parameters = append(parameters, uid)
	}
	condBuilder.WriteString(" and b.url like ?")
	condBuilder.WriteString(" order by a.mtime desc")
	parameters = append(parameters, fmt.Sprintf("%%%s%%", uriPattern))
	var cond = condBuilder.String()
	return s.doQueryPagenote(cond, ", b.title, b.url", "left join hlc_page b on a.page = b.id", parameters, true)
}

func buildPagenoteQuery(conds, fields, joins string) string {
	var queryBuilder bytes.Buffer
	queryBuilder.WriteString(`select 
a.id, a.anchor, a.start, a.startOffset, a.end, a.endOffset, a.page, a.author, a.option, a.text `)
	queryBuilder.WriteString(fields)
	queryBuilder.WriteString(` from hlc_range a `)
	queryBuilder.WriteString(joins)
	queryBuilder.WriteString(` where 1=1 `)
	queryBuilder.WriteString(conds)
	return queryBuilder.String()
}

func (s *SqliteStorage) doQueryPagenote(
	conds string,
	additionalFields string,
	joins string,
	parameters []interface{},
	withText bool) (result PagenoteDict, related PagenoteAddon, err error,
) {
	result = make(PagenoteDict)
	query := buildPagenoteQuery(conds, additionalFields, joins)
	nFields := 0
	if additionalFields != "" {
		nFields = len(strings.Split(additionalFields, ",")) - 1
		related = make(PagenoteAddon)
	}
	err = util.QueryDb(s.DB, query,
		parameters, func(rowno int, rows *sql.Rows) error {
			var id, startOffset, endOffset, page, author uint32
			var anchor, start, end, option, text string
			var valuePtrs = []interface{}{&id, &anchor, &start, &startOffset,
				&end, &endOffset, &page, &author, &option, &text}
			var additional []interface{}
			if nFields > 0 {
				additional = make([]interface{}, nFields, nFields)
				for i := range additional {
					valuePtrs = append(valuePtrs, &additional[i])
				}
			}
			err := rows.Scan(valuePtrs...)
			if err != nil {
				return err
			}

			if !withText {
				text = ""
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
				Text:        text,
				Option:      option,
			}
			note.Highlights = append(note.Highlights, meta)
			if nFields > 0 {
				related[id] = additional
			}
			return nil
		})
	return
}

// QueryPagenote returns all Pagenote of user @p uid for a page @p pid
func (s *SqliteStorage) QueryPagenote(uid uint32, pid uint32) (PagenoteDict, error) {
	if uid == 0 && pid == 0 {
		return PagenoteDict{}, errors.New("uid and url cannot both be 0")
	}
	var condBuilder bytes.Buffer
	var parameters = []interface{}{}
	if uid > 0 {
		condBuilder.WriteString(" and author = ?")
		parameters = append(parameters, uid)
	}
	if pid > 0 {
		condBuilder.WriteString(" and page = ?")
		parameters = append(parameters, pid)
	}
	var cond = condBuilder.String()
	result, _, err := s.doQueryPagenote(cond, "", "", parameters, false)
	return result, err
}

// NewRangeMeta creates new range meta, returns error if failed.
func (s *SqliteStorage) NewRangeMeta(
	uid uint32, pid uint32, m *hlcmsg.RangeMeta) (uint32, error,
) {
	r, err := s.DB.Exec(
		`insert into hlc_range(anchor, start, startOffset, end, 
endOffset, text, page, author, option)
values (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		m.Anchor, m.Start, m.StartOffset, m.End, m.EndOffset,
		m.Text, pid, uid, m.Option)
	if err != nil {
		return 0, err
	}
	lastID, err := r.LastInsertId()
	return uint32(lastID), err
}

//DeleteRangeMeta delete given meta, returns error if failed.
func (s *SqliteStorage) DeleteRangeMeta(id uint32) error {
	if id < 1 {
		return errors.New("invalid range meta id, should be > 0")
	}
	_, err := s.DB.Exec(`delete from hlc_range where id=?`, id)
	return err
}

// QueryPageID returns page id for a given URI, 0 if not found.
func (s *SqliteStorage) QueryPageID(url string) uint32 {
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

// QueryPage returns URI of a page id
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

// NewPage create new record for URI
func (s *SqliteStorage) NewPage(title, url string) (id uint32) {
	rst, err := s.DB.Exec(
		"insert into hlc_page (title, url) values (?, ?)",
		title, url,
	)
	if err != nil {
		log.Warn("ignored error:", err)
		return
	}
	lastID, err := rst.LastInsertId()
	if err != nil {
		log.Warn("ignored error:", err)
		id = 0
	} else {
		id = uint32(lastID)
	}
	return
}

// NewUser creates new user record, returns created user id
func (s *SqliteStorage) NewUser(name, email, password, slt string) (id uint32) {
	r, err := s.DB.Exec(
		"insert into hlc_user (name, email, password, _slt) values (?, ?, ?, ?)",
		name, email, password, slt,
	)
	if err != nil {
		log.Warn("ignored error: ", err)
		return
	}
	lastID, err := r.LastInsertId()
	if err != nil {
		log.Warn("ignored error: ", err)
		return
	}
	id = uint32(lastID)
	return
}

// QueryUser returns user id for given user handle and password
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
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if err := initDb(db); err != nil {
		return nil, err
	}
	return db, nil
}
