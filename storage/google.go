package storage

import (
	"database/sql"

	"github.com/lpimem/hlcsrv/util"
)

func (s *SqliteStorage) QueryUidByGoogleId(gid string) uint32 {
	const query = "select uid from hlc_google_auth where google_id=?"
	var id uint64
	util.QueryDb(s.DB, query, []interface{}{gid}, func(idx int, rows *sql.Rows) error {
		return rows.Scan(&id)
	})
	return uint32(id)
}
