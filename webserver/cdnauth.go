package webserver

import (
	"database/sql"
	"encoding/json"
	"log"

	"slices"

	"github.com/iotames/easyserver/httpsvr"
)

func cdnauth(ctx httpsvr.Context) {
	var err error
	var hdrb []byte
	q := ctx.Request.URL.Query()
	request_headers := ""
	hdrb, err = json.Marshal(ctx.Request.Header)
	if err == nil {
		request_headers = string(hdrb)
	}

	areq := CdnAuthRequest{
		Id:            q.Get("clientrequestid"),
		Ip:            q.Get("clientip"),
		XForwardedFor: q.Get("clientxforwardedfor"),
		UserAgent:     q.Get("clientua"),
		Referer:       q.Get("clientreferer"),
		RequestUrl:    q.Get("requesturl"),
		Headers:       request_headers,
		RawUrl:        ctx.Request.URL.String(),
	}

	// IP白名单PASS
	okips, _ := GetDB().GetIpWhiteList()
	if slices.Contains(okips, areq.Ip) {
		AddRequest(areq)
		success(ctx)
		return
	}
	// 添加到黑名单BLOCK

	AddRequest(areq)
	success(ctx)
}

type CdnAuthRequest struct {
	Id, Ip, XForwardedFor, UserAgent, Referer, RequestUrl, Headers, RawUrl string
}

func AddRequest(areq CdnAuthRequest) error {
	var err error
	var sqlresult sql.Result
	d := GetDB()
	sqlresult, err = d.AddRequest(
		areq.Id,
		areq.Ip,
		areq.XForwardedFor,
		areq.UserAgent,
		areq.Referer,
		areq.RequestUrl,
		areq.Headers,
		areq.RawUrl,
	)
	if err != nil {
		log.Println("sqlresult error:", err)
	} else {
		var n int64
		n, err = sqlresult.RowsAffected()
		log.Println("SUCCESS sqlresult:", n, err)
	}
	return err
}
