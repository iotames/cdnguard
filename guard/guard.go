package guard

import (
	"log"
	"slices"

	"github.com/iotames/cdnguard/model"
)

func GuardPass(hreq model.HttpRequest, callback func(pass bool)) error {
	// PASS IP白名单的请求，直接通过
	okips := model.GetIpWhiteList()
	if slices.Contains(okips, hreq.Ip) {
		callback(true)
		log.Println("info:ip whitelist PASS:", hreq.Ip)
		// log.Println("error: AddRequest sqlresult Fail:", err)
		return model.AddRequestPass(hreq)
	}

	// 前置过滤器，放在白名单之后，黑名单之前
	// BLOCK 过滤器拦截异常请求
	if block_type, isBlock := FilterBlock(hreq); isBlock {
		callback(false)
		log.Println("block: URL ends with .php:", hreq.RequestUrl)
		return model.AddRequestBlock(hreq, block_type)
	}

	// BLOCK 拦截IP黑名单
	blackips := model.GetIpBlackList()
	if slices.Contains(blackips, hreq.Ip) {
		callback(false)
		log.Println("error:ip blacklist Block:", hreq.Ip)
		return model.AddRequestBlock(hreq, model.BLOCK_TYPE_BLACK)
	}
	// PASS 默认通过
	callback(true)
	return model.AddRequestPass(hreq)
}
