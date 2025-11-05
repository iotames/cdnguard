package qiniu

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/iotames/cdnguard/model"
	"github.com/qiniu/go-sdk/v7/storage"
)

var result sql.Result

// var hashNext bool

// SyncFilesInfo 同步某个存储空间下的文件列表信息到数据表中
func (q QiniuCdn) SyncFilesInfo(bucketName string, bucketId int) error {
	var err error

	if bucketId == -1 {
		return fmt.Errorf("bucketName(%s) not in BUCKET_NAME_LIST. Please look the file: .env", bucketName)
	}

	// 参数准备
	// marker 上一次列举返回的位置标记，作为本次列举的起点信息。默认值为空字符串。
	// prefix 指定前缀，只有资源名匹配该前缀的资源会被列出。默认值为空字符串。
	// delimiter 指定目录分隔符，列出所有公共前缀（模拟列出目录效果）。默认值为空字符串。
	limit := 1000
	// prefix := "qiniu"
	prefix := ""
	delimiter := ""
	//初始列举marker为空
	beginMarker := ""
	hasNext := false
	id := 0
	lastSyncLogFilesLen := 0
	lastMarker := ""
	err = model.GetLastSyncLog(bucketId, &hasNext, &lastMarker, &lastSyncLogFilesLen, &id)
	if err != nil {
		return fmt.Errorf("SyncFilesInfo GetLastSyncLog error: %v,", err)
	}
	log.Println("lastSyncLogFilesLen:", lastSyncLogFilesLen)

	sort := 0
	bucketManager := storage.NewBucketManager(q.auth, &q.conf)

	// 检查是否有新增文件
	if id > 0 {
		if !hasNext {
			var nextMarker string
			var entries []storage.ListItem
			// 通过API网络请求接口，查看文件列表数据
			entries, _, nextMarker, hasNext, err = bucketManager.ListFiles(bucketName, prefix, delimiter, lastMarker, limit)
			// CursorMarker 为每次请求返回的最后一条数据，【既定规则下的Base64 编码】：base64Encode(fmt.Sprintf(`{"c":0,"k":"%s"}`, "file_key"))
			// TODO 校验七牛云在指定bucket文件列表中，同步的最后一个文件在【既定规则下的Base64 编码】，是否与同步日志的buket_file_sync_log表中的cursor_marker一致。
			if err != nil {
				// API网络请求获取新文件列表数据失败
				return fmt.Errorf("last file_sync_log get new ListFiles error: %v", err)
			}

			lenFiles := len(entries)
			if lenFiles == 0 {
				return fmt.Errorf("SyncFilesInfo error: 没有新增文件可以同步。")
			}

			// // 因为有了新增文件。所以继续同步前，先更新 file_sync_log 表最近那条记录：set has_next=true
			// _, err = model.UpdateSyncLogStatus(id, true)
			// if err != nil {
			// 	return fmt.Errorf("UpdateSyncLogStatus error: %v", err)
			// }

			// 有新增文件，但已经没有下一页了。
			if !hasNext {
				// 只有一页新增文件，直接保存并返回
				return q.saveFilesToDb(bucketId, nextMarker, hasNext, entries, &sort)
			}

			// 有新增文件，且还有下一页
			// 保存新增文件到数据库，并添加同步记录
			err = q.saveFilesToDb(bucketId, nextMarker, hasNext, entries, &sort)
			if err != nil {
				return fmt.Errorf("同步新增文件列表到数据库发生异常: %v", err)
			}
			// 修改新一轮数据同步的起始标记位
			beginMarker = nextMarker
		} else {
			// 修改新一轮数据同步的起始标记位
			beginMarker = lastMarker
		}
	}
	log.Println("-----AsyncFiles Start-----", bucketName, beginMarker)
	q.syncFilesBatch(bucketManager, bucketId, bucketName, prefix, delimiter, beginMarker, limit, &sort)
	return nil
}

// syncFilesBatch 封装文件同步的批处理逻辑
func (q QiniuCdn) syncFilesBatch(bucketManager *storage.BucketManager, bucketId int, bucketName, prefix, delimiter, marker string, limit int, sort *int) {
	// TODO 数据库操作错误时使用 break退出循环，但没有回滚机制。如果插入部分批次后失败，数据可能不一致。
	var err error
	var hasNext bool
	var nextMarker string

	for {
		nextMarker, hasNext, err = q.syncByApi(bucketManager, bucketId, bucketName, prefix, delimiter, marker, limit, sort)
		if err != nil {
			log.Println("syncFilesBatch error,", err)
			break
		}
		if !hasNext {
			log.Println("-----SyncFilesInfo End-----")
			break
		}
		marker = nextMarker
	}
}

func (q QiniuCdn) syncByApi(bucketManager *storage.BucketManager, bucketId int, bucketName, prefix, delimiter, marker string, limit int, sort *int) (nextMarker string, hashNext bool, err error) {
	// 从起始标记位获取文件列表
	var entries []storage.ListItem
	// TODO 如果按同步时间逆序，可能后面新增的文件会无法继续同步。
	entries, _, nextMarker, hashNext, err = bucketManager.ListFiles(bucketName, prefix, delimiter, marker, limit)
	if err != nil {
		err = fmt.Errorf("api request error:%v", err)
		return
	}
	log.Println("hashNext:", hashNext)
	log.Println("nextMarker:", nextMarker, len(nextMarker))
	err = q.saveFilesToDb(bucketId, nextMarker, hashNext, entries, sort)
	return
}

// saveFilesToDb 同步保存新增文件列表到数据库，并添加同步记录
func (q QiniuCdn) saveFilesToDb(bucketId int, lastMarker string, hashNext bool, entries []storage.ListItem, sort *int) (err error) {
	// 把文件列表数据批量插入数据库
	result, err = model.BatchInsertQiniuFiles(bucketId, entries)
	if err != nil {
		err = fmt.Errorf("BatchInsertQiniuFiles error: %v,", err)
		return
	} else {
		rowsaffected, _ := result.RowsAffected()
		lastinsertid, _ := result.LastInsertId()
		log.Println("BatchInsertQiniuFiles result:", rowsaffected, lastinsertid)
	}
	// 添加一条同步记录
	result, err = model.AddSyncLog(bucketId, hashNext, lastMarker, len(entries), *sort)
	if err != nil {
		err = fmt.Errorf("AddSyncLog error: %v", err)
		return
	} else {
		rowsaffected, _ := result.RowsAffected()
		lastinsertid, _ := result.LastInsertId()
		log.Println("AddSyncLog result:", rowsaffected, lastinsertid)
	}
	*sort++
	return
}

func (q QiniuCdn) ShowFilesInfo(bucketName, lastCursor string, limit int) error {
	if limit == 0 {
		limit = 1000
	}
	bucketManager := storage.NewBucketManager(q.auth, &q.conf)
	entries, _, nextMarker, hasNext, err := bucketManager.ListFiles(bucketName, "", "", lastCursor, limit)
	if err != nil {
		return fmt.Errorf("api request error:%v", err)
	}
	log.Println("ShowFilesInfo--hashNext:", hasNext)
	for i, entry := range entries {
		base64Src := fmt.Sprintf(`{"c":0,"k":"%s"}`, entry.Key)
		fcursor := base64.StdEncoding.EncodeToString([]byte(base64Src))
		log.Printf("ShowFilesInfo--entry-i[%d]-key(%s)-size(%d)-hash(%s)-cursor(%s)\n", i, entry.Key, entry.Fsize, entry.Hash, fcursor)
		// log.Println("", entry.Key, "file cursor:", fcursor)
	}
	log.Println("ShowFilesInfo--nextMarker:", nextMarker)
	log.Println("ShowFilesInfo--hashNext:", hasNext)
	return nil
}
