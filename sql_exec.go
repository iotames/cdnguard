package cdnguard

import "time"

// Prune. 定期对冗长的requests和block_requests数据表进行清理
func Prune(d time.Duration) error {
	// TODO
	// 请求头没有accept-language可能是爬虫
	return nil
}
