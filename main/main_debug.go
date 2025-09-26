package main

import (
	"log"
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
}
