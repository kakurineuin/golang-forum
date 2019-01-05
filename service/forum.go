package service

import (
	"github.com/jinzhu/gorm"
	"github.com/kakurineuin/golang-forum/database"
	"github.com/kakurineuin/golang-forum/model"
	"github.com/kakurineuin/golang-forum/sql"
)

// ForumService 處理論壇相關功能請求的 service。
type ForumService struct {
	DAO *database.DAO
}

// FindForumStatistics 查詢論壇統計資料。
func (s ForumService) FindForumStatistics() (forumStatistics model.ForumStatistics, err error) {
	err = s.DAO.DB.Raw(sql.SqlTemplate["FindForumStatistics"]).Scan(&forumStatistics).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return model.ForumStatistics{}, err
	}

	return forumStatistics, nil
}
