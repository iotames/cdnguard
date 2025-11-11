package migrate

import (
	"net/http"

	"github.com/iotames/cdnguard/log"
	"github.com/iotames/easyserver/httpsvr"
)

type FileMigrate struct {
	fromDomain, toDomain, referer, fromBucket, toBucket string
}

func NewFileMigrate(fromDomain, toDomain, referer, fromBucket, toBucket string) *FileMigrate {
	log.Info("NewFileMigrate", "fromDomain", fromDomain, "toDomain", toDomain, "referer", referer, "fromBucket", fromBucket, "toBucket", toBucket)
	return &FileMigrate{fromDomain: fromDomain, toDomain: toDomain, referer: referer, fromBucket: fromBucket, toBucket: toBucket}
}

func (f FileMigrate) Handler(w http.ResponseWriter, r *http.Request, dataFlow *httpsvr.DataFlow) (next bool) {
	f.Copy(r)
	return true
}

func (f FileMigrate) Copy(r *http.Request) {
	log.Debug("migrate copy", "fromBucket", f.fromBucket, "toBucket", f.toBucket, "host", r.Host, "hostname", r.URL.Hostname(), "path", r.URL.Path)
	// TODO 调用七牛云API复制文件

}

// func Move() {}

// HTTP 302 临时重定向。注意不要使用HTTP 301 永久重定向
// HTTP/1.1 302 Found
// Location: https://xx.yyy.com/b2.jpg

// HTTP 301 永久重定向
// HTTP/1.1 301 Moved Permanently
// Location: https://xx.yyy.com/b2.jpg

// type Migrate struct {
// 	FromDomain, ToDomain string
// }

// func NewMigrate(fromDomain, toDomain string) *Migrate {
// 	return &Migrate{FromDomain: fromDomain, ToDomain: toDomain}
// }
