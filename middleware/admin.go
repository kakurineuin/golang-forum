package middleware

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// Admin 限制只有系統管理員才能使用功能的 middleware。
func Admin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)

		// 若未登入。
		if user == nil {
			c.JSON(http.StatusForbidden, map[string]interface{}{
				"message": "權限不足。",
			})
			return nil
		}

		claims := user.Claims.(jwt.MapClaims)
		role := claims["role"].(string)

		// 若不是系統管理員。
		if role != "admin" {
			c.JSON(http.StatusForbidden, map[string]interface{}{
				"message": "權限不足。",
			})
			return nil
		}

		return next(c)
	}
}
