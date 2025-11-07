package main

import (
	"log"
	"time"

	"github.com/iotames/cdnguard"
	"github.com/iotames/cdnguard/cdnapi"
)

func extCmdRun() bool {
	// TODO gdb.Stats() 的调用应该放在WebServer的API接口中，才是所要获取的连接池的实际信息
	if DbStats {
		gdb.Stats()
		return true
	}
	if StatisEveryDay {
		// 记录开始时间
		startTime := time.Now()
		log.Println("Begin Statis ......")
		if rownum, err := cdnguard.StatisEveryDay(); err != nil {
			panic(err)
		} else {
			costTime := time.Since(startTime) // 正确计算耗时
			log.Println("Statis every day AffectedRowNum:", rownum, "costTime:", costTime)
		}
		return true
	}
	if Prune {
		if err := gdb.Prune(); err != nil {
			panic(err)
		}
		return true
	}
	if SyncBucketFiles {
		capi := cdnapi.NewCdnApi(CdnName, QiniuAccessKey, QiniuSecretKey, BucketNameList)
		err := capi.SyncFilesInfo(BucketName)
		if err != nil {
			panic(err)
		}
		return true
	}
	if ShowBucketFiles {
		capi := cdnapi.NewCdnApi(CdnName, QiniuAccessKey, QiniuSecretKey, BucketNameList)
		err := capi.ShowFilesInfo(BucketName, LastCursorMarker)
		if err != nil {
			panic(err)
		}
		return true
	}
	if AddBlackIps {
		cdnguard.AddBlackIpList(RequestLimit)
		return true
	}
	if Debug {
		debug()
		return true
	}
	return false
}
