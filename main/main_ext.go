package main

import (
	"github.com/iotames/cdnguard"
	"github.com/iotames/cdnguard/cdnapi"
)

func extCmdRun() bool {
	if Prune {
		if err := gdb.Prune(); err != nil {
			panic(err)
		}
		return true
	}
	if SyncBucketFiles {
		capi := cdnapi.NewCdnApi(CdnName, QiniuAccessKey, QiniuSecretKey, BucketNameList)
		capi.SyncFiles(BucketName)
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
