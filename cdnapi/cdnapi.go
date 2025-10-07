package cdnapi

import (
	"github.com/iotames/cdnguard/cdnapi/qiniu"
)

type CdnApi struct {
	cdnName string
	key     string
	secret  string
}

func NewCdnApi(cdnName string, key, secret string) *CdnApi {
	return &CdnApi{cdnName, key, secret}
}

func (c CdnApi) SyncFiles(bucketName string) {
	if c.cdnName == "qiniu" {
		qiniu := qiniu.NewQiniuCdn(c.key, c.secret)
		qiniu.SyncFiles(bucketName)
	}
}
