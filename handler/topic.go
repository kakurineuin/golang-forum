package handler

import (
	"github.com/kakurineuin/golang-forum/model"
	"github.com/kakurineuin/golang-forum/service"
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// TopicHandler 處理主題相關功能請求的 handler。
type TopicHandler struct {
	TopicService *service.TopicService
}

// FindTopicsStatistics 查詢主題統計資料。
func (h TopicHandler) FindTopicsStatistics(c echo.Context) (err error) {
	golangStatistics, nodeJSStatistics, err := h.TopicService.FindTopicsStatistics()

	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"golang": golangStatistics,
		"nodeJS": nodeJSStatistics,
	})
}

// FindTopics 查詢主題列表。
func (h TopicHandler) FindTopics(c echo.Context) (err error) {
	category := c.Param("category")
	searchTopic := c.QueryParam("searchTopic")
	offset, err := strconv.Atoi(c.QueryParam("offset"))

	if err != nil {
		return
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))

	if err != nil {
		return
	}

	c.Logger().Infof("category: %v, searchTopic: %v, offset: %v, limit: %v", category, searchTopic, offset, limit)
	topics, totalCount, err := h.TopicService.FindTopics(category, searchTopic, offset, limit)

	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"topics":     topics,
		"totalCount": totalCount,
	})
}

// CreatePost 新增文章。
func (h TopicHandler) CreatePost(c echo.Context) (err error) {
	post := new(model.Post)

	if err = c.Bind(post); err != nil {
		return
	}

	if err = c.Validate(post); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	err = h.TopicService.CreatePost(c.Param("category"), post)

	if err != nil {
		return
	}

	message := ""

	if post.ReplyPostID == nil {
		message = "新增主題成功。"
	} else {
		message = "回覆成功。"
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": message,
		"post":    post,
	})
}

// FindTopic 查詢某個主題的討論文章。
func (h TopicHandler) FindTopic(c echo.Context) (err error) {
	category := c.Param("category")
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))

	if err != nil {
		return
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))

	if err != nil {
		return
	}

	c.Logger().Infof("category: %v, id: %v, offset: %v, limit: %v", category, id, offset, limit)
	findPostsResults, totalCount, err := h.TopicService.FindTopic(category, id, offset, limit)

	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"posts":      findPostsResults,
		"totalCount": totalCount,
	})
}

// UpdatePost 修改文章。
func (h TopicHandler) UpdatePost(c echo.Context) (err error) {
	postOnUpdate := new(model.PostOnUpdate)

	if err = c.Bind(postOnUpdate); err != nil {
		return
	}

	if err = c.Validate(postOnUpdate); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return
	}

	post, err := h.TopicService.UpdatePost(c.Param("category"), id, *postOnUpdate, getUserID(c))

	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "修改文章成功。",
		"post":    post,
	})
}

// DeletePost 刪除文章。
func (h TopicHandler) DeletePost(c echo.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return
	}

	post, err := h.TopicService.DeletePost(c.Param("category"), id, getUserID(c))

	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "刪除文章成功。",
		"post":    post,
	})
}

func getUserID(c echo.Context) int {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return int(claims["id"].(float64))
}
