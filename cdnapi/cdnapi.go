package cdnapi

import (
	"github.com/iotames/cdnguard/cdnapi/qiniu"
)

type CdnApi struct {
	cdnName        string
	key            string
	secret         string
	bucketNameList []string
}

func NewCdnApi(cdnName string, key, secret string, buckets []string) *CdnApi {
	return &CdnApi{cdnName, key, secret, buckets}
}

func (c CdnApi) SyncFiles(bucketName string) {
	if c.cdnName == "qiniu" {
		qiniu := qiniu.NewQiniuCdn(c.key, c.secret, c.bucketNameList)
		qiniu.SyncFiles(bucketName)
	}
}
