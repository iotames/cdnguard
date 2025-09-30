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

	// IP白名单PASS
	okips := GetIpWhiteList()
	if slices.Contains(okips, areq.Ip) {
		success(ctx)
		log.Println("info:ip white list PASS:", areq.Ip)
		AddRequest(areq, false)
		return
	}

	// IP黑名单BLOCK
	blackips := GetIpBlackList()
	if slices.Contains(blackips, areq.Ip) {
		fail(ctx)
		log.Println("error:ip blacklist Block:", areq.Ip)
		AddRequest(areq, true)
		return
	}

	AddRequest(areq, false)
	success(ctx)
}

type CdnAuthRequest struct {
	Id, Ip, XForwardedFor, UserAgent, Referer, RequestUrl, Headers, RawUrl string
}

func AddRequest(areq CdnAuthRequest, block bool) error {
	var err error
	d := GetDB()
	var ua, referer any
	if areq.UserAgent == "" {
		ua = nil
	} else {
		ua = areq.UserAgent
	}
	if areq.Referer == "" {
		referer = nil
	} else {
		referer = areq.Referer
	}
	insertvals := []any{areq.Id, areq.Ip, areq.XForwardedFor, ua, referer, areq.RequestUrl, areq.Headers, areq.RawUrl}
	if block {
		_, err = d.AddBlockRequest(insertvals...)
	} else {
		_, err = d.AddRequest(insertvals...)
	}
	if err != nil {
		log.Println("error: AddRequest sqlresult Fail:", err)
	}
	return err
}
