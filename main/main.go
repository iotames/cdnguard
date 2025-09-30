package main

import (
	"github.com/iotames/cdnguard/db"
)

func main() {
	// 关闭整个d连接池
	d := db.GetDb(nil)
	defer d.CloseDb()
	if extCmdRun() {
		return
	}
	runserver()
}

func init() {
	parseRunArgs()
	dbinit()
}
