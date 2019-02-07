package storage

import (
	"database/sql"
	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/util"
)

type user struct {
}

// UserID is user's primary key
type UserID uint32

// UserInfo stores information about one user
type UserInfo struct {
	ID    UserID
	Name  string
	Email string
	CTime string
}

// User provides queries to user table
var User user

const defaultQueryLimit = 50

func (*user) All(limit, offset int) ([]UserInfo, error) {
	if offset < 0 {
		offset = 0
	}
	if limit < 0 {
		limit = defaultQueryLimit
	}
	query := "select id, name, email, ctime from hlc_user order by ctime asc limit ? offset ?"
	result := []UserInfo{}
	err := util.QueryDb(storage.DB, query, []interface{}{limit, offset},
		func(i int, r *sql.Rows) error {
			var u UserInfo
			if err := r.Scan(&u.ID, &u.Name, &u.Email, &u.CTime); err != nil {
				return err
			}
			result = append(result, u)
			log.Debugf("%d User %v", i, u)
			return nil
		},
	)
	return result, err
}

func (*user) Get(id UserID) (UserInfo, error) {
	return UserInfo{}, nil
}

func (*user) QueryByName(name string) ([]UserInfo, error) {
	return nil, nil
}

func (*user) QueryByEmail(email string) ([]UserInfo, error) {
	return nil, nil
}

// not needed as all users comes from Google Account
// func (*user) New(name, email, passwd string) (UserID, error) {
// 	return 0, nil
// }
