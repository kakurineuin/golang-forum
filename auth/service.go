package auth

import (
	fe "github.com/kakurineuin/golang-forum/error"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// roleUser 表示角色是一般使用者。
const roleUser string = "user"

// Service 處理請求的 service。
type Service struct {
	DB *gorm.DB
}

// Register 註冊。
func (s Service) Register(userProfile *UserProfile) (err error) {

	// 檢查是否已有相同使用者名稱。
	count := 0

	if err = s.DB.Table("user_profile").
		Where("username = ?", userProfile.Username).
		Count(&count).Error; err != nil {
		return
	}

	if count > 0 {
		return fe.CustomError{http.StatusBadRequest, "此使用者名稱已被使用。"}
	}

	if err = s.DB.Table("user_profile").
		Where("email = ?", userProfile.Email).
		Count(&count).Error; err != nil {
		return
	}

	if count > 0 {
		return fe.CustomError{http.StatusBadRequest, "此 email 已被使用。"}
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

	return s.DB.Create(userProfile).Error
}

// Login 登入。
func (s Service) Login(loginRequest LoginRequest) (userProfile UserProfile, err error) {

	// 查詢帳號。
	if err = s.DB.Where("email = ?", loginRequest.Email).First(&userProfile).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return UserProfile{}, fe.CustomError{http.StatusNotFound, "查無此 email 帳號。"}
		}

		return UserProfile{}, err
	}

	// 核對密碼。
	err = bcrypt.CompareHashAndPassword([]byte(*userProfile.Password), []byte(*loginRequest.Password))
	if err != nil {
		return UserProfile{}, fe.CustomError{http.StatusBadRequest, "密碼錯誤。"}
	}

	return
}
