package main

import (
	"log"

	"github.com/iotames/cdnguard"
)

func debug() {
	txt, err := gdb.GetDbSizeText()
	log.Printf("------debug---GetDbSizeText(%s)---\n", txt)
	if err != nil {
		panic(err)
	}
	list, err := gdb.GetIpWhiteList()
	if err != nil {
		panic(err)
	}
	log.Printf("------debug---GetIpWhiteList(%v)---\n", list)
	ips, err := cdnguard.GetTopRequestIpToday(1, 5)
	if err != nil {
		panic(err)
	}
	for _, ip := range ips {
		log.Printf("---------RequestIpBlackList--ip( %s )--count(%d)-----\n", ip.Ip, ip.RequestCount)
	}
}
