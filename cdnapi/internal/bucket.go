package internal

import (
	"fmt"
)

// GetBucketId 仅支持单个CDN服务商的数据处理。当前仅接入了七牛云的API。
// 因为数据列表仅有bucket_id标记同一个CDN服务商的数据桶ID，没办法区分不同服务商。
func GetBucketId(bucketName string, buckets []string) (bucketId int, err error) {
	bucketId = -1
	for i, bname := range buckets {
		if bname == bucketName {
			bucketId = i
			break
		}
	}
	if bucketId == -1 {
		return -1, fmt.Errorf("the cmd ArgPair[bucketname=%s] value not in BUCKET_NAME_LIST(%v). Please look the file: .env", bucketName, buckets)
	}
	return bucketId, nil
}
