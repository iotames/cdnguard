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

// 搬迁文件
// 搬迁规则 Copy
func (f FileMigrate) Handler(w http.ResponseWriter, r *http.Request, dataFlow *httpsvr.DataFlow) (next bool) {
	f.Copy(r)
	return true
}

// 搬迁文件
// 搬迁规则
func (f FileMigrate) Copy(r *http.Request) {
	requestReferer := r.Header.Get("referer")
	log.Debug("migrate copy", "fromBucket", f.fromBucket, "toBucket", f.toBucket, "host", r.Host, "hostname", r.URL.Hostname(), "path", r.URL.Path, "requestReferer", requestReferer, "migrateReferer", f.referer)

	// if f.referer == requestReferer {
	// 	// TODO 调用七牛云API复制文件
	// }
	// 添加日志 file_opt_log: file_key, opt_type(0copy,1move,3delete), status(1success|0fail), file_size, upload_time, created_at, updated_at, qiniu_etag, md5, from_bucket, to_bucket
	// 以后会从日志记录中删除原有的文件

}

// func Move() {}

type HttpLocation struct {
	fromHost, toHost string
}

func NewHttpLocation(fromHost, toHost string) *HttpLocation {
	log.Info("NewHttpLocation", "fromHost", fromHost, "toHost", toHost)
	return &HttpLocation{fromHost: fromHost, toHost: toHost}
}

// 已搬迁，访问搬迁后的bucket
// 未搬迁，跳转至原来的bucket
func (l HttpLocation) Handler(w http.ResponseWriter, r *http.Request, dataFlow *httpsvr.DataFlow) (next bool) {
	// HTTP 302 临时重定向。注意不要使用HTTP 301 永久重定向
	// HTTP/1.1 302 Found
	// Location: https://xx.yyy.com/b2.jpg

	// HTTP 301 永久重定向
	// HTTP/1.1 301 Moved Permanently
	// Location: https://xx.yyy.com/b2.jpg
	log.Debug("HttpLocation", "fromHost", l.fromHost, "toHost", l.toHost, "host", r.Host, "hostname", r.URL.Hostname(), "path", r.URL.Path)
	if r.Host == l.fromHost {
		// 搬迁后的bucket
		w.Header().Set("Location", "https://"+l.toHost+r.URL.Path)
		w.WriteHeader(302)
	}
	return true
}
