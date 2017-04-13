package storage

import (
	"database/sql"

	"time"

	"github.com/lpimem/hlcsrv/util"
)

func QuerySession(
	sid string,
	uid uint32,
) (*time.Time, error) {
	const query = "select last_access from hlc_session where id = ? and uid = ?"
	var lastAccess time.Time
	err := util.QueryDb(storage.DB, query,
		[]interface{}{sid, uid},
		func(idx int, rows *sql.Rows) error {
			return rows.Scan(&lastAccess)
		})
	if err != nil {
		return nil, err
	}
	return &lastAccess, err
}

func UpdateSession(sid string, uid uint32) {
	storage.Upsert("hlc_session",
		[]string{"id", "uid", "last_access"},
		[]string{"id"},
		[]interface{}{
			sid, uid, time.Now(),
		},
		[]interface{}{sid})
}
