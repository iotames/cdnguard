package cdnguard

import (
	"fmt"
	"log"
	"time"

	"github.com/iotames/cdnguard/db"
)

type IpRequest struct {
	Ip           string `db:"ip"`
	RequestCount int    `db:"request_count"`
}

func AddBlackIpList() error {
	// 获取当天1点后，5点前，网络请求的referer为空，且请求最多的IP
	ips, err := GetTopRequestIpToday(1, 5)
	if err != nil {
		return err
	}
	d := db.GetDb(nil)
	// 在指定时间段内，请求数超过1600，且，则封IP
	max := 66
	blackTitle := fmt.Sprintf("请求数超过%d，且referer为空", max)
	for i, ip := range ips {
		if ip.RequestCount > max {
			d.AddIpToBlackList(ip.Ip, blackTitle, 2)
			log.Printf("---AddBlackIpList--i(%d)--IP(%s)---count(%d)", i, ip.Ip, ip.RequestCount)
		}
	}
	// 10分钟或者1小时，更新一次IP黑名单。go func(){}添加到黑名单
	// TODO 请求头没有accept-language可能是爬虫
	return nil
}

// GetTopRequestIpToday 获取今天指定时间范围内，请求数最高的IP列表
func GetTopRequestIpToday(hourBegin, hourBefor int) (ipreqs []IpRequest, err error) {
	d := db.GetDb(nil)
	now := time.Now()
	startAt := time.Date(now.Year(), now.Month(), now.Day(), hourBegin, 0, 0, 0, now.Location())
	endAt := time.Date(now.Year(), now.Month(), now.Day(), hourBefor, 0, 0, 0, now.Location())
	// startAt := `CURRENT_DATE + INTERVAL '1 hours'`
	// endAt := `CURRENT_DATE + INTERVAL '5 hours'`
	err = d.GetTopRequestIps(&ipreqs, startAt, endAt)
	return ipreqs, err
}
