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
	category := c.Param("category")
	offset := c.QueryParam("offset")
	limit := c.QueryParam("limit")
	c.Logger().Infof("category: %v, offset: %v, limit: %v", category, offset, limit)

	posts := []Post{}
	err = h.DB.Table("post_" + category).
		Where("is_topic = 1").
		Order("id desc").
		Offset(offset).
		Limit(limit).
		Find(&posts).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"posts": &posts,
	})
}

// CreatePost 新增文章。
func (h Handler) CreatePost(c echo.Context) (err error) {
	post := new(Post)

	if err = c.Bind(post); err != nil {
		return
	}

	if err = c.Validate(post); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	category := c.Param("category")

	if err = h.DB.Table("post_" + category).Create(post).Error; err != nil {
		return
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "新增文章成功。",
		"post":    post,
	})
}
