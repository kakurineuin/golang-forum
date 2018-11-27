package post

import (
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// Handler 處理請求的 handler。
type Handler struct {
	Service *Service
}

// FindTopicsStatistics 查詢主題統計資料。
func (h Handler) FindTopicsStatistics(c echo.Context) (err error) {
	golangStatistics, nodeJSStatistics, err := h.Service.FindTopicsStatistics()

	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"golang": golangStatistics,
		"nodeJS": nodeJSStatistics,
	})
}

// FindTopics 查詢主題列表。
func (h Handler) FindTopics(c echo.Context) (err error) {
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
	topics, totalCount, err := h.Service.FindTopics(category, searchTopic, offset, limit)

	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"topics":     topics,
		"totalCount": totalCount,
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

	err = h.Service.CreatePost(c.Param("category"), post)

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
func (h Handler) FindTopic(c echo.Context) (err error) {
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

	findPostsResults, totalCount, err := h.Service.FindTopic(category, id, offset, limit)

	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"posts":      findPostsResults,
		"totalCount": totalCount,
	})
}

// UpdatePost 修改文章。
func (h Handler) UpdatePost(c echo.Context) (err error) {
	postOnUpdate := new(PostOnUpdate)

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

	post, err := h.Service.UpdatePost(c.Param("category"), id, *postOnUpdate, getUserID(c))

	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "修改文章成功。",
		"post":    post,
	})
}

func getUserID(c echo.Context) int {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return int(claims["id"].(float64))
}
