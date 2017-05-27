package storage

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/util"
)

type permission struct {
}

// Permission is singleton to perform permission queries.
var Permission permission

func (*permission) ForUser(uid uint32) ([]string, error) {
	uris := []string{}
	err := util.QueryDb(storage.DB,
		"select uri from permission where user=?",
		[]interface{}{uid},
		func(rowNo int, rows *sql.Rows) error {
			var uri string
			rows.Scan(&uri)
			uris = append(uris, uri)
			return nil
		})
	return uris, err
}

func (*permission) ToURI(uri string) ([]uint32, error) {
	users := []uint32{}
	var err error
	uri, err = validatePermissionURI(uri)
	if err != nil {
		return users, err
	}
	err = util.QueryDb(storage.DB,
		"select user from permission where uri=?",
		[]interface{}{uri},
		func(rowNo int, rows *sql.Rows) error {
			var uid uint32
			if err := rows.Scan(&uid); err != nil {
				return err
			}
			users = append(users, uid)
			return nil
		})
	return users, err
}

func validatePermissionURI(uri string) (string, error) {
	var err error
	if strings.Index(uri, "/") != 0 {
		err = errors.New("permission URI must start with /")
		return uri, err
	}
	if len(uri) <= 1 {
		err = errors.New("permission URI must contains at least one component")
		return uri, err
	}
	if !strings.HasSuffix(uri, "/") {
		uri = uri + "/"
	}
	return uri, err
}

func (p *permission) Grant(uid uint32, uri string) error {
	var err error
	uri, err = validatePermissionURI(uri)
	if err != nil {
		return err
	}
	if acc, err := p.HasAccess(uid, uri); nil != err || acc {
		return err
	}
	_, err = storage.DB.Exec(`insert into permission (user, uri) values (?, ?)`, uid, uri)
	if err != nil {
		log.Errorf("cannot grant access for %d, %s. Reason: %s", uid, uri, err)
	}
	return err
}

func (*permission) Revoke(uid uint32, uri string) error {
	// var err error
	// uri, err = validatePermissionURI(uri)
	// if err != nil {
	// 	return err
	// }
	// _, err = storage.DB.Exec("delete from permission where user = ? and uri = ?", uid, uri)
	return errors.New("Revoke is not supported")
}

// HasAccess returns true if there is at least one record in permission table of which
// the uri values is a prefix of the parameter uri for user with uid
// Note: this function should only be used for non-administrator roles.
func (*permission) HasAccess(uid uint32, uri string) (bool, error) {
	var err error
	if "/" == uri {
		return false, err
	}
	uri, err = validatePermissionURI(uri)
	if err != nil {
		return false, err
	}
	query, values := buildHasAccessQuery(uid, uri)
	if err != nil {
		return false, err
	}
	var result bool
	err = util.QueryDb(storage.DB, query, values,
		func(rowNo int, rows *sql.Rows) error {
			var count int
			if err := rows.Scan(&count); err != nil {
				return err
			}
			log.Debug(count, " permission records for ", values)
			result = count > 0
			return nil
		},
	)
	return result, err
}

func buildHasAccessQuery(uid uint32, uri string) (query string, values []interface{}) {
	query = "select count(*) from permission where user = ? and instr(?, uri)"
	values = []interface{}{uid, uri}
	return query, values
}
