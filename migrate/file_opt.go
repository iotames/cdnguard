package migrate

import (
	"github.com/iotames/cdnguard/cdnapi"
	"github.com/iotames/cdnguard/log"
	"github.com/iotames/cdnguard/model"
)

// 使用事务保证操作的原子性
// https://developer.qiniu.com/kodo/1250/batch
func (m FileMigrate) Migrate(capi *cdnapi.CdnApi) error {

	migrateFiles := []model.MigrateFile{}
	err := model.GetMigrateFiles(&migrateFiles)
	if err != nil {
		return err
	}
	fileKeys := []string{}
	for _, migrateFile := range migrateFiles {
		log.Debug("FileMigrate.Migrate append fileKeys", "fileKey", migrateFile.FileKey)
		fileKeys = append(fileKeys, migrateFile.FileKey)
	}

	return capi.MigrateFiles("copy", m.fromBucket, m.toBucket, fileKeys, func(fkey string, err error) {
		if err != nil {
			log.Error("FileMigrate.Migrate error", "fileKey", fkey, "err", err)
		} else {
			model.LogFileMigrate("copy", fkey, m.fromBucket, m.toBucket)
			log.Info("SUCCESS! FileMigrate.Migrate Copy Done", "fileKey", fkey)
		}
	})
}

// 删除文件
func (m FileMigrate) Delete(capi *cdnapi.CdnApi) error {
	var deleteFiles []string
	err := model.GetDeleteFiles(&deleteFiles)
	if err != nil {
		return err
	}
	log.Debug("FileMigrate.Delete", "bucketName", m.fromBucket, "deleteFiles", deleteFiles)
	return capi.DeleteFiles(m.fromBucket, deleteFiles, func(fkey string, err error) {
		if err != nil {
			log.Error("FileMigrate.DeleteFiles error", "fileKey", fkey, "err", err)
		} else {
			model.LogFileMigrate("delete", fkey, m.fromBucket, m.toBucket)
			log.Info("SUCCESS! FileMigrate.Delete Done", "bucketName", m.fromBucket, "fileKey", fkey)
		}
	})
}
