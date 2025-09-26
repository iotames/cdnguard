package db

import (
	"database/sql"
)

func (d DB) CreateTables() (sql.Result, error) {
	return d.ExecSqlFile("tables_init.sql")
}

func (d DB) AddRequest(args ...any) (sql.Result, error) {
	return d.ExecSqlFile("requests_insert.sql", args...)
}
