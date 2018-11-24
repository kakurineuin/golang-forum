package post

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/beevik/etree"
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

// Service 處理請求的 Service。
type Service struct {
	DB *gorm.DB
}

// FindTopicsStatistics 查詢主題統計資料。
func (s Service) FindTopicsStatistics() (golangStatistics, nodeJSStatistics *Statistics, err error) {
	golangStatistics = new(Statistics)
	nodeJSStatistics = new(Statistics)

	// 查詢 golang 文章統計資料。
	err = s.DB.Raw(sqlTemplate["FindTopicsGolangStatistics"]).Scan(golangStatistics).Error

	if err != nil {
		return nil, nil, err
	}

	// 查詢 Node.js 文章統計資料。
	err = s.DB.Raw(sqlTemplate["FindTopicsNodeJSStatistics"]).Scan(nodeJSStatistics).Error

	if err != nil {
		return nil, nil, err
	}

	return
}

// FindTopics 查詢主題列表。
func (s Service) FindTopics(category string, offset, limit int) (topics []Topic, totalCount int, err error) {
	table, err := getTable(category)

	if err != nil {
		return nil, 0, err
	}

	sql := fmt.Sprintf(sqlTemplate["FindTopics"], table, table)
	rows, err := s.DB.Raw(sql, offset, limit).Rows()
	defer rows.Close()

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, 0, err
	}

	for rows.Next() {
		var topic Topic
		s.DB.ScanRows(rows, &topic)
		topics = append(topics, topic)
	}

	// 查詢總筆數。
	sql = fmt.Sprintf(sqlTemplate["FindTopicsTotalCount"], table, table)
	row := s.DB.Raw(sql).Row()
	row.Scan(&totalCount)
	return
}

// CreatePost 新增文章。
func (s Service) CreatePost(category string, post *Post) (err error) {
	return s.DB.Table("post_" + category).Create(post).Error
}

// FindTopic 查詢某個主題的討論文章。
func (s Service) FindTopic(category string, id, offset, limit int) (findPostsResults []FindPostsResult, totalCount int, err error) {
	table, err := getTable(category)

	if err != nil {
		return nil, 0, err
	}

	sql := fmt.Sprintf(sqlTemplate["FindTopic"], table, table)
	rows, err := s.DB.Raw(sql, id, id, offset, limit).Rows()
	defer rows.Close()

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, 0, err
	}

	for rows.Next() {
		var findPostsResult FindPostsResult
		s.DB.ScanRows(rows, &findPostsResult)
		findPostsResults = append(findPostsResults, findPostsResult)
	}

	// 查詢總筆數。
	sql = fmt.Sprintf(sqlTemplate["FindTopicTotalCount"], table, table)
	row := s.DB.Raw(sql, id, id).Row()
	row.Scan(&totalCount)
	return
}

// UpdatePost 修改文章。
func (s Service) UpdatePost(category string, id int, postOnUpdate PostOnUpdate, userID int) (post *Post, err error) {
	post = new(Post)

	// 查詢原本文章。
	err = s.DB.Table("post_" + category).First(post, id).Error

	if err != nil {
		return nil, err
	}

	// 不能修改別人的文章。
	if *post.UserProfileID != userID {
		err = echo.NewHTTPError(http.StatusBadRequest, "不能修改別人的文章。")
		return nil, err
	}

	// 修改文章。
	post.Content = postOnUpdate.Content
	err = s.DB.Table("post_" + category).Save(post).Error

	if err != nil {
		return nil, err
	}

	return
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
