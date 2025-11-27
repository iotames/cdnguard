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
func GetMigrateFiles(fileKeys *[]string) error {
	// return getDB().GetMany("SELECT file_key, status FROM public.qiniu_cdnauth_file_migrate_list WHERE status = 0", dest)
	return getDB().GetAllBySqlFile("get_migrate_files.sql", fileKeys)
}

// GetDeleteFiles 获取待删除文件
func GetDeleteFiles(fileKeys *[]string) error {
	// return getDB().GetMany("SELECT file_key FROM public.qiniu_cdnauth_file_migrate_list WHERE status = 1", fileKeys)
	return getDB().GetAllBySqlFile("get_delete_files.sql", fileKeys)
}

// func BeginTx() (tx *sql.Tx, err error) {
// 	return getDB().GetSqlDB().Begin()
// }

// LogFileMigrate 变更迁移状态，添加文件操作记录。-1 操作失败0未开始，1copy成功，2move成功，3原文件已删除
// 仅opt=copy时，支持传入addPreDir参数。其他情况请传入空字符串
func LogFileMigrate(opt, file_key, from_bucket, to_bucket, addPreDir string) error {
	var err error
	tx := getDB().GetSqlDB()
	// tx, err := BeginTx()
	// if err != nil {
	// 	return err
	// }
	if opt == "copy" {
		result, err = tx.Exec(`UPDATE qiniu_cdnauth_file_migrate_list SET status=$1, updated_at=CURRENT_TIMESTAMP WHERE file_key=$2`, 1, file_key)
		if err != nil {
			// tx.Rollback()
			panic(err)
		}
		result, err = tx.Exec(`INSERT INTO qiniu_cdnauth_file_opt_log (file_key, opt_type, state, from_bucket, to_bucket, add_pre_dir)VALUES ($1, $2, $3, $4, $5, $6);`, file_key, 1, true, from_bucket, to_bucket, addPreDir)
		if err != nil {
			// tx.Rollback()
			panic(err)
		}
	}
	if opt == "delete" {
		_, err = tx.Exec(`UPDATE qiniu_cdnauth_file_migrate_list SET status=$1, updated_at=CURRENT_TIMESTAMP WHERE file_key=$2`, 3, file_key)
		if err != nil {
			panic(err)
		}
		_, err = tx.Exec(`INSERT INTO qiniu_cdnauth_file_opt_log (file_key, opt_type, state, from_bucket, to_bucket)VALUES ($1, $2, $3, $4, $5);`, file_key, 3, true, from_bucket, to_bucket)
		if err != nil {
			panic(err)
		}
	}
	// else {
	// 	tx.Rollback()
	// }
	// return tx.Commit()
	return err
}
