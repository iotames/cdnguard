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
func GetLastSyncLog(bucketId int, hasNext *bool, CursorMarker *string, id *int) error {
	sql := `SELECT has_next, cursor_marker, id FROM qiniu_cdnauth_file_sync_log WHERE bucket_id = $1 ORDER BY id DESC LIMIT 1`
	return getDB().GetOne(sql, []any{hasNext, CursorMarker, id}, bucketId)
}

func BatchInsertQiniuFiles(bucketId int, files []storage.ListItem) (result sql.Result, err error) {
	sql := `INSERT INTO qiniu_cdnauth_files (file_key, file_size, file_hash, mime_type, file_type, upload_time, bucket_id, status, request_count, data_raw) VALUES %s`
	sqlvalues := []string{}
	sqli := 1
	sqlargs := []interface{}{}
	for _, f := range files {
		// 上传时间，单位：100纳秒，其值去掉低七位即为Unix时间戳。
		upload_time := time.Unix(f.PutTime/10000000, 0)
		fjson, _ := json.Marshal(f)
		// fmt.Printf("-----hash(%s)%d----mimetype(%s)%d------\n", f.Hash, len(f.Hash), f.MimeType, len(f.MimeType))
		sqlargs = append(sqlargs, f.Key, f.Fsize, f.Hash, f.MimeType, f.Type, upload_time, bucketId, f.Status, 0, fjson)
		sqlvalues = append(sqlvalues, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)", sqli, sqli+1, sqli+2, sqli+3, sqli+4, sqli+5, sqli+6, sqli+7, sqli+8, sqli+9))
		sqli += 10
	}
	sql = fmt.Sprintf(sql, strings.Join(sqlvalues, ",")+";")
	return getDB().Exec(sql, sqlargs...)
}
