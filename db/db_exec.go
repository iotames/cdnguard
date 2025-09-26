package db

import (
	"database/sql"
)

func (d DB) Exec(sqltxt string, args ...interface{}) (sql.Result, error) {
	return d.edb.Exec(sqltxt, args...)
}

func (d DB) ExecSqlFile(filename string, args ...any) (sql.Result, error) {
	var err error
	var sqltxt string
	sqltxt, err = d.sqlDir.GetSQL(filename)
	if err != nil {
		return nil, err
	}
	return d.Exec(sqltxt, args...)
}
