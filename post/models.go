package post

import (
	"time"
)

// Post 文章。
type Post struct {
	ID            int
	UserProfileID int    `json:"user_profile_id"`
	ReplyPostID   int    `json:"reply_post_id"`
	Title         string `json:"title" valid:"required~標題必填。,stringlength(1|30)~標題長度須在1至30之間。"`
	Content       string `json:"content" valid:"required~內文必填。,stringlength(1|500)~內文長度須在1至500之間。"`
	IsTopic       int    `json:"is_topic"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
