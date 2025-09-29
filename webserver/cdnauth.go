package webserver

import (
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
	// TODO 请求头没有accept-language可能是爬虫

	// IP白名单PASS
	okips, _ := GetDB().GetIpWhiteList()
	if slices.Contains(okips, areq.Ip) {
		AddRequest(areq)
		success(ctx)
		return
	}

	// IP黑名单BLOCK
	okips, _ = GetDB().GetIpBlackList()
	if slices.Contains(okips, areq.Ip) {
		fail(ctx)
		log.Println("error:ip blacklist Block:", areq.Ip)
		// TODO 添加到block_requests
		return
	}

	// 10分钟或者1小时，更新一次IP黑名单。go func(){}添加到黑名单
	AddRequest(areq)
	success(ctx)
}

type CdnAuthRequest struct {
	Id, Ip, XForwardedFor, UserAgent, Referer, RequestUrl, Headers, RawUrl string
}

func AddRequest(areq CdnAuthRequest) error {
	var err error
	d := GetDB()
	_, err = d.AddRequest(
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
		log.Println("error: AddRequest sqlresult Fail:", err)
	}
	return err
}
