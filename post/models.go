package post

import (
	"time"
)

// PostStatistic 首頁的文章統計資料。
// 其中 TopicCount 和 ReplyCount 為了在查無資料時是 0，所以使用 int 而不是 *int。
type PostStatistic struct {
	TopicCount      int        `json:"topicCount"`
	ReplyCount      int        `json:"replyCount"`
	LastPostAccount *string    `json:"lastPostAccount"`
	LastPostTime    *time.Time `json:"lastPostTime"`
}

// Post 文章。
type Post struct {
	ID            *int       `json:"id"`
	UserProfileID *int       `json:"userProfileID"`
	ReplyPostID   *int       `json:"replyPostID"`
	Title         *string    `json:"title" valid:"required~標題必填。,stringlength(1|30)~標題長度須在1至30之間。"`
	Content       *string    `json:"content" valid:"required~內文必填。,stringlength(1|500)~內文長度須在1至500之間。"`
	CreatedAt     *time.Time `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt"`
}
