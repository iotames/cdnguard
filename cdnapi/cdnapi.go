package cdnapi

import (
	"log"

	"github.com/iotames/cdnguard/cdnapi/internal/qiniu"
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
	var err error
	if c.cdnName == "qiniu" {
		qiniu := qiniu.NewQiniuCdn(c.key, c.secret, c.bucketNameList)
		err = qiniu.SyncFilesInfo(bucketName)
	}
	if err != nil {
		panic(err)
	}
	log.Println("-----SyncFilesInfo End-----")
}
