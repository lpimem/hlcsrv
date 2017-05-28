package storage

import (
	"database/sql"
	"github.com/go-playground/log"
	"github.com/lpimem/hlcsrv/util"
)

type restriction struct{}

// Restriction provided queries to restriction table
var Restriction restriction

func (r *restriction) Add(uri string) error {
	var err error
	if uri, err = validatePermissionURI(uri); err != nil {
		return err
	}
	if has, err := r.Has(uri); has || err != nil {
		if err != nil {
			log.Errorf("restricted.Add: %s", err)
		}
		return err
	}
	_, err = storage.DB.Exec("insert into `restriction` (uri) values (?)", uri)
	return err
}

func (r *restriction) Remove(uri string) error {
	var err error
	if has, err := r.Has(uri); !has || err != nil {
		return err
	}
	_, err = storage.DB.Exec("delete from `restriction` where uri = ?", uri)
	return err
}

func (*restriction) Has(uri string) (result bool, err error) {
	if uri, err = validatePermissionURI(uri); err != nil {
		return true, err
	}
	err = util.QueryDb(storage.DB, "select count(*) from `restriction` where instr(?, uri)", []interface{}{uri}, func(i int, r *sql.Rows) error {
		var count int
		if err := r.Scan(&count); err != nil {
			return err
		}
		result = count > 0
		return nil
	})
	return result, err
}

func (*restriction) All() (result []string, err error) {
	result = make([]string, 0)
	err = util.QueryDb(storage.DB, "select uri from `restriction`", nil, func(i int, r *sql.Rows) error {
		var uri string
		if err := r.Scan(&uri); err != nil {
			return nil
		}
		result = append(result, uri)
		return nil
	})
	return result, err
}
