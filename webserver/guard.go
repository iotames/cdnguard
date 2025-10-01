package webserver

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
	if FilterBlock(hreq) {
		callback(false)
		log.Println("block: URL ends with .php:", hreq.RequestUrl)
		return model.AddRequestBlock(hreq, model.BLOCK_TYPE_RULE)
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

func FilterBlock(hreq model.HttpRequest) bool {
	// 前置过滤器由多个过滤规则组成，每个过滤规则，包含：过滤规则，规则编码，规则标题，是否加入黑名单
	// 如果URL以.php结尾，则拦截。未避免误伤自己的测试人员，放宽到2分钟10个请求再拉黑名单
	if len(hreq.RequestUrl) >= 4 && hreq.RequestUrl[len(hreq.RequestUrl)-4:] == ".php" {
		// 拦截不合法的请求地址
		// TODO IP超过10次此类请求再拉黑
		return true
	}
	return false
}

// 102.90.101.150 Wildto/1 CFNetwork/3826.600.41 Darwin/24.6.0
// 216.73.216.38 Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; ClaudeBot/1.0; +claudebot@anthropic.com)
// 66.249.77.42 Googlebot-Image/1.0
