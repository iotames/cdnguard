package model

import (
	"database/sql"
	"log"
	"reflect"
)

var result sql.Result

// StatisRequestEveryDay 统计每天网络请求
// 每天凌晨 0:10 都会执行一个StatisRequestEveryDay定时任务：
// 1. 判断qiniu_cdnauth_statis表是否已存在statis_date为昨天的数据，没有就进行SQL统计
// 2. qiniu_cdnauth_requests 关联查询qiniu_cdnauth_files，得出昨天一整天的时间范围内，请求总次数request_count，请求的文件总大小request_size，结果以bucket_id进行分组
// 3. qiniu_cdnauth_block_requests关联查询qiniu_cdnauth_files，得出昨天一整天的时间范围内，拦截的请求总次数blocked_count，拦截请求的文件总大小blocked_size，还有各种拦截类型的次数统计，blocked_black_count，blocked_scanvul_count，blocked_webspider_count，blocked_useragent_count 结果以bucket_id进行分组
// 4. 把查询的结果整理成数据插入qiniu_cdnauth_statis数据表中。帮我把SQL语句写出来。
func StatisRequestEveryDay() (rownum int64, err error) {
	d := getDB()
	// 如执行了多条SQL语句，则sql.Result.RowsAffected() 只能获取最后一条SQL语句影响的行数
	result, err = d.ExecSqlFile("insert_statis.sql")
	if err != nil {
		return
	}
	rownum, err = result.RowsAffected()
	log.Printf("-----model.StatisRequestEveryDay---RowsAffected(%d)---\n", rownum)
	return
}

func DbStats() map[string]interface{} {
	d := getDB()
	stats := d.Stats()
	// 使用反射将 sql.DBStats 结构体中的字段转换为 map[string]interface{}
	result := make(map[string]interface{})
	// 获取 stats 的反射值
	v := reflect.ValueOf(stats)
	t := reflect.TypeOf(stats)
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		valuue := v.Field(i).Interface()
		result[field.Name] = valuue
	}
	return result
}

// stats := DBStats{
// 	MaxOpenConnections: db.maxOpen,
// 	Idle:            len(db.freeConn),
// 	OpenConnections: db.numOpen,
// 	InUse:           db.numOpen - len(db.freeConn),
// 	WaitCount:         db.waitCount,
// 	WaitDuration:      time.Duration(wait),
// 	MaxIdleClosed:     db.maxIdleClosed,
// 	MaxIdleTimeClosed: db.maxIdleTimeClosed,
// 	MaxLifetimeClosed: db.maxLifetimeClosed,
// }
