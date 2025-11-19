package cdnapi

import (
	"fmt"
	"log"

	"github.com/iotames/cdnguard/cdnapi/internal"
	"github.com/iotames/cdnguard/cdnapi/internal/qiniu"
)

type CdnApi struct {
	cdnName        string
	key            string
	secret         string
	bucketNameList []string
}

// NewCdnApi 仅支持单个CDN服务商的数据处理。当前仅接入了七牛云的API。
// 因为数据列表仅有bucket_id标记同一个CDN服务商的数据桶ID，没办法区分不同服务商。
func NewCdnApi(cdnName string, key, secret string, buckets []string) *CdnApi {
	return &CdnApi{cdnName, key, secret, buckets}
}

// SyncFilesInfo 同步文件列表信息
func (c CdnApi) SyncFilesInfo(bucketName string) error {
	var err error
	var bucketId int
	bucketId, err = internal.GetBucketId(bucketName, c.bucketNameList)
	if bucketId == -1 || err != nil {
		return err
	}

	if c.cdnName == "qiniu" {
		qiniu := qiniu.NewQiniuCdn(c.key, c.secret, c.bucketNameList)
		err = qiniu.SyncFilesInfo(bucketName, bucketId)
	}
	if err != nil {
		return err
	}
	log.Println("-----SyncFilesInfo End-----")
	return nil
}

// ShowFilesInfo 显示文件列表信息。lastCursor表示从上一个文件的marker标记位之后开始显示。
func (c CdnApi) ShowFilesInfo(bucketName, lastCursor string) error {
	var err error
	if c.cdnName == "qiniu" {
		qiniu := qiniu.NewQiniuCdn(c.key, c.secret, c.bucketNameList)
		err = qiniu.ShowBucketFilesInfo(bucketName, lastCursor, 1000)
	}
	if err != nil {
		return err
	}
	return nil
}

func (c CdnApi) copyFiles(fromBucket, toBucket string, fileKeys []string, callback func(fkey string, err error)) error {
	if c.cdnName == "qiniu" {
		qiniu := qiniu.NewQiniuCdn(c.key, c.secret, c.bucketNameList)
		qiniu.BatchCopyFile(fromBucket, toBucket, fileKeys, callback)
		return nil
	}
	return fmt.Errorf("copy files only support qiniu")
}

func (c CdnApi) MigrateFiles(fromBucket, toBucket string, fileKeys []string, callback func(fkey string, err error)) error {
	return c.copyFiles(fromBucket, toBucket, fileKeys, callback)
}
