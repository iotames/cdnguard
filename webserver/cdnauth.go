package webserver

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/iotames/cdnguard/db"
	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
	_ "github.com/lib/pq"
)

func cdnauth(ctx httpsvr.Context) {
	var err error
	var hdrb []byte
	var sqlresult sql.Result
	d := db.GetDb(nil)

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
}
