package handler

import (
	"github.com/kakurineuin/golang-forum/model"
	"github.com/kakurineuin/golang-forum/service"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// AuthHandler 處理 auth 相關功能請求的 handler。
type AuthHandler struct {
	AuthService *service.AuthService
	JwtSecret   string
}

// Register 註冊。
func (h AuthHandler) Register(c echo.Context) (err error) {
	userProfile := new(model.UserProfile)

	if err = c.Bind(userProfile); err != nil {
		return
	}

	if err = c.Validate(userProfile); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	if err = h.AuthService.Register(userProfile); err != nil {
		return
	}

	userProfile.Password = nil // 密碼不能傳到前端。
	return returnResponse(c, *userProfile, "註冊成功。", h.JwtSecret)
}

// Login 登入。
func (h AuthHandler) Login(c echo.Context) (err error) {
	loginRequest := new(model.LoginRequest)

	if err = c.Bind(loginRequest); err != nil {
		return
	}

	if err = c.Validate(loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	userProfile, err := h.AuthService.Login(*loginRequest)

	if err != nil {
		return
	}

	userProfile.Password = nil // 密碼不能傳到前端。
	return returnResponse(c, userProfile, "登入成功。", h.JwtSecret)
}

func createToken(userProfile model.UserProfile, jwtSecret string) (string, int64, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	exp := time.Now().Add(time.Hour * 8760).Unix()

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userProfile.Id
	claims["username"] = userProfile.Username
	claims["email"] = userProfile.Email
	claims["role"] = userProfile.Role
	claims["exp"] = exp

	// Generate encoded token.
	tokenString, err := token.SignedString([]byte(jwtSecret))
	return tokenString, exp, err
}

func returnResponse(
	c echo.Context, userProfile model.UserProfile, message, jwtSecret string) (err error) {
	token, exp, err := createToken(userProfile, jwtSecret)

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
