package webserver

import (
	"fmt"

	"github.com/iotames/easyserver/httpsvr"
)

func NewWebServer(port int) *httpsvr.EasyServer {
	s := httpsvr.NewEasyServer(fmt.Sprintf(":%d", port))
	s.AddMiddleHead(httpsvr.NewMiddleCORS("*"))
	setRouter(s)
	return s
}
