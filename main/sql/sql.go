package sql

import (
	"embed"
)

//go:embed *.sql
var sqlFS embed.FS

func GetSqlFs() embed.FS {
	return sqlFS
}

// func lsDir() []string {
// 	entries, err := sqlFS.ReadDir(".")
// 	if err != nil {
// 		panic(err)
// 	}
// 	var filenames []string
// 	for _, entry := range entries {
// 		filenames = append(filenames, entry.Name())
// 		if entry.IsDir() {
// 			// 读取子目录中的文件
// 			subEntries, err := sqlFS.ReadDir(entry.Name())
// 			if err != nil {
// 				panic(err)
// 			}
// 			// 将子目录中的文件名添加到列表中
// 			for _, subEntry := range subEntries {
// 				filenames = append(filenames, entry.Name()+"/"+subEntry.Name())
// 			}
// 		}
// 	}
// 	return filenames
// }
