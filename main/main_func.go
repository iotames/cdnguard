package main

import (
	"flag"
	"log"

	"github.com/iotames/cdnguard"
	"github.com/iotames/cdnguard/db"
	cdnsql "github.com/iotames/cdnguard/main/sql"
	"github.com/iotames/cdnguard/migrate"
	"github.com/iotames/cdnguard/webserver"
	"github.com/iotames/easyconf"
)

const DEFALUT_SQL_DIR = "sql"

var gdb *db.DB
var SqlDir string
var DbDriverName, DbHost, DbUser, DbPassword, DbName string
var CdnName, BucketName, QiniuAccessKey, QiniuSecretKey string
var LastCursorMarker string
var DbPort, WebPort, RequestLimit int
var Debug, Prune, AddBlackIps, SyncBucketFiles, ShowBucketFiles, StatisEveryDay, DbStats bool
var FileMigrate, FileDelete bool
var BucketNameList []string
var migrateFromHost, migrateToHost, migrateReferer, fromBucket, toBucket string

func dbinit() {
	gdb = db.NewDb(DbDriverName, DbHost, DbUser, DbPassword, DbName, DbPort)
	sqldir := cdnguard.NewScriptDir(cdnsql.GetSqlFs(), SqlDir, DEFALUT_SQL_DIR)
	gdb.SetSqlDir(sqldir)
	_, err := gdb.CreateTables()
	if err != nil {
		panic(err)
	}
	log.Println("数据库初始化完成", DbHost, DbPort, DbName)
	// 调用GetDb方法，传入gdb。以便其他模块使用GetDb(nil)获取全局单例
	db.GetDb(gdb)
}

func runserver() {
	var err error
	s := webserver.NewWebServer(WebPort)
	s.AddMiddleHead(migrate.NewFileMigrate(migrateFromHost, migrateToHost, migrateReferer, fromBucket, toBucket))
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
	cf.IntVar(&RequestLimit, "REQUEST_LIMIT", 1600, "单位时间内IP最大请求数限制", "可能是凌晨1-5点")
	cf.StringVar(&SqlDir, "SQL_DIR", DEFALUT_SQL_DIR, "sql文件目录")
	cf.StringVar(&QiniuAccessKey, "QINIU_ACCESS_KEY", "", "七牛AccessKey")
	cf.StringVar(&QiniuSecretKey, "QINIU_SECRET_KEY", "", "七牛SecretKey")
	cf.StringListVar(&BucketNameList, "BUCKET_NAME_LIST", []string{"bucket123"}, "可用的空间名列表。逗号分隔。固定好顺序不要变，有需要可往后添加。因为数据表存储的bucket_id和顺序有关。")
	cf.StringVar(&migrateFromHost, "MIGRATE_FROM_HOST", "", "迁移源主机")
	cf.StringVar(&migrateToHost, "MIGRATE_TO_HOST", "", "迁移目标主机")
	cf.StringVar(&migrateReferer, "MIGRATE_REFERER", "", "通过Referer迁移")
	cf.StringVar(&fromBucket, "MIGRATE_FROM_BUCKET", "", "迁移源空间")
	cf.StringVar(&toBucket, "MIGRATE_TO_BUCKET", "", "迁移目标空间")
	if err := cf.Parse(false); err != nil {
		log.Fatal(err)
	}
}

func parseCmd() {
	flag.BoolVar(&Debug, "debug", false, "debug mode")
	flag.BoolVar(&Prune, "prune", false, "prune db")
	flag.BoolVar(&AddBlackIps, "addblackips", false, "Add IP list to Black IP List")
	flag.BoolVar(&SyncBucketFiles, "syncbucketfiles", false, "sync bucket files")
	flag.BoolVar(&ShowBucketFiles, "showbucketfiles", false, "show bucket files")
	flag.BoolVar(&StatisEveryDay, "statiseveryday", false, "statis request data every day")
	flag.BoolVar(&DbStats, "dbstats", false, "db stats")
	flag.StringVar(&CdnName, "cdnname", "qiniu", "cdn name")
	flag.StringVar(&BucketName, "bucketname", "", "bucket name")
	flag.StringVar(&LastCursorMarker, "lastcursor", "", "上一次列表的最后一条数据标记")
	flag.BoolVar(&FileMigrate, "filemigrate", false, "file migrate")
	flag.BoolVar(&FileDelete, "filedelete", false, "file delete")
	flag.Parse()
}

func parseRunArgs() {
	parseConf()
	parseCmd()
}
