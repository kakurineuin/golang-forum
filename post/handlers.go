package post

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"net/http"
)

// Handler 處理請求的 handler。
type Handler struct {
	DB *gorm.DB
}

// FindPosts 查詢文章。
func (h Handler) FindPosts(c echo.Context) (err error) {
	findPostsRequest := new(FindPostsRequest)

	if err = c.Bind(findPostsRequest); err != nil {
		return
	}

	if err = c.Validate(findPostsRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	posts := []Post{}
	err = h.DB.Table("post_" + findPostsRequest.Category).
		Where("is_topic = 1").
		Order("id desc").
		Offset(findPostsRequest.Offset).
		Limit(findPostsRequest.Limit).
		Find(&posts).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"posts": &posts,
	})
}
