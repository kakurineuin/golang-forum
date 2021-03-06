package model

import (
	"time"
)

// Statistics 首頁的各類型文章統計資料。
// 其中 TopicCount 和 ReplyCount 為了在查無資料時是 0，所以使用 int 而不是 *int。
type Statistics struct {
	TopicCount       int        `json:"topicCount"`
	ReplyCount       int        `json:"replyCount"`
	LastPostUsername *string    `json:"lastPostUsername"`
	LastPostTime     *time.Time `json:"lastPostTime"`
}

// Topic 主題資料，查詢主題列表時所用。
type Topic struct {
	Id                 *int       `json:"id"`
	Topic              *string    `json:"topic"`
	ReplyCount         *int       `json:"replyCount"`
	CreatedAt          *time.Time `json:"createdAt"`
	Username           *string    `json:"username"`
	LastReplyCreatedAt *time.Time `json:"lastReplyCreatedAt"`
	LastReplyUsername  *string    `json:"lastReplyUsername"`
}

// Post 文章。
type Post struct {
	Id            *int       `json:"id"`
	UserProfileId *int       `json:"userProfileId"`
	ReplyPostId   *int       `json:"replyPostId"`
	Topic         *string    `json:"topic" valid:"required~主題必填。,stringlength(1|30)~主題長度須在1至30之間。"`
	Content       *string    `json:"content" valid:"required~內文必填。,stringlength(1|20000)~內文長度須在1至20000之間。"`
	CreatedAt     *time.Time `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt"`
	DeletedAt     *time.Time `json:"deletedAt"`
}

// FindPostsResult 查詢某個主題討論串的結果。
type FindPostsResult struct {
	Id        *int       `json:"id"`
	Topic     *string    `json:"topic"`
	Content   *string    `json:"content"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	Username  *string    `json:"username"`
	Role      *string    `json:"role"`
}

// PostOnUpdate 修改的文章。
type PostOnUpdate struct {
	Content *string `json:"content" valid:"required~內文必填。,stringlength(1|20000)~內文長度須在1至20000之間。"`
}
