package auth

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// roleUser 表示角色是一般使用者。
const roleUser string = "user"

// Handler 處理請求的 handler。
type Handler struct {
	DB *gorm.DB
}

// Register 註冊。
func (h Handler) Register(c echo.Context) (err error) {
	userProfile := new(UserProfile)

	if err = c.Bind(userProfile); err != nil {
		return
	}

	if err = c.Validate(userProfile); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	// 檢查是否已有相同帳號。
	count := 0

	if err = h.DB.Table("user_profile").
		Where("account = ?", userProfile.Account).
		Count(&count).Error; err != nil {
		return
	}

	if count > 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "此帳號已被使用。",
		})
	}

	if err = h.DB.Table("user_profile").
		Where("email = ?", userProfile.Email).
		Count(&count).Error; err != nil {
		return
	}

	if count > 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "此 email 已被使用。",
		})
	}

	// 核對密碼。
	hash, err := bcrypt.GenerateFromPassword([]byte(userProfile.Password), bcrypt.DefaultCost)

	if err != nil {
		return
	}

	userProfile.Password = string(hash)
	userProfile.Role = roleUser

	if err = h.DB.Create(userProfile).Error; err != nil {
		return
	}

	userProfile.Password = "" // 密碼不能傳到前端。
	return returnTokenAndUserProfile(c, *userProfile, "註冊成功。")
}

// Login 登入。
func (h Handler) Login(c echo.Context) (err error) {
	loginRequest := new(LoginRequest)

	if err = c.Bind(loginRequest); err != nil {
		return
	}

	if err := c.Validate(loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	var userProfile UserProfile

	// 查詢帳號。
	if err = h.DB.Where("email = ?", loginRequest.Email).First(&userProfile).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"message": "查無此 email 帳號。",
			})
		}

		return
	}

	// 核對密碼。
	err = bcrypt.CompareHashAndPassword([]byte(userProfile.Password), []byte(loginRequest.Password))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "密碼錯誤。",
		})
	}

	userProfile.Password = "" // 密碼不能傳到前端。
	return returnTokenAndUserProfile(c, userProfile, "登入成功。")
}

// Logout 登出。
func (h Handler) Logout(c echo.Context) error {
	// TODO: 待實做。
	return nil
}

func (h Handler) Restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	c.Logger().Info(user)
	claims := user.Claims.(jwt.MapClaims)
	c.Logger().Info(claims)
	account := claims["account"].(string)
	c.Logger().Info(account)
	return c.String(http.StatusOK, "Welcome "+account+"!")
}

func createToken(userProfile UserProfile) (string, int64, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	exp := time.Now().Add(time.Minute * 10).Unix()

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["account"] = userProfile.Account
	claims["email"] = userProfile.Email
	claims["exp"] = exp

	// Generate encoded token.
	tokenString, err := token.SignedString([]byte("golang_secret"))
	return tokenString, exp, err
}

func returnTokenAndUserProfile(
	c echo.Context, userProfile UserProfile, message string) (err error) {
	token, exp, err := createToken(userProfile)

	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     message,
		"userProfile": userProfile,
		"token":       token,
		"exp":         exp,
	})
}
