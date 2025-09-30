package db

import (
	"database/sql"
	"log"
)

func (d DB) CreateTables() (sql.Result, error) {
	return d.ExecSqlFile("init_tables.sql")
}

// AddRequest 添加网络请求记录
func (d DB) AddRequest(args ...any) (sql.Result, error) {
	return d.ExecSqlFile("requests_insert.sql", args...)
}

// AddBlockRequest 添加被拦截的请求记录
func (d DB) AddBlockRequest(args ...any) (sql.Result, error) {
	return d.ExecSqlFile("block_requests_insert.sql", args...)
}

// GetTopRequestIps 获取指定时间范围内，请求数最高的IP列表
// topIps: ip, request_count
// created_at > CURRENT_DATE + INTERVAL '1 hours'
// created_at < CURRENT_DATE + INTERVAL '5 hours'
func (d DB) GetTopRequestIps(topIps any, startAt, endAt string) error {
	return d.GetAllBySqlFile("get_top_request_ips.sql", &topIps, startAt, endAt)
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

// Prune. 定期对冗长的requests和block_requests数据表进行清理。防止数据过于庞大
func (d DB) Prune() error {
	// 如执行了多条SQL语句，则sql.Result.RowsAffected() 只能获取最后一条SQL语句影响的行数
	result, err := d.ExecSqlFile("prune.sql")
	if err != nil {
		return err
	}
	rownum, err := result.RowsAffected()
	log.Printf("-----PruneRequests---RowsAffected(%d)---\n", rownum)
	return err
}
