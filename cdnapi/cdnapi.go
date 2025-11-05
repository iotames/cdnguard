package cdnapi

import (
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
