package webserver

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
)

func cdnauth(ctx httpsvr.Context) {
	var err error
	var hdrb []byte
	var sqlresult sql.Result
	d := GetDB()
	q := ctx.Request.URL.Query()
	request_headers := ""
	hdrb, err = json.Marshal(ctx.Request.Header)
	if err == nil {
		request_headers = string(hdrb)
	}
	sqlresult, err = d.AddRequest(
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
