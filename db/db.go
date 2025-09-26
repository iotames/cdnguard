package db

import (
	"sync"

	"github.com/iotames/cdnguard/contract"
	"github.com/iotames/easydb"
	_ "github.com/lib/pq"
)

type DB struct {
	edb    *easydb.EasyDb
	sqlDir contract.ISqlDir
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

func (d *DB) SetSqlDir(sqldir contract.ISqlDir) {
	d.sqlDir = sqldir
}

func (d DB) CloseDb() error {
	return d.edb.CloseDb()
}
