package service

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/kakurineuin/golang-forum/database"
	fe "github.com/kakurineuin/golang-forum/error"
	"github.com/kakurineuin/golang-forum/model"
	"github.com/kakurineuin/golang-forum/sql"
	"net/http"
	"strings"
	"time"
)

// TopicService 處理請求的 service。
type TopicService struct {
	DAO *database.DAO
}

// FindTopicsStatistics 查詢主題統計資料。
func (s TopicService) FindTopicsStatistics() (golangStatistics, nodeJSStatistics model.Statistics, err error) {

	// 查詢 golang 文章統計資料。
	err = s.DAO.DB.Raw(sql.SqlTemplate["FindTopicsGolangStatistics"]).Scan(&golangStatistics).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return model.Statistics{}, model.Statistics{}, err
	}

	// 查詢 Node.js 文章統計資料。
	err = s.DAO.DB.Raw(sql.SqlTemplate["FindTopicsNodeJSStatistics"]).Scan(&nodeJSStatistics).Error

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

	statement := fmt.Sprintf(sql.SqlTemplate["FindTopics"], table, table)
	rows, err := s.DAO.DB.Raw(statement, searchTopic, offset, limit).Rows()
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
	statement = fmt.Sprintf(sql.SqlTemplate["FindTopicsTotalCount"], table, table)
	row := s.DAO.DB.Raw(statement, searchTopic).Row()
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

	statement := fmt.Sprintf(sql.SqlTemplate["FindTopic"], table, table)
	rows, err := s.DAO.DB.Raw(statement, id, id, offset, limit).Rows()
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
	statement = fmt.Sprintf(sql.SqlTemplate["FindTopicTotalCount"], table, table)
	row := s.DAO.DB.Raw(statement, id, id).Row()
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
