package handler

import (
	"github.com/kakurineuin/golang-forum/service"
	"github.com/labstack/echo"
	"net/http"

	"time"
)

// ForumHandler 處理論壇相關功能請求的 handler。
type ForumHandler struct {
	ForumService *service.ForumService
}

// FindForumStatistics 查詢論壇統計資料。
func (h ForumHandler) FindForumStatistics(c echo.Context) (err error) {
	forumStatistics, err := h.ForumService.FindForumStatistics()

	if err != nil {
		return
	}

	time.Sleep(5 * time.Second)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"forumStatistics": forumStatistics,
	})
}
