package storage

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/log"
)

type Executable interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

/**Upsert performs `update on insert duplicated` query.
Note, this method uses unchecked string format to build the query.
You should never use user's input as the value for table, fields, and keys.
*/
func Upsert(
	s Executable,
	table string,
	fields []string,
	keys []string,
	parameters []interface{},
	keyValues []interface{},
) (r sql.Result, err error) {
	if len(fields) != len(parameters) {
		return nil, errors.New("Numbers of fields and parameters do not match.")
	}
	if len(keys) != len(keyValues) {
		return nil, errors.New("Numbers of keys and keyValues do not match.")
	}
	if len(keys) < 1 {
		return nil, errors.New("Numbers of keys cannot be 0")
	}

	var (
		snippetUpdateFields string
		snippetUpdateCond   string
		snippetInsertFields string
		snippetInsertValues string
		query               string
		queryParameters     []interface{}
	)
	const template = `
		UPDATE %s
		SET %s
		WHERE %s;
		INSERT INTO %s (%s)
		SELECT %s
		WHERE (Select Changes() = 0);`

	snippetInsertFields = strings.Join(fields, ",")
	snippetInsertValues = strings.Repeat("?,", len(fields)-1) + "?"

	var updateFieldsBuilder bytes.Buffer
	var fieldsNumber = len(fields) - 1
	for i, f := range fields {
		updateFieldsBuilder.WriteString(f)
		updateFieldsBuilder.WriteString(" = ?")
		if i < fieldsNumber {
			updateFieldsBuilder.WriteString((","))
		}
	}
	snippetUpdateFields = updateFieldsBuilder.String()

	var updateCondBuilder bytes.Buffer
	var keysNumber = len(keys) - 1
	for i, k := range keys {
		updateCondBuilder.WriteString(k)
		updateCondBuilder.WriteString(" = ?")
		if i < keysNumber {
			updateCondBuilder.WriteString((","))
		}
	}
	snippetUpdateCond = updateCondBuilder.String()

	query = fmt.Sprintf(template, table, snippetUpdateFields, snippetUpdateCond,
		table, snippetInsertFields, snippetInsertValues)
	log.Debug(query)

	queryParameters = append(parameters, keyValues...)
	queryParameters = append(queryParameters, parameters...)

	r, err = s.Exec(query, queryParameters...)
	return
}
