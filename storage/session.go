package storage

import (
	"database/sql"

	"time"

	"errors"

	"github.com/lpimem/hlcsrv/util"
)

func QuerySession(
	sid string,
	uid uint32,
) (*time.Time, error) {
	const query = "select last_access from hlc_session where id = ? and uid = ?"
	var lastAccess time.Time
	var miss = true
	err := util.QueryDb(storage.DB, query,
		[]interface{}{sid, uid},
		func(idx int, rows *sql.Rows) error {
			miss = false
			return rows.Scan(&lastAccess)
		})
	if err != nil || miss {
		return nil, err
	}
	return &lastAccess, err
}

func QuerySessionByUid(
	uid uint32,
) (*struct {
	Sid        string
	LastAccess *time.Time
}, error) {
	const query = "select id, last_access from hlc_session where uid = ?"
	var sid string = ""
	var lastAccess time.Time
	err := util.QueryDb(storage.DB, query,
		[]interface{}{uid},
		func(idx int, rows *sql.Rows) error {
			return rows.Scan(&sid, &lastAccess)
		})
	if err != nil {
		return nil, err
	}
	if sid == "" {
		return nil, nil
	}
	var result = struct {
		Sid        string
		LastAccess *time.Time
	}{sid, &lastAccess}
	return &result, nil
}

func UpdateSession(sid string, uid uint32) error {
	if uid < 1 {
		return errors.New("uid cannot be 0")
	}
	if sid == "" {
		return errors.New("sid cannot be empty")
	}
	_, err := storage.Upsert(
		"hlc_session",
		[]string{"id", "uid", "last_access"},
		[]string{"id"},
		[]interface{}{
			sid, uid, time.Now(),
		},
		[]interface{}{sid})
	return err
}
