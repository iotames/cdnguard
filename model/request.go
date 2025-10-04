package model

// 拦截请求类别为：IP黑名单拦截
const BLOCK_TYPE_BLACK = 0

// 拦截请求类别为：漏洞扫描
const BLOCK_TYPE_SCAN_VUL = 1

// 拦截请求类别为：网络爬虫
const BLOCK_TYPE_SPIDER = 2

// 拦截请求类别为：异常的UserAgent
const BLOCK_TYPE_USERAGENT = 3

type HttpRequest struct {
	Id, Ip, XForwardedFor, UserAgent, Referer, RequestUrl, Headers, RawUrl string
}

// addRequest 添加网络请求记录
// block: 是否阻断拦截请求。
// block_type: 拦截阻断的理由类别。0=IP黑名单拦截，1=漏洞扫描拦截，2=网络爬虫
func addRequest(areq HttpRequest, block bool, block_type int) error {
	var err error
	d := getDB()
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

// AddRequestPass 记录通过的请求
func AddRequestPass(areq HttpRequest) error {
	return addRequest(areq, false, 0)
}

// AddRequestBlock 记录被拦截阻断的请求
// block_type: 拦截阻断的理由类别。0=IP黑名单拦截，1=规则拦截，2=网络爬虫
func AddRequestBlock(areq HttpRequest, block_type int) error {
	return addRequest(areq, true, block_type)
}
