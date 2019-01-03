package service

import (
	"errors"
	"fmt"
	"github.com/kakurineuin/golang-forum/model"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kakurineuin/golang-forum/database"

	"github.com/beevik/etree"
	"github.com/jinzhu/gorm"
	fe "github.com/kakurineuin/golang-forum/error"
)

var sqlTemplate = make(map[string]string)

func init() {
	pwd, _ := os.Getwd()
	directory := filepath.Base(pwd)
	sqlTemplatePath := ""

	switch directory {
	case "golang-forum":
		sqlTemplatePath = "sql/template.xml"
	default:
		// 執行測試時的路徑。
		sqlTemplatePath = "../../sql/template.xml"
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

// TopicService 處理請求的 service。
type TopicService struct {
	DAO *database.DAO
}

// FindForumStatistics 查詢論壇統計資料。
func (s TopicService) FindForumStatistics() (forumStatistics model.ForumStatistics, err error) {
	err = s.DAO.DB.Raw(sqlTemplate["FindForumStatistics"]).Scan(&forumStatistics).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return model.ForumStatistics{}, err
	}

	return forumStatistics, nil
}

// FindTopicsStatistics 查詢主題統計資料。
func (s TopicService) FindTopicsStatistics() (golangStatistics, nodeJSStatistics model.Statistics, err error) {

	// 查詢 golang 文章統計資料。
	err = s.DAO.DB.Raw(sqlTemplate["FindTopicsGolangStatistics"]).Scan(&golangStatistics).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return model.Statistics{}, model.Statistics{}, err
	}

	// 查詢 Node.js 文章統計資料。
	err = s.DAO.DB.Raw(sqlTemplate["FindTopicsNodeJSStatistics"]).Scan(&nodeJSStatistics).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return model.Statistics{}, model.Statistics{}, err
	}

	return golangStatistics, nodeJSStatistics, nil
}

// FindTopics 查詢主題列表。
func (s TopicService) FindTopics(category, searchTopic string, offset, limit int) (topics []model.Topic, totalCount int, err error) {
	topics = make([]model.Topic, 0)
	searchTopic = "%" + strings.TrimSpace(searchTopic) + "%"

	table, err := getTable(category)

	if err != nil {
		return topics, 0, err
	}

	sql := fmt.Sprintf(sqlTemplate["FindTopics"], table, table)
	rows, err := s.DAO.DB.Raw(sql, searchTopic, offset, limit).Rows()
	defer rows.Close()

	if err != nil {
		return topics, 0, err
	}

	for rows.Next() {
		var topic model.Topic
		s.DAO.DB.ScanRows(rows, &topic)
		topics = append(topics, topic)
	}

	// 查詢總筆數。
	sql = fmt.Sprintf(sqlTemplate["FindTopicsTotalCount"], table, table)
	row := s.DAO.DB.Raw(sql, searchTopic).Row()
	row.Scan(&totalCount)
	return
}

// CreatePost 新增文章。
func (s TopicService) CreatePost(category string, post *model.Post) (err error) {
	return s.DAO.WithinTransaction(func(tx *gorm.DB) error {
		return tx.Table("post_" + category).Create(post).Error
	})
}

// FindTopic 查詢某個主題的討論文章。
func (s TopicService) FindTopic(category string, id, offset, limit int) (findPostsResults []model.FindPostsResult, totalCount int, err error) {
	table, err := getTable(category)

	if err != nil {
		return nil, 0, err
	}

	sql := fmt.Sprintf(sqlTemplate["FindTopic"], table, table)
	rows, err := s.DAO.DB.Raw(sql, id, id, offset, limit).Rows()
	defer rows.Close()

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, 0, err
	}

	for rows.Next() {
		var findPostsResult model.FindPostsResult
		s.DAO.DB.ScanRows(rows, &findPostsResult)
		findPostsResults = append(findPostsResults, findPostsResult)
	}

	// 查詢總筆數。
	sql = fmt.Sprintf(sqlTemplate["FindTopicTotalCount"], table, table)
	row := s.DAO.DB.Raw(sql, id, id).Row()
	row.Scan(&totalCount)
	return
}

// UpdatePost 修改文章。
func (s TopicService) UpdatePost(category string, id int, postOnUpdate model.PostOnUpdate, userID int) (post model.Post, err error) {

	// 查詢原本文章。
	err = s.DAO.DB.Table("post_"+category).First(&post, id).Error

	if err != nil {
		return model.Post{}, err
	}

	// 不能修改已刪除的文章。
	if post.DeletedAt != nil {
		return model.Post{}, fe.CustomError{
			HTTPStatusCode: http.StatusBadRequest,
			Message:        "不能修改已刪除的文章。",
		}
	}

	// 不能修改別人的文章。
	if *post.UserProfileID != userID {
		return model.Post{}, fe.CustomError{
			HTTPStatusCode: http.StatusBadRequest,
			Message:        "不能修改別人的文章。",
		}
	}

	// 修改文章。
	post.Content = postOnUpdate.Content
	err = s.DAO.WithinTransaction(func(tx *gorm.DB) error {
		return tx.Table("post_" + category).Save(&post).Error
	})

	if err != nil {
		return model.Post{}, err
	}

	return
}

// DeletePost 刪除文章，不是真的刪除，而是修改文章內容和刪除時間欄位。
func (s TopicService) DeletePost(category string, id, userID int) (post model.Post, err error) {

	// 查詢原本文章。
	err = s.DAO.DB.Table("post_"+category).First(&post, id).Error

	if err != nil {
		return model.Post{}, err
	}

	user := model.UserProfile{}
	err = s.DAO.DB.First(&user, userID).Error

	if err != nil {
		return model.Post{}, err
	}

	// 不是系統管理員則不能刪除別人的文章。
	if *user.Role != "admin" && *post.UserProfileID != userID {
		return model.Post{}, fe.CustomError{
			HTTPStatusCode: http.StatusBadRequest,
			Message:        "不能刪除別人的文章。",
		}
	}

	// 不是真的刪除，而是修改文章內容並更新刪除時間欄位。
	content := "此篇文章已被刪除。"
	post.Content = &content
	deleteAt := time.Now()
	post.DeletedAt = &deleteAt
	err = s.DAO.WithinTransaction(func(tx *gorm.DB) error {
		return tx.Table("post_" + category).Save(&post).Error
	})

	if err != nil {
		return model.Post{}, err
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
