package model

import (
	"github.com/iotames/cdnguard/db"
)

// 只调用db包的方法或数据
// model层是业务逻辑层的调用对象。是业务层和数据层的桥梁

func getDB() *db.DB {
	return db.GetDb(nil)
}
