package main

import (
	"github.com/iotames/cdnguard/db"
)

var Debug bool

func main() {
	// 关闭整个d连接池
	d := db.GetDb(nil)
	defer d.CloseDb()
	if Debug {
		debug()
		return
	}
	runserver()
}

func init() {

	parseArgs()
	dbinit()
}
