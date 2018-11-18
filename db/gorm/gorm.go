package gorm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // for mysql。
)

// DB 資料庫物件。
var DB *gorm.DB

// InitDB 初始化資料庫連線。
func InitDB(user string, password string, dbname string) {
	var err error
	dbSource := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		user, password, dbname)
	DB, err = gorm.Open("mysql", dbSource)

	if err != nil {
		panic(err)
	}

	DB.SingularTable(true)
}
