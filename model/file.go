package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/storage"
)

func AddSyncLog(bucketId int, hasNext bool, cursorMarker string, sizeLen, sort int) (result sql.Result, err error) {
	sql := `INSERT INTO qiniu_cdnauth_file_sync_log (bucket_id,has_next,cursor_marker,size_len,sort) VALUES ($1,$2,$3,$4,$5)`
	return getDB().Exec(sql, bucketId, hasNext, cursorMarker, sizeLen, sort)
}

// 获取最近一条FileSyncLog同步记录
func GetLastSyncLog(bucketId int, hasNext *bool, CursorMarker *string, sizeLen, id *int) error {
	sql := `SELECT has_next, cursor_marker, size_len, id FROM qiniu_cdnauth_file_sync_log WHERE bucket_id = $1 ORDER BY id DESC LIMIT 1`
	return getDB().GetOne(sql, []any{hasNext, CursorMarker, sizeLen, id}, bucketId)
}

// // 更新同步日志状态
// func UpdateSyncLogStatus(id int, hasNext bool) (result sql.Result, err error) {
// 	sql := `UPDATE qiniu_cdnauth_file_sync_log SET has_next=$1, updated_at=CURRENT_TIMESTAMP WHERE id=$3`
// 	return getDB().Exec(sql, hasNext, id)
// }

func BatchInsertQiniuFiles(bucketId int, files []storage.ListItem) (result sql.Result, err error) {
	sql := `INSERT INTO qiniu_cdnauth_files (file_key, file_size, file_hash, md5, mime_type, file_type, upload_time, bucket_id, status, request_count, data_raw) VALUES %s`
	sqlvalues := []string{}
	sqli := 1
	sqlargs := []interface{}{}
	for _, f := range files {
		// 上传时间，单位：100纳秒，其值去掉低七位即为Unix时间戳。
		upload_time := time.Unix(f.PutTime/10000000, 0)
		fjson, _ := json.Marshal(f)
		// fmt.Printf("-----hash(%s)%d--md5(%s)--mimetype(%s)%d---filekey(%s)--uploadTime(%s)--\n", f.Hash, len(f.Hash), f.Md5, f.MimeType, len(f.MimeType), f.Key, upload_time)
		sqlargs = append(sqlargs, f.Key, f.Fsize, f.Hash, f.Md5, f.MimeType, f.Type, upload_time, bucketId, f.Status, 0, fjson)
		sqlvalues = append(sqlvalues, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)", sqli, sqli+1, sqli+2, sqli+3, sqli+4, sqli+5, sqli+6, sqli+7, sqli+8, sqli+9, sqli+10))
		sqli += 11
	}
	sql = fmt.Sprintf(sql, strings.Join(sqlvalues, ",")+";")
	return getDB().Exec(sql, sqlargs...)
}

type MigrateFile struct {
	FileUrl    string `db:"file_url"`
	FileKey    string `db:"file_key"`
	Status     int    `db:"status"`
	FromBucket string `db:"from_bucket"`
}

// GetMigrateFiles 获取待迁移文件
func GetMigrateFiles(dest *[]MigrateFile) error {
	return getDB().GetMany("SELECT * FROM public.qiniu_cdnauth_file_migrate_list WHERE status = 0", dest)
}

// UpdateMigrateFileStatus 变更迁移状态：-1 操作失败0未开始，1copy成功，2move成功，3原文件已删除
func UpdateMigrateFileStatus(fileKey string, status int) (result sql.Result, err error) {
	sql := `UPDATE qiniu_cdnauth_file_migrate_list SET status=$1, updated_at=CURRENT_TIMESTAMP WHERE file_key=$2`
	// 添加文件操作记录
	return getDB().Exec(sql, status, fileKey)
}
