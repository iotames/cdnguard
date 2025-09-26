package main

import (
	"log"

	"github.com/iotames/cdnguard"
	"github.com/iotames/cdnguard/db"
	cdnsql "github.com/iotames/cdnguard/main/sql"
	"github.com/iotames/cdnguard/webserver"
	"github.com/iotames/easyconf"
	_ "github.com/lib/pq"
)

const DEFALUT_SQL_DIR = "sql"

var SqlDir string
var DbDriverName, DbHost, DbUser, DbPassword, DbName string
var DbPort, WebPort int

func dbinit() {
	d := db.NewDb(DbDriverName, DbHost, DbUser, DbPassword, DbName, DbPort)
	sqldir := cdnguard.NewScriptDir(cdnsql.GetSqlFs(), SqlDir, DEFALUT_SQL_DIR)
	d.SetSqlDir(sqldir)
	d.CreateTables()
	db.GetDb(d)
}

func runserver() {
	var err error
	s := webserver.NewWebServer(WebPort)
	if err = s.ListenAndServe(); err != nil {
		panic(err)
	}
}

func parseArgs() {
	cf := easyconf.NewConf()
	cf.StringVar(&DbDriverName, "DB_DRIVER_NAME", "postgres", "数据库驱动名称")
	cf.StringVar(&DbHost, "DB_HOST", "127.0.0.1", "数据库主机地址")
	cf.StringVar(&DbUser, "DB_USER", "postgres", "数据库用户名")
	cf.StringVar(&DbPassword, "DB_PASSWORD", "postgres", "数据库密码")
	cf.StringVar(&DbName, "DB_NAME", "postgres", "数据库名称")
	cf.IntVar(&DbPort, "DB_PORT", 5432, "数据库端口")
	cf.IntVar(&WebPort, "WEB_PORT", 1212, "web服务端口")
	cf.StringVar(&SqlDir, "SQL_DIR", DEFALUT_SQL_DIR, "sql文件目录")
	if err := cf.Parse(); err != nil {
		log.Fatal(err)
	}
}
