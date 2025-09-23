package cdnguard

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/iotames/miniutils"
)

type ScriptDir struct {
	DefaultDir, CustomDir string
	embedFS               embed.FS
}

var onesd *ScriptDir
var once sync.Once

func GetScriptDir(sd *ScriptDir) *ScriptDir {
	once.Do(func() {
		onesd = sd
	})
	return onesd
}

func NewScriptDir(customDir, defaultDir string, embedFs embed.FS) *ScriptDir {
	return &ScriptDir{DefaultDir: defaultDir, CustomDir: customDir, embedFS: embedFs}
}

// GetSQL 获取sql文本
// replaceList 字符串列表，依次替换SQL文本中的?占位符
// TODO 需要强调占位符与通配符的区别，比如%和_在LIKE子句中不是占位符，而是通配符，需要和参数化查询中的占位符区分开。
func (s ScriptDir) GetSQL(fpath string, replaceList ...string) (string, error) {
	sqlTxt, err := s.getSqlText(fpath)
	if err != nil {
		return "", err
	}
	for _, rerplaceStr := range replaceList {
		sqlTxt = strings.Replace(sqlTxt, "?", rerplaceStr, 1)
	}
	return sqlTxt, nil
}

// getSqlText 获取sql文本
// 优先从custom/sql自定义目录读取sql。如找不到SQL文件，则从默认的sql目录中读取。如再找不到文件，则从内嵌文件中读取。
func (s ScriptDir) getSqlText(fpath string) (sqlTxt string, err error) {
	defaultFilePath := filepath.Join(s.DefaultDir, fpath)
	customDirPath := filepath.Join(s.CustomDir, fpath)
	// 优先从custom/sql自定义目录读取sql。如找不到SQL文件，则从默认的sql目录中读取。
	sqlTxt, err = GetTextByFilePath(defaultFilePath, customDirPath)
	if sqlTxt == "" {
		// 找不到文件，从内嵌的文件中读取sql文件
		var sqlBytes []byte
		sqlBytes, err = s.embedFS.ReadFile(fpath)
		if err != nil {
			return
		}
		sqlTxt = string(sqlBytes)
	}
	return
}

func getExistFile(defaultFilePath string, customFilePath string) string {
	if miniutils.IsPathExists(customFilePath) {
		return customFilePath
	}
	if miniutils.IsPathExists(defaultFilePath) {
		return defaultFilePath
	}
	return ""
}

func GetTextByFilePath(defaultFilePath string, customFilePath string) (content string, err error) {
	fpath := getExistFile(defaultFilePath, customFilePath)
	if fpath == "" {
		return "", fmt.Errorf("file not found: %s", defaultFilePath)
	}
	return miniutils.ReadFileToString(fpath)
}
