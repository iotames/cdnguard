package migrate

import (
	"github.com/iotames/cdnguard/cdnapi"
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
		fileKeys = append(fileKeys, migrateFile.FileKey)
	}
	return capi.MigrateFiles(m.fromBucket, m.toBucket, fileKeys, func(fkey string, err error) {
		// TODO 变更迁移状态，添加文件操作记录
	})
}
