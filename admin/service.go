package admin

import (
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/kakurineuin/golang-forum/database"
)

// Service 處理請求的 Service。
type Service struct {
	DAO *database.DAO
}

// FindUsers 查詢使用者。
func (s Service) FindUsers(searchUser string, offset, limit int) (users []User, totalCount int, err error) {
	users = make([]User, 0)
	searchUser = strings.TrimSpace(searchUser)
	DB := s.DAO.DB.Table("user_profile").Offset(offset).Limit(limit).Order("username asc")

	if searchUser != "" {
		DB = DB.Where("username like ?", "%"+searchUser+"%")
	}

	err = DB.Find(&users).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return users, 0, err
	}

	// 查詢總筆數。
	DB = s.DAO.DB.Table("user_profile")

	if searchUser != "" {
		DB = DB.Where("username like ?", "%"+searchUser+"%")
	}

	err = DB.Count(&totalCount).Error

	if err != nil {
		return users, 0, err
	}

	return users, totalCount, nil
}

// DisableUser 停用使用者。
func (s Service) DisableUser(id int) (user User, err error) {
	err = s.DAO.WithinTransaction(func(tx *gorm.DB) error {
		return tx.Table("user_profile").Where("id = ?", id).Update("is_disabled", 1).Error
	})

	if err != nil {
		return User{}, err
	}

	err = s.DAO.DB.Table("user_profile").First(&user, id).Error

	if err != nil {
		return User{}, err
	}

	return user, nil
}
