package db

import (
	"sync"

	"github.com/iotames/easydb"
)

type ISqlDir interface {
	GetSQL(fpath string, replaceList ...string) (string, error)
}

type DB struct {
	edb    *easydb.EasyDb
	sqlDir ISqlDir
}

var once sync.Once
var db *DB

func GetDb(oncedb *DB) *DB {
	once.Do(func() {
		db = oncedb
	})
	if db == nil {
		panic("db is nil")
	}
	return db
}

func NewDb(driverName, dbHost, dbUser, dbPassword, dbName string, dbPort int) *DB {
	var err error
	d := easydb.NewEasyDb(driverName, dbHost, dbUser, dbPassword, dbName, dbPort)
	// 测试连接d
	if err = d.Ping(); err != nil {
		panic(err)
	}
	return &DB{d, nil}
}

func (d *DB) SetSqlDir(sqldir ISqlDir) {
	d.sqlDir = sqldir
}

func (d DB) CloseDb() error {
	return d.edb.CloseDb()
}

// SELECT client_ip, COUNT(*) AS request_count
// FROM public.qiniu_cdnauth_requests
// WHERE created_at >= NOW() - INTERVAL '10 minutes'
// --WHERE created_at >= NOW() - INTERVAL '1 hour'
// GROUP BY client_ip
// ORDER BY request_count DESC
// LIMIT 10;
// -- 最近10分钟内网络请求最频繁的前10名IP
// -- 最近1小时内网络请求最频繁的前10名IP
