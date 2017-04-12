package storage

import (
	"database/sql"

	"errors"

	"github.com/lpimem/hlcsrv/security"
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

func (s *SqliteStorage) NewUserByGoogleId(gid string, email string) (newUserId uint32, err error) {
	const query = `insert into hlc_user (email, password, _slt) values (?, ?, ?);`
	const query_2 = `insert into hlc_google_auth(google_id, uid) values (?, ?);`
	const passwdStrength = 32
	var passwd = security.RandStringBytesMaskImprSrc(passwdStrength)
	var hash, slt = security.Hash(passwd)
	err = util.InTxWithDB(s.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			r, err := tx.Exec(query, email, hash, slt)
			if err != nil {
				return err
			}
			uid, err := r.LastInsertId()
			newUserId = uint32(uid)
			return nil
		},
		func(tx *sql.Tx) error {
			if newUserId <= 0 {
				return errors.New("new user id should be > 0")
			}
			if _, err := tx.Exec(query_2, gid, newUserId); err != nil {
				return err
			}
			return nil
		},
	})
	return
}
