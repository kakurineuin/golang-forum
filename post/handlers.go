package post

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/beevik/etree"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

var sqlTemplate = make(map[string]string)

func init() {
	pwd, _ := os.Getwd()
	directory := filepath.Base(pwd)
	sqlTemplatePath := ""

	switch directory {
	case "post":
		sqlTemplatePath = "../sql/template.xml"
	case "forum":
		sqlTemplatePath = "sql/template.xml"
	default:
		fmt.Println("============== directory", directory)
	}

	doc := etree.NewDocument()

	if err := doc.ReadFromFile(sqlTemplatePath); err != nil {
		panic(err)
	}

	sqls := doc.SelectElement("Sqls")
	for _, sql := range sqls.SelectElements("Sql") {
		name := sql.SelectAttrValue("name", "")
		sqlTemplate[name] = sql.Text()
	}
}

// Handler 處理請求的 handler。
type Handler struct {
	DB *gorm.DB
}

// FindTopicsStatistics 查詢主題統計資料。
func (h Handler) FindTopicsStatistics(c echo.Context) (err error) {

	// 查詢 golang 文章統計資料。
	var golangStatistics Statistics

	h.DB.Raw(sqlTemplate["FindTopicsGolangStatistics"]).Scan(&golangStatistics)

	// 查詢 Node.js 文章統計資料。
	var nodeJSStatistics Statistics

	h.DB.Raw(sqlTemplate["FindTopicsNodeJSStatistics"]).Scan(&nodeJSStatistics)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"golang": golangStatistics,
		"nodeJS": nodeJSStatistics,
	})
}

// FindTopics 查詢主題列表。
func (h Handler) FindTopics(c echo.Context) (err error) {
	category := c.Param("category")
	offset := c.QueryParam("offset")
	limit := c.QueryParam("limit")
	c.Logger().Infof("category: %v, offset: %v, limit: %v", category, offset, limit)

	table, err := getTable(category)
	if err != nil {
		return
	}

	sql := fmt.Sprintf(sqlTemplate["FindTopics"], table, table)
	rows, err := h.DB.Raw(sql, offset, limit).Rows()
	defer rows.Close()

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return
	}

	topics := make([]Topic, 0)

	for rows.Next() {
		var topic Topic
		h.DB.ScanRows(rows, &topic)
		topics = append(topics, topic)
	}

	// 查詢總筆數。
	totalCount := 0
	sql = fmt.Sprintf(sqlTemplate["FindTopicsTotalCount"], table, table)
	row := h.DB.Raw(sql).Row()
	row.Scan(&totalCount)

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

	if err = h.DB.Table("post_" + c.Param("category")).Create(post).Error; err != nil {
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
	id := c.Param("id")
	offset := c.QueryParam("offset")
	limit := c.QueryParam("limit")
	c.Logger().Infof("category: %v, id: %v, offset: %v, limit: %v", category, id, offset, limit)

	table, err := getTable(category)

	if err != nil {
		return
	}

	sql := fmt.Sprintf(sqlTemplate["FindTopic"], table, table)
	rows, err := h.DB.Raw(sql, id, id, offset, limit).Rows()
	defer rows.Close()

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return
	}

	findPostsResults := make([]FindPostsResult, 0)

	for rows.Next() {
		var findPostsResult FindPostsResult
		h.DB.ScanRows(rows, &findPostsResult)
		findPostsResults = append(findPostsResults, findPostsResult)
	}

	// 查詢總筆數。
	totalCount := 0
	sql = fmt.Sprintf(sqlTemplate["FindTopicTotalCount"], table, table)
	row := h.DB.Raw(sql, id, id).Row()
	row.Scan(&totalCount)

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

	// 查詢原本文章。
	post := new(Post)
	err = h.DB.Table("post_"+c.Param("category")).
		First(post, c.Param("id")).
		Error

	if err != nil {
		return
	}

	userID := getUserID(c)

	// 不能修改別人的文章。
	if *post.UserProfileID != userID {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "不能別人的文章。",
		})
	}

	// 修改文章。
	post.Content = postOnUpdate.Content
	err = h.DB.Table("post_" + c.Param("category")).Save(post).Error

	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "修改文章成功。",
		"post":    post,
	})
}

func getTable(category string) (string, error) {
	switch category {
	case "golang":
		return "post_golang", nil
	case "nodejs":
		return "post_nodejs", nil
	default:
		return "", errors.New("category is error")
	}
}

func getUserID(c echo.Context) int {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return int(claims["id"].(float64))
}
