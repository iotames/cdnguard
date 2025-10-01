package cdnguard

import (
	"fmt"
	"log"

	"github.com/iotames/cdnguard/model"
)

func AddBlackIpList(requestLimit int) error {
	// 获取当天1点后，5点前，网络请求的referer为空，且请求最多的IP
	ips, err := model.GetTopRequestIpToday(1, 5)
	if err != nil {
		return err
	}
	// 在指定时间段内，请求数超过1600，且，则封IP
	blackTitle := fmt.Sprintf("请求数超过%d，且referer为空", requestLimit)
	for i, ip := range ips {
		if ip.RequestCount > requestLimit {
			ip.AddIpToBlackList(blackTitle, model.BLACK_TYPE_OVER_LIMIT)
			log.Printf("---AddBlackIpList--i(%d)--IP(%s)---count(%d)", i, ip.Ip, ip.RequestCount)
		}
	}
	// 10分钟或者1小时，更新一次IP黑名单。go func(){}添加到黑名单
	// TODO 请求头没有accept-language可能是爬虫
	return nil
}
