package qiniu

import (
	"strings"
)

// https://developer.qiniu.com/kodo/1250/batch

// BatchCopyFile 批量复制文件。保持文件key不变，从源bucket复制到目标bucket。
func (q QiniuCdn) BatchCopyFile(bucketSrc, bucketDest string, fileKeys []string, callback func(fkey string, err error), addPreDir string) {
	if strings.HasPrefix(addPreDir, `/`) || strings.HasSuffix(addPreDir, `/`) {
		panic("addPreDir must not start or end with /")
	}
	var err error
	for _, fileKey := range fileKeys {
		destFileKey := fileKey
		if addPreDir != "" {
			destFileKey = addPreDir + "/" + fileKey
		}
		err = q.bucketManager.Copy(bucketSrc, fileKey, bucketDest, destFileKey, false)
		callback(fileKey, err)
	}
}

// BatchDeleteFile 批量删除文件。
func (q QiniuCdn) BatchDeleteFile(bucketName string, fileKeys []string, callback func(fkey string, err error)) {
	var err error
	for _, fileKey := range fileKeys {
		err = q.bucketManager.Delete(bucketName, fileKey)
		callback(fileKey, err)
	}
}
