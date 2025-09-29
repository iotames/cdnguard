package webserver

import (
	"github.com/iotames/cdnguard/db"
	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
)

var router = map[string]func(ctx httpsvr.Context){
	"/cdnauth": cdnauth,
}

func GetDB() *db.DB {
	return db.GetDb(nil)
}

func success(ctx httpsvr.Context) {
	ctx.Writer.Write(response.NewApiDataOk("success").Bytes())
}

func fail(ctx httpsvr.Context) {
	// status=404
	ctx.Writer.Write(response.NewApiDataUnauthorized().Bytes())
}
