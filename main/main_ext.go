package main

import (
	"github.com/iotames/cdnguard"
)

func extCmdRun() bool {
	if Prune {
		if err := gdb.Prune(); err != nil {
			panic(err)
		}
		return true
	}
	if AddBlackIps {
		cdnguard.AddBlackIpList()
		return true
	}
	if Debug {
		debug()
		return true
	}
	return false
}
