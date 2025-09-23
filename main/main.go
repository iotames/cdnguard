package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/iotames/cdnguard"
	cdnsql "github.com/iotames/cdnguard/main/sql"
	"github.com/iotames/easyconf"
	"github.com/iotames/easydb"
	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
	_ "github.com/lib/pq"
)

const DEFALUT_SQL_DIR = "sql"

var SqlDir string
var DbDriverName, DbHost, DbUser, DbPassword, DbName string
var DbPort, WebPort int
var d *easydb.EasyDb

func main() {
	var err error
	var hdrb []byte
	var sqlresult sql.Result
	// 关闭整个d连接池
	defer d.CloseDb()
	s := httpsvr.NewEasyServer(fmt.Sprintf(":%d", WebPort))
	s.AddMiddleware(httpsvr.NewMiddleCORS("*"))
	s.AddHandler("GET", "/cdnauth", func(ctx httpsvr.Context) {
		sql := `INSERT INTO qiniu_cdnauth_requests (request_id,client_ip,x_forwarded_for,user_agent,http_referer,request_url,request_headers,raw_url) VALUES ($1,$2,$3,$4,$5,$6,$7, $8)`
		q := ctx.Request.URL.Query()
		request_headers := ""
		hdrb, err = json.Marshal(ctx.Request.Header)
		if err == nil {
			request_headers = string(hdrb)
		}
		sqlresult, err = d.Exec(sql,
			q.Get("clientrequestid"),
			q.Get("clientip"),
			q.Get("clientxforwardedfor"),
			q.Get("clientua"),
			q.Get("clientreferer"),
			q.Get("requesturl"),
			request_headers,
			ctx.Request.URL.String(),
		)
		if err != nil {
			log.Println("sqlresult error:", err)
		} else {
			var n int64
			n, err = sqlresult.RowsAffected()
			log.Println("SUCCESS sqlresult:", n, err)
		}
		ctx.Writer.Write(response.NewApiDataOk("success").Bytes())
	})
	s.ListenAndServe()
}

func init() {
	var err error
	cf := easyconf.NewConf()
	cf.StringVar(&DbDriverName, "DB_DRIVER_NAME", "postgres", "数据库驱动名称")
	cf.StringVar(&DbHost, "DB_HOST", "172.16.160.12", "数据库主机地址")
	cf.StringVar(&DbUser, "DB_USER", "postgres", "数据库用户名")
	cf.StringVar(&DbPassword, "DB_PASSWORD", "postgres", "数据库密码")
	cf.StringVar(&DbName, "DB_NAME", "postgres", "数据库名称")
	cf.IntVar(&DbPort, "DB_PORT", 5432, "数据库端口")
	cf.IntVar(&WebPort, "WEB_PORT", 1212, "web服务端口")
	cf.StringVar(&SqlDir, "SQL_DIR", DEFALUT_SQL_DIR, "sql文件目录")
	cf.Parse()

	sqldir := cdnguard.NewScriptDir(SqlDir, DEFALUT_SQL_DIR, cdnsql.GetSqlFs())
	sqlCreateTable, err := sqldir.GetSQL("tables_init.sql")
	if err != nil {
		panic(err)
	}

	d = easydb.NewEasyDb(DbDriverName, DbHost, DbUser, DbPassword, DbName, DbPort)
	// 测试连接d
	if err = d.Ping(); err != nil {
		log.Fatal(err)
		panic(err)
	}
	_, err = d.Exec(sqlCreateTable)
	if err != nil {
		panic(err)
	}
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
