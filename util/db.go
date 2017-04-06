package util

import "database/sql"

func InTxWithDB(db *sql.DB, ops []func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	return WithInTx(tx, ops)
}

func WithInTx(tx *sql.Tx, ops []func(tx *sql.Tx) error) error {
	for _, op := range ops {
		if err := op(tx); err != nil {
			return err
		}
	}
	return nil
}

func QueryDb(db *sql.DB, query string, args []interface{}, handler func(rowNo int, rows *sql.Rows) error) error {
	rows, err := db.Query(query, args...)
	return iterateRows(rows, err, handler)

}

func QueryTx(tx *sql.DB, query string, args []interface{}, handler func(rowNo int, rows *sql.Rows) error) error {
	rows, err := tx.Query(query, args...)
	return iterateRows(rows, err, handler)
}

func iterateRows(rows *sql.Rows, err error, handler func(rowNo int, rows *sql.Rows) error) error {
	if err != nil {
		return err
	}
	defer rows.Close()
	var current = 0
	for rows.Next() {
		err = handler(current, rows)
		if err != nil {
			return err
		}
		current += 1
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}
