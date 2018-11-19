package post

import (
	"time"
)

// Statistics 首頁的文章統計資料。
// 其中 TopicCount 和 ReplyCount 為了在查無資料時是 0，所以使用 int 而不是 *int。
type Statistics struct {
	TopicCount      int        `json:"topicCount"`
	ReplyCount      int        `json:"replyCount"`
	LastPostAccount *string    `json:"lastPostAccount"`
	LastPostTime    *time.Time `json:"lastPostTime"`
}

// Topic 主題資料，查詢主題列表時所用。
type Topic struct {
	ID                 *int       `json:"id"`
	Topic              *string    `json:"topic"`
	ReplyCount         *int       `json:"replyCount"`
	CreatedAt          *time.Time `json:"createdAt"`
	Account            *string    `json:"account"`
	LastReplyCreatedAt *time.Time `json:"lastReplyCreatedAt"`
	LastReplyAccount   *string    `json:"lastReplyAccount"`
}

// Post 文章。
type Post struct {
	ID            *int       `json:"id"`
	UserProfileID *int       `json:"userProfileID"`
	ReplyPostID   *int       `json:"replyPostID"`
	Topic         *string    `json:"topic" valid:"required~主題必填。,stringlength(1|30)~主題長度須在1至30之間。"`
	Content       *string    `json:"content" valid:"required~內文必填。,stringlength(1|500)~內文長度須在1至500之間。"`
	CreatedAt     *time.Time `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt"`
}

// FindPostsResult 查詢某個主題討論串的結果。
type FindPostsResult struct {
	ID        *int       `json:"id"`
	Topic     *string    `json:"topic"`
	Content   *string    `json:"content"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	Account   *string    `json:"account"`
	Role      *string    `json:"role"`
}
