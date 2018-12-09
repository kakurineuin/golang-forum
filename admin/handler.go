package admin

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// Handler 處理請求的 handler。
type Handler struct {
	Service *Service
}

// FindUsers 查詢使用者。
func (h Handler) FindUsers(c echo.Context) (err error) {
	searchUser := c.QueryParam("searchUser")
	offset, err := strconv.Atoi(c.QueryParam("offset"))

	if err != nil {
		return
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))

	if err != nil {
		return
	}

	c.Logger().Infof("searchUser: %v, offset: %v, limit: %v", searchUser, offset, limit)
	users, totalCount, err := h.Service.FindUsers(searchUser, offset, limit)

	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"users":      users,
		"totalCount": totalCount,
	})
}

// DisableUser 停用使用者。
func (h Handler) DisableUser(c echo.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return
	}

	user, err := h.Service.DisableUser(id)

	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "停用使用者成功。",
		"user":    user,
	})
}
