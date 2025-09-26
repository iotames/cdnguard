package webserver

import (
	"fmt"

	"github.com/iotames/easyserver/httpsvr"
	_ "github.com/lib/pq"
)

func NewWebServer(port int) *httpsvr.EasyServer {
	s := httpsvr.NewEasyServer(fmt.Sprintf(":%d", port))
	s.AddMiddleware(httpsvr.NewMiddleCORS("*"))
	for k, v := range router {
		s.AddHandler("GET", k, v)
	}
	return s
}
