package webserver

import (
	"github.com/iotames/cdnguard/db"
	"github.com/iotames/easyserver/httpsvr"
)

var router = map[string]func(ctx httpsvr.Context){
	"/cdnauth": cdnauth,
}

func GetDB() *db.DB {
	return db.GetDb(nil)
}
