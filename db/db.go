package db

import (
	"database/sql"
	"sync"
	"time"

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

func (d *DB) GetSqlDB() *sql.DB {
	return d.edb.GetSqlDB()
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
	// 设置合理的连接池参数
	// SHOW max_connections; // 检查 PostgreSQL 最大连接数
	// SELECT * FROM pg_stat_activity; 检查是否有其他应用占用连接
	// SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE state = 'idle' AND query_start < NOW() - INTERVAL '10 minutes'; // 杀死空闲或长时间运行的连接
	// db.SetMaxOpenConns(20)  // 最大打开连接数
	// db.SetMaxIdleConns(5)   // 最大空闲连接数
	// db.SetConnMaxLifetime(time.Hour) // 连接最大生命周期
	dbb := d.GetSqlDB()
	dbb.SetMaxOpenConns(2000)
	dbb.SetConnMaxLifetime(time.Minute * 10)
	dd := &DB{d, nil, cf}
	return dd
}

func (d *DB) SetSqlDir(sqldir contract.ISqlDir) {
	d.sqlDir = sqldir
}

func (d *DB) Stats() sql.DBStats {
	// 获取连接统计信息
	dbb := d.edb.GetSqlDB()
	return dbb.Stats()
}

func (d DB) CloseDb() error {
	return d.edb.CloseDb()
}
