package util

import (
	"database/sql"
	"io/ioutil"
	"os"
	"testing"

	"github.com/go-playground/log"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func TestInTxWithDB(t *testing.T) {
	type testcase struct {
		name    string
		queries []string
		commit  bool
	}
	tcs := []*testcase{
		{"valid transaction", []string{
			`insert into hlc_user(name, email, password, _slt)
			values ("ExAm1", "example2@gmail.com", "unsafe", "unsafe");`,
			`insert into hlc_google_auth(google_id, uid)
			values ("fakeid", 11);`,
		}, true},
		{"invalid transaction", []string{
			`insert into hlc_user(name, email, password, _slt)
			values ("ExAm2", "example3@gmail.com", "unsafe", "unsafe");`,
			`insert into hlc_google_auth(google_id, uid)
			values ("fakeid", 12);`,
		}, false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			var (
				cnt_user, cnt_google_auth     uint64
				cnt_user_2, cnt_google_auth_2 uint64
				err                           error
				committed                     bool
				consist                       bool
			)
			if cnt_user, err = count(db, "hlc_user"); err != nil {
				t.Error(err)
				t.Fail()
				return
			}
			if cnt_google_auth, err = count(db, "hlc_google_auth"); err != nil {
				t.Error(err)
				t.Fail()
				return
			}
			err = InTxWithDB(db, []func(tx *sql.Tx) error{
				func(tx *sql.Tx) error {
					if _, err := tx.Exec(tc.queries[0]); err != nil {
						return err
					}
					return nil
				},
				func(tx *sql.Tx) error {
					if _, err := tx.Exec(tc.queries[1]); err != nil {
						return err
					}
					return nil
				},
			})
			if err != nil {
				log.Warn(err)
			}
			if cnt_user_2, err = count(db, "hlc_user"); err != nil {
				t.Error(err)
				t.Fail()
				return
			}
			if cnt_google_auth_2, err = count(db, "hlc_google_auth"); err != nil {
				t.Error(err)
				t.Fail()
				return
			}
			committed = cnt_user < cnt_user_2 && cnt_google_auth < cnt_google_auth_2
			consist = cnt_user_2-cnt_user == cnt_google_auth_2-cnt_google_auth
			if tc.commit != (committed && consist) {
				if err != nil {
					t.Error(err)
				}
				t.Fail()
			}
		})
	}
}

func count(db *sql.DB, tb string) (uint64, error) {
	rows, err := db.Query("select count(*) from " + tb)
	if err != nil {
		return 0, err
	}
	var c uint64
	c = 0
	err = IterateRows(rows, nil, func(rowNo int, rows *sql.Rows) error {
		return rows.Scan(&c)
	})
	if err != nil {
		return 0, err
	}
	return c, err
}

func init() {
	var err error
	dbpath := GetAbsRunDirPath() + "/db/unittest.db"
	log.Info(dbpath)
	if _, err = os.Stat(dbpath); os.IsExist(err) {
		os.Remove(dbpath)
	}
	db, err = sql.Open("sqlite3", dbpath)
	if err != nil {
		log.Error(err)
		return
	}

	fpath := GetAbsRunDirPath() + "/db/tables.sql"
	createTables, err := ioutil.ReadFile(fpath)
	if err != nil {
		log.Error("Cannot init db: ", err)
		log.Error("    Current dir: ", GetAbsRunDirPath())
		log.Error("    file path:", fpath)
		return
	}
	_, err = db.Exec(string(createTables))
	if err != nil {
		log.Error(err)
		return
	}
	_, err = db.Exec(`insert into hlc_user(id, name, email, password, _slt)
		values (1, "Bob", "bob@example.com", "unsafe", "unsafe");
		values (10, "ExAm", "example@gmail.com", "unsafe", "unsafe");

	insert into hlc_page(id, title, url)
		values (1, "example", "http://example.com");

	insert into hlc_range(id, anchor, start, startOffset, end, endOffset, text, page, author)
		values (1, "#c", "#c/1", 0, "#c/12", 32, "This is the selected text", 1, 1);

	insert into hlc_google_auth(google_id, uid) values ("100000", 1);
	insert into hlc_google_auth(google_id, uid) values ("example@gmail.com", 10);

	insert into hlc_session(id, uid, last_access) values ("fake_session_id", 10, CURRENT_TIMESTAMP);
	`)
	if err != nil {
		log.Error(err)
	}
}
