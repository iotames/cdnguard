package webserver

import (
	"net/http"

	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
)

var router = map[string]func(ctx httpsvr.Context){
	"/cdnauth": cdnauth,
}

func success(ctx httpsvr.Context) {
	ctx.Writer.Write(response.NewApiDataOk("success").Bytes())
}

func fail(ctx httpsvr.Context) {
	ctx.Writer.WriteHeader(http.StatusUnauthorized)
	ctx.Writer.Write(response.NewApiDataUnauthorized().Bytes())
}
