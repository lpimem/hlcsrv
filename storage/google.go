package storage

import (
	"database/sql"

	"errors"

	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/security"
	"github.com/lpimem/hlcsrv/util"
)

// GoogleTokenClaim represents fields extracted from an IDToken
type GoogleTokenClaim struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Sub           string `json:"sub"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// return uid for a a given user's google id
func GetOrCreateUidForGoogleUser(profile *GoogleTokenClaim) (uint32, error) {
	var uid uint32
	var err error
	uid = storage.QueryUidByGoogleId(profile.Sub)
	log.Debug("user id for google user ", profile.Sub, " not found, creating new")
	if uid < 1 {
		uid, err = storage.NewUserByGoogleProfile(profile)
		if err != nil {
			log.Info("cannot create user for google id ", profile.Sub, ": ", err)
			uid = storage.QueryUidByGoogleId(profile.Sub)
			if uid > 0 {
				err = nil
			}
		}
	}
	return uid, err
}

// return user id for a given google id
func (s *SqliteStorage) QueryUidByGoogleId(gid string) uint32 {
	const query = "select uid from hlc_google_auth where google_id=?"
	var id uint64
	util.QueryDb(s.DB, query, []interface{}{gid}, func(idx int, rows *sql.Rows) error {
		return rows.Scan(&id)
	})
	return uint32(id)
}

// create user with given google id
func (s *SqliteStorage) NewUserByGoogleProfile(profile *GoogleTokenClaim) (newUserId uint32, err error) {
	const query = `insert into hlc_user (name, email, password, _slt) values (?, ?, ?, ?);`
	const query_2 = `insert into hlc_google_auth(google_id, uid, picture) values (?, ?, ?);`
	const passwdStrength = 32
	var passwd = security.RandStringBytesMaskImprSrc(passwdStrength)
	var hash, slt = security.Hash(passwd)
	err = util.InTxWithDB(s.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			log.Debugf("%s | %s, %s, ...", query, profile.Name, profile.Email)
			r, err := tx.Exec(query, profile.Name, profile.Email, hash, slt)
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
			if _, err := tx.Exec(query_2, profile.Sub, newUserId, profile.Picture); err != nil {
				return err
			}
			return nil
		},
	})
	return
}
