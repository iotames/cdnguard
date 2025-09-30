package cdnguard

import (
	"log"

	"github.com/iotames/cdnguard/db"
)

type IpRequest struct {
	Ip           string `db:"ip"`
	RequestCount int    `db:"request_count"`
}

func AddBlackIpList() error {
	d := db.GetDb(nil)
	var ips []IpRequest
	startAt := `CURRENT_DATE + INTERVAL '1 hours'`
	endAt := `CURRENT_DATE + INTERVAL '5 hours'`
	err := d.GetTopRequestIps(&ips, startAt, endAt)
	if err != nil {
		return err
	}
	for _, ip := range ips {
		log.Println("---------RequestIpBlackList---------", ip, ip, ip.RequestCount)
	}
	// 10分钟或者1小时，更新一次IP黑名单。go func(){}添加到黑名单
	// TODO 请求头没有accept-language可能是爬虫
	return nil
}
