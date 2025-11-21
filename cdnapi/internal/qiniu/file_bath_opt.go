package qiniu

// BatchCopyFile 批量复制文件。保持文件key不变，从源bucket复制到目标bucket。
func (q QiniuCdn) BatchCopyFile(bucketSrc, bucketDest string, fileKeys []string, callback func(fkey string, err error)) {
	var err error
	for _, fileKey := range fileKeys {
		err = q.bucketManager.Copy(bucketSrc, fileKey, bucketDest, fileKey, false)
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
