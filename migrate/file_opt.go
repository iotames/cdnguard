package migrate

import (
	"github.com/iotames/cdnguard/cdnapi"
	"github.com/iotames/cdnguard/log"
	"github.com/iotames/cdnguard/model"
)

func GetMigrateFiles() (fileKeys []string, err error) {
	err = model.GetMigrateFiles(&fileKeys)
	if err != nil {
		return
	}
	log.Debug("FileMigrate.GetMigrateFiles", "fileKeys", fileKeys)
	return
}

// Migrate 迁移文件
func (m FileMigrate) Migrate(capi *cdnapi.CdnApi, addPreDir string) error {
	fileKeys, err := GetMigrateFiles()
	if err != nil {
		return err
	}
	return capi.MigrateFiles("copy", m.fromBucket, m.toBucket, fileKeys, func(fkey string, err error) {
		if err != nil {
			log.Error("FileMigrate.Migrate error", "fileKey", fkey, "err", err)
		} else {
			model.LogFileMigrate("copy", fkey, m.fromBucket, m.toBucket, addPreDir)
			log.Info("SUCCESS! FileMigrate.Migrate Copy Done", "fileKey", fkey)
		}
	}, addPreDir)
}

func GetDeleteFiles() (fileKeys []string, err error) {
	err = model.GetDeleteFiles(&fileKeys)
	if err != nil {
		return
	}
	log.Debug("FileMigrate.GetDeleteFiles", "fileKeys", fileKeys)
	return
}

// Delete 删除文件
func (m FileMigrate) Delete(capi *cdnapi.CdnApi) error {
	deleteFiles, err := GetDeleteFiles()
	if err != nil {
		return err
	}
	log.Debug("FileMigrate.Delete", "bucketName", m.fromBucket)
	return capi.DeleteFiles(m.fromBucket, deleteFiles, func(fkey string, err error) {
		if err != nil {
			log.Error("FileMigrate.DeleteFiles error", "fileKey", fkey, "err", err)
		} else {
			model.LogFileMigrate("delete", fkey, m.fromBucket, m.toBucket, "")
			log.Info("SUCCESS! FileMigrate.Delete Done", "bucketName", m.fromBucket, "fileKey", fkey)
		}
	})
}
