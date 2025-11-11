package qiniu

import (
	"encoding/base64"
	"fmt"
	"log"
	"log/slog"

	"github.com/qiniu/go-sdk/v7/storage"
)

type QiniuBucketFile struct {
	bucketManager     *storage.BucketManager
	prefix, delimiter string
	limit             int
	debug             bool
}

func NewQiniuBucketFile(bucketManager *storage.BucketManager, debug bool) *QiniuBucketFile {
	return &QiniuBucketFile{
		bucketManager: bucketManager,
		limit:         1000,
		debug:         debug,
	}
}

func (bf *QiniuBucketFile) SetPrefix(prefix string) {
	bf.prefix = prefix
}

func (bf *QiniuBucketFile) SetDelimiter(delimiter string) {
	bf.delimiter = delimiter
}

func (bf *QiniuBucketFile) SetLimit(limit int) {
	bf.limit = limit
}

func (bf QiniuBucketFile) ListFiles(bucketName, lastCursorMarker string) (entries []storage.ListItem, commonPrefixes []string, nextMarker string, hasNext bool, err error) {
	bucketManager := bf.bucketManager
	slog.Info("QiniuBucketFile.ListFiles Begin:", "bucketName", bucketName, "prefix", bf.prefix, "delimiter", bf.delimiter, "limit", bf.limit, "lastCursorMarker", lastCursorMarker)
	entries, commonPrefixes, nextMarker, hasNext, err = bucketManager.ListFiles(bucketName, bf.prefix, bf.delimiter, lastCursorMarker, bf.limit)
	if err != nil {
		err = fmt.Errorf("api request error:%v", err)
		return
	}
	if bf.debug {
		for i, entry := range entries {
			// {"c":0,"k":"filename.xlsx"}
			base64Src := fmt.Sprintf(`{"c":0,"k":"%s"}`, entry.Key)
			fcursor := base64.StdEncoding.EncodeToString([]byte(base64Src))
			log.Printf("QiniuBucketFile.ListFiles--entry-i[%d]-key(%s)-size(%d)-hash(%s)-cursor(%s)\n", i, entry.Key, entry.Fsize, entry.Hash, fcursor)
			// log.Println("", entry.Key, "file cursor:", fcursor)
		}
	}
	slog.Info("QiniuBucketFile.ListFiles End:", "nextMarker", nextMarker, "nextMarkerLen", len(nextMarker), "hasNext", hasNext)
	return
}

func (bf QiniuBucketFile) Copy(srcBucket, srcKey, destBucket, destKey string) error {
	// TODO 数据库里添加操作日志
	return bf.bucketManager.Copy(srcBucket, srcKey, destBucket, destKey, false)
}

func (bf QiniuBucketFile) Move(srcBucket, srcKey, destBucket, destKey string) error {
	// TODO 数据库里添加操作日志
	return bf.bucketManager.Move(srcBucket, srcKey, destBucket, destKey, false)
}
