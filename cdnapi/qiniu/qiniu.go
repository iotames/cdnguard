package qiniu

import (
	"fmt"
	"log"
	"os"

	"database/sql"

	"github.com/iotames/cdnguard/model"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

// # bucketname为空间名称。不添加参数，默认为：wildto。可用值为: wildto, wildto-private, buerdiy, buerdiy-staging, santic-pan, sagriatech-private, santic, newwildto

// '空间名：0wildto，1wildto-private';
const BUCKET_WILDTO = 0
const BUCKET_WILDTO_PRIVATE = 1
const BUCKET_BUERDIY = 2
const BUCKET_BUERDIY_STAGING = 3
const BUCKET_SANTIC_PAN = 4
const BUCKET_SAGRIATECH_PRIVATE = 5
const BUCKET_SANTIC = 6
const BUCKET_NEWWILDTO = 7

func GetBucketId(bucketName string) int {
	switch bucketName {
	case "wildto":
		return BUCKET_WILDTO
	case "wildto-private":
		return BUCKET_WILDTO_PRIVATE
	case "buerdiy":
		return BUCKET_BUERDIY
	case "buerdiy-staging":
		return BUCKET_BUERDIY_STAGING
	case "santic-pan":
		return BUCKET_SANTIC_PAN
	case "sagriatech-private":
		return BUCKET_SAGRIATECH_PRIVATE
	case "santic":
		return BUCKET_SANTIC
	case "newwildto":
		return BUCKET_NEWWILDTO
	default:
		return -1
	}
	return -1
}

type QiniuCdn struct {
	auth *auth.Credentials
	conf storage.Config
}

func NewQiniuCdn(key, secret string) *QiniuCdn {
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Zone=&storage.ZoneHuabei
	return &QiniuCdn{
		auth: auth.New(key, secret),
		conf: cfg,
	}
}

func (q QiniuCdn) SyncFiles(bucketName string) {
	var result sql.Result
	var err error
	bucketId := GetBucketId(bucketName)
	if bucketId < 0 {
		panic("bucketName error")
	}
	// marker 上一次列举返回的位置标记，作为本次列举的起点信息。默认值为空字符串。
	// prefix 指定前缀，只有资源名匹配该前缀的资源会被列出。默认值为空字符串。
	// delimiter 指定目录分隔符，列出所有公共前缀（模拟列出目录效果）。默认值为空字符串。

	bucketManager := storage.NewBucketManager(q.auth, &q.conf)
	limit := 1000
	// prefix := "qiniu"
	prefix := ""
	delimiter := ""
	//初始列举marker为空
	marker := ""
	hasNext := false
	id := 0
	err = model.GetLastSyncLog(bucketId, &hasNext, &marker, &id)
	if err != nil {
		fmt.Println("GetLastSyncLog error,", err)
		return
	}
	if !hasNext && id > 0 {
		fmt.Println("-----AsyncFiles End-----")
		return
	}
	fmt.Println("-----AsyncFiles Start-----", bucketName, marker)

	sort := 0
	for {
		entries, _, nextMarker, hashNext, err := bucketManager.ListFiles(bucketName, prefix, delimiter, marker, limit)
		if err != nil {
			errmsg := fmt.Sprintf("api request error:%v", err)
			log.Println(errmsg)
			break
		}
		result, err = model.BatchInsertQiniuFiles(bucketId, entries)
		if err != nil {
			fmt.Println("BatchInsertQiniuFiles error,", err)
			break
		} else {
			rowsaffected, _ := result.RowsAffected()
			lastinsertid, _ := result.LastInsertId()
			fmt.Println("BatchInsertQiniuFiles result:", rowsaffected, lastinsertid)
		}

		fmt.Println("hashNext:", hashNext)
		fmt.Println("nextMarker:", nextMarker, len(nextMarker))

		result, err = model.AddSyncLog(bucketId, hashNext, nextMarker, len(entries), sort)
		sort++
		if err != nil {
			fmt.Println("AddSyncLog error,", err)
			break
		} else {
			rowsaffected, _ := result.RowsAffected()
			lastinsertid, _ := result.LastInsertId()
			fmt.Println("AddSyncLog result:", rowsaffected, lastinsertid)
		}
		if hashNext {
			marker = nextMarker
		} else {
			fmt.Println("-----AsyncFiles End-----")
			break
		}
	}
}

func ListBucket(bucket, key, secret string) {
	mac := auth.New(key, secret)

	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Zone=&storage.ZoneHuabei
	bucketManager := storage.NewBucketManager(mac, &cfg)

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
