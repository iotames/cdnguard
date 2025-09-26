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

func (d DB) GetDbSizeText() (string, error) {
	var sizetxt string
	err := d.GetOneBySqlFile("database_size.sql", []any{&sizetxt})
	return sizetxt, err
}

func (d DB) GetIpWhiteList() ([]string, error) {
	var ip_list []string
	err := d.GetAllBySqlFileReplace("ip_list.sql", &ip_list, "qiniu_cdnauth_ip_white_list")
	return ip_list, err
}

func (d DB) GetIpBlackList() ([]string, error) {
	var ip_list []string
	err := d.GetAllBySqlFileReplace("ip_list.sql", &ip_list, "qiniu_cdnauth_ip_black_list")
	return ip_list, err
}
