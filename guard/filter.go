package guard

import (
	"github.com/iotames/cdnguard/model"
)

func FilterBlock(hreq model.HttpRequest) (int, bool) {
	// 前置过滤器由多个过滤规则组成，每个过滤规则，包含：过滤规则，规则编码，规则标题，是否加入黑名单
	// 如果URL以.php结尾，则拦截。未避免误伤自己的测试人员，放宽到2分钟10个请求再拉黑名单
	if len(hreq.RequestUrl) >= 4 && hreq.RequestUrl[len(hreq.RequestUrl)-4:] == ".php" {
		// 拦截不合法的请求地址
		// TODO IP超过10次此类请求再拉黑
		return model.BLOCK_TYPE_SCAN_VUL, true
	}
	ualimiter := NewUserAgentLimiter()
	if ualimiter.IsSpider(hreq.UserAgent) {
		return model.BLOCK_TYPE_SPIDER, true
	}
	if !ualimiter.Allow(hreq.UserAgent) {
		return model.BLOCK_TYPE_USERAGENT, true
	}
	return 0, false
}
