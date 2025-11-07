package webserver

import (
	"log"
	"net"
	"net/http"

	"github.com/iotames/cdnguard/model"
	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
)

func dbstats(ctx httpsvr.Context) {
	remoteAddr := ctx.Request.RemoteAddr
	var host string
	// 尝试分割主机和端口
	h, port, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		// 如果分割失败，可能没有端口号，直接使用整个地址
		host = remoteAddr
	} else {
		host = h
	}
	log.Printf("---webserver(/api/local/dbstats)--remoteAddr(%s)--host(%s)--port(%s)--\n", remoteAddr, host, port)
	// 检查是否为本地地址（拒绝非本地访问）
	if host != "127.0.0.1" && host != "::1" {
		fail(ctx)
		return
	}
	b := response.NewApiData(model.DbStats(), "success", http.StatusOK).Bytes()
	ctx.Writer.Write(b)
}
