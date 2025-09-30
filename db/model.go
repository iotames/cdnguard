package db

type HttpRequest struct {
	Id, Ip, XForwardedFor, UserAgent, Referer, RequestUrl, Headers, RawUrl string
}

func AddRequest(areq HttpRequest, block bool, block_type int) error {
	var err error
	d := db
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
		// 如果是阻断请求，添加阻断的理由类别，记录到阻断请求表
		insertvals = append(insertvals, block_type)
		_, err = d.AddBlockRequest(insertvals...)
	} else {
		_, err = d.AddRequest(insertvals...)
	}
	return err
}
