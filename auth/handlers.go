package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

// JwtSecret JWT secret key。
const JwtSecret = "die_meere"

// Handler 處理請求的 handler。
type Handler struct {
	Service *Service
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

	if err = h.Service.Register(userProfile); err != nil {
		return
	}

	userProfile.Password = nil // 密碼不能傳到前端。
	return returnTokenAndUserProfile(c, *userProfile, "註冊成功。")
}

// Login 登入。
func (h Handler) Login(c echo.Context) (err error) {
	loginRequest := new(LoginRequest)

	if err = c.Bind(loginRequest); err != nil {
		return
	}

	if err = c.Validate(loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	userProfile, err := h.Service.Login(*loginRequest)

	if err != nil {
		return
	}

	userProfile.Password = nil // 密碼不能傳到前端。
	return returnTokenAndUserProfile(c, userProfile, "登入成功。")
}

func createToken(userProfile UserProfile) (string, int64, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	exp := time.Now().Add(time.Hour * 72).Unix()

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userProfile.ID
	claims["username"] = userProfile.Username
	claims["email"] = userProfile.Email
	claims["exp"] = exp

	// Generate encoded token.
	tokenString, err := token.SignedString([]byte(JwtSecret))
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
