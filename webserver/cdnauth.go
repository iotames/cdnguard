package webserver

import (
	"encoding/json"
	"log"

	"github.com/iotames/cdnguard/guard"
	"github.com/iotames/cdnguard/model"
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

	hreq := model.HttpRequest{
		Id:            q.Get("clientrequestid"),
		Ip:            q.Get("clientip"),
		XForwardedFor: q.Get("clientxforwardedfor"),
		UserAgent:     q.Get("clientua"),
		Referer:       q.Get("clientreferer"),
		RequestUrl:    q.Get("requesturl"),
		Headers:       request_headers,
		RawUrl:        ctx.Request.URL.String(),
	}
	guard.GuardPass(hreq, func(pass bool) {
		if pass {
			success(ctx)
		} else {
			fail(ctx)
		}
	}, func(err error) {
		if err != nil {
			log.Printf("----error--cdnauth--GuardPass--SaveRequestError(%v)---", err)
		}
	})
}
