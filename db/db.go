package db

import (
	"sync"

	"github.com/iotames/cdnguard/contract"
	"github.com/iotames/easydb"
	_ "github.com/lib/pq"
)

// DB结构体和方法，只给main,model调用
type DB struct {
	edb     *easydb.EasyDb
	sqlDir  contract.ISqlDir
	dsnConf *easydb.DsnConf
}

var once sync.Once
var db *DB

// DB结构体和方法，只给main,model调用
func GetDb(oncedb *DB) *DB {
	once.Do(func() {
		db = oncedb
	})
	if db == nil {
		panic("db is nil")
	}
	return db
}

// DB结构体和方法，只给main,model调用
func NewDb(driverName, dbHost, dbUser, dbPassword, dbName string, dbPort int) *DB {
	var err error
	cf := easydb.NewDsnConf(driverName, dbHost, dbUser, dbPassword, dbName, dbPort)
	d := easydb.NewEasyDbByConf(*cf)
	// 测试连接d
	if err = d.Ping(); err != nil {
		panic(err)
	}
	return &DB{d, nil, cf}
}

func (d *DB) SetSqlDir(sqldir contract.ISqlDir) {
	d.sqlDir = sqldir
}

func (d DB) CloseDb() error {
	return d.edb.CloseDb()
}
