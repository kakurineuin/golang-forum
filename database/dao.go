package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // for mysql。
)

// DB 資料庫物件。
var db *gorm.DB

type (
	DBFunc func(tx *gorm.DB) error // 以 gorm tx 物件為傳入參數的函式，此函式內自行實作 CRUD 等操作。
	DAO    struct {
		DB *gorm.DB
	}
)

// InitDAO 初始化資料庫存取物件。
func InitDAO(user string, password string, dbname string) *DAO {
	var err error
	dbSource := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		user, password, dbname)
	db, err = gorm.Open("mysql", dbSource)

	if err != nil {
		panic(err)
	}

	db.SingularTable(true)
	return &DAO{
		DB: db,
	}
}

// WithinTransaction 會開啟交易來執行傳入的 DBFunc。當執行新增、修改、刪除時就使用此函式。
func (dao DAO) WithinTransaction(fn DBFunc) (err error) {
	tx := dao.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err = tx.Error; err != nil {
		return err
	}

	if err = fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
