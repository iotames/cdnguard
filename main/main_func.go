package main

import (
	"flag"
	"log"

	"github.com/iotames/cdnguard"
	"github.com/iotames/cdnguard/db"
	cdnsql "github.com/iotames/cdnguard/main/sql"
	"github.com/iotames/cdnguard/webserver"
	"github.com/iotames/easyconf"
)

const DEFALUT_SQL_DIR = "sql"

var gdb *db.DB
var SqlDir string
var DbDriverName, DbHost, DbUser, DbPassword, DbName string
var DbPort, WebPort int
var Prune bool

func dbinit() {
	gdb = db.NewDb(DbDriverName, DbHost, DbUser, DbPassword, DbName, DbPort)
	sqldir := cdnguard.NewScriptDir(cdnsql.GetSqlFs(), SqlDir, DEFALUT_SQL_DIR)
	gdb.SetSqlDir(sqldir)
	_, err := gdb.CreateTables()
	if err != nil {
		panic(err)
	}
	log.Println("数据库初始化完成")
	// 调用GetDb方法，传入gdb。以便其他模块使用GetDb(nil)获取全局单例
	db.GetDb(gdb)
}

func runserver() {
	var err error
	s := webserver.NewWebServer(WebPort)
	if err = s.ListenAndServe(); err != nil {
		panic(err)
	}
}

func parseConf() {
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

func parseCmd() {
	flag.BoolVar(&Debug, "debug", false, "debug mode")
	flag.BoolVar(&Prune, "prune", false, "prune db")
	flag.Parse()
}

func parseRunArgs() {
	parseConf()
	parseCmd()
}
