package qiniu

import (
	"fmt"
	"os"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

type QiniuCdn struct {
	auth    *auth.Credentials
	conf    storage.Config
	buckets []string
}

func NewQiniuCdn(key, secret string, buckets []string) *QiniuCdn {
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Zone=&storage.ZoneHuabei
	return &QiniuCdn{
		auth:    auth.New(key, secret),
		conf:    cfg,
		buckets: buckets,
	}
}

func (q QiniuCdn) ListBucket(bucket string) {
	bucketManager := storage.NewBucketManager(q.auth, &q.conf)

	//列举所有文件
	prefix, delimiter, marker := "", "", ""
	entries, err := bucketManager.ListBucket(bucket, prefix, delimiter, marker)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ListBucket: %v\n", err)
		os.Exit(1)
	}
	i := 0
	for listItem := range entries {
		i++
		fmt.Printf("--[%d]--dir(%s)--Marker(%s)--Item(%+v)--\n", i, listItem.Dir, listItem.Marker, listItem.Item)
	}
}
