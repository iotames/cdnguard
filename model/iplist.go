package model

import (
	"time"
)

// 黑名单类别为手动添加
const BLACK_TYPE_MAN = 0

// 黑名单类别为规则拦截
const BLACK_TYPE_RULE = 1

// 黑名单类别为请求过于频繁
const BLACK_TYPE_OVER_LIMIT = 2

// IP请求统计
type IpRequest struct {
	Ip           string `db:"ip"`
	RequestCount int    `db:"request_count"`
}

// AddIpToBlackList 把当IP添加到IP黑名单
func (i IpRequest) AddIpToBlackList(title string, black_type int) {
	getDB().AddIpToBlackList(i.Ip, title, black_type)
}

// GetTopRequestIpToday 获取今天指定时间范围内，请求数最高的IP列表
func GetTopRequestIpToday(hourBegin, hourBefor int) (ipreqs []IpRequest, err error) {
	// startAt := `CURRENT_DATE + INTERVAL '1 hours'`
	// endAt := `CURRENT_DATE + INTERVAL '5 hours'`
	now := time.Now()
	startAt := time.Date(now.Year(), now.Month(), now.Day(), hourBegin, 0, 0, 0, now.Location())
	endAt := time.Date(now.Year(), now.Month(), now.Day(), hourBefor, 0, 0, 0, now.Location())
	err = getDB().GetTopRequestIps(&ipreqs, startAt, endAt)
	return ipreqs, err
}

// var ipWhiteList []string

// GetIpWhiteList 获取IP白名单
// TODO 可以N小时更新一次IP白名单
func GetIpWhiteList() []string {
	// once.Do(func() {
	ipWhiteList, _ := getDB().GetIpWhiteList()
	// })
	return ipWhiteList
}

// GetIpBlackList 获取IP黑名单
// TODO 可以15分钟更新一次IP黑名单，可以不用每次网络请求都执行SQL查询来获取IP黑名单列表
func GetIpBlackList() []string {
	ipBlackList, _ := getDB().GetIpBlackList()
	return ipBlackList
}
