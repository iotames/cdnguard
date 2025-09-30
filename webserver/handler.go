package webserver

import (
	// "once"
	"net/http"

	"github.com/iotames/cdnguard/db"
	"github.com/iotames/easyserver/httpsvr"
	"github.com/iotames/easyserver/response"
)

// var once sync.Once
// var ipWhiteList []string

// GetIpWhiteList 获取IP白名单
// TODO 可以N小时更新一次IP白名单
func GetIpWhiteList() []string {
	// once.Do(func() {
	ipWhiteList, _ := GetDB().GetIpWhiteList()
	// })
	return ipWhiteList
}

// GetIpBlackList 获取IP黑名单
// TODO 可以15分钟更新一次IP黑名单，可以不用每次网络请求都执行SQL查询来获取IP黑名单列表
func GetIpBlackList() []string {
	ipBlackList, _ := GetDB().GetIpBlackList()
	return ipBlackList
}

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
	ctx.Writer.WriteHeader(http.StatusUnauthorized)
	ctx.Writer.Write(response.NewApiDataUnauthorized().Bytes())
}
