package cdnguard

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type ScriptDir struct {
	DefaultDir, CustomDir string
	embedFS               embed.FS
}

var onesd *ScriptDir
var once sync.Once

// GetScriptDir 获取脚本目录单例
func GetScriptDir(sd *ScriptDir) *ScriptDir {
	once.Do(func() {
		onesd = sd
	})
	if onesd == nil {
		panic("ScriptDir is nil")
	}
	return onesd
}

// NewScriptDir 初始化一个程序运行时所需的外部脚本文件目录。
func NewScriptDir(customDir, defaultDir string, embedFs embed.FS) *ScriptDir {
	return &ScriptDir{DefaultDir: defaultDir, CustomDir: customDir, embedFS: embedFs}
}

// GetSQL 获取sql文本
// replaceList 字符串列表，依次替换SQL文本中的?占位符
// TODO 需要强调占位符与通配符的区别，比如%和_在LIKE子句中不是占位符，而是通配符，需要和参数化查询中的占位符区分开。
func (s ScriptDir) GetSQL(fpath string, replaceList ...string) (string, error) {
	sqlTxt, err := s.GetScriptText(fpath)
	if err != nil {
		return "", err
	}
	for _, rerplaceStr := range replaceList {
		sqlTxt = strings.Replace(sqlTxt, "?", rerplaceStr, 1)
	}
	return sqlTxt, nil
}

func (s ScriptDir) OkDir(d string) error {
	info, err := os.Stat(d)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("%s is not a directory", d)
		}
		return err
	}
	if os.IsNotExist(err) {
		return fmt.Errorf("dir(%s) does not exist", d)
	}
	return err
}

func (s ScriptDir) OkNormalFile(d string) error {
	info, err := os.Stat(d)
	if err == nil {
		if info.IsDir() {
			return fmt.Errorf("%s is not a directory", d)
		}
		return err
	}
	if os.IsNotExist(err) {
		return fmt.Errorf("file(%s) does not exist", d)
	}
	return err
}

// GetScriptText 获取脚本文件的纯文本内容
// 优先从CustomDir自定义目录查找文件。如找不到，则从DefaultDir默认目录查找。如再找不到文件，则从内嵌文件中读取。
func (s ScriptDir) GetScriptText(fpath string) (stxt string, err error) {
	var b []byte
	customDirPath := filepath.Join(s.CustomDir, fpath)
	defaultFilePath := filepath.Join(s.DefaultDir, fpath)
	realfpath := s.getExistFile(customDirPath, defaultFilePath)
	if realfpath != "" {
		// 找到目标文件
		b, err = os.ReadFile(realfpath)
		stxt = string(b)
	}
	if stxt != "" && err == nil {
		// 文件内容不为空，读取没有错误
		return stxt, err
	}
	// 找不到文件，获取读取了文件内容为空。从内嵌的文件中读取sql文件

	b, err = s.embedFS.ReadFile(fpath)
	if err != nil {
		return stxt, err
	}
	stxt = string(b)
	return stxt, err
}

func (s ScriptDir) getExistFile(customFilePath string, defaultFilePath string) string {
	if s.OkNormalFile(customFilePath) == nil {
		return customFilePath
	}
	if s.OkNormalFile(defaultFilePath) == nil {
		return defaultFilePath
	}
	return ""
}
