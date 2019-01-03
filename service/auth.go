package service

import (
	"github.com/kakurineuin/golang-forum/model"
	"net/http"

	"github.com/kakurineuin/golang-forum/database"

	"github.com/jinzhu/gorm"
	fe "github.com/kakurineuin/golang-forum/error"
	"golang.org/x/crypto/bcrypt"
)

// roleUser 表示角色是一般使用者。
const roleUser string = "user"

// AuthService 處理請求的 service。
type AuthService struct {
	DAO *database.DAO
}

// Register 註冊。
func (s AuthService) Register(userProfile *model.UserProfile) (err error) {

	// 檢查是否已有相同使用者名稱。
	count := 0

	if err = s.DAO.DB.Table("user_profile").
		Where("username = ?", userProfile.Username).
		Count(&count).Error; err != nil {
		return
	}

	if count > 0 {
		return fe.CustomError{
			HTTPStatusCode: http.StatusBadRequest,
			Message:        "此使用者名稱已被使用。",
		}
	}

	if err = s.DAO.DB.Table("user_profile").
		Where("email = ?", userProfile.Email).
		Count(&count).Error; err != nil {
		return
	}

	if count > 0 {
		return fe.CustomError{
			HTTPStatusCode: http.StatusBadRequest,
			Message:        "此 email 已被使用。",
		}
	}

	// 加密密碼。
	hash, err := bcrypt.GenerateFromPassword([]byte(*userProfile.Password), bcrypt.DefaultCost)

	if err != nil {
		return
	}

	hashString := string(hash)
	userProfile.Password = &hashString
	role := roleUser
	userProfile.Role = &role
	isDisabled := 0
	userProfile.IsDisabled = &isDisabled

	return s.DAO.WithinTransaction(func(tx *gorm.DB) error {
		return tx.Create(userProfile).Error
	})
}

// Login 登入。
func (s AuthService) Login(loginRequest model.LoginRequest) (userProfile model.UserProfile, err error) {

	// 檢查帳號是否存在。
	if err = s.DAO.DB.Where("email = ?", loginRequest.Email).First(&userProfile).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return model.UserProfile{}, fe.CustomError{
				HTTPStatusCode: http.StatusNotFound,
				Message:        "查無此 email 帳號。",
			}
		}

		return model.UserProfile{}, err
	}

	// 檢查帳號是否已被停用。
	if *userProfile.IsDisabled == 1 {
		return model.UserProfile{}, fe.CustomError{
			HTTPStatusCode: http.StatusForbidden,
			Message:        "此帳號已被停用。",
		}
	}

	// 核對密碼。
	err = bcrypt.CompareHashAndPassword([]byte(*userProfile.Password), []byte(*loginRequest.Password))
	if err != nil {
		return model.UserProfile{}, fe.CustomError{
			HTTPStatusCode: http.StatusBadRequest,
			Message:        "密碼錯誤。",
		}
	}

	return
}
