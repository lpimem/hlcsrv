package storage

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lpimem/hlcsrv/util"
)

// QuerySession checks if a session exists. If so, returns the last access time of the session, or error.
func QuerySession(
	sid string,
	uid UserID,
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

// QuerySessionByUID check if a user has a session record
func QuerySessionByUID(
	uid UserID,
) (*struct {
	Sid        string
	LastAccess *time.Time
}, error) {
	const query = "select id, last_access from hlc_session where uid = ?"
	var sid string
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

// UpdateSession refreshes a session
func UpdateSession(sid string, uid UserID) error {
	if uid < 1 {
		return errors.New("uid cannot be 0")
	}
	if sid == "" {
		return errors.New("sid cannot be empty")
	}
	_, err := Upsert(
		storage.DB,
		"hlc_session",
		[]string{"id", "uid", "last_access"},
		[]string{"id"},
		[]interface{}{
			sid, uid, time.Now(),
		},
		[]interface{}{sid})
	return err
}

// DeleteSession removes session with the given id.
func DeleteSession(sid string) error {
	_, err := storage.DB.Exec("delete from hlc_session where id = ?", sid)
	return err
}
