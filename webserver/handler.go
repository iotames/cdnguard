package webserver

import (
	"net/http"

	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
)

func setRouter(s *httpsvr.EasyServer) {
	s.AddHandler("GET", "/cdnauth", cdnauth)
	s.AddHandler("GET", "/api/local/dbstats", dbstats)
}

func success(ctx httpsvr.Context) {
	ctx.Writer.Write(response.NewApiDataOk("success").Bytes())
}

func fail(ctx httpsvr.Context) {
	ctx.Writer.WriteHeader(http.StatusUnauthorized)
	ctx.Writer.Write(response.NewApiDataUnauthorized().Bytes())
}
