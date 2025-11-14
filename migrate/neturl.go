package migrate

import (
	"strings"
)

// GetDomainByUrl. 获取url网址的域名。
// the arg url startwith http, //, / ; return like: "www.baidu.com", "baidu.com", ""
func GetDomainByUrl(url string) string {
	urlS := strings.Split(url, "/")

	if (strings.HasPrefix(url, "http") || strings.HasPrefix(url, "//")) && len(urlS) > 2 {
		return urlS[2]
	}
	if strings.HasPrefix(url, "/") && len(urlS) > 1 {
		return urlS[1]
	}
	return ""
}
